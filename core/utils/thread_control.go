package utils

import (
	"context"
	"sync"
	"sync/atomic"
)

var _ ThreadControl = &threadControl{}

// ThreadControl is a helper for managing a group of goroutines.
type ThreadControl interface {
	// Go starts a goroutine and tracks the lifetime of the goroutine.
	Go(fn func(context.Context))
	// Close cancels the goroutines and waits for all of them to exit.
	Close()
}

func NewThreadControl() *threadControl {
	tc := &threadControl{
		stop: make(chan struct{}),
	}

	return tc
}

type threadControl struct {
	threadsWG sync.WaitGroup

	running atomic.Int32

	stop StopChan
}

func (tc *threadControl) Go(fn func(context.Context)) {
	tc.add()
	go func() {
		defer tc.done()
		ctx, cancel := tc.stop.NewCtx()
		defer cancel()
		fn(ctx)
	}()
}

func (tc *threadControl) Close() {
	close(tc.stop)
	tc.threadsWG.Wait()
}

func (tc *threadControl) add() {
	tc.running.Add(1)
	tc.threadsWG.Add(1)
}

func (tc *threadControl) done() {
	tc.running.Add(-1)
	tc.threadsWG.Done()
}
