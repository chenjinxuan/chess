package signal

import (
	"chess/common/log"
	"chess/common/helper"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	Wg sync.WaitGroup
	// server close signal
	Die = make(chan struct{})
)

// handle unix signals
func Handler() {
	defer helper.PrintPanicStack()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM)

	for {
		msg := <-ch
		switch msg {
		case syscall.SIGTERM: // 关闭agent
			close(Die)
			log.Info("sigterm received")
			log.Info("waiting for agents close, please wait...")
			Wg.Wait()
			log.Info("agent shutdown.")
			os.Exit(0)
		}
	}
}
