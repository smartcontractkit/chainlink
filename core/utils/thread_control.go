package utils

import (
	"context"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

var _ ThreadControl = &threadControl{}

// ThreadControl is a helper for managing a group of goroutines.
type ThreadControl interface {
	// Go starts a goroutine and tracks the lifetime of the goroutine.
	Go(fn func(context.Context))
	// GoCtx starts a goroutine with a given context and tracks the lifetime of the goroutine.
	GoCtx(ctx context.Context, fn func(context.Context))
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
	stop      services.StopChan
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

func (tc *threadControl) GoCtx(ctx context.Context, fn func(context.Context)) {
	tc.threadsWG.Add(1)
	go func() {
		defer tc.threadsWG.Done()
		// Create a new context that is cancelled when either parent context is cancelled or stop is closed.
		ctx2, cancel := tc.stop.Ctx(ctx)
		defer cancel()
		fn(ctx2)
	}()
}

func (tc *threadControl) Close() {
	close(tc.stop)
	tc.threadsWG.Wait()
}
