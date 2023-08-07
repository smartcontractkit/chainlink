package logprovider

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
	keepersflows "github.com/smartcontractkit/ocr2keepers/pkg/v3/flows"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

const (
	performedBuffer = 32
)

type UpkeepStateReader interface {
	// SelectByUpkeepIDsAndBlockRange retrieves upkeep states for provided upkeep ids and block range, the result is currently not in particular order
	SelectByUpkeepIDsAndBlockRange(upkeepIds []*big.Int, start, end int64) ([]*ocr2keepers.UpkeepPayload, []*ocr2keepers.UpkeepState, error)
}

type RecoveryOptions struct {
	Interval       time.Duration
	GCInterval     time.Duration
	TTL            time.Duration
	LookbackBlocks int64
	BatchSize      int32
}

func (o *RecoveryOptions) defaults() {
	if o.Interval == 0 {
		o.Interval = 30 * time.Second
	}
	if o.GCInterval == 0 {
		o.GCInterval = o.Interval*10 + o.Interval/2
	}
	if o.TTL == 0 {
		o.TTL = 24*time.Hour - time.Second
	}
	if o.LookbackBlocks == 0 {
		o.LookbackBlocks = 512
	}
	if o.BatchSize == 0 {
		o.BatchSize = 10
	}
}

type logRecoverer struct {
	lggr logger.Logger

	cancel context.CancelFunc

	lock *sync.RWMutex

	opts RecoveryOptions

	pending []ocr2keepers.UpkeepPayload

	filterStore UpkeepFilterStore
	visited     map[string]time.Time

	poller          logpoller.LogPoller
	upkeepStates    UpkeepStateReader
	registryAddress common.Address
}

var _ keepersflows.RecoverableProvider = &logRecoverer{}

func NewLogRecoverer(
	lggr logger.Logger,
	poller logpoller.LogPoller,
	upkeepStates UpkeepStateReader,
	registryAddress common.Address,
	filterStore UpkeepFilterStore,
	opts *RecoveryOptions,
) (*logRecoverer, error) {
	if opts == nil {
		opts = new(RecoveryOptions)
	}
	opts.defaults()

	return &logRecoverer{
		lggr:            lggr.Named("LogRecoverer"),
		opts:            *opts,
		lock:            &sync.RWMutex{},
		pending:         make([]ocr2keepers.UpkeepPayload, 0),
		filterStore:     filterStore,
		visited:         make(map[string]time.Time),
		poller:          poller,
		registryAddress: registryAddress,
	}, nil
}

func recoveryFilterName(addr common.Address) string {
	return logpoller.FilterName("KeepersRegistry LogRecoverer", addr)
}

func (r *logRecoverer) Start(pctx context.Context) error {
	ctx, cancel := context.WithCancel(pctx)
	r.lock.Lock()
	r.cancel = cancel
	interval := r.opts.Interval
	gcInterval := r.opts.GCInterval
	r.lock.Unlock()

	if err := r.registerFilters(ctx); err != nil {
		return err
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	gcTicker := time.NewTicker(gcInterval)
	defer gcTicker.Stop()

	r.lggr.Debug("Starting log recoverer")

	for {
		select {
		case <-ticker.C:
			r.recover(ctx)
		case <-gcTicker.C:
			r.clean(ctx)
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

func (r *logRecoverer) GetRecoverables() ([]ocr2keepers.UpkeepPayload, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if len(r.pending) == 0 {
		return nil, nil
	}

	pending := r.pending
	r.pending = make([]ocr2keepers.UpkeepPayload, 0)

	for _, p := range pending {
		r.visited[p.ID] = time.Now()
	}

	return pending, nil
}

func (r *logRecoverer) registerFilters(ctx context.Context) error {
	return r.poller.RegisterFilter(logpoller.Filter{
		Name: recoveryFilterName(r.registryAddress),
		EventSigs: []common.Hash{
			// listening to dedup key added event
			iregistry21.IKeeperRegistryMasterDedupKeyAdded{}.Topic(),
		},
		Addresses: []common.Address{r.registryAddress},
	})
}

func (r *logRecoverer) clean(ctx context.Context) {
	r.lock.Lock()
	defer r.lock.Unlock()

	cleaned := 0
	for id, t := range r.visited {
		if time.Since(t) > r.opts.TTL {
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
	logs, err = r.filterPerformed(ctx, start, end, f, logs)
	if err != nil {
		return fmt.Errorf("failed to filter performed: %w", err)
	}
	logs, err = r.filterIneligible(start, end, f, logs)
	if err != nil {
		return fmt.Errorf("failed to filter ineligible: %w", err)
	}

	for _, log := range logs {
		trigger := logToTrigger(f.upkeepID, log)
		// TODO: align payload creation
		r.pending = append(r.pending, ocr2keepers.UpkeepPayload{
			Upkeep: ocr2keepers.ConfiguredUpkeep{
				ID:   f.upkeepID.Bytes(),
				Type: 1,
			},
			Trigger: trigger,
		})
	}

	return nil
}

func UpkeepWorkID(t ocr2keepers.Trigger) string {
	return fmt.Sprintf("%+v", t)
}

func (r *logRecoverer) filterPerformed(ctx context.Context, start, end int64, f upkeepFilter, logs []logpoller.Log) ([]logpoller.Log, error) {
	workIDs := make(map[string]logpoller.Log)
	for _, log := range logs {
		trigger := logToTrigger(f.upkeepID, log)
		workIDs[UpkeepWorkID(trigger)] = log
	}
	performedLogs, err := r.poller.LogsWithSigs(
		start,
		end+performedBuffer,
		[]common.Hash{
			iregistry21.IKeeperRegistryMasterDedupKeyAdded{}.Topic(),
		},
		r.registryAddress,
		pg.WithParentCtx(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to collect DedupKeyAdded (performed) logs from log poller", err)
	}
	for _, log := range performedLogs {
		topics := log.GetTopics()
		if len(topics) < 2 {
			r.lggr.Debugw("unexpected log topics", "topics", topics)
		}
		key := hexutil.Encode(topics[1].Bytes())
		if _, ok := workIDs[key]; ok {
			delete(workIDs, key)
		}
	}
	logs = make([]logpoller.Log, 0, len(workIDs))
	for _, log := range workIDs {
		logs = append(logs, log)
	}

	return logs, nil
}

func (r *logRecoverer) filterIneligible(start, end int64, f upkeepFilter, logs []logpoller.Log) ([]logpoller.Log, error) {
	payloads := make([]ocr2keepers.UpkeepPayload, 0)
	// TODO uncomment when upkeep states store is ready
	// payloads, _, err := r.upkeepStates.SelectByWorkID()
	// if err != nil {
	// 	return nil, fmt.Errorf("could not read upkeep states: %w", err)
	// }

	results := make([]logpoller.Log, 0)
	for _, log := range logs {
		trigger := logToTrigger(f.upkeepID, log)
		for _, payload := range payloads {
			// TODO: payload.WorkID
			if ts, ok := r.visited[payload.ID]; ok && time.Since(ts) < r.opts.TTL {
				continue // we already visited this log
			}
			r.lggr.Debugw("recovered missed log", "trigger", trigger, "upkeepID", f.upkeepID)
			results = append(results, log)
			break
		}
	}

	return results, nil
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
	batchSize := int(atomic.LoadInt32(&r.opts.BatchSize))

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

// getRecoveryOffsetBlock returns the max block number that the recoverer will try to fetch logs.
func (r *logRecoverer) getRecoveryOffsetBlock(latest int64) int64 {
	lookbackBlocks := atomic.LoadInt64(&r.opts.LookbackBlocks)
	return latest - lookbackBlocks
}

func logToTrigger(id *big.Int, log logpoller.Log) ocr2keepers.Trigger {
	return ocr2keepers.NewTrigger(
		log.BlockNumber,
		log.BlockHash.Hex(),
		LogTriggerExtension{
			TxHash:   log.TxHash.Hex(),
			LogIndex: log.LogIndex,
		},
	)
}
