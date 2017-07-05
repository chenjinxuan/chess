package main

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "chess/auth/proto"
	"testing"
)

const (
	address = "localhost:50006"
)

func TestAuthUUID(t *testing.T) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address)
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewAuthServiceClient(conn)
	test_uuid := "CA761232-ED42-11CE-BACD-00AA0057B223"
	// Contact the server and print out its response.
	r, err := c.Auth(context.Background(), &pb.Auth_Certificate{Type: pb.Auth_UUID, Proof: []byte(test_uuid)})
	if err != nil {
		t.Fatalf("could not query: %v", err)
	}
	t.Log(r)

	r, err = c.Auth(context.Background(), &pb.Auth_Certificate{Type: pb.Auth_UUID, Proof: []byte(test_uuid + "XXX")})
	if err != nil {
		t.Fatalf("could not query: %v", err)
	}
	t.Log(r)
}
