package main

import (
	"chess/common/config"
	"chess/common/db"
       "chess/common/consul"
	"chess/models"
)

func main() {
	err := consul.InitConsulClientViaEnv()
	if err != nil {
		panic(err)
	}

	// TODO 换皮配置分发，可存储到mongodb
	err = config.Api.Import()
	if err != nil {
		panic(err)
	}

	//InitRpcWrapper()
	db.InitMySQL()
	db.InitMongo()
	models.Init()

	InitRouter()
}
