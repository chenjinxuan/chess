package main

import (
	pb "chess/srv/srv-auth/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"testing"
)

const (
	address = "localhost:50001"
)

func TestAuth(t *testing.T) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewAuthServiceClient(conn)

	// Contact the server and print out its response.
	r, err := c.Auth(context.Background(), &pb.AuthArgs{
		UserId:     1,
		AppFrom:    "ios",
		AppVer:     100,
		AppChannel: "AppStore",
		UniqueId:   "xxxx-xxxx",
		Token:      "CA761232-ED42-11CE-BACD-00AA0057B223",
	})
	if err != nil {
		t.Fatalf("could not query: %v", err)
	}
	t.Log(r)

}
