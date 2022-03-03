package headtracker

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	httypes "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/types"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
)

type headSaver struct {
	orm    ORM
	config Config
	logger logger.Logger
	heads  Heads
}

func NewHeadSaver(lggr logger.Logger, orm ORM, config Config) httypes.HeadSaver {
	return &headSaver{
		orm:    orm,
		config: config,
		logger: lggr.Named(logger.HeadSaver),
		heads:  NewHeads(),
	}
}

func (hs *headSaver) Save(ctx context.Context, head *evmtypes.Head) error {
	if err := hs.orm.IdempotentInsertHead(ctx, head); err != nil {
		return err
	}

	historyDepth := uint(hs.config.EvmHeadTrackerHistoryDepth())
	hs.heads.AddHeads(historyDepth, head)

	return hs.orm.TrimOldHeads(ctx, historyDepth)
}

func (hs *headSaver) LoadFromDB(ctx context.Context) (chain *evmtypes.Head, err error) {
	historyDepth := uint(hs.config.EvmHeadTrackerHistoryDepth())
	heads, err := hs.orm.LatestHeads(ctx, historyDepth)
	if err != nil {
		return nil, err
	}

	hs.heads.AddHeads(historyDepth, heads...)
	return hs.heads.LatestHead(), nil
}

func (hs *headSaver) LatestHeadFromDB(ctx context.Context) (head *evmtypes.Head, err error) {
	return hs.orm.LatestHead(ctx)
}

func (hs *headSaver) LatestChain() *evmtypes.Head {
	head := hs.heads.LatestHead()
	if head == nil {
		return nil
	}
	if head.ChainLength() < hs.config.EvmFinalityDepth() {
		hs.logger.Debugw("chain shorter than EvmFinalityDepth", "chainLen", head.ChainLength(), "evmFinalityDepth", hs.config.EvmFinalityDepth())
	}
	return head
}

func (hs *headSaver) Chain(hash common.Hash) *evmtypes.Head {
	return hs.heads.HeadByHash(hash)
}

var NullSaver httypes.HeadSaver = &nullSaver{}

type nullSaver struct{}

func (*nullSaver) Save(ctx context.Context, head *evmtypes.Head) error          { return nil }
func (*nullSaver) LoadFromDB(ctx context.Context) (*evmtypes.Head, error)       { return nil, nil }
func (*nullSaver) LatestHeadFromDB(ctx context.Context) (*evmtypes.Head, error) { return nil, nil }
func (*nullSaver) LatestChain() *evmtypes.Head                                  { return nil }
func (*nullSaver) Chain(hash common.Hash) *evmtypes.Head                        { return nil }
