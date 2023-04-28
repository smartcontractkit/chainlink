package txmgr

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/label"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
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
type EthResender[CHAIN_ID txmgrtypes.ID, ADDR types.Hashable, TX_HASH types.Hashable, BLOCK_HASH types.Hashable, SEQ txmgrtypes.Sequence] struct {
	txStore             txmgrtypes.TxStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, txmgrtypes.NewTx[ADDR, TX_HASH], *evmtypes.Receipt, EthTx[ADDR, TX_HASH], EthTxAttempt[ADDR, TX_HASH], SEQ]
	ethClient           evmclient.Client
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
		ethClient,
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
func (er *EthResender[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ]) Start() {
	er.logger.Debugf("Enabled with poll interval of %s and age threshold of %s", er.interval, er.config.TxResendAfterThreshold())
	go er.runLoop()
}

// Stop is a comment which satisfies the linter
func (er *EthResender[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ]) Stop() {
	er.cancel()
	<-er.chDone
}

func (er *EthResender[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ]) runLoop() {
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

func (er *EthResender[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ]) resendUnconfirmed() error {
	enabledAddresses, err := er.ks.EnabledAddressesForChain(er.chainID)
	if err != nil {
		return errors.Wrapf(err, "EthResender failed getting enabled keys for chain %s", er.chainID.String())
	}
	ageThreshold := er.config.TxResendAfterThreshold()
	maxInFlightTransactions := er.config.MaxInFlightTransactions()
	olderThan := time.Now().Add(-ageThreshold)
	var allAttempts []EthTxAttempt[ADDR, TX_HASH]
	for _, k := range enabledAddresses {
		var attempts []EthTxAttempt[ADDR, TX_HASH]
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
	reqs, err := batchSendTransactions(ctx, er.txStore, allAttempts, batchSize, er.logger, er.ethClient)
	if err != nil {
		return errors.Wrap(err, "failed to re-send transactions")
	}
	logResendResult(er.logger, reqs)

	return nil
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

func (er *EthResender[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ]) logStuckAttempts(attempts []EthTxAttempt[ADDR, TX_HASH], fromAddress ADDR) {
	if time.Since(er.lastAlertTimestamps[fromAddress.String()]) >= unconfirmedTxAlertLogFrequency {
		oldestAttempt, exists := findOldestUnconfirmedAttempt(attempts)
		if exists {
			// Wait at least 2 times the EthTxResendAfterThreshold to log critical with an unconfirmedTxAlertDelay
			if time.Since(oldestAttempt.CreatedAt) > er.config.TxResendAfterThreshold()*2 {
				er.lastAlertTimestamps[fromAddress.String()] = time.Now()
				er.logger.Errorw("TxAttempt has been unconfirmed for more than: ", er.config.TxResendAfterThreshold()*2,
					"txID", oldestAttempt.EthTxID, "GasPrice", oldestAttempt.GasPrice, "GasTipCap", oldestAttempt.GasTipCap, "GasFeeCap", oldestAttempt.GasFeeCap,
					"BroadcastBeforeBlockNum", oldestAttempt.BroadcastBeforeBlockNum, "Hash", oldestAttempt.Hash, "fromAddress", fromAddress)
			}
		}
	}
}

func findOldestUnconfirmedAttempt[ADDR types.Hashable, TX_HASH types.Hashable](attempts []EthTxAttempt[ADDR, TX_HASH]) (EthTxAttempt[ADDR, TX_HASH], bool) {
	var oldestAttempt EthTxAttempt[ADDR, TX_HASH]
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
