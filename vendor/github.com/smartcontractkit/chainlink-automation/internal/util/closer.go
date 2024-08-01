package util

import (
	"context"
	"sync"
)

type Closer struct {
	cancel context.CancelFunc
	lock   sync.Mutex
}

func (c *Closer) Store(cancel context.CancelFunc) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.cancel != nil {
		return false
	}
	c.cancel = cancel
	return true
}

func (c *Closer) Close() bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.cancel != nil {
		c.cancel()
		c.cancel = nil
		return true
	}
	return false
}
