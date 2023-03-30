package builder

import (
	"github.com/smartcontractkit/sqlx"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/forwarders"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

func BuildNewTxm(
	db *sqlx.DB,
	cfg txmgr.Config,
	client evmclient.Client,
	lggr logger.Logger,
	logPoller logpoller.LogPoller,
	keyStore keystore.Eth,
	eventBroadcaster pg.EventBroadcaster,
	estimator gas.EvmFeeEstimator,
) (txm txmgr.TxManager[*types.Address, *types.TxHash, *types.BlockHash],
	err error,
) {
	var fwdMgr txmgrtypes.ForwarderManager[*types.Address]

	if cfg.EvmUseForwarders() {
		fwdMgr = forwarders.NewFwdMgr(db, client, logPoller, lggr, cfg)
	} else {
		lggr.Info("EvmForwarderManager: Disabled")
	}
	checker := &txmgr.CheckerFactory{Client: client}
	// create tx attempt builder
	txAttemptBuilder := txmgr.NewEvmTxAttemptBuilder(*client.ChainID(), cfg, keyStore, estimator)
	txStorageService := txmgr.NewTxStorageService(db, lggr, cfg)
	txNonceSyncer := txmgr.NewNonceSyncer(txStorageService, lggr, client, keyStore)

	addresses, err := keyStore.EnabledAddressesForChain(client.ChainID())
	if err != nil {
		return nil, err
	}
	ethBroadcaster := txmgr.NewEthBroadcaster(txStorageService, client, cfg, keyStore, eventBroadcaster, addresses, txAttemptBuilder, txNonceSyncer, lggr, checker, cfg.EvmNonceAutoSync())
	ethConfirmer := txmgr.NewEthConfirmer(txStorageService, client, cfg, keyStore, addresses, txAttemptBuilder, lggr)
	txm = txmgr.NewTxm(db, client, cfg, keyStore, eventBroadcaster, lggr, checker, fwdMgr, txAttemptBuilder, txStorageService, txNonceSyncer, *ethBroadcaster, *ethConfirmer)
	return txm, nil
}
