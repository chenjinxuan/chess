package main

import (
	"chess/common/config"
	"chess/common/db"
        "chess/common/consul"
	"chess/models"
        "chess/api/redis"
    "chess/common/services"
    "chess/common/log"
    "os"
    "flag"
    "net/http"
    "fmt"
)

const (
    SERVICE_NAME        = "api"
)
var (
    port             = flag.Int("port", 13333, "listen port")
    address          = flag.String("address","192.168.60.164","external address")
    service          = flag.String("service","api-1","external address")
    logPath          = flag.String("log_path", "./log_data", "service id")
    debug            = flag.Bool("debug", false, "debug")
)


func main() {
	err := consul.InitConsulClient("127.0.0.1:8500","local","","")
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
