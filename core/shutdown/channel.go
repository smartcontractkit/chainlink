package shutdown

import (
	"os"
	ossignal "os/signal"
	"sync"
	"syscall"
)

type signal struct {
	ch       chan struct{}
	stopOnce sync.Once
}

type Signal interface {
	Wait() <-chan struct{}
}

func NewSignal() Signal {
	sigs := make(chan os.Signal, 1)
	ossignal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	p := &signal{
		ch: make(chan struct{}),
	}
	go func() {
		<-sigs
		p.Stop()
	}()
	return p
}

func (p *signal) Stop() {
	p.stopOnce.Do(func() {
		close(p.ch)
	})
}

func (p *signal) Wait() <-chan struct{} {
	return p.ch
}
