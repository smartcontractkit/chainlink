package soltxm

import (
	"context"
	"errors"
	"sync"

	"github.com/gagliardetto/solana-go"
	"golang.org/x/exp/maps"
)

type pendingTxContext struct {
	cancelBy map[solana.Signature]context.CancelFunc
	lock     sync.RWMutex
}

func newPendingTxContext() *pendingTxContext {
	return &pendingTxContext{
		cancelBy: map[solana.Signature]context.CancelFunc{},
	}
}

func (c *pendingTxContext) Insert(sig solana.Signature, cancel context.CancelFunc) error {
	// already exists
	c.lock.RLock()
	if c.cancelBy[sig] != nil {
		c.lock.RUnlock()
		return errors.New("signature already exists")
	}
	c.lock.RUnlock()

	// save cancel func
	c.lock.Lock()
	c.cancelBy[sig] = cancel
	c.lock.Unlock()
	return nil
}

func (c *pendingTxContext) Cancel(sig solana.Signature) {
	// already cancelled
	c.lock.RLock()
	if c.cancelBy[sig] == nil {
		c.lock.RUnlock()
		return
	}
	c.lock.RUnlock()

	// call cancel func + remove from map
	c.lock.Lock()
	c.cancelBy[sig]() // cancel context
	delete(c.cancelBy, sig)
	c.lock.Unlock()
	return
}

func (c *pendingTxContext) FetchAndUpdateInflight() []solana.Signature {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return maps.Keys(c.cancelBy)
}

type pendingTxContextWithProm struct {
	pendingTx *pendingTxContext
	chainID   string
}

func newPendingTxContextWithProm(id string) *pendingTxContextWithProm {
	return &pendingTxContextWithProm{
		chainID:   id,
		pendingTx: newPendingTxContext(),
	}
}

func (c *pendingTxContextWithProm) Insert(sig solana.Signature, cancel context.CancelFunc) error {
	return c.pendingTx.Insert(sig, cancel)
}

// Success - tx included in block and confirmed
func (c *pendingTxContextWithProm) Success(sig solana.Signature) {
	promSolTxmSuccessfulTxs.WithLabelValues(c.chainID).Add(1)
	c.pendingTx.Cancel(sig)
	return
}

// Revert - tx included in on chain but failed execution or simulation indicates will fail execution
func (c *pendingTxContextWithProm) Revert(sig solana.Signature) {
	promSolTxmRevertedTxs.WithLabelValues(c.chainID).Add(1)
	c.pendingTx.Cancel(sig)
	return
}

// Failed - tx failed sending to chain or failed simulation with unexpected reason
func (c *pendingTxContextWithProm) Failed(sig solana.Signature) {
	promSolTxmFailedTxs.WithLabelValues(c.chainID).Add(1)
	c.pendingTx.Cancel(sig)
	return
}

// Cancel - tx retry timed out, was not picked up by the network and confirmed in time
func (c *pendingTxContextWithProm) Cancel(sig solana.Signature) {
	promSolTxmTimedOutTxs.WithLabelValues(c.chainID).Add(1)
	c.pendingTx.Cancel(sig)
	return
}

func (c *pendingTxContextWithProm) FetchAndUpdateInflight() []solana.Signature {
	sigs := c.pendingTx.FetchAndUpdateInflight()
	promSolTxmInflightTxs.WithLabelValues(c.chainID).Set(float64(len(sigs)))
	return sigs
}
