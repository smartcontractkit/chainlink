package logprovider

import (
	"encoding/hex"
	"math/big"
	"sort"
	"sync"
	"sync/atomic"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/prommetrics"
)

const (
	defaultLogLimitHigh = 10
)

type BufferedLog struct {
	ID  *big.Int
	Log logpoller.Log
}

type LogBuffer interface {
	// Enqueue adds logs to the buffer and might also drop logs if the limit for the
	// given upkeep was exceeded. Returns the number of logs that were added and number of logs that were  dropped.
	Enqueue(id *big.Int, logs ...logpoller.Log) (added int, dropped int)
	// Dequeue pulls logs from the buffer that are within the given block window,
	// with a maximum number of logs per upkeep and a total maximum number of logs to return.
	// It also accepts a function to select upkeeps.
	// Returns logs (associated to upkeeps) and the number of remaining
	// logs in that window for the involved upkeeps.
	Dequeue(block int64, blockRate, upkeepLimit, maxResults int, upkeepSelector func(id *big.Int) bool) ([]BufferedLog, int)
	// SetConfig sets the buffer size and the maximum number of logs to keep for each upkeep.
	SetConfig(lookback, maxUpkeepLogs int)
}

func DefaultUpkeepSelector(id *big.Int) bool {
	return true
}

type logBuffer struct {
	lggr logger.Logger
	// max number of logs to keep in the buffer for each upkeep per block
	maxUpkeepLogs *atomic.Int32
	// number of blocks to keep in the buffer
	bufferSize *atomic.Int32
	// last block number seen by the buffer
	lastBlockSeen *atomic.Int64
	// map of upkeep id to its buffer
	upkeepBuffers map[string]*upkeepLogBuffer
	lock          sync.RWMutex
}

func NewLogBuffer(lggr logger.Logger, size, upkeepLogLimit int) LogBuffer {
	s := new(atomic.Int32)
	s.Add(int32(size))
	l := new(atomic.Int32)
	l.Add(int32(upkeepLogLimit))
	return &logBuffer{
		lggr:          lggr.Named("KeepersRegistry.LogEventBufferV1"),
		maxUpkeepLogs: l,
		bufferSize:    s,
		lastBlockSeen: new(atomic.Int64),
		upkeepBuffers: make(map[string]*upkeepLogBuffer),
	}
}

func (b *logBuffer) SetConfig(lookback, logLimitHigh int) {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.bufferSize.Store(int32(lookback))
	b.maxUpkeepLogs.Store(int32(logLimitHigh))

	for _, ub := range b.upkeepBuffers {
		ub.setConfig(logLimitHigh)
	}
}

// Enqueue adds logs to the buffer and might also drop logs if the limit for the
// given upkeep was exceeded. It will create a new buffer if it does not exist.
// Returns the number of logs that were added and number of logs that were  dropped.
func (b *logBuffer) Enqueue(uid *big.Int, logs ...logpoller.Log) (int, int) {
	buf, ok := b.getUpkeepBuffer(uid)
	if !ok || buf == nil {
		buf = newUpkeepLogBuffer(b.lggr, uid, int(b.maxUpkeepLogs.Load()*b.bufferSize.Load()))
		b.setUpkeepBuffer(uid, buf)
	}
	lastBlockSeen := latestBlockNumber(logs...)
	if b.lastBlockSeen.Load() < lastBlockSeen {
		b.lastBlockSeen.Store(lastBlockSeen)
	}
	blockThreshold := b.lastBlockSeen.Load() - int64(b.bufferSize.Load())
	if blockThreshold <= 0 {
		blockThreshold = 1
	}
	return buf.enqueue(blockThreshold, logs...)
}

// Dequeue greedly pulls logs from the buffers.
// Returns logs and the number of remaining logs in the buffer.
func (b *logBuffer) Dequeue(block int64, blockRate, upkeepLimit, maxResults int, upkeepSelector func(id *big.Int) bool) ([]BufferedLog, int) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	start, end := BlockWindow(block, blockRate)
	result, remaining := b.tryDequeue(start, end, upkeepLimit, maxResults, upkeepSelector)
	// if there are still logs to pull, try to dequeue again
	// TODO: check if we should limit the number of iterations
	for len(result) < maxResults && remaining > 0 {
		nextResults, nextRemaining := b.tryDequeue(start, end, upkeepLimit, maxResults-len(result), upkeepSelector)
		result = append(result, nextResults...)
		remaining = nextRemaining
	}

	return result, remaining
}

// tryDequeue pulls logs from the buffers, according to the given selector, in block range [start,end]
// with minimum number of results per upkeep and the total capacity for results.
// Returns logs and the number of remaining logs in the buffer.
func (b *logBuffer) tryDequeue(start, end int64, minUpkeepLogs, capacity int, upkeepSelector func(id *big.Int) bool) ([]BufferedLog, int) {
	var result []BufferedLog
	var remainingLogs int
	for _, buf := range b.upkeepBuffers {
		if !upkeepSelector(buf.id) {
			// if the upkeep is not selected, skip it
			continue
		}
		if capacity == 0 {
			// if there is no more capacity for results, just count the remaining logs
			remainingLogs += buf.sizeOfWindow(start, end)
			continue
		}
		if minUpkeepLogs > capacity {
			// if there are more logs to fetch than the capacity, fetch the minimum
			minUpkeepLogs = capacity
		}
		logs, remaining := buf.dequeue(start, end, minUpkeepLogs)
		for _, l := range logs {
			result = append(result, BufferedLog{ID: buf.id, Log: l})
			capacity--
		}
		remainingLogs += remaining
	}
	return result, remainingLogs
}

func (b *logBuffer) getUpkeepBuffer(uid *big.Int) (*upkeepLogBuffer, bool) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	ub, ok := b.upkeepBuffers[uid.String()]
	return ub, ok
}

func (b *logBuffer) setUpkeepBuffer(uid *big.Int, buf *upkeepLogBuffer) {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.upkeepBuffers[uid.String()] = buf
}

type upkeepLogBuffer struct {
	lggr logger.Logger

	id      *big.Int
	maxLogs *atomic.Int32

	q       []logpoller.Log
	visited map[string]int64
	lock    sync.RWMutex
}

func newUpkeepLogBuffer(lggr logger.Logger, id *big.Int, maxLogs int) *upkeepLogBuffer {
	limit := new(atomic.Int32)
	limit.Add(int32(maxLogs))
	return &upkeepLogBuffer{
		lggr:    lggr.With("id", id.String()),
		id:      id,
		maxLogs: limit,
		q:       make([]logpoller.Log, 0, maxLogs),
		visited: make(map[string]int64),
	}
}

func (ub *upkeepLogBuffer) setConfig(maxLogs int) {
	ub.maxLogs.Store(int32(maxLogs))
}

// size returns the total number of logs in the buffer.
func (ub *upkeepLogBuffer) size() int {
	ub.lock.RLock()
	defer ub.lock.RUnlock()

	return len(ub.q)
}

// size returns the total number of logs in the buffer.
func (ub *upkeepLogBuffer) sizeOfWindow(start, end int64) int {
	ub.lock.RLock()
	defer ub.lock.RUnlock()

	size := 0
	for _, l := range ub.q {
		if l.BlockNumber >= start && l.BlockNumber <= end {
			size++
		}
	}
	return size
}

// dequeue pulls logs from the buffer that are within the given block range,
// with a limit of logs to pull. Returns logs and the number of remaining logs in the buffer.
func (ub *upkeepLogBuffer) dequeue(start, end int64, limit int) ([]logpoller.Log, int) {
	ub.lock.Lock()
	defer ub.lock.Unlock()

	if len(ub.q) == 0 {
		return nil, 0
	}

	var results []logpoller.Log
	var remaining int
	updatedLogs := make([]logpoller.Log, 0)
	for _, l := range ub.q {
		if l.BlockNumber >= start && l.BlockNumber <= end {
			if len(results) < limit {
				results = append(results, l)
				continue
			}
			remaining++
		}
		updatedLogs = append(updatedLogs, l)
	}

	if len(results) > 0 {
		ub.q = updatedLogs
	}

	ub.lggr.Debugf("Dequeued %d logs, remaining %d", len(results), remaining)
	prommetrics.AutomationLogsInLogBuffer.Sub(float64(len(results)))

	return results, remaining
}

// enqueue adds logs to the buffer and might also drop logs if the limit for the
// given upkeep was exceeded. Additionally, it will drop logs that are older than blockThreshold.
// Returns the number of logs that were added and number of logs that were  dropped.
func (ub *upkeepLogBuffer) enqueue(blockThreshold int64, logsToAdd ...logpoller.Log) (int, int) {
	ub.lock.Lock()
	defer ub.lock.Unlock()

	logs := ub.q
	var added int
	for _, log := range logsToAdd {
		if log.BlockNumber < blockThreshold {
			ub.lggr.Debugw("Skipping log from old block", "blockThreshold", blockThreshold, "logBlock", log.BlockNumber)
			continue
		}
		logid := logID(log)
		if _, ok := ub.visited[logid]; ok {
			ub.lggr.Debugw("Skipping known log", "blockThreshold", blockThreshold, "logBlock", log.BlockNumber)
			continue
		}
		added++
		if len(logs) == 0 {
			// if the buffer is empty, just add the log
			logs = append(logs, log)
		} else {
			// otherwise, find the right index to insert the log
			// to keep the buffer sorted
			// TODO: check what is better: 1. maintain sorted slice; 2. sort once at the end
			i, _ := sort.Find(len(logs), func(i int) int {
				return LogComparator(log, logs[i])
			})
			if i == len(logs) {
				logs = append(logs, log)
			} else {
				logs = append(logs[:i], append([]logpoller.Log{log}, logs[i:]...)...)
			}
		}
		ub.visited[logid] = log.BlockNumber
	}
	ub.q = logs

	var dropped int
	if added > 0 {
		dropped = ub.clean(blockThreshold)
	}

	ub.lggr.Debugf("Enqueued %d logs, dropped %d with blockThreshold %d", added, dropped, blockThreshold)
	prommetrics.AutomationLogsInLogBuffer.Add(float64(added))

	return added, dropped
}

// clean removes logs that are older than blockThreshold and drops logs if the limit for the
// given upkeep was exceeded. Returns the number of logs that were dropped.
func (ub *upkeepLogBuffer) clean(blockThreshold int64) int {
	maxLogs := int(ub.maxLogs.Load())

	// sort.SliceStable(updated, func(i, j int) bool {
	// 	return LogSorter(updated[i], updated[j])
	// })
	updated := make([]logpoller.Log, 0)
	var dropped int
	for _, l := range ub.q {
		if l.BlockNumber > blockThreshold {
			if len(updated) < maxLogs {
				updated = append(updated, l)
			} else {
				prommetrics.AutomationLogsInLogBuffer.Dec()
				// TODO: check if we should clean visited as well
				ub.lggr.Debugw("Reached log buffer limits, dropping log", "blockNumber", l.BlockNumber,
					"blockHash", l.BlockHash, "txHash", l.TxHash, "logIndex", l.LogIndex, "len updated", len(updated), "maxLogs", maxLogs)
				dropped++
			}
		} else {
			prommetrics.AutomationLogsInLogBuffer.Dec()
			// old logs are ignored and removed from visited
			ub.lggr.Debugw("Dropping old log", "blockNumber", l.BlockNumber, "blockThreshold", blockThreshold, "logIndex", l.LogIndex)
			logid := logID(l)
			delete(ub.visited, logid)
		}
	}

	ub.lggr.Debugw("Cleaned logs", "dropped", dropped, "blockThreshold", blockThreshold, "len updated", len(updated), "len ub.q", len(ub.q), "maxLogs", maxLogs)

	ub.q = updated

	for lid, block := range ub.visited {
		if block <= blockThreshold {
			delete(ub.visited, lid)
		}
	}

	return dropped
}

// logID returns a unique identifier for a log, which is an hex string
// of ocr2keepers.LogTriggerExtension.LogIdentifier()
func logID(l logpoller.Log) string {
	ext := ocr2keepers.LogTriggerExtension{
		Index: uint32(l.LogIndex),
	}
	copy(ext.TxHash[:], l.TxHash[:])
	copy(ext.BlockHash[:], l.BlockHash[:])
	return hex.EncodeToString(ext.LogIdentifier())
}

// latestBlockNumber returns the latest block number from the given logs
func latestBlockNumber(logs ...logpoller.Log) int64 {
	var latest int64
	for _, l := range logs {
		if l.BlockNumber > latest {
			latest = l.BlockNumber
		}
	}
	return latest
}
