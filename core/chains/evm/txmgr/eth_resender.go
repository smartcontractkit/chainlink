package txmgr

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/chains/evm/label"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/utils"
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
	db        *sqlx.DB
	ethClient evmclient.Client
	ks        KeyStore
	chainID   big.Int
	interval  time.Duration
	config    Config
	logger    logger.Logger

	ctx    context.Context
	cancel context.CancelFunc
	chDone chan struct{}
}

// NewEthResender creates a new concrete EthResender
func NewEthResender(lggr logger.Logger, db *sqlx.DB, ethClient evmclient.Client, ks KeyStore, pollInterval time.Duration, config Config) *EthResender {
	if config.EthTxResendAfterThreshold() == 0 {
		panic("EthResender requires a non-zero threshold")
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &EthResender{
		db,
		ethClient,
		ks,
		*ethClient.ChainID(),
		pollInterval,
		config,
		lggr.Named("EthResender"),
		ctx,
		cancel,
		make(chan struct{}),
	}
}

// Start is a comment which satisfies the linter
func (er *EthResender) Start() {
	er.logger.Debugf("Enabled with poll interval of %s and age threshold of %s", er.interval, er.config.EthTxResendAfterThreshold())
	go er.runLoop()
}

// Stop is a comment which satisfies the linter
func (er *EthResender) Stop() {
	er.cancel()
	<-er.chDone
}

func (er *EthResender) runLoop() {
	defer close(er.chDone)

	if err := er.resendUnconfirmed(); err != nil {
		er.logger.Warnw("Failed to resend unconfirmed transactions", "err", err)
	}

	ticker := time.NewTicker(utils.WithJitter(er.interval))
	defer ticker.Stop()
	for {
		select {
		case <-er.ctx.Done():
			return
		case <-ticker.C:
			if err := er.resendUnconfirmed(); err != nil {
				er.logger.Warnw("Failed to resend unconfirmed transactions", "err", err)
			}
		}
	}
}

func (er *EthResender) resendUnconfirmed() error {
	keys, err := er.ks.EnabledKeysForChain(&er.chainID)
	if err != nil {
		return errors.Wrapf(err, "EthResender failed getting enabled keys for chain %s", er.chainID.String())
	}
	ageThreshold := er.config.EthTxResendAfterThreshold()
	maxInFlightTransactions := er.config.EvmMaxInFlightTransactions()
	olderThan := time.Now().Add(-ageThreshold)
	var allAttempts []EthTxAttempt
	for _, k := range keys {
		var attempts []EthTxAttempt
		attempts, err = FindEthTxAttemptsRequiringResend(er.db, olderThan, maxInFlightTransactions, er.chainID, k.Address)
		if err != nil {
			return errors.Wrap(err, "failed to FindEthTxAttemptsRequiringResend")
		}

		allAttempts = append(allAttempts, attempts...)
	}

	if len(allAttempts) == 0 {
		return nil
	}
	er.logger.Infow(fmt.Sprintf("Re-sending %d unconfirmed transactions that were last sent over %s ago. These transactions are taking longer than usual to be mined. %s", len(allAttempts), ageThreshold, label.NodeConnectivityProblemWarning), "n", len(allAttempts))

	batchSize := int(er.config.EvmRPCDefaultBatchSize())
	ctx, cancel := context.WithTimeout(er.ctx, batchSendTransactionTimeout)
	defer cancel()
	reqs, err := batchSendTransactions(ctx, er.db, allAttempts, batchSize, er.logger, er.ethClient)
	if err != nil {
		return errors.Wrap(err, "failed to re-send transactions")
	}
	logResendResult(er.logger, reqs)

	return nil
}

// FindEthTxAttemptsRequiringResend returns the highest priced attempt for each
// eth_tx that was last sent before or at the given time (up to limit)
func FindEthTxAttemptsRequiringResend(db *sqlx.DB, olderThan time.Time, maxInFlightTransactions uint32, chainID big.Int, address common.Address) (attempts []EthTxAttempt, err error) {
	var limit null.Uint32
	if maxInFlightTransactions > 0 {
		limit = null.Uint32From(maxInFlightTransactions)
	}
	// this select distinct works because of unique index on eth_txes
	// (evm_chain_id, from_address, nonce)
	err = db.Select(&attempts, `
SELECT DISTINCT ON (nonce) eth_tx_attempts.*
FROM eth_tx_attempts
JOIN eth_txes ON eth_txes.id = eth_tx_attempts.eth_tx_id AND eth_txes.state IN ('unconfirmed', 'confirmed_missing_receipt')
WHERE eth_tx_attempts.state <> 'in_progress' AND eth_txes.broadcast_at <= $1 AND evm_chain_id = $2 AND from_address = $3
ORDER BY eth_txes.nonce ASC, eth_tx_attempts.gas_price DESC, eth_tx_attempts.gas_tip_cap DESC
LIMIT $4
`, olderThan, chainID.String(), address, limit)

	return attempts, errors.Wrap(err, "FindEthTxAttemptsRequiringResend failed to load eth_tx_attempts")
}

func logResendResult(lggr logger.Logger, reqs []rpc.BatchElem) {
	var nNew int
	var nFatal int
	for _, req := range reqs {
		serr := evmclient.NewSendError(req.Error)
		if serr == nil {
			nNew++
		} else if serr.Fatal() {
			nFatal++
		}
	}
	lggr.Debugw("Completed", "n", len(reqs), "nNew", nNew, "nFatal", nFatal)
}
