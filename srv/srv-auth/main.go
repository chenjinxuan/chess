package main

import (
	"net"
	"os"

	"chess/common/log"
	"chess/common/define"
	"chess/common/consul"
        "chess/common/config"
	"chess/common/db"
	pb "chess/srv/srv-auth/proto"
	"google.golang.org/grpc"
	cli "gopkg.in/urfave/cli.v2"
        "chess/models"
        "chess/common/services"
    "fmt"
    "net/http"
)

const (
    SERVICE_NAME = "authserver"
)
func main() {
	app := &cli.App{
		Name:    "auth",
		Usage:   "auth service",
		Version: "2.0",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "port",
				Value: 50011,
				Usage: "listening address:port",
			},
			&cli.StringFlag{
			    Name:  "service-id",
			    Value: "auth-1",
			    Usage: "service id",
			},
			&cli.StringFlag{
			    Name:  "address",
			    Value: "192.168.60.164",
			    Usage: "external address",
			},
		},
		Action: func(c *cli.Context) error {
			// TODO 从consul读取配置，初始化数据库连接
		        err := consul.InitConsulClient("127.0.0.1:8500","local","","")
			if err != nil {
				panic(err)
			}
			err = config.SrvAuth.Import()
			if err != nil {
				panic(err)
			}

			//InitRpcWrapper()
			db.InitMySQL()
			db.InitMongo()
			models.Init()

		       err = services.Register(c.String("service-id"), define.SRV_NAME_AUTH, c.String("address"), c.Int("port"), c.Int("port")+10, []string{"master"})
		    if err != nil {
			log.Error(err)
			os.Exit(-1)
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

			// 注册服务
			s := grpc.NewServer()
			ins := &server{}
			ins.init()
			pb.RegisterAuthServiceServer(s, ins)
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
