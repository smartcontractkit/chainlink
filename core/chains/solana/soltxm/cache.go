package soltxm

import (
	"context"
	"errors"
	"sync"

	"github.com/gagliardetto/solana-go"
	solanaClient "github.com/smartcontractkit/chainlink-solana/pkg/solana/client"
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

type ValidClient struct {
	tc     func() (solanaClient.ReaderWriter, error)
	client solanaClient.ReaderWriter
	lock   sync.RWMutex
}

func NewValidClient(tc func() (solanaClient.ReaderWriter, error)) *ValidClient {
	return &ValidClient{
		tc: tc,
	}
}

// Get a new client if it doesnt already exist
func (vc *ValidClient) Get() (solanaClient.ReaderWriter, error) {
	vc.lock.RLock()
	exist := vc.client != nil
	vc.lock.RUnlock()

	if !exist {
		client, err := vc.tc()
		if err != nil {
			return nil, err
		}
		vc.lock.Lock()
		vc.client = client
		vc.lock.Unlock()
	}

	vc.lock.RLock()
	defer vc.lock.RUnlock()
	return vc.client, nil
}

// Clear the existing client
func (vc *ValidClient) Clear() {
	vc.lock.Lock()
	defer vc.lock.Unlock()
	vc.client = nil
}
