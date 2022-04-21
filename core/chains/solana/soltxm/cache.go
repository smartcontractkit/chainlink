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
	chainID string
	cache   map[solana.Signature]context.CancelFunc
	lock    sync.RWMutex
}

func NewTxCache(id string) *TxCache {
	return &TxCache{
		chainID: id,
		cache:   map[solana.Signature]context.CancelFunc{},
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

// Success - tx included in block and confirmed
func (c *TxCache) Success(sig solana.Signature) {
	promSolTxmSuccessfulTxs.WithLabelValues(c.chainID).Add(1)
	c.cancel(sig)
	return
}

// Revert - tx included in block but failed execution
func (c *TxCache) Revert(sig solana.Signature) {
	promSolTxmRevertedTxs.WithLabelValues(c.chainID).Add(1)
	c.cancel(sig)
	return
}

// Failed - tx failed sending to chain or failed simulation
func (c *TxCache) Failed(sig solana.Signature) {
	promSolTxmFailedTxs.WithLabelValues(c.chainID).Add(1)
	c.cancel(sig)
	return
}

// Cancel - tx retry timed out, was not picked up by the network and confirmed in time
func (c *TxCache) Cancel(sig solana.Signature) {
	promSolTxmTimedOutTxs.WithLabelValues(c.chainID).Add(1)
	c.cancel(sig)
	return
}

func (c *TxCache) cancel(sig solana.Signature) {
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
	sigs := maps.Keys(c.cache)
	c.lock.RUnlock()

	promSolTxmInflightTxs.WithLabelValues(c.chainID).Set(float64(len(sigs)))
	return sigs
}

type ValidClient struct {
	tc     func() (solanaClient.ReaderWriter, error)
	client solanaClient.ReaderWriter
	lock   sync.Mutex
}

func NewValidClient(tc func() (solanaClient.ReaderWriter, error)) *ValidClient {
	return &ValidClient{
		tc: tc,
	}
}

// Get a new client if it doesnt already exist
func (vc *ValidClient) Get() (solanaClient.ReaderWriter, error) {
	vc.lock.Lock()
	defer vc.lock.Unlock()

	if vc.client == nil {
		client, err := vc.tc()
		if err != nil {
			return nil, err
		}
		vc.client = client
	}
	return vc.client, nil
}

// Clear the existing client
func (vc *ValidClient) Clear() {
	vc.lock.Lock()
	defer vc.lock.Unlock()
	vc.client = nil
}
