package logprovider

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"math/big"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	DefaultRecoveryInterval = 5 * time.Second
	RecoveryCacheTTL        = 10*time.Minute - time.Second
	GCInterval              = RecoveryCacheTTL

	recoveryBatchSize  = 10
	recoveryLogsBuffer = int64(50)
)

type LogRecoverer interface {
	ocr2keepers.RecoverableProvider
	GetProposalData(context.Context, ocr2keepers.CoordinatedBlockProposal) ([]byte, error)

	Start(context.Context) error
	io.Closer
}

type visitedRecord struct {
	visitedAt time.Time
	payload   ocr2keepers.UpkeepPayload
}

type logRecoverer struct {
	lggr logger.Logger

	cancel context.CancelFunc

	lookbackBlocks *atomic.Int64
	blockTime      *atomic.Int64

	interval time.Duration
	lock     sync.RWMutex

	pending []ocr2keepers.UpkeepPayload
	visited map[string]visitedRecord

	filterStore UpkeepFilterStore
	states      core.UpkeepStateReader
	packer      LogDataPacker
	poller      logpoller.LogPoller
	client      client.Client
}

var _ LogRecoverer = &logRecoverer{}

func NewLogRecoverer(lggr logger.Logger, poller logpoller.LogPoller, client client.Client, stateStore core.UpkeepStateReader, packer LogDataPacker, filterStore UpkeepFilterStore, opts LogTriggersOptions) *logRecoverer {
	rec := &logRecoverer{
		lggr: lggr.Named("LogRecoverer"),

		blockTime:      &atomic.Int64{},
		lookbackBlocks: &atomic.Int64{},
		interval:       opts.ReadInterval * 5,

		pending:     make([]ocr2keepers.UpkeepPayload, 0),
		visited:     make(map[string]visitedRecord),
		poller:      poller,
		filterStore: filterStore,
		states:      stateStore,
		packer:      packer,
		client:      client,
	}

	rec.lookbackBlocks.Store(opts.LookbackBlocks)
	rec.blockTime.Store(int64(defaultBlockTime))

	return rec
}

func (r *logRecoverer) Start(pctx context.Context) error {
	ctx, cancel := context.WithCancel(context.Background())

	r.lock.Lock()
	if r.cancel != nil {
		r.lock.Unlock()
		cancel() // Cancel the created context
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
			gcTicker := time.NewTicker(utils.WithJitter(GCInterval))
			defer gcTicker.Stop()

			for {
				select {
				case <-ticker.C:
					if err := r.recover(ctx); err != nil {
						r.lggr.Warnw("failed to recover logs", "err", err)
					}
				case <-gcTicker.C:
					r.clean(ctx)
					gcTicker.Reset(utils.WithJitter(GCInterval))
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

func (r *logRecoverer) GetProposalData(ctx context.Context, proposal ocr2keepers.CoordinatedBlockProposal) ([]byte, error) {
	switch core.GetUpkeepType(proposal.UpkeepID) {
	case ocr2keepers.LogTrigger:
		return r.getLogTriggerCheckData(ctx, proposal)
	default:
		return []byte{}, errors.New("not a log trigger upkeep ID")
	}
}

func (r *logRecoverer) getLogTriggerCheckData(ctx context.Context, proposal ocr2keepers.CoordinatedBlockProposal) ([]byte, error) {
	if !r.filterStore.Has(proposal.UpkeepID.BigInt()) {
		return nil, fmt.Errorf("filter not found for upkeep %v", proposal.UpkeepID)
	}
	latest, err := r.poller.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		return nil, err
	}

	start, offsetBlock := r.getRecoveryWindow(latest)
	if proposal.Trigger.LogTriggerExtension == nil {
		return nil, errors.New("missing log trigger extension")
	}
	logBlock := int64(proposal.Trigger.LogTriggerExtension.BlockNumber)
	if logBlock == 0 {
		var number *big.Int
		number, _, err = core.GetTxBlock(ctx, r.client, proposal.Trigger.LogTriggerExtension.TxHash)
		if err != nil {
			return nil, err
		}
		if number == nil {
			return nil, errors.New("failed to get tx block")
		}
		logBlock = number.Int64()
	}
	if isRecoverable := logBlock < offsetBlock && logBlock > start; !isRecoverable {
		return nil, errors.New("log block is not recoverable")
	}
	upkeepStates, err := r.states.SelectByWorkIDs(ctx, proposal.WorkID)
	if err != nil {
		return nil, err
	}

	for _, upkeepState := range upkeepStates {
		switch upkeepState {
		case ocr2keepers.Performed, ocr2keepers.Ineligible:
			return nil, errors.New("upkeep state is not recoverable")
		default:
			// we can proceed
		}
	}

	var filter upkeepFilter
	r.filterStore.RangeFiltersByIDs(func(i int, f upkeepFilter) {
		filter = f
	}, proposal.UpkeepID.BigInt())

	if len(filter.addr) == 0 {
		return nil, fmt.Errorf("invalid filter found for upkeepID %s", proposal.UpkeepID.String())
	}
	if filter.configUpdateBlock > uint64(logBlock) {
		return nil, fmt.Errorf("log block %d is before the filter configUpdateBlock %d for upkeepID %s", logBlock, filter.configUpdateBlock, proposal.UpkeepID.String())
	}

	logs, err := r.poller.LogsWithSigs(logBlock-1, logBlock+1, filter.topics, common.BytesToAddress(filter.addr), pg.WithParentCtx(ctx))
	if err != nil {
		return nil, fmt.Errorf("could not read logs: %w", err)
	}
	logs = filter.Select(logs...)

	for _, log := range logs {
		trigger := logToTrigger(log)
		// use coordinated proposal block number as checkblock/hash
		trigger.BlockHash = proposal.Trigger.BlockHash
		trigger.BlockNumber = proposal.Trigger.BlockNumber
		wid := core.UpkeepWorkID(proposal.UpkeepID, trigger)
		if wid == proposal.WorkID {
			r.lggr.Debugw("found log for proposal", "upkeepId", proposal.UpkeepID, "trigger.ext", trigger.LogTriggerExtension)
			checkData, err := r.packer.PackLogData(log)
			if err != nil {
				return nil, fmt.Errorf("failed to pack log data: %w", err)
			}
			return checkData, nil
		}
	}
	return nil, fmt.Errorf("no log found for upkeepID %v and trigger %+v", proposal.UpkeepID, proposal.Trigger)
}

func (r *logRecoverer) GetRecoveryProposals(ctx context.Context) ([]ocr2keepers.UpkeepPayload, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if len(r.pending) == 0 {
		return nil, nil
	}

	logsCount := map[string]int{}

	var results, pending []ocr2keepers.UpkeepPayload
	for _, payload := range r.pending {
		uid := payload.UpkeepID.String()
		if logsCount[uid] >= AllowedLogsPerUpkeep {
			pending = append(pending, payload)
			continue
		}
		logsCount[uid]++
		results = append(results, payload)
	}
	r.pending = pending

	r.lggr.Debugf("found %d pending payloads", len(pending))

	return results, nil
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
	if start < 0 {
		start = 0
	}

	filters := r.getFilterBatch(offsetBlock)
	if len(filters) == 0 {
		return nil
	}

	r.lggr.Debugw("recovering logs", "filters", filters, "startBlock", start, "offsetBlock", offsetBlock, "latestBlock", latest)

	var wg sync.WaitGroup
	for _, f := range filters {
		wg.Add(1)
		go func(f upkeepFilter) {
			defer wg.Done()
			if err := r.recoverFilter(ctx, f, start, offsetBlock); err != nil {
				r.lggr.Debugw("error recovering filter", "err", err.Error())
			}
		}(f)
	}
	wg.Wait()

	return nil
}

// recoverFilter recovers logs for a single upkeep filter.
func (r *logRecoverer) recoverFilter(ctx context.Context, f upkeepFilter, startBlock, offsetBlock int64) error {
	start := f.lastRePollBlock
	// ensure we don't recover logs from before the filter was created
	// NOTE: we expect that filter with configUpdateBlock > offsetBlock were already filtered out.
	if configUpdateBlock := int64(f.configUpdateBlock); start < configUpdateBlock {
		start = configUpdateBlock
	}
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
	logs = f.Select(logs...)

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

	states, err := r.states.SelectByWorkIDs(ctx, workIDs...)
	if err != nil {
		return fmt.Errorf("could not read states: %w", err)
	}
	if len(logs) != len(states) {
		return fmt.Errorf("log and state count mismatch: %d != %d", len(logs), len(states))
	}
	filteredLogs := r.filterFinalizedStates(f, logs, states)

	added, alreadyPending := r.populatePending(f, filteredLogs)
	if added > 0 {
		r.lggr.Debugw("found missed logs", "added", added, "alreadyPending", alreadyPending, "upkeepID", f.upkeepID)
	}
	r.filterStore.UpdateFilters(func(uf1, uf2 upkeepFilter) upkeepFilter {
		uf1.lastRePollBlock = end
		return uf1
	}, f)

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
		// Set the checkBlock and Hash to zero so that the checkPipeline uses the latest block
		trigger.BlockHash = [32]byte{}
		trigger.BlockNumber = 0
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
		r.visited[wid] = visitedRecord{
			visitedAt: time.Now(),
			payload:   payload,
		}
		r.pending = append(r.pending, payload)
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
	blocksInDay := int64(24*time.Hour) / blockTime
	return latest - blocksInDay, latest - lookbackBlocks
}

// getFilterBatch returns a batch of filters that are ready to be recovered.
func (r *logRecoverer) getFilterBatch(offsetBlock int64) []upkeepFilter {
	filters := r.filterStore.GetFilters(func(f upkeepFilter) bool {
		// ensure we work only on filters that are ready to be recovered
		// no need to recover in case f.configUpdateBlock is after offsetBlock
		return f.lastRePollBlock <= offsetBlock && int64(f.configUpdateBlock) <= offsetBlock
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
		i, err := r.randIntn(len(filters))
		if err != nil {
			r.lggr.Debugw("error generating random number", "error", err.Error())
			continue
		}
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

func (r *logRecoverer) randIntn(limit int) (int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(limit)))
	if err != nil {
		return 0, err
	}

	return int(n.Int64()), nil
}

func logToTrigger(log logpoller.Log) ocr2keepers.Trigger {
	t := ocr2keepers.NewTrigger(
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

func (r *logRecoverer) clean(ctx context.Context) {
	r.lock.RLock()
	var expired []string
	for id, t := range r.visited {
		if time.Since(t.visitedAt) > RecoveryCacheTTL {
			expired = append(expired, id)
		}
	}
	r.lock.RUnlock()
	lggr := r.lggr.With("where", "clean")
	if len(expired) == 0 {
		lggr.Debug("no expired upkeeps")
		return
	}
	cleaned, err := r.tryExpire(ctx, expired...)
	if err != nil {
		lggr.Warnw("failed to clean visited upkeeps", "err", err)
	}
	if len(expired) > 0 {
		lggr.Debugw("expired upkeeps", "expired", len(expired), "cleaned", cleaned)
	}
}

func (r *logRecoverer) tryExpire(ctx context.Context, ids ...string) (int, error) {
	latestBlock, err := r.poller.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		return 0, fmt.Errorf("failed to get latest block: %w", err)
	}
	sort.Slice(ids, func(i, j int) bool {
		return ids[i] < ids[j]
	})
	states, err := r.states.SelectByWorkIDs(ctx, ids...)
	if err != nil {
		return 0, fmt.Errorf("failed to get states: %w", err)
	}
	lggr := r.lggr.With("where", "clean")
	start, _ := r.getRecoveryWindow(latestBlock)
	r.lock.Lock()
	defer r.lock.Unlock()
	var removed int
	for i, state := range states {
		switch state {
		case ocr2keepers.UnknownState:
			// in case the state is unknown, we can't be sure if the upkeep was performed or not
			// so we push it back to the pending list
			rec, ok := r.visited[ids[i]]
			if !ok {
				// in case it was removed by another thread
				continue
			}
			if logBlock := rec.payload.Trigger.LogTriggerExtension.BlockNumber; int64(logBlock) < start {
				// we can't recover this log anymore, so we remove it from the visited list
				lggr.Debugw("removing expired log: old block", "upkeepID", rec.payload.UpkeepID,
					"logBlock", logBlock, "start", start)
				delete(r.visited, ids[i])
				removed++
				continue
			}
			r.pending = append(r.pending, rec.payload)
			rec.visitedAt = time.Now()
			r.visited[ids[i]] = rec
		default:
			delete(r.visited, ids[i])
			removed++
		}
	}

	return removed, nil
}
