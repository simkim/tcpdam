package main

import (
	"flag"
	"net"

	"github.com/simkim/tcpdam"
)

var (
	localAddr        = flag.String("l", ":9999", "local address")
	remoteAddr       = flag.String("r", "127.0.0.1:80", "remote address")
	maxParkedProxies = flag.Int("max-parked", 0, "maximum parked connections")
)

func main() {
	flag.Parse()
	laddr, err := net.ResolveTCPAddr("tcp", *localAddr)
	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		panic(err)
	}

	dam := tcpdam.NewDam(remoteAddr)
	dam.MaxParkedProxies = *maxParkedProxies
	go dam.ListenSignal()

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			panic(err)
		}
		p := &tcpdam.Proxy{
			Lconn: conn,
		}
		dam.Push(p)
	}
}
