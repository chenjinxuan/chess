package models

import (
	"chess/common/db"
	"chess/common/log"
	"database/sql"
)

const (
	MongoDBStr       = "chess"
	MongoColSession  = "session"
	MongoColUserTask = "user_task"
	MongoColGambling = "gambling"
        MongoColUserBag  = "user_bag"
)

var (
	Mysql MySQLDB
	Mongo MongoDB
)

type MySQLDB struct {
	Chess *sql.DB
}

type MongoDB struct {
	Chess *db.MongoDB
}

func Init() {
	var err error
	Mysql.Chess, err = db.D("chess")
	if err != nil {
		log.Warnf("db.D(\"chess\") Error(%s)", err)
	}

	Mongo.Chess, err = db.M("chess")
	if err != nil {
		log.Warnf("db.M(\"chess\") Error(%s)", err)
	}
}
