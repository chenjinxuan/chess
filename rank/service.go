package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"golang.org/x/net/context"
)

import (
	. "rank/proto"
)

const (
	SERVICE = "[RANK]"
)

const (
	BOLTDB_FILE    = "/data/RANK-DUMP.DAT"
	BOLTDB_BUCKET  = "RANKING"
	CHANGES_SIZE   = 65536
	CHECK_INTERVAL = time.Minute // if ranking has changed, how long to check
)

var (
	OK                    = &Ranking_Nil{}
	ERROR_NAME_NOT_EXISTS = errors.New("name not exists")
)

type server struct {
	ranks   map[uint64]*RankSet
	pending chan uint64
	sync.RWMutex
}

func (s *server) init() {
	s.ranks = make(map[uint64]*RankSet)
	s.pending = make(chan uint64, CHANGES_SIZE)
	s.restore()
	go s.persistence_task()
}

func (s *server) lock_read(f func()) {
	s.RLock()
	defer s.RUnlock()
	f()
}

func (s *server) lock_write(f func()) {
	s.Lock()
	defer s.Unlock()
	f()
}

func (s *server) RankChange(ctx context.Context, p *Ranking_Change) (*Ranking_Nil, error) {
	// check name existence
	var rs *RankSet
	s.lock_write(func() {
		rs = s.ranks[p.SetId]
		if rs == nil {
			rs = NewRankSet()
			s.ranks[p.SetId] = rs
		}
	})

	// apply update on the rankset
	rs.Update(p.UserId, p.Score)
	s.pending <- p.SetId
	return OK, nil
}

func (s *server) QueryRankRange(ctx context.Context, p *Ranking_Range) (*Ranking_RankList, error) {
	var rs *RankSet
	s.lock_read(func() {
		rs = s.ranks[p.SetId]
	})

	if rs == nil {
		return nil, ERROR_NAME_NOT_EXISTS
	}

	ids, cups := rs.GetList(int(p.A), int(p.B))
	return &Ranking_RankList{UserIds: ids, Scores: cups}, nil
}

func (s *server) QueryUsers(ctx context.Context, p *Ranking_Users) (*Ranking_UserList, error) {
	var rs *RankSet
	s.lock_read(func() {
		rs = s.ranks[p.SetId]
	})

	if rs == nil {
		return nil, ERROR_NAME_NOT_EXISTS
	}

	ranks := make([]int32, 0, len(p.UserIds))
	scores := make([]int32, 0, len(p.UserIds))
	for _, id := range p.UserIds {
		rank, score := rs.Rank(id)
		ranks = append(ranks, rank)
		scores = append(scores, score)
	}
	return &Ranking_UserList{Ranks: ranks, Scores: scores}, nil
}

func (s *server) DeleteSet(ctx context.Context, p *Ranking_SetId) (*Ranking_Nil, error) {
	s.lock_write(func() {
		delete(s.ranks, p.SetId)
	})
	return OK, nil
}

func (s *server) DeleteUser(ctx context.Context, p *Ranking_DeleteUserRequest) (*Ranking_Nil, error) {
	var rs *RankSet
	s.lock_read(func() {
		rs = s.ranks[p.SetId]
	})
	if rs == nil {
		return nil, ERROR_NAME_NOT_EXISTS
	}
	rs.Delete(p.UserId)
	return OK, nil
}

// persistence ranking tree into db
func (s *server) persistence_task() {
	timer := time.After(CHECK_INTERVAL)
	db := s.open_db()
	changes := make(map[uint64]bool)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	for {
		select {
		case key := <-s.pending:
			changes[key] = true
		case <-timer:
			s.dump(db, changes)
			if len(changes) > 0 {
				log.Infof("perisisted %v rankset:", len(changes))
			}
			changes = make(map[uint64]bool)
			timer = time.After(CHECK_INTERVAL)
		case nr := <-sig:
			s.dump(db, changes)
			db.Close()
			log.Info(nr)
			os.Exit(0)
		}
	}
}

func (s *server) open_db() *bolt.DB {
	db, err := bolt.Open(BOLTDB_FILE, 0600, nil)
	if err != nil {
		log.Panic(err)
		os.Exit(-1)
	}
	// create bulket
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(BOLTDB_BUCKET))
		if err != nil {
			log.Panicf("create bucket: %s", err)
			os.Exit(-1)
		}
		return nil
	})
	return db
}

func (s *server) dump(db *bolt.DB, changes map[uint64]bool) {
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BOLTDB_BUCKET))
		for k := range changes {
			// marshal
			var rs *RankSet
			s.lock_read(func() {
				rs = s.ranks[k]
			})

			if rs == nil { // rankset deletion
				b.Delete([]byte(fmt.Sprint(k)))
			} else { // serialization and save
				bin, err := rs.Marshal()
				if err != nil {
					log.Error(err)
					continue
				}
				b.Put([]byte(fmt.Sprint(k)), bin)
			}
		}
		return nil
	})
}

func (s *server) restore() {
	// restore data from db file
	db := s.open_db()
	defer db.Close()
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BOLTDB_BUCKET))
		b.ForEach(func(k, v []byte) error {
			rs := NewRankSet()
			err := rs.Unmarshal(v)
			if err != nil {
				log.Panic("rank data corrupted:", err)
				os.Exit(-1)
			}
			id, err := strconv.ParseUint(string(k), 0, 64)
			if err != nil {
				log.Panic("rank data corrupted:", err)
				os.Exit(-1)
			}
			s.ranks[id] = rs
			return nil
		})
		return nil
	})
}
