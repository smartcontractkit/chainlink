package bulletprooftxmanager

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gorm.io/gorm"
)

//go:generate mockery --name ReaperConfig --output ./mocks/ --case=underscore

// ReaperConfig is the config subset used by the reaper
type ReaperConfig interface {
	EthTxReaperThreshold() time.Duration
	EthFinalityDepth() uint
}

// Reaper handles periodic database cleanup for BPTXM
type Reaper struct {
	db        *gorm.DB
	config    ReaperConfig
	log       *logger.Logger
	blockNums utils.Mailbox

	chStop chan struct{}
	chDone chan struct{}
}

// NewReaper instantiates a new reaper object
func NewReaper(db *gorm.DB, config ReaperConfig) *Reaper {
	return &Reaper{
		db,
		config,
		logger.CreateLogger(logger.Default.With("id", "bptxm_reaper")),
		*utils.NewMailbox(1),
		make(chan struct{}),
		make(chan struct{}),
	}
}

// Start the reaper. Should only be called once.
func (r *Reaper) Start() {
	r.log.Debugf("EthTxReaper: started with age threshold %v", r.config.EthTxReaperThreshold())
	go r.runLoop()
}

// Stop the reaper. Should only be called once.
func (r *Reaper) Stop() {
	r.log.Debug("EthTxReaper: stopping")
	close(r.chStop)
	<-r.chDone
}

func (r *Reaper) runLoop() {
	defer close(r.chDone)
	for {
		select {
		case <-r.chStop:
			return
		case <-r.blockNums.Notify():
			n, exists := r.blockNums.Retrieve()
			if !exists {
				continue
			}
			latestBlockNum := n.(int64)
			err := r.ReapEthTxes(latestBlockNum)
			if err != nil {
				r.log.Error("Reaper: unable to reap old eth_txes: ", err)
			}

		}
	}
}

// Deliver should be called on every new highest block number
func (r *Reaper) Deliver(latestBlockNum int64) {
	r.blockNums.Deliver(latestBlockNum)
}

// ReapEthTxes deletes old eth_txes
func (r *Reaper) ReapEthTxes(headNum int64) error {
	threshold := r.config.EthTxReaperThreshold()
	if threshold == 0 {
		r.log.Debug("Reaper: ETH_TX_REAPER_THRESHOLD set to 0; skipping reapEthTxes")
		return nil
	}
	minBlockNumberToKeep := headNum - int64(r.config.EthFinalityDepth())
	mark := time.Now()
	timeThreshold := mark.Add(-threshold)

	r.log.Infow(fmt.Sprintf("Reaper: reaping old eth_txes created before %s", timeThreshold.Format(time.RFC3339)), "timeThreshold", timeThreshold, "minBlockNumberToKeep", minBlockNumberToKeep)

	// Delete old confirmed eth_txes
	// NOTE that this relies on foreign key triggers automatically removing
	// the eth_tx_attempts and eth_receipts linked to every eth_tx
	err := orm.Batch(func(_, limit uint) (count uint, err error) {
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
AND eth_txes.state = 'confirmed'`, minBlockNumberToKeep, limit, timeThreshold)
		if res.Error != nil {
			return count, res.Error
		}
		return uint(res.RowsAffected), res.Error
	})
	if err != nil {
		return errors.Wrap(err, "Reaper#reapEthTxes batch delete of confirmed eth_txes failed")
	}
	// Delete old 'fatal_error' eth_txes
	err = orm.Batch(func(_, limit uint) (count uint, err error) {
		res := r.db.Exec(`
DELETE FROM eth_txes
WHERE created_at < ?
AND state = 'fatal_error'`, timeThreshold)
		if res.Error != nil {
			return count, res.Error
		}
		return uint(res.RowsAffected), res.Error
	})
	if err != nil {
		return errors.Wrap(err, "Reaper#reapEthTxes batch delete of fatally errored eth_txes failed")
	}

	r.log.Debugf("Reaper: completed in %v", time.Since(mark))

	return nil
}
