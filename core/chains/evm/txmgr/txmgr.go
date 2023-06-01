package txmgr

import (
	"github.com/smartcontractkit/sqlx"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/forwarders"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

func NewTxm(
	db *sqlx.DB,
	cfg Config,
	client evmclient.Client,
	lggr logger.Logger,
	logPoller logpoller.LogPoller,
	keyStore keystore.Eth,
	eventBroadcaster pg.EventBroadcaster,
	estimator gas.EvmFeeEstimator,
) (txm EvmTxManager,
	err error,
) {
	var fwdMgr EvmFwdMgr

	if cfg.EvmUseForwarders() {
		fwdMgr = forwarders.NewFwdMgr(db, client, logPoller, lggr, cfg)
	} else {
		lggr.Info("EvmForwarderManager: Disabled")
	}
	checker := &CheckerFactory{Client: client}
	// create tx attempt builder
	txAttemptBuilder := NewEvmTxAttemptBuilder(*client.ConfiguredChainID(), cfg, keyStore, estimator)
	txStore := NewTxStore(db, lggr, cfg)
	txNonceSyncer := NewNonceSyncer(txStore, lggr, client, keyStore)

	txmCfg := NewEvmTxmConfig(cfg)       // wrap Evm specific config
	txmClient := NewEvmTxmClient(client) // wrap Evm specific client
	ethBroadcaster := NewEvmBroadcaster(txStore, txmClient, txmCfg, keyStore, eventBroadcaster, txAttemptBuilder, txNonceSyncer, lggr, checker, cfg.EvmNonceAutoSync())
	ethConfirmer := NewEvmConfirmer(txStore, txmClient, txmCfg, keyStore, txAttemptBuilder, lggr)
	var ethResender *EvmResender
	if cfg.EthTxResendAfterThreshold() > 0 {
		ethResender = NewEvmResender(lggr, txStore, txmClient, keyStore, txmgr.DefaultResenderPollInterval, txmCfg)
	}
	txm = NewEvmTxm(txmClient.ConfiguredChainID(), txmCfg, keyStore, lggr, checker, fwdMgr, txAttemptBuilder, txStore, txNonceSyncer, ethBroadcaster, ethConfirmer, ethResender)
	return txm, nil
}
