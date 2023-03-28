package evm

import (
	"fmt"

	"github.com/smartcontractkit/sqlx"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmconfig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/forwarders"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func newEvmTxm(
	db *sqlx.DB,
	cfg evmconfig.ChainScopedConfig,
	client evmclient.Client,
	lggr logger.Logger,
	logPoller logpoller.LogPoller,
	opts ChainSetOpts,
) (txm txmgr.TxManager[*types.Address, *types.TxHash, *types.BlockHash],
	estimator gas.EvmFeeEstimator,
	err error,
) {
	chainID := cfg.ChainID()
	if !cfg.EVMRPCEnabled() {
		txm = &txmgr.NullTxManager[*types.Address, *types.TxHash, *types.BlockHash]{ErrMsg: fmt.Sprintf("Ethereum is disabled for chain %d", chainID)}
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
		var fwdMgr txmgrtypes.ForwarderManager[*types.Address]

		if cfg.EvmUseForwarders() {
			fwdMgr = forwarders.NewFwdMgr(db, client, logPoller, lggr, cfg)
		} else {
			lggr.Info("EvmForwarderManager: Disabled")
		}

		checker := &txmgr.CheckerFactory{Client: client}
		// create tx attempt builder
		txAttemptBuilder := txmgr.NewEvmTxAttemptBuilder(*client.ChainID(), cfg, opts.KeyStore, estimator)
		txStorageService := txmgr.NewTxStorageService(db, lggr, cfg)
		txNonceSyncer := txmgr.NewNonceSyncer(txStorageService, lggr, client, opts.KeyStore)

		addresses, err := opts.KeyStore.EnabledAddressesForChain(client.ChainID())
		if err != nil {
			return nil, nil, err
		}
		ethBroadcaster := txmgr.NewEthBroadcaster(txStorageService, client, cfg, opts.KeyStore, opts.EventBroadcaster, addresses, txAttemptBuilder, txNonceSyncer, lggr, checker, cfg.EvmNonceAutoSync())
		ethConfirmer := txmgr.NewEthConfirmer(txStorageService, client, cfg, opts.KeyStore, addresses, txAttemptBuilder, lggr)
		txm = txmgr.NewTxm(db, client, cfg, opts.KeyStore, opts.EventBroadcaster, lggr, checker, fwdMgr, txAttemptBuilder, txStorageService, txNonceSyncer, *ethBroadcaster, *ethConfirmer)
	} else {
		txm = opts.GenTxManager(chainID)
	}

	return txm, estimator, nil
}
