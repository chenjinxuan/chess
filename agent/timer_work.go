package main

import (
	. "chess/agent/types"
	"chess/common/log"
)

var (
	rpmLimit int
)

func initTimer(rpm_limit int) {
	rpmLimit = rpm_limit
}

// 玩家1分钟定时器
func timer_work(sess *Session, out *Buffer) {
	defer func() {
		sess.PacketCount1Min = 0
	}()

	// 发包频率控制，太高的RPS直接踢掉
	if sess.PacketCount1Min > rpmLimit {
		sess.Flag |= SESS_KICKED_OUT
		log.Infof("RPM --- userid:%d count1m:%d total:%d", sess.UserId, sess.PacketCount1Min, sess.PacketCount)
		return
	}
}
