package soltxm

import (
	"context"
	"errors"
	"sync"

	"github.com/gagliardetto/solana-go"
	solanaClient "github.com/smartcontractkit/chainlink-solana/pkg/solana/client"
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
