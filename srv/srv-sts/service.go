package main

import (
	"chess/common/define"
	"chess/common/log"
	. "chess/srv/srv-sts/proto"
	"chess/srv/srv-sts/redis"
	"encoding/json"
	"golang.org/x/net/context"
)

type server struct {
}

func (s *server) init() {
}

func (s *server) GameInfo(ctx context.Context, args *GameTableInfoArgs) (*StsRes, error) {
	log.Debug("gameInfo receive.")
	//判断数据是否是否收到
	if args.RoomId != 0 {
		//存入redis 队列
		dataByte, err := json.Marshal(args)
		if err != nil {
			log.Errorf(" GameInfo err %s", err)
			return &StsRes{Ret: 0, Msg: "recive fail."}, err
		}
		data := string(dataByte)
		sts_redis.Redis.Sts.Lpush(define.STS_GAME_INFO_REDIS_KEY, data)
		return &StsRes{Ret: 1, Msg: ""}, nil

	}
	return &StsRes{Ret: 0, Msg: "recive fail."}, nil
}
