package main

import _ "github.com/gonet2/libs/statsd-pprof"

func main() {
	arch := &Archiver{}
	arch.init()
	<-arch.stop
}
