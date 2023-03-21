package evm

import (
	"fmt"

	"github.com/smartcontractkit/sqlx"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	evmconfig "github.com/smartcontractkit/chainlink/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/core/chains/evm/forwarders"
	"github.com/smartcontractkit/chainlink/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/logger"
)

func newEvmTxm(
	db *sqlx.DB,
	cfg evmconfig.ChainScopedConfig,
	client evmclient.Client,
	lggr logger.Logger,
	logPoller logpoller.LogPoller,
	opts ChainSetOpts,
) txmgr.TxManager {
	chainID := cfg.ChainID()
	var txm txmgr.TxManager
	if !cfg.EVMRPCEnabled() {
		txm = &txmgr.NullTxManager{ErrMsg: fmt.Sprintf("Ethereum is disabled for chain %d", chainID)}
	} else if opts.GenTxManager == nil {
		lggr = lggr.Named("Txm")
		lggr.Infow("Initializing EVM transaction manager",
			"gasBumpTxDepth", cfg.EvmGasBumpTxDepth(),
			"maxInFlightTransactions", cfg.EvmMaxInFlightTransactions(),
			"maxQueuedTransactions", cfg.EvmMaxQueuedTransactions(),
			"nonceAutoSync", cfg.EvmNonceAutoSync(),
			"gasLimitDefault", cfg.EvmGasLimitDefault(),
		)

		// build estimator from factory
		estimator := gas.NewEstimator(lggr, client, cfg)

		// build forwarder manager for evm txm
		var fwdMgr *forwarders.FwdMgr
		if cfg.EvmUseForwarders() {
			fwdMgr = forwarders.NewFwdMgr(db, client, logPoller, lggr, cfg)
		} else {
			lggr.Info("EvmForwarderManager: Disabled")
		}

		// create attempt builder
		attemptBuilder := txmgr.NewEvmAttemptBuilder(*client.ChainID(), cfg, opts.KeyStore, estimator)

		checker := &txmgr.CheckerFactory{Client: client}
		txm = txmgr.NewTxm(db, client, cfg, opts.KeyStore, opts.EventBroadcaster, lggr, checker, fwdMgr, attemptBuilder)
	} else {
		txm = opts.GenTxManager(chainID)
	}
	return txm
}
