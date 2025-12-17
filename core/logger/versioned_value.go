package logger

import (
	"sync"
)

// VersionedValue is a simple thread-safe implementation of versioned value storage.
type VersionedValue[T any] struct {
	mu      sync.RWMutex
	value   T
	version int64
}

func (v *VersionedValue[T]) Store(p T) int64 {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.value = p
	v.version++
	return v.version
}

func (v *VersionedValue[T]) Load() (T, int64) {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.value, v.version
}
