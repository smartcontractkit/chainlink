package soltxm

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/google/uuid"
	"golang.org/x/exp/maps"
)

type PendingTxContext interface {
	New(sig solana.Signature, cancel context.CancelFunc) (uuid.UUID, error)
	Add(id uuid.UUID, sig solana.Signature) error
	Remove(sig solana.Signature) uuid.UUID
	ListAll() []solana.Signature
	Expired(sig solana.Signature, lifespan time.Duration) bool
	// state change hooks
	OnSuccess(sig solana.Signature) uuid.UUID
	OnError(sig solana.Signature, errType int) uuid.UUID // match err type using enum
}

var _ PendingTxContext = &pendingTxContext{}

type pendingTxContext struct {
	cancelBy  map[uuid.UUID]context.CancelFunc
	timestamp map[uuid.UUID]time.Time
	sigToId   map[solana.Signature]uuid.UUID
	idToSigs  map[uuid.UUID][]solana.Signature
	lock      sync.RWMutex
}

func newPendingTxContext() *pendingTxContext {
	return &pendingTxContext{
		cancelBy:  map[uuid.UUID]context.CancelFunc{},
		timestamp: map[uuid.UUID]time.Time{},
		sigToId:   map[solana.Signature]uuid.UUID{},
		idToSigs:  map[uuid.UUID][]solana.Signature{},
	}
}

func (c *pendingTxContext) New(sig solana.Signature, cancel context.CancelFunc) (uuid.UUID, error) {
	// validate signature does not exist
	c.lock.RLock()
	if _, exists := c.sigToId[sig]; exists {
		c.lock.RUnlock()
		return uuid.UUID{}, errors.New("signature already exists")
	}
	c.lock.RUnlock()

	// upgrade to write lock if sig does not exist
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, exists := c.sigToId[sig]; exists {
		return uuid.UUID{}, errors.New("signature already exists")
	}
	// save cancel func
	id := uuid.New()
	c.cancelBy[id] = cancel
	c.timestamp[id] = time.Now()
	c.sigToId[sig] = id
	c.idToSigs[id] = []solana.Signature{sig}
	return id, nil
}

func (c *pendingTxContext) Add(id uuid.UUID, sig solana.Signature) error {
	// already exists
	c.lock.RLock()
	if _, exists := c.sigToId[sig]; exists {
		c.lock.RUnlock()
		return errors.New("signature already exists")
	}
	if _, exists := c.idToSigs[id]; !exists {
		c.lock.RUnlock()
		return errors.New("id does not exist")
	}
	c.lock.RUnlock()

	// upgrade to write lock if sig does not exist
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, exists := c.sigToId[sig]; exists {
		return errors.New("signature already exists")
	}
	if _, exists := c.idToSigs[id]; !exists {
		return errors.New("id does not exist - tx likely confirmed by other signature")
	}
	// save signature
	c.sigToId[sig] = id
	c.idToSigs[id] = append(c.idToSigs[id], sig)
	return nil
}

// returns the id if removed (otherwise returns 0-id)
func (c *pendingTxContext) Remove(sig solana.Signature) (id uuid.UUID) {
	// check if already cancelled
	c.lock.RLock()
	id, sigExists := c.sigToId[sig]
	if !sigExists {
		c.lock.RUnlock()
		return id
	}
	if _, idExists := c.idToSigs[id]; !idExists {
		c.lock.RUnlock()
		return id
	}
	c.lock.RUnlock()

	// upgrade to write lock if sig does not exist
	c.lock.Lock()
	defer c.lock.Unlock()
	id, sigExists = c.sigToId[sig]
	if !sigExists {
		return id
	}
	sigs, idExists := c.idToSigs[id]
	if !idExists {
		return id
	}

	// call cancel func + remove from map
	c.cancelBy[id]() // cancel context
	delete(c.cancelBy, id)
	delete(c.timestamp, id)
	delete(c.idToSigs, id)
	for _, s := range sigs {
		delete(c.sigToId, s)
	}
	return id
}

func (c *pendingTxContext) ListAll() []solana.Signature {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return maps.Keys(c.sigToId)
}

// Expired returns if the timeout for trying to confirm a signature has been reached
func (c *pendingTxContext) Expired(sig solana.Signature, lifespan time.Duration) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	id, exists := c.sigToId[sig]
	if !exists {
		return false // return expired = false if timestamp does not exist (likely cleaned up by something else previously)
	}

	timestamp, exists := c.timestamp[id]
	if !exists {
		return false // return expired = false if timestamp does not exist (likely cleaned up by something else previously)
	}

	return time.Since(timestamp) > lifespan
}

func (c *pendingTxContext) OnSuccess(sig solana.Signature) uuid.UUID {
	return c.Remove(sig)
}

func (c *pendingTxContext) OnError(sig solana.Signature, _ int) uuid.UUID {
	return c.Remove(sig)
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

func (c *pendingTxContextWithProm) New(sig solana.Signature, cancel context.CancelFunc) (uuid.UUID, error) {
	return c.pendingTx.New(sig, cancel)
}

func (c *pendingTxContextWithProm) Add(id uuid.UUID, sig solana.Signature) error {
	return c.pendingTx.Add(id, sig)
}

func (c *pendingTxContextWithProm) Remove(sig solana.Signature) uuid.UUID {
	return c.pendingTx.Remove(sig)
}

func (c *pendingTxContextWithProm) ListAll() []solana.Signature {
	sigs := c.pendingTx.ListAll()
	promSolTxmPendingTxs.WithLabelValues(c.chainID).Set(float64(len(sigs)))
	return sigs
}

func (c *pendingTxContextWithProm) Expired(sig solana.Signature, lifespan time.Duration) bool {
	return c.pendingTx.Expired(sig, lifespan)
}

// Success - tx included in block and confirmed
func (c *pendingTxContextWithProm) OnSuccess(sig solana.Signature) uuid.UUID {
	id := c.pendingTx.OnSuccess(sig) // empty ID indicates already previously removed
	if id != uuid.Nil {              // increment if tx was not removed
		promSolTxmSuccessTxs.WithLabelValues(c.chainID).Add(1)
	}
	return id
}

func (c *pendingTxContextWithProm) OnError(sig solana.Signature, errType int) uuid.UUID {
	// special RPC rejects transaction (signature will not be valid)
	if errType == TxFailReject {
		promSolTxmRejectTxs.WithLabelValues(c.chainID).Add(1)
		promSolTxmErrorTxs.WithLabelValues(c.chainID).Add(1)
		return uuid.Nil
	}

	id := c.pendingTx.OnError(sig, errType) // empty ID indicates already removed
	if id != uuid.Nil {
		switch errType {
		case TxFailRevert:
			promSolTxmRevertTxs.WithLabelValues(c.chainID).Add(1)
		case TxFailDrop:
			promSolTxmDropTxs.WithLabelValues(c.chainID).Add(1)
		case TxFailSimRevert:
			promSolTxmSimRevertTxs.WithLabelValues(c.chainID).Add(1)
		case TxFailSimOther:
			promSolTxmSimOtherTxs.WithLabelValues(c.chainID).Add(1)
		}
		// increment total errors
		promSolTxmErrorTxs.WithLabelValues(c.chainID).Add(1)
	}

	return id
}
