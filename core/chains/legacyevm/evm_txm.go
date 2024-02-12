package legacyevm

import (
	"github.com/jmoiron/sqlx"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmconfig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	optimistictxm "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr/optimistictxm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func newEvmTxm(
	db *sqlx.DB,
	cfg evmconfig.EVM,
	evmRPCEnabled bool,
	databaseConfig txmgr.DatabaseConfig,
	listenerConfig txmgr.ListenerConfig,
	client evmclient.Client,
	lggr logger.Logger,
	logPoller logpoller.LogPoller,
	opts ChainRelayExtenderConfig,
) (txm optimistictxm.TxManager,
	estimator gas.EvmFeeEstimator,
	err error,
) {
	chainID := cfg.ChainID()

	lggr = lggr.Named("Txm")
	lggr.Infow("Initializing EVM transaction manager",
		"bumpTxDepth", cfg.GasEstimator().BumpTxDepth(),
		"maxInFlightTransactions", cfg.Transactions().MaxInFlight(),
		"maxQueuedTransactions", cfg.Transactions().MaxQueued(),
		"nonceAutoSync", cfg.NonceAutoSync(),
		"limitDefault", cfg.GasEstimator().LimitDefault(),
	)

	// build estimator from factory
	if opts.GenGasEstimator == nil {
		estimator = gas.NewEstimator(lggr, client, cfg, cfg.GasEstimator())
	} else {
		estimator = opts.GenGasEstimator(chainID)
	}

	if opts.GenTxManager == nil {
		txm, err = txmgr.NewOptimisticTxm(
			db,
			cfg,
			txmgr.NewEvmTxmFeeConfig(cfg.GasEstimator()),
			cfg.Transactions(),
			databaseConfig,
			listenerConfig,
			client,
			lggr,
			logPoller,
			opts.KeyStore,
			estimator)
	}
	return
}
