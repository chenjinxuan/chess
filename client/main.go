package main

import (
	"chess/agent/misc/packet"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	server := "127.0.0.1:8898"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	defer conn.Close()

	fmt.Println("connect success")

	go sender(conn)

	for {
		time.Sleep(1 * 1e9)
	}
}

type Ping struct {
	Name    string
	Message string
}

func sender(conn net.Conn) {
	for i := 0; i < 1; i++ {
		data := packet.Writer()
		// 写包序号，自增
		data.WriteU32(uint32(i + 1))
		// 写 协议号 payload
		_data := packet.Pack(int16(10), Ping{Name: "test", Message: "Ping"}, data)

		cache := make([]byte, 65535+2)

		sz := len(_data)
		fmt.Print(sz)
		binary.BigEndian.PutUint16(cache, uint16(sz))
		copy(cache[2:], _data)
		fmt.Print(len(cache))
		fmt.Print(cache[:sz+2])
		// write data
		n, err := conn.Write(cache[:sz+2])
		if err != nil {
			fmt.Printf("Error send reply data, bytes: %v reason: %v", n, err)
		}
	}
}
