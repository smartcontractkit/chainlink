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
	HistoryDepth    = 128
	PollLogInterval = time.Second
	PollLogLookBack = 300
	Separator       = "|"
	ChannelSize     = 20
)

type BlockKey struct {
	block int64
	hash  common.Hash
}

func (bk *BlockKey) getBlockKey() ocr2keepers.BlockKey {
	return ocr2keepers.BlockKey(fmt.Sprintf("%d%s%s", bk.block, Separator, bk.hash.Hex()))
}

type HeadProvider struct {
	sync                       utils.StartStopOnce
	mu                         sync.RWMutex
	ctx                        context.Context
	cancel                     context.CancelFunc
	ht                         httypes.HeadTracker
	hb                         httypes.HeadBroadcaster
	lp                         logpoller.LogPoller
	headC                      chan int64
	unsubscribe                func()
	latestBlocksFromLogPoller  []BlockKey
	latestBlockFromBroadcaster int64
	subscribers                map[int]chan ocr2keepers.BlockHistory
	maxSubId                   int
	lggr                       logger.Logger
}

func NewHeadProvider(c evm.Chain, lggr logger.Logger) *HeadProvider {
	return &HeadProvider{
		ht:          c.HeadTracker(),
		hb:          c.HeadBroadcaster(),
		lp:          c.LogPoller(),
		headC:       make(chan int64, ChannelSize),
		subscribers: map[int]chan ocr2keepers.BlockHistory{},
		lggr:        lggr.Named("HeadProvider"),
	}
}

func (hw *HeadProvider) Start(ctx context.Context) error {
	return hw.sync.StartOnce("HeadProvider", func() error {
		hw.mu.Lock()
		defer hw.mu.Unlock()
		_, hw.unsubscribe = hw.hb.Subscribe(&headWrapper{c: hw.headC})
		hw.ctx, hw.cancel = context.WithCancel(context.Background())

		// poll from head broadcaster channel and push to subscribers
		{
			go func(ctx context.Context) {
				for {
					select {
					// on fast chain, are we concerned that we will run this too often?
					case bk := <-hw.headC:
						hw.latestBlockFromBroadcaster = bk

						var keys [HistoryDepth]BlockKey
						var idx = HistoryDepth - 1
						hw.mu.Lock()
						// find the latest HistoryDepth blocks which are smaller than latestBlockFromBroadcaster
						length := len(hw.latestBlocksFromLogPoller)
						for i := length - 1; i >= 0; i-- {
							b := hw.latestBlocksFromLogPoller[i]
							if b.block > hw.latestBlockFromBroadcaster {
								continue
							}
							keys[idx] = b
							idx--
							if idx < 0 {
								break
							}
						}
						// does it matter if it doesn't have enough blocks to fill the entire history?

						history := getBlockHistory(keys)
						// send history to all subscribers
						for _, subC := range hw.subscribers {
							subC <- history
						}

						hw.mu.Unlock()
					case <-ctx.Done():
						return
					}
				}
			}(hw.ctx)
		}

		// poll logs from log poller at an interval and
		{
			go func(cx context.Context) {
				ticker := time.NewTicker(PollLogInterval)
				for {
					select {
					case <-ticker.C:
						h, err := hw.lp.LatestBlock(pg.WithParentCtx(ctx))
						if err != nil {
							hw.lggr.Errorf("failed to get latest block", err)
						}

						var blocks []uint64
						for i := 0; i < PollLogLookBack; i++ {
							blocks = append(blocks, uint64(h-int64(i)))
						}

						// request the past LOOK_BACK blocks from log poller
						// returned blocks are in ASC order
						logpollerBlocks, err := hw.lp.GetBlocksRange(cx, blocks, pg.WithParentCtx(ctx))
						var tmp []BlockKey
						hw.mu.Lock()
						for _, b := range logpollerBlocks {
							tmp = append(tmp, BlockKey{
								block: b.BlockNumber,
								hash:  b.BlockHash,
							})
						}
						hw.latestBlocksFromLogPoller = tmp
						hw.mu.Unlock()
						if err != nil {
							hw.lggr.Errorf("failed to get blocks range", err)
						}
					case <-cx.Done():
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
	c chan int64
}

func (w *headWrapper) OnNewLongestChain(_ context.Context, head *evmtypes.Head) {
	var b int64
	if head != nil {
		b = head.Number
	}
	w.c <- b
}

func getBlockHistory(keys [HistoryDepth]BlockKey) ocr2keepers.BlockHistory {
	var blockKeys []ocr2keepers.BlockKey
	for _, k := range keys {
		blockKeys = append(blockKeys, k.getBlockKey())
	}
	return blockKeys
}
