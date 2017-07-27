package packet

import (
	"encoding/binary"
	//"fmt"
	"github.com/golang/protobuf/proto"
)

// 加上协议号
func Pack(tos int16, msg proto.Message) []byte {
	data, _ := proto.Marshal(msg)

	cache := make([]byte, PACKET_LIMIT+2)
	binary.BigEndian.PutUint16(cache, uint16(tos))

	copy(cache[2:], data)

	//fmt.Println("conn write -----------", cache[:len(data)+2])
	return cache[:len(data)+2]
}
