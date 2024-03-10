package headtracker

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink/v2/common/types"
	httypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

type headSaver struct {
	orm      ORM
	config   Config
	htConfig HeadTrackerConfig
	logger   logger.Logger
	heads    Heads
}

var _ commontypes.HeadSaver[*evmtypes.Head, common.Hash] = (*headSaver)(nil)

func NewHeadSaver(lggr logger.Logger, orm ORM, config Config, htConfig HeadTrackerConfig) httypes.HeadSaver {
	return &headSaver{
		orm:      orm,
		config:   config,
		htConfig: htConfig,
		logger:   logger.Named(lggr, "HeadSaver"),
		heads:    NewHeads(),
	}
}

func (hs *headSaver) Save(ctx context.Context, head *evmtypes.Head) error {
	if err := hs.orm.IdempotentInsertHead(ctx, head); err != nil {
		return err
	}

	historyDepth := uint(hs.htConfig.HistoryDepth())
	hs.heads.AddHeads(historyDepth, head)

	return hs.orm.TrimOldHeads(ctx, historyDepth)
}

func (hs *headSaver) Load(ctx context.Context) (chain *evmtypes.Head, err error) {
	historyDepth := uint(hs.htConfig.HistoryDepth())
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
	if head.ChainLength() < hs.config.FinalityDepth() {
		hs.logger.Debugw("chain shorter than FinalityDepth", "chainLen", head.ChainLength(), "evmFinalityDepth", hs.config.FinalityDepth())
	}
	return head
}

func (hs *headSaver) Chain(hash common.Hash) *evmtypes.Head {
	return hs.heads.HeadByHash(hash)
}

var NullSaver httypes.HeadSaver = &nullSaver{}

type nullSaver struct{}

func (*nullSaver) Save(ctx context.Context, head *evmtypes.Head) error          { return nil }
func (*nullSaver) Load(ctx context.Context) (*evmtypes.Head, error)             { return nil, nil }
func (*nullSaver) LatestHeadFromDB(ctx context.Context) (*evmtypes.Head, error) { return nil, nil }
func (*nullSaver) LatestChain() *evmtypes.Head                                  { return nil }
func (*nullSaver) Chain(hash common.Hash) *evmtypes.Head                        { return nil }
