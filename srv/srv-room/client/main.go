package main

import (
	"google.golang.org/grpc"
	pb "chess/srv/srv-room/proto"
	"fmt"
	"google.golang.org/grpc/metadata"
	"golang.org/x/net/context"
	"io"
	"chess/common/log"
	"encoding/binary"
	"time"
)

var (
	target = "192.168.40.157:20001"
)

func main() {
	// 用户id
	uid := time.Now().Second()

	player := NewPlayer()
	player.Id = int32(uid)

	conn, err := grpc.Dial(target, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	cli := pb.NewRoomServiceClient(conn)

	// 开启到游戏服的流
	ctx := metadata.NewContext(context.Background(), metadata.New(map[string]string{"userid": fmt.Sprint(uid)}))
	stream, err := cli.Stream(ctx)
	if err != nil {
		panic(err)
	}

	player.Stream = stream
	//player.Die = make(chan struct{})

	defer func() {
		//close(player.Die)
		if player.Stream != nil {
			player.Stream.CloseSend()
		}
	}()

	// 读取room返回消息的goroutine
	fetcher_task := func(p *Player) {
		for {
			frame, err := p.Stream.Recv()
			if err == io.EOF { // 流关闭
				log.Debug(err)
				return
			}
			if err != nil {
				log.Error(err)
				return
			}

			// 读协议号
			c := int16(binary.BigEndian.Uint16(frame.Message[:2]))
			p.HandleMQ(c, frame.Message[2:])

			//select {
			//case <-p.Die:
			//}
		}
	}
	go fetcher_task(player)
	player.CmdLoop()
}