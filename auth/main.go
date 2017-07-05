package main

import (
	"net"
	"os"

	log "chess/common/log"
	_ "github.com/gonet2/libs/statsd-pprof"
	"google.golang.org/grpc"
)

import (
	pb "chess/grpc/auth"
)

func main() {
	// 监听
	lis, err := net.Listen("tcp", _port)
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
	s.Serve(lis)
}
