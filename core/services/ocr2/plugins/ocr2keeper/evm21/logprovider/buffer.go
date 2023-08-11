package logprovider

import (
	"math/big"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// fetchedLog holds the log and the ID of the upkeep
type fetchedLog struct {
	id  *big.Int
	log logpoller.Log
}

// fetchedBlock holds the logs fetched for a block
type fetchedBlock struct {
	blockNumber int64
	// logs is the logs fetched for the block and haven't been visited yet
	logs []fetchedLog
	// visited is the logs fetched for the block and have been visited.
	// We keep them in order to avoid fetching them again.
	visited []fetchedLog
}

// Reset resets the block to the given block number, if the block is newer than the current block.
func (b fetchedBlock) Reset(block int64) (fetchedBlock, int64) {
	if b.blockNumber < block {
		return fetchedBlock{
			blockNumber: block,
		}, b.blockNumber
	}
	return b, b.blockNumber
}

// Has returns true if the block has the log,
// and the number of logs for that upkeep in the block.
func (b fetchedBlock) Has(id *big.Int, log logpoller.Log) (bool, int) {
	allLogs := append(b.logs, b.visited...)
	upkeepLogs := 0
	for _, l := range allLogs {
		if l.id.Cmp(id) != 0 {
			continue
		}
		upkeepLogs++
		if l.log.BlockNumber == log.BlockNumber && l.log.TxHash == log.TxHash && l.log.LogIndex == log.LogIndex {
			return true, upkeepLogs
		}
	}
	return false, upkeepLogs
}

// logEventBuffer is a circular/ring buffer of fetched logs.
// Each entry in the buffer represents a block,
// and holds the logs fetched for that block.
type logEventBuffer struct {
	lggr logger.Logger
	lock sync.RWMutex
	// size is the number of blocks supported by the buffer
	size int32

	maxBlockLogs, maxUpkeepLogsPerBlock int32
	// blocks is the circular buffer of fetched blocks
	blocks []fetchedBlock
	// latestBlock is the latest block number seen
	latestBlock int64
}

func newLogEventBuffer(lggr logger.Logger, size, maxBlockLogs, maxUpkeepLogsPerBlock int) *logEventBuffer {
	return &logEventBuffer{
		lggr:                  lggr.Named("KeepersRegistry.LogEventBuffer"),
		size:                  int32(size),
		blocks:                make([]fetchedBlock, size),
		maxBlockLogs:          int32(maxBlockLogs),
		maxUpkeepLogsPerBlock: int32(maxUpkeepLogsPerBlock),
	}
}

func (b *logEventBuffer) latestBlockSeen() int64 {
	return atomic.LoadInt64(&b.latestBlock)
}

func (b *logEventBuffer) bufferSize() int {
	return int(atomic.LoadInt32(&b.size))
}

// enqueue adds logs (if not exist) to the buffer, returning the number of logs added
func (b *logEventBuffer) enqueue(id *big.Int, logs ...logpoller.Log) int {
	b.lock.Lock()
	defer b.lock.Unlock()

	lggr := b.lggr.With("id", id.String())

	maxBlockLogs := int(b.maxBlockLogs)
	maxUpkeepLogs := int(b.maxUpkeepLogsPerBlock)

	latestBlock := b.latestBlockSeen()
	added := 0
	for _, log := range logs {
		if log.BlockNumber == 0 {
			// invalid log
			continue
		}
		i := b.blockNumberIndex(log.BlockNumber)
		block, prevBlock := b.blocks[i].Reset(log.BlockNumber)
		if prevBlock > 0 {
			if prevBlock > log.BlockNumber {
				lggr.Debugw("Skipping old log", "currentBlock", block.blockNumber, "newBlock", log.BlockNumber)
				continue
			} else if prevBlock < log.BlockNumber {
				lggr.Debugw("Overriding block", "prevBlock", prevBlock, "newBlock", log.BlockNumber)
			}
		}
		if len(block.logs)+1 > maxBlockLogs {
			lggr.Debugw("Reached max logs number per block, dropping log", "blockNumber", log.BlockNumber,
				"blockHash", log.BlockHash, "txHash", log.TxHash, "logIndex", log.LogIndex)
			continue
		}
		if has, upkeepLogs := block.Has(id, log); has {
			// lggr.Debugw("Skipping existing log", "blockNumber", log.BlockNumber,
			// 	"blockHash", log.BlockHash, "txHash", log.TxHash, "logIndex", log.LogIndex)
			continue
		} else if upkeepLogs+1 > maxUpkeepLogs {
			lggr.Debugw("Reached max logs number per upkeep, dropping log", "blockNumber", log.BlockNumber,
				"blockHash", log.BlockHash, "txHash", log.TxHash, "logIndex", log.LogIndex)
			continue
		}
		// lggr.Debugw("Adding log", "i", i, "blockBlock", block.blockNumber, "logBlock", log.BlockNumber, "id", id)
		block.logs = append(block.logs, fetchedLog{id: id, log: log})
		b.blocks[i] = block
		added++
		if log.BlockNumber > latestBlock {
			latestBlock = log.BlockNumber
		}
	}

	if latestBlock > b.latestBlockSeen() {
		atomic.StoreInt64(&b.latestBlock, latestBlock)
	}
	if added > 0 {
		lggr.Debugw("Added logs to buffer", "addedLogs", added, "latestBlock", latestBlock)
	}

	return added
}

// peek returns the logs in range [latestBlock-blocks, latestBlock]
func (b *logEventBuffer) peek(blocks int) []fetchedLog {
	latestBlock := b.latestBlockSeen()
	if latestBlock == 0 {
		return nil
	}
	if blocks > int(latestBlock) {
		blocks = int(latestBlock) - 1
	}

	return b.peekRange(latestBlock-int64(blocks), latestBlock)
}

// peekRange returns the logs between start and end inclusive.
func (b *logEventBuffer) peekRange(start, end int64) []fetchedLog {
	b.lock.RLock()
	defer b.lock.RUnlock()

	blocksInRange := b.getBlocksInRange(int(start), int(end))

	var results []fetchedLog
	for _, block := range blocksInRange {
		// double checking that we don't have any gaps in the range
		if block.blockNumber < start || block.blockNumber > end {
			continue
		}
		results = append(results, block.logs...)
	}

	sort.SliceStable(results, func(i, j int) bool {
		return results[i].log.BlockNumber < results[j].log.BlockNumber
	})

	b.lggr.Debugw("Peeked logs", "results", len(results), "start", start, "end", end)

	return results
}

// peek returns the logs in range [latestBlock-blocks, latestBlock]
func (b *logEventBuffer) dequeue(blocks int) []fetchedLog {
	latestBlock := b.latestBlockSeen()
	if latestBlock == 0 {
		return nil
	}
	if blocks > int(latestBlock) {
		blocks = int(latestBlock) - 1
	}
	return b.dequeueRange(latestBlock-int64(blocks), latestBlock)
}

// dequeueRange returns the logs between start and end inclusive.
func (b *logEventBuffer) dequeueRange(start, end int64) []fetchedLog {
	b.lock.Lock()
	defer b.lock.Unlock()

	blocksInRange := b.getBlocksInRange(int(start), int(end))

	var results []fetchedLog
	for i, block := range blocksInRange {
		// double checking that we don't have any gaps in the range
		if block.blockNumber < start || block.blockNumber > end {
			continue
		}
		results = append(results, block.logs...)
		block.visited = append(block.visited, block.logs...)
		block.logs = nil
		b.blocks[i] = block
	}

	sort.SliceStable(results, func(i, j int) bool {
		return results[i].log.BlockNumber < results[j].log.BlockNumber
	})

	b.lggr.Debugw("Dequeued logs", "results", len(results), "start", start, "end", end)

	return results
}

// getBlocksInRange returns the blocks between start and end.
// NOTE: this function should be called with the lock held
func (b *logEventBuffer) getBlocksInRange(start, end int) []fetchedBlock {
	var blocksInRange []fetchedBlock
	start, end = b.normalRange(start, end)
	if start == -1 || end == -1 {
		// invalid range
		return blocksInRange
	}
	if start < end {
		return b.blocks[start:end]
	}
	// in case we get circular range such as [0, 1, end, ... , start, ..., size-1]
	// we need to return the blocks in two ranges: [start, size-1] and [0, end]
	blocksInRange = append(blocksInRange, b.blocks[start:]...)
	blocksInRange = append(blocksInRange, b.blocks[:end]...)

	return blocksInRange
}

// normalRange returns the normalized range of start and end,
// aligned with buffer sizes.
func (b *logEventBuffer) normalRange(start, end int) (int, int) {
	if end < start || end == 0 {
		// invalid range
		return -1, -1
	}
	size := b.bufferSize()
	if start == 0 {
		// we reduce start by 1 to make it easier to calculate the index,
		// but we need to ensure we don't go below 0.
		start++
	}
	if start == end {
		// ensure we have at least one block in range
		end++
	}
	if end-start > size {
		// ensure we don't have more than the buffer size
		start = (end - size) + 1
	}
	start = (start - 1) % size
	end = end % size

	return start, end
}

// blockNumberIndex returns the index of the block in the buffer
func (b *logEventBuffer) blockNumberIndex(bn int64) int {
	return int(bn-1) % b.bufferSize()
}
