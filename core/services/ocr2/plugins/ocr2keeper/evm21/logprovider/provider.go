package logprovider

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"math/big"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
	keepersflows "github.com/smartcontractkit/ocr2keepers/pkg/v3/flows"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_utils_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

const (
	logTriggerType = 1
)

var (
	ErrHeadNotAvailable   = fmt.Errorf("head not available")
	ErrBlockLimitExceeded = fmt.Errorf("block limit exceeded")
)

// LogTriggerConfig is an alias for log trigger config.
type LogTriggerConfig automation_utils_2_1.LogTriggerConfig

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
	GetLogs(context.Context) ([]ocr2keepers.UpkeepPayload, error)
}

type LogEventProviderTest interface {
	LogEventProvider
	ReadLogs(ctx context.Context, force bool, ids ...*big.Int) error
	CurrentPartitionIdx() uint64
}

var _ keepersflows.PayloadBuilder = &logEventProvider{}
var _ keepersflows.LogEventProvider = &logEventProvider{}

// logEventProvider manages log filters for upkeeps and enables to read the log events.
type logEventProvider struct {
	lggr logger.Logger

	cancel context.CancelFunc

	poller logpoller.LogPoller

	packer LogDataPacker

	lock sync.RWMutex

	filterStore UpkeepFilterStore
	buffer      *logEventBuffer

	opts *LogEventProviderOptions

	currentPartitionIdx uint64
}

func New(lggr logger.Logger, poller logpoller.LogPoller, packer LogDataPacker, filterStore UpkeepFilterStore, opts *LogEventProviderOptions) *logEventProvider {
	if opts == nil {
		opts = new(LogEventProviderOptions)
	}
	opts.Defaults()
	return &logEventProvider{
		packer:      packer,
		lggr:        lggr.Named("KeepersRegistry.LogEventProvider"),
		buffer:      newLogEventBuffer(lggr, opts.LogBufferSize, opts.BufferMaxBlockSize, opts.AllowedLogsPerBlock),
		poller:      poller,
		lock:        sync.RWMutex{},
		opts:        opts,
		filterStore: filterStore,
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

	if p.cancel != nil {
		p.cancel()
	}
	return nil
}

func (p *logEventProvider) BuildPayload(ctx context.Context, proposal ocr2keepers.CoordinatedProposal) (ocr2keepers.UpkeepPayload, error) {
	// TODO: implement
	return ocr2keepers.UpkeepPayload{}, nil
}

func (p *logEventProvider) GetLogs(context.Context) ([]ocr2keepers.UpkeepPayload, error) {
	latest := p.buffer.latestBlockSeen()
	diff := latest - p.opts.LogBlocksLookback
	if diff < 0 {
		diff = latest
	}
	logs := p.buffer.dequeue(int(diff))

	var payloads []ocr2keepers.UpkeepPayload
	for _, l := range logs {
		log := l.log
		trig := ocr2keepers.NewTrigger(
			log.BlockNumber,
			log.BlockHash.Hex(),
			core.LogTriggerExtension{
				TxHash:   log.TxHash.Hex(),
				LogIndex: log.LogIndex,
			},
		)
		checkData, err := p.packer.PackLogData(log)
		if err != nil {
			p.lggr.Warnw("failed to pack log data", "err", err, "log", log)
			continue
		}

		payload, err := core.NewUpkeepPayload(l.id, logTriggerType, trig, checkData)
		if err != nil {
			// skip invalid payloads
			continue
		}

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

	err = p.readLogs(ctx, latest, entries)
	p.updateEntriesLastPoll(entries)
	// p.lggr.Debugw("read logs for entries", "latestBlock", latest, "entries", len(entries), "err", err)
	if err != nil {
		return fmt.Errorf("fetched logs with errors: %w", err)
	}

	return nil
}

func (p *logEventProvider) CurrentPartitionIdx() uint64 {
	return atomic.LoadUint64(&p.currentPartitionIdx)
}

// scheduleReadJobs starts a scheduler that pushed ids to readQ for reading logs in the background.
func (p *logEventProvider) scheduleReadJobs(pctx context.Context, execute func([]*big.Int)) error {
	ctx, cancel := context.WithCancel(pctx)
	defer cancel()

	ticker := time.NewTicker(p.opts.ReadInterval)
	defer ticker.Stop()

	h := sha256.New()

	partitionIdx := p.CurrentPartitionIdx()

	for {
		select {
		case <-ticker.C:
			ids := p.getPartitionIds(h, int(partitionIdx))
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
			atomic.StoreUint64(&p.currentPartitionIdx, partitionIdx)
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
	numOfPartitions := p.filterStore.Size() / p.opts.ReadMaxBatchSize
	if numOfPartitions < 1 {
		numOfPartitions = 1
	}
	partition = partition % numOfPartitions

	ids := p.filterStore.GetIDs(func(f upkeepFilter) bool {
		if len(f.addr) == 0 {
			return false
		}
		n, err := hashFn.Write(f.addr)
		if err != nil || n == 0 {
			p.lggr.Warnw("failed to hash upkeep address", "err", err, "addr", hexutil.Encode(f.addr))
			return false
		}
		h := hashFn.Sum(nil)
		defer hashFn.Reset()
		// taking only 6 bytes to avoid working with big numbers
		i := big.NewInt(0).SetBytes(h[len(h)-6:])
		return int(i.Int64())%numOfPartitions == partition
	})

	return ids
}

func (p *logEventProvider) updateEntriesLastPoll(entries []upkeepFilter) {
	p.filterStore.UpdateFilters(func(orig, f upkeepFilter) upkeepFilter {
		if f.lastPollBlock > orig.lastPollBlock {
			orig.lastPollBlock = f.lastPollBlock
		}
		return orig
	}, entries...)
}

// getEntries returns the filters for the given upkeepIDs,
// returns empty filter for inactive upkeeps.
func (p *logEventProvider) getEntries(latestBlock int64, force bool, ids ...*big.Int) []upkeepFilter {
	var filters []upkeepFilter
	p.filterStore.RangeFiltersByIDs(func(i int, f upkeepFilter) {
		if len(f.addr) == 0 { // not found
			p.lggr.Debugw("upkeep filter not found", "upkeep", f.upkeepID.String())
			filters = append(filters, f)
			return
		}
		if !force && f.lastPollBlock > latestBlock {
			p.lggr.Debugw("already polled latest block", "entry.lastPollBlock", f.lastPollBlock, "latestBlock", latestBlock, "upkeep", f.upkeepID.String())
			filters = append(filters, upkeepFilter{upkeepID: f.upkeepID})
			return
		}
		// cloning struct to be thread safe
		topics := make([]common.Hash, len(f.topics))
		copy(topics, f.topics)
		addr := make([]byte, len(f.addr))
		copy(addr, f.addr)
		filters = append(filters, upkeepFilter{
			upkeepID:        f.upkeepID,
			topics:          topics,
			addr:            addr,
			lastPollBlock:   f.lastPollBlock,
			lastRePollBlock: f.lastRePollBlock,
			blockLimiter:    f.blockLimiter,
		})
	}, ids...)

	return filters
}

// readLogs calls log poller to get the logs for the given upkeep entries.
// we use p.opts.LookbackBuffer to check for reorgs based logs.
//
// TODO: batch entries by contract address and call log poller once per contract address
// NOTE: the entries are already grouped by contract address
func (p *logEventProvider) readLogs(ctx context.Context, latest int64, entries []upkeepFilter) (merr error) {
	// mainLggr := p.lggr.With("latestBlock", latest)
	logBlocksLookback := p.opts.LogBlocksLookback
	maxBurst := int(logBlocksLookback*2 + 1)

	for _, entry := range entries {
		if len(entry.addr) == 0 {
			continue
		}

		// start should either be the last block polled for the entry or the
		// lookback range in the case this is the first time the entry is polled
		start := entry.lastPollBlock
		if start == 0 || start < latest-logBlocksLookback {
			// long range or first time polling,
			// using a larger lookback and burst
			start = latest - logBlocksLookback*2

			// start should not be negative in the special case of an empty or
			// new blockchain (this is the case for the simulated chain)
			if start < 0 {
				start = 0
			}

			// override the burst limit to the max to allow the long range scan
			entry.blockLimiter.SetBurst(maxBurst)
		}

		// apply a reservation for the entry block range
		// this reservation may be cancelled later in the case of overriding
		// the block limitation with maxBurst
		// fmt.Printf("reserving: %d for %s\n", int(latest-start), entry.upkeepID)
		resv := entry.blockLimiter.ReserveN(time.Now(), int(latest-start))
		// fmt.Printf("reservation limit for %s: %d\n", entry.upkeepID, resv.DelayFrom(time.Now())/time.Millisecond)
		if !resv.OK() {
			// the block range cannot be processed for the event either because
			// it is above the configured rate or burst limit
			merr = errors.Join(merr, fmt.Errorf("%w: %s", ErrBlockLimitExceeded, entry.upkeepID.String()))

			continue
		}

		// adding a buffer outside the reserved block range to check for reorgs
		start = start - p.opts.LookbackBuffer
		if start < 0 {
			start = 0
		}

		logs, err := p.poller.LogsWithSigs(start, latest, entry.topics, common.BytesToAddress(entry.addr), pg.WithParentCtx(ctx))
		if err != nil {
			// cancel limit reservation as we failed to get logs
			resv.Cancel()

			if ctx.Err() != nil {
				return errors.Join(merr, ctx.Err())
			}

			merr = errors.Join(merr, fmt.Errorf("failed to get logs for upkeep %s: %w", entry.upkeepID.String(), err))

			// continue processing entries as this is not a hard error
			continue
		}

		// if this limiter's burst was set to the max reset it and cancel the
		// reservation to allow further processing without exceeding the rate
		// limit
		if entry.blockLimiter.Burst() == maxBurst {
			// cancel the reservation as we are resetting the burst
			resv.Cancel()

			entry.blockLimiter.SetBurst(p.opts.BlockLimitBurst)
		}

		// add logs returned from the poller to the queue and return the total
		// number of unique logs added to the queue
		// added :=
		p.buffer.enqueue(entry.upkeepID, logs...)

		// fmt.Printf("added: %d; logs: %d\n", added, len(logs))

		// if no new logs were added or no logs were found, update the last poll block
		// if added > 0 || len(logs) == 0 {
		entry.lastPollBlock = latest

		p.filterStore.UpdateFilters(func(original upkeepFilter, updated upkeepFilter) upkeepFilter {
			original.lastPollBlock = updated.lastPollBlock

			// fmt.Printf("updating %s to %d\n", updated.upkeepID, updated.lastPollBlock)

			return original
		}, entry)
		// }
	}

	return merr
}
