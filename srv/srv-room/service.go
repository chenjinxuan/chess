package main

import (
	"chess/common/helper"
	"chess/common/log"
	"chess/srv/srv-room/client_handler"
	pb "chess/srv/srv-room/proto"
	"chess/srv/srv-room/registry"
	. "chess/srv/srv-room/texas_holdem"
	"encoding/binary"
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"io"
	"strconv"
)

const (
	SESS_KICKED_OUT = 0x1 // 踢掉

	SERVICE_NAME        = "room"
	DEFAULT_CH_IPC_SIZE = 16 // 默认玩家异步IPC消息队列大小
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
				log.Error(err)
				return
			}
			select {
			case ch <- in:
			case <-sess_die:
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
	ch_ipc := make(chan *pb.Room_Frame, DEFAULT_CH_IPC_SIZE)

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

	// todo get user info from mysql

	// player init
	player := NewPlayer(userid, ch_ipc)
	player.Chips = 10000

	// register user
	registry.Register(player.Id, ch_ipc)
	log.Debug("userid:", player.Id, " logged in")

	defer func() {
		registry.Unregister(player.Id, ch_ipc)
		close(sess_die)
		player.Leave()
		log.Debug("stream end userid:", player.Id)
	}()

	// >> main message loop <<
	for {
		select {
		case frame, ok := <-ch_agent: // frames from agent
			if !ok { // EOF
				return nil
			}
			switch frame.Type {
			case pb.Room_Message: // the passthrough message from client->agent->room
				// locate handler by proto number
				c := int16(binary.BigEndian.Uint16(frame.Message[:2]))
				handle := client_handler.Handlers[c]
				if handle == nil {
					log.Error("service not bind:", c)
					return ERROR_SERVICE_NOT_BIND

				}

				// handle request
				ret := handle(player, frame.Message[2:])

				// construct frame & return message from logic
				if ret != nil {
					if err := stream.Send(&pb.Room_Frame{Type: pb.Room_Message, Message: ret}); err != nil {
						log.Error(err)
						return err
					}
				}

				// session control by logic
				if player.Flag&SESS_KICKED_OUT != 0 { // logic kick out
					if err := stream.Send(&pb.Room_Frame{Type: pb.Room_Kick}); err != nil {
						log.Error(err)
						return err
					}
					return nil
				}
			case pb.Room_Ping:
				if err := stream.Send(&pb.Room_Frame{Type: pb.Room_Ping, Message: frame.Message}); err != nil {
					log.Error(err)
					return err
				}
				log.Debug("pong")
			default:
				log.Error("incorrect frame type:", frame.Type)
				return ERROR_INCORRECT_FRAME_TYPE
			}
		case frame := <-ch_ipc: // forward async messages from interprocess(goroutines) communication
			if err := stream.Send(frame); err != nil {
				log.Error(err)
				return err
			}
		}
	}
}

func (s *server) RoomInfo(ctx context.Context, args *pb.RoomInfoArgs) (*pb.RoomInfoRes, error) {
	return &pb.RoomInfoRes{}, nil
}
