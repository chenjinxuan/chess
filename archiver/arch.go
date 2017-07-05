package main

import (
	"encoding/binary"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	nsq "github.com/bitly/go-nsq"
	"github.com/boltdb/bolt"
)

const (
	DEFAULT_NSQLOOKUPD   = "http://172.17.42.1:4161"
	ENV_NSQLOOKUPD       = "NSQLOOKUPD_HOST"
	TOPIC                = "REDOLOG"
	CHANNEL              = "ARCH"
	SERVICE              = "[ARCH]"
	REDO_TIME_FORMAT     = "REDO-2006-01-02T15:04:05.RDO"
	REDO_ROTATE_INTERVAL = 24 * time.Hour
	BOLTDB_BUCKET        = "REDOLOG"
	DATA_DIRECTORY       = "/data/"
	BATCH_SIZE           = 1024
	SYNC_INTERVAL        = 10 * time.Millisecond
)

type Archiver struct {
	pending chan []byte
	stop    chan bool
}

func (arch *Archiver) init() {
	arch.pending = make(chan []byte, BATCH_SIZE)
	arch.stop = make(chan bool)
	cfg := nsq.NewConfig()
	consumer, err := nsq.NewConsumer(TOPIC, CHANNEL, cfg)
	if err != nil {
		log.Panic(err)
		os.Exit(-1)
	}

	// message process
	consumer.AddHandler(nsq.HandlerFunc(func(msg *nsq.Message) error {
		arch.pending <- msg.Body
		return nil
	}))

	// read environtment variable
	addresses := []string{DEFAULT_NSQLOOKUPD}
	if env := os.Getenv(ENV_NSQLOOKUPD); env != "" {
		addresses = strings.Split(env, ";")
	}

	// connect to nsqlookupd
	log.Debug("connect to nsqlookupds ip:", addresses)
	if err := consumer.ConnectToNSQLookupds(addresses); err != nil {
		log.Error(err)
		return
	}
	log.Info("nsqlookupd connected")

	go arch.archive_task()
}

func (arch *Archiver) archive_task() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM)
	timer := time.After(REDO_ROTATE_INTERVAL)
	sync_ticker := time.NewTicker(SYNC_INTERVAL)
	db := arch.new_redolog()
	key := make([]byte, 8)
	for {
		select {
		case <-sync_ticker.C:
			n := len(arch.pending)
			if n == 0 {
				continue
			}

			db.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte(BOLTDB_BUCKET))
				for i := 0; i < n; i++ {
					id, err := b.NextSequence()
					if err != nil {
						log.Error(err)
						continue
					}
					binary.BigEndian.PutUint64(key, uint64(id))
					if err = b.Put(key, <-arch.pending); err != nil {
						log.Error(err)
						continue
					}
				}
				return nil
			})
		case <-timer:
			db.Close()
			// rotate redolog
			db = arch.new_redolog()
			timer = time.After(REDO_ROTATE_INTERVAL)
		case <-sig:
			db.Close()
			log.Info("SIGTERM")
			os.Exit(0)
		}
	}
}

func (arch *Archiver) new_redolog() *bolt.DB {
	file := DATA_DIRECTORY + time.Now().Format(REDO_TIME_FORMAT)
	log.Info(file)
	db, err := bolt.Open(file, 0600, nil)
	if err != nil {
		log.Panic(err)
		os.Exit(-1)
	}
	// create bulket
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(BOLTDB_BUCKET))
		if err != nil {
			log.Errorf("create bucket: %s", err)
			return err
		}
		return nil
	})
	return db
}
