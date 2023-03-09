package evm

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/sqlx"

	txmgrtypes "github.com/smartcontractkit/chainlink/common/txmgr/types"
	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	evmconfig "github.com/smartcontractkit/chainlink/core/chains/evm/config"
	httypes "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var _ httypes.HeadTrackable = &evmTxm{}

// EVM specific wrapper to hold the core TxMgr object underneath
type evmTxm struct {
	httypes.HeadTrackable
	services.ServiceCtx
	utils.StartStopOnce

	// core txm object being wrapped
	txm txmgr.TxManager[*evmtypes.Head]
}

func (e *evmTxm) OnNewLongestChain(ctx context.Context, evmHead *evmtypes.Head) {
	e.txm.OnNewLongestChain(ctx, NewHeadViewImpl(evmHead))
}

func (e *evmTxm) Start(ctx context.Context) (err error) {
	return e.txm.Start(ctx)
}

func (e *evmTxm) Close() error {
	return e.txm.Close()
}

func (e *evmTxm) Ready() error {
	return e.txm.Ready()
}

func (e *evmTxm) Healthy() error {
	return e.txm.Healthy()
}

func newEvmTxm(
	db *sqlx.DB,
	cfg evmconfig.ChainScopedConfig,
	client evmclient.Client,
	lggr logger.Logger,
	logPoller logpoller.LogPoller,
	opts ChainSetOpts,
) evmTxm {
	chainID := cfg.ChainID()
	var txm txmgr.TxManager[*evmtypes.Head]
	if !cfg.EVMRPCEnabled() {
		txm = &txmgr.NullTxManager[*evmtypes.Head]{ErrMsg: fmt.Sprintf("Ethereum is disabled for chain %d", chainID)}
	} else if opts.GenTxManager == nil {
		checker := &txmgr.CheckerFactory{Client: client}
		txm = txmgr.NewTxm(db, client, cfg, opts.KeyStore, opts.EventBroadcaster, lggr, checker, logPoller)
	} else {
		txm = opts.GenTxManager(chainID)
	}
	return evmTxm{txm: txm}
}

var _ txmgrtypes.HeadView[*evmtypes.Head] = &headViewImpl{}

// Evm implementation for the generic HeadView interface
type headViewImpl struct {
	// txmgrtypes.HeadView[*evmtypes.Head]
	evmHead *evmtypes.Head
}

func NewHeadViewImpl(head *evmtypes.Head) txmgrtypes.HeadView[*evmtypes.Head] {
	return &headViewImpl{evmHead: head}
}

func (head *headViewImpl) BlockNumber() int64 {
	return head.evmHead.Number
}

// ChainLength returns the length of the chain followed by recursively looking up parents
func (head *headViewImpl) ChainLength() uint32 {
	return head.evmHead.ChainLength()
}

// EarliestInChain recurses through parents until it finds the earliest one
func (head *headViewImpl) EarliestInChain() txmgrtypes.HeadView[*evmtypes.Head] {
	return NewHeadViewImpl(head.evmHead.EarliestInChain())
}

func (head *headViewImpl) Hash() common.Hash {
	return head.evmHead.Hash
}

func (head *headViewImpl) Parent() txmgrtypes.HeadView[*evmtypes.Head] {
	if head.evmHead.Parent == nil {
		return nil
	}

	return NewHeadViewImpl(head.evmHead.Parent)
}

// HashAtHeight returns the hash of the block at the given height, if it is in the chain.
// If not in chain, returns the zero hash
func (head *headViewImpl) HashAtHeight(blockNum int64) common.Hash {
	return head.evmHead.HashAtHeight(blockNum)
}

func (head *headViewImpl) GetNativeHead() *evmtypes.Head {
	return head.evmHead
}
