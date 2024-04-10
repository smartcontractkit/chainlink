package logprovider

import (
	"math"
	"math/big"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/prommetrics"
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
	SetConfig(lookback, blockRate, logLimit uint32)
	// NumOfUpkeeps returns the number of upkeeps that are being tracked by the buffer.
	NumOfUpkeeps() int
	// SyncFilters removes upkeeps that are not in the filter store.
	SyncFilters(filterStore UpkeepFilterStore) error
}

func DefaultUpkeepSelector(id *big.Int) bool {
	return true
}

type logBufferOptions struct {
	// number of blocks to keep in the buffer
	lookback *atomic.Uint32
	// blockRate is the number of blocks per window
	blockRate *atomic.Uint32
	// max number of logs to keep in the buffer for each upkeep per window (LogLimit*10)
	windowLimit *atomic.Uint32
}

func newLogBufferOptions(lookback, blockRate, logLimit uint32) *logBufferOptions {
	opts := &logBufferOptions{
		windowLimit: new(atomic.Uint32),
		lookback:    new(atomic.Uint32),
		blockRate:   new(atomic.Uint32),
	}
	opts.override(lookback, blockRate, logLimit)

	return opts
}

func (o *logBufferOptions) override(lookback, blockRate, logLimit uint32) {
	o.windowLimit.Store(logLimit * 10)
	o.lookback.Store(lookback)
	o.blockRate.Store(blockRate)
}

func (o *logBufferOptions) windows() int {
	return int(math.Ceil(float64(o.lookback.Load()) / float64(o.blockRate.Load())))
}

type logBuffer struct {
	lggr logger.Logger
	opts *logBufferOptions
	// last block number seen by the buffer
	lastBlockSeen *atomic.Int64
	// map of upkeep id to its queue
	queues map[string]*upkeepLogQueue
	lock   sync.RWMutex
}

func NewLogBuffer(lggr logger.Logger, lookback, blockRate, logLimit uint32) LogBuffer {
	return &logBuffer{
		lggr:          lggr.Named("KeepersRegistry.LogEventBufferV1"),
		opts:          newLogBufferOptions(lookback, blockRate, logLimit),
		lastBlockSeen: new(atomic.Int64),
		queues:        make(map[string]*upkeepLogQueue),
	}
}

// Enqueue adds logs to the buffer and might also drop logs if the limit for the
// given upkeep was exceeded. It will create a new buffer if it does not exist.
// Returns the number of logs that were added and number of logs that were  dropped.
func (b *logBuffer) Enqueue(uid *big.Int, logs ...logpoller.Log) (int, int) {
	buf, ok := b.getUpkeepQueue(uid)
	if !ok || buf == nil {
		buf = newUpkeepLogQueue(b.lggr, uid, b.opts)
		b.setUpkeepQueue(uid, buf)
	}
	latestBlock := latestBlockNumber(logs...)
	if b.lastBlockSeen.Load() < latestBlock {
		b.lastBlockSeen.Store(latestBlock)
	}
	blockThreshold := b.lastBlockSeen.Load() - int64(b.opts.lookback.Load())
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

	start, end := getBlockWindow(block, blockRate)
	return b.dequeue(start, end, upkeepLimit, maxResults, upkeepSelector)
}

// dequeue pulls logs from the buffers, depends the given selector (upkeepSelector),
// in block range [start,end] with minimum number of results per upkeep (upkeepLimit)
// and the maximum number of results (capacity).
// Returns logs and the number of remaining logs in the buffer for the given range and selector.
// NOTE: this method is not thread safe and should be called within a lock.
func (b *logBuffer) dequeue(start, end int64, upkeepLimit, capacity int, upkeepSelector func(id *big.Int) bool) ([]BufferedLog, int) {
	var result []BufferedLog
	var remainingLogs int
	for _, q := range b.queues {
		if !upkeepSelector(q.id) {
			// if the upkeep is not selected, skip it
			continue
		}
		logsInRange := q.sizeOfRange(start, end)
		if logsInRange == 0 {
			// if there are no logs in the range, skip the upkeep
			continue
		}
		if capacity == 0 {
			// if there is no more capacity for results, just count the remaining logs
			remainingLogs += logsInRange
			continue
		}
		if upkeepLimit > capacity {
			// adjust limit if it is higher than the actual capacity
			upkeepLimit = capacity
		}
		logs, remaining := q.dequeue(start, end, upkeepLimit)
		for _, l := range logs {
			result = append(result, BufferedLog{ID: q.id, Log: l})
			capacity--
		}
		remainingLogs += remaining
	}
	return result, remainingLogs
}

func (b *logBuffer) SetConfig(lookback, blockRate, logLimit uint32) {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.opts.override(lookback, blockRate, logLimit)
}

func (b *logBuffer) NumOfUpkeeps() int {
	b.lock.RLock()
	defer b.lock.RUnlock()

	return len(b.queues)
}

func (b *logBuffer) SyncFilters(filterStore UpkeepFilterStore) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	for upkeepID := range b.queues {
		uid := new(big.Int)
		_, ok := uid.SetString(upkeepID, 10)
		if ok && !filterStore.Has(uid) {
			// remove upkeep that is not in the filter store
			delete(b.queues, upkeepID)
		}
	}

	return nil
}

func (b *logBuffer) getUpkeepQueue(uid *big.Int) (*upkeepLogQueue, bool) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	ub, ok := b.queues[uid.String()]
	return ub, ok
}

func (b *logBuffer) setUpkeepQueue(uid *big.Int, buf *upkeepLogQueue) {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.queues[uid.String()] = buf
}

// TODO (AUTO-9256) separate files

// logTriggerState represents the state of a log in the buffer.
type logTriggerState uint8

const (
	// the log was dropped due to buffer limits
	logTriggerStateDropped logTriggerState = iota
	// the log was enqueued by the buffer
	logTriggerStateEnqueued
	// the log was visited/dequeued from the buffer
	logTriggerStateDequeued
)

// logTriggerStateEntry represents the state of a log in the buffer and the block number of the log.
// TODO (AUTO-10013) handling of reorgs might require to store the block hash as well.
type logTriggerStateEntry struct {
	state logTriggerState
	block int64
}

// upkeepLogQueue is a priority queue for logs associated to a specific upkeep.
// It keeps track of the logs that were already visited and the capacity of the queue.
type upkeepLogQueue struct {
	lggr logger.Logger

	id   *big.Int
	opts *logBufferOptions

	// logs is the buffer of logs for the upkeep
	logs []logpoller.Log
	// states keeps track of the state of the logs that are known to the queue
	// and the block number they were seen at
	states map[string]logTriggerStateEntry
	lock   sync.RWMutex
}

func newUpkeepLogQueue(lggr logger.Logger, id *big.Int, opts *logBufferOptions) *upkeepLogQueue {
	maxLogs := int(opts.windowLimit.Load()) * opts.windows() // limit per window * windows
	return &upkeepLogQueue{
		lggr:   lggr.With("upkeepID", id.String()),
		id:     id,
		opts:   opts,
		logs:   make([]logpoller.Log, 0, maxLogs),
		states: make(map[string]logTriggerStateEntry),
	}
}

// sizeOfRange returns the number of logs in the buffer that are within the given block range.
func (q *upkeepLogQueue) sizeOfRange(start, end int64) int {
	q.lock.RLock()
	defer q.lock.RUnlock()

	size := 0
	for _, l := range q.logs {
		if l.BlockNumber >= start && l.BlockNumber <= end {
			size++
		}
	}
	return size
}

// dequeue pulls logs from the buffer that are within the given block range,
// with a limit of logs to pull. Returns logs and the number of remaining logs in the buffer.
func (q *upkeepLogQueue) dequeue(start, end int64, limit int) ([]logpoller.Log, int) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if len(q.logs) == 0 {
		return nil, 0
	}

	var results []logpoller.Log
	var remaining int
	updatedLogs := make([]logpoller.Log, 0)
	for _, l := range q.logs {
		if l.BlockNumber >= start && l.BlockNumber <= end {
			if len(results) < limit {
				results = append(results, l)
				lid := logID(l)
				if s, ok := q.states[lid]; ok {
					s.state = logTriggerStateDequeued
					q.states[lid] = s
				}
				continue
			}
			remaining++
		}
		updatedLogs = append(updatedLogs, l)
	}

	if len(results) > 0 {
		q.logs = updatedLogs
		q.lggr.Debugw("Dequeued logs", "start", start, "end", end, "limit", limit, "results", len(results), "remaining", remaining)
	}

	prommetrics.AutomationLogBufferFlow.WithLabelValues(prommetrics.LogBufferFlowDirectionEgress).Add(float64(len(results)))

	return results, remaining
}

// enqueue adds logs to the buffer and might also drop logs if the limit for the
// given upkeep was exceeded. Additionally, it will drop logs that are older than blockThreshold.
// Returns the number of logs that were added and number of logs that were  dropped.
func (q *upkeepLogQueue) enqueue(blockThreshold int64, logsToAdd ...logpoller.Log) (int, int) {
	q.lock.Lock()
	defer q.lock.Unlock()

	logs := q.logs
	var added int
	for _, log := range logsToAdd {
		if log.BlockNumber < blockThreshold {
			// q.lggr.Debugw("Skipping log from old block", "blockThreshold", blockThreshold, "logBlock", log.BlockNumber, "logIndex", log.LogIndex)
			continue
		}
		lid := logID(log)
		if _, ok := q.states[lid]; ok {
			// q.lggr.Debugw("Skipping known log", "blockThreshold", blockThreshold, "logBlock", log.BlockNumber, "logIndex", log.LogIndex)
			continue
		}
		q.states[lid] = logTriggerStateEntry{state: logTriggerStateEnqueued, block: log.BlockNumber}
		added++
		logs = append(logs, log)
	}
	q.logs = logs

	var dropped int
	if added > 0 {
		q.orderLogs()
		dropped = q.clean(blockThreshold)
		q.lggr.Debugw("Enqueued logs", "added", added, "dropped", dropped, "blockThreshold", blockThreshold, "q size", len(q.logs), "visited size", len(q.states))
	}

	prommetrics.AutomationLogBufferFlow.WithLabelValues(prommetrics.LogBufferFlowDirectionIngress).Add(float64(added))
	prommetrics.AutomationLogBufferFlow.WithLabelValues(prommetrics.LogBufferFlowDirectionDropped).Add(float64(dropped))

	return added, dropped
}

// orderLogs sorts the logs in the buffer.
// NOTE: this method is not thread safe and should be called within a lock.
func (q *upkeepLogQueue) orderLogs() {
	// sort logs by block number, tx hash and log index
	// to keep the q sorted and to ensure that logs can be
	// grouped by block windows for the cleanup
	sort.SliceStable(q.logs, func(i, j int) bool {
		return LogSorter(q.logs[i], q.logs[j])
	})
}

// clean removes logs that are older than blockThreshold and drops logs if the limit for the
// given upkeep was exceeded. Returns the number of logs that were dropped.
// NOTE: this method is not thread safe and should be called within a lock.
func (q *upkeepLogQueue) clean(blockThreshold int64) int {
	var dropped, expired int
	blockRate := int(q.opts.blockRate.Load())
	windowLimit := int(q.opts.windowLimit.Load())
	updated := make([]logpoller.Log, 0)
	// helper variables to keep track of the current window capacity
	currentWindowCapacity, currentWindowStart := 0, int64(0)
	for _, l := range q.logs {
		if blockThreshold > l.BlockNumber { // old log, removed
			prommetrics.AutomationLogBufferFlow.WithLabelValues(prommetrics.LogBufferFlowDirectionExpired).Inc()
			// q.lggr.Debugw("Expiring old log", "blockNumber", l.BlockNumber, "blockThreshold", blockThreshold, "logIndex", l.LogIndex)
			logid := logID(l)
			delete(q.states, logid)
			expired++
			continue
		}
		start, _ := getBlockWindow(l.BlockNumber, blockRate)
		if start != currentWindowStart {
			// new window, reset capacity
			currentWindowStart = start
			currentWindowCapacity = 0
		}
		currentWindowCapacity++
		// if capacity has been reached, drop the log
		if currentWindowCapacity > windowLimit {
			lid := logID(l)
			if s, ok := q.states[lid]; ok {
				s.state = logTriggerStateDropped
				q.states[lid] = s
			}
			dropped++
			prommetrics.AutomationLogBufferFlow.WithLabelValues(prommetrics.LogBufferFlowDirectionDropped).Inc()
			q.lggr.Debugw("Reached log buffer limits, dropping log", "blockNumber", l.BlockNumber,
				"blockHash", l.BlockHash, "txHash", l.TxHash, "logIndex", l.LogIndex, "len updated", len(updated),
				"currentWindowStart", currentWindowStart, "currentWindowCapacity", currentWindowCapacity,
				"maxLogsPerWindow", windowLimit, "blockRate", blockRate)
			continue
		}
		updated = append(updated, l)
	}

	if dropped > 0 || expired > 0 {
		q.lggr.Debugw("Cleaned logs", "dropped", dropped, "expired", expired, "blockThreshold", blockThreshold, "len updated", len(updated), "len before", len(q.logs))
		q.logs = updated
	}

	q.cleanStates(blockThreshold)

	return dropped
}

// cleanStates removes states that are older than blockThreshold.
// NOTE: this method is not thread safe and should be called within a lock.
func (q *upkeepLogQueue) cleanStates(blockThreshold int64) {
	for lid, s := range q.states {
		if s.block <= blockThreshold {
			delete(q.states, lid)
		}
	}
}

// getBlockWindow returns the start and end block of the window for the given block.
func getBlockWindow(block int64, blockRate int) (start int64, end int64) {
	windowSize := int64(blockRate)
	if windowSize == 0 {
		return block, block
	}
	start = block - (block % windowSize)
	end = start + windowSize - 1
	return
}
