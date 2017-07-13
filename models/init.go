package models

import (
	"chess/common/db"
	"chess/common/log"
	"database/sql"
)

var (
	Chess      *sql.DB
	ChessMongo *db.MongoDB
)

func Init() {
	var err error
	Chess, err = db.D("wstool")
	if err != nil {
		log.Warnf("db.D Error(%s)", err)
	}

	ChessMongo, err = db.M("wstool")
	if err != nil {
		log.Warnf("db.M Error(%s)", err)
	}
}
