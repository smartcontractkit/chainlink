package txmgr

import (
	"math/big"
	"time"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/forwarders"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

// NewTxm constructs the necessary dependencies for the EvmTxm (broadcaster, confirmer, etc) and returns a new EvmTxManager
func NewTxm(
	db *sqlx.DB,
	cfg Config,
	txConfig config.Transactions,
	dbConfig DatabaseConfig,
	listenerConfig ListenerConfig,
	client evmclient.Client,
	lggr logger.Logger,
	logPoller logpoller.LogPoller,
	keyStore keystore.Eth,
	eventBroadcaster pg.EventBroadcaster,
	estimator gas.EvmFeeEstimator,
) (txm EvmTxManager,
	err error,
) {
	var fwdMgr EvmFwdMgr

	if txConfig.ForwardersEnabled() {
		fwdMgr = forwarders.NewFwdMgr(db, client, logPoller, lggr, cfg, dbConfig)
	} else {
		lggr.Info("EvmForwarderManager: Disabled")
	}
	checker := &CheckerFactory{Client: client}
	// create tx attempt builder
	txAttemptBuilder := NewEvmTxAttemptBuilder(*client.ConfiguredChainID(), cfg, keyStore, estimator)
	txStore := NewTxStore(db, lggr, dbConfig)
	txNonceSyncer := NewNonceSyncer(txStore, lggr, client, keyStore)

	txmCfg := NewEvmTxmConfig(cfg)       // wrap Evm specific config
	txmClient := NewEvmTxmClient(client) // wrap Evm specific client
	ethBroadcaster := NewEvmBroadcaster(txStore, txmClient, txmCfg, txConfig, listenerConfig, keyStore, eventBroadcaster, txAttemptBuilder, txNonceSyncer, lggr, checker, cfg.EvmNonceAutoSync())
	ethConfirmer := NewEvmConfirmer(txStore, txmClient, txmCfg, txConfig, dbConfig, keyStore, txAttemptBuilder, lggr)
	var ethResender *EvmResender
	if txConfig.ResendAfterThreshold() > 0 {
		ethResender = NewEvmResender(lggr, txStore, txmClient, keyStore, txmgr.DefaultResenderPollInterval, txmCfg, txConfig)
	}
	txm = NewEvmTxm(txmClient.ConfiguredChainID(), txmCfg, txConfig, keyStore, lggr, checker, fwdMgr, txAttemptBuilder, txStore, txNonceSyncer, ethBroadcaster, ethConfirmer, ethResender)
	return txm, nil
}

// NewEvmTxm creates a new concrete EvmTxm
func NewEvmTxm(
	chainId *big.Int,
	cfg txmgrtypes.TransactionManagerConfig,
	txCfg txmgrtypes.TransactionManagerTransactionsConfig,
	keyStore EvmKeyStore,
	lggr logger.Logger,
	checkerFactory EvmTransmitCheckerFactory,
	fwdMgr EvmFwdMgr,
	txAttemptBuilder EvmTxAttemptBuilder,
	txStore TxStore,
	nonceSyncer EvmNonceSyncer,
	broadcaster *EvmBroadcaster,
	confirmer *EvmConfirmer,
	resender *EvmResender,
) *EvmTxm {
	return txmgr.NewTxm(chainId, cfg, txCfg, keyStore, lggr, checkerFactory, fwdMgr, txAttemptBuilder, txStore, nonceSyncer, broadcaster, confirmer, resender)
}

// NewEvnResender creates a new concrete EvmResender
func NewEvmResender(
	lggr logger.Logger,
	txStore TransactionStore,
	client TransactionClient,
	ks EvmKeyStore,
	pollInterval time.Duration,
	config EvmResenderConfig,
	txConfig txmgrtypes.ResenderTransactionsConfig,
) *EvmResender {
	return txmgr.NewResender(lggr, txStore, client, ks, pollInterval, config, txConfig)
}

// NewEvmReaper instantiates a new EVM-specific reaper object
func NewEvmReaper(lggr logger.Logger, store txmgrtypes.TxHistoryReaper[*big.Int], config EvmReaperConfig, txConfig txmgrtypes.ReaperTransactionsConfig, chainID *big.Int) *EvmReaper {
	return txmgr.NewReaper(lggr, store, config, txConfig, chainID)
}

// NewEvmConfirmer instantiates a new EVM confirmer
func NewEvmConfirmer(
	txStore TxStore,
	client EvmTxmClient,
	config txmgrtypes.ConfirmerConfig,
	txConfig txmgrtypes.ConfirmerTransactionsConfig,
	dbConfig txmgrtypes.ConfirmerDatabaseConfig,
	keystore EvmKeyStore,
	txAttemptBuilder EvmTxAttemptBuilder,
	lggr logger.Logger,
) *EvmConfirmer {
	return txmgr.NewConfirmer(txStore, client, config, txConfig, dbConfig, keystore, txAttemptBuilder, lggr, func(r *evmtypes.Receipt) bool { return r == nil })
}

// NewEvmBroadcaster returns a new concrete EvmBroadcaster
func NewEvmBroadcaster(
	txStore TransactionStore,
	client TransactionClient,
	config txmgrtypes.BroadcasterConfig,
	txConfig txmgrtypes.BroadcasterTransactionsConfig,
	listenerConfig txmgrtypes.BroadcasterListenerConfig,
	keystore EvmKeyStore,
	eventBroadcaster pg.EventBroadcaster,
	txAttemptBuilder EvmTxAttemptBuilder,
	nonceSyncer EvmNonceSyncer,
	logger logger.Logger,
	checkerFactory EvmTransmitCheckerFactory,
	autoSyncNonce bool,
) *EvmBroadcaster {
	return txmgr.NewBroadcaster(txStore, client, config, txConfig, listenerConfig, keystore, eventBroadcaster, txAttemptBuilder, nonceSyncer, logger, checkerFactory, autoSyncNonce, stringToGethAddress)
}
