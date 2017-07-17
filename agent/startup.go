package main

import (
	"chess/common/services"
)

func startup(names []string) {
	go sig_handler()
	services.Discover(names)
}
