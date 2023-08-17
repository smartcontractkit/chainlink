package logprovider

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

var (
	ErrNotFound             = errors.New("not found")
	DefaultRecoveryInterval = 5 * time.Second
	RecoveryCacheTTL        = 24*time.Hour - time.Second

	recoveryBatchSize  = 10
	recoveryLogsBuffer = int64(10)
)

type LogRecoverer interface {
	ocr2keepers.RecoverableProvider
	ocr2keepers.PayloadBuilder

	Start(context.Context) error
	io.Closer
}

type logRecoverer struct {
	lggr logger.Logger

	cancel context.CancelFunc

	lookbackBlocks *atomic.Int64
	blockTime      *atomic.Int64

	interval time.Duration
	lock     sync.RWMutex

	pending []ocr2keepers.UpkeepPayload
	visited map[string]time.Time

	filterStore UpkeepFilterStore
	states      core.UpkeepStateReader
	packer      LogDataPacker
	poller      logpoller.LogPoller
}

var _ LogRecoverer = &logRecoverer{}

func NewLogRecoverer(lggr logger.Logger, poller logpoller.LogPoller, stateStore core.UpkeepStateReader, packer LogDataPacker, filterStore UpkeepFilterStore, interval time.Duration, lookbackBlocks int64) *logRecoverer {
	if interval == 0 {
		interval = DefaultRecoveryInterval
	}

	rec := &logRecoverer{
		lggr: lggr.Named("LogRecoverer"),

		blockTime:      &atomic.Int64{},
		lookbackBlocks: &atomic.Int64{},
		interval:       interval,

		pending:     make([]ocr2keepers.UpkeepPayload, 0),
		visited:     make(map[string]time.Time),
		poller:      poller,
		filterStore: filterStore,
		states:      stateStore,
		packer:      packer,
	}

	rec.lookbackBlocks.Store(lookbackBlocks)
	rec.blockTime.Store(defaultBlockTime)

	return rec
}

func (r *logRecoverer) Start(pctx context.Context) error {
	ctx, cancel := context.WithCancel(context.Background())

	r.lock.Lock()
	if r.cancel != nil {
		r.lock.Unlock()
		return errors.New("already started")
	}
	r.cancel = cancel
	r.lock.Unlock()

	blockTimeResolver := newBlockTimeResolver(r.poller)
	blockTime, err := blockTimeResolver.BlockTime(ctx, defaultSampleSize)
	if err != nil {
		// TODO: TBD exit or just log a warning
		// return fmt.Errorf("failed to compute block time: %w", err)
		r.lggr.Warnw("failed to compute block time", "err", err)
	}
	if blockTime > 0 {
		r.blockTime.Store(int64(blockTime))
	}

	r.lggr.Infow("starting log recoverer", "blockTime", r.blockTime.Load(), "lookbackBlocks", r.lookbackBlocks.Load(), "interval", r.interval)

	{
		go func(ctx context.Context, interval time.Duration) {
			ticker := time.NewTicker(interval)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					r.recover(ctx)
				case <-ctx.Done():
					return
				}
			}
		}(ctx, r.interval)
	}

	return nil
}

func (r *logRecoverer) Close() error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if cancel := r.cancel; cancel != nil {
		r.cancel = nil
		cancel()
	} else {
		return errors.New("already stopped")
	}
	return nil
}

func (r *logRecoverer) BuildPayload(ctx context.Context, proposal ocr2keepers.CoordinatedBlockProposal) (ocr2keepers.UpkeepPayload, error) {
	switch core.GetUpkeepType(proposal.UpkeepID) {
	case ocr2keepers.LogTrigger:
		return r.buildLogTriggerPayload(ctx, proposal)
	default:
		return ocr2keepers.UpkeepPayload{}, errors.New("not a log trigger upkeep ID")
	}
	return ocr2keepers.UpkeepPayload{}, nil
}

func (r *logRecoverer) buildLogTriggerPayload(ctx context.Context, proposal ocr2keepers.CoordinatedBlockProposal) (ocr2keepers.UpkeepPayload, error) {
	// TODO should we be querying the filter store with something other than upkeep ID?
	if r.filterStore.Has(proposal.UpkeepID.BigInt()) {

		latest, err := r.poller.LatestBlock(pg.WithParentCtx(ctx))
		if err != nil {
			return ocr2keepers.UpkeepPayload{}, err
		}

		start, offsetBlock := r.getRecoveryWindow(latest)
		block := int64(proposal.Trigger.LogTriggerExtension.BlockNumber)
		if isRecoverable := block < offsetBlock && block > start; isRecoverable {
			upkeepStates, err := r.states.SelectByWorkIDsInRange(ctx, int64(proposal.Trigger.LogTriggerExtension.BlockNumber)-1, offsetBlock, proposal.WorkID)
			if err != nil {
				return ocr2keepers.UpkeepPayload{}, err
			}

			for _, upkeepState := range upkeepStates {
				switch upkeepState {
				case ocr2keepers.Performed, ocr2keepers.Ineligible:
					return ocr2keepers.UpkeepPayload{}, nil
				default:
					// we can proceed
				}
			}

			var filter upkeepFilter
			r.filterStore.RangeFiltersByIDs(func(i int, f upkeepFilter) {
				filter = f
			}, proposal.UpkeepID.BigInt())

			if len(filter.addr) == 0 {
				return ocr2keepers.UpkeepPayload{}, errors.New("filter not found, upkeep is inactive") // TODO fix error msg
			}

			logs, err := r.poller.LogsWithSigs(int64(proposal.Trigger.LogTriggerExtension.BlockNumber)-1, offsetBlock, filter.topics, common.BytesToAddress(filter.addr), pg.WithParentCtx(ctx))
			if err != nil {
				return ocr2keepers.UpkeepPayload{}, fmt.Errorf("could not read logs: %w", err)
			}

			for _, log := range logs {
				trigger := logToTrigger(log)
				upkeepId := &ocr2keepers.UpkeepIdentifier{}
				// TODO do we need to use the filter upkeepID for correctness, or can we remove this block and use the upkeep ID on the proposal?
				ok := upkeepId.FromBigInt(filter.upkeepID)
				if !ok {
					r.lggr.Warnw("failed to convert upkeepID to UpkeepIdentifier", "upkeepID", filter.upkeepID)
					continue
				}
				wid := core.UpkeepWorkID(*upkeepId, trigger)
				if wid == proposal.WorkID {
					checkData, err := r.packer.PackLogData(log)
					if err != nil {
						r.lggr.Warnw("failed to pack log data", "err", err, "log", log)
						continue
					}

					return core.NewUpkeepPayload(proposal.UpkeepID.BigInt(), trigger, checkData)
				}
			}
		}
	}
	return ocr2keepers.UpkeepPayload{}, nil
}

func (r *logRecoverer) BuildPayloads(ctx context.Context, proposals ...ocr2keepers.CoordinatedBlockProposal) ([]ocr2keepers.UpkeepPayload, error) {
	// TODO: implement
	return []ocr2keepers.UpkeepPayload{}, nil
}

func (r *logRecoverer) GetRecoveryProposals(ctx context.Context) ([]ocr2keepers.UpkeepPayload, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if len(r.pending) == 0 {
		return nil, nil
	}

	pending := make([]ocr2keepers.UpkeepPayload, len(r.pending))
	copy(pending, r.pending)

	r.pending = make([]ocr2keepers.UpkeepPayload, 0)

	return pending, nil
}

func (r *logRecoverer) clean(ctx context.Context) {
	r.lock.Lock()
	defer r.lock.Unlock()

	cleaned := 0
	for id, t := range r.visited {
		if time.Since(t) > RecoveryCacheTTL {
			delete(r.visited, id)
			cleaned++
		}
	}

	if cleaned > 0 {
		r.lggr.Debugw("gc: cleaned visited upkeeps", "cleaned", cleaned)
	}
}

func (r *logRecoverer) recover(ctx context.Context) error {
	latest, err := r.poller.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrHeadNotAvailable, err)
	}
	start, offsetBlock := r.getRecoveryWindow(latest)
	if offsetBlock < 0 {
		// too soon to recover, we don't have enough blocks
		return nil
	}

	filters := r.getFilterBatch(offsetBlock)
	if len(filters) == 0 {
		return nil
	}

	// r.lggr.Debugw("recovering logs", "filters", filters, "startBlock", start, "offsetBlock", offsetBlock, "latestBlock", latest)

	var wg sync.WaitGroup
	for _, f := range filters {
		wg.Add(1)
		go func(f upkeepFilter) {
			defer wg.Done()
			r.recoverFilter(ctx, f, start, offsetBlock)
		}(f)
	}
	wg.Wait()

	return nil
}

// recoverFilter recovers logs for a single upkeep filter.
func (r *logRecoverer) recoverFilter(ctx context.Context, f upkeepFilter, startBlock, offsetBlock int64) error {
	start := f.lastRePollBlock
	if start < startBlock {
		start = startBlock
	}
	end := start + recoveryLogsBuffer
	if end > offsetBlock {
		end = offsetBlock
	}
	// we expect start to be > offsetBlock in any case
	logs, err := r.poller.LogsWithSigs(start, end, f.topics, common.BytesToAddress(f.addr), pg.WithParentCtx(ctx))
	if err != nil {
		return fmt.Errorf("could not read logs: %w", err)
	}

	workIDs := make([]string, 0)
	for _, log := range logs {
		trigger := logToTrigger(log)
		upkeepId := &ocr2keepers.UpkeepIdentifier{}
		ok := upkeepId.FromBigInt(f.upkeepID)
		if !ok {
			r.lggr.Warnw("failed to convert upkeepID to UpkeepIdentifier", "upkeepID", f.upkeepID)
			continue
		}
		workIDs = append(workIDs, core.UpkeepWorkID(*upkeepId, trigger))
	}

	states, err := r.states.SelectByWorkIDsInRange(ctx, start, end, workIDs...)
	if err != nil {
		return fmt.Errorf("could not read states: %w", err)
	}
	if len(logs) != len(states) {
		return fmt.Errorf("log and state count mismatch: %d != %d", len(logs), len(states))
	}
	filteredLogs := r.filterFinalizedStates(f, logs, states)

	added, alreadyPending := r.populatePending(f, filteredLogs)
	if added > 0 {
		r.lggr.Debugw("recovered logs", "count", added, "upkeepID", f.upkeepID)
	} else if alreadyPending == 0 {
		// no logs found or still in process, update the lastRePollBlock for this upkeep
		r.filterStore.UpdateFilters(func(uf1, uf2 upkeepFilter) upkeepFilter {
			uf1.lastRePollBlock = end
			return uf1
		}, f)
	}

	return nil
}

// populatePending adds the logs to the pending list if they are not already pending.
// returns the number of logs added and the number of logs that were already pending.
func (r *logRecoverer) populatePending(f upkeepFilter, filteredLogs []logpoller.Log) (int, int) {
	r.lock.Lock()
	defer r.lock.Unlock()

	pendingSizeBefore := len(r.pending)
	alreadyPending := 0
	for _, log := range filteredLogs {
		trigger := logToTrigger(log)
		upkeepId := &ocr2keepers.UpkeepIdentifier{}
		ok := upkeepId.FromBigInt(f.upkeepID)
		if !ok {
			r.lggr.Warnw("failed to convert upkeepID to UpkeepIdentifier", "upkeepID", f.upkeepID)
			continue
		}
		wid := core.UpkeepWorkID(*upkeepId, trigger)
		if _, ok := r.visited[wid]; ok {
			alreadyPending++
			continue
		}
		checkData, err := r.packer.PackLogData(log)
		if err != nil {
			r.lggr.Warnw("failed to pack log data", "err", err, "log", log)
			continue
		}
		payload, err := core.NewUpkeepPayload(f.upkeepID, trigger, checkData)
		if err != nil {
			r.lggr.Warnw("failed to create payload", "err", err, "log", log)
			continue
		}
		// r.lggr.Debugw("adding a payload to pending", "payload", payload)
		r.visited[wid] = time.Now()
		r.pending = append(r.pending, payload)
	}
	if len(r.pending) == 0 {
		return 0, 0
	}
	return len(r.pending) - pendingSizeBefore, alreadyPending
}

// filterFinalizedStates filters out the log upkeeps that have already been completed (performed or ineligible).
func (r *logRecoverer) filterFinalizedStates(f upkeepFilter, logs []logpoller.Log, states []ocr2keepers.UpkeepState) []logpoller.Log {
	filtered := make([]logpoller.Log, 0)

	for i, log := range logs {
		state := states[i]
		if state != ocr2keepers.UnknownState {
			continue
		}
		filtered = append(filtered, log)
	}

	return filtered
}

// getRecoveryWindow returns the block range of which the recoverer will try work on
func (r *logRecoverer) getRecoveryWindow(latest int64) (int64, int64) {
	lookbackBlocks := r.lookbackBlocks.Load()
	blockTime := r.blockTime.Load()
	start := int64(24*time.Hour) / blockTime
	return latest - start, latest - lookbackBlocks
}

// getFilterBatch returns a batch of filters that are ready to be recovered.
func (r *logRecoverer) getFilterBatch(offsetBlock int64) []upkeepFilter {
	filters := r.filterStore.GetFilters(func(f upkeepFilter) bool {
		if f.lastRePollBlock > offsetBlock {
			return false
		}
		return true
	})

	sort.Slice(filters, func(i, j int) bool {
		return filters[i].lastRePollBlock < filters[j].lastRePollBlock
	})

	return r.selectFilterBatch(filters)
}

// selectFilterBatch selects a batch of filters to be recovered.
// Half of the batch is selected randomly, the other half is selected
// in order of the oldest lastRePollBlock.
func (r *logRecoverer) selectFilterBatch(filters []upkeepFilter) []upkeepFilter {
	batchSize := recoveryBatchSize

	if len(filters) < batchSize {
		return filters
	}
	results := filters[:batchSize/2]
	filters = filters[batchSize/2:]

	for len(results) < batchSize && len(filters) != 0 {
		i := rand.Intn(len(filters))
		results = append(results, filters[i])
		if i == 0 {
			filters = filters[1:]
		} else if i == len(filters)-1 {
			filters = filters[:i]
		} else {
			filters = append(filters[:i], filters[i+1:]...)
		}
	}

	return results
}

func logToTrigger(log logpoller.Log) ocr2keepers.Trigger {
	t := ocr2keepers.NewTrigger(
		// TODO: use zero values or latest block
		ocr2keepers.BlockNumber(log.BlockNumber),
		log.BlockHash,
	)
	t.LogTriggerExtension = &ocr2keepers.LogTriggerExtension{
		TxHash:      log.TxHash,
		Index:       uint32(log.LogIndex),
		BlockHash:   log.BlockHash,
		BlockNumber: ocr2keepers.BlockNumber(log.BlockNumber),
	}
	return t
}
