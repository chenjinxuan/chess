package main

import (
	"net"
	"os"

	"chess/common/log"
	//"chess/common/config"
	//"chess/common/db"
	cli "gopkg.in/urfave/cli.v2"
	"google.golang.org/grpc"
	pb "chess/srv/srv-room/proto"
)


func main() {
	app := &cli.App{
		Name:    "auth",
		Usage:   "auth service",
		Version: "2.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "port",
				Value: ":20001",
				Usage: "listening address:port",
			},
		},
		Action: func(c *cli.Context) error {
			// TODO 从consul读取配置，初始化数据库连接
			//err := config.InitConsulClientViaEnv()
			//if err != nil {
			//	panic(err)
			//}
			//err = config.Api.Import()
			//if err != nil {
			//	panic(err)
			//}
			//
			////InitRpcWrapper()
			//db.InitMySQL()
			//db.InitMongo()
			//models.Init()


			// 监听
			lis, err := net.Listen("tcp", c.String("port"))
			if err != nil {
				log.Error(err)
				os.Exit(-1)
			}
			log.Info("listening on ", lis.Addr())

			// 注册服务
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
