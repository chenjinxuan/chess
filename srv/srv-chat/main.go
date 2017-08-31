package main

import (
	"net"
	"net/http"
	"os"
	"time"

	"chess/common/consul"
	"chess/common/define"
	"chess/common/services"
	"chess/srv/srv-chat/kafka"

	cli "gopkg.in/urfave/cli.v2"

	log "github.com/Sirupsen/logrus"
	"google.golang.org/grpc"

	pb "chess/srv/srv-chat/proto"
	"fmt"
)

func main() {
	//log.AddHook(logrushooks.LineNoHook{})

	//go func() {
	//	log.Info(http.ListenAndServe("0.0.0.0:6060", nil))
	//}()
	app := &cli.App{
		Name: "chat",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "service-id",
				Value: "chat-1",
				Usage: "service id",
			},
			&cli.StringFlag{
				Name:  "address",
				Value: "127.0.0.1",
				Usage: "external address",
			},
			&cli.IntFlag{
				Name:  "port",
				Value: 13001,
				Usage: "listening port",
			},
			&cli.IntFlag{
				Name:  "check-port",
				Value: 13101,
				Usage: "consul health check listening port",
			},
			&cli.StringFlag{
				Name:  "kafka-bucket",
				Value: "kafka-bucket",
				Usage: "key with kafka offset",
			},
			&cli.StringFlag{
				Name:  "chat-topic",
				Value: "chat_updates",
				Usage: "chat topic in kafka",
			},
			&cli.StringSliceFlag{
				Name:  "kafka-brokers",
				Value: cli.NewStringSlice("127.0.0.1:9092"),
				Usage: "kafka brokers address",
			},
			&cli.StringFlag{
				Name:  "boltdb",
				Value: "/data/CHAT.DAT",
				Usage: "chat snapshot file",
			},
			&cli.StringFlag{
				Name:  "bucket",
				Value: "EPS",
				Usage: "bucket name",
			},
			&cli.IntFlag{
				Name:  "retention",
				Value: 1024,
				Usage: "retention number of messags for each endpoints",
			},
			&cli.DurationFlag{
				Name:  "write-interval",
				Value: 10 * time.Minute,
				Usage: "chat message persistence interval",
			},
		},

		Action: func(c *cli.Context) error {
			log.Println("boltdb:", c.String("boltdb"))
			log.Println("kafka-brokers:", c.StringSlice("kafka-brokers"))
			log.Println("chat-topic", c.String("chat-topic"))
			log.Println("bucket:", c.String("bucket"))
			log.Println("retention:", c.Int("retention"))
			log.Println("write-interval:", c.Duration("write-interval"))
			log.Println("kafka-bucket", c.String("kafka-bucket"))

			// 从consul读取配置，初始化数据库连接
			err := consul.InitConsulClientViaEnv()
			if err != nil {
				panic(err)
			}

			// consul 服务注册
			err = services.Register(c.String("service-id"), define.SRV_NAME_CHAT, c.String("address"), c.Int("port"), c.Int("check-port"), []string{"master"})
			if err != nil {
				panic(err)
			}

			// consul 健康检查
			http.HandleFunc("/check", consulCheck)
			go http.ListenAndServe(fmt.Sprintf(":%d", c.Int("check-port")), nil)

			// 监听
			laddr := fmt.Sprintf(":%d", c.Int("port"))
			lis, err := net.Listen("tcp", laddr)
			if err != nil {
				log.Panic(err)
				os.Exit(-1)
			}
			log.Info("listening on:", lis.Addr())

			kafka.Init(c)
			// 注册服务
			s := grpc.NewServer()
			ins := &server{}
			ins.init(c)
			pb.RegisterChatServiceServer(s, ins)
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
