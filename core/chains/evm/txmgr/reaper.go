package txmgr

import (
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//go:generate mockery --quiet --name ReaperConfig --output ./mocks/ --case=underscore

// ReaperConfig is the config subset used by the reaper
type ReaperConfig interface {
	EthTxReaperInterval() time.Duration
	EthTxReaperThreshold() time.Duration
	EvmFinalityDepth() uint32
}

// Reaper handles periodic database cleanup for Txm
type Reaper struct {
	db             *sqlx.DB
	config         ReaperConfig
	chainID        utils.Big
	log            logger.Logger
	latestBlockNum atomic.Int64
	trigger        chan struct{}
	chStop         chan struct{}
	chDone         chan struct{}
}

// NewReaper instantiates a new reaper object
func NewReaper(lggr logger.Logger, db *sqlx.DB, config ReaperConfig, chainID big.Int) *Reaper {
	r := &Reaper{
		db,
		config,
		*utils.NewBig(&chainID),
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
	r.log.Debugf("TxmReaper: started with age threshold %v and interval %v", r.config.EthTxReaperThreshold(), r.config.EthTxReaperInterval())
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
	ticker := time.NewTicker(utils.WithJitter(r.config.EthTxReaperInterval()))
	defer ticker.Stop()
	for {
		select {
		case <-r.chStop:
			return
		case <-ticker.C:
			r.work()
			ticker.Reset(utils.WithJitter(r.config.EthTxReaperInterval()))
		case <-r.trigger:
			r.work()
			ticker.Reset(utils.WithJitter(r.config.EthTxReaperInterval()))
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
	threshold := r.config.EthTxReaperThreshold()
	if threshold == 0 {
		r.log.Debug("TxmReaper: ETH_TX_REAPER_THRESHOLD set to 0; skipping ReapEthTxes")
		return nil
	}
	minBlockNumberToKeep := headNum - int64(r.config.EvmFinalityDepth())
	mark := time.Now()
	timeThreshold := mark.Add(-threshold)

	r.log.Debugw(fmt.Sprintf("TxmReaper: reaping old eth_txes created before %s", timeThreshold.Format(time.RFC3339)), "ageThreshold", threshold, "timeThreshold", timeThreshold, "minBlockNumberToKeep", minBlockNumberToKeep)

	// Delete old confirmed eth_txes
	// NOTE that this relies on foreign key triggers automatically removing
	// the eth_tx_attempts and eth_receipts linked to every eth_tx
	err := pg.Batch(func(_, limit uint) (count uint, err error) {
		res, err := r.db.Exec(`
WITH old_enough_receipts AS (
	SELECT tx_hash FROM eth_receipts
	WHERE block_number < $1
	ORDER BY block_number ASC, id ASC
	LIMIT $2
)
DELETE FROM eth_txes
USING old_enough_receipts, eth_tx_attempts
WHERE eth_tx_attempts.eth_tx_id = eth_txes.id
AND eth_tx_attempts.hash = old_enough_receipts.tx_hash
AND eth_txes.created_at < $3
AND eth_txes.state = 'confirmed'
AND evm_chain_id = $4`, minBlockNumberToKeep, limit, timeThreshold, r.chainID)
		if err != nil {
			return count, errors.Wrap(err, "ReapEthTxes failed to delete old confirmed eth_txes")
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return count, errors.Wrap(err, "ReapEthTxes failed to get rows affected")
		}
		return uint(rowsAffected), err
	})
	if err != nil {
		return errors.Wrap(err, "TxmReaper#reapEthTxes batch delete of confirmed eth_txes failed")
	}
	// Delete old 'fatal_error' eth_txes
	err = pg.Batch(func(_, limit uint) (count uint, err error) {
		res, err := r.db.Exec(`
DELETE FROM eth_txes
WHERE created_at < $1
AND state = 'fatal_error'
AND evm_chain_id = $2`, timeThreshold, r.chainID)
		if err != nil {
			return count, errors.Wrap(err, "ReapEthTxes failed to delete old fatally errored eth_txes")
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return count, errors.Wrap(err, "ReapEthTxes failed to get rows affected")
		}
		return uint(rowsAffected), err
	})
	if err != nil {
		return errors.Wrap(err, "TxmReaper#reapEthTxes batch delete of fatally errored eth_txes failed")
	}

	r.log.Debugf("TxmReaper: ReapEthTxes completed in %v", time.Since(mark))

	return nil
}
