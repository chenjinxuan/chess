package main

import (
	"errors"

	"chess/common/log"
)

import (
	pb "chess/agent/proto"
	. "chess/agent/types"
)

var (
	ERROR_STREAM_NOT_OPEN = errors.New("stream not opened yet")
)

// forward messages to room server
func forward(sess *Session, p []byte) error {
	frame := &pb.Room_Frame{
		Type:    pb.Room_Message,
		Message: p,
	}

	// check stream
	if sess.Stream == nil {
		return ERROR_STREAM_NOT_OPEN
	}

	// forward the frame to game
	if err := sess.Stream.Send(frame); err != nil {
		log.Error(err)
		return err
	}
	return nil
}
