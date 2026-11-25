// Package netfault simulates network failures in tests.
package netfault

import (
	"io"
	"net"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

// Proxy is a TCP hop between a client and a server.
//
// Blackhole makes it behave like a vanished host: bytes stop flowing but no socket
// is closed, so neither peer gets a FIN or RST and has to detect the failure itself.
type Proxy struct {
	ln     net.Listener
	target string

	mu         sync.Mutex
	conns      []net.Conn
	pairs      []*pair
	blackholed bool
}

type pair struct{ cut chan struct{} }

func NewProxy(t *testing.T, target string) *Proxy {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	p := &Proxy{ln: ln, target: target}
	t.Cleanup(p.close)
	go p.acceptLoop()
	return p
}

func (p *Proxy) Addr() string { return p.ln.Addr().String() }

// Blackhole silently drops traffic on every open connection and on any
// connection made while it is in effect.
func (p *Proxy) Blackhole() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.blackholed = true
	for _, pr := range p.pairs {
		pr.cutOnce()
	}
}

// Restore lets new connections through again. Connections already blackholed stay dead.
func (p *Proxy) Restore() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.blackholed = false
}

func (p *Proxy) acceptLoop() {
	for {
		client, err := p.ln.Accept()
		if err != nil {
			return
		}
		server, err := net.Dial("tcp", p.target)
		if err != nil {
			client.Close()
			continue
		}
		pr := &pair{cut: make(chan struct{})}
		p.mu.Lock()
		p.conns = append(p.conns, client, server)
		p.pairs = append(p.pairs, pr)
		if p.blackholed {
			pr.cutOnce()
		}
		p.mu.Unlock()
		go pr.forward(server, client)
		go pr.forward(client, server)
	}
}

func (pr *pair) cutOnce() {
	select {
	case <-pr.cut:
	default:
		close(pr.cut)
	}
}

// forward stops reading once cut, so the sender's buffers fill like they would towards a dead host.
func (pr *pair) forward(dst io.Writer, src io.Reader) {
	buf := make([]byte, 32<<10)
	for {
		select {
		case <-pr.cut:
			return
		default:
		}
		n, err := src.Read(buf)
		if err != nil {
			return
		}
		select {
		case <-pr.cut:
			return
		default:
		}
		if _, err := dst.Write(buf[:n]); err != nil {
			return
		}
	}
}

func (p *Proxy) close() {
	p.ln.Close()
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, c := range p.conns {
		c.Close()
	}
}
