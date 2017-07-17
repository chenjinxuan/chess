package main

import (
	"chess/agent/utils"
	"chess/common/log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	wg sync.WaitGroup
	// server close signal
	die = make(chan struct{})
)

// handle unix signals
func sig_handler() {
	defer utils.PrintPanicStack()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM)

	for {
		msg := <-ch
		switch msg {
		case syscall.SIGTERM: // 关闭agent
			close(die)
			log.Info("sigterm received")
			log.Info("waiting for agents close, please wait...")
			wg.Wait()
			log.Info("agent shutdown.")
			os.Exit(0)
		}
	}
}
