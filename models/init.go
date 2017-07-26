package models

import (
    "chess/common/db"
    "chess/common/log"
    "database/sql"
)

var (
	ChessMysql      MySQLDB
	ChessMongo      *db.MongoDB
        ChessRedis      RedisDB
)


type RedisDB struct {
    Main             *db.Redis
    Order            *db.Redis
    Sms              *db.Redis
    Captcha          *db.Redis
    Login            *db.Redis
    Config           *db.Redis
}

type MySQLDB struct {
    Main  *sql.DB
    Logs  *sql.DB
    Games *sql.DB
}
func Init() {
	var err error
	ChessMysql.Main, err = db.D("main")
	if err != nil {
		log.Warnf("db.D Error(%s)", err)
	}

	ChessMongo, err = db.M("main")
	if err != nil {
		log.Warnf("db.M Error(%s)", err)
	}

        ChessRedis.Main, err = db.R("main")
	if err != nil {
	    log.Warnf("db.R Error(%s)", err)
	}
	ChessRedis.Login, err = db.R("login")
	if err != nil {
	    log.Warnf("db.R Error(%s)", err)
	}
}
