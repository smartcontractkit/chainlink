package shutdown

import (
	"os"
	ossignal "os/signal"
	"sync"
	"syscall"
)

var HardPanic bool

func init() {
	if os.Getenv("ENABLE_HARD_PANIC") == "true" {
		HardPanic = true
	}
}

type signal struct {
	ch       chan struct{}
	stopOnce sync.Once
}

type Signal interface {
	Panic()
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

func (p *signal) Panic() {
	if HardPanic {
		panic("panic")
	}
	p.Stop()
}

func (p *signal) Wait() <-chan struct{} {
	return p.ch
}
