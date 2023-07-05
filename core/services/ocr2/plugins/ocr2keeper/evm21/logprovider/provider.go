package logprovider

import (
	"context"
	"crypto/sha256"
	"fmt"
	"hash"
	"math/big"
	"runtime"
	"sync"
	"time"

	"go.uber.org/multierr"
	"golang.org/x/time/rate"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

const (
	BlockLimitExceeded = "block limit exceeded"
	logTriggerType     = 1
)

var (
	ErrHeadNotAvailable = fmt.Errorf("head not available")
)

// LogTriggerConfig is an alias for log trigger config.
type LogTriggerConfig i_keeper_registry_master_wrapper_2_1.KeeperRegistryBase21LogTriggerConfig

// upkeepFilterEntry holds the upkeep filter, rate limiter and last polled block.
type upkeepFilterEntry struct {
	id     *big.Int
	filter logpoller.Filter
	cfg    LogTriggerConfig
	// lastPollBlock is the last block number the logs were fetched for this upkeep
	lastPollBlock int64
	// blockLimiter is used to limit the number of blocks to fetch logs for an upkeep
	blockLimiter *rate.Limiter
}

type LogEventProvider interface {
	// Start starts the log event provider.
	Start(ctx context.Context) error
	// Close closes the log event provider.
	Close() error
	// RegisterFilter registers the filter (if valid) for the given upkeepID.
	RegisterFilter(upkeepID *big.Int, cfg LogTriggerConfig) error
	// UnregisterFilter removes the filter for the given upkeepID.
	UnregisterFilter(upkeepID *big.Int) error
	// GetLogs returns the logs in the given range.
	GetLogs() ([]UpkeepPayload, error)
}

type LogEventProviderTest interface {
	LogEventProvider
	ReadLogs(ctx context.Context, force bool, ids ...*big.Int) error
}

// logEventProvider manages log filters for upkeeps and enables to read the log events.
type logEventProvider struct {
	lggr logger.Logger

	cancel context.CancelFunc

	poller logpoller.LogPoller

	packer LogDataPacker

	lock   sync.RWMutex
	active map[string]upkeepFilterEntry

	buffer *logEventBuffer

	opts *LogEventProviderOptions
}

func New(lggr logger.Logger, poller logpoller.LogPoller, packer LogDataPacker, opts *LogEventProviderOptions) *logEventProvider {
	if opts == nil {
		opts = new(LogEventProviderOptions)
	}
	opts.Defaults()
	return &logEventProvider{
		packer: packer,
		lggr:   lggr.Named("KeepersRegistry.LogEventProvider"),
		buffer: newLogEventBuffer(lggr, opts.LogBufferSize, opts.BufferMaxBlockSize, opts.AllowedLogsPerBlock),
		poller: poller,
		lock:   sync.RWMutex{},
		active: make(map[string]upkeepFilterEntry),
		opts:   opts,
	}
}

func (p *logEventProvider) Start(pctx context.Context) error {
	ctx, cancel := context.WithCancel(pctx)
	defer cancel()

	p.lock.Lock()
	p.cancel = cancel
	p.lock.Unlock()

	readQ := make(chan []*big.Int, 32)

	for i := 0; i < p.opts.Readers; i++ {
		go p.startReader(ctx, readQ)
	}

	return p.scheduleReadJobs(ctx, func(ids []*big.Int) {
		select {
		case readQ <- ids:
		case <-ctx.Done():
		default:
			p.lggr.Warnw("readQ is full, dropping ids", "ids", ids)
		}
	})
}

func (p *logEventProvider) Close() error {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.active = make(map[string]upkeepFilterEntry)
	if p.cancel != nil {
		p.cancel()
	}
	return nil
}

func (p *logEventProvider) GetLogs() ([]UpkeepPayload, error) {
	latest := p.buffer.latestBlockSeen()
	diff := latest - p.opts.LogBlocksLookback
	if diff < 0 {
		diff = latest
	}
	logs := p.buffer.dequeue(int(diff))

	var payloads []UpkeepPayload
	for _, l := range logs {
		log := l.log
		logExtension := fmt.Sprintf("%s:%d", log.TxHash.Hex(), uint(log.LogIndex))
		trig := NewTrigger(log.BlockNumber, log.BlockHash.Hex(), logExtension)
		checkData, err := p.packer.PackLogData(log)
		if err != nil {
			p.lggr.Warnw("failed to pack log data", "err", err, "log", log)
			continue
		}
		payload := NewUpkeepPayload(l.id, logTriggerType, trig, checkData)
		payloads = append(payloads, payload)
	}

	return payloads, nil
}

// ReadLogs fetches the logs for the given upkeeps.
func (p *logEventProvider) ReadLogs(ctx context.Context, force bool, ids ...*big.Int) error {
	latest, err := p.poller.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrHeadNotAvailable, err)
	}
	entries := p.getEntries(latest, force, ids...)

	err = p.readLogs(ctx, latest, entries...)
	p.updateEntriesLastPoll(entries)
	// p.lggr.Debugw("read logs for entries", "latestBlock", latest, "entries", len(entries), "err", err)
	if err != nil {
		return fmt.Errorf("fetched logs with errors: %w", err)
	}

	return nil
}

// scheduleReadJobs starts a scheduler that pushed ids to readQ for reading logs in the background.
func (p *logEventProvider) scheduleReadJobs(pctx context.Context, execute func([]*big.Int)) error {
	ctx, cancel := context.WithCancel(pctx)
	defer cancel()

	ticker := time.NewTicker(p.opts.ReadInterval)
	defer ticker.Stop()

	h := sha256.New()
	partitionIdx := 0

	for {
		select {
		case <-ticker.C:
			ids := p.getPartitionIds(h, partitionIdx)
			if len(ids) > 0 {
				maxBatchSize := p.opts.ReadMaxBatchSize
				for len(ids) > maxBatchSize {
					batch := ids[:maxBatchSize]
					execute(batch)
					ids = ids[maxBatchSize:]
					runtime.Gosched()
				}
				execute(ids)
			}
			partitionIdx++
		case <-ctx.Done():
			return nil
		}
	}
}

// startReader starts a reader that reads logs from the ids coming from readQ.
func (p *logEventProvider) startReader(pctx context.Context, readQ <-chan []*big.Int) {
	ctx, cancel := context.WithCancel(pctx)
	defer cancel()

	lggr := p.lggr.With("where", "reader")

	for {
		select {
		case batch := <-readQ:
			if err := p.ReadLogs(ctx, true, batch...); err != nil {
				lggr.Warnw("failed to read logs", "err", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

// getPartitionIds returns the upkeepIDs for the given partition and the number of partitions.
// Partitioning is done by hashing the upkeepID and taking the modulus of the number of partitions.
func (p *logEventProvider) getPartitionIds(hashFn hash.Hash, partition int) []*big.Int {
	p.lock.RLock()
	defer p.lock.RUnlock()

	numOfPartitions := len(p.active) / p.opts.ReadMaxBatchSize
	if numOfPartitions < 1 {
		numOfPartitions = 1
	}
	partition = partition % numOfPartitions

	var ids []*big.Int
	for _, entry := range p.active {
		if len(entry.filter.Addresses) == 0 {
			continue
		}
		n, err := hashFn.Write(entry.filter.Addresses[0].Bytes())
		if err != nil || n == 0 {
			p.lggr.Warnw("failed to hash upkeep address", "err", err, "addr", entry.filter.Addresses[0])
			continue
		}
		h := hashFn.Sum(nil)
		// taking only 6 bytes to avoid working with big numbers
		i := big.NewInt(0).SetBytes(h[len(h)-6:])
		if int(i.Int64())%numOfPartitions == partition {
			ids = append(ids, entry.id)
		}
		hashFn.Reset()
	}

	return ids
}

func (p *logEventProvider) updateEntriesLastPoll(entries []*upkeepFilterEntry) {
	p.lock.Lock()
	defer p.lock.Unlock()

	for _, entry := range entries {
		// for successful queries, the last poll block was updated
		orig, ok := p.active[entry.id.String()]
		if !ok {
			continue // entry was removed
		}
		if entry.lastPollBlock == orig.lastPollBlock {
			continue
		}
		orig.lastPollBlock = entry.lastPollBlock
		p.active[entry.id.String()] = orig
	}
}

// getEntries returns the filters for the given upkeepIDs,
// returns empty filter for inactive upkeeps.
func (p *logEventProvider) getEntries(latestBlock int64, force bool, ids ...*big.Int) []*upkeepFilterEntry {
	p.lock.RLock()
	defer p.lock.RUnlock()

	var filters []*upkeepFilterEntry
	for _, id := range ids {
		entry, ok := p.active[id.String()]
		if !ok { // entry not found, could be inactive upkeep
			p.lggr.Debugw("upkeep filter not found", "upkeep", id.String())
			filters = append(filters, &upkeepFilterEntry{id: id})
			continue
		}
		if !force && entry.lastPollBlock > latestBlock {
			p.lggr.Debugw("already polled latest block", "entry.lastPollBlock", entry.lastPollBlock, "latestBlock", latestBlock, "upkeep", id.String())
			filters = append(filters, &upkeepFilterEntry{id: id, lastPollBlock: entry.lastPollBlock})
			continue
		}
		// recreating the struct to be thread safe
		filters = append(filters, &upkeepFilterEntry{
			id:            id,
			filter:        p.newLogFilter(id, entry.cfg),
			lastPollBlock: entry.lastPollBlock,
			blockLimiter:  entry.blockLimiter,
		})
	}

	return filters
}

// readLogs calls log poller to get the logs for the given upkeep entries.
// we use p.opts.LookbackBuffer to check for reorgs based logs.
//
// TODO: batch entries by contract address and call log poller once per contract address
// NOTE: the entries are already grouped by contract address
func (p *logEventProvider) readLogs(ctx context.Context, latest int64, entries ...*upkeepFilterEntry) (merr error) {
	// mainLggr := p.lggr.With("latestBlock", latest)
	logBlocksLookback := p.opts.LogBlocksLookback
	maxBurst := int(logBlocksLookback*2 + 1)

	for _, entry := range entries {
		if len(entry.filter.Addresses) == 0 {
			continue
		}
		// lggr := mainLggr.With("upkeep", entry.id.String(), "addrs", entry.filter.Addresses, "sigs", entry.filter.EventSigs)
		start := entry.lastPollBlock
		if start == 0 || start < latest-logBlocksLookback {
			// long range or first time polling,
			// using a larger lookback and burst
			start = latest - logBlocksLookback*2
			entry.blockLimiter.SetBurst(maxBurst)
		}
		resv := entry.blockLimiter.ReserveN(time.Now(), int(latest-start))
		if !resv.OK() {
			merr = multierr.Append(merr, fmt.Errorf("%s: %s", BlockLimitExceeded, entry.id.String()))
			continue
		}
		start = start - p.opts.LookbackBuffer // adding a buffer to check for reorgs
		if start < 0 {
			start = 0
		}
		// lggr = lggr.With("startBlock", start)
		logs, err := p.poller.LogsWithSigs(start, latest, entry.filter.EventSigs, entry.filter.Addresses[0], pg.WithParentCtx(ctx))
		if err != nil {
			resv.Cancel() // cancels limit reservation as we failed to get logs
			if ctx.Err() != nil {
				return multierr.Append(merr, ctx.Err())
			}
			merr = multierr.Append(merr, fmt.Errorf("failed to get logs for upkeep %s: %w", entry.id.String(), err))
			continue
		}
		// if this limiter's burst was set to the max,
		// we need to reset it
		if entry.blockLimiter.Burst() == maxBurst {
			resv.Cancel() // cancel the reservation as we are resetting the burst
			entry.blockLimiter.SetBurst(p.opts.BlockLimitBurst)
		}
		added := p.buffer.enqueue(entry.id, logs...)
		// if we added logs or couldn't find, update the last poll block
		if added > 0 || len(logs) == 0 {
			entry.lastPollBlock = latest
		}
		// if n := len(logs); n > 0 {
		// 	lggr.Debugw("got logs for upkeep", "logs", n, "added", added)
		// }
	}

	return merr
}
