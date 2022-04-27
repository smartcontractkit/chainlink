package soltxm

import (
	"context"
	"errors"
	"sync"

	"github.com/gagliardetto/solana-go"
	"golang.org/x/exp/maps"
)

type TxProcesses struct {
	chainID  string
	inflight map[solana.Signature]context.CancelFunc
	lock     sync.RWMutex
}

func NewTxProcesses(id string) *TxProcesses {
	return &TxProcesses{
		chainID:  id,
		inflight: map[solana.Signature]context.CancelFunc{},
	}
}

func (c *TxProcesses) Insert(sig solana.Signature, cancel context.CancelFunc) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.inflight[sig] != nil {
		return errors.New("signature already exists")
	}
	c.inflight[sig] = cancel
	return nil
}

// Success - tx included in block and confirmed
func (c *TxProcesses) Success(sig solana.Signature) {
	promSolTxmSuccessfulTxs.WithLabelValues(c.chainID).Add(1)
	c.cancel(sig)
	return
}

// Revert - tx included in block but failed execution
func (c *TxProcesses) Revert(sig solana.Signature) {
	promSolTxmRevertedTxs.WithLabelValues(c.chainID).Add(1)
	c.cancel(sig)
	return
}

// Failed - tx failed sending to chain or failed simulation
func (c *TxProcesses) Failed(sig solana.Signature) {
	promSolTxmFailedTxs.WithLabelValues(c.chainID).Add(1)
	c.cancel(sig)
	return
}

// Cancel - tx retry timed out, was not picked up by the network and confirmed in time
func (c *TxProcesses) Cancel(sig solana.Signature) {
	promSolTxmTimedOutTxs.WithLabelValues(c.chainID).Add(1)
	c.cancel(sig)
	return
}

func (c *TxProcesses) cancel(sig solana.Signature) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// already cancelled
	if c.inflight[sig] == nil {
		return
	}

	c.inflight[sig]() // cancel context
	delete(c.inflight, sig)
	return
}

func (c *TxProcesses) FetchAndUpdateInflight() []solana.Signature {
	c.lock.RLock()
	sigs := maps.Keys(c.inflight)
	c.lock.RUnlock()

	promSolTxmInflightTxs.WithLabelValues(c.chainID).Set(float64(len(sigs)))
	return sigs
}

type LazyLoad[T any] struct {
	f     func() (T, error)
	state T
	lock  sync.RWMutex
	once  sync.Once
}

func NewLazyLoad[T any](f func() (T, error)) *LazyLoad[T] {
	return &LazyLoad[T]{
		f: f,
	}
}

func (l *LazyLoad[T]) Get() (out T, err error) {

	// fetch only once (or whenever cleared)
	l.lock.Lock()
	l.once.Do(func() {
		l.state, err = l.f()
	})
	l.lock.Unlock()

	// if err, clear so next get will retry
	if err != nil {
		l.Clear()
	}

	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.state, err
}

func (l *LazyLoad[T]) Clear() {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.once = sync.Once{}
}
