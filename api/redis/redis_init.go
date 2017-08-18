package redis

import (
	"chess/common/db"
	"chess/common/log"
)

var (
	Redis RedisDB
)

type RedisDB struct {
	Main    *db.Redis
	Login   *db.Redis
	Captcha *db.Redis
	Sms     *db.Redis
}

func InitApiRedis() {
	var err error
	Redis.Main, err = db.R("main")
	if err != nil {
		log.Warnf("db.R(\"chess\") Error(%s)", err)
	}
	Redis.Login, err = db.R("login")
	if err != nil {
		log.Warnf("db.R(\"login\") Error(%s)", err)
	}
	Redis.Captcha, err = db.R("captcha")
	if err != nil {
		log.Warnf("db.R(\"captcha\") Error(%s)", err)
	}
	Redis.Sms, err = db.R("sms")
	if err != nil {
		log.Warnf("db.R(\"sms\") Error(%s)", err)
	}

}
