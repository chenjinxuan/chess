package main

import (
	. "chess/srv/srv-chat/proto"
	"testing"
	"io"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"fmt"
	"time"
)

const (
	address = "127.0.0.1:30001"
	ChatId = "1"
)

var (
	conn *grpc.ClientConn
	err  error
)

func init(){
	// Set up a connection to the server.
	conn, err = grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("did not connect: %v", err)
	}
}

func TestServer_Reg(t *testing.T) {
	c := NewChatServiceClient(conn)

	// Contact the server and print out its response.
	_, err = c.Reg(context.Background(), &Chat_Id{Id: ChatId})
	if err != nil {
		t.Logf("could not query: %v", err)
	}
}

func TestChat(t *testing.T) {
	count:=10

	go recv(&Chat_Consumer{Id: ChatId}, count)
	go recv(&Chat_Consumer{Id: ChatId}, count)

	send(&Chat_Message{Id: ChatId, Body: []byte("Hello")}, count)
	time.Sleep(1*time.Minute)
}


func send(m *Chat_Message, count int) {
	c := NewChatServiceClient(conn)
	for {
		if count == 0 {
			return
		}
		_, err := c.Send(context.Background(), m)
		if err != nil {
			fmt.Println("c.Send error: ", err)
		} else {
			fmt.Println("send:", m)
		}

		count--
	}
}

func recv(chat *Chat_Consumer, count int) {
	c := NewChatServiceClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	stream, err := c.Subscribe(ctx, chat)
	if err != nil {
		fmt.Println("c.Subscribe error: ", err)
		return
	}

	for {
		fmt.Println(count)
		if count == 0 {
			return
		}
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("stream.Recv error: ", err)
			return
		}
		fmt.Println("recv:", message)
		count--
		cancel() // recv should continue until error
	}

	fmt.Println("recv exit.....")
}

//func TestServer_Subscribe(t *testing.T) {
//	c := NewChatServiceClient(conn)
//	ctx, cancel := context.WithCancel(context.Background())
//	stream, err := c.Subscribe(ctx, &Chat_Consumer{Id:"1", From:-1})
//	if err != nil {
//		t.Fatal(err)
//	}
//	count := 10
//	for {
//		if count == 0 {
//			return
//		}
//		message, err := stream.Recv()
//		if err == io.EOF {
//			break
//		}
//		if err != nil {
//			t.Log(err)
//			return
//		}
//		println("recv:", count)
//		t.Log("recv:", message)
//		count--
//		cancel() // recv should continue until error
//	}
//}