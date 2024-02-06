package txm

import (
	"context"
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
)

type ReaperTxStore interface {
	ReapTxs(context.Context, time.Time, *big.Int) error
}

type ReaperConfig interface {
	ReaperInterval() time.Duration
	ReaperThreshold() time.Duration
}

type Reaper struct {
	lggr           logger.Logger
	txStore        ReaperTxStore
	config         ReaperConfig
	chainID        *big.Int
	latestBlockNum atomic.Int64
	trigger        chan struct{}
	chStop         services.StopChan
	chDone         chan struct{}
}

func NewReaper(lggr logger.Logger, txStore ReaperTxStore, config ReaperConfig, chainID *big.Int) *Reaper {
	r := &Reaper{
		lggr:           logger.Named(lggr, "Reaper"),
		txStore:        txStore,
		config:         config,
		chainID:        chainID,
		latestBlockNum: atomic.Int64{},
		trigger:        make(chan struct{}, 1),
		chStop:         make(services.StopChan),
		chDone:         make(chan struct{}),
	}
	r.latestBlockNum.Store(-1)
	return r
}

// Start the reaper. Should only be called once.
func (r *Reaper) Start() {
	r.lggr.Debugf("started with age threshold %v and interval %v", r.config.ReaperThreshold(), r.config.ReaperInterval())
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
	ticker := time.NewTicker(utils.WithJitter(r.config.ReaperInterval()))
	defer ticker.Stop()
	for {
		select {
		case <-r.chStop:
			return
		case <-ticker.C:
			r.work()
			ticker.Reset(utils.WithJitter(r.config.ReaperInterval()))
		case <-r.trigger:
			r.work()
			ticker.Reset(utils.WithJitter(r.config.ReaperInterval()))
		}
	}
}

func (r *Reaper) work() {
	latestBlockNum := r.latestBlockNum.Load()
	if latestBlockNum < 0 {
		return
	}
	err := r.ReapTxes(latestBlockNum)
	if err != nil {
		r.lggr.Error("unable to reap old txes: ", err)
	}
}

// SetLatestBlockNum should be called on every new highest block number
func (r *Reaper) SetLatestBlockNum(latestBlockNum int64) {
	if latestBlockNum < 0 {
		panic(fmt.Sprintf("latestBlockNum must be 0 or greater, got: %d", latestBlockNum))
	}
	was := r.latestBlockNum.Swap(latestBlockNum)
	if was < 0 {
		// Run reaper once on startup
		r.trigger <- struct{}{}
	}
}

// ReapTxes deletes old txes
func (r *Reaper) ReapTxes(headNum int64) error {
	ctx, cancel := r.chStop.NewCtx()
	defer cancel()
	threshold := r.config.ReaperThreshold()
	if threshold == 0 {
		r.lggr.Debug("Transactions.ReaperThreshold  set to 0; skipping ReapTxes")
		return nil
	}
	mark := time.Now()
	timeThreshold := mark.Add(-threshold)

	r.lggr.Debugw(fmt.Sprintf("reaping old txes created before %s", timeThreshold.Format(time.RFC3339)), "ageThreshold", threshold, "timeThreshold", timeThreshold)

	if err := r.txStore.ReapTxs(ctx, timeThreshold, r.chainID); err != nil {
		return err
	}

	r.lggr.Debugf("ReapTxes completed in %v", time.Since(mark))

	return nil
}
