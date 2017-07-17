package client_handler

import (
	"chess/srv/srv-room/misc/packet"
	pb "chess/srv/srv-room/proto"
	. "chess/srv/srv-room/types"
	"github.com/golang/protobuf/proto"
)

//----------------------------------- ping
func P_room_ping_req(sess *Session, data []byte) []byte {
	tbl := &pb.AutoId{}
	proto.Unmarshal(data[2:], tbl)
	return packet.Pack(Code["room_ping_ack"], tbl)
}
