// 消息封包格式
// +----------------------------------------------------------------+
// | SIZE(2) | TIMESTAMP(4) | PROTO(2) | PAYLOAD(SIZE-6)            |
// +----------------------------------------------------------------+
package main

import (
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"

	pb "chess/protobuf-test/proto"
)

func main() {
	data := write()
	read(data)
}

func read(data []byte) {
	sz := binary.BigEndian.Uint16(data[:2])
	fmt.Println(sz)

	// 读数据包序号
	t := binary.BigEndian.Uint32(data[2:6])
	fmt.Println(t)

	// 读协议号
	p := binary.BigEndian.Uint16(data[6:8])
	fmt.Println(p)

	// 读消息体
	stReceive := &pb.RoomGetTableAck{}
	pData := data[8 : sz+2]
	//protobuf解码
	err := proto.Unmarshal(pData, stReceive)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(*stReceive)
	}
}

func write() []byte {
	// 写数据包序号
	tBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(tBuf, uint32(1))
	fmt.Println(tBuf)

	// 写协议号
	pBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(pBuf, uint16(10))
	fmt.Println(pBuf)

	// 写消息体
	bBuf, _ := proto.Marshal(&pb.RoomGetTableAck{
		BaseAck: &pb.BaseAck{Ret:20},
		Table: &pb.TableInfo{
			Id: "100",
			Cards: []*pb.CardInfo{&pb.CardInfo{Suit:10}},
			Players: []*pb.PlayerInfo{
				&pb.PlayerInfo{Id:10},
			},
		},
	})
	fmt.Println(bBuf)

	data := make([]byte, 65535+2)
	sz := len(bBuf) + 6
	binary.BigEndian.PutUint16(data, uint16(sz))
	copy(data[2:], tBuf)
	copy(data[6:], pBuf)
	copy(data[8:], bBuf)

	return data[:sz+2]
}
