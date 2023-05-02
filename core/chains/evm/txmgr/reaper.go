package txmgr

import (
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// Reaper handles periodic database cleanup for Txm
type Reaper struct {
	store          txmgrtypes.ReaperTxStore[*big.Int]
	config         EvmReaperConfig
	chainID        *big.Int
	log            logger.Logger
	latestBlockNum atomic.Int64
	trigger        chan struct{}
	chStop         chan struct{}
	chDone         chan struct{}
}

// NewReaper instantiates a new reaper object
func NewReaper(lggr logger.Logger, store txmgrtypes.ReaperTxStore[*big.Int], config EvmReaperConfig, chainID *big.Int) *Reaper {
	r := &Reaper{
		store,
		config,
		chainID,
		lggr.Named("txm_reaper"),
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
	r.log.Debugf("TxmReaper: started with age threshold %v and interval %v", r.config.TxReaperThreshold(), r.config.TxReaperInterval())
	go r.runLoop()
}

// Stop the reaper. Should only be called once.
func (r *Reaper) Stop() {
	r.log.Debug("TxmReaper: stopping")
	close(r.chStop)
	<-r.chDone
}

func (r *Reaper) runLoop() {
	defer close(r.chDone)
	ticker := time.NewTicker(utils.WithJitter(r.config.TxReaperInterval()))
	defer ticker.Stop()
	for {
		select {
		case <-r.chStop:
			return
		case <-ticker.C:
			r.work()
			ticker.Reset(utils.WithJitter(r.config.TxReaperInterval()))
		case <-r.trigger:
			r.work()
			ticker.Reset(utils.WithJitter(r.config.TxReaperInterval()))
		}
	}
}

func (r *Reaper) work() {
	latestBlockNum := r.latestBlockNum.Load()
	if latestBlockNum < 0 {
		return
	}
	err := r.ReapEthTxes(latestBlockNum)
	if err != nil {
		r.log.Error("TxmReaper: unable to reap old eth_txes: ", err)
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

// ReapEthTxes deletes old eth_txes
func (r *Reaper) ReapEthTxes(headNum int64) error {
	threshold := r.config.TxReaperThreshold()
	if threshold == 0 {
		r.log.Debug("TxmReaper: EVM.Transactions.ReaperThreshold  set to 0; skipping ReapEthTxes")
		return nil
	}
	minBlockNumberToKeep := headNum - int64(r.config.FinalityDepth())
	mark := time.Now()
	timeThreshold := mark.Add(-threshold)

	r.log.Debugw(fmt.Sprintf("TxmReaper: reaping old eth_txes created before %s", timeThreshold.Format(time.RFC3339)), "ageThreshold", threshold, "timeThreshold", timeThreshold, "minBlockNumberToKeep", minBlockNumberToKeep)

	if err := r.store.ReapConfirmedOrFatalErrorTxs(minBlockNumberToKeep, timeThreshold, r.chainID); err != nil {
		return err
	}

	r.log.Debugf("TxmReaper: ReapEthTxes completed in %v", time.Since(mark))

	return nil
}
