package main

import (
	"chess/srv/srv-room/redis"
	redislib "github.com/garyburd/redigo/redis"
	"fmt"
	"chess/common/log"
	. "chess/srv/srv-room/texas_holdem"
	"time"
)

const(
	SESSION_STATUS_LOGIN = "login"
	SESSION_STATUS_LOGOUT = "logout"
)

type Session struct {
	Uid int "redis:uid"
	SrvId string "redis:srv_id"
	UniqueId string "redis:unique_id"
	Status string "redis:status"
	//Created int64  "redis:created"
	//Updated int64  "redis:updated"
}

func NewSession(uid int)*Session {
	return &Session{Uid: uid}
}

func (s *Session) key() string {
	return fmt.Sprintf("chess_session_%d", s.Uid)
}

func (s *Session) kickedOutKey() string {
	return fmt.Sprintf("kick_%d_%s_%s", s.Uid, s.SrvId, s.UniqueId)
}

func (s *Session) Get() error {
	conn := redis.Chess.GetConn()
	defer conn.Close()

	v, err := redislib.Values(conn.Do("HGETALL", s.key()))
	if err != nil {
		log.Errorf("Get Session(%d) Fail(%s)", s.Uid, err)
		return err
	}

	if err = redislib.ScanStruct(v, s); err != nil {
		log.Errorf("Get Session(%d) Fail(%s)", s.Uid, err)
		return err
	}
	return nil
}

func (s *Session) Set() error {
	conn := redis.Chess.GetConn()
	defer conn.Close()
	_, err := conn.Do("HMSET", redislib.Args{}.Add(s.key()).AddFlat(s)...)
	if err != nil {
		log.Errorf("Set Session(%d) Fail(%s)", s.Uid, err)
		return err
	}
	return nil
}

func (s *Session) Reset() error {
	conn := redis.Chess.GetConn()
	defer conn.Close()

	tmp := &Session{}
	v, err := redislib.Values(conn.Do("HGETALL", s.key()))
	if err != nil {
		log.Errorf("Get Session(%d) Fail(%s)", s.Uid, err)
		return err
	}

	if err = redislib.ScanStruct(v, tmp); err != nil {
		log.Errorf("Get Session(%d) Fail(%s)", s.Uid, err)
		return err
	}

	// 当前登录
	if tmp.UniqueId == s.UniqueId && tmp.SrvId == s.SrvId {
		log.Debugf("玩家%d注销登录", s.Uid)
		s.Status = SESSION_STATUS_LOGOUT
		_, err := conn.Do("HMSET", redislib.Args{}.Add(s.key()).AddFlat(s)...)
		if err != nil {
			log.Errorf("Set Session(%d) Fail(%s)", s.Uid, err)
			return err
		}

	}

	return nil
}
func (s *Session) DelKickedOutQueue() error {
	return redis.Chess.Del(s.kickedOutKey())
}

func (s *Session) NotifyKickedOut() error {
	log.Debugf("踢出玩家%d通知, queue key: %s",s.Uid,s.kickedOutKey())
	return redis.Chess.Lpush(s.kickedOutKey(), s.SrvId + s.UniqueId)
}

func (s *Session) KickedOutLoop(p *Player, sess_die chan struct{}) {
	for {
		log.Debugf("玩家%d 取踢出通知队列...", p.Id)
		conn := redis.Chess.GetConn()

		res, err := redislib.Strings(conn.Do("BRPOP", s.kickedOutKey(), 600))
		if err != nil {
			if err == redislib.ErrNil {
				log.Debugf("玩家%d KickedOutQueueTimeout...", p.Id)
			} else {
				log.Errorf("Get KickedOut Info Fail(%s)", err)
			}
		}

		conn.Close()


		if err == nil && res[1] == s.SrvId+s.UniqueId {
			log.Debugf("get KickedOut info(%s)", res)
			p.Flag |= SESS_KICKED_OUT
			return
		}

		select {
		case <- sess_die:
			log.Debugf("玩家%d KickedOutQueue 退出...", p.Id)
			return
		case <-time.After(1*time.Second):
			continue
		}
	}
}