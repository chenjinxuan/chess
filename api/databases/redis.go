package databases

import (
	"errors"
	redislib "github.com/garyburd/redigo/redis"
	"strings"
	"time"
	"treasure/log"
)

var (
	Redis *RedisDB
)

type RedisDB struct {
	Main             *RedisPool
	Order            *RedisPool
	Sms              *RedisPool
	Captcha          *RedisPool
	Login            *RedisPool
	UserSts          *RedisPool
	CouponLuckyTimes *RedisPool
	Config           *RedisPool
	CouponNotify     *RedisPool
	CouponScript     *RedisPool
}

type RedisPool struct {
	Host        string
	Password    string
	Db          int
	Pool        *redislib.Pool
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
	WaitIdle    bool
}

func NewRedisPool(host, password string, db, maxIdle, maxActive int, idleTimeout time.Duration, waitIdle bool) *RedisPool {
	redis := new(RedisPool)
	redis.Host = host
	redis.Password = password
	redis.Db = db
	redis.MaxIdle = maxIdle
	redis.MaxActive = maxActive
	redis.IdleTimeout = idleTimeout
	redis.WaitIdle = waitIdle

	pool := &redislib.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: time.Duration(idleTimeout) * time.Second,
		Wait:        waitIdle,
		Dial: func() (redislib.Conn, error) {
			c, err := redislib.Dial("tcp", host)
			if err != nil {
				// TODO: log
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					// TODO: log
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redislib.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				// TODO: log
				log.Log.Errorf("redis:test on borrow fail(%s)", err)
				return err
			}
			return nil
		},
	}

	redis.Pool = pool
	return redis
}

func (r *RedisPool) GetConn() redislib.Conn {
	if r.Pool == nil {
		return nil
	}
	conn := r.Pool.Get()
	conn.Do("SELECT", r.Db)
	return conn
}

func (r *RedisPool) Get(key string) (string, error) {
	conn := r.GetConn()
	defer conn.Close()

	return redislib.String(conn.Do("get", key))
}

func (r *RedisPool) GetInt(key string) (int64, error) {
	conn := r.GetConn()
	defer conn.Close()

	return redislib.Int64(conn.Do("get", key))
}

func (r *RedisPool) Set(key string, val string) error {
	conn := r.GetConn()
	defer conn.Close()

	_, err := conn.Do("SET", key, val)
	return err
}

func (r *RedisPool) Unshift(key string, val string) error {
	conn := r.GetConn()
	defer conn.Close()

	_, err := conn.Do("LPUSH", key, val)
	return err
}
func (r *RedisPool) Sadd(key string, val string) error {
	conn := r.GetConn()
	defer conn.Close()

	_, err := conn.Do("SADD", key, val)
	return err
}
func (r *RedisPool) Spop(key string) (int64, error) {
	conn := r.GetConn()
	defer conn.Close()

	return redislib.Int64(conn.Do("SPOP", key))

}
func (r *RedisPool) Scard(key string) (int64, error) {
	conn := r.GetConn()
	defer conn.Close()

	return redislib.Int64(conn.Do("SCARD", key))
}

func (r *RedisPool) Setex(key string, val string, time int64) error {
	conn := r.GetConn()
	defer conn.Close()

	_, err := conn.Do("SETEX", key, time, val)
	return err
}

func (r *RedisPool) SetInt(key string, val int) error {
	conn := r.GetConn()
	defer conn.Close()

	_, err := conn.Do("SET", key, val)
	return err
}

func (r *RedisPool) SetIntEx(key string, val int, time int64) error {
	conn := r.GetConn()
	defer conn.Close()

	_, err := conn.Do("SETEX", key, time, val)
	return err
}

func (r *RedisPool) Del(key string) error {
	conn := r.GetConn()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}

func (r *RedisPool) Incr(key string) (int64, error) {
	conn := r.GetConn()
	defer conn.Close()

	return redislib.Int64(conn.Do("INCR", key))
}

func (r *RedisPool) Decr(key string) (int64, error) {
	conn := r.GetConn()
	defer conn.Close()

	return redislib.Int64(conn.Do("DECR", key))
}

func (r *RedisPool) Exists(key string) (bool, error) {
	conn := r.GetConn()
	defer conn.Close()

	return redislib.Bool(conn.Do("EXISTS", key))
}

func (r *RedisPool) Expire(key string, expire int) error {
	conn := r.GetConn()
	defer conn.Close()

	_, err := conn.Do("EXPIRE", key, expire)
	return err
}

func (r *RedisPool) GetTTL(key string) (int64, error) {
	conn := r.GetConn()
	defer conn.Close()

	return redislib.Int64(conn.Do("ttl", key))
}

func (r *RedisPool) GetVersion() (string, error) {
	conn := r.GetConn()
	defer conn.Close()
	version := ""

	res, err := redislib.String(conn.Do("info"))
	if err != nil {
		return "", err
	}

	resArr := strings.Split(res, "\r\n")
	for _, item := range resArr {
		keyVal := strings.Split(item, ":")
		if len(keyVal) == 2 {
			if keyVal[0] == "redis_version" {
				version = keyVal[1]
			}
		}
	}

	if version == "" {
		return "", errors.New("Could not find redis version")
	}

	return version, nil
}
