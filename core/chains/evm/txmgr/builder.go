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
	chainConfig ChainConfig,
	fCfg FeeConfig,
	txConfig config.Transactions,
	dbConfig DatabaseConfig,
	listenerConfig ListenerConfig,
	client evmclient.Client,
	lggr logger.Logger,
	logPoller logpoller.LogPoller,
	keyStore keystore.Eth,
	eventBroadcaster pg.EventBroadcaster,
	estimator gas.EvmFeeEstimator,
) (txm TxManager,
	err error,
) {
	var fwdMgr FwdMgr

	if txConfig.ForwardersEnabled() {
		fwdMgr = forwarders.NewFwdMgr(db, client, logPoller, lggr, chainConfig, dbConfig)
	} else {
		lggr.Info("EvmForwarderManager: Disabled")
	}
	checker := &CheckerFactory{Client: client}
	// create tx attempt builder
	txAttemptBuilder := NewEvmTxAttemptBuilder(*client.ConfiguredChainID(), fCfg, keyStore, estimator)
	txStore := NewTxStore(db, lggr, dbConfig)
	txNonceSyncer := NewNonceSyncer(txStore, lggr, client)

	txmCfg := NewEvmTxmConfig(chainConfig) // wrap Evm specific config
	feeCfg := NewEvmTxmFeeConfig(fCfg)     // wrap Evm specific config
	txmClient := NewEvmTxmClient(client)   // wrap Evm specific client
	ethBroadcaster := NewEvmBroadcaster(txStore, txmClient, txmCfg, feeCfg, txConfig, listenerConfig, keyStore, eventBroadcaster, txAttemptBuilder, txNonceSyncer, lggr, checker, chainConfig.NonceAutoSync())
	ethConfirmer := NewEvmConfirmer(txStore, txmClient, txmCfg, feeCfg, txConfig, dbConfig, keyStore, txAttemptBuilder, lggr)
	var ethResender *Resender
	if txConfig.ResendAfterThreshold() > 0 {
		ethResender = NewEvmResender(lggr, txStore, txmClient, keyStore, txmgr.DefaultResenderPollInterval, chainConfig, txConfig)
	}
	txm = NewEvmTxm(txmClient.ConfiguredChainID(), txmCfg, txConfig, keyStore, lggr, checker, fwdMgr, txAttemptBuilder, txStore, txNonceSyncer, ethBroadcaster, ethConfirmer, ethResender)
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
	nonceSyncer NonceSyncer,
	broadcaster *Broadcaster,
	confirmer *Confirmer,
	resender *Resender,
) *Txm {
	return txmgr.NewTxm(chainId, cfg, txCfg, keyStore, lggr, checkerFactory, fwdMgr, txAttemptBuilder, txStore, nonceSyncer, broadcaster, confirmer, resender)
}

// NewEvnResender creates a new concrete EvmResender
func NewEvmResender(
	lggr logger.Logger,
	txStore TransactionStore,
	client TransactionClient,
	ks KeyStore,
	pollInterval time.Duration,
	config EvmResenderConfig,
	txConfig txmgrtypes.ResenderTransactionsConfig,
) *Resender {
	return txmgr.NewResender(lggr, txStore, client, ks, pollInterval, config, txConfig)
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

// NewEvmBroadcaster returns a new concrete EvmBroadcaster
func NewEvmBroadcaster(
	txStore TransactionStore,
	client TransactionClient,
	chainConfig txmgrtypes.BroadcasterChainConfig,
	feeConfig txmgrtypes.BroadcasterFeeConfig,
	txConfig txmgrtypes.BroadcasterTransactionsConfig,
	listenerConfig txmgrtypes.BroadcasterListenerConfig,
	keystore KeyStore,
	eventBroadcaster pg.EventBroadcaster,
	txAttemptBuilder TxAttemptBuilder,
	nonceSyncer NonceSyncer,
	logger logger.Logger,
	checkerFactory TransmitCheckerFactory,
	autoSyncNonce bool,
) *Broadcaster {
	return txmgr.NewBroadcaster(txStore, client, chainConfig, feeConfig, txConfig, listenerConfig, keystore, eventBroadcaster, txAttemptBuilder, nonceSyncer, logger, checkerFactory, autoSyncNonce, stringToGethAddress, evmtypes.GenerateNextNonce)
}
