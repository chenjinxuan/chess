package main

import (
	"chess/common/config"
	"chess/common/consul"
	"chess/common/db"
	"chess/common/log"
	"chess/models"
	"net"
	"net/http"
	"os"
	//"chess/common/db"
	"chess/common/define"
	"chess/common/services"
	pb "chess/srv/srv-room/proto"
	"chess/srv/srv-room/redis"
	"chess/srv/srv-room/signal"
	"fmt"
	"google.golang.org/grpc"
	cli "gopkg.in/urfave/cli.v2"
)

var Cfg = new(Config)

type Config struct {
	ServiceId string
	Address   string
	Port      int
}

func main() {
	app := &cli.App{
		Name:    "room",
		Usage:   "room service",
		Version: "1.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "service-id",
				Value: "room-1",
				Usage: "service id",
			},
			&cli.StringFlag{
				Name:  "address",
				Value: "127.0.0.1",
				Usage: "external address",
			},
			&cli.IntFlag{
				Name:  "port",
				Value: 20001,
				Usage: "listening port",
			},
			&cli.StringSliceFlag{
				Name:  "services",
				Value: cli.NewStringSlice("centre", "chat"),
				Usage: "auto-discovering services",
			},
		},
		Action: func(c *cli.Context) error {
			Cfg.ServiceId = c.String("service-id")
			Cfg.Address = c.String("address")
			Cfg.Port = c.Int("port")

			// 从consul读取配置，初始化数据库连接
			err := consul.InitConsulClientViaEnv()
			if err != nil {
				panic(err)
			}

			err = config.SrvRoom.Import()
			if err != nil {
				panic(err)
			}

			db.InitMySQL()
			db.InitRedis()
			db.InitMongo()
			models.Init()
			redis.Init()

			// consul 服务注册
			err = services.Register(c.String("service-id"), define.SRV_NAME_ROOM, c.String("address"), c.Int("port"), c.Int("port")+100, []string{"master"})
			if err != nil {
				panic(err)
			}
			// init services
			services.Discover(c.StringSlice("services"))

			go signal.Handler()

			// consul 健康检查
			http.HandleFunc("/check", consulCheck)
			go http.ListenAndServe(fmt.Sprintf(":%d", c.Int("port")+100), nil)

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
			pb.RegisterRoomServiceServer(s, ins)
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
