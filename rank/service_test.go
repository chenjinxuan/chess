package main

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "proto"
	"testing"
)

const (
	address = "localhost:50001"
	KEY     = 0
)

func TestRankChange(t *testing.T) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address)
	if err != nil {
		t.Fatal("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewRankingServiceClient(conn)

	COUNT := 5000
	// Contact the server and print out its response.
	for i := 1; i < COUNT; i++ {
		_, err = c.RankChange(context.Background(), &pb.Ranking_Change{int32(i), int32(i), KEY})
		if err != nil {
			t.Fatalf("could not query: %v", err)
		}
		if i%1000 == 0 {
			list, err := c.QueryRankRange(context.Background(), &pb.Ranking_Range{1, 100, KEY})
			if err != nil {
				t.Fatalf("could not query: %v", err)
			}
			t.Log(list)
		}
	}
}
