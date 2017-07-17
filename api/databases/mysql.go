package databases

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"treasure/config"
)

var (
	MySQL *MySQLDB
)

type MySQLDB struct {
	Main  *sql.DB
	Logs  *sql.DB
	Games *sql.DB
}

func NewMySQL(username, password, host, schema string) (*MySQLDB, error) {
	MySQL = new(MySQLDB)

	mainConnString := GenMysqlConnString(
		username,
		password,
		host,
		schema,
	) + "&loc=Asia%2FChongqing"

	var err error
	MySQL.Main, err = sql.Open("mysql", mainConnString)
	if err == nil {
		return MySQL, MySQL.Main.Ping()
	}
	return nil, err
}

func InitMySQL(c *config.Config) (*MySQLDB, error) {
	var err error
	mysql := new(MySQLDB)

	// init main
	mainConfig, ok := c.Databases.MySQL.Server["main"]
	if !ok {
		return nil, config.FormatError
	}
	mainConnString := GenMysqlConnString(
		mainConfig.Username,
		mainConfig.Password,
		mainConfig.Host,
		mainConfig.Db,
	) + "&loc=Asia%2FChongqing"
	mysql.Main, err = sql.Open("mysql", mainConnString)
	if err != nil {
		return nil, err
	}
	err = mysql.Main.Ping()
	if err != nil {
		return nil, err
	}

	// init logs
	logsConfig, ok := c.Databases.MySQL.Server["logs"]
	if !ok {
		return nil, config.FormatError
	}
	logsConnString := GenMysqlConnString(
		logsConfig.Username,
		logsConfig.Password,
		logsConfig.Host,
		logsConfig.Db,
	) + "&loc=Asia%2FChongqing"
	mysql.Logs, err = sql.Open("mysql", logsConnString)
	if err != nil {
		return nil, err
	}
	err = mysql.Logs.Ping()
	if err != nil {
		return nil, err
	}

	// init games
	gamesConfig, ok := c.Databases.MySQL.Server["games"]
	if !ok {
		return nil, config.FormatError
	}
	gamesConnString := GenMysqlConnString(
		gamesConfig.Username,
		gamesConfig.Password,
		gamesConfig.Host,
		gamesConfig.Db,
	) + "&loc=Asia%2FChongqing"
	mysql.Games, err = sql.Open("mysql", gamesConnString)
	if err != nil {
		return nil, err
	}
	err = mysql.Games.Ping()
	if err != nil {
		return nil, err
	}

	return mysql, nil
}

func GenMysqlConnString(dbUser, dbPass, dbHost, dbName string) string {
	mysqlConnFormat := "%v:%v@tcp(%v)/%v?charset=utf8&parseTime=true"
	mysqlConnString := fmt.Sprintf(mysqlConnFormat, dbUser, dbPass, dbHost, dbName)
	return mysqlConnString
}
