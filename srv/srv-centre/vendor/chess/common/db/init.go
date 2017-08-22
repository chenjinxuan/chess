package db

import (
	"database/sql"
)

var (
	mysqlMap map[string]*sql.DB
	redisMap map[string]*Redis
	mongoMap map[string]*MongoDB
)

func init() {
	mysqlMap = make(map[string]*sql.DB)
	redisMap = make(map[string]*Redis)
	mongoMap = make(map[string]*MongoDB)
}
