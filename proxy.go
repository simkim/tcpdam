package tcpdam

import (
	"fmt"
	"io"
	"sync"
)

type Proxy struct {
	lconn, rconn io.ReadWriteCloser
	dam          *Dam
}

func (p *Proxy) Flush() {
	fmt.Println("Dial connection")
	rconn, err := p.dam.Dial()
	if err != nil {
		panic(err)
	}
	p.rconn = rconn

	defer p.lconn.Close()
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
	fmt.Println("Flushing done")
}

func (p *Proxy) Start() {
	// Should slowly read and buffer lconn to not timeout
}
