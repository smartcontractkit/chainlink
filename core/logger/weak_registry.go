package logger

import (
	"slices"
	"sync"
	"time"
	"weak"
)

const registryCleanupInterval = time.Minute * 1

// WeakRegistry is a thread-safe storage of weak references allowing updates to all live entries.
// It also periodically cleans up references that have been garbage collected.
type WeakRegistry[T any] struct {
	mu          sync.Mutex
	entries     []weak.Pointer[T]
	cleanupWg   sync.WaitGroup
	cleanupStop chan struct{}
}

func NewWeakRegistry[T any]() *WeakRegistry[T] {
	registry := &WeakRegistry[T]{cleanupStop: make(chan struct{})}
	registry.startPeriodicCleanup()
	return registry
}

func (r *WeakRegistry[T]) Add(entry *T) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries = append(r.entries, weak.Make(entry))
}

func (r *WeakRegistry[T]) Update(f func(*T)) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, p := range r.entries {
		v := p.Value()
		if v != nil {
			f(v)
		}
	}
}

func (r *WeakRegistry[T]) cleanup() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries = slices.DeleteFunc(r.entries, func(p weak.Pointer[T]) bool {
		return p.Value() == nil
	})
	r.entries = slices.Clip(r.entries)
}

func (r *WeakRegistry[T]) startPeriodicCleanup() {
	r.cleanupWg.Go(func() {
		ticker := time.NewTicker(registryCleanupInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				r.cleanup()
			case <-r.cleanupStop:
				return
			}
		}
	})
}

func (r *WeakRegistry[T]) Close() {
	close(r.cleanupStop)
	r.cleanupWg.Wait()
}
