package legacyevm

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"

	evmclient "github.com/smartcontractkit/chainlink-evm/pkg/client"
	evmconfig "github.com/smartcontractkit/chainlink-evm/pkg/config"
	"github.com/smartcontractkit/chainlink-evm/pkg/gas"
	"github.com/smartcontractkit/chainlink-evm/pkg/gas/rollups"
	evmheads "github.com/smartcontractkit/chainlink-evm/pkg/heads"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
)

func newEvmTxm(
	ds sqlutil.DataSource,
	cfg evmconfig.EVM,
	databaseConfig txmgr.DatabaseConfig,
	listenerConfig txmgr.ListenerConfig,
	client evmclient.Client,
	lggr logger.Logger,
	logPoller logpoller.LogPoller,
	opts ChainRelayOpts,
	headTracker evmheads.Tracker,
	estimator gas.EvmFeeEstimator,
) (txm txmgr.TxManager,
	err error,
) {
	chainID := cfg.ChainID()

	lggr = logger.Named(lggr, "Txm")
	lggr.Infow("Initializing EVM transaction manager",
		"bumpTxDepth", cfg.GasEstimator().BumpTxDepth(),
		"maxInFlightTransactions", cfg.Transactions().MaxInFlight(),
		"maxQueuedTransactions", cfg.Transactions().MaxQueued(),
		"nonceAutoSync", cfg.NonceAutoSync(),
		"limitDefault", cfg.GasEstimator().LimitDefault(),
	)

	if opts.GenTxManager == nil {
		var txmv2 txmgr.TxManager
		if cfg.Transactions().TransactionManagerV2().Enabled() {
			txmv2, err = txmgr.NewTxmV2(
				ds,
				cfg,
				txmgr.NewEvmTxmFeeConfig(cfg.GasEstimator()),
				cfg.Transactions(),
				cfg.Transactions().TransactionManagerV2(),
				client,
				lggr,
				logPoller,
				opts.KeyStore,
				estimator,
			)
			if cfg.Transactions().TransactionManagerV2().DualBroadcast() == nil || !*cfg.Transactions().TransactionManagerV2().DualBroadcast() {
				return txmv2, err
			}
		}
		txm, err = txmgr.NewTxm(
			ds,
			cfg,
			txmgr.NewEvmTxmFeeConfig(cfg.GasEstimator()),
			cfg.Transactions(),
			cfg.NodePool().Errors(),
			databaseConfig,
			listenerConfig,
			client,
			lggr,
			logPoller,
			opts.KeyStore,
			estimator,
			headTracker,
			txmv2)
	} else {
		txm = opts.GenTxManager(chainID)
	}
	return
}

func newGasEstimator(
	cfg evmconfig.EVM,
	client evmclient.Client,
	lggr logger.Logger,
	opts ChainRelayOpts,
	clientsByChainID map[string]rollups.DAClient,
) (estimator gas.EvmFeeEstimator, err error) {
	lggr = logger.Named(lggr, "Txm")
	chainID := cfg.ChainID()
	// build estimator from factory
	if opts.GenGasEstimator == nil {
		if estimator, err = gas.NewEstimator(lggr, client, cfg.ChainType(), chainID, cfg.GasEstimator(), clientsByChainID); err != nil {
			return nil, fmt.Errorf("failed to initialize estimator: %w", err)
		}
	} else {
		estimator = opts.GenGasEstimator(chainID)
	}
	return
}
