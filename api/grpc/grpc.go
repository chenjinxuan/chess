package grpc

import (
   pb "chess/api/proto"
    "chess/common/config"
    "chess/common/services"
)
var (
    AuthClient pb.AuthServiceClient
)

func Init()  {
    services.Discover(config.C.GrpcServer)
    auth:=services.GetService(config.C.GrpcServer[0])
    AuthClient= pb.NewAuthServiceClient(auth)
}