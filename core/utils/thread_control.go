package utils

import (
	"context"
	"sync"
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
	stop      StopChan
}

func (tc *threadControl) Go(fn func(context.Context)) {
	tc.threadsWG.Add(1)
	go func() {
		defer tc.threadsWG.Done()
		ctx, cancel := tc.stop.NewCtx()
		defer cancel()
		fn(ctx)
	}()
}

func (tc *threadControl) Close() {
	close(tc.stop)
	tc.threadsWG.Wait()
}
