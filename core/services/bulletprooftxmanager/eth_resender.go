package bulletprooftxmanager

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/null"
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
	db        *gorm.DB
	ethClient eth.Client
	interval  time.Duration
	config    Config

	chStop chan struct{}
	chDone chan struct{}
}

func NewEthResender(db *gorm.DB, ethClient eth.Client, pollInterval time.Duration, config Config) *EthResender {
	if config.EthTxResendAfterThreshold() == 0 {
		panic("EthResender requires a non-zero threshold")
	}
	return &EthResender{
		db,
		ethClient,
		pollInterval,
		config,
		make(chan struct{}),
		make(chan struct{}),
	}
}

func (er *EthResender) Start() {
	logger.Infof("EthResender: Enabled with poll interval of %s and age threshold of %s", er.interval, er.config.EthTxResendAfterThreshold())
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
	ageThreshold := er.config.EthTxResendAfterThreshold()
	maxInFlightTransactions := er.config.EthMaxInFlightTransactions()

	olderThan := time.Now().Add(-ageThreshold)
	attempts, err := FindEthTxesRequiringResend(er.db, olderThan, maxInFlightTransactions)
	if err != nil {
		return errors.Wrap(err, "failed to findEthTxAttemptsRequiringReceiptFetch")
	}

	if len(attempts) == 0 {
		return nil
	}

	logger.Infow(fmt.Sprintf("EthResender: re-sending %d unconfirmed transactions that were last sent over %s ago. These transactions are taking longer than usual to be mined. %s", len(attempts), ageThreshold, EthNodeConnectivityProblemLabel), "n", len(attempts))

	reqs := make([]rpc.BatchElem, len(attempts))
	ethTxIDs := make([]int64, len(attempts))
	for i, attempt := range attempts {
		ethTxIDs[i] = attempt.EthTxID
		req := rpc.BatchElem{
			Method: "eth_sendRawTransaction",
			Args:   []interface{}{hexutil.Encode(attempt.SignedRawTx)},
			Result: &common.Hash{},
		}
		reqs[i] = req
	}

	now := time.Now()
	batchSize := int(er.config.EthRPCDefaultBatchSize())
	if batchSize == 0 {
		batchSize = len(reqs)
	}
	for i := 0; i < len(reqs); i += batchSize {
		j := i + batchSize
		if j > len(reqs) {
			j = len(reqs)
		}

		logger.Debugw(fmt.Sprintf("EthResender: batch resending transactions %v thru %v", i, j))

		ctx, cancel := eth.DefaultQueryCtx()
		defer cancel()
		if err := er.ethClient.RoundRobinBatchCallContext(ctx, reqs[i:j]); err != nil {
			return errors.Wrap(err, "failed to re-send transactions")
		}

		if err := er.updateBroadcastAts(now, ethTxIDs[i:j]); err != nil {
			return errors.Wrap(err, "failed to update last succeeded on attempts")
		}
	}

	logResendResult(reqs)

	return nil
}

// FindEthTxesRequiringResend returns the highest priced attempt for each
// eth_tx that was last sent before or at the given time (up to limit)
func FindEthTxesRequiringResend(db *gorm.DB, olderThan time.Time, maxInFlightTransactions uint32) (attempts []models.EthTxAttempt, err error) {
	var limit null.Uint32
	if maxInFlightTransactions > 0 {
		limit = null.Uint32From(maxInFlightTransactions)
	}
	err = db.Raw(`
SELECT DISTINCT ON (eth_tx_id) eth_tx_attempts.*
FROM eth_tx_attempts
JOIN eth_txes ON eth_txes.id = eth_tx_attempts.eth_tx_id AND eth_txes.state IN ('unconfirmed', 'confirmed_missing_receipt')
WHERE eth_tx_attempts.state <> 'in_progress' AND eth_txes.broadcast_at <= ?
ORDER BY eth_tx_attempts.eth_tx_id ASC, eth_txes.nonce ASC, eth_tx_attempts.gas_price DESC
LIMIT ?
`, olderThan, limit).
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

func logResendResult(reqs []rpc.BatchElem) {
	var nNew int
	var nFatal int
	for _, req := range reqs {
		serr := eth.NewSendError(req.Error)
		if serr == nil {
			nNew++
		} else if serr.Fatal() {
			nFatal++
		}
	}
	logger.Debugw("EthResender: completed", "n", len(reqs), "nNew", nNew, "nFatal", nFatal)
}
