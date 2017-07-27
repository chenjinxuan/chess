package main

import (
	"chess/common/config"
	"chess/common/db"
        "chess/common/consul"
	"chess/models"
        "chess/api/redis"
)

func main() {
	err := consul.InitConsulClient("127.0.0.1:8500","local","","")
	if err != nil {
		panic(err)
	}


	// TODO 换皮配置分发，可存储到mongodb
	err = config.Api.Import()
	if err != nil {
		panic(err)
	}

	//InitRpcWrapper()
        config.InitConfig()
	db.InitMySQL()
	db.InitMongo()
        db.InitRedis()
        //init redis
        redis.InitApiRedis()
	models.Init()
	InitRouter()
}


