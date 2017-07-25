package main

import (
	"code.google.com/p/go.net/websocket"
	"encoding/binary"
	"io"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	. "chess/agent/types"
	"chess/agent/utils"

	"chess/common/consul"
	"chess/common/log"
	"github.com/xtaci/kcp-go"
	cli "gopkg.in/urfave/cli.v2"
)

var config *Config

type Config struct {
	wslisten                      string
	listen                        string
	readDeadline                  time.Duration
	sockbuf                       int
	udp_sockbuf                   int
	txqueuelen                    int
	dscp                          int
	sndwnd                        int
	rcvwnd                        int
	mtu                           int
	nodelay, interval, resend, nc int
}

func main() {

	// to catch all uncaught panic
	defer utils.PrintPanicStack()

	// open profiling
	//go http.ListenAndServe("0.0.0.0:6060", nil)
	app := &cli.App{
		Name:    "agent",
		Usage:   "a gateway for games with stream multiplexing",
		Version: "2.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "tcp-listen",
				Value: ":8888",
				Usage: "tpc listening address:port",
			},
			&cli.StringFlag{
				Name:  "ws-listen",
				Value: ":8899",
				Usage: "websocket listening address:port",
			},
			&cli.StringSliceFlag{
				Name:  "services",
				Value: cli.NewStringSlice("room"),
				Usage: "auto-discovering services",
			},
			&cli.DurationFlag{
				Name:  "read-deadline",
				Value: 15 * time.Second,
				Usage: "per connection read timeout",
			},
			&cli.IntFlag{
				Name:  "txqueuelen",
				Value: 128,
				Usage: "per connection output message queue, packet will be dropped if exceeds",
			},
			&cli.IntFlag{
				Name:  "sockbuf",
				Value: 32767,
				Usage: "per connection tcp socket buffer",
			},
			&cli.IntFlag{
				Name:  "udp-sockbuf",
				Value: 4194304,
				Usage: "UDP listener socket buffer",
			},
			&cli.IntFlag{
				Name:  "udp-sndwnd",
				Value: 32,
				Usage: "per connection UDP send window",
			},
			&cli.IntFlag{
				Name:  "udp-rcvwnd",
				Value: 32,
				Usage: "per connection UDP recv window",
			},
			&cli.IntFlag{
				Name:  "udp-mtu",
				Value: 1280,
				Usage: "MTU of UDP packets, without IP(20) + UDP(8)",
			},
			&cli.IntFlag{
				Name:  "dscp",
				Value: 46,
				Usage: "set DSCP(6bit)",
			},
			&cli.IntFlag{
				Name:  "nodelay",
				Value: 1,
				Usage: "ikcp_nodelay()",
			},
			&cli.IntFlag{
				Name:  "interval",
				Value: 20,
				Usage: "ikcp_nodelay()",
			},
			&cli.IntFlag{
				Name:  "resend",
				Value: 1,
				Usage: "ikcp_nodelay()",
			},
			&cli.IntFlag{
				Name:  "nc",
				Value: 1,
				Usage: "ikcp_nodelay()",
			},
			&cli.IntFlag{
				Name:  "rpm-limit",
				Value: 200,
				Usage: "per connection rpm limit",
			},
		},
		Action: func(c *cli.Context) error {
			log.Info("tpc-listen:", c.String("tcp-listen"))
			log.Info("ws-listen:", c.String("ws-listen"))
			log.Info("services:", c.StringSlice("services"))
			log.Info("read-deadline:", c.Duration("read-deadline"))
			log.Info("txqueuelen:", c.Int("txqueuelen"))
			log.Info("sockbuf:", c.Int("sockbuf"))
			log.Info("udp-sockbuf:", c.Int("udp-sockbuf"))
			log.Info("udp-sndwnd:", c.Int("udp-sndwnd"))
			log.Info("udp-rcvwnd:", c.Int("udp-rcvwnd"))
			log.Info("udp-mtu:", c.Int("udp-mtu"))
			log.Info("dscp:", c.Int("dscp"))
			log.Info("rpm-limit:", c.Int("rpm-limit"))
			log.Info("nodelay parameters:", c.Int("nodelay"), c.Int("interval"), c.Int("resend"), c.Int("nc"))

			//setup net param
			config = &Config{
				wslisten:     c.String("ws-listen"),
				listen:       c.String("tcp-listen"),
				readDeadline: c.Duration("read-deadline"),
				sockbuf:      c.Int("sockbuf"),
				udp_sockbuf:  c.Int("udp-sockbuf"),
				txqueuelen:   c.Int("txqueuelen"),
				dscp:         c.Int("dscp"),
				sndwnd:       c.Int("udp-sndwnd"),
				rcvwnd:       c.Int("udp-rcvwnd"),
				mtu:          c.Int("udp-mtu"),
				nodelay:      c.Int("nodelay"),
				interval:     c.Int("interval"),
				resend:       c.Int("resend"),
				nc:           c.Int("nc"),
			}

			// init consul config
			err := consul.InitConsulClientViaEnv()
			if err != nil {
				log.Error(err)
				os.Exit(-1)
			}

			// init services
			startup(c.StringSlice("services"))
			// init timer
			initTimer(c.Int("rpm-limit"))

			// listeners
			go tcpServer()
			go udpServer()
			go wsServer()

			// wait forever
			select {}
		},
	}
	app.Run(os.Args)
}

func wsServer() {
	//wsSrv := websocket.Server{}
	//wsSrv.Config.Header = http.Header{
	//	"Access-Control-Allow-Origin": {"*"},
	//}
	//wsSrv.Handler = websocket.Handler(wsHandler)
	//
	//httpServeMux := http.NewServeMux()
	//httpServeMux.HandleFunc("/web", wsSrv.ServeHTTP)
	//
	//server := &http.Server{
	//	Handler: httpServeMux,
	//	ReadTimeout: time.Duration(config.readDeadline) * time.Second,
	//	WriteTimeout: time.Duration(20) * time.Second,
	//}
	//l, err := net.Listen("tcp", config.wslisten)
	//if err != nil {
	//	checkError(err)
	//}
	//
	//log.Info("websocket listening on:", config.wslisten)
	//
	//if err := server.Serve(l); err != nil {
	//	checkError(err)
	//}

	//http.Handle("/", http.FileServer(http.Dir("."))) // <-- note this line

	http.Handle("/", websocket.Handler(wsHandler))

	if err := http.ListenAndServe(":8899", nil); err != nil {
		log.Error("ListenAndServe:", err)
	}
}

func wsHandler(ws *websocket.Conn) {
	ws.PayloadType = websocket.BinaryFrame
	handleClient(ws)
}

func tcpServer() {
	// resolve address & start listening
	tcpAddr, err := net.ResolveTCPAddr("tcp4", config.listen)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	log.Info("listening on:", listener.Addr())

	// loop accepting
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Warn("accept failed:", err)
			continue
		}
		// set socket read buffer
		conn.SetReadBuffer(config.sockbuf)
		// set socket write buffer
		conn.SetWriteBuffer(config.sockbuf)
		// start a goroutine for every incoming connection for reading
		go handleClient(conn)
	}
}

func udpServer() {
	l, err := kcp.Listen(config.listen)
	checkError(err)
	log.Info("udp listening on:", l.Addr())
	lis := l.(*kcp.Listener)

	if err := lis.SetReadBuffer(config.sockbuf); err != nil {
		log.Error("SetReadBuffer", err)
	}
	if err := lis.SetWriteBuffer(config.sockbuf); err != nil {
		log.Error("SetWriteBuffer", err)
	}
	if err := lis.SetDSCP(config.dscp); err != nil {
		log.Error("SetDSCP ", err)
	}

	// loop accepting
	for {
		conn, err := lis.AcceptKCP()
		if err != nil {
			log.Warn("accept failed:", err)
			continue
		}
		// set kcp parameters
		conn.SetWindowSize(config.sndwnd, config.rcvwnd)
		conn.SetNoDelay(config.nodelay, config.interval, config.resend, config.nc)
		conn.SetStreamMode(true)
		conn.SetMtu(config.mtu)

		// start a goroutine for every incoming connection for reading
		go handleClient(conn)
	}
}

// PIPELINE #1: handleClient
// the goroutine is used for reading incoming PACKETS
// each packet is defined as :
// | 2B size |     DATA       |
//
func handleClient(conn net.Conn) {
	defer utils.PrintPanicStack()
	defer conn.Close()
	// for reading the 2-Byte header
	header := make([]byte, 2)
	// the input channel for agent()
	in := make(chan []byte)
	defer func() {
		close(in) // session will close
	}()

	// create a new session object for the connection
	// and record it's IP address
	var sess Session
	//host, port, err := net.SplitHostPort(conn.RemoteAddr().String())
	//if err != nil {
	//	log.Error("cannot get remote address:", err)
	//	return
	//}
	//sess.IP = net.ParseIP(host)
	//log.Infof("new connection from:%v port:%v", host, port)

	// session die signal, will be triggered by agent()
	sess.Die = make(chan struct{})

	// create a write buffer
	out := new_buffer(conn, sess.Die, config.txqueuelen)
	go out.start()

	// start agent for PACKET processing
	wg.Add(1)
	go agent(&sess, in, out)

	// read loop
	for {
		// solve dead link problem:
		// physical disconnection without any communcation between client and server
		// will cause the read to block FOREVER, so a timeout is a rescue.
		conn.SetReadDeadline(time.Now().Add(config.readDeadline))

		// read 2B header
		n, err := io.ReadFull(conn, header)
		if err != nil {
			log.Warnf("read header failed, ip:%v reason:%v size:%v", sess.IP, err, n)
			return
		}
		size := binary.BigEndian.Uint16(header)

		// alloc a byte slice of the size defined in the header for reading data
		payload := make([]byte, size)
		n, err = io.ReadFull(conn, payload)
		if err != nil {
			log.Warnf("read payload failed, ip:%v reason:%v size:%v", sess.IP, err, n)
			return
		}
		log.Debug("read header:", header, " payload:", payload)

		// deliver the data to the input queue of agent()
		select {
		case in <- payload: // payload queued
		case <-sess.Die:
			log.Warnf("connection closed by logic, flag:%v ip:%v", sess.Flag, sess.IP)
			return
		}
	}
}

func checkError(err error) {
	if err != nil {
		log.Error(err)
		os.Exit(-1)
	}
}
