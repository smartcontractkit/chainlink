package evm

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	httypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	// PollLogInterval is the interval to query log poller
	PollLogInterval = time.Second
	// CleanUpInterval is the interval for cleaning up block maps
	CleanUpInterval = 15 * time.Minute
	// Separator is the separator for block key
	Separator = "|"
	// ChannelSize represents the channel size for head broadcaster
	ChannelSize = 20
)

type BlockKey struct {
	block int64
	hash  common.Hash
}

func (bk *BlockKey) getBlockKey() ocr2keepers.BlockKey {
	return ocr2keepers.BlockKey(fmt.Sprintf("%d%s%s", bk.block, Separator, bk.hash.Hex()))
}

type HeadProvider struct {
	sync                  utils.StartStopOnce
	mu                    sync.RWMutex
	ctx                   context.Context
	cancel                context.CancelFunc
	hb                    httypes.HeadBroadcaster
	lp                    logpoller.LogPoller
	headC                 chan BlockKey
	unsubscribe           func()
	subscribers           map[int]chan ocr2keepers.BlockHistory
	blocksFromPoller      map[int64]common.Hash
	blocksFromBroadcaster map[int64]common.Hash
	maxSubId              int
	lastClearedBlock      int64
	lastSentBlock         int64
	blockHistorySize      int64
	lggr                  logger.Logger
}

func NewHeadProvider(c evm.Chain, blockHistorySize int64, lggr logger.Logger) *HeadProvider {
	return &HeadProvider{
		hb:                    c.HeadBroadcaster(),
		lp:                    c.LogPoller(),
		headC:                 make(chan BlockKey, ChannelSize),
		subscribers:           map[int]chan ocr2keepers.BlockHistory{},
		blocksFromPoller:      map[int64]common.Hash{},
		blocksFromBroadcaster: map[int64]common.Hash{},
		blockHistorySize:      blockHistorySize,
		lggr:                  lggr.Named("HeadProvider"),
	}
}

func (hw *HeadProvider) Start(_ context.Context) error {
	hw.lggr.Info("Head Provider started.")
	return hw.sync.StartOnce("HeadProvider", func() error {
		hw.mu.Lock()
		defer hw.mu.Unlock()
		_, hw.unsubscribe = hw.hb.Subscribe(&headWrapper{headC: hw.headC})
		hw.ctx, hw.cancel = context.WithCancel(context.Background())

		// poll from head broadcaster channel and push to subscribers
		{
			go func(ctx context.Context) {
				for {
					select {
					case bk := <-hw.headC:
						hw.mu.Lock()
						if hw.lastClearedBlock == 0 {
							hw.lastClearedBlock = bk.block - 1
							hw.lggr.Infof("lastClearedBlock is %d", hw.lastClearedBlock)
						}
						hw.blocksFromBroadcaster[bk.block] = bk.hash
						hw.lggr.Infof("blocksFromBroadcaster block %d hash is %s", bk.block, bk.hash.String())

						var keys []BlockKey
						// populate keys slice in block DES order
						for i := int64(0); i < hw.blockHistorySize; i-- {
							if h1, ok1 := hw.blocksFromPoller[bk.block-i]; ok1 {
								// if a block exists in log poller, use block data from log poller
								keys = append(keys, BlockKey{
									block: bk.block - i,
									hash:  h1,
								})
							} else if h2, ok2 := hw.blocksFromBroadcaster[bk.block-i]; ok2 {
								// if a block only exists in broadcaster, use data from broadcaster
								keys = append(keys, BlockKey{
									block: bk.block - i,
									hash:  h2,
								})
							}
							hw.lggr.Infof("block %d is missing", bk.block-i)
							// if a block does not exist in both log poller and broadcaster, skip
						}

						hw.lastSentBlock = bk.block
						hw.lggr.Infof("lastSentBlock is %d", hw.lastSentBlock)
						history := getBlockHistory(keys)
						// send history to all subscribers
						for _, subC := range hw.subscribers {
							subC <- history
						}
						hw.lggr.Infof("published block history with length %d to %d subscriber(s)", len(history), len(hw.subscribers))

						hw.mu.Unlock()
					case <-ctx.Done():
						return
					}
				}
			}(hw.ctx)
		}

		// poll logs from log poller at an interval and update block map
		{
			go func(ctx context.Context) {
				ticker := time.NewTicker(PollLogInterval)
				for {
					select {
					case <-ticker.C:
						h, err := hw.lp.LatestBlock(pg.WithParentCtx(ctx))
						hw.lggr.Infof("latest block from log poller is %d", h)
						if err != nil {
							hw.lggr.Errorf("failed to get latest block", err)
						}

						var blocks []uint64
						for i := int64(0); i < hw.blockHistorySize; i++ {
							blocks = append(blocks, uint64(h-i))
						}

						// request the past LOOK_BACK blocksFromPoller from log poller
						// returned blocksFromPoller are in ASC order
						logpollerBlocks, err := hw.lp.GetBlocksRange(ctx, blocks, pg.WithParentCtx(ctx))
						hw.mu.Lock()
						for _, b := range logpollerBlocks {
							hw.blocksFromPoller[b.BlockNumber] = b.BlockHash
						}
						hw.mu.Unlock()
						if err != nil {
							hw.lggr.Errorf("failed to get blocksFromPoller range", err)
						}
					case <-ctx.Done():
						ticker.Stop()
						return
					}
				}
			}(hw.ctx)
		}

		// clean up block maps
		{
			go func(ctx context.Context) {
				ticker := time.NewTicker(CleanUpInterval)
				for {
					select {
					case <-ticker.C:
						hw.mu.Lock()
						hw.lggr.Infof("start clearing blocks from %d to %d", hw.lastClearedBlock+1, hw.lastSentBlock-hw.blockHistorySize)
						for i := hw.lastClearedBlock + 1; i <= hw.lastSentBlock-hw.blockHistorySize; i++ {
							delete(hw.blocksFromPoller, i)
							delete(hw.blocksFromBroadcaster, i)
						}
						hw.lastClearedBlock = hw.lastSentBlock - hw.blockHistorySize
						hw.lggr.Infof("lastClearedBlock is set to %d", hw.lastClearedBlock)
						hw.mu.Unlock()
					case <-ctx.Done():
						ticker.Stop()
						return
					}
				}
			}(hw.ctx)
		}

		return nil
	})
}

func (hw *HeadProvider) Close() error {
	return hw.sync.StopOnce("HeadProvider", func() error {
		hw.mu.Lock()
		defer hw.mu.Unlock()

		close(hw.headC)
		hw.cancel()
		hw.unsubscribe()
		return nil
	})
}

func (hw *HeadProvider) Subscribe() (int, chan ocr2keepers.BlockHistory, error) {
	hw.mu.Lock()
	defer hw.mu.Unlock()

	hw.maxSubId++
	subId := hw.maxSubId
	newC := make(chan ocr2keepers.BlockHistory, ChannelSize)
	hw.subscribers[subId] = newC

	return subId, newC, nil
}

func (hw *HeadProvider) Unsubscribe(subId int) error {
	hw.mu.Lock()
	defer hw.mu.Unlock()

	c, ok := hw.subscribers[subId]
	if !ok {
		return fmt.Errorf("subscriber %d does not exist", subId)
	}

	close(c)
	delete(hw.subscribers, subId)
	return nil
}

type headWrapper struct {
	headC chan BlockKey
}

func (w *headWrapper) OnNewLongestChain(_ context.Context, head *evmtypes.Head) {
	if head != nil {
		w.headC <- BlockKey{
			block: head.Number,
			hash:  head.BlockHash(),
		}
	}
}

func getBlockHistory(keys []BlockKey) ocr2keepers.BlockHistory {
	var blockKeys []ocr2keepers.BlockKey
	for _, k := range keys {
		blockKeys = append(blockKeys, k.getBlockKey())
	}
	return blockKeys
}
