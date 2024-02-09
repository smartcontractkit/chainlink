package txm

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

type ReaperTxStore interface {
	ReapTxs(context.Context, time.Time, evmtypes.Nonce, *big.Int) error
}

type ReaperClient interface {
	SequenceAt(context.Context, common.Address, *big.Int) (evmtypes.Nonce, error)
}

type ReaperConfig struct {
	ReaperInterval  time.Duration
	ReaperThreshold time.Duration
}

type Reaper struct {
	lggr    logger.Logger
	txStore ReaperTxStore
	client  ReaperClient
	ks      KeyStore
	config  ReaperConfig
	chainID *big.Int
	trigger chan struct{}
	chStop  services.StopChan
	chDone  chan struct{}
}

func NewReaper(lggr logger.Logger, txStore ReaperTxStore, config ReaperConfig, chainID *big.Int, client ReaperClient, ks KeyStore) *Reaper {
	r := &Reaper{
		lggr:    logger.Named(lggr, "Reaper"),
		txStore: txStore,
		client:  client,
		ks:      ks,
		config:  config,
		chainID: chainID,
		trigger: make(chan struct{}, 1),
		chStop:  make(services.StopChan),
		chDone:  make(chan struct{}),
	}
	return r
}

// Start the reaper. Should only be called once.
func (r *Reaper) Start() {
	r.lggr.Debugf("started with age threshold %v and interval %v", r.config.ReaperThreshold, r.config.ReaperInterval)
	go r.runLoop()
}

// Stop the reaper. Should only be called once.
func (r *Reaper) Stop() {
	r.lggr.Debug("stopping")
	close(r.chStop)
	<-r.chDone
}

func (r *Reaper) runLoop() {
	defer close(r.chDone)
	ticker := time.NewTicker(utils.WithJitter(r.config.ReaperInterval))
	defer ticker.Stop()
	for {
		select {
		case <-r.chStop:
			return
		case <-ticker.C:
			r.work()
			ticker.Reset(utils.WithJitter(r.config.ReaperInterval))
		case <-r.trigger:
			r.work()
			ticker.Reset(utils.WithJitter(r.config.ReaperInterval))
		}
	}
}

func (r *Reaper) work() {
	err := r.ReapTxs()
	if err != nil {
		r.lggr.Error("unable to reap old txs: ", err)
	}
}

// ReapTxs deletes old txs
func (r *Reaper) ReapTxs() error {
	ctx, cancel := r.chStop.NewCtx()
	defer cancel()
	threshold := r.config.ReaperThreshold
	if threshold == 0 {
		r.lggr.Debug("ReaperThreshold  set to 0; skipping ReapTxs")
		return nil
	}
	mark := time.Now()
	timeThreshold := mark.Add(-threshold)

	r.lggr.Debugw(fmt.Sprintf("reaping old txs created before %s", timeThreshold.Format(time.RFC3339)), "ageThreshold", threshold, "timeThreshold", timeThreshold)

	// TODO: get all addresses instead of enabled ones
	enabledAddresses, err := r.ks.EnabledAddressesForChain(r.chainID)
	if err != nil {
		return fmt.Errorf("Reaper failed getting enabled keys for chain %s: %w", r.chainID.String(), err)
	}
	for _, address := range enabledAddresses {
		nonce, err := r.client.SequenceAt(ctx, address, r.chainID)
		if err != nil {
			r.lggr.Errorw("Error occurred while fetching sequence for address. Skipping reaping.", "address", address, "err", err)
			continue
		}
		if err := r.txStore.ReapTxs(ctx, timeThreshold, nonce, r.chainID); err != nil {
			return err
		}
	}

	r.lggr.Debugf("ReapTxs completed in %v", time.Since(mark))

	return nil
}
