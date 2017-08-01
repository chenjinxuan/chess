package main

import (
	"chess/common/consul"
	"chess/common/log"
	"chess/common/config"
	"chess/common/db"
	"chess/models"
	"net"
	"net/http"
	"os"
	//"chess/common/db"
	"chess/common/services"
	pb "chess/srv/srv-room/proto"
	"fmt"
	"google.golang.org/grpc"
	cli "gopkg.in/urfave/cli.v2"
)

func main() {
	app := &cli.App{
		Name:    "auth",
		Usage:   "auth service",
		Version: "2.0",
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
		},
		Action: func(c *cli.Context) error {
			// TODO 从consul读取配置，初始化数据库连接
			err := consul.InitConsulClientViaEnv()
			if err != nil {
				panic(err)
			}

			err = config.SrvRoom.Import()
			if err != nil {
				panic(err)
			}

			db.InitMySQL()
			//db.InitMongo()
			models.Init()

			// consul 服务注册
			err = services.Register(c.String("service-id"), SERVICE_NAME, c.String("address"), c.Int("port"), c.Int("port")+10, []string{"master"})
			if err != nil {
				panic(err)
			}

			// consul 健康检查
			http.HandleFunc("/check", consulCheck)
			go http.ListenAndServe(fmt.Sprintf(":%d", c.Int("port")+10), nil)

			// grpc监听
			laddr := fmt.Sprintf(":%d", c.Int("port"))
			lis, err := net.Listen("tcp", laddr)
			if err != nil {
				log.Error(err)
				os.Exit(-1)
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
