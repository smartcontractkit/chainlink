package logprovider

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"io"
	"math/big"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_utils_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/core"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	LogProviderServiceName = "LogEventProvider"

	ErrHeadNotAvailable   = fmt.Errorf("head not available")
	ErrBlockLimitExceeded = fmt.Errorf("block limit exceeded")

	// AllowedLogsPerUpkeep is the maximum number of logs allowed per upkeep every single call.
	AllowedLogsPerUpkeep = 5
	// MaxPayloads is the maximum number of payloads to return per call.
	MaxPayloads = 100

	readJobQueueSize = 64
	readLogsTimeout  = 10 * time.Second

	readMaxBatchSize = 32
	// reorgBuffer is the number of blocks to add as a buffer to the block range when reading logs.
	reorgBuffer   = int64(32)
	readerThreads = 4
)

// LogTriggerConfig is an alias for log trigger config.
type LogTriggerConfig automation_utils_2_1.LogTriggerConfig

type FilterOptions struct {
	UpkeepID      *big.Int
	TriggerConfig LogTriggerConfig
	UpdateBlock   uint64
}

type LogTriggersLifeCycle interface {
	// RegisterFilter registers the filter (if valid) for the given upkeepID.
	RegisterFilter(ctx context.Context, opts FilterOptions) error
	// UnregisterFilter removes the filter for the given upkeepID.
	UnregisterFilter(upkeepID *big.Int) error
}
type LogEventProvider interface {
	ocr2keepers.LogEventProvider
	LogTriggersLifeCycle

	RefreshActiveUpkeeps(ids ...*big.Int) ([]*big.Int, error)

	Start(context.Context) error
	io.Closer
}

type LogEventProviderTest interface {
	LogEventProvider
	ReadLogs(ctx context.Context, ids ...*big.Int) error
	CurrentPartitionIdx() uint64
}

var _ LogEventProvider = &logEventProvider{}
var _ LogEventProviderTest = &logEventProvider{}

// logEventProvider manages log filters for upkeeps and enables to read the log events.
type logEventProvider struct {
	services.StateMachine
	threadCtrl utils.ThreadControl

	lggr logger.Logger

	poller logpoller.LogPoller

	packer LogDataPacker

	lock         sync.RWMutex
	registerLock sync.Mutex

	filterStore UpkeepFilterStore
	buffer      *logEventBuffer

	opts LogTriggersOptions

	currentPartitionIdx uint64
}

func NewLogProvider(lggr logger.Logger, poller logpoller.LogPoller, packer LogDataPacker, filterStore UpkeepFilterStore, opts LogTriggersOptions) *logEventProvider {
	return &logEventProvider{
		threadCtrl:  utils.NewThreadControl(),
		lggr:        lggr.Named("KeepersRegistry.LogEventProvider"),
		packer:      packer,
		buffer:      newLogEventBuffer(lggr, int(opts.LookbackBlocks), defaultNumOfLogUpkeeps, defaultFastExecLogsHigh),
		poller:      poller,
		opts:        opts,
		filterStore: filterStore,
	}
}

func (p *logEventProvider) Start(context.Context) error {
	return p.StartOnce(LogProviderServiceName, func() error {

		readQ := make(chan []*big.Int, readJobQueueSize)

		p.lggr.Infow("starting log event provider", "readInterval", p.opts.ReadInterval, "readMaxBatchSize", readMaxBatchSize, "readers", readerThreads)

		for i := 0; i < readerThreads; i++ {
			p.threadCtrl.Go(func(ctx context.Context) {
				p.startReader(ctx, readQ)
			})
		}

		p.threadCtrl.Go(func(ctx context.Context) {
			lggr := p.lggr.With("where", "scheduler")

			p.scheduleReadJobs(ctx, func(ids []*big.Int) {
				select {
				case readQ <- ids:
				case <-ctx.Done():
				default:
					lggr.Warnw("readQ is full, dropping ids", "ids", ids)
				}
			})
		})

		return nil
	})
}

func (p *logEventProvider) Close() error {
	return p.StopOnce(LogProviderServiceName, func() error {
		p.threadCtrl.Close()
		return nil
	})
}

func (p *logEventProvider) HealthReport() map[string]error {
	return map[string]error{LogProviderServiceName: p.Healthy()}
}

func (p *logEventProvider) GetLatestPayloads(ctx context.Context) ([]ocr2keepers.UpkeepPayload, error) {
	latest, err := p.poller.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrHeadNotAvailable, err)
	}
	start := latest.BlockNumber - p.opts.LookbackBlocks
	if start <= 0 {
		start = 1
	}
	logs := p.buffer.dequeueRange(start, latest.BlockNumber, AllowedLogsPerUpkeep, MaxPayloads)

	// p.lggr.Debugw("got latest logs from buffer", "latest", latest, "diff", diff, "logs", len(logs))

	var payloads []ocr2keepers.UpkeepPayload
	for _, l := range logs {
		log := l.log
		trig := logToTrigger(log)
		checkData, err := p.packer.PackLogData(log)
		if err != nil {
			p.lggr.Warnw("failed to pack log data", "err", err, "log", log)
			continue
		}
		payload, err := core.NewUpkeepPayload(l.upkeepID, trig, checkData)
		if err != nil {
			p.lggr.Warnw("failed to create upkeep payload", "err", err, "id", l.upkeepID, "trigger", trig, "checkData", checkData)
			continue
		}

		payloads = append(payloads, payload)
	}

	return payloads, nil
}

// ReadLogs fetches the logs for the given upkeeps.
func (p *logEventProvider) ReadLogs(pctx context.Context, ids ...*big.Int) error {
	ctx, cancel := context.WithTimeout(pctx, readLogsTimeout)
	defer cancel()

	latest, err := p.poller.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrHeadNotAvailable, err)
	}
	if latest.BlockNumber == 0 {
		return fmt.Errorf("%w: %s", ErrHeadNotAvailable, "latest block is 0")
	}
	filters := p.getFilters(latest.BlockNumber, ids...)

	err = p.readLogs(ctx, latest.BlockNumber, filters)
	p.updateFiltersLastPoll(filters)
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
func (p *logEventProvider) scheduleReadJobs(pctx context.Context, execute func([]*big.Int)) {
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
				maxBatchSize := readMaxBatchSize
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
			return
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
			if err := p.ReadLogs(ctx, batch...); err != nil {
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
	numOfPartitions := p.filterStore.Size() / readMaxBatchSize
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

func (p *logEventProvider) updateFiltersLastPoll(entries []upkeepFilter) {
	p.filterStore.UpdateFilters(func(orig, f upkeepFilter) upkeepFilter {
		if f.lastPollBlock > orig.lastPollBlock {
			orig.lastPollBlock = f.lastPollBlock
			if f.lastPollBlock%10 == 0 {
				// print log occasionally to avoid spamming logs
				p.lggr.Debugw("Updated lastPollBlock", "lastPollBlock", f.lastPollBlock, "upkeepID", f.upkeepID)
			}
		}
		return orig
	}, entries...)
}

// getFilters returns the filters for the given upkeepIDs,
// returns empty filter for inactive upkeeps.
func (p *logEventProvider) getFilters(latestBlock int64, ids ...*big.Int) []upkeepFilter {
	var filters []upkeepFilter
	p.filterStore.RangeFiltersByIDs(func(i int, f upkeepFilter) {
		if len(f.addr) == 0 { // not found
			p.lggr.Debugw("upkeep filter not found", "upkeep", f.upkeepID.String())
			filters = append(filters, f)
			return
		}
		if f.configUpdateBlock > uint64(latestBlock) {
			p.lggr.Debugw("upkeep config update block was created after latestBlock", "upkeep", f.upkeepID.String(), "configUpdateBlock", f.configUpdateBlock, "latestBlock", latestBlock)
			filters = append(filters, upkeepFilter{upkeepID: f.upkeepID})
			return
		}
		if f.lastPollBlock > latestBlock {
			p.lggr.Debugw("already polled latest block", "entry.lastPollBlock", f.lastPollBlock, "latestBlock", latestBlock, "upkeep", f.upkeepID.String())
			filters = append(filters, upkeepFilter{upkeepID: f.upkeepID})
			return
		}
		filters = append(filters, f.Clone())
	}, ids...)

	return filters
}

// readLogs calls log poller to get the logs for the given upkeep entries.
//
// Exploratory: batch filters by contract address and call log poller once per contract address
// NOTE: the filters are already grouped by contract address
func (p *logEventProvider) readLogs(ctx context.Context, latest int64, filters []upkeepFilter) (merr error) {
	lookbackBlocks := p.opts.LookbackBlocks
	if latest < lookbackBlocks {
		// special case of a new blockchain (e.g. simulated chain)
		lookbackBlocks = latest - 1
	}
	// maxBurst will be used to increase the burst limit to allow a long range scan
	maxBurst := int(lookbackBlocks + 1)

	for i, filter := range filters {
		if len(filter.addr) == 0 {
			continue
		}
		start := filter.lastPollBlock
		// range should not exceed [lookbackBlocks, latest]
		if start < latest-lookbackBlocks {
			start = latest - lookbackBlocks
			filter.blockLimiter.SetBurst(maxBurst)
		}

		resv := filter.blockLimiter.ReserveN(time.Now(), int(latest-start))
		if !resv.OK() {
			merr = errors.Join(merr, fmt.Errorf("%w: %s", ErrBlockLimitExceeded, filter.upkeepID.String()))
			continue
		}
		// adding a buffer to check for reorged logs.
		start = start - reorgBuffer
		// make sure start of the range is not before the config update block
		if configUpdateBlock := int64(filter.configUpdateBlock); start < configUpdateBlock {
			start = configUpdateBlock
		}
		// query logs based on contract address, event sig, and blocks
		logs, err := p.poller.LogsWithSigs(start, latest, []common.Hash{filter.topics[0]}, common.BytesToAddress(filter.addr), pg.WithParentCtx(ctx))
		if err != nil {
			// cancel limit reservation as we failed to get logs
			resv.Cancel()
			if ctx.Err() != nil {
				// exit if the context was canceled
				return merr
			}
			merr = errors.Join(merr, fmt.Errorf("failed to get logs for upkeep %s: %w", filter.upkeepID.String(), err))
			continue
		}
		filteredLogs := filter.Select(logs...)

		// if this limiter's burst was set to the max ->
		// reset it and cancel the reservation to allow further processing
		if filter.blockLimiter.Burst() == maxBurst {
			resv.Cancel()
			filter.blockLimiter.SetBurst(p.opts.BlockLimitBurst)
		}

		p.buffer.enqueue(filter.upkeepID, filteredLogs...)

		// Update the lastPollBlock for filter in slice this is then
		// updated into filter store in updateFiltersLastPoll
		filters[i].lastPollBlock = latest
	}

	return merr
}
