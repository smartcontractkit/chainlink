package evm

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"

	"github.com/smartcontractkit/chainlink-evm/pkg/heads"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink-evm/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	// cleanUpInterval is the interval for cleaning up block maps
	cleanUpInterval = 15 * time.Minute
	// channelSize represents the channel size for head broadcaster
	channelSize = 100
	// lookbackDepth decides valid trigger block lookback range
	lookbackDepth = 1024
	// blockHistorySize decides the block history size sent to subscribers
	blockHistorySize = int64(256)
)

var (
	BlockSubscriberServiceName = "BlockSubscriber"
)

type BlockSubscriber struct {
	services.StateMachine
	threadCtrl utils.ThreadControl

	mu               sync.RWMutex
	hb               heads.Broadcaster
	lp               logpoller.LogPoller
	headC            chan *types.Head
	unsubscribe      func()
	subscribers      map[int]chan ocr2keepers.BlockHistory
	blocks           map[int64]string
	maxSubId         int
	lastClearedBlock int64
	lastSentBlock    int64
	latestBlock      atomic.Pointer[ocr2keepers.BlockKey]
	blockHistorySize int64
	blockSize        int64
	finalityDepth    uint32
	lggr             logger.Logger
}

func (bs *BlockSubscriber) LatestBlock() *ocr2keepers.BlockKey {
	return bs.latestBlock.Load()
}

var _ ocr2keepers.BlockSubscriber = &BlockSubscriber{}

func NewBlockSubscriber(hb heads.Broadcaster, lp logpoller.LogPoller, finalityDepth uint32, lggr logger.Logger) *BlockSubscriber {
	return &BlockSubscriber{
		threadCtrl:       utils.NewThreadControl(),
		hb:               hb,
		lp:               lp,
		headC:            make(chan *types.Head, channelSize),
		subscribers:      map[int]chan ocr2keepers.BlockHistory{},
		blocks:           map[int64]string{},
		blockHistorySize: blockHistorySize,
		blockSize:        lookbackDepth,
		finalityDepth:    finalityDepth,
		latestBlock:      atomic.Pointer[ocr2keepers.BlockKey]{},
		lggr:             logger.Named(lggr, "BlockSubscriber"),
	}
}

func (bs *BlockSubscriber) getBlockRange(ctx context.Context) ([]uint64, error) {
	h, err := bs.lp.LatestBlock(ctx)
	if err != nil {
		return nil, err
	}
	latestBlockNumber := h.BlockNumber
	bs.lggr.Infof("latest block from log poller is %d", latestBlockNumber)

	var blocks []uint64
	for i := bs.blockSize - 1; i >= 0; i-- {
		if latestBlockNumber-i > 0 {
			blocks = append(blocks, uint64(latestBlockNumber-i))
		}
	}
	return blocks, nil
}

func (bs *BlockSubscriber) initializeBlocks(ctx context.Context, blocks []uint64) error {
	logpollerBlocks, err := bs.lp.GetBlocksRange(ctx, blocks)
	if err != nil {
		return err
	}
	for i, b := range logpollerBlocks {
		if i == 0 {
			bs.lastClearedBlock = b.BlockNumber - 1
			bs.lggr.Infof("lastClearedBlock is %d", bs.lastClearedBlock)
		}
		bs.blocks[b.BlockNumber] = b.BlockHash.Hex()
	}
	bs.lggr.Infof("initialize with %d blocks", len(logpollerBlocks))
	return nil
}

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

	bs.lggr.Debugf("start clearing blocks from %d to %d", bs.lastClearedBlock+1, bs.lastSentBlock-bs.blockSize)
	for i := bs.lastClearedBlock + 1; i <= bs.lastSentBlock-bs.blockSize; i++ {
		delete(bs.blocks, i)
	}
	bs.lastClearedBlock = bs.lastSentBlock - bs.blockSize
	bs.lggr.Infof("lastClearedBlock is set to %d", bs.lastClearedBlock)
}

func (bs *BlockSubscriber) initialize(ctx context.Context) {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	// initialize the blocks map with the recent blockSize blocks
	blocks, err := bs.getBlockRange(ctx)
	if err != nil {
		bs.lggr.Errorf("failed to get block range; error %v", err)
	}
	err = bs.initializeBlocks(ctx, blocks)
	if err != nil {
		bs.lggr.Errorf("failed to get log poller blocks; error %v", err)
	}
	_, bs.unsubscribe = bs.hb.Subscribe(&headWrapper{headC: bs.headC, lggr: bs.lggr})
}

func (bs *BlockSubscriber) Start(ctx context.Context) error {
	return bs.StartOnce(BlockSubscriberServiceName, func() error {
		bs.lggr.Info("block subscriber started.")
		bs.initialize(ctx)
		// poll from head broadcaster channel and push to subscribers
		bs.threadCtrl.Go(func(ctx context.Context) {
			for {
				select {
				case h := <-bs.headC:
					if h != nil {
						bs.processHead(h)
					}
				case <-ctx.Done():
					return
				}
			}
		})
		// cleanup old blocks
		bs.threadCtrl.Go(func(ctx context.Context) {
			ticker := time.NewTicker(cleanUpInterval)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					bs.cleanup()
				case <-ctx.Done():
					return
				}
			}
		})

		return nil
	})
}

func (bs *BlockSubscriber) Close() error {
	return bs.StopOnce(BlockSubscriberServiceName, func() error {
		bs.lggr.Info("stop block subscriber")
		bs.threadCtrl.Close()
		bs.unsubscribe()
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

func (bs *BlockSubscriber) processHead(h *types.Head) {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	// head parent is a linked list with EVM finality depth
	// when re-org happens, new heads will have pointers to the new blocks
	i := int64(0)
	for cp := h; cp != nil; cp = cp.Parent.Load() {
		// we don't stop when a matching (block number/hash) entry is seen in the map because parent linked list may be
		// cut short during a re-org if head broadcaster backfill is not complete. This can cause some re-orged blocks
		// left in the map. for example, re-org happens for block 98, 99, 100. next head 101 from broadcaster has parent list
		// of 100, so block 100 and 101 are updated. when next head 102 arrives, it has full parent history of finality depth.
		// if we stop when we see a block number/hash match, we won't look back and correct block 98 and 99.
		// hence, we make a compromise here and check previous max(finality depth, blockSize) blocks and update the map.
		existingHash, ok := bs.blocks[cp.Number]
		if !ok {
			bs.lggr.Debugf("filling block %d with new hash %s", cp.Number, cp.Hash.Hex())
		} else if existingHash != cp.Hash.Hex() {
			bs.lggr.Warnf("overriding block %d old hash %s with new hash %s due to re-org", cp.Number, existingHash, cp.Hash.Hex())
		}
		bs.blocks[cp.Number] = cp.Hash.Hex()
		i++
		if i > int64(bs.finalityDepth) || i > bs.blockSize {
			break
		}
	}
	bs.lggr.Debugf("blocks block %d hash is %s", h.Number, h.Hash.Hex())

	history := bs.buildHistory(h.Number)
	block := &ocr2keepers.BlockKey{
		Number: ocr2keepers.BlockNumber(h.Number),
	}
	copy(block.Hash[:], h.Hash[:])
	bs.latestBlock.Store(block)
	bs.lastSentBlock = h.Number
	// send history to all subscribers
	for _, subC := range bs.subscribers {
		// wrapped in a select to not get blocked by certain subscribers
		select {
		case subC <- history:
		default:
			bs.lggr.Warnf("subscriber channel is full, dropping block history with length %d", len(history))
		}
	}

	bs.lggr.Debugf("published block history with length %d and latestBlock %d to %d subscriber(s)", len(history), bs.latestBlock.Load(), len(bs.subscribers))
}

func (bs *BlockSubscriber) queryBlocksMap(bn int64) (string, bool) {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	v, ok := bs.blocks[bn]
	return v, ok
}

type headWrapper struct {
	headC chan *types.Head
	lggr  logger.Logger
}

func (w *headWrapper) OnNewLongestChain(_ context.Context, head *types.Head) {
	if head != nil {
		select {
		case w.headC <- head:
		default:
			w.lggr.Debugf("head channel is full, discarding head %+v", head)
		}
	}
}
