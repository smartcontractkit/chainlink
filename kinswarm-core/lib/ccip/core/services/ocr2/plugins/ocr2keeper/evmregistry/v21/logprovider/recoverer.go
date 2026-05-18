package logprovider

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"

	"github.com/smartcontractkit/chainlink-automation/pkg/v3/random"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/core"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/prommetrics"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	LogRecovererServiceName = "LogRecoverer"

	// RecoveryInterval is the interval at which the recovery scanning processing is triggered
	RecoveryInterval = 5 * time.Second
	// RecoveryCacheTTL is the time to live for the recovery cache
	RecoveryCacheTTL = 10 * time.Minute
	// GCInterval is the interval at which the recovery cache is cleaned up
	GCInterval = RecoveryCacheTTL - time.Second
	// MaxProposals is the maximum number of proposals that can be returned by GetRecoveryProposals
	MaxProposals = 20
	// recoveryBatchSize is the number of filters to recover in a single batch
	recoveryBatchSize = 10
	// recoveryLogsBuffer is the number of blocks to be used as a safety buffer when reading logs
	recoveryLogsBuffer = int64(200)
	recoveryLogsBurst  = int64(500)
	// blockTimeUpdateCadence is the cadence at which the chain's blocktime is re-calculated
	blockTimeUpdateCadence = 10 * time.Minute
	// maxPendingPayloadsPerUpkeep is the number of logs we can have pending for a single upkeep
	// at any given time
	maxPendingPayloadsPerUpkeep = 500
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
	services.StateMachine
	threadCtrl utils.ThreadControl

	lggr logger.SugaredLogger

	lookbackBlocks *atomic.Int64
	blockTime      *atomic.Int64

	interval time.Duration
	lock     sync.RWMutex

	pending []ocr2keepers.UpkeepPayload
	visited map[string]visitedRecord

	filterStore       UpkeepFilterStore
	states            core.UpkeepStateReader
	packer            LogDataPacker
	poller            logpoller.LogPoller
	client            client.Client
	blockTimeResolver *blockTimeResolver

	finalityDepth int64
}

var _ LogRecoverer = &logRecoverer{}

func NewLogRecoverer(lggr logger.Logger, poller logpoller.LogPoller, client client.Client, stateStore core.UpkeepStateReader, packer LogDataPacker, filterStore UpkeepFilterStore, opts LogTriggersOptions) *logRecoverer {
	rec := &logRecoverer{
		lggr: logger.Sugared(lggr).Named(LogRecovererServiceName),

		threadCtrl: utils.NewThreadControl(),

		blockTime:      new(atomic.Int64),
		lookbackBlocks: new(atomic.Int64),
		interval:       opts.ReadInterval * 5,

		pending:           make([]ocr2keepers.UpkeepPayload, 0),
		visited:           make(map[string]visitedRecord),
		poller:            poller,
		filterStore:       filterStore,
		states:            stateStore,
		packer:            packer,
		client:            client,
		blockTimeResolver: newBlockTimeResolver(poller),

		finalityDepth: opts.FinalityDepth,
	}

	rec.lookbackBlocks.Store(opts.LookbackBlocks)
	rec.blockTime.Store(int64(defaultBlockTime))

	return rec
}

// Start starts the log recoverer, which runs 3 threads in the background:
// 1. Recovery thread: scans for logs that were missed by the log poller
// 2. Cleanup thread: cleans up the cache of logs that were already processed
// 3. Block time thread: updates the block time of the chain
func (r *logRecoverer) Start(ctx context.Context) error {
	return r.StartOnce(LogRecovererServiceName, func() error {
		r.updateBlockTime(ctx)

		r.lggr.Infow("starting log recoverer", "blockTime", r.blockTime.Load(), "lookbackBlocks", r.lookbackBlocks.Load(), "interval", r.interval)

		r.threadCtrl.Go(func(ctx context.Context) {
			recoveryTicker := time.NewTicker(r.interval)
			defer recoveryTicker.Stop()

			for {
				select {
				case <-recoveryTicker.C:
					if err := r.recover(ctx); err != nil {
						r.lggr.Warnw("failed to recover logs", "err", err)
					}
				case <-ctx.Done():
					return
				}
			}
		})

		r.threadCtrl.Go(func(ctx context.Context) {
			cleanupTicker := services.NewTicker(GCInterval)
			defer cleanupTicker.Stop()

			for {
				select {
				case <-cleanupTicker.C:
					r.clean(ctx)
					cleanupTicker.Reset()
				case <-ctx.Done():
					return
				}
			}
		})

		r.threadCtrl.Go(func(ctx context.Context) {
			blockTimeTicker := time.NewTicker(blockTimeUpdateCadence)
			defer blockTimeTicker.Stop()

			for {
				select {
				case <-blockTimeTicker.C:
					r.updateBlockTime(ctx)
					blockTimeTicker.Reset(utils.WithJitter(blockTimeUpdateCadence))
				case <-ctx.Done():
					return
				}
			}
		})

		return nil
	})
}

func (r *logRecoverer) Close() error {
	return r.StopOnce(LogRecovererServiceName, func() error {
		r.threadCtrl.Close()
		return nil
	})
}

func (r *logRecoverer) HealthReport() map[string]error {
	return map[string]error{LogRecovererServiceName: r.Healthy()}
}

func (r *logRecoverer) GetProposalData(ctx context.Context, proposal ocr2keepers.CoordinatedBlockProposal) ([]byte, error) {
	switch core.GetUpkeepType(proposal.UpkeepID) {
	case types.LogTrigger:
		return r.getLogTriggerCheckData(ctx, proposal)
	default:
		return []byte{}, errors.New("not a log trigger upkeep ID")
	}
}

func (r *logRecoverer) getLogTriggerCheckData(ctx context.Context, proposal ocr2keepers.CoordinatedBlockProposal) ([]byte, error) {
	if !r.filterStore.Has(proposal.UpkeepID.BigInt()) {
		return nil, fmt.Errorf("filter not found for upkeep %v", proposal.UpkeepID)
	}
	latest, err := r.poller.LatestBlock(ctx)
	if err != nil {
		return nil, err
	}

	start, offsetBlock := r.getRecoveryWindow(latest.BlockNumber)
	if proposal.Trigger.LogTriggerExtension == nil {
		return nil, errors.New("missing log trigger extension")
	}

	// Verify the log is still present on chain, not reorged and is within recoverable range
	// Do not trust the logBlockNumber from proposal since it's not included in workID
	logBlockHash := common.BytesToHash(proposal.Trigger.LogTriggerExtension.BlockHash[:])
	bn, bh, err := core.GetTxBlock(ctx, r.client, proposal.Trigger.LogTriggerExtension.TxHash)
	if err != nil {
		return nil, err
	}
	if bn == nil {
		return nil, errors.New("failed to get tx block")
	}
	if bh.Hex() != logBlockHash.Hex() {
		return nil, errors.New("log tx reorged")
	}
	logBlock := bn.Int64()
	if isRecoverable := logBlock < offsetBlock && logBlock > start; !isRecoverable {
		return nil, errors.New("log block is not recoverable")
	}

	// Check if the log was already performed or ineligible
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

	logs, err := r.poller.LogsWithSigs(ctx, logBlock-1, logBlock+1, filter.topics, common.BytesToAddress(filter.addr))
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
	latestBlock, err := r.poller.LatestBlock(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrHeadNotAvailable, err)
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	if len(r.pending) == 0 {
		return nil, nil
	}

	allLogsCounter := 0
	logsCount := map[string]int{}

	r.sortPending(uint64(latestBlock.BlockNumber))

	var results, pending []ocr2keepers.UpkeepPayload
	for _, payload := range r.pending {
		if allLogsCounter >= MaxProposals {
			// we have enough proposals, the rest are pushed back to pending
			pending = append(pending, payload)
			continue
		}
		uid := payload.UpkeepID.String()
		if logsCount[uid] >= AllowedLogsPerUpkeep {
			// we have enough proposals for this upkeep, the rest are pushed back to pending
			pending = append(pending, payload)
			continue
		}
		results = append(results, payload)
		logsCount[uid]++
		allLogsCounter++
	}

	r.pending = pending
	prommetrics.AutomationRecovererPendingPayloads.Set(float64(len(r.pending)))

	r.lggr.Debugf("found %d recoverable payloads", len(results))

	return results, nil
}

func (r *logRecoverer) recover(ctx context.Context) error {
	latest, err := r.poller.LatestBlock(ctx)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrHeadNotAvailable, err)
	}

	start, offsetBlock := r.getRecoveryWindow(latest.BlockNumber)
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
	start := f.lastRePollBlock + 1 // NOTE: we expect f.lastRePollBlock + 1 <= offsetBlock, as others would have been filtered out
	// ensure we don't recover logs from before the filter was created
	if configUpdateBlock := int64(f.configUpdateBlock); start < configUpdateBlock {
		// NOTE: we expect that configUpdateBlock <= offsetBlock, as others would have been filtered out
		start = configUpdateBlock
	}
	if start < startBlock {
		start = startBlock
	}
	end := start + recoveryLogsBuffer
	if offsetBlock-end > 100*recoveryLogsBuffer {
		// If recoverer is lagging by a lot (more than 100x recoveryLogsBuffer), allow
		// a range of recoveryLogsBurst
		// Exploratory: Store lastRePollBlock in DB to prevent bursts during restarts
		// (while also taking into account existing pending payloads)
		end = start + recoveryLogsBurst
	}
	if end > offsetBlock {
		end = offsetBlock
	}
	// we expect start to be > offsetBlock in any case
	logs, err := r.poller.LogsWithSigs(ctx, start, end, f.topics, common.BytesToAddress(f.addr))
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

	added, alreadyPending, ok := r.populatePending(f, filteredLogs)
	if added > 0 {
		r.lggr.Debugw("found missed logs", "added", added, "alreadyPending", alreadyPending, "upkeepID", f.upkeepID)
		prommetrics.AutomationRecovererMissedLogs.Add(float64(added))
	}
	if !ok {
		r.lggr.Debugw("failed to add all logs to pending", "upkeepID", f.upkeepID)
		return nil
	}
	r.filterStore.UpdateFilters(func(uf1, uf2 upkeepFilter) upkeepFilter {
		uf1.lastRePollBlock = end
		r.lggr.Debugw("Updated lastRePollBlock", "lastRePollBlock", end, "upkeepID", uf1.upkeepID)
		return uf1
	}, f)

	return nil
}

// populatePending adds the logs to the pending list if they are not already pending.
// returns the number of logs added, the number of logs that were already pending,
// and a flag that indicates whether some errors happened while we are trying to add to pending q.
func (r *logRecoverer) populatePending(f upkeepFilter, filteredLogs []logpoller.Log) (int, int, bool) {
	r.lock.Lock()
	defer r.lock.Unlock()

	pendingSizeBefore := len(r.pending)
	alreadyPending := 0
	errs := make([]error, 0)
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
		if err := r.addPending(payload); err != nil {
			errs = append(errs, err)
		} else {
			r.visited[wid] = visitedRecord{
				visitedAt: time.Now(),
				payload:   payload,
			}
		}
	}
	return len(r.pending) - pendingSizeBefore, alreadyPending, len(errs) == 0
}

// filterFinalizedStates filters out the log upkeeps that have already been completed (performed or ineligible).
func (r *logRecoverer) filterFinalizedStates(_ upkeepFilter, logs []logpoller.Log, states []ocr2keepers.UpkeepState) []logpoller.Log {
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
	start := latest - blocksInDay
	// Exploratory: Instead of subtracting finality depth to account for finalized performs
	// keep two pointers of lastRePollBlock for soft and hard finalization, i.e. manage
	// unfinalized perform logs better
	end := latest - lookbackBlocks - r.finalityDepth
	if start > end {
		// In this case, allow starting from more than a day behind
		start = end
	}
	return start, end
}

// getFilterBatch returns a batch of filters that are ready to be recovered.
func (r *logRecoverer) getFilterBatch(offsetBlock int64) []upkeepFilter {
	filters := r.filterStore.GetFilters(func(f upkeepFilter) bool {
		// ensure we work only on filters that are ready to be recovered
		// no need to recover in case f.configUpdateBlock is after offsetBlock
		return f.lastRePollBlock < offsetBlock && int64(f.configUpdateBlock) <= offsetBlock
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
			r.lggr.Debugw("error generating random number", "err", err.Error())
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
	err := r.tryExpire(ctx, expired...)
	if err != nil {
		lggr.Warnw("failed to clean visited upkeeps", "err", err)
	}
}

func (r *logRecoverer) tryExpire(ctx context.Context, ids ...string) error {
	latestBlock, err := r.poller.LatestBlock(ctx)
	if err != nil {
		return fmt.Errorf("failed to get latest block: %w", err)
	}
	sort.Slice(ids, func(i, j int) bool {
		return ids[i] < ids[j]
	})
	states, err := r.states.SelectByWorkIDs(ctx, ids...)
	if err != nil {
		return fmt.Errorf("failed to get states: %w", err)
	}
	lggr := r.lggr.With("where", "clean")
	start, _ := r.getRecoveryWindow(latestBlock.BlockNumber)
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
					"latestBlock", latestBlock, "logBlock", logBlock, "start", start)
				r.removePending(rec.payload.WorkID)
				delete(r.visited, ids[i])
				removed++
				continue
			}
			if err := r.addPending(rec.payload); err == nil {
				rec.visitedAt = time.Now()
				r.visited[ids[i]] = rec
			}
		default:
			delete(r.visited, ids[i])
			removed++
		}
	}

	if removed > 0 {
		lggr.Debugw("expired upkeeps", "expired", len(ids), "cleaned", removed)
	}

	return nil
}

// addPending adds a payload to the pending list if it's not already there.
// NOTE: the lock must be held before calling this function.
func (r *logRecoverer) addPending(payload ocr2keepers.UpkeepPayload) error {
	var exist bool
	pending := r.pending
	upkeepPayloads := 0
	for _, p := range pending {
		if bytes.Equal(p.UpkeepID[:], payload.UpkeepID[:]) {
			upkeepPayloads++
		}
		if p.WorkID == payload.WorkID {
			exist = true
		}
	}
	if upkeepPayloads >= maxPendingPayloadsPerUpkeep {
		return fmt.Errorf("upkeep %v has too many payloads in pending queue", payload.UpkeepID)
	}
	if !exist {
		r.pending = append(pending, payload)
		prommetrics.AutomationRecovererPendingPayloads.Inc()
	}
	return nil
}

// removePending removes a payload from the pending list.
// NOTE: the lock must be held before calling this function.
func (r *logRecoverer) removePending(workID string) {
	updated := make([]ocr2keepers.UpkeepPayload, 0, len(r.pending))
	for _, p := range r.pending {
		if p.WorkID != workID {
			updated = append(updated, p)
		} else {
			prommetrics.AutomationRecovererPendingPayloads.Dec()
		}
	}
	r.pending = updated
}

// sortPending sorts the pending list by a random order based on the normalized latest block number.
// Divided by 10 to ensure that nodes with similar block numbers won't end up with different order.
// NOTE: the lock must be held before calling this function.
func (r *logRecoverer) sortPending(latestBlock uint64) {
	normalized := latestBlock / 100
	if normalized == 0 {
		normalized = 1
	}
	randSeed := random.GetRandomKeySource(nil, normalized)

	shuffledIDs := make(map[string]string, len(r.pending))
	for _, p := range r.pending {
		shuffledIDs[p.WorkID] = random.ShuffleString(p.WorkID, randSeed)
	}

	sort.SliceStable(r.pending, func(i, j int) bool {
		return shuffledIDs[r.pending[i].WorkID] < shuffledIDs[r.pending[j].WorkID]
	})
}

func (r *logRecoverer) updateBlockTime(ctx context.Context) {
	blockTime, err := r.blockTimeResolver.BlockTime(ctx, defaultSampleSize)
	if err != nil {
		r.lggr.Warnw("failed to compute block time", "err", err)
		return
	}
	if blockTime > 0 {
		currentBlockTime := r.blockTime.Load()
		newBlockTime := int64(blockTime)
		if currentBlockTime > 0 && (int64(math.Abs(float64(currentBlockTime-newBlockTime)))*100/currentBlockTime) > 20 {
			r.lggr.Warnf("updating blocktime from %d to %d, this change is larger than 20%%", currentBlockTime, newBlockTime)
		} else {
			r.lggr.Debugf("updating blocktime from %d to %d", currentBlockTime, newBlockTime)
		}
		r.blockTime.Store(newBlockTime)
	}
}
