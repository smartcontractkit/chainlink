package headtracker

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
)

// HeadSaver maintains chains persisted in DB. All methods are thread-safe.
type HeadSaver interface {
	// Save updates the latest block number, if indeed the latest, and persists
	// this number in case of reboot.
	Save(ctx context.Context, head *eth.Head) error
	// LoadFromDB loads latest EvmHeadTrackerHistoryDepth heads, returns the latest chain.
	LoadFromDB(ctx context.Context) (*eth.Head, error)
	// LatestHeadFromDB returns the highest seen head from DB.
	LatestHeadFromDB(ctx context.Context) (*eth.Head, error)
	// LatestChain returns the block header with the highest number that has been seen, or nil.
	LatestChain() *eth.Head
	// Chain returns a head for the specified hash, or nil.
	Chain(hash common.Hash) *eth.Head
}

type headSaver struct {
	orm    ORM
	config Config
	logger logger.Logger
	heads  Heads
}

func NewHeadSaver(lggr logger.Logger, orm ORM, config Config) HeadSaver {
	return &headSaver{
		orm:    orm,
		config: config,
		logger: lggr.Named(logger.HeadSaver),
		heads:  NewHeads(),
	}
}

func (hs *headSaver) Save(ctx context.Context, head *eth.Head) error {
	if err := hs.orm.IdempotentInsertHead(ctx, head); err != nil {
		return err
	}

	historyDepth := uint(hs.config.EvmHeadTrackerHistoryDepth())
	hs.heads.AddHeads(historyDepth, head)

	return hs.orm.TrimOldHeads(ctx, historyDepth)
}

func (hs *headSaver) LoadFromDB(ctx context.Context) (chain *eth.Head, err error) {
	historyDepth := uint(hs.config.EvmHeadTrackerHistoryDepth())
	heads, err := hs.orm.LatestHeads(ctx, historyDepth)
	if err != nil {
		return nil, err
	}

	hs.heads.AddHeads(historyDepth, heads...)
	return hs.heads.LatestHead(), nil
}

func (hs *headSaver) LatestHeadFromDB(ctx context.Context) (head *eth.Head, err error) {
	return hs.orm.LatestHead(ctx)
}

func (hs *headSaver) LatestChain() *eth.Head {
	head := hs.heads.LatestHead()
	if head == nil {
		return nil
	}
	if head.ChainLength() < hs.config.EvmFinalityDepth() {
		hs.logger.Debugw("chain shorter than EvmFinalityDepth", "chainLen", head.ChainLength(), "evmFinalityDepth", hs.config.EvmFinalityDepth())
	}
	return head
}

func (hs *headSaver) Chain(hash common.Hash) *eth.Head {
	return hs.heads.HeadByHash(hash)
}
