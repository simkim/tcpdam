package tcpdam

import (
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/op/go-logging"
)

type Dam struct {
	MaxParkedProxies  int
	open              bool
	listenAddr        *string
	remoteAddr        *string
	raddr             *net.TCPAddr
	listener          *net.TCPListener
	parkedProxies     []*Proxy
	parkedProxiesLock *sync.Mutex
	parkedProxiesCond *sync.Cond
	shouldQuitCond    sync.Cond
	quit              chan bool
	sigs              chan os.Signal
	Logger            *logging.Logger
}

func NewDam(listenAddr *string, remoteAddr *string) *Dam {
	mutex := &sync.Mutex{}
	return &Dam{
		MaxParkedProxies:  0,
		open:              false,
		listenAddr:        listenAddr,
		remoteAddr:        remoteAddr,
		parkedProxies:     make([]*Proxy, 0),
		parkedProxiesLock: mutex,
		parkedProxiesCond: sync.NewCond(mutex),
		Logger:            logging.MustGetLogger("dam"),
		shouldQuitCond:    sync.Cond{L: &sync.Mutex{}},
	}
}

func (dam *Dam) Dial() (*net.TCPConn, error) {
	rconn, err := net.DialTCP("tcp", nil, dam.raddr)
	if err != nil {
		dam.Logger.Warningf("Can't connect to upstream : %s\n", err.Error())
	}
	return rconn, err
}

func (dam *Dam) Push(p *Proxy) {
	p.Dam = dam
	p.Logger = dam.Logger

	dam.parkedProxiesLock.Lock()

	for !dam.open && dam.MaxParkedProxies != 0 && len(dam.parkedProxies) >= dam.MaxParkedProxies {
		dam.Logger.Debugf("Too many connections, waiting free slots : %d >= %d\n", len(dam.parkedProxies), dam.MaxParkedProxies)
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
	if !dam.open {
		dam.Logger.Debugf("Already closed")
		return
	}
	dam.open = false
	dam.Logger.Notice("Close dam")
}

func (dam *Dam) Open() error {
	if dam.open {
		dam.Logger.Debugf("Already opened")
		return nil
	}
	dam.Logger.Debugf("Resolving %s\n", *dam.remoteAddr)
	var err error
	dam.raddr, err = net.ResolveTCPAddr("tcp", *dam.remoteAddr)
	if err != nil {
		dam.Logger.Warningf("Can't resolve remote addr %s\n", err.Error())
		return err
	}

	dam.Logger.Notice("Open dam")
	dam.open = true
	dam.Flush()
	return nil
}

func (dam *Dam) Flush() {
	dam.Logger.Debug("Flushing dam requested")

	dam.parkedProxiesLock.Lock()
	dam.Logger.Debug("Flushing dam")
	for _, proxy := range dam.parkedProxies {
		go proxy.Flush()
	}
	dam.parkedProxies = make([]*Proxy, 0)
	dam.parkedProxiesLock.Unlock()

	dam.parkedProxiesCond.Broadcast()
}

func (dam *Dam) Start() {
	laddr, err := net.ResolveTCPAddr("tcp", *dam.listenAddr)
	if err != nil {
		panic(err)
	}

	dam.listener, err = net.ListenTCP("tcp", laddr)
	if err != nil {
		panic(err)
	}
	dam.quit = make(chan bool, 1)
	defer dam.StopListeningSignal()
	go dam.ListenSignal()
	for {
		conn, err := dam.listener.AcceptTCP()
		if err != nil {
			select {
			case <-dam.quit:
				dam.Logger.Debug("Received quit -> return")
				return
			default:
				panic(err)
			}
		}
		p := &Proxy{
			Lconn: conn,
		}
		dam.Push(p)
	}
}

func (dam *Dam) Stop() {
	dam.Logger.Info("Stop requested")
	dam.quit <- true
	dam.listener.Close()
}

func (dam *Dam) WaitEmpty() {
}

func (dam *Dam) ListenSignal() {
	dam.sigs = make(chan os.Signal, 1)
	signal.Notify(dam.sigs, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-dam.sigs
		switch sig {
		case syscall.SIGTERM:
			dam.Open()
			dam.Stop()
		case syscall.SIGINT:
			dam.Open()
			dam.Stop()
		case syscall.SIGUSR1:
			dam.Open()
			break
		case syscall.SIGUSR2:
			dam.Close()
			break
		}
	}
}

func (dam *Dam) StopListeningSignal() {
	signal.Stop(dam.sigs)
	close(dam.sigs)
}
