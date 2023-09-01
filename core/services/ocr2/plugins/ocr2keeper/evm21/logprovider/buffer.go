package logprovider

import (
	"math/big"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var (
	// allowedLogsPerBlock is the maximum number of logs allowed per upkeep in a block.
	allowedLogsPerBlock = 128
	// bufferMaxBlockSize is the maximum number of blocks in the buffer.
	bufferMaxBlockSize = 1024
)

// fetchedLog holds the log and the ID of the upkeep
type fetchedLog struct {
	upkeepID *big.Int
	log      logpoller.Log
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

// Has returns true if the block has the log,
// and the number of logs for that upkeep in the block.
func (b fetchedBlock) Has(id *big.Int, log logpoller.Log) (bool, int) {
	allLogs := append(b.logs, b.visited...)
	upkeepLogs := 0
	for _, l := range allLogs {
		if l.upkeepID.Cmp(id) != 0 {
			continue
		}
		upkeepLogs++
		if l.log.BlockNumber == log.BlockNumber && l.log.TxHash == log.TxHash && l.log.LogIndex == log.LogIndex {
			return true, upkeepLogs
		}
	}
	return false, upkeepLogs
}

func (b fetchedBlock) Clone() fetchedBlock {
	logs := make([]fetchedLog, len(b.logs))
	copy(logs, b.logs)
	visited := make([]fetchedLog, len(b.visited))
	copy(visited, b.visited)
	return fetchedBlock{
		blockNumber: b.blockNumber,
		logs:        logs,
		visited:     visited,
	}
}

// logEventBuffer is a circular/ring buffer of fetched logs.
// Each entry in the buffer represents a block,
// and holds the logs fetched for that block.
type logEventBuffer struct {
	lggr logger.Logger
	lock sync.RWMutex
	// size is the number of blocks supported by the buffer
	size int32

	maxBlockLogs, maxUpkeepLogsPerBlock int
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
		maxBlockLogs:          maxBlockLogs,
		maxUpkeepLogsPerBlock: maxUpkeepLogsPerBlock,
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
		currentBlock := b.blocks[i]
		if currentBlock.blockNumber < log.BlockNumber {
			lggr.Debugw("Got log on a new block", "prevBlock", currentBlock.blockNumber, "newBlock", log.BlockNumber)
			currentBlock.blockNumber = log.BlockNumber
			currentBlock.logs = nil
			currentBlock.visited = nil
		} else if currentBlock.blockNumber > log.BlockNumber {
			// not expected to happen
			lggr.Debugw("Skipping log from old block", "currentBlock", currentBlock.blockNumber, "newBlock", log.BlockNumber)
			continue
		}
		if len(currentBlock.logs)+1 > maxBlockLogs {
			lggr.Debugw("Reached max logs number per block, dropping log", "blockNumber", log.BlockNumber,
				"blockHash", log.BlockHash, "txHash", log.TxHash, "logIndex", log.LogIndex)
			continue
		}
		if has, upkeepLogs := currentBlock.Has(id, log); has {
			// Skipping existing log
			continue
		} else if upkeepLogs+1 > maxUpkeepLogs {
			lggr.Debugw("Reached max logs number per upkeep, dropping log", "blockNumber", log.BlockNumber,
				"blockHash", log.BlockHash, "txHash", log.TxHash, "logIndex", log.LogIndex)
			continue
		}
		// lggr.Debugw("Adding log", "i", i, "blockBlock", currentBlock.blockNumber, "logBlock", log.BlockNumber, "id", id)
		currentBlock.logs = append(currentBlock.logs, fetchedLog{upkeepID: id, log: log})
		b.blocks[i] = currentBlock
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

// dequeueRange returns the logs between start and end inclusive.
func (b *logEventBuffer) dequeueRange(start, end int64, upkeepLimit int) []fetchedLog {
	b.lock.Lock()
	defer b.lock.Unlock()

	blocksInRange := b.getBlocksInRange(int(start), int(end))
	fetchedBlocks := make([]fetchedBlock, 0, len(blocksInRange))
	for _, block := range blocksInRange {
		// Create clone of the blocks as they get processed and update underlying b.blocks
		fetchedBlocks = append(fetchedBlocks, block.Clone())
	}

	// Sort the blocks in reverse order of block number so that latest logs
	// are preferred while dequeueing.
	sort.SliceStable(fetchedBlocks, func(i, j int) bool {
		return fetchedBlocks[i].blockNumber > fetchedBlocks[j].blockNumber
	})

	logsCount := map[string]int{}
	var results []fetchedLog
	for _, block := range fetchedBlocks {
		// double checking that we don't have any gaps in the range
		if block.blockNumber < start || block.blockNumber > end {
			continue
		}
		var remainingLogs, blockResults []fetchedLog
		for _, log := range block.logs {
			if logsCount[log.upkeepID.String()] >= upkeepLimit {
				remainingLogs = append(remainingLogs, log)
				continue
			}
			logsCount[log.upkeepID.String()]++
			blockResults = append(blockResults, log)
		}
		if len(blockResults) == 0 {
			continue
		}
		results = append(results, blockResults...)
		block.visited = append(block.visited, blockResults...)
		block.logs = remainingLogs
		b.blocks[b.blockNumberIndex(block.blockNumber)] = block
	}

	if len(results) > 0 {
		b.lggr.Debugw("Dequeued logs", "results", len(results), "start", start, "end", end)
	}

	return results
}

// getBlocksInRange returns the blocks between start and end.
// NOTE: this function should be called with the lock held
func (b *logEventBuffer) getBlocksInRange(start, end int) []fetchedBlock {
	var blocksInRange []fetchedBlock
	start, end = b.blockRangeToIndices(start, end)
	if start == -1 || end == -1 {
		// invalid range
		return blocksInRange
	}
	if start <= end {
		// Normal range, need to return indices from start to end(inclusive)
		return b.blocks[start : end+1]
	}
	// in case we get circular range such as [0, 1, end, ... , start, ..., size-1]
	// we need to return the blocks in two ranges: [0, end](inclusive) and [start, size-1]
	blocksInRange = append(blocksInRange, b.blocks[start:]...)
	blocksInRange = append(blocksInRange, b.blocks[:end+1]...)

	return blocksInRange
}

// blockRangeToIndices returns the normalized range of start to end block range,
// to indices aligned with buffer size. Note ranges inclusive of start, end indices.
func (b *logEventBuffer) blockRangeToIndices(start, end int) (int, int) {
	latest := b.latestBlockSeen()
	if end > int(latest) {
		// Limit end of range to latest block seen
		end = int(latest)
	}
	if end < start || start == 0 || end == 0 {
		// invalid range
		return -1, -1
	}
	size := b.bufferSize()
	if end-start >= size {
		// If range requires more than buffer size blocks, only to return
		// last size blocks as that's the max the buffer stores.
		start = (end - size) + 1
	}
	return b.blockNumberIndex(int64(start)), b.blockNumberIndex(int64(end))
}

// blockNumberIndex returns the index of the block in the buffer
func (b *logEventBuffer) blockNumberIndex(bn int64) int {
	return int(bn-1) % b.bufferSize()
}
