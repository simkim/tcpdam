package main

import (
	"flag"
	"net"

	"github.com/simkim/tcpdam"
)

var (
	localAddr  = flag.String("l", ":9999", "local address")
	remoteAddr = flag.String("r", "127.0.0.1:80", "remote address")
)

func main() {
	flag.Parse()
	laddr, err := net.ResolveTCPAddr("tcp", *localAddr)
	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		panic(err)
	}

	dam := tcpdam.NewDam(remoteAddr)
	go dam.ListenSignal()

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			panic(err)
		}
		p := dam.NewProxy(conn)
		go p.Start()
	}
}
