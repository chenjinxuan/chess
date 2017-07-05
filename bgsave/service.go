package main

import (
	pb "bgsave/proto"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/fzzy/radix/extra/cluster"
	"github.com/golang/snappy"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/vmihailenco/msgpack.v2"
)

const (
	SERVICE             = "[BGSAVE]"
	SAVE_DELAY          = 100 * time.Millisecond
	COUNT_DELAY         = 1 * time.Minute
	DEFAULT_REDIS_HOST  = "172.17.42.1:7100"
	DEFAULT_MONGODB_URL = "mongodb://172.17.42.1/mydb"
	ENV_REDIS_HOST      = "REDIS_HOST"
	ENV_MONGODB_URL     = "MONGODB_URL"
	ENV_SNAPPY          = "ENABLE_SNAPPY"
	BUFSIZ              = 65536
)

type server struct {
	wait          chan string
	redis_client  *cluster.Cluster
	db            *mgo.Database
	enable_snappy bool
}

func (s *server) init() {
	// snappy
	if env := os.Getenv(ENV_SNAPPY); env != "" {
		s.enable_snappy = true
	}

	// read redis host
	redis_host := DEFAULT_REDIS_HOST
	if env := os.Getenv(ENV_REDIS_HOST); env != "" {
		redis_host = env
	}
	// start connection to redis cluster
	client, err := cluster.NewCluster(redis_host)
	if err != nil {
		log.Panic(err)
		os.Exit(-1)
	}
	s.redis_client = client

	// read mongodb host
	mongodb_url := DEFAULT_MONGODB_URL
	if env := os.Getenv(ENV_MONGODB_URL); env != "" {
		mongodb_url = env
	}

	// start connection to mongodb
	sess, err := mgo.Dial(mongodb_url)
	if err != nil {
		log.Panic(err)
		os.Exit(-1)
	}
	// database is provided in url
	s.db = sess.DB("")

	// wait chan
	s.wait = make(chan string, BUFSIZ)
	go s.loader_task()
}

func (s *server) MarkDirty(ctx context.Context, in *pb.BgSave_Key) (*pb.BgSave_NullResult, error) {
	s.wait <- in.Name
	return &pb.BgSave_NullResult{}, nil
}

func (s *server) MarkDirties(ctx context.Context, in *pb.BgSave_Keys) (*pb.BgSave_NullResult, error) {
	for k := range in.Names {
		s.wait <- in.Names[k]
	}
	return &pb.BgSave_NullResult{}, nil
}

// background loader, copy chan into map, execute dump every SAVE_DELAY
func (s *server) loader_task() {
	dirty := make(map[string]bool)
	timer := time.After(SAVE_DELAY)
	timer_count := time.After(COUNT_DELAY)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM)

	var count uint64

	for {
		select {
		case key := <-s.wait:
			dirty[key] = true
		case <-timer:
			if len(dirty) > 0 {
				count += uint64(len(dirty))
				s.dump(dirty)
				dirty = make(map[string]bool)
			}
			timer = time.After(SAVE_DELAY)
		case <-timer_count:
			log.Info("num records saved:", count)
			timer_count = time.After(COUNT_DELAY)
		case <-sig:
			if len(dirty) > 0 {
				s.dump(dirty)
			}
			log.Info("SIGTERM")
			os.Exit(0)
		}
	}
}

// dump all dirty data into backend database
func (s *server) dump(dirty map[string]bool) {
	for k := range dirty {
		raw, err := s.redis_client.Cmd("GET", k).Bytes()
		if err != nil {
			log.Error(err)
			continue
		}

		// snappy
		if s.enable_snappy {
			if dec, err := snappy.Decode(nil, raw); err == nil {
				raw = dec
			} else {
				log.Error(err)
				continue
			}
		}

		// unpack message from msgpack format
		var record map[string]interface{}
		err = msgpack.Unmarshal(raw, &record)
		if err != nil {
			log.Error(err)
			continue
		}

		// split key into TABLE NAME and RECORD ID
		strs := strings.Split(k, ":")
		if len(strs) != 2 { // log the wrong key
			log.Error("cannot split key", k)
			continue
		}
		tblname, id_str := strs[0], strs[1]
		// save data to mongodb
		id, err := strconv.Atoi(id_str)
		if err != nil {
			log.Error(err)
			continue
		}

		_, err = s.db.C(tblname).Upsert(bson.M{"Id": id}, record)
		if err != nil {
			log.Error(err)
			continue
		}
	}
}
