package main

import (
"chess/common/config"
"chess/common/define"
"chess/common/log"
"chess/models"
"chess/srv/srv-auth/redis"
"errors"
"golang.org/x/net/context"
"regexp"
"strconv"
. "chess/srv/srv-task/proto"
)

var (
    ERROR_METHOD_NOT_SUPPORTED = errors.New("method not supoorted")
)
var (
    uuid_regexp = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
)

type server struct {
}

func (s *server) init() {
}


func (s *server) GameOver(ctx context.Context,args *GameInfoArgs) (*BaseRes,error){
     //判断数据是否是否收到
    if args.Winner != 0 {
	go s.loopProcessing(args)
	return &BaseRes{Ret:1,Msg:""},nil
	
    }
    return  &BaseRes{Ret:0,Msg:"recive fail."},nil
}

func (s *server) PlayerEvent(ctx context.Context,args *PlayerActionArgs)(*BaseRes,error){
    //判断数据是否是否收到
    if args.Id != 0 {
	go s.loopProcessing(args)
	return &BaseRes{Ret:1,Msg:""},nil

    }
    return  &BaseRes{Ret:0,Msg:"recive fail."},nil
}

func (s *server) loopProcessing(args *GameInfoArgs) {
    
}

func ()  {
    
}