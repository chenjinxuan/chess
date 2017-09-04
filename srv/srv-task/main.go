package main

import (
	"chess/common/config"
	"chess/common/consul"
	"chess/common/db"
	"chess/common/define"
	"chess/common/log"
	"chess/common/services"
	"chess/models"
	"chess/srv/srv-task/handler"
	pb "chess/srv/srv-task/proto"
	"chess/srv/srv-task/redis"
	"fmt"
	"google.golang.org/grpc"
	cli "gopkg.in/urfave/cli.v2"
	"net"
	"net/http"
	"os"
)

func main() {
	app := &cli.App{
		Name:    "task",
		Usage:   "task service",
		Version: "2.0",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "port",
				Value: 15001,
				Usage: "listening address:port",
			},
			&cli.IntFlag{
				Name:  "check-port",
				Value: 15101,
				Usage: "consul health check listening address:port",
			},
			&cli.StringFlag{
				Name:  "service-id",
				Value: "task-1",
				Usage: "service id",
			},
			&cli.StringFlag{
				Name:  "address",
				Value: "127.0.0.1",
				Usage: "external address",
			},
		},
		Action: func(c *cli.Context) error {
			// 从consul读取配置，初始化数据库连接
			err := consul.InitConsulClientViaEnv()
			if err != nil {
				panic(err)
			}
			err = config.SrvTask.Import()
			if err != nil {
				panic(err)
			}

			//InitRpcWrapper()
			db.InitMySQL()
			db.InitMongo()
			db.InitRedis()
			task_redis.InitTaskRedis()
			models.Init()

			err = services.Register(c.String("service-id"), define.SRV_NAME_TASK, c.String("address"), c.Int("port"), c.Int("check-port"), []string{"master"})
			if err != nil {
				log.Error(err)
				os.Exit(-1)
			}
			// consul 健康检查
			http.HandleFunc("/check", consulCheck)
			go http.ListenAndServe(fmt.Sprintf(":%d", c.Int("check-port")), nil)
			// grpc监听
			laddr := fmt.Sprintf(":%d", c.Int("port"))
			lis, err := net.Listen("tcp", laddr)
			if err != nil {
				log.Error(err)
				os.Exit(-1)
			}
			log.Info("listening on ", lis.Addr())
			//初始化handler
			mgr := handler.GetTaskHandlerMgr()
			if mgr == nil {
				log.Errorf("Get GetTaskHandlerMgr fail")
				os.Exit(-1)
			}
			mgr.Loop()
			mgr.SubLoop()
			//初始化update
			upsetMgr := handler.GetTaskUpsetMgr()
			if upsetMgr == nil {
				log.Errorf("Get GetTaskHandlerMgr fail")
				os.Exit(-1)
			}
			upsetMgr.Loop()
			upsetMgr.SubLoop()
		        ChargeMgr:=handler.GetTaskChargeMgr()
			if ChargeMgr == nil {
			    log.Errorf("Get GetTaskHandlerMgr fail")
			    os.Exit(-1)
			}
		        ChargeMgr.Loop()
		    	ChargeMgr.SubLoop()
			// 注册服务
			s := grpc.NewServer()
			ins := &server{}
			ins.init()
			pb.RegisterTaskServiceServer(s, ins)
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
