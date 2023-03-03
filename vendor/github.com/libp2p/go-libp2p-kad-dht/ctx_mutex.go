package dht

import (
	"context"
)

type ctxMutex chan struct{}

func newCtxMutex() ctxMutex {
	return make(ctxMutex, 1)
}

func (m ctxMutex) Lock(ctx context.Context) error {
	select {
	case m <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (m ctxMutex) Unlock() {
	select {
	case <-m:
	default:
		panic("not locked")
	}
}
