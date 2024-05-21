package logprovider

import (
	"encoding/hex"
	"math/big"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/smartcontractkit/chainlink-automation/pkg/v3/random"
	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/prommetrics"
)

var (
	// defaultFastExecLogsHigh is the default upper bound / maximum number of logs that Automation is committed to process for each upkeep,
	// based on available capacity, i.e. if there are no logs from other upkeeps.
	// Used by Log buffer to limit the number of logs we are saving in memory for each upkeep in a block
	defaultFastExecLogsHigh = 32
	// defaultNumOfLogUpkeeps is the default number of log upkeeps supported by the registry.
	defaultNumOfLogUpkeeps = 50
)

// fetchedLog holds the log and the ID of the upkeep
type fetchedLog struct {
	upkeepID *big.Int
	log      logpoller.Log
	// cachedLogID is the cached log identifier, used for sorting.
	// It is calculated lazily, and cached for performance.
	cachedLogID string
}

func (l *fetchedLog) getLogID() string {
	if len(l.cachedLogID) == 0 {
		ext := ocr2keepers.LogTriggerExtension{
			Index: uint32(l.log.LogIndex),
		}
		copy(ext.TxHash[:], l.log.TxHash[:])
		copy(ext.BlockHash[:], l.log.BlockHash[:])
		l.cachedLogID = hex.EncodeToString(ext.LogIdentifier())
	}
	return l.cachedLogID
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

func (b *fetchedBlock) Append(lggr logger.Logger, fl fetchedLog, maxBlockLogs, maxUpkeepLogs int) (fetchedLog, bool) {
	has, upkeepLogs := b.has(fl.upkeepID, fl.log)
	if has {
		// Skipping known logs
		return fetchedLog{}, false
	}
	// lggr.Debugw("Adding log", "i", i, "blockBlock", currentBlock.blockNumber, "logBlock", log.BlockNumber, "id", id)
	b.logs = append(b.logs, fl)

	// drop logs if we reached limits.
	if upkeepLogs+1 > maxUpkeepLogs {
		// in case we have logs overflow for a particular upkeep, we drop a log of that upkeep,
		// based on shared, random (per block) order of the logs in the block.
		b.Sort()
		var dropped fetchedLog
		currentLogs := make([]fetchedLog, 0, len(b.logs)-1)
		for _, l := range b.logs {
			if dropped.upkeepID == nil && l.upkeepID.Cmp(fl.upkeepID) == 0 {
				dropped = l
				continue
			}
			currentLogs = append(currentLogs, l)
		}
		b.logs = currentLogs
		return dropped, true
	} else if len(b.logs)+len(b.visited) > maxBlockLogs {
		// in case we have logs overflow in the buffer level, we drop a log based on
		// shared, random (per block) order of the logs in the block.
		b.Sort()
		dropped := b.logs[0]
		b.logs = b.logs[1:]
		return dropped, true
	}

	return fetchedLog{}, true
}

// Has returns true if the block has the log,
// and the number of logs for that upkeep in the block.
func (b fetchedBlock) has(id *big.Int, log logpoller.Log) (bool, int) {
	allLogs := append(b.logs, b.visited...)
	upkeepLogs := 0
	for _, l := range allLogs {
		if l.upkeepID.Cmp(id) != 0 {
			continue
		}
		upkeepLogs++
		if l.log.BlockHash == log.BlockHash && l.log.TxHash == log.TxHash && l.log.LogIndex == log.LogIndex {
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

// Sort by log identifiers, shuffled using a pseduorandom souce that is shared across all nodes
// for a given block.
func (b *fetchedBlock) Sort() {
	randSeed := random.GetRandomKeySource(nil, uint64(b.blockNumber))

	shuffledLogIDs := make(map[string]string, len(b.logs))
	for _, log := range b.logs {
		logID := log.getLogID()
		shuffledLogIDs[logID] = random.ShuffleString(logID, randSeed)
	}

	sort.SliceStable(b.logs, func(i, j int) bool {
		return shuffledLogIDs[b.logs[i].getLogID()] < shuffledLogIDs[b.logs[j].getLogID()]
	})
}

// logEventBuffer is a circular/ring buffer of fetched logs.
// Each entry in the buffer represents a block,
// and holds the logs fetched for that block.
type logEventBuffer struct {
	lggr logger.Logger
	lock sync.RWMutex
	// size is the number of blocks supported by the buffer
	size int32

	numOfLogUpkeeps, fastExecLogsHigh uint32
	// blocks is the circular buffer of fetched blocks
	blocks []fetchedBlock
	// latestBlock is the latest block number seen
	latestBlock int64
}

func newLogEventBuffer(lggr logger.Logger, size, numOfLogUpkeeps, fastExecLogsHigh int) *logEventBuffer {
	return &logEventBuffer{
		lggr:             lggr.Named("KeepersRegistry.LogEventBuffer"),
		size:             int32(size),
		blocks:           make([]fetchedBlock, size),
		numOfLogUpkeeps:  uint32(numOfLogUpkeeps),
		fastExecLogsHigh: uint32(fastExecLogsHigh),
	}
}

func (b *logEventBuffer) latestBlockSeen() int64 {
	return atomic.LoadInt64(&b.latestBlock)
}

func (b *logEventBuffer) bufferSize() int {
	return int(atomic.LoadInt32(&b.size))
}

func (b *logEventBuffer) SetLimits(numOfLogUpkeeps, fastExecLogsHigh int) {
	atomic.StoreUint32(&b.numOfLogUpkeeps, uint32(numOfLogUpkeeps))
	atomic.StoreUint32(&b.fastExecLogsHigh, uint32(fastExecLogsHigh))
}

// enqueue adds logs (if not exist) to the buffer, returning the number of logs added
// minus the number of logs dropped.
func (b *logEventBuffer) enqueue(id *big.Int, logs ...logpoller.Log) int {
	b.lock.Lock()
	defer b.lock.Unlock()

	lggr := b.lggr.With("id", id.String())

	maxBlockLogs := int(atomic.LoadUint32(&b.fastExecLogsHigh) * atomic.LoadUint32(&b.numOfLogUpkeeps))
	maxUpkeepLogs := int(atomic.LoadUint32(&b.fastExecLogsHigh))

	latestBlock := b.latestBlockSeen()
	added, dropped := 0, 0

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
		droppedLog, ok := currentBlock.Append(lggr, fetchedLog{upkeepID: id, log: log}, maxBlockLogs, maxUpkeepLogs)
		if !ok {
			// Skipping known logs
			continue
		}
		if droppedLog.upkeepID != nil {
			dropped++
			lggr.Debugw("Reached log buffer limits, dropping log", "blockNumber", droppedLog.log.BlockNumber,
				"blockHash", droppedLog.log.BlockHash, "txHash", droppedLog.log.TxHash, "logIndex", droppedLog.log.LogIndex,
				"upkeepID", droppedLog.upkeepID.String())
		}
		added++
		b.blocks[i] = currentBlock

		if log.BlockNumber > latestBlock {
			latestBlock = log.BlockNumber
		}
	}

	if latestBlock > b.latestBlockSeen() {
		atomic.StoreInt64(&b.latestBlock, latestBlock)
	}
	if added > 0 {
		lggr.Debugw("Added logs to buffer", "addedLogs", added, "dropped", dropped, "latestBlock", latestBlock)
		prommetrics.AutomationLogBufferFlow.WithLabelValues(prommetrics.LogBufferFlowDirectionIngress).Add(float64(added))
		prommetrics.AutomationLogBufferFlow.WithLabelValues(prommetrics.LogBufferFlowDirectionDropped).Add(float64(dropped))
	}

	return added - dropped
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
func (b *logEventBuffer) dequeueRange(start, end int64, upkeepLimit, totalLimit int) []fetchedLog {
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
	totalCount := 0
	var results []fetchedLog
	for _, block := range fetchedBlocks {
		if block.blockNumber < start || block.blockNumber > end {
			// double checking that we don't have any gaps in the range
			continue
		}
		if totalCount >= totalLimit {
			// reached total limit, no need to process more blocks
			break
		}
		// Sort the logs in random order that is shared across all nodes.
		// This ensures that nodes across the network will process the same logs.
		block.Sort()
		var remainingLogs, blockResults []fetchedLog
		for _, log := range block.logs {
			if totalCount >= totalLimit {
				remainingLogs = append(remainingLogs, log)
				continue
			}
			if logsCount[log.upkeepID.String()] >= upkeepLimit {
				remainingLogs = append(remainingLogs, log)
				continue
			}
			blockResults = append(blockResults, log)
			logsCount[log.upkeepID.String()]++
			totalCount++
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
		prommetrics.AutomationLogBufferFlow.WithLabelValues(prommetrics.LogBufferFlowDirectionEgress).Add(float64(len(results)))
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
