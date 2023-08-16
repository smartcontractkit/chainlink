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
	RecoveryBatchSize       = 10
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

	lookbackBlocks int64
	interval       time.Duration
	lock           sync.RWMutex

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

	return &logRecoverer{
		lggr: lggr.Named("LogRecoverer"),

		lookbackBlocks: lookbackBlocks,
		interval:       interval,

		pending:     make([]ocr2keepers.UpkeepPayload, 0),
		visited:     make(map[string]time.Time),
		poller:      poller,
		filterStore: filterStore,
		states:      stateStore,
		packer:      packer,
	}
}

func (r *logRecoverer) Start(pctx context.Context) error {
	ctx, cancel := context.WithCancel(pctx)
	r.lock.Lock()
	r.cancel = cancel
	interval := r.interval
	r.lock.Unlock()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	r.lggr.Debug("Starting log recoverer")

	for {
		select {
		case <-ticker.C:
			r.recover(ctx)
		case <-ctx.Done():
			return nil
		}
	}
}

func (r *logRecoverer) Close() error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.cancel != nil {
		r.cancel()
	}
	return nil
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

	pending := r.pending
	r.pending = make([]ocr2keepers.UpkeepPayload, 0)

	for _, p := range pending {
		if _, ok := r.visited[p.WorkID]; !ok {
			r.visited[p.WorkID] = time.Now()
		}
	}

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
	offsetBlock := r.getRecoveryOffsetBlock(latest)
	if offsetBlock < 0 {
		// too soon to recover, we don't have enough blocks
		return nil
	}

	filters := r.getFiltersBatch(offsetBlock)
	if len(filters) == 0 {
		return nil
	}

	r.lggr.Debugw("recovering logs", "filters", filters)

	var wg sync.WaitGroup
	for _, f := range filters {
		wg.Add(1)
		go func(f upkeepFilter) {
			defer wg.Done()
			r.recoverFilter(ctx, f)
		}(f)
	}
	wg.Wait()

	return nil
}

func (r *logRecoverer) recoverFilter(ctx context.Context, f upkeepFilter) error {
	start, end := f.lastRePollBlock, f.lastRePollBlock+10
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

	filteredLogs := make([]logpoller.Log, 0)
	for i, log := range logs {
		state := states[i]
		if state != ocr2keepers.UnknownState {
			continue
		}
		filteredLogs = append(filteredLogs, log)
	}

	for _, log := range filteredLogs {
		trigger := logToTrigger(log)
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
		r.pending = append(r.pending, payload)
	}

	return nil
}

// getRecoveryOffsetBlock returns the max block number that the recoverer will try to fetch logs.
func (r *logRecoverer) getRecoveryOffsetBlock(latest int64) int64 {
	lookbackBlocks := atomic.LoadInt64(&r.lookbackBlocks)
	return latest - lookbackBlocks
}

func (r *logRecoverer) getFiltersBatch(offsetBlock int64) []upkeepFilter {
	filters := r.filterStore.GetFilters(func(f upkeepFilter) bool {
		if f.lastRePollBlock >= offsetBlock {
			return false
		}
		return true
	})

	sort.Slice(filters, func(i, j int) bool {
		return filters[i].lastRePollBlock < filters[j].lastRePollBlock
	})

	return r.selectBatch(filters)
}

func (r *logRecoverer) selectBatch(filters []upkeepFilter) []upkeepFilter {
	batchSize := RecoveryBatchSize

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
