package util

import (
	"sync"
	"time"
)

type Cleanable[T any] interface {
	ClearExpired()
}

func NewIntervalCacheCleaner[T any](interval time.Duration) *IntervalCacheCleaner[T] {
	return &IntervalCacheCleaner[T]{
		interval: interval,
		stop:     make(chan struct{}),
	}
}

type IntervalCacheCleaner[T any] struct {
	interval time.Duration
	stopper  sync.Once
	stop     chan struct{}
}

func (ic *IntervalCacheCleaner[T]) Run(c Cleanable[T]) {
	ticker := time.NewTicker(ic.interval)
	for {
		select {
		case <-ticker.C:
			c.ClearExpired()
		case <-ic.stop:
			ticker.Stop()
			return
		}
	}
}

func (ic *IntervalCacheCleaner[T]) Stop() {
	ic.stopper.Do(func() {
		close(ic.stop)
	})
}
