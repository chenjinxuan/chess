package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"chess/common/config"
)

func InitMySQL() {
	var err error
	for dbname, server := range config.Db.MySQL.Server {
		mysqlMap[dbname], err = NewMySQL(server.Username, server.Password, server.Host, server.Db,
			config.Db.MySQL.Setting.ParseTime, config.Db.MySQL.Setting.Loc)
		if err != nil {
			panic(err)
		}
	}
}

func D(dbname string) (*sql.DB, error) {
	db, ok := mysqlMap[dbname]
	if !ok || db == nil {
		return nil, fmt.Errorf("%s DB Not Found", dbname)
	}

	return db, db.Ping()
}

func NewMySQL(username, password, host, dbname string, parseTime bool, loc string) (*sql.DB, error) {
	connStr := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8&parseTime=%v&loc=%s",
		username, password, host, dbname, parseTime, loc)

	m, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}

	return m, m.Ping()
}
