package task_redis

import (
	"chess/common/db"
	"chess/common/log"
)

var (
	Redis RedisDB
)

type RedisDB struct {
	Task *db.Redis
}

func InitTaskRedis() {
	var err error
	Redis.Task, err = db.R("task")
	if err != nil {
		log.Warnf("db.R(\"Task\") Error(%s)", err)
	}
}
