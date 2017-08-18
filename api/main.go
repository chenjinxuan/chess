package main

import (
	"chess/api/redis"
	"chess/common/config"
	"chess/common/consul"
	"chess/common/db"
	"chess/common/log"
	"chess/common/services"
	"chess/models"
	"flag"
	"fmt"
	"net/http"
	"os"
)

const (
	SERVICE_NAME = "api"
)

var (
	port    = flag.Int("port", 13333, "listen port")
	address = flag.String("address", "192.168.60.164", "external address")
	service = flag.String("service", "api-1", "external address")
	logPath = flag.String("log_path", "./log_data", "service id")
	debug   = flag.Bool("debug", false, "debug")
)

func main() {
	err := consul.InitConsulClient("192.168.40.117:8500", "lan-dc1", "", "")
	if err != nil {
		panic(err)
	}

	// TODO 换皮配置分发，CONSUL
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

	// consul 服务注册
	err = services.Register(*service, SERVICE_NAME, *address, *port+10, *port+10, []string{"master"})
	if err != nil {
		log.Error(err)
		os.Exit(-1)
	}
	// consul 健康检查
	http.HandleFunc("/check", consulCheck)
	go http.ListenAndServe(fmt.Sprintf(":%d", *port+10), nil)
	services.Discover(config.C.GrpcServer)
	InitRouter()
}

func consulCheck(w http.ResponseWriter, r *http.Request) {
	//log.Info("Consul Health Check!")
	fmt.Fprintln(w, "consulCheck")
}
