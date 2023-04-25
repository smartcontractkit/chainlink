package evm

import (
	"fmt"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/builder"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmconfig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func newEvmTxm(
	db *sqlx.DB,
	cfg evmconfig.ChainScopedConfig,
	client evmclient.Client,
	lggr logger.Logger,
	logPoller logpoller.LogPoller,
	opts ChainSetOpts,
) (txm txmgr.EvmTxManager,
	estimator gas.EvmFeeEstimator,
	err error,
) {
	chainID := cfg.ChainID()
	if !cfg.EVMRPCEnabled() {
		txm = &txmgr.NullEvmTxManager{ErrMsg: fmt.Sprintf("Ethereum is disabled for chain %d", chainID)}
		return txm, nil, nil
	}

	lggr = lggr.Named("Txm")
	lggr.Infow("Initializing EVM transaction manager",
		"gasBumpTxDepth", cfg.EvmGasBumpTxDepth(),
		"maxInFlightTransactions", cfg.EvmMaxInFlightTransactions(),
		"maxQueuedTransactions", cfg.EvmMaxQueuedTransactions(),
		"nonceAutoSync", cfg.EvmNonceAutoSync(),
		"gasLimitDefault", cfg.EvmGasLimitDefault(),
	)

	// build estimator from factory
	if opts.GenGasEstimator == nil {
		estimator = gas.NewEstimator(lggr, client, cfg)
	} else {
		estimator = opts.GenGasEstimator(chainID)
	}

	if opts.GenTxManager == nil {
		txm, err = builder.NewTxm(
			db,
			cfg,
			client,
			lggr,
			logPoller,
			opts.KeyStore,
			opts.EventBroadcaster,
			estimator)
	} else {
		txm = opts.GenTxManager(chainID)
	}
	return
}
