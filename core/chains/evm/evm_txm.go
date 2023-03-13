package evm

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/sqlx"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	evmconfig "github.com/smartcontractkit/chainlink/core/chains/evm/config"
	httypes "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
)

// GenericTxManager is type alias for txmgr.TxManager.
// This is necessary because Golang doesn't allow embedding
// txmgr.TxManager directly inside evmTxm struct.
type GenericTxManager = txmgr.TxManager

var _ httypes.HeadTrackable = &evmTxm{}

// evmTxm is an evm wrapper over the generic TxManager interface
type evmTxm struct {
	httypes.HeadTrackable
	GenericTxManager
}

func (e evmTxm) OnNewLongestChain(ctx context.Context, head *evmtypes.Head) {
	e.GenericTxManager.OnNewLongestChain(ctx, head)
}

func newEvmTxm(
	db *sqlx.DB,
	cfg evmconfig.ChainScopedConfig,
	client evmclient.Client,
	lggr logger.Logger,
	logPoller logpoller.LogPoller,
	opts ChainSetOpts,
) *evmTxm {
	chainID := cfg.ChainID()
	var txm txmgr.TxManager
	if !cfg.EVMRPCEnabled() {
		txm = &txmgr.NullTxManager{ErrMsg: fmt.Sprintf("Ethereum is disabled for chain %d", chainID)}
	} else if opts.GenTxManager == nil {
		checker := &txmgr.CheckerFactory{Client: client}
		txm = txmgr.NewTxm(db, client, cfg, opts.KeyStore, opts.EventBroadcaster, lggr, checker, logPoller)
	} else {
		txm = opts.GenTxManager(chainID)
	}
	return &evmTxm{GenericTxManager: txm}
}
