package soltxm

import (
	"context"
	"errors"
	"sync"

	"github.com/gagliardetto/solana-go"
	"golang.org/x/exp/maps"
)

type TxCache struct {
	cache map[solana.Signature]context.CancelFunc
	lock  sync.RWMutex
}

func NewTxCache() *TxCache {
	return &TxCache{
		cache: map[solana.Signature]context.CancelFunc{},
	}
}

func (c *TxCache) Insert(sig solana.Signature, cancel context.CancelFunc) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.cache[sig] != nil {
		return errors.New("signature already exists")
	}
	c.cache[sig] = cancel
	return nil
}

func (c *TxCache) Cancel(sig solana.Signature) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// already cancelled
	if c.cache[sig] == nil {
		return
	}

	c.cache[sig]() // cancel context
	delete(c.cache, sig)
	return
}

func (c *TxCache) List() []solana.Signature {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return maps.Keys(c.cache)
}
