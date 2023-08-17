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
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_utils_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
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
	// GetLatestPayloads returns the logs in the given range.
	GetLatestPayloads(context.Context) ([]ocr2keepers.UpkeepPayload, error)
}

type LogEventProviderTest interface {
	LogEventProvider
	ReadLogs(ctx context.Context, force bool, ids ...*big.Int) error
	CurrentPartitionIdx() uint64
}

var _ ocr2keepers.PayloadBuilder = &logEventProvider{}
var _ ocr2keepers.LogEventProvider = &logEventProvider{}

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

func (p *logEventProvider) BuildPayloads(ctx context.Context, proposals ...ocr2keepers.CoordinatedBlockProposal) ([]ocr2keepers.UpkeepPayload, error) {
	// TODO: implement
	return []ocr2keepers.UpkeepPayload{}, nil
}

func (p *logEventProvider) GetLatestPayloads(context.Context) ([]ocr2keepers.UpkeepPayload, error) {
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
			ocr2keepers.BlockNumber(log.BlockNumber),
			log.BlockHash,
		)
		trig.LogTriggerExtension = &ocr2keepers.LogTriggerExtension{
			TxHash: log.TxHash,
			Index:  uint32(log.LogIndex),
		}
		checkData, err := p.packer.PackLogData(log)
		if err != nil {
			p.lggr.Warnw("failed to pack log data", "err", err, "log", log)
			continue
		}

		payload, err := core.NewUpkeepPayload(l.id, trig, checkData)
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
				if ctx.Err() != nil {
					return
				}
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
//
// TODO: batch entries by contract address and call log poller once per contract address
// NOTE: the entries are already grouped by contract address
func (p *logEventProvider) readLogs(ctx context.Context, latest int64, entries []upkeepFilter) (merr error) {
	lookbackBlocks := p.opts.LogBlocksLookback
	if latest < lookbackBlocks {
		// special case of an empty or new blockchain (e.g. simulated chain)
		lookbackBlocks = latest
	}
	// maxBurst will be used to increase the burst limit to allow a long range scan
	maxBurst := int(lookbackBlocks + 1)

	for _, entry := range entries {
		if len(entry.addr) == 0 {
			continue
		}
		// start should either be the last block polled for the entry or the
		// lookback range in the case this is the first time the entry is polled
		start := entry.lastPollBlock
		if start == 0 || start < latest-lookbackBlocks {
			start = latest - lookbackBlocks
			entry.blockLimiter.SetBurst(maxBurst)
		}

		resv := entry.blockLimiter.ReserveN(time.Now(), int(latest-start))
		if !resv.OK() {
			merr = errors.Join(merr, fmt.Errorf("%w: %s", ErrBlockLimitExceeded, entry.upkeepID.String()))
			continue
		}
		// adding a buffer to check for reorged logs.
		start = start - p.opts.LookbackBuffer
		if start < 0 {
			start = 0
		}

		logs, err := p.poller.LogsWithSigs(start, latest, entry.topics, common.BytesToAddress(entry.addr), pg.WithParentCtx(ctx))
		if err != nil {
			// cancel limit reservation as we failed to get logs
			resv.Cancel()
			// exit if the context was canceled
			if ctx.Err() != nil {
				return merr
			}

			merr = errors.Join(merr, fmt.Errorf("failed to get logs for upkeep %s: %w", entry.upkeepID.String(), err))

			continue
		}
		// if this limiter's burst was set to the max ->
		// reset it and cancel the reservation to allow further processing
		if entry.blockLimiter.Burst() == maxBurst {
			resv.Cancel()
			entry.blockLimiter.SetBurst(p.opts.BlockLimitBurst)
		}

		p.buffer.enqueue(entry.upkeepID, logs...)

		entry.lastPollBlock = latest
	}

	return merr
}
