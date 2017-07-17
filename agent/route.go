package main

import (
	"chess/agent/client_handler"
	. "chess/agent/types"
	"chess/agent/utils"
	"chess/common/log"
	"encoding/binary"
	"time"
)

// route client protocol
func route(sess *Session, p []byte) []byte {
	start := time.Now()
	defer utils.PrintPanicStack(sess, p)

	log.Debug("binary.BigEndian.Uint32 before", p[:4])
	// 解密
	//if sess.Flag&SESS_ENCRYPT != 0 {
	//	sess.Decoder.XORKeyStream(p, p)
	//}

	if len(p) < 6 {
		log.Error("packet length error")
		sess.Flag |= SESS_KICKED_OUT
		return nil
	}

	// 读客户端数据包序列号(1,2,3...)
	// 客户端发送的数据包必须包含一个自增的序号，必须严格递增
	// 加密后，可避免重放攻击-REPLAY-ATTACK
	log.Debug("binary.BigEndian.Uint32", p[:4])
	seq_id := binary.BigEndian.Uint32(p[:4])

	// 数据包序列号验证
	if seq_id != sess.PacketCount {
		log.Errorf("illegal packet sequence id:%v should be:%v size:%v", seq_id, sess.PacketCount, len(p)-6)
		sess.Flag |= SESS_KICKED_OUT
		return nil
	}

	// 读协议号
	b := int16(binary.BigEndian.Uint16(p[4:6]))
	if _, ok := client_handler.RCode[b]; !ok {
		log.Error("protocol number not defined.")
		sess.Flag |= SESS_KICKED_OUT
		return nil
	}

	// 根据协议号断做服务划分
	// 协议号的划分采用分割协议区间, 用户可以自定义多个区间，用于转发到不同的后端服务
	var ret []byte
	if b > 2000 && b < 3000 {
		if err := forward(sess, p[4:]); err != nil {
			log.Errorf("service id:%v execute failed, error:%v", b, err)
			sess.Flag |= SESS_KICKED_OUT
			return nil
		}
	} else {
		if h := client_handler.Handlers[b]; h != nil {
			ret = h(sess, p[4:])
		} else {
			log.Errorf("service id:%v not bind", b)
			sess.Flag |= SESS_KICKED_OUT
			return nil
		}
	}

	elasped := time.Now().Sub(start)
	if b != 0 { // 排除心跳包日志
		log.Infof("REQ --- cost:%d api:%s code:%d", elasped, client_handler.RCode[b], b)
	}
	return ret
}
