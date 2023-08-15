package evm

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"

	httypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	// cleanUpInterval is the interval for cleaning up block maps
	cleanUpInterval = 15 * time.Minute
	// channelSize represents the channel size for head broadcaster
	channelSize = 100
	// lookbackDepth decides valid trigger block lookback range
	lookbackDepth = 1024
	// blockHistorySize decides the block history size
	blockHistorySize = int64(128)
	// cleanUpInterval is the interval for cleaning up block maps
	cleanUpInterval = 15 * time.Minute
	// channelSize represents the channel size for head broadcaster
	channelSize = 100
	// lookbackDepth decides valid trigger block lookback range
	lookbackDepth = 1024
	// blockHistorySize decides the block history size
	blockHistorySize = int64(128)
)

type BlockSubscriber struct {
	sync             utils.StartStopOnce
	mu               sync.RWMutex
	ctx              context.Context
	cancel           context.CancelFunc
	hb               httypes.HeadBroadcaster
	lp               logpoller.LogPoller
	headC            chan *evmtypes.Head
	unsubscribe      func()
	subscribers      map[int]chan ocr2keepers.BlockHistory
	blocks           map[int64]string
	maxSubId         int
	lastClearedBlock int64
	lastSentBlock    int64
	latestBlock      atomic.Int64
	blockHistorySize int64
	blockSize        int64
	lggr             logger.Logger
	sync             utils.StartStopOnce
	mu               sync.RWMutex
	ctx              context.Context
	cancel           context.CancelFunc
	hb               httypes.HeadBroadcaster
	lp               logpoller.LogPoller
	headC            chan *evmtypes.Head
	unsubscribe      func()
	subscribers      map[int]chan ocr2keepers.BlockHistory
	blocks           map[int64]string
	maxSubId         int
	lastClearedBlock int64
	lastSentBlock    int64
	latestBlock      atomic.Int64
	blockHistorySize int64
	blockSize        int64
	lggr             logger.Logger
}

func NewBlockSubscriber(hb httypes.HeadBroadcaster, lp logpoller.LogPoller, lggr logger.Logger) *BlockSubscriber {
func NewBlockSubscriber(hb httypes.HeadBroadcaster, lp logpoller.LogPoller, lggr logger.Logger) *BlockSubscriber {
	return &BlockSubscriber{
		hb:               hb,
		lp:               lp,
		headC:            make(chan *evmtypes.Head, channelSize),
		subscribers:      map[int]chan ocr2keepers.BlockHistory{},
		blocks:           map[int64]string{},
		blockHistorySize: blockHistorySize,
		blockSize:        lookbackDepth,
		latestBlock:      atomic.Int64{},
		lggr:             lggr.Named("BlockSubscriber"),
	}
}

func (bs *BlockSubscriber) getBlockRange(ctx context.Context) ([]uint64, error) {
	h, err := bs.lp.LatestBlock(pg.WithParentCtx(ctx))
func (bs *BlockSubscriber) getBlockRange(ctx context.Context) ([]uint64, error) {
	h, err := bs.lp.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		return nil, err
	}
	bs.lggr.Infof("latest block from log poller is %d", h)
	bs.lggr.Infof("latest block from log poller is %d", h)

	var blocks []uint64
	for i := bs.blockSize - 1; i >= 0; i-- {
		if h-i > 0 {
			blocks = append(blocks, uint64(h-i))
		}
	}
	for i := bs.blockSize - 1; i >= 0; i-- {
		if h-i > 0 {
			blocks = append(blocks, uint64(h-i))
		}
	}
	return blocks, nil
}

func (bs *BlockSubscriber) initializeBlocks(blocks []uint64) error {
	logpollerBlocks, err := bs.lp.GetBlocksRange(bs.ctx, blocks, pg.WithParentCtx(bs.ctx))
func (bs *BlockSubscriber) initializeBlocks(blocks []uint64) error {
	logpollerBlocks, err := bs.lp.GetBlocksRange(bs.ctx, blocks, pg.WithParentCtx(bs.ctx))
	if err != nil {
		return err
	}
	for i, b := range logpollerBlocks {
		if i == 0 {
			bs.lastClearedBlock = b.BlockNumber - 1
			bs.lggr.Infof("lastClearedBlock is %d", bs.lastClearedBlock)
		}
		bs.blocks[b.BlockNumber] = b.BlockHash.Hex()
	for i, b := range logpollerBlocks {
		if i == 0 {
			bs.lastClearedBlock = b.BlockNumber - 1
			bs.lggr.Infof("lastClearedBlock is %d", bs.lastClearedBlock)
		}
		bs.blocks[b.BlockNumber] = b.BlockHash.Hex()
	}
	bs.lggr.Infof("initialize with %d blocks", len(logpollerBlocks))
	bs.lggr.Infof("initialize with %d blocks", len(logpollerBlocks))
	return nil
}

func (bs *BlockSubscriber) buildHistory(block int64) ocr2keepers.BlockHistory {
	var keys []ocr2keepers.BlockKey
func (bs *BlockSubscriber) buildHistory(block int64) ocr2keepers.BlockHistory {
	var keys []ocr2keepers.BlockKey
	// populate keys slice in block DES order
	for i := int64(0); i < bs.blockHistorySize; i++ {
		if block-i > 0 {
			if h, ok := bs.blocks[block-i]; ok {
				keys = append(keys, ocr2keepers.BlockKey{
					Number: ocr2keepers.BlockNumber(block - i),
					Hash:   common.HexToHash(h),
				})
			} else {
				bs.lggr.Debugf("block %d is missing", block-i)
			}
		}
	}
	return keys
}

func (bs *BlockSubscriber) cleanup() {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	hw.lggr.Infof("start clearing blocks from %d to %d", hw.lastClearedBlock+1, hw.lastSentBlock-hw.blockHistorySize)
	for i := hw.lastClearedBlock + 1; i <= hw.lastSentBlock-hw.blockHistorySize; i++ {
		delete(hw.blocksFromPoller, i)
		delete(hw.blocksFromBroadcaster, i)
	}
	hw.lastClearedBlock = hw.lastSentBlock - hw.blockHistorySize
	hw.lggr.Infof("lastClearedBlock is set to %d", hw.lastClearedBlock)
}

func (hw *BlockSubscriber) Start(_ context.Context) error {
	hw.lggr.Info("block subscriber started.")
	return hw.sync.StartOnce("BlockSubscriber", func() error {
		hw.mu.Lock()
		defer hw.mu.Unlock()
		_, hw.unsubscribe = hw.hb.Subscribe(&headWrapper{headC: hw.headC})
		hw.ctx, hw.cancel = context.WithCancel(context.Background())

		// poll from head broadcaster channel and push to subscribers
		{
			go func(ctx context.Context) {
				for {
					select {
					case h := <-bs.headC:
						bs.processHead(h)
					case <-ctx.Done():
						return
					}
				}
			}(bs.ctx)
		}

		// clean up block maps
		{
			go func(ctx context.Context) {
				ticker := time.NewTicker(cleanUpInterval)
				for {
					select {
					case <-ticker.C:
						bs.cleanup()
					case <-ctx.Done():
						ticker.Stop()
						return
					}
				}
			}(bs.ctx)
		}

		return nil
	})
}

func (bs *BlockSubscriber) Close() error {
	bs.lggr.Info("stop block subscriber")
	return bs.sync.StopOnce("BlockSubscriber", func() error {
		bs.mu.Lock()
		defer bs.mu.Unlock()

		close(hw.headC)
		hw.cancel()
		hw.unsubscribe()
		return nil
	})
}

func (bs *BlockSubscriber) Subscribe() (int, chan ocr2keepers.BlockHistory, error) {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	bs.maxSubId++
	subId := bs.maxSubId
	newC := make(chan ocr2keepers.BlockHistory, channelSize)
	bs.subscribers[subId] = newC
	bs.lggr.Infof("new subscriber %d", subId)

	return subId, newC, nil
}

func (bs *BlockSubscriber) Unsubscribe(subId int) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	c, ok := bs.subscribers[subId]
	if !ok {
		return fmt.Errorf("subscriber %d does not exist", subId)
	}

	close(c)
	delete(bs.subscribers, subId)
	bs.lggr.Infof("subscriber %d unsubscribed", subId)
	return nil
}

type headWrapper struct {
	headC chan BlockKey
}

func (w *headWrapper) OnNewLongestChain(_ context.Context, head *evmtypes.Head) {
	w.lggr.Debugf("OnNewLongestChain called with new head %+v", head)

	if head != nil {
		w.headC <- BlockKey{
			block: head.Number,
			hash:  head.BlockHash(),
		}

	}
}
