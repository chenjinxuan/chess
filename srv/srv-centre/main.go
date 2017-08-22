package main

import (
	"chess/common/config"
	"chess/common/consul"
	"chess/common/db"
	"chess/common/define"
	"chess/common/helper"
	"chess/common/log"
	"chess/common/services"
	pb "chess/srv/srv-centre/proto"
	"chess/srv/srv-centre/redis"
	"fmt"
	"google.golang.org/grpc"
	cli "gopkg.in/urfave/cli.v2"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var Cfg = new(Config)

type Config struct {
	ServiceId string
	Address   string
	Port      int
}

func main() {
	app := &cli.App{
		Name:    "centre",
		Usage:   "centre service",
		Version: "1.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "service-id",
				Value: "centre-1",
				Usage: "service id",
			},
			&cli.StringFlag{
				Name:  "address",
				Value: "127.0.0.1",
				Usage: "external address",
			},
			&cli.IntFlag{
				Name:  "port",
				Value: 10001,
				Usage: "listening port",
			},
		},
		Action: func(c *cli.Context) error {
			Cfg.ServiceId = c.String("service-id")
			Cfg.Address = c.String("address")
			Cfg.Port = c.Int("port")

			// TODO 从consul读取配置，初始化数据库连接
			err := consul.InitConsulClientViaEnv()
			if err != nil {
				panic(err)
			}

			// consul 服务注册
			err = services.Register(c.String("service-id"), define.SRV_NAME_CENTRE, c.String("address"), c.Int("port"), c.Int("port")+10, []string{"master"})
			if err != nil {
				panic(err)
			}

			err = config.SrvCentre.Import()
			if err != nil {
				panic(err)
			}

			db.InitRedis()
			redis.Init()

			// consul 健康检查
			http.HandleFunc("/check", consulCheck)
			go http.ListenAndServe(fmt.Sprintf(":%d", c.Int("port")+10), nil)

			// grpc监听
			laddr := fmt.Sprintf(":%d", c.Int("port"))
			lis, err := net.Listen("tcp", laddr)
			if err != nil {
				panic(err)
			}
			log.Info("listening on ", lis.Addr())

			// 注册grpc服务
			s := grpc.NewServer()
			ins := &server{}
			ins.init()
			go signalHandler(ins)
			pb.RegisterCentreServiceServer(s, ins)
			// 开始服务
			return s.Serve(lis)
		},
	}
	app.Run(os.Args)
}

func consulCheck(w http.ResponseWriter, r *http.Request) {
	//log.Info("Consul Health Check!")
	fmt.Fprintln(w, "consulCheck")
}

func signalHandler(s *server) {
	defer helper.PrintPanicStack()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT)

	for {
		msg := <-ch
		switch msg {
		case syscall.SIGTERM, syscall.SIGINT:// 关闭room
			fmt.Println("waiting for set server data to redis, please wait...")
			s.Close()
			fmt.Println("centre shutdown.")
			os.Exit(0)
		case syscall.SIGHUP:
			return
		}
	}
}