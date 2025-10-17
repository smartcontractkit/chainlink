package utils

import "sync"

// KeyedMutex allows to lock based on particular values
type KeyedMutex struct {
	mutexes sync.Map
}

// LockInt64 locks the value for read/write
func (m *KeyedMutex) LockInt64(key int64) func() {
	value, _ := m.mutexes.LoadOrStore(key, new(sync.Mutex))
	mtx := value.(*sync.Mutex)
	mtx.Lock()

	return mtx.Unlock
}
