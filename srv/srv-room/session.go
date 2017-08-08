package main

import (
	"chess/srv/srv-room/redis"
	redislib "github.com/garyburd/redigo/redis"
	"fmt"
	"chess/common/log"
	. "chess/srv/srv-room/texas_holdem"
)

const(
	SESSION_STATUS_LOGIN = "login"
	SESSION_STATUS_LOGOUT = "logout"

	QUIT_LOOP_MSG = "quit_loop"
)

type Session struct {
	TraceId string "redis:trace_id"
	Uid int "redis:uid"
	SrvId string "redis:srv_id"
	UniqueId string "redis:unique_id"
	Status string "redis:status"
	//Created int64  "redis:created"
	//Updated int64  "redis:updated"
}

func NewSession(uid int)*Session {
	return &Session{
		Uid: uid,
	}
}

func (s *Session) key() string {
	return fmt.Sprintf("chess_session_%d", s.Uid)
}

func (s *Session) kickedOutKey() string {
	return fmt.Sprintf("kick_%d_%s", s.Uid, s.TraceId)
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
	if tmp.TraceId == s.TraceId {
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
	redis.Chess.Lpush(s.kickedOutKey(), QUIT_LOOP_MSG)
	return redis.Chess.Del(s.kickedOutKey())
}

func (s *Session) NotifyKickedOut() error {
	log.Debugf("踢出玩家%d通知, queue key: %s",s.Uid,s.kickedOutKey())
	return redis.Chess.Lpush(s.kickedOutKey(), s.SrvId + s.UniqueId + s.TraceId)
}

func (s *Session) KickedOutLoop(p *Player, sess_die chan struct{}) {
	for {
		select {
		case <- sess_die:
			log.Debugf("退出登录踢出监听协程: 玩家%d 设备号%s ", p.Id, s.UniqueId)
			return
		default:
			conn := redis.Chess.GetConn()

			res, err := redislib.Strings(conn.Do("BRPOP", s.kickedOutKey(), 600))
			if err != nil {
				if err == redislib.ErrNil {
					log.Debugf("取登录踢出消息超时: 玩家%d 设备号%s ", p.Id, s.UniqueId)
				} else {
					log.Debugf("取登录踢出消息失败（%s）: 玩家%d 设备号%s ", err, p.Id, s.UniqueId)
				}
			}

			conn.Close()

			if err == nil {
				if res[1] == s.SrvId+s.UniqueId+s.TraceId {
					log.Debugf("收到登录踢出通知(%s)", res)
					p.Flag |= SESS_KICKED_OUT
					return
				}
				if res[1] == QUIT_LOOP_MSG {
					log.Debugf("退出登录踢出监听协程: 玩家%d 设备号%s ", p.Id, s.UniqueId)
					return
				}

			}
		}
	}
}