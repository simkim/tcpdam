package tcpdam

import (
	"fmt"
	"io"
	"sync"
)

type Proxy struct {
	Lconn, Rconn io.ReadWriteCloser
	Dam          *Dam
}

func (p *Proxy) Flush() {
	fmt.Println("Dial connection")
	Rconn, err := p.Dam.Dial()
	if err != nil {
		panic(err)
	}
	p.Rconn = Rconn

	defer p.Lconn.Close()
	defer p.Rconn.Close()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		io.Copy(p.Lconn, p.Rconn)
		p.Lconn.Close()
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		io.Copy(p.Rconn, p.Lconn)
		p.Rconn.Close()
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("Flushing done")
}

func (p *Proxy) Start() {
	// Should slowly read and buffer Lconn to not timeout
}
