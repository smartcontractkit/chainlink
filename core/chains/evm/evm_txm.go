package evm

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/sqlx"

	txmgrtypes "github.com/smartcontractkit/chainlink/common/txmgr/types"
	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	evmconfig "github.com/smartcontractkit/chainlink/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/core/chains/evm/forwarders"
	"github.com/smartcontractkit/chainlink/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
)

func newEvmTxm(
	db *sqlx.DB,
	cfg evmconfig.ChainScopedConfig,
	client evmclient.Client,
	lggr logger.Logger,
	logPoller logpoller.LogPoller,
	opts ChainSetOpts,
) txmgr.TxManager[*types.Address, *types.TxHash] {
	chainID := cfg.ChainID()
	var txm txmgr.TxManager[*types.Address, *types.TxHash]
	if !cfg.EVMRPCEnabled() {
		txm = &txmgr.NullTxManager[*types.Address, *types.TxHash]
		{
		ErrMsg:
			fmt.Sprintf("Ethereum is disabled for chain %d", chainID)
		}
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

		var fwdMgr txmgrtypes.ForwarderManager[common.Address]

		if cfg.EvmUseForwarders() {
			fwdMgr = forwarders.NewFwdMgr(db, client, logPoller, lggr, cfg)
		} else {
			lggr.Info("EvmForwarderManager: Disabled")
		}

		checker := &txmgr.CheckerFactory{Client: client}
		txm = txmgr.NewTxm[*types.Address, *types.TxHash](db, client, cfg, opts.KeyStore, opts.EventBroadcaster, lggr, checker, estimator, fwdMgr)
	} else {
		txm = opts.GenTxManager(chainID)
	}
	return txm
}
