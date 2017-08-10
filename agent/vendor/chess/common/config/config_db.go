package config

import (
	. "chess/common/consul"
	"encoding/json"
	"time"
)

var Db = new(DbConfig)


type DbConfig struct {
	MySQL MySQLConf `json:"mysql"`
	Mongo MongoConf `json:"mongo"`
	Redis RedisConf `json:"redis"`
}

type MySQLConf struct {
	Setting MySQLSetting           `json:"setting"`
	Server  map[string]MySQLServer `json:"server"`
}

type MySQLSetting struct {
	ParseTime bool   `json:"parse_time"`
	Loc       string `json:"loc"`
}

type MySQLServer struct {
	Db       string `json:"db"`
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type MongoConf struct {
	Setting MongoSetting           `json:"setting"`
	Server  map[string]MongoServer `json:"server"`
}

type MongoSetting struct {
	DialTimeout time.Duration `json:"dial_timeout"`
}

type MongoServer struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type RedisConf struct {
	Setting RedisSetting           `json:"setting"`
	Server  map[string]RedisServer `json:"server"`
}

type RedisSetting struct {
	MaxIdle     int           `json:"max_idle"`
	MaxActive   int           `json:"max_active"`
	IdleTimeout time.Duration `json:"idle_timeout"`
	WaitIdle    bool          `json:"wait_idle"`
}

type RedisServer struct {
	Db       int    `json:"db"`
	Host     string `json:"host"`
	Password string `json:"password"`
}

func (c *DbConfig) Import(srvName string) error {
	key := srvName + "/db"

	val, err := ConsulClient.Key(key, "")
	if err != nil {
		return err
	}
	//ConsulClient.KeyWatch(key, &val)
	Db = c

	return json.Unmarshal([]byte(val), c)
}
