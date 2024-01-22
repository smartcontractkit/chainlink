package txmgr

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

// Reaper handles periodic database cleanup for Txm
type Reaper[CHAIN_ID types.ID] struct {
	store          txmgrtypes.TxHistoryReaper[CHAIN_ID]
	config         txmgrtypes.ReaperChainConfig
	txConfig       txmgrtypes.ReaperTransactionsConfig
	chainID        CHAIN_ID
	log            logger.Logger
	latestBlockNum atomic.Int64
	trigger        chan struct{}
	chStop         services.StopChan
	chDone         chan struct{}
}

// NewReaper instantiates a new reaper object
func NewReaper[CHAIN_ID types.ID](lggr logger.Logger, store txmgrtypes.TxHistoryReaper[CHAIN_ID], config txmgrtypes.ReaperChainConfig, txConfig txmgrtypes.ReaperTransactionsConfig, chainID CHAIN_ID) *Reaper[CHAIN_ID] {
	r := &Reaper[CHAIN_ID]{
		store,
		config,
		txConfig,
		chainID,
		logger.Named(lggr, "Reaper"),
		atomic.Int64{},
		make(chan struct{}, 1),
		make(services.StopChan),
		make(chan struct{}),
	}
	r.latestBlockNum.Store(-1)
	return r
}

// Start the reaper. Should only be called once.
func (r *Reaper[CHAIN_ID]) Start() {
	r.log.Debugf("started with age threshold %v and interval %v", r.txConfig.ReaperThreshold(), r.txConfig.ReaperInterval())
	go r.runLoop()
}

// Stop the reaper. Should only be called once.
func (r *Reaper[CHAIN_ID]) Stop() {
	r.log.Debug("stopping")
	close(r.chStop)
	<-r.chDone
}

func (r *Reaper[CHAIN_ID]) runLoop() {
	defer close(r.chDone)
	ticker := time.NewTicker(utils.WithJitter(r.txConfig.ReaperInterval()))
	defer ticker.Stop()
	for {
		select {
		case <-r.chStop:
			return
		case <-ticker.C:
			r.work()
			ticker.Reset(utils.WithJitter(r.txConfig.ReaperInterval()))
		case <-r.trigger:
			r.work()
			ticker.Reset(utils.WithJitter(r.txConfig.ReaperInterval()))
		}
	}
}

func (r *Reaper[CHAIN_ID]) work() {
	latestBlockNum := r.latestBlockNum.Load()
	if latestBlockNum < 0 {
		return
	}
	err := r.ReapTxes(latestBlockNum)
	if err != nil {
		r.log.Error("unable to reap old txes: ", err)
	}
}

// SetLatestBlockNum should be called on every new highest block number
func (r *Reaper[CHAIN_ID]) SetLatestBlockNum(latestBlockNum int64) {
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
func (r *Reaper[CHAIN_ID]) ReapTxes(headNum int64) error {
	ctx, cancel := r.chStop.NewCtx()
	defer cancel()
	threshold := r.txConfig.ReaperThreshold()
	if threshold == 0 {
		r.log.Debug("Transactions.ReaperThreshold  set to 0; skipping ReapTxes")
		return nil
	}
	minBlockNumberToKeep := headNum - int64(r.config.FinalityDepth())
	mark := time.Now()
	timeThreshold := mark.Add(-threshold)

	r.log.Debugw(fmt.Sprintf("reaping old txes created before %s", timeThreshold.Format(time.RFC3339)), "ageThreshold", threshold, "timeThreshold", timeThreshold, "minBlockNumberToKeep", minBlockNumberToKeep)

	if err := r.store.ReapTxHistory(ctx, minBlockNumberToKeep, timeThreshold, r.chainID); err != nil {
		return err
	}

	r.log.Debugf("ReapTxes completed in %v", time.Since(mark))

	return nil
}
