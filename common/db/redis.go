package db

import (
	"chess/common/config"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

type Redis struct {
	Host, Pwd string
	Db        int
	Pool      *redis.Pool
}

func R(dbname string) (*Redis, error) {
	r, ok := redisMap[dbname]
	if !ok || r == nil {
		return nil, fmt.Errorf("%s Redis Not Found", dbname)
	}
	return r, nil
}

func InitRedis() {
	for dbname, c := range config.Db.Redis.Server {
		r := NewRedis(c.Host, c.Db, c.Password)
		redisMap[dbname] = r
	}
}

func NewRedis(host string, db int, password string) (r *Redis) {
	r = &Redis{Host: host, Db: db, Pwd: password}
	maxIdle := config.Db.Redis.Setting.MaxIdle
	maxActive := config.Db.Redis.Setting.MaxActive
	idleTimeout := config.Db.Redis.Setting.IdleTimeout
	waitIdle := config.Db.Redis.Setting.WaitIdle
	r.Pool = &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: idleTimeout * time.Second,
		Wait:        waitIdle,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", r.Host)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", r.Pwd); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return
}

func (r *Redis) GetConn() redis.Conn {
	if r.Pool == nil {
		return nil
	}
	conn := r.Pool.Get()
	conn.Do("SELECT", r.Db)
	return conn
}

func (r *Redis) Get(key string) (string, error) {
	conn := r.GetConn()
	defer conn.Close()

	return redis.String(conn.Do("GET", key))
}

func (r *Redis) GetInt(key string) (int, error) {
	conn := r.GetConn()
	defer conn.Close()

	return redis.Int(conn.Do("GET", key))
}

func (r *Redis) Set(key string, value string) error {
	conn := r.GetConn()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	return err
}

func (r *Redis) Setex(key string, value string, time int64) error {
	conn := r.GetConn()
	defer conn.Close()

	_, err := conn.Do("SETEX", key, time, value)
	return err
}
func (r *Redis) SetexInt(key string, value int, time int64) error {
	conn := r.GetConn()
	defer conn.Close()

	_, err := conn.Do("SETEX", key, time, value)
	return err
}
func (r *Redis) Exists(key string) (bool, error) {
	conn := r.GetConn()
	defer conn.Close()

	return redis.Bool(conn.Do("EXISTS", key))
}

func (r *Redis) Del(key string) error {
	conn := r.GetConn()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}

func (r *Redis) Expire(key string, time int64) error {
	conn := r.GetConn()
	defer conn.Close()

	_, err := conn.Do("EXPIRE", key, time)
	return err
}

func (r *Redis) EasyScan(keypattern string) ([]string, error) {
	conn := r.GetConn()
	defer conn.Close()

	iter := 0

	var keys []string
	for {
		if arr, err := redis.Values(conn.Do("SCAN", iter, "MATCH", keypattern)); err != nil {
			return keys, err
		} else {
			iter, _ = redis.Int(arr[0], nil)
			keys, _ = redis.Strings(arr[1], nil)
		}

		if iter == 0 {
			break
		}
	}

	return keys, nil
}

//小心使用，如果遇到操作自身key值的，考虑好逻辑
func (r *Redis) ScanWalk(keypattern string, f func(string, int) bool) {
	conn := r.GetConn()
	defer conn.Close()

	iter := 0
	_iter := 0

	var keys []string
	for {
		keys = []string{}
		if arr, err := redis.Values(conn.Do("SCAN", iter, "MATCH", keypattern)); err != nil {
			return
		} else {
			iter, _ = redis.Int(arr[0], nil)
			keys, _ = redis.Strings(arr[1], nil)
		}

		for _, key := range keys {
			isnext := f(key, _iter)
			_iter++
			if !isnext {
				return
			}
		}

		if iter == 0 {
			break
		}
	}

	return
}

func (r *Redis) Lpush(key string, val string) error {
	conn := r.GetConn()
	defer conn.Close()

	_, err := conn.Do("LPUSH", key, val)
	return err
}
func (r *Redis) Brpop(key string,time int) ([]string,error) {
    conn := r.GetConn()
    defer conn.Close()

    res, err := redis.Strings(conn.Do("BRPOP", key, time))
    return res,err
}
