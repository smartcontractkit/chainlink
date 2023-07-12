package txmgr

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	clienttypes "github.com/smartcontractkit/chainlink/v2/common/chains/client"
	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/label"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	// pollInterval is the maximum amount of time in addition to
	// EthTxResendAfterThreshold that we will wait before resending an attempt
	DefaultResenderPollInterval = 5 * time.Second

	// Alert interval for unconfirmed transaction attempts
	unconfirmedTxAlertLogFrequency = 2 * time.Minute

	// timeout value for batchSendTransactions
	batchSendTransactionTimeout = 30 * time.Second
)

// EthResender periodically picks up transactions that have been languishing
// unconfirmed for a configured amount of time without being sent, and sends
// their highest priced attempt again. This helps to defend against geth/parity
// silently dropping txes, or txes being ejected from the mempool.
//
// Previously we relied on the bumper to do this for us implicitly but there
// can occasionally be problems with this (e.g. abnormally long block times, or
// if gas bumping is disabled)
type Resender[
	CHAIN_ID types.ID,
	ADDR types.Hashable,
	TX_HASH types.Hashable,
	BLOCK_HASH types.Hashable,
	SEQ types.Sequence,
	FEE feetypes.Fee,
] struct {
	txStore             txmgrtypes.TransactionStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, SEQ, FEE]
	client              txmgrtypes.TransactionClient[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	ks                  txmgrtypes.KeyStore[ADDR, CHAIN_ID, SEQ]
	chainID             CHAIN_ID
	interval            time.Duration
	config              txmgrtypes.ResenderChainConfig
	txConfig            txmgrtypes.ResenderTransactionsConfig
	logger              logger.Logger
	lastAlertTimestamps map[string]time.Time

	ctx    context.Context
	cancel context.CancelFunc
	chDone chan struct{}
}

func NewResender[
	CHAIN_ID types.ID,
	ADDR types.Hashable,
	TX_HASH types.Hashable,
	BLOCK_HASH types.Hashable,
	SEQ types.Sequence,
	FEE feetypes.Fee,
](
	lggr logger.Logger,
	txStore txmgrtypes.TransactionStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, SEQ, FEE],
	client txmgrtypes.TransactionClient[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	ks txmgrtypes.KeyStore[ADDR, CHAIN_ID, SEQ],
	pollInterval time.Duration,
	config txmgrtypes.ResenderChainConfig,
	txConfig txmgrtypes.ResenderTransactionsConfig,
) *Resender[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] {
	if txConfig.ResendAfterThreshold() == 0 {
		panic("Resender requires a non-zero threshold")
	}
	// todo: add context to evmTxStore
	ctx, cancel := context.WithCancel(context.Background())
	return &Resender[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]{
		txStore,
		client,
		ks,
		client.ConfiguredChainID(),
		pollInterval,
		config,
		txConfig,
		lggr.Named("Resender"),
		make(map[string]time.Time),
		ctx,
		cancel,
		make(chan struct{}),
	}
}

// Start is a comment which satisfies the linter
func (er *Resender[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) Start() {
	er.logger.Debugf("Enabled with poll interval of %s and age threshold of %s", er.interval, er.txConfig.ResendAfterThreshold())
	go er.runLoop()
}

// Stop is a comment which satisfies the linter
func (er *Resender[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) Stop() {
	er.cancel()
	<-er.chDone
}

func (er *Resender[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) runLoop() {
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

func (er *Resender[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) resendUnconfirmed() error {
	enabledAddresses, err := er.ks.EnabledAddressesForChain(er.chainID)
	if err != nil {
		return errors.Wrapf(err, "EthResender failed getting enabled keys for chain %s", er.chainID.String())
	}
	ageThreshold := er.txConfig.ResendAfterThreshold()
	maxInFlightTransactions := er.txConfig.MaxInFlight()
	olderThan := time.Now().Add(-ageThreshold)
	var allAttempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	for _, k := range enabledAddresses {
		var attempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
		attempts, err = er.txStore.FindTxAttemptsRequiringResend(olderThan, maxInFlightTransactions, er.chainID, k)
		if err != nil {
			return errors.Wrap(err, "failed to FindEthTxAttemptsRequiringResend")
		}
		er.logStuckAttempts(attempts, k)

		allAttempts = append(allAttempts, attempts...)
	}

	if len(allAttempts) == 0 {
		for k := range er.lastAlertTimestamps {
			er.lastAlertTimestamps[k] = time.Now()
		}
		return nil
	}
	er.logger.Infow(fmt.Sprintf("Re-sending %d unconfirmed transactions that were last sent over %s ago. These transactions are taking longer than usual to be mined. %s", len(allAttempts), ageThreshold, label.NodeConnectivityProblemWarning), "n", len(allAttempts))

	batchSize := int(er.config.RPCDefaultBatchSize())
	ctx, cancel := context.WithTimeout(er.ctx, batchSendTransactionTimeout)
	defer cancel()
	txErrTypes, _, err := er.client.BatchSendTransactions(ctx, er.txStore.UpdateBroadcastAts, allAttempts, batchSize, er.logger)
	if err != nil {
		return errors.Wrap(err, "failed to re-send transactions")
	}
	logResendResult(er.logger, txErrTypes)

	return nil
}

func logResendResult(lggr logger.Logger, codes []clienttypes.SendTxReturnCode) {
	var nNew int
	var nFatal int
	for _, c := range codes {
		if c == clienttypes.Successful {
			nNew++
		} else if c == clienttypes.Fatal {
			nFatal++
		}
	}
	lggr.Debugw("Completed", "n", len(codes), "nNew", nNew, "nFatal", nFatal)
}

func (er *Resender[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) logStuckAttempts(attempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], fromAddress ADDR) {
	if time.Since(er.lastAlertTimestamps[fromAddress.String()]) >= unconfirmedTxAlertLogFrequency {
		oldestAttempt, exists := findOldestUnconfirmedAttempt(attempts)
		if exists {
			// Wait at least 2 times the EthTxResendAfterThreshold to log critical with an unconfirmedTxAlertDelay
			if time.Since(oldestAttempt.CreatedAt) > er.txConfig.ResendAfterThreshold()*2 {
				er.lastAlertTimestamps[fromAddress.String()] = time.Now()
				er.logger.Errorw("TxAttempt has been unconfirmed for more than max duration", "maxDuration", er.txConfig.ResendAfterThreshold()*2,
					"txID", oldestAttempt.TxID, "txFee", oldestAttempt.TxFee,
					"BroadcastBeforeBlockNum", oldestAttempt.BroadcastBeforeBlockNum, "Hash", oldestAttempt.Hash, "fromAddress", fromAddress)
			}
		}
	}
}

func findOldestUnconfirmedAttempt[
	CHAIN_ID types.ID,
	ADDR types.Hashable,
	TX_HASH, BLOCK_HASH types.Hashable,
	SEQ types.Sequence,
	FEE feetypes.Fee,
](attempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) (txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], bool) {
	var oldestAttempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	if len(attempts) < 1 {
		return oldestAttempt, false
	}
	oldestAttempt = attempts[0]
	for i := 1; i < len(attempts); i++ {
		if oldestAttempt.CreatedAt.Sub(attempts[i].CreatedAt) <= 0 {
			oldestAttempt = attempts[i]
		}
	}
	return oldestAttempt, true
}
