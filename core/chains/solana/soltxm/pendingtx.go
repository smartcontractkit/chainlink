package soltxm

import (
	"context"
	"errors"
	"sync"

	"github.com/gagliardetto/solana-go"
	"golang.org/x/exp/maps"
)

type PendingTxContext interface {
	Add(sig solana.Signature, cancel context.CancelFunc) error
	Remove(sig solana.Signature)
	ListAll() []solana.Signature
	// state change hooks
	OnSuccess(sig solana.Signature)
	OnError(sig solana.Signature, errType int) // match err type using enum
}

var _ PendingTxContext = &pendingTxContext{}

type pendingTxContext struct {
	cancelBy map[solana.Signature]context.CancelFunc
	lock     sync.RWMutex
}

func newPendingTxContext() *pendingTxContext {
	return &pendingTxContext{
		cancelBy: map[solana.Signature]context.CancelFunc{},
	}
}

func (c *pendingTxContext) Add(sig solana.Signature, cancel context.CancelFunc) error {
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

func (c *pendingTxContext) Remove(sig solana.Signature) {
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

func (c *pendingTxContext) ListAll() []solana.Signature {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return maps.Keys(c.cancelBy)
}

func (c *pendingTxContext) OnSuccess(sig solana.Signature) {
	c.Remove(sig)
	return
}

func (c *pendingTxContext) OnError(sig solana.Signature, _ int) {
	c.Remove(sig)
	return
}

var _ PendingTxContext = &pendingTxContextWithProm{}

type pendingTxContextWithProm struct {
	pendingTx *pendingTxContext
	chainID   string
}

const (
	TxFailRevert = iota
	TxFailReject
	TxFailDrop
	TxFailSimRevert
	TxFailSimOther
)

func newPendingTxContextWithProm(id string) *pendingTxContextWithProm {
	return &pendingTxContextWithProm{
		chainID:   id,
		pendingTx: newPendingTxContext(),
	}
}

func (c *pendingTxContextWithProm) Add(sig solana.Signature, cancel context.CancelFunc) error {
	return c.pendingTx.Add(sig, cancel)
}

func (c *pendingTxContextWithProm) Remove(sig solana.Signature) {
	c.pendingTx.Remove(sig)
	return
}

func (c *pendingTxContextWithProm) ListAll() []solana.Signature {
	sigs := c.pendingTx.ListAll()
	promSolTxmPendingTxs.WithLabelValues(c.chainID).Set(float64(len(sigs)))
	return sigs
}

// Success - tx included in block and confirmed
func (c *pendingTxContextWithProm) OnSuccess(sig solana.Signature) {
	promSolTxmSuccessTxs.WithLabelValues(c.chainID).Add(1)
	c.pendingTx.OnSuccess(sig)
	return
}

func (c *pendingTxContextWithProm) OnError(sig solana.Signature, errType int) {
	switch errType {
	case TxFailRevert:
		promSolTxmRevertTxs.WithLabelValues(c.chainID).Add(1)
	case TxFailReject:
		promSolTxmRejectTxs.WithLabelValues(c.chainID).Add(1)
	case TxFailDrop:
		promSolTxmDropTxs.WithLabelValues(c.chainID).Add(1)
	case TxFailSimRevert:
		promSolTxmSimRevertTxs.WithLabelValues(c.chainID).Add(1)
	case TxFailSimOther:
		promSolTxmSimOtherTxs.WithLabelValues(c.chainID).Add(1)
	}
	// increment total errors
	promSolTxmErrorTxs.WithLabelValues(c.chainID).Add(1)
	c.pendingTx.OnError(sig, errType)
}
