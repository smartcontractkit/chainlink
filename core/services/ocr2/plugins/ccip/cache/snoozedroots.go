package cache

import (
	"sync"
	"time"
)

type SnoozedRoots interface {
	Get(k [32]byte) (time.Time, bool)
	Set(k [32]byte, v time.Time)
}

type SnoozedRootsInMem struct {
	mem map[[32]byte]time.Time
	mu  *sync.RWMutex
}

func NewSnoozedRootsInMem() *SnoozedRootsInMem {
	return &SnoozedRootsInMem{
		mem: make(map[[32]byte]time.Time),
		mu:  &sync.RWMutex{},
	}
}

func (c *SnoozedRootsInMem) Get(k [32]byte) (time.Time, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.mem[k]
	return v, ok
}

func (c *SnoozedRootsInMem) Set(k [32]byte, v time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.mem[k] = v
}
