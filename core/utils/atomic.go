package utils

import (
	"sync/atomic"
)

type AtomicBool struct {
	atomic.Value
}

func (b *AtomicBool) Get() bool {
	if as, is := b.Load().(bool); is {
		return as
	}
	return false
}
func (b *AtomicBool) Set(val bool) {
	b.Store(val)
}
