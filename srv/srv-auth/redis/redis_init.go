package redis

import(
    "chess/common/log"
    "chess/common/db"
)

var (
    Redis RedisDB
)

type RedisDB struct {
    Login *db.Redis
}

func InitAuthRedis() {
    var err error
    Redis.Login, err = db.R("login")
    if err != nil {
	log.Warnf("db.R(\"login\") Error(%s)", err)
    }
}