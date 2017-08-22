package main

import (
	"chess/common/log"
	"golang.org/x/net/context"
	pb "chess/srv/srv-centre/proto"
	"chess/srv/srv-centre/redis"
	"sync"
	"errors"
	"encoding/json"
	redislib "github.com/garyburd/redigo/redis"
	"fmt"
)
var ErrArgsInvalid = errors.New("args invalid")

const (
	RedisKey = "srv-centre-data"
)

type server struct{
	mu *sync.RWMutex `json:"-"`

	Services map[string]map[int32]*pb.RoomInfo `json:"services"` // 对应room服的房间信息
	Summary map[int32]*pb.RoomInfo `json:"summary"` // 全服
}

func (s *server) init() {
	s.mu = new(sync.RWMutex)

	// get data from redis
	str, err := redis.Chess.Get(RedisKey)
	if err != nil && err != redislib.ErrNil {
		log.Errorf("Get server data fail(%s)", err)
		return
	}
	if str == "" {
		s.Services = make(map[string]map[int32]*pb.RoomInfo)
		s.Summary = make(map[int32]*pb.RoomInfo)
	} else {
		err = json.Unmarshal([]byte(str), s)
		if err != nil {
			log.Errorf("json.Unmarshal server data fail(%s)", err)
			return
		}
	}

	log.Debugf("server init: %+v", s)
}

func (s *server) Close() {
	// set data to redis
	strBytes, err := json.Marshal(s)
	if err != nil {
		fmt.Printf("json.Marshal server data fail(%s)", err)
	}

	err = redis.Chess.Set(RedisKey, string(strBytes))
	if err != nil {
		fmt.Printf("Set server data fail(%s)", err)
	} else {
		fmt.Println("Set server data to redis success!")
	}
}

func (s *server) RoomList(ctx context.Context, args *pb.RoomListArgs) (*pb.RoomListRes, error){
	s.mu.RLock()
	defer s.mu.RUnlock()

	return &pb.RoomListRes{
		List: s.Summary,
	}, nil
}

// 更新房间信息
func (s *server) UpdateRoomInfo(ctx context.Context, args *pb.UpdateRoomInfoArgs) (*pb.BaseRes, error){
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Debug("UpdateRoomInfo args: ", args)

	if args.ServiceId == "" || args.RoomId <= 0 || args.RoomInfo == nil {
		log.Debug("args invalid: %+v", args)
		return &pb.BaseRes{}, ErrArgsInvalid
	}

	if s.Services[args.ServiceId] == nil {
		s.Services[args.ServiceId] = make(map[int32]*pb.RoomInfo)
	}

	if s.Services[args.ServiceId][args.RoomId] == nil {
		s.Services[args.ServiceId][args.RoomId] = &pb.RoomInfo{}
	}

	s.Services[args.ServiceId][args.RoomId] = args.RoomInfo

	// 更新所有服务
	if s.Summary[args.RoomId] == nil {
		s.Summary[args.RoomId] = &pb.RoomInfo{}
	}

	s.Summary[args.RoomId].PlayerNumber = 0
	s.Summary[args.RoomId].TableNumber = 0
	for _, service := range s.Services {
		for rid, room := range service {
			if rid == args.RoomId {
				s.Summary[args.RoomId].PlayerNumber += room.PlayerNumber
				s.Summary[args.RoomId].TableNumber += room.TableNumber
			}
		}
	}

	return &pb.BaseRes{Ret:1}, nil
}
