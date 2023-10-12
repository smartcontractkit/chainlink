package txmgr

import (
	"fmt"
	"sync/atomic"
	"time"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// Reaper handles periodic database cleanup for Txm
type Reaper struct {
	store          txmgrtypes.TxHistoryReaper
	finalityDepth  uint32
	txConfig       txmgrtypes.ReaperTransactionsConfig
	log            logger.Logger
	latestBlockNum atomic.Int64
	trigger        chan struct{}
	chStop         chan struct{}
	chDone         chan struct{}
}

// NewReaper instantiates a new reaper object
func NewReaper(lggr logger.Logger, store txmgrtypes.TxHistoryReaper, finalityDepth uint32, txConfig txmgrtypes.ReaperTransactionsConfig) *Reaper {
	r := &Reaper{
		store,
		finalityDepth,
		txConfig,
		lggr.Named("Reaper"),
		atomic.Int64{},
		make(chan struct{}, 1),
		make(chan struct{}),
		make(chan struct{}),
	}
	r.latestBlockNum.Store(-1)
	return r
}

// Start the reaper. Should only be called once.
func (r *Reaper) Start() {
	r.log.Debugf("started with age threshold %v and interval %v", r.txConfig.ReaperThreshold(), r.txConfig.ReaperInterval())
	go r.runLoop()
}

// Stop the reaper. Should only be called once.
func (r *Reaper) Stop() {
	r.log.Debug("stopping")
	close(r.chStop)
	<-r.chDone
}

func (r *Reaper) runLoop() {
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

func (r *Reaper) work() {
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
	ctx, cancel := utils.StopChan(r.chStop).NewCtx()
	defer cancel()
	threshold := r.txConfig.ReaperThreshold()
	if threshold == 0 {
		r.log.Debug("Transactions.ReaperThreshold  set to 0; skipping ReapTxes")
		return nil
	}
	minBlockNumberToKeep := headNum - int64(r.finalityDepth)
	mark := time.Now()
	timeThreshold := mark.Add(-threshold)

	r.log.Debugw(fmt.Sprintf("reaping old txes created before %s", timeThreshold.Format(time.RFC3339)), "ageThreshold", threshold, "timeThreshold", timeThreshold, "minBlockNumberToKeep", minBlockNumberToKeep)

	if err := r.store.ReapTxHistory(ctx, minBlockNumberToKeep, timeThreshold); err != nil {
		return err
	}

	r.log.Debugf("ReapTxes completed in %v", time.Since(mark))

	return nil
}
