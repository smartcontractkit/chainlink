package bulletprooftxmanager

import (
	"fmt"
	"math/big"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/utils"
	"go.uber.org/atomic"
	"gorm.io/gorm"
)

//go:generate mockery --name ReaperConfig --output ./mocks/ --case=underscore

// ReaperConfig is the config subset used by the reaper
type ReaperConfig interface {
	EthTxReaperInterval() time.Duration
	EthTxReaperThreshold() time.Duration
	EvmFinalityDepth() uint32
}

// Reaper handles periodic database cleanup for BPTXM
type Reaper struct {
	db             *gorm.DB
	config         ReaperConfig
	chainID        utils.Big
	log            logger.Logger
	latestBlockNum *atomic.Int64
	trigger        chan struct{}
	chStop         chan struct{}
	chDone         chan struct{}
}

// NewReaper instantiates a new reaper object
func NewReaper(lggr logger.Logger, db *gorm.DB, config ReaperConfig, chainID big.Int) *Reaper {
	return &Reaper{
		db,
		config,
		*utils.NewBig(&chainID),
		lggr.With("id", "bptxm_reaper"),
		atomic.NewInt64(-1),
		make(chan struct{}, 1),
		make(chan struct{}),
		make(chan struct{}),
	}
}

// Start the reaper. Should only be called once.
func (r *Reaper) Start() {
	r.log.Debugf("BPTXMReaper: started with age threshold %v and interval %v", r.config.EthTxReaperThreshold(), r.config.EthTxReaperInterval())
	go r.runLoop()
}

// Stop the reaper. Should only be called once.
func (r *Reaper) Stop() {
	r.log.Debug("BPTXMReaper: stopping")
	close(r.chStop)
	<-r.chDone
}

func (r *Reaper) runLoop() {
	defer close(r.chDone)
	ticker := time.NewTicker(r.config.EthTxReaperInterval())
	defer ticker.Stop()
	for {
		select {
		case <-r.chStop:
			return
		case <-ticker.C:
			r.work()
		case <-r.trigger:
			r.work()
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
		r.log.Error("BPTXMReaper: unable to reap old eth_txes: ", err)
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
	threshold := r.config.EthTxReaperThreshold()
	if threshold == 0 {
		r.log.Debug("BPTXMReaper: ETH_TX_REAPER_THRESHOLD set to 0; skipping ReapEthTxes")
		return nil
	}
	minBlockNumberToKeep := headNum - int64(r.config.EvmFinalityDepth())
	mark := time.Now()
	timeThreshold := mark.Add(-threshold)

	r.log.Debugw(fmt.Sprintf("BPTXMReaper: reaping old eth_txes created before %s", timeThreshold.Format(time.RFC3339)), "ageThreshold", threshold, "timeThreshold", timeThreshold, "minBlockNumberToKeep", minBlockNumberToKeep)

	// Delete old confirmed eth_txes
	// NOTE that this relies on foreign key triggers automatically removing
	// the eth_tx_attempts and eth_receipts linked to every eth_tx
	err := postgres.Batch(func(_, limit uint) (count uint, err error) {
		res := r.db.Exec(`
WITH old_enough_receipts AS (
	SELECT tx_hash FROM eth_receipts
	WHERE block_number < ?
	ORDER BY block_number ASC, id ASC
	LIMIT ?
)
DELETE FROM eth_txes
USING old_enough_receipts, eth_tx_attempts
WHERE eth_tx_attempts.eth_tx_id = eth_txes.id
AND eth_tx_attempts.hash = old_enough_receipts.tx_hash
AND eth_txes.created_at < ?
AND eth_txes.state = 'confirmed'
AND evm_chain_id = ?`, minBlockNumberToKeep, limit, timeThreshold, r.chainID)
		if res.Error != nil {
			return count, res.Error
		}
		return uint(res.RowsAffected), res.Error
	})
	if err != nil {
		return errors.Wrap(err, "BPTXMReaper#reapEthTxes batch delete of confirmed eth_txes failed")
	}
	// Delete old 'fatal_error' eth_txes
	err = postgres.Batch(func(_, limit uint) (count uint, err error) {
		res := r.db.Exec(`
DELETE FROM eth_txes
WHERE created_at < ?
AND state = 'fatal_error'
AND evm_chain_id = ?`, timeThreshold, r.chainID)
		if res.Error != nil {
			return count, res.Error
		}
		return uint(res.RowsAffected), res.Error
	})
	if err != nil {
		return errors.Wrap(err, "BPTXMReaper#reapEthTxes batch delete of fatally errored eth_txes failed")
	}

	r.log.Debugf("BPTXMReaper: ReapEthTxes completed in %v", time.Since(mark))

	return nil
}
