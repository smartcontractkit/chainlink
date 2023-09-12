package evm

import (
	"github.com/smartcontractkit/sqlx"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmconfig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

type TxmCfg struct {
	EVM              evmconfig.EVM
	DB               config.Database
	LogPoller        logpoller.LogPoller
	GasEstimator     gas.EvmFeeEstimator
	KeyStore         keystore.Eth
	EventBroadcaster pg.EventBroadcaster
}

func newEvmTxm(
	db *sqlx.DB,
	cfg TxmCfg,
	client evmclient.Client,
	lggr logger.Logger,
) (txm txmgr.TxManager,
	err error,
) {

	lggr = lggr.Named("Txm")
	lggr.Infow("Initializing EVM transaction manager",
		"bumpTxDepth", cfg.EVM.GasEstimator().BumpTxDepth(),
		"maxInFlightTransactions", cfg.EVM.Transactions().MaxInFlight(),
		"maxQueuedTransactions", cfg.EVM.Transactions().MaxQueued(),
		"nonceAutoSync", cfg.EVM.NonceAutoSync(),
		"limitDefault", cfg.EVM.GasEstimator().LimitDefault(),
	)

	return txmgr.NewTxm(
		db,
		cfg.EVM,
		txmgr.NewEvmTxmFeeConfig(cfg.EVM.GasEstimator()),
		cfg.EVM.Transactions(),
		cfg.DB,
		cfg.DB.Listener(),
		client,
		lggr,
		cfg.LogPoller,
		cfg.KeyStore,
		cfg.EventBroadcaster,
		cfg.GasEstimator)

}
