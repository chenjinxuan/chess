package main

import (
	"chess/common/define"
	"chess/common/log"
	. "chess/srv/srv-task/proto"
	"chess/srv/srv-task/redis"
	"encoding/json"
	"golang.org/x/net/context"
	"strconv"
)

type server struct {
}

func (s *server) init() {
}

func (s *server) GameOver(ctx context.Context, args *GameTableInfoArgs) (*TaskRes, error) {
	log.Debugf("gameOver receive.", args.TableId, args.End)
	//判断数据是否是否收到
	if args.RoomId != 0 {
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
	log.Debugf("PlayerEvent receive.", args.Id)
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

func (s *server) UpsetTask(ctx context.Context, args *UpsetTaskArgs) (*TaskRes, error) {
	log.Debugf("Upset receive.")
	//判断数据是否是否收到
	if args.Id != 0 {
		//存入redis 队列
		_ = task_redis.Redis.Task.Sadd(define.TaskUpsetRedisKey, strconv.Itoa(int(args.Id)))
		return &TaskRes{Ret: 1, Msg: ""}, nil
	}
	return &TaskRes{Ret: 0, Msg: "recive fail."}, nil
}

func (s *server) IncrUserBag(ctx context.Context, args *UpdateBagArgs) (*TaskRes, error) {
	log.Debug(" UpsetUserBag receive.")
	if args.UserId != 0 {
		//存入redis 队列
		dataByte, err := json.Marshal(args)
		if err != nil {
			log.Errorf("UpsetUserBag err %s", err)
			return &TaskRes{Ret: 0, Msg: "recive fail."}, err
		}
		data := string(dataByte)
		task_redis.Redis.Task.Lpush(define.TaskUserBagRedisKey, data)
		return &TaskRes{Ret: 1, Msg: ""}, nil
	}
	return &TaskRes{Ret: 0, Msg: "recive fail."}, nil
}
