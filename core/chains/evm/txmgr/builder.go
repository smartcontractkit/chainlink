package txmgr

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	"github.com/smartcontractkit/chainlink-framework/chains/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink-framework/chains/txmgr/types"
	"github.com/smartcontractkit/chainlink-framework/metrics"

	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/config"
	"github.com/smartcontractkit/chainlink-evm/pkg/config/chaintype"
	"github.com/smartcontractkit/chainlink-evm/pkg/forwarders"
	"github.com/smartcontractkit/chainlink-evm/pkg/gas"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink-evm/pkg/txm"
	"github.com/smartcontractkit/chainlink-evm/pkg/txm/clientwrappers"
	"github.com/smartcontractkit/chainlink-evm/pkg/txm/storage"
	"github.com/smartcontractkit/chainlink-evm/pkg/types"
)

type latestAndFinalizedBlockHeadTracker interface {
	LatestAndFinalizedBlock(ctx context.Context) (latest, finalized *types.Head, err error)
}

// NewTxm constructs the necessary dependencies for the EvmTxm (broadcaster, confirmer, etc) and returns a new EvmTxManager
func NewTxm(
	ds sqlutil.DataSource,
	chainConfig ChainConfig,
	fCfg FeeConfig,
	txConfig config.Transactions,
	clientErrors config.ClientErrors,
	dbConfig DatabaseConfig,
	listenerConfig ListenerConfig,
	client client.Client,
	lggr logger.Logger,
	logPoller logpoller.LogPoller,
	keyStore keys.ChainStore,
	estimator gas.EvmFeeEstimator,
	headTracker latestAndFinalizedBlockHeadTracker,
	txmv2wrapper TxManager,
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
	metrics, err := NewEVMTxmMetrics(chainID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize EVM TXM metrics: %w", err)
	}
	evmBroadcaster := NewEvmBroadcaster(txStore, txmClient, txmCfg, feeCfg, txConfig, listenerConfig, keyStore, txAttemptBuilder, lggr, checker, chainConfig.NonceAutoSync(), chainConfig.ChainType(), metrics)
	evmTracker := NewEvmTracker(txStore, keyStore, chainID, lggr)
	stuckTxDetector := NewStuckTxDetector(lggr, client.ConfiguredChainID(), chainConfig.ChainType(), fCfg.PriceMax(), txConfig.AutoPurge(), estimator, txStore, client)
	evmConfirmer := NewEvmConfirmer(txStore, txmClient, feeCfg, txConfig, dbConfig, keyStore, txAttemptBuilder, lggr, stuckTxDetector, metrics)
	evmFinalizer := NewEvmFinalizer(lggr, client.ConfiguredChainID(), chainConfig.RPCDefaultBatchSize(), txConfig.ForwardersEnabled(), txStore, txmClient, headTracker, metrics)
	var evmResender *Resender
	if txConfig.ResendAfterThreshold() > 0 {
		evmResender = NewEvmResender(lggr, txStore, txmClient, evmTracker, keyStore, txmgr.DefaultResenderPollInterval, chainConfig, txConfig)
	}
	txm = NewEvmTxm(chainID, txmCfg, txConfig, keyStore, lggr, checker, fwdMgr, txAttemptBuilder, txStore, evmBroadcaster, evmConfirmer, evmResender, evmTracker, evmFinalizer, txmv2wrapper)
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
	finalizer Finalizer,
	txmv2wrapper TxManager,
) *Txm {
	return txmgr.NewTxm(chainId, cfg, txCfg, keyStore, lggr, checkerFactory, fwdMgr, txAttemptBuilder, txStore, broadcaster, confirmer, resender, tracker, finalizer, client.NewTxError, txmv2wrapper)
}

func NewTxmV2(
	ds sqlutil.DataSource,
	chainConfig ChainConfig,
	fCfg FeeConfig,
	txConfig config.Transactions,
	txmV2Config config.TransactionManagerV2,
	client client.Client,
	lggr logger.Logger,
	logPoller logpoller.LogPoller,
	keyStore keys.ChainStore,
	estimator gas.EvmFeeEstimator,
) (TxManager, error) {
	var fwdMgr *forwarders.FwdMgr
	if txConfig.ForwardersEnabled() {
		fwdMgr = forwarders.NewFwdMgr(ds, client, logPoller, lggr, chainConfig)
	} else {
		lggr.Info("ForwarderManager: Disabled")
	}

	chainID := client.ConfiguredChainID()

	var stuckTxDetector txm.StuckTxDetector
	if txConfig.AutoPurge().Enabled() {
		stuckTxDetectorConfig := txm.StuckTxDetectorConfig{
			BlockTime:             *txmV2Config.BlockTime(),
			StuckTxBlockThreshold: *txConfig.AutoPurge().Threshold(),
			DetectionURL:          txConfig.AutoPurge().DetectionApiUrl().String(),
		}
		stuckTxDetector = txm.NewStuckTxDetector(lggr, chainConfig.ChainType(), stuckTxDetectorConfig)
	}

	attemptBuilder := txm.NewAttemptBuilder(fCfg.PriceMaxKey, estimator, keyStore)
	inMemoryStoreManager := storage.NewInMemoryStoreManager(lggr, chainID)
	config := txm.Config{
		EIP1559:   fCfg.EIP1559DynamicFees(),
		BlockTime: *txmV2Config.BlockTime(),
		//nolint:gosec // reuse existing config until migration
		RetryBlockThreshold: uint16(fCfg.BumpThreshold()),
		EmptyTxLimitDefault: fCfg.LimitDefault(),
	}
	var c txm.Client
	if txmV2Config.DualBroadcast() != nil && *txmV2Config.DualBroadcast() {
		c = clientwrappers.NewDualBroadcastClient(client, keyStore, txmV2Config.CustomURL())
	} else {
		c = clientwrappers.NewChainClient(client)
	}
	t := txm.NewTxm(lggr, chainID, c, attemptBuilder, inMemoryStoreManager, stuckTxDetector, config, keyStore)
	return txm.NewTxmOrchestrator(lggr, chainID, t, inMemoryStoreManager, fwdMgr, keyStore, attemptBuilder), nil
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
func NewEvmReaper(lggr logger.Logger, store txmgrtypes.TxHistoryReaper[*big.Int], txConfig txmgrtypes.ReaperTransactionsConfig, chainID *big.Int) *Reaper {
	return txmgr.NewReaper(lggr, store, txConfig, chainID)
}

// NewEvmConfirmer instantiates a new EVM confirmer
func NewEvmConfirmer(
	txStore TxStore,
	client TxmClient,
	feeConfig txmgrtypes.ConfirmerFeeConfig,
	txConfig txmgrtypes.ConfirmerTransactionsConfig,
	dbConfig txmgrtypes.ConfirmerDatabaseConfig,
	keystore KeyStore,
	txAttemptBuilder TxAttemptBuilder,
	lggr logger.Logger,
	stuckTxDetector StuckTxDetector,
	metrics metrics.GenericTXMMetrics,
) *Confirmer {
	return txmgr.NewConfirmer(txStore, client, feeConfig, txConfig, dbConfig, keystore, txAttemptBuilder, lggr, func(r *types.Receipt) bool { return r == nil }, stuckTxDetector, metrics)
}

// NewEvmTracker instantiates a new EVM tracker for abandoned transactions
func NewEvmTracker(
	txStore TxStore,
	keyStore KeyStore,
	chainID *big.Int,
	lggr logger.Logger,
) *Tracker {
	return txmgr.NewTracker[*big.Int, common.Address, common.Hash, common.Hash, *types.Receipt](txStore, keyStore, chainID, lggr)
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
	chainType chaintype.ChainType,
	metrics metrics.GenericTXMMetrics,
) *Broadcaster {
	nonceTracker := NewNonceTracker(logger, txStore, client)
	return txmgr.NewBroadcaster(txStore, client, chainConfig, feeConfig, txConfig, listenerConfig, keystore, txAttemptBuilder, nonceTracker, logger, checkerFactory, autoSyncNonce, string(chainType), metrics)
}
