package evm

import (
	"context"
	"fmt"
	"sync"
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
	initialBlockSize int64
	lggr             logger.Logger
}

func NewBlockSubscriber(hb httypes.HeadBroadcaster, lp logpoller.LogPoller, blockHistorySize int64, initialBlockSize int64, lggr logger.Logger) *BlockSubscriber {
	return &BlockSubscriber{
		hb:               hb,
		lp:               lp,
		headC:            make(chan *evmtypes.Head, channelSize),
		subscribers:      map[int]chan ocr2keepers.BlockHistory{},
		blocks:           map[int64]string{},
		blockHistorySize: blockHistorySize,
		initialBlockSize: initialBlockSize,
		latestBlock:      atomic.Int64{},
		lggr:             lggr.Named("BlockSubscriber"),
	}
}

func (hw *BlockSubscriber) getBlockRange(ctx context.Context) ([]uint64, error) {
	h, err := hw.lp.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		return nil, err
	}
	hw.lggr.Infof("latest block from log poller is %d", h)

	var blocks []uint64
	for i := hw.initialBlockSize - 1; i >= 0; i-- {
		blocks = append(blocks, uint64(h-i))
	}
	return blocks, nil
}

func (hw *BlockSubscriber) initializeBlocks(blocks []uint64) error {
	logpollerBlocks, err := hw.lp.GetBlocksRange(hw.ctx, blocks, pg.WithParentCtx(hw.ctx))
	if err != nil {
		return err
	}
	for i, b := range logpollerBlocks {
		if i == 0 {
			hw.lastClearedBlock = b.BlockNumber - 1
			hw.lggr.Infof("lastClearedBlock is %d", hw.lastClearedBlock)
		}
		hw.blocks[b.BlockNumber] = b.BlockHash.Hex()
	}
	hw.lggr.Infof("initialize with %d blocks", len(logpollerBlocks))
	return nil
}

func (hw *BlockSubscriber) buildHistory(block int64) ocr2keepers.BlockHistory {
	var keys []ocr2keepers.BlockKey
	// populate keys slice in block DES order
	for i := int64(0); i < hw.blockHistorySize; i++ {
		if h, ok := hw.blocks[block-i]; ok {
			keys = append(keys, ocr2keepers.BlockKey{
				Number: ocr2keepers.BlockNumber(block - i),
				Hash:   common.HexToHash(h),
			})
		} else {
			hw.lggr.Infof("block %d is missing", block-i)
		}
	}
	return keys
}

func (hw *BlockSubscriber) cleanup() {
	hw.mu.Lock()
	defer hw.mu.Unlock()

	hw.lggr.Infof("start clearing blocks from %d to %d", hw.lastClearedBlock+1, hw.lastSentBlock-hw.initialBlockSize)
	for i := hw.lastClearedBlock + 1; i <= hw.lastSentBlock-hw.initialBlockSize; i++ {
		delete(hw.blocks, i)
	}
	hw.lastClearedBlock = hw.lastSentBlock - hw.initialBlockSize
	hw.lggr.Infof("lastClearedBlock is set to %d", hw.lastClearedBlock)
}

func (hw *BlockSubscriber) Start(ctx context.Context) error {
	hw.lggr.Info("block subscriber started.")
	return hw.sync.StartOnce("BlockSubscriber", func() error {
		hw.mu.Lock()
		defer hw.mu.Unlock()
		hw.ctx, hw.cancel = context.WithCancel(ctx)
		// initialize the blocks map with the recent initialBlockSize blocks
		blocks, err := hw.getBlockRange(hw.ctx)
		if err != nil {
			hw.lggr.Errorf("failed to get block range", err)
		}
		err = hw.initializeBlocks(blocks)
		if err != nil {
			hw.lggr.Errorf("failed to get log poller blocks", err)
		}

		_, hw.unsubscribe = hw.hb.Subscribe(&headWrapper{headC: hw.headC})

		// poll from head broadcaster channel and push to subscribers
		{
			go func(ctx context.Context) {
				for {
					select {
					case h := <-hw.headC:
						hw.processHead(h)
					case <-ctx.Done():
						return
					}
				}
			}(hw.ctx)
		}

		// clean up block maps
		{
			go func(ctx context.Context) {
				ticker := time.NewTicker(cleanUpInterval)
				for {
					select {
					case <-ticker.C:
						hw.cleanup()
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

func (hw *BlockSubscriber) Close() error {
	hw.lggr.Info("stop block subscriber")
	return hw.sync.StopOnce("BlockSubscriber", func() error {
		hw.mu.Lock()
		defer hw.mu.Unlock()

		close(hw.headC)
		hw.cancel()
		hw.unsubscribe()
		return nil
	})
}

func (hw *BlockSubscriber) Subscribe() (int, chan ocr2keepers.BlockHistory, error) {
	hw.mu.Lock()
	defer hw.mu.Unlock()

	hw.maxSubId++
	subId := hw.maxSubId
	newC := make(chan ocr2keepers.BlockHistory, channelSize)
	hw.subscribers[subId] = newC
	hw.lggr.Infof("new subscriber %d", subId)

	return subId, newC, nil
}

func (hw *BlockSubscriber) Unsubscribe(subId int) error {
	hw.mu.Lock()
	defer hw.mu.Unlock()

	c, ok := hw.subscribers[subId]
	if !ok {
		return fmt.Errorf("subscriber %d does not exist", subId)
	}

	close(c)
	delete(hw.subscribers, subId)
	hw.lggr.Infof("subscriber %d unsubscribed", subId)
	return nil
}

func (hw *BlockSubscriber) processHead(h *evmtypes.Head) {
	hw.mu.Lock()
	defer hw.mu.Unlock()
	// head parent is a linked list with EVM finality depth
	// when re-org happens, new heads will have pointers to the new blocks
	for cp := h; cp != nil; cp = cp.Parent {
		if cp != h && hw.blocks[cp.Number] != cp.Hash.Hex() {
			hw.lggr.Warnf("overriding block %d old hash %s with new hash %s due to re-org", cp.Number, hw.blocks[cp.Number], cp.Hash.Hex())
		}
		hw.blocks[cp.Number] = cp.Hash.Hex()
	}
	hw.lggr.Debugf("blocks block %d hash is %s", h.Number, h.Hash.Hex())

	history := hw.buildHistory(h.Number)

	hw.latestBlock.Store(h.Number)
	hw.lastSentBlock = h.Number
	// send history to all subscribers
	for _, subC := range hw.subscribers {
		// wrapped in a select to not get blocked by certain subscribers
		select {
		case subC <- history:
		default:
			hw.lggr.Warnf("subscriber channel is full, dropping block history with length %d", len(history))
		}
	}

	hw.lggr.Debugf("published block history with length %d and latestBlock %d to %d subscriber(s)", len(history), hw.latestBlock.Load(), len(hw.subscribers))
}

func (hw *BlockSubscriber) queryBlocksMap(bn int64) (string, bool) {
	hw.mu.RLock()
	defer hw.mu.RUnlock()
	v, ok := hw.blocks[bn]
	return v, ok
}

type headWrapper struct {
	headC chan *evmtypes.Head
}

func (w *headWrapper) OnNewLongestChain(_ context.Context, head *evmtypes.Head) {
	if head != nil {
		w.headC <- head
	}
}
