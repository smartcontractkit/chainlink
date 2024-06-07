package headtracker

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/common/headtracker"
	commontypes "github.com/smartcontractkit/chainlink/v2/common/headtracker/types"
	httypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

type headSaver struct {
	orm      ORM
	config   commontypes.Config
	htConfig commontypes.HeadTrackerConfig
	logger   logger.Logger
	heads    Heads
}

var _ headtracker.HeadSaver[*evmtypes.Head, common.Hash] = (*headSaver)(nil)

func NewHeadSaver(lggr logger.Logger, orm ORM, config commontypes.Config, htConfig commontypes.HeadTrackerConfig) httypes.HeadSaver {
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

	hs.heads.AddHeads(head)

	return nil
}

func (hs *headSaver) Load(ctx context.Context, latestFinalized int64) (chain *evmtypes.Head, err error) {
	minBlockNumber := hs.calculateMinBlockToKeep(latestFinalized)
	heads, err := hs.orm.LatestHeads(ctx, minBlockNumber)
	if err != nil {
		return nil, err
	}

	hs.heads.AddHeads(heads...)
	return hs.heads.LatestHead(), nil
}

func (hs *headSaver) calculateMinBlockToKeep(latestFinalized int64) int64 {
	return max(latestFinalized-int64(hs.htConfig.HistoryDepth()), 0)
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

func (hs *headSaver) MarkFinalized(ctx context.Context, finalized *evmtypes.Head) error {
	minBlockToKeep := hs.calculateMinBlockToKeep(finalized.BlockNumber())
	if !hs.heads.MarkFinalized(finalized.BlockHash(), minBlockToKeep) {
		return fmt.Errorf("failed to find %s block in the canonical chain to mark it as finalized", finalized)
	}

	return hs.orm.TrimOldHeads(ctx, minBlockToKeep)
}

var NullSaver httypes.HeadSaver = &nullSaver{}

type nullSaver struct{}

func (*nullSaver) Save(ctx context.Context, head *evmtypes.Head) error { return nil }
func (*nullSaver) Load(ctx context.Context, latestFinalized int64) (*evmtypes.Head, error) {
	return nil, nil
}
func (*nullSaver) LatestHeadFromDB(ctx context.Context) (*evmtypes.Head, error) { return nil, nil }
func (*nullSaver) LatestChain() *evmtypes.Head                                  { return nil }
func (*nullSaver) Chain(hash common.Hash) *evmtypes.Head                        { return nil }
func (*nullSaver) MarkFinalized(ctx context.Context, latestFinalized *evmtypes.Head) error {
	return nil
}
