package sts_redis

import (
	"chess/common/db"
	"chess/common/log"
)

var (
	Redis RedisDB
)

type RedisDB struct {
	Sts *db.Redis
}

func InitStsRedis() {
	var err error
	Redis.Sts, err = db.R("sts")
	if err != nil {
		log.Warnf("db.R(\"sts\") Error(%s)", err)
	}
}
