package tcpdam

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type Dam struct {
	remoteAddr     *string
	raddr          *net.TCPAddr
	waitingProxies []*Proxy
}

func NewDam(remoteAddr *string) *Dam {
	return &Dam{
		remoteAddr:     remoteAddr,
		waitingProxies: make([]*Proxy, 0),
	}
}

func (dam *Dam) NewProxy(lconn *net.TCPConn) *Proxy {
	p := &Proxy{
		lconn: lconn,
		dam:   dam,
	}
	dam.waitingProxies = append(dam.waitingProxies, p)
	return p
}

func (dam *Dam) Dial() (*net.TCPConn, error) {
	rconn, err := net.DialTCP("tcp", nil, dam.raddr)
	if err != nil {
		fmt.Printf("Can't connect to upstream : %s\n", err.Error())
	}
	return rconn, err
}

func (dam *Dam) Flush() {
	fmt.Println("Flush requested")

	fmt.Printf("Resolving %s\n", *dam.remoteAddr)
	var err error
	dam.raddr, err = net.ResolveTCPAddr("tcp", *dam.remoteAddr)
	if err != nil {
		fmt.Printf("Can't resolve remote addr %s\n", err.Error())
		return
	}

	fmt.Println("Flushing dam")
	for _, proxy := range dam.waitingProxies {
		go proxy.Flush()
	}
	dam.waitingProxies = make([]*Proxy, 0)
}

func (dam *Dam) ListenSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP)
	for {
		<-sigs
		dam.Flush()
	}
}
