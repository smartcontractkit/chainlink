package utils

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

var (
	_ ThreadControl = &threadControl{}

	ErrThreadLimitReached = errors.New("thread limit reached")
)

// ThreadControl is a helper for managing a group of goroutines.
type ThreadControl interface {
	// Go starts a goroutine and tracks the lifetime of the goroutine.
	Go(fn func(context.Context)) error
	// Close cancels the context and waits for all tracked goroutines to exit.
	Close()
}

func NewThreadControl(pctx context.Context, limit int) *threadControl {
	ctx, cancel := context.WithCancel(pctx)
	tc := &threadControl{
		ctx:    ctx,
		cancel: cancel,
	}
	tc.limit.Store(int32(limit))

	return tc
}

type threadControl struct {
	threadsWG sync.WaitGroup

	limit   atomic.Int32
	running atomic.Int32

	ctx    context.Context
	cancel context.CancelFunc
}

func (tc *threadControl) Go(fn func(context.Context)) error {
	if err := tc.add(); err != nil {
		return err
	}
	go func(ctx context.Context) {
		defer tc.done()
		fn(ctx)
	}(tc.ctx)

	return nil
}

func (tc *threadControl) Close() {
	tc.cancel()
	tc.threadsWG.Wait()
}

func (tc *threadControl) add() error {
	if tc.running.Add(1) > tc.limit.Load() {
		tc.running.Add(-1)
		return ErrThreadLimitReached
	}
	tc.threadsWG.Add(1)
	return nil
}

func (tc *threadControl) done() {
	tc.running.Add(-1)
	tc.threadsWG.Done()
}
