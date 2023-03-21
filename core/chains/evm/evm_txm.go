package evm

import (
	"fmt"

	"github.com/smartcontractkit/sqlx"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	evmconfig "github.com/smartcontractkit/chainlink/core/chains/evm/config"
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
		checker := &txmgr.CheckerFactory{Client: client}
		txm = txmgr.NewTxm[*types.Address, *types.TxHash](
			db,
			client,
			cfg,
			opts.KeyStore,
			opts.EventBroadcaster,
			lggr,
			checker,
			logPoller,
		)
	} else {
		txm = opts.GenTxManager(chainID)
	}
	return txm
}
