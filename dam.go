package tcpdam

import (
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/op/go-logging"
)

type Dam struct {
	open            bool
	listenAddr      string
	remoteAddr      string
	raddr           *net.TCPAddr
	listener        *net.TCPListener
	parkedProxies   chan *Proxy
	flushingProxies chan bool
	close           chan bool
	quit            chan bool
	sigs            chan os.Signal
	Logger          *logging.Logger
}

func NewDam(listenAddr string, remoteAddr string, maxParked int, maxFlushing int) *Dam {
	return &Dam{
		open:            false,
		listenAddr:      listenAddr,
		remoteAddr:      remoteAddr,
		parkedProxies:   make(chan *Proxy, maxParked),
		flushingProxies: make(chan bool, maxFlushing),
		close:           make(chan bool),
		Logger:          logging.MustGetLogger("dam"),
	}
}

func (dam *Dam) Dial() (*net.TCPConn, error) {
	rconn, err := net.DialTCP("tcp", nil, dam.raddr)
	if err != nil {
		dam.Logger.Warningf("Can't connect to upstream : %s\n", err.Error())
	}
	return rconn, err
}

func (dam *Dam) Flushed(p *Proxy) {
	<-dam.flushingProxies
}

func (dam *Dam) Push(conn *net.TCPConn) {
	p := &Proxy{
		Lconn:  conn,
		Dam:    dam,
		Logger: dam.Logger,
	}
	dam.parkedProxies <- p
}

func (dam *Dam) Close() {
	if !dam.open {
		dam.Logger.Debugf("Already closed")
		return
	}
	dam.open = false
	dam.close <- true
	dam.Logger.Notice("Close dam")
}

func (dam *Dam) Open() error {
	if dam.open {
		dam.Logger.Debug("Already opened")
		return nil
	}
	dam.Logger.Debugf("Resolving %s\n", dam.remoteAddr)
	var err error
	dam.raddr, err = net.ResolveTCPAddr("tcp", dam.remoteAddr)
	if err != nil {
		dam.Logger.Warningf("Can't resolve remote addr %s\n", err.Error())
		return err
	}

	dam.Logger.Noticef("Open dam (%d parked)", len(dam.parkedProxies))
	dam.open = true
	go dam.Flush()
	return nil
}

func (dam *Dam) Flush() {
	dam.Logger.Debug("Flushing dam requested")

	for {
		dam.Logger.Debug("Flushing dam ...")
		select {
		case <-dam.close:
			goto end
		case p := <-dam.parkedProxies:
			dam.flushingProxies <- true
			go p.Flush()
		}
	}
end:
	dam.Logger.Debug("Flushing dam done")
}

func (dam *Dam) Start() error {
	laddr, err := net.ResolveTCPAddr("tcp", dam.listenAddr)
	if err != nil {
		dam.Logger.Errorf("Can't resolve listen address: %s", err.Error())
		return err
	}

	dam.listener, err = net.ListenTCP("tcp", laddr)
	if err != nil {
		dam.Logger.Errorf("Can't listen: %s", err.Error())
		return err
	}
	dam.quit = make(chan bool, 1)
	defer dam.StopListeningSignal()
	go dam.ListenSignal()
	for {
		delay := time.Duration(1) * time.Second
		deadline := time.Now().Add(delay)
		dam.listener.SetDeadline(deadline)
		conn, err := dam.listener.AcceptTCP()
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				select {
				case <-dam.quit:
					dam.Logger.Debug("Received quit -> return")
					dam.listener.Close()
					dam.waitEmpty()
					return nil
				default:
					continue
				}
			} else {
				return err
			}
		}
		dam.Push(conn)
	}
}

func (dam *Dam) Stop() {
	dam.Logger.Info("Stop requested")
	dam.quit <- true
}

func (dam *Dam) waitEmpty() {
	dam.Logger.Debug("Wait the dam to become empty")
	for len(dam.flushingProxies) > 0 {
		dam.Logger.Debug("Wait the dam to become empty loop")
		time.Sleep(1 * time.Second)
	}
	dam.Logger.Debug("Wait the dam to become empty loop done")
}

func (dam *Dam) ListenSignal() {
	dam.sigs = make(chan os.Signal, 1)
	signal.Notify(dam.sigs, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-dam.sigs
		switch sig {
		case syscall.SIGTERM, syscall.SIGINT:
			dam.Open()
			dam.Stop()
		case syscall.SIGUSR1:
			dam.Open()
		case syscall.SIGUSR2:
			dam.Close()
		}
	}
}

func (dam *Dam) StopListeningSignal() {
	signal.Stop(dam.sigs)
	close(dam.sigs)
}
