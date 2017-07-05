package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/xtaci/chat/kafka"

	"github.com/Shopify/sarama"
	log "github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"golang.org/x/net/context"

	. "github.com/xtaci/chat/proto"

	cli "gopkg.in/urfave/cli.v2"
)

var (
	OK                   = &Chat_Nil{}
	ERROR_ALREADY_EXISTS = errors.New("id already exists")
	ERROR_NOT_EXISTS     = errors.New("id not exists")
)

type Consumer struct {
	offset   int64 // next message offset
	pushFunc func(msg *Chat_Message)
}

// endpoint definition
type EndPoint struct {
	retention   int
	StartOffset int64 // offset of the first message
	Inbox       []Chat_Message
	consumers   map[uint64]*Consumer
	chReady     chan struct{}
	die         chan struct{}
	mu          sync.Mutex
}

func newEndPoint(retention int) *EndPoint {
	ep := &EndPoint{}
	ep.retention = retention
	ep.chReady = make(chan struct{}, 1)
	ep.consumers = make(map[uint64]*Consumer)
	ep.StartOffset = 1
	ep.die = make(chan struct{})
	go ep.pushTask()
	return ep
}

// push a message to this endpoint
func (ep *EndPoint) push(msg *Chat_Message) {
	if len(ep.Inbox) > ep.retention {
		ep.Inbox = append(ep.Inbox[1:], *msg)
		ep.StartOffset++
	} else {
		ep.Inbox = append(ep.Inbox, *msg)
	}
	ep.notifyConsumers()
}

// closes this endpoint
func (ep *EndPoint) close() {
	close(ep.die)
}

func (ep *EndPoint) notifyConsumers() {
	select {
	case ep.chReady <- struct{}{}:
	default:
	}
}

func (ep *EndPoint) pushTask() {
	for {
		select {
		case <-ep.chReady:
			ep.mu.Lock()
			for _, consumer := range ep.consumers {
				idx := consumer.offset - ep.StartOffset
				if idx < 0 { // lag behind many
					idx = 0
				}
				for i := idx; i < int64(len(ep.Inbox)); i++ {
					ep.Inbox[i].Offset = i + ep.StartOffset
					consumer.pushFunc(&ep.Inbox[i])
				}
				consumer.offset = ep.StartOffset + int64(len(ep.Inbox))
			}
			ep.mu.Unlock()
		case <-ep.die:
		}
	}
}

// server definition
type server struct {
	consumerid_autoinc uint64
	kafkaOffset        int64
	offsetBucket       string
	retention          int
	boltdb             string
	bucket             string
	interval           time.Duration
	eps                map[uint64]*EndPoint // end-point-s
	sync.RWMutex
}

func (s *server) init(c *cli.Context) {
	s.retention = c.Int("retention")
	s.boltdb = c.String("boltdb")
	s.bucket = c.String("bucket")
	s.interval = c.Duration("write-interval")
	s.offsetBucket = c.String("kafka-bucket")

	s.eps = make(map[uint64]*EndPoint)
	s.restore()
	go s.receive()
	go s.persistence_task()
}

func (s *server) read_ep(id uint64) *EndPoint {
	s.RLock()
	defer s.RUnlock()
	return s.eps[id]
}

// subscribe to an endpoint & receive server streams
func (s *server) Subscribe(p *Chat_Consumer, stream ChatService_SubscribeServer) error {
	ep := s.read_ep(p.Id)
	if ep == nil {
		log.Errorf("cannot find endpoint %v", p)
		return ERROR_NOT_EXISTS
	}

	consumerid := atomic.AddUint64(&s.consumerid_autoinc, 1)
	e := make(chan error, 1)

	// activate consumer
	ep.mu.Lock()

	// from newest
	if p.From == -1 {
		p.From = ep.StartOffset + int64(len(ep.Inbox))
	}
	ep.consumers[consumerid] = &Consumer{p.From, func(msg *Chat_Message) {
		if err := stream.Send(msg); err != nil {
			select {
			case e <- err:
			default:
			}
		}
	}}
	ep.mu.Unlock()
	defer func() {
		ep.mu.Lock()
		delete(ep.consumers, consumerid)
		ep.mu.Unlock()
	}()

	ep.notifyConsumers()

	select {
	case <-stream.Context().Done():
	case err := <-e:
		return err
	}
	return nil
}

func (s *server) receive() {
	consumer, err := kafka.NewConsumer()
	if err != nil {
		log.Fatalln(err)
	}

	defer consumer.Close()
	partionConsumer, err := consumer.ConsumePartition(kafka.ChatTopic, 0, s.kafkaOffset)
	log.Info("kafkaOffset ", s.kafkaOffset)
	if err != nil {
		log.Fatalln(err)
	}
	defer partionConsumer.Close()
	for {
		select {
		case msg := <-partionConsumer.Messages():
			log.WithField("msg", msg).WithField("OFFSET", s.kafkaOffset).WithField("IsNew", s.kafkaOffset < msg.Offset).Info("Receive")
			if s.kafkaOffset < msg.Offset {
				chat := new(Chat_Message)
				json.Unmarshal(msg.Key, &chat.Id)
				chat.Body = msg.Value
				ep := s.read_ep(chat.Id)
				s.Lock()
				if ep != nil {
					ep.mu.Lock()
					ep.push(chat)
					ep.mu.Unlock()
				}
				s.kafkaOffset = msg.Offset
				s.Unlock()
			}
		}
	}

}

func (s *server) Reg(ctx context.Context, p *Chat_Id) (*Chat_Nil, error) {
	s.Lock()
	defer s.Unlock()
	ep := s.eps[p.Id]
	if ep != nil {
		log.Errorf("id already exists:%v", p.Id)
		return nil, ERROR_ALREADY_EXISTS
	}

	s.eps[p.Id] = newEndPoint(s.retention)
	return OK, nil
}

func (s *server) Query(ctx context.Context, crange *Chat_ConsumeRange) (*Chat_List, error) {
	ep := s.read_ep(crange.Id)
	if ep == nil {
		return nil, ERROR_NOT_EXISTS
	}

	ep.mu.Lock()
	defer ep.mu.Unlock()

	if crange.From < ep.StartOffset {
		crange.From = ep.StartOffset
	}

	if crange.To > ep.StartOffset+int64(len(ep.Inbox))-1 {
		crange.To = ep.StartOffset + int64(len(ep.Inbox)) - 1
	}

	list := &Chat_List{}
	if crange.To > crange.From {
		return list, nil
	}

	for i := crange.From; i <= crange.To; i++ {
		msg := ep.Inbox[i-ep.StartOffset]
		msg.Offset = i
		list.Messages = append(list.Messages, &msg)
	}

	return list, nil
}

func (s *server) Latest(ctx context.Context, crange *Chat_ConsumeLatest) (*Chat_List, error) {
	ep := s.read_ep(crange.Id)
	if ep == nil {
		return nil, ERROR_NOT_EXISTS
	}

	ep.mu.Lock()
	defer ep.mu.Unlock()

	list := &Chat_List{}
	i := len(ep.Inbox) - int(crange.Length)
	if i < 0 {
		i = 0
	}
	for ; i < len(ep.Inbox); i++ {
		offset := int64(i) + ep.StartOffset
		msg := ep.Inbox[i]
		msg.Offset = offset
		list.Messages = append(list.Messages, &msg)
	}
	return list, nil
}

// persistence endpoints into db
func (s *server) persistence_task() {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	db := s.open_db()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	for {
		select {
		case <-ticker.C:
			s.dump(db)
		case nr := <-sig:
			s.dump(db)
			db.Close()
			log.Info(nr)
			os.Exit(0)
		}
	}
}

func (s *server) open_db() *bolt.DB {
	db, err := bolt.Open(s.boltdb, 0600, nil)
	if err != nil {
		log.Panic(err)
		os.Exit(-1)
	}
	// create bulket
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(s.bucket))
		if err != nil {
			log.Panicf("create bucket: %s", err)
			os.Exit(-1)
		}
		_, err = tx.CreateBucketIfNotExists([]byte(s.offsetBucket))
		if err != nil {
			log.Panicf("create bucket: %s", err)
			os.Exit(-1)
		}
		return nil
	})
	return db
}

func (s *server) dump(db *bolt.DB) {
	// save offset
	db.Update(func(tx *bolt.Tx) error {
		s.Lock()
		// write kafka offset
		b := tx.Bucket([]byte(s.offsetBucket))
		bin, _ := json.Marshal(s.kafkaOffset)
		if err := b.Put([]byte(s.offsetBucket), bin); err != nil {
			log.Error(err)
		}

		// write endpoints
		b = tx.Bucket([]byte(s.bucket))
		eps := make(map[uint64]*EndPoint)
		for k, v := range s.eps {
			eps[k] = v
		}

		for k, ep := range eps {
			ep.mu.Lock()
			if bin, err := json.Marshal(ep); err != nil {
				log.Error("cannot marshal:", err)
			} else if err := b.Put([]byte(fmt.Sprint(k)), bin); err != nil {
				log.Error(err)
			}
			ep.mu.Unlock()
		}
		s.Unlock()
		return nil
	})
}

func (s *server) restore() {
	// restore data from db file
	db := s.open_db()
	defer db.Close()
	count := 0
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.offsetBucket))
		s.kafkaOffset = sarama.OffsetNewest
		b.ForEach(func(k, v []byte) error {
			json.Unmarshal(v, &s.kafkaOffset)
			return nil
		})

		b = tx.Bucket([]byte(s.bucket))
		b.ForEach(func(k, v []byte) error {
			ep := newEndPoint(s.retention)
			ep.mu.Lock()
			if err := json.Unmarshal(v, &ep); err != nil {
				log.Fatalln("chat data corrupted:", err)
			}

			id, err := strconv.ParseUint(string(k), 0, 64)
			if err != nil {
				log.Fatalln("chat data corrupted:", err)
			}

			// settings
			if len(ep.Inbox) > s.retention {
				remove := len(ep.Inbox) - s.retention
				if remove > 0 {
					ep.Inbox = ep.Inbox[remove:]
				}
			}
			s.eps[id] = ep
			count++
			ep.mu.Unlock()
			return nil
		})
		return nil
	})

	log.Infof("restored %v chats", count)
}
