package grpc

import (
	pb "chess/api/proto"
	"chess/common/services"
        "chess/common/define"
)


func GetAuthGrpc() (AuthClient pb.AuthServiceClient,ret int) {
	auth := services.GetService(define.SRV_NAME_AUTH)
	if auth==nil {
	    ret = 0
	    return nil,ret
	}
	AuthClient = pb.NewAuthServiceClient(auth)
    	ret=1
	return AuthClient ,ret
}

func GetCentreGrpc() (CentreClient pb.CentreServiceClient,ret int) {
	centre := services.GetService(define.SRV_NAME_CENTRE)
	if centre==nil {
	    ret = 0
	    return nil,ret
	}
	CentreClient = pb.NewCentreServiceClient(centre)
   	 ret=1
	return CentreClient,ret
}


func GetTaskGrpc() (TaskClient pb.TaskServiceClient,ret int) {
    task:= services.GetService(define.SRV_NAME_TASK)
    if task==nil {
	ret = 0
	return nil,ret
    }
    TaskClient = pb.NewTaskServiceClient(task)
    ret=1
    return TaskClient,ret
}
