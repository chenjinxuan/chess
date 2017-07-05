package main

import (
	"net"
	"net/http"
	"os"
	"time"

	"github.com/xtaci/chat/kafka"

	cli "gopkg.in/urfave/cli.v2"

	log "github.com/Sirupsen/logrus"
	"github.com/xtaci/logrushooks"
	"google.golang.org/grpc"

	pb "github.com/xtaci/chat/proto"
)

func main() {
	log.AddHook(logrushooks.LineNoHook{})

	go func() {
		log.Info(http.ListenAndServe("0.0.0.0:6060", nil))
	}()
	app := &cli.App{
		Name: "chat",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "listen",
				Value: ":10000",
				Usage: "listening address:port",
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
			log.Println("listen:", c.String("listen"))
			log.Println("boltdb:", c.String("boltdb"))
			log.Println("kafka-brokers:", c.StringSlice("kafka-brokers"))
			log.Println("chat-topic", c.String("chat-topic"))
			log.Println("bucket:", c.String("bucket"))
			log.Println("retention:", c.Int("retention"))
			log.Println("write-interval:", c.Duration("write-interval"))
			log.Println("kafka-bucket", c.String("kafka-bucket"))
			log.Println("kafka-brokers", c.StringSlice("kafka-brokers"))
			// 监听
			lis, err := net.Listen("tcp", c.String("listen"))
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
