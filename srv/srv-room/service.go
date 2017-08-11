package main

import (
	"chess/common/define"
	"chess/common/helper"
	"chess/common/log"
	"chess/srv/srv-room/client_handler"
	"chess/srv/srv-room/misc/packet"
	pb "chess/srv/srv-room/proto"
	"chess/srv/srv-room/registry"
	. "chess/srv/srv-room/texas_holdem"
	"encoding/binary"
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"io"
	"strconv"
	"fmt"
	"time"
	"chess/srv/srv-room/signal"
)

var (
	ERROR_INCORRECT_FRAME_TYPE = errors.New("incorrect frame type")
	ERROR_SERVICE_NOT_BIND     = errors.New("service not bind")
)

type server struct{}

func (s *server) init() {
	// Todo 从mysql取房间列表
	InitRoomList()
}

// PIPELINE #1 stream receiver
// this function is to make the stream receiving SELECTABLE
func (s *server) recv(stream pb.RoomService_StreamServer, sess_die chan struct{}) chan *pb.Room_Frame {
	ch := make(chan *pb.Room_Frame, 1)
	go func() {
		defer func() {
			close(ch)
		}()
		for {
			in, err := stream.Recv()
			if err == io.EOF { // client closed
				return
			}

			if err != nil {
				log.Error("stream.Recv() ", err)
				return
			}
			select {
			case ch <- in:
			case <-sess_die:
				return
			}
		}
	}()
	return ch
}

// PIPELINE #2 stream processing
// the center of room logic
func (s *server) Stream(stream pb.RoomService_StreamServer) error {
	defer helper.PrintPanicStack()

	sess_die := make(chan struct{})
	ch_agent := s.recv(stream, sess_die)

	// read metadata from context
	md, ok := metadata.FromContext(stream.Context())
	if !ok {
		log.Error("cannot read metadata from context")
		return ERROR_INCORRECT_FRAME_TYPE
	}
	// read key
	if len(md["userid"]) == 0 {
		log.Error("cannot read key:userid from metadata")
		return ERROR_INCORRECT_FRAME_TYPE
	}
	// parse userid
	userid, err := strconv.Atoi(md["userid"][0])
	if err != nil {
		log.Error(err)
		return ERROR_INCORRECT_FRAME_TYPE
	}
	uniqueId := md["unique_id"][0]
	serviceId := md["service_id"][0]
	isReconnect, _ := strconv.Atoi(md["is_reconnect"][0])

	// 是否已登录
	sess := NewSession(userid)
	err = sess.Get()
	if err != nil {
		return err
	}

	// 已登录 踢出
	if sess.Status == SESSION_STATUS_LOGIN {
		err = sess.NotifyKickedOut()
		if err != nil {
			log.Error("sess.NotifyKickedOut: ", err)
			return err
		}
	}

	var player *Player

	// 断线重连
	if isReconnect == 1 {
		tmp := registry.Query(userid)
		if _player, ok := tmp.(*Player); ok {
			player = _player
			if player.Flag&define.PLAYER_DISCONNECT != 0 {
				player.Flag |= define.PLAYER_LOGIN
				player.Stream = stream
			}


		} else {
			// @todo 未找到玩家处理

			log.Debug("断线重连---未找到玩家")
			return nil
		}
	} else { // 正常登录
		// player init and register
		player = NewPlayer(userid, stream)
		registry.Register(player.Id, player)
	}

	// 保存当前登录状态
	sess.TraceId = helper.Md5(fmt.Sprintf("%d-%d", userid, time.Now().Unix()))
	sess.SrvId = serviceId
	sess.UniqueId = uniqueId
	sess.Status = SESSION_STATUS_LOGIN
	err = sess.Set()
	if err != nil {
		registry.Unregister(player.Id, player)
		log.Error("sess.Set: ", err)
		return err
	}
	// 读取踢出通知
	go sess.KickedOutLoop(player, sess_die)


	log.Debugf("玩家%d登录成功，设备号：%s", player.Id, uniqueId)

	signal.SessWg.Add(1)
	defer func() {
		close(sess_die)
		// 注销登录状态
		sess.Reset()
		sess.DelKickedOutQueue()
		log.Debugf("玩家%d登出，设备号：%s", player.Id, uniqueId)
		signal.SessWg.Done()
	}()

	// >> main message loop <<
	for {
		select {
		case frame, ok := <-ch_agent: // frames from agent
			if !ok { // EOF
				player.Disconnect()
				return nil
			}
			switch frame.Type {
			case pb.Room_Message: // the passthrough message from client->agent->room
				// locate handler by proto number
				c := int16(binary.BigEndian.Uint16(frame.Message[:2]))
				handle := client_handler.Handlers[c]
				if handle == nil {
					log.Error("service not bind:", c)
					player.Disconnect()
					return ERROR_SERVICE_NOT_BIND
				}

				// handle request
				ret := handle(player, frame.Message[2:])

				// construct frame & return message from logic
				if ret != nil {
					if err := stream.Send(&pb.Room_Frame{Type: pb.Room_Message, Message: ret}); err != nil {
						log.Error(err)
						player.Disconnect()
						return err
					}
				}

			case pb.Room_Ping:
				if err := stream.Send(&pb.Room_Frame{Type: pb.Room_Ping, Message: frame.Message}); err != nil {
					log.Error(err)
					player.Disconnect()
					return err
				}
				//log.Debugf("玩家%d pong...", player.Id)
			default:
				player.Disconnect()
				log.Error("incorrect frame type:", frame.Type)
				return ERROR_INCORRECT_FRAME_TYPE
			}
		case <-signal.SessDie:
			registry.Unregister(player.Id, player)
			log.Debugf("玩家%d, Receive signal.SessDie", player.Id)
			return nil
		}

		// session control by logic
		if player.Flag&define.PLAYER_KICKED_OUT != 0 { // logic kick out
			registry.Unregister(player.Id, player)
			player.Leave()

			if err := stream.Send(&pb.Room_Frame{
				Type: pb.Room_Kick,
				Message: packet.Pack(
					define.Code["kicked_out_ack"],
					&pb.KickedOutAck{BaseAck: &pb.BaseAck{Ret: 1}},
				),
			}); err != nil {
				log.Error(err)
				return err
			}
			log.Debugf("玩家%d被踢出.", player.Id)
			return nil
		}
	}
}

func (s *server) RoomInfo(ctx context.Context, args *pb.RoomInfoArgs) (*pb.RoomInfoRes, error) {
	return &pb.RoomInfoRes{}, nil
}
