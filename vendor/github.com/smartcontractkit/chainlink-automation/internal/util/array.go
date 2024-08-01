package util

import "sync"

type SyncedArray[T any] struct {
	data []T
	mu   sync.RWMutex
}

func NewSyncedArray[T any]() *SyncedArray[T] {
	return &SyncedArray[T]{
		data: []T{},
	}
}

func (a *SyncedArray[T]) Append(vals ...T) *SyncedArray[T] {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.data = append(a.data, vals...)
	return a
}

func (a *SyncedArray[T]) Values() []T {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.data
}

func Unflatten[T any](b []T, size int) (groups [][]T) {
	for i := 0; i < len(b); i += size {
		j := i + size
		if j > len(b) {
			j = len(b)
		}
		groups = append(groups, b[i:j])
	}
	return
}
