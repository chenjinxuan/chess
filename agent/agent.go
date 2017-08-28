package main

import (
	pb "chess/agent/proto"
	. "chess/agent/types"
	"chess/agent/utils"
	"chess/common/define"
	"chess/common/log"
	"chess/common/services"
	"golang.org/x/net/context"
	"time"
)

// PIPELINE #2: agent
// all the packets from handleClient() will be handled
func agent(sess *Session, in chan []byte, out *Buffer) {
	defer wg.Done() // will decrease waitgroup by one, useful for manual server shutdown
	defer utils.PrintPanicStack()

	// init session
	sess.MQ = make(chan pb.Room_Frame, 512)
	sess.ConnectTime = time.Now()
	sess.LastPacketTime = time.Now()
	// minute timer
	min_timer := time.After(time.Minute)

	// cleanup work
	defer func() {
		close(sess.Die)
		if sess.Stream != nil {
			sess.Stream.CloseSend()
		}
	}()

	// >> the main message loop <<
	// handles 4 types of message:
	//  1. from client
	//  2. from game service
	//  3. timer
	//  4. server shutdown signal
	for {
		select {
		case msg, ok := <-in: // packet from network
			if !ok {
				return
			}

			sess.PacketCount++
			sess.PacketCount1Min++
			sess.PacketTime = time.Now()

			if result := route(sess, msg); result != nil {
				out.send(sess, result)
			}
			sess.LastPacketTime = sess.PacketTime
		case frame := <-sess.MQ: // packets from room
			switch frame.Type {
			case pb.Room_Message:
				out.send(sess, frame.Message)
			case pb.Room_Kick: // 登录踢出
				out.send(sess, frame.Message)
				// 拉黑token
				authConn, authServiceId := services.GetService2(define.SRV_NAME_AUTH)
				if authConn == nil {
					log.Error("cannot get auth service:", authServiceId)
				} else {
					authCli := pb.NewAuthServiceClient(authConn)
					_, err := authCli.BlackToken(context.Background(), &pb.BlackTokenArgs{Code: define.AuthKickedOut, Token: sess.Token})
					if err != nil {
						log.Error("authCli.Auth: ", err)
					}
				}
				sess.Flag |= SESS_KICKED_OUT
			}
		case <-min_timer: // minutes timer
			timer_work(sess, out)
			min_timer = time.After(time.Minute)
		case <-die: // server is shuting down...
			sess.Flag |= SESS_KICKED_OUT
		}

		// see if the player should be kicked out.
		if sess.Flag&SESS_KICKED_OUT != 0 {
			log.Debug("player should be kicked out.")
			return
		}
	}
}
