package txmgr

import (
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/forwarders"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

// NewTxm constructs the necessary dependencies for the EvmTxm (broadcaster, confirmer, etc) and returns a new EvmTxManager
func NewTxm(
	ds sqlutil.DataSource,
	chainConfig ChainConfig,
	fCfg FeeConfig,
	txConfig config.Transactions,
	clientErrors config.ClientErrors,
	dbConfig DatabaseConfig,
	listenerConfig ListenerConfig,
	client evmclient.Client,
	lggr logger.Logger,
	logPoller logpoller.LogPoller,
	keyStore keystore.Eth,
	estimator gas.EvmFeeEstimator,
) (txm TxManager,
	err error,
) {
	var fwdMgr FwdMgr

	if txConfig.ForwardersEnabled() {
		fwdMgr = forwarders.NewFwdMgr(ds, client, logPoller, lggr, chainConfig)
	} else {
		lggr.Info("EvmForwarderManager: Disabled")
	}
	checker := &CheckerFactory{Client: client}
	// create tx attempt builder
	txAttemptBuilder := NewEvmTxAttemptBuilder(*client.ConfiguredChainID(), fCfg, keyStore, estimator)
	txStore := NewTxStore(ds, lggr)
	txmCfg := NewEvmTxmConfig(chainConfig)             // wrap Evm specific config
	feeCfg := NewEvmTxmFeeConfig(fCfg)                 // wrap Evm specific config
	txmClient := NewEvmTxmClient(client, clientErrors) // wrap Evm specific client
	chainID := txmClient.ConfiguredChainID()
	evmBroadcaster := NewEvmBroadcaster(txStore, txmClient, txmCfg, feeCfg, txConfig, listenerConfig, keyStore, txAttemptBuilder, lggr, checker, chainConfig.NonceAutoSync())
	evmTracker := NewEvmTracker(txStore, keyStore, chainID, lggr)
	evmConfirmer := NewEvmConfirmer(txStore, txmClient, txmCfg, feeCfg, txConfig, dbConfig, keyStore, txAttemptBuilder, lggr)
	var evmResender *Resender
	if txConfig.ResendAfterThreshold() > 0 {
		evmResender = NewEvmResender(lggr, txStore, txmClient, evmTracker, keyStore, txmgr.DefaultResenderPollInterval, chainConfig, txConfig)
	}
	txm = NewEvmTxm(chainID, txmCfg, txConfig, keyStore, lggr, checker, fwdMgr, txAttemptBuilder, txStore, evmBroadcaster, evmConfirmer, evmResender, evmTracker)
	return txm, nil
}

// NewEvmTxm creates a new concrete EvmTxm
func NewEvmTxm(
	chainId *big.Int,
	cfg txmgrtypes.TransactionManagerChainConfig,
	txCfg txmgrtypes.TransactionManagerTransactionsConfig,
	keyStore KeyStore,
	lggr logger.Logger,
	checkerFactory TransmitCheckerFactory,
	fwdMgr FwdMgr,
	txAttemptBuilder TxAttemptBuilder,
	txStore TxStore,
	broadcaster *Broadcaster,
	confirmer *Confirmer,
	resender *Resender,
	tracker *Tracker,
) *Txm {
	return txmgr.NewTxm(chainId, cfg, txCfg, keyStore, lggr, checkerFactory, fwdMgr, txAttemptBuilder, txStore, broadcaster, confirmer, resender, tracker)
}

// NewEvmResender creates a new concrete EvmResender
func NewEvmResender(
	lggr logger.Logger,
	txStore TransactionStore,
	client TransactionClient,
	tracker *Tracker,
	ks KeyStore,
	pollInterval time.Duration,
	config EvmResenderConfig,
	txConfig txmgrtypes.ResenderTransactionsConfig,
) *Resender {
	return txmgr.NewResender(lggr, txStore, client, tracker, ks, pollInterval, config, txConfig)
}

// NewEvmReaper instantiates a new EVM-specific reaper object
func NewEvmReaper(lggr logger.Logger, store txmgrtypes.TxHistoryReaper[*big.Int], config EvmReaperConfig, txConfig txmgrtypes.ReaperTransactionsConfig, chainID *big.Int) *Reaper {
	return txmgr.NewReaper(lggr, store, config, txConfig, chainID)
}

// NewEvmConfirmer instantiates a new EVM confirmer
func NewEvmConfirmer(
	txStore TxStore,
	client TxmClient,
	chainConfig txmgrtypes.ConfirmerChainConfig,
	feeConfig txmgrtypes.ConfirmerFeeConfig,
	txConfig txmgrtypes.ConfirmerTransactionsConfig,
	dbConfig txmgrtypes.ConfirmerDatabaseConfig,
	keystore KeyStore,
	txAttemptBuilder TxAttemptBuilder,
	lggr logger.Logger,
) *Confirmer {
	return txmgr.NewConfirmer(txStore, client, chainConfig, feeConfig, txConfig, dbConfig, keystore, txAttemptBuilder, lggr, func(r *evmtypes.Receipt) bool { return r == nil })
}

// NewEvmTracker instantiates a new EVM tracker for abandoned transactions
func NewEvmTracker(
	txStore TxStore,
	keyStore KeyStore,
	chainID *big.Int,
	lggr logger.Logger,
) *Tracker {
	return txmgr.NewTracker(txStore, keyStore, chainID, lggr)
}

// NewEvmBroadcaster returns a new concrete EvmBroadcaster
func NewEvmBroadcaster(
	txStore TransactionStore,
	client TransactionClient,
	chainConfig txmgrtypes.BroadcasterChainConfig,
	feeConfig txmgrtypes.BroadcasterFeeConfig,
	txConfig txmgrtypes.BroadcasterTransactionsConfig,
	listenerConfig txmgrtypes.BroadcasterListenerConfig,
	keystore KeyStore,
	txAttemptBuilder TxAttemptBuilder,
	logger logger.Logger,
	checkerFactory TransmitCheckerFactory,
	autoSyncNonce bool,
) *Broadcaster {
	nonceTracker := NewNonceTracker(logger, txStore, client)
	return txmgr.NewBroadcaster(txStore, client, chainConfig, feeConfig, txConfig, listenerConfig, keystore, txAttemptBuilder, nonceTracker, logger, checkerFactory, autoSyncNonce)
}
