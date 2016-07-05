package tcpdam

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Proxy struct {
	lconn, rconn io.ReadWriteCloser
	dam          *Dam
	errsig       chan bool
	erred        bool
}

func (p *Proxy) err(s string, err error) {
	if p.erred {
		return
	}
	if err != io.EOF {
		fmt.Printf("%s %s", s, err)
	}
	p.errsig <- true
	p.erred = true
}

func (p *Proxy) pipe(src, dst io.ReadWriteCloser) {
	buff := make([]byte, 0xffff)
	for {
		n, err := src.Read(buff)
		if err != nil {
			if err == io.EOF {
				dst.Close()
			}
			p.err("Can't read", err)
			return
		}
		b := buff[:n]
		n2, err := dst.Write(b)
		if err != nil {
			p.err("can't write", err)
			return
		}
		if n != n2 {
			panic(fmt.Sprintf("read buffer %d not fully written %d", n, n2))
		}
	}
}

func (p *Proxy) Flush(rconn *net.TCPConn) {
	defer p.lconn.Close()
	p.rconn = rconn
	defer p.rconn.Close()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		io.Copy(p.lconn, p.rconn)
		p.lconn.Close()
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		io.Copy(p.rconn, p.lconn)
		p.rconn.Close()
		wg.Done()
	}()
	wg.Wait()
	fmt.Println("done")
}

func (p *Proxy) Start() {
	// Should slowly read and buffer lconn to not timeout
}
