package main

import (
	"chess/common/log"
	pb "chess/srv/srv-room/proto"
	"encoding/binary"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"time"
)

var (
	//target = "192.168.40.157:20001"
	//serviceId = "room-1"
	target    = "192.168.40.157:30001"
	serviceId = "room-2"
)

func main() {
	// 用户id
	//uid := time.Now().Second()
	uid := 10000001

	player := NewPlayer()
	player.Id = int32(uid)

	conn, err := grpc.Dial(target, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	cli := pb.NewRoomServiceClient(conn)

	// 开启到游戏服的流
	ctx := metadata.NewContext(
		context.Background(),
		metadata.New(map[string]string{
			"userid":       fmt.Sprint(uid),
			"service_name": "room",
			"service_id":   serviceId,
			"unique_id":    fmt.Sprintf("xxxx-xxxxx-%d", time.Now().Unix()),
			//"unique_id":    fmt.Sprintf("xxxx-xxxxx-%d", 123),
		}),
	)

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
			if frame.Type == pb.Room_Message || frame.Type == pb.Room_Kick {
				c := int16(binary.BigEndian.Uint16(frame.Message[:2]))
				p.HandleMQ(c, frame.Message[2:])
			}

			//select {
			//case <-p.Die:
			//}
		}
	}
	go fetcher_task(player)

	// ping
	go func(p *Player) {
		for {
			frame := &pb.Room_Frame{
				Type:    pb.Room_Ping,
				Message: []byte{},
			}

			// check stream
			if p.Stream == nil {
				return
			}

			if err := p.Stream.Send(frame); err != nil {
				log.Error("Send room ping frame error:", err)
				return
			}

			time.Sleep(5 * time.Second)
		}
	}(player)

	player.CmdLoop()
}
