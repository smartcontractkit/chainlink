package logger

import (
	"slices"
	"sync"
	"time"
	"weak"
)

const registerCleanupInterval = time.Minute * 5

// Register holds weak references to children T and allows updating them.
// It periodically cleans up children that have been garbage collected.
type Register[T any] struct {
	mu          sync.RWMutex
	children    []weak.Pointer[T]
	cleanupWg   sync.WaitGroup
	stopCleanup chan struct{}
}

func (r *Register[T]) Add(child *T) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.children = append(r.children, weak.Make(child))
}

func (r *Register[T]) Update(f func(*T)) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, p := range r.children {
		cw := p.Value()
		if cw != nil {
			f(cw)
		}
	}
}

func (r *Register[T]) cleanup() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.children = slices.DeleteFunc(r.children, func(p weak.Pointer[T]) bool {
		return p.Value() == nil
	})
	r.children = slices.Clip(r.children)
}

func (r *Register[T]) startPeriodicCleanup() {
	r.cleanupWg.Go(func() {
		ticker := time.NewTicker(registerCleanupInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				r.cleanup()
			case <-r.stopCleanup:
				return
			}
		}
	})
}

func (r *Register[T]) Close() {
	close(r.stopCleanup)
	r.cleanupWg.Wait()
}
