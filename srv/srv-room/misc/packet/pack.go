package packet

import (
	"encoding/binary"
	"github.com/golang/protobuf/proto"
)

// 加上包头 2字节 size
func FastPack(data []byte) []byte {
	cache := make([]byte, PACKET_LIMIT+2)
	sz := len(data)
	binary.BigEndian.PutUint16(cache, uint16(sz))
	copy(cache[2:], data)
	return cache[:sz+2]
}

// 加上包头 2字节 size
func Pack(tos int16, msg proto.Message) []byte {
	data, _ := proto.Marshal(msg)

	cache := make([]byte, PACKET_LIMIT+2)
	sz := len(data) + 2
	binary.BigEndian.PutUint16(cache, uint16(sz))

	tosBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(tosBuf, uint16(tos))

	copy(cache[2:4], tosBuf)
	copy(cache[4:], data)

	return cache[:sz+2]
}
