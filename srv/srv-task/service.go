package main

import (
	"chess/common/define"
	"chess/common/log"
	. "chess/srv/srv-task/proto"
	"chess/srv/srv-task/redis"
	"encoding/json"
	"errors"
	"golang.org/x/net/context"
	"regexp"
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

func (s *server) GameOver(ctx context.Context, args *GameInfoArgs) (*TaskRes, error) {
	log.Debug("gameOver receive.")
	//判断数据是否是否收到
	if args.Winner != 0 {
		//存入redis 队列
		dataByte, err := json.Marshal(args)
		if err != nil {
			log.Errorf(" GameOver err %s", err)
			return &TaskRes{Ret: 0, Msg: "recive fail."}, err
		}
		data := string(dataByte)
		task_redis.Redis.Task.Lpush(define.TaskLoopHandleGameOverRedisKey, data)
		return &TaskRes{Ret: 1, Msg: ""}, nil

	}
	return &TaskRes{Ret: 0, Msg: "recive fail."}, nil
}

func (s *server) PlayerEvent(ctx context.Context, args *PlayerActionArgs) (*TaskRes, error) {
	log.Debug("PlayerEvent receive.")
	//判断数据是否是否收到
	if args.Id != 0 {
		//存入redis 队列
		dataByte, err := json.Marshal(args)
		if err != nil {
			log.Errorf("PlayerEvent err %s", err)
			return &TaskRes{Ret: 0, Msg: "recive fail."}, err
		}
		data := string(dataByte)
		task_redis.Redis.Task.Lpush(define.TaskLoopHandlePlayerEventRedisKey, data)
		return &TaskRes{Ret: 1, Msg: ""}, nil
	}
	return &TaskRes{Ret: 0, Msg: "recive fail."}, nil
}
