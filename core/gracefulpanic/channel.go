package gracefulpanic

import "sync"

type signal struct {
	ch        chan struct{}
	panicOnce sync.Once
}

type Signal interface {
	Panic()
	Wait() <-chan struct{}
}

func NewSignal() Signal {
	return &signal{
		ch: make(chan struct{}),
	}
}

func (p *signal) Panic() {
	p.panicOnce.Do(func() {
		go close(p.ch)
	})
}

func (p *signal) Wait() <-chan struct{} {
	return p.ch
}
