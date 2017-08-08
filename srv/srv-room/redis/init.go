package redis

import (
	"chess/common/db"
	"chess/common/log"
)

var (
	Chess *db.Redis
)

func Init() {
	var err error
	Chess, err = db.R("chess")
	if err != nil {
		log.Warnf("db.R(\"chess\") Error(%s)", err)
	}
}
