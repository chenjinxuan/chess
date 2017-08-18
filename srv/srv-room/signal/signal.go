package signal

import (
	"chess/common/helper"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	TableWg sync.WaitGroup
	SessWg  sync.WaitGroup
	// server close signal
	TableDie = make(chan struct{})
	SessDie  = make(chan struct{})
)

// handle unix signals
func Handler() {
	defer helper.PrintPanicStack()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT)

	for {
		msg := <-ch
		switch msg {
		case syscall.SIGTERM, syscall.SIGINT: // 关闭room
			close(TableDie)
			fmt.Println("waiting for table close, please wait...")
			TableWg.Wait()
			fmt.Println("all tables closed.")

			close(SessDie)
			fmt.Println("waiting for session close, please wait...")
			SessWg.Wait()
			fmt.Println("all session closed.")

			fmt.Println("room shutdown.")
			os.Exit(0)
		case syscall.SIGHUP:
			return
		}
	}
}
