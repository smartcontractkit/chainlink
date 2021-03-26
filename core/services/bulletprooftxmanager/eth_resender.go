package bulletprooftxmanager

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gorm.io/gorm"
)

// pollInterval is the maximum amount of time in addition to
// EthTxResendAfterThreshold that we will wait before resending an attempt
const defaultResenderPollInterval = 5 * time.Second

// EthResender periodically picks up transactions that have been languishing
// unconfirmed for a configured amount of time without being sent, and sends
// their highest priced attempt again. This helps to defend against geth/parity
// silently dropping txes, or txes being ejected from the mempool.
//
// Previously we relied on the bumper to do this for us implicitly but there
// can occasionally be problems with this (e.g. abnormally long block times, or
// if gas bumping is disabled)
type EthResender struct {
	db           *gorm.DB
	ethClient    eth.Client
	interval     time.Duration
	ageThreshold time.Duration

	chStop chan struct{}
	chDone chan struct{}
}

func NewEthResender(db *gorm.DB, ethClient eth.Client, pollInterval, ethTxResendAfterThreshold time.Duration) *EthResender {
	if ethTxResendAfterThreshold == 0 {
		panic("EthResender requires a non-zero threshold")
	}
	return &EthResender{
		db,
		ethClient,
		pollInterval,
		ethTxResendAfterThreshold,
		make(chan struct{}),
		make(chan struct{}),
	}
}

func (er *EthResender) Start() {
	logger.Infof("EthResender: Enabled with poll interval of %s and age threshold of %s", er.interval, er.ageThreshold)
	go er.runLoop()
}

func (er *EthResender) Stop() {
	close(er.chStop)
	<-er.chDone
}

func (er *EthResender) runLoop() {
	defer close(er.chDone)

	if err := er.resendUnconfirmed(); err != nil {
		logger.Warnw("EthResender: failed to resend unconfirmed transactions", "err", err)
	}

	ticker := time.NewTicker(utils.WithJitter(er.interval))
	defer ticker.Stop()
	for {
		select {
		case <-er.chStop:
			return
		case <-ticker.C:
			if err := er.resendUnconfirmed(); err != nil {
				logger.Warnw("EthResender: failed to resend unconfirmed transactions", "err", err)
			}
		}
	}
}

func (er *EthResender) resendUnconfirmed() error {
	olderThan := time.Now().Add(-er.ageThreshold)
	attempts, err := FindEthTxesRequiringResend(er.db, olderThan)
	if err != nil {
		return errors.Wrap(err, "failed to findEthTxAttemptsRequiringReceiptFetch")
	}

	if len(attempts) == 0 {
		return nil
	}

	logger.Debugw(fmt.Sprintf("EthResender: re-sending %d transactions that were last sent over %s ago", len(attempts), er.ageThreshold), "n", len(attempts))

	var reqs []rpc.BatchElem
	for _, attempt := range attempts {
		req := rpc.BatchElem{
			Method: "eth_sendRawTransaction",
			Args:   []interface{}{hexutil.Encode(attempt.SignedRawTx)},
			Result: &common.Hash{},
		}

		reqs = append(reqs, req)
	}

	now := time.Now()
	if err := er.ethClient.RoundRobinBatchCallContext(context.Background(), reqs); err != nil {
		return errors.Wrap(err, "failed to re-send transactions")
	}

	var succeeded []int64
	for i, req := range reqs {
		if req.Error == nil {
			succeeded = append(succeeded, attempts[i].EthTxID)
		}
	}

	if err := er.updateBroadcastAts(now, succeeded); err != nil {
		return errors.Wrap(err, "failed to update last succeeded on attempts")
	}
	nSuccess := len(succeeded)
	nErrored := len(attempts) - nSuccess

	logger.Debugw("EthResender: completed", "nSuccess", nSuccess, "nErrored", nErrored)

	return nil
}

// FindEthTxesRequiringResend returns the highest priced attempt for each
// eth_tx that was last sent before or at the given time
func FindEthTxesRequiringResend(db *gorm.DB, olderThan time.Time) (attempts []models.EthTxAttempt, err error) {
	err = db.Raw(`
SELECT DISTINCT ON (eth_tx_id) eth_tx_attempts.*
FROM eth_tx_attempts
JOIN eth_txes ON eth_txes.id = eth_tx_attempts.eth_tx_id AND eth_txes.state IN ('unconfirmed', 'confirmed_missing_receipt')
WHERE eth_tx_attempts.state <> 'in_progress' AND eth_txes.broadcast_at <= ?
ORDER BY eth_tx_attempts.eth_tx_id ASC, eth_txes.nonce ASC, eth_tx_attempts.gas_price DESC
`, olderThan).
		Find(&attempts).Error

	return
}

func (er *EthResender) updateBroadcastAts(now time.Time, etxIDs []int64) error {
	// Deliberately do nothing on NULL broadcast_at because that indicates the
	// tx has been moved into a state where broadcast_at is not relevant, e.g.
	// fatally errored.
	//
	// Since we may have raced with the EthConfirmer (totally OK since highest
	// priced transaction always wins) we only want to update broadcast_at if
	// our version is later.
	return er.db.Exec(`UPDATE eth_txes SET broadcast_at = ? WHERE id = ANY(?) AND broadcast_at < ?`, now, pq.Array(etxIDs), now).Error
}
