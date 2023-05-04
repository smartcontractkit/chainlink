package txmgr

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	clienttypes "github.com/smartcontractkit/chainlink/v2/common/chains/client"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/label"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// pollInterval is the maximum amount of time in addition to
// EthTxResendAfterThreshold that we will wait before resending an attempt
const DefaultResenderPollInterval = 5 * time.Second

// Alert interval for unconfirmed transaction attempts
const unconfirmedTxAlertLogFrequency = 2 * time.Minute

// EthResender periodically picks up transactions that have been languishing
// unconfirmed for a configured amount of time without being sent, and sends
// their highest priced attempt again. This helps to defend against geth/parity
// silently dropping txes, or txes being ejected from the mempool.
//
// Previously we relied on the bumper to do this for us implicitly but there
// can occasionally be problems with this (e.g. abnormally long block times, or
// if gas bumping is disabled)
type EthResender[
	CHAIN_ID txmgrtypes.ID,
	ADDR types.Hashable,
	TX_HASH types.Hashable,
	BLOCK_HASH types.Hashable,
	SEQ txmgrtypes.Sequence,
	FEE txmgrtypes.Fee,
	R txmgrtypes.ChainReceipt[TX_HASH],
	ADD any,
] struct {
	txStore             txmgrtypes.TxStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD]
	client              TxmClient[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD]
	ks                  txmgrtypes.KeyStore[ADDR, CHAIN_ID, SEQ]
	chainID             CHAIN_ID
	interval            time.Duration
	config              EvmResenderConfig
	logger              logger.Logger
	lastAlertTimestamps map[string]time.Time

	ctx    context.Context
	cancel context.CancelFunc
	chDone chan struct{}
}

// NewEthResender creates a new concrete EthResender
func NewEthResender(
	lggr logger.Logger,
	txStore EvmTxStore,
	ethClient evmclient.Client, ks EvmKeyStore,
	pollInterval time.Duration,
	config EvmResenderConfig,
) *EvmResender {
	if config.TxResendAfterThreshold() == 0 {
		panic("EthResender requires a non-zero threshold")
	}
	// todo: add context to evmTxStore
	ctx, cancel := context.WithCancel(context.Background())
	return &EvmResender{
		txStore,
		NewEvmTxmClient(ethClient),
		ks,
		ethClient.ConfiguredChainID(),
		pollInterval,
		config,
		lggr.Named("EthResender"),
		make(map[string]time.Time),
		ctx,
		cancel,
		make(chan struct{}),
	}
}

// Start is a comment which satisfies the linter
func (er *EthResender[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE, R, ADD]) Start() {
	er.logger.Debugf("Enabled with poll interval of %s and age threshold of %s", er.interval, er.config.TxResendAfterThreshold())
	go er.runLoop()
}

// Stop is a comment which satisfies the linter
func (er *EthResender[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE, R, ADD]) Stop() {
	er.cancel()
	<-er.chDone
}

func (er *EthResender[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE, R, ADD]) runLoop() {
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

func (er *EthResender[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE, R, ADD]) resendUnconfirmed() error {
	enabledAddresses, err := er.ks.EnabledAddressesForChain(er.chainID)
	if err != nil {
		return errors.Wrapf(err, "EthResender failed getting enabled keys for chain %s", er.chainID.String())
	}
	ageThreshold := er.config.TxResendAfterThreshold()
	maxInFlightTransactions := er.config.MaxInFlightTransactions()
	olderThan := time.Now().Add(-ageThreshold)
	var allAttempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD]
	for _, k := range enabledAddresses {
		var attempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD]
		attempts, err = er.txStore.FindEthTxAttemptsRequiringResend(olderThan, maxInFlightTransactions, er.chainID, k)
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
	txErrTypes, _, err := er.client.BatchSendTransactions(ctx, er.txStore, allAttempts, batchSize, er.logger)
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

func (er *EthResender[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE, R, ADD]) logStuckAttempts(attempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], fromAddress ADDR) {
	if time.Since(er.lastAlertTimestamps[fromAddress.String()]) >= unconfirmedTxAlertLogFrequency {
		oldestAttempt, exists := findOldestUnconfirmedAttempt(attempts)
		if exists {
			// Wait at least 2 times the EthTxResendAfterThreshold to log critical with an unconfirmedTxAlertDelay
			if time.Since(oldestAttempt.CreatedAt) > er.config.TxResendAfterThreshold()*2 {
				er.lastAlertTimestamps[fromAddress.String()] = time.Now()
				er.logger.Errorw("TxAttempt has been unconfirmed for more than: ", er.config.TxResendAfterThreshold()*2,
					"txID", oldestAttempt.TxID, "txFee", oldestAttempt.Fee(),
					"BroadcastBeforeBlockNum", oldestAttempt.BroadcastBeforeBlockNum, "Hash", oldestAttempt.Hash, "fromAddress", fromAddress)
			}
		}
	}
}

func findOldestUnconfirmedAttempt[
	CHAIN_ID txmgrtypes.ID,
	ADDR types.Hashable,
	TX_HASH, BLOCK_HASH types.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH],
	FEE txmgrtypes.Fee,
	ADD any,
](attempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD]) (txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], bool) {
	var oldestAttempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD]
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
