package tcpdam

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Dam struct {
	MaxParkedProxies  int
	open              bool
	remoteAddr        *string
	raddr             *net.TCPAddr
	parkedProxies     []*Proxy
	parkedProxiesLock *sync.Mutex
	parkedProxiesCond *sync.Cond
}

func NewDam(remoteAddr *string) *Dam {
	mutex := &sync.Mutex{}
	return &Dam{
		MaxParkedProxies:  0,
		open:              false,
		remoteAddr:        remoteAddr,
		parkedProxies:     make([]*Proxy, 0),
		parkedProxiesLock: mutex,
		parkedProxiesCond: sync.NewCond(mutex),
	}
}

func (dam *Dam) Dial() (*net.TCPConn, error) {
	rconn, err := net.DialTCP("tcp", nil, dam.raddr)
	if err != nil {
		fmt.Printf("Can't connect to upstream : %s\n", err.Error())
	}
	return rconn, err
}

func (dam *Dam) Push(p *Proxy) {
	p.Dam = dam
	dam.parkedProxiesLock.Lock()

	for !dam.open && dam.MaxParkedProxies != 0 && len(dam.parkedProxies) >= dam.MaxParkedProxies {
		fmt.Printf("Wait to create proxy : %d >= %d\n", len(dam.parkedProxies), dam.MaxParkedProxies)
		dam.parkedProxiesCond.Wait()
	}

	if dam.open {
		go p.Flush()
	} else {
		dam.parkedProxies = append(dam.parkedProxies, p)
	}
	dam.parkedProxiesLock.Unlock()
}

func (dam *Dam) Close() {
	dam.open = false
	fmt.Println("[X] Close dam")
}

func (dam *Dam) Open() error {

	fmt.Printf("Resolving %s\n", *dam.remoteAddr)
	var err error
	dam.raddr, err = net.ResolveTCPAddr("tcp", *dam.remoteAddr)
	if err != nil {
		fmt.Printf("Can't resolve remote addr %s\n", err.Error())
		return err
	}

	fmt.Println("[O] Open dam")
	dam.open = true
	dam.Flush()
	return nil
}

func (dam *Dam) Flush() {
	fmt.Println("Flush requested")

	dam.parkedProxiesLock.Lock()
	fmt.Println("Flushing dam")
	for _, proxy := range dam.parkedProxies {
		go proxy.Flush()
	}
	dam.parkedProxies = make([]*Proxy, 0)
	dam.parkedProxiesLock.Unlock()

	dam.parkedProxiesCond.Broadcast()
}

func (dam *Dam) ListenSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGUSR1, syscall.SIGUSR2)
	for {
		sig := <-sigs
		switch sig {
		case syscall.SIGUSR1:
			dam.Open()
			break
		case syscall.SIGUSR2:
			dam.Close()
			break
		}
	}
}
