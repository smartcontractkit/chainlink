package evm

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"golang.org/x/time/rate"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

// TODO: configurable values or based on block time
const (
	// logRetention is the amount of time to retain logs for.
	// 5 minutes is the desired retention time for logs, but we add an extra 10 minutes buffer.
	logRetention            = (time.Minute * 5) + (time.Minute * 10)
	logBlocksLookback int64 = 256
	lookbackBuffer    int64 = 10
)

// TODO: configurable values or based on block time
var (
	blockRateLimit  = rate.Every(time.Second)
	blockLimitBurst = 32
	logsRateLimit   = rate.Every(time.Second)
	logsLimitBurst  = 4
	queryWorkers    = 4
)

// LogTriggerConfig is an alias for log trigger config.
type LogTriggerConfig = i_keeper_registry_master_wrapper_2_1.KeeperRegistryBase21LogTriggerConfig

// upkeepFilterEntry holds the upkeep filter, rate limiter and last polled block.
type upkeepFilterEntry struct {
	id     *big.Int
	filter logpoller.Filter
	cfg    LogTriggerConfig
	// lastPollBlock is the last block number the logs were fetched for this upkeep
	lastPollBlock int64
	// blockLimiter is used to limit the number of blocks to fetch logs for an upkeep
	blockLimiter *rate.Limiter
	// logsLimiter is used to limit the number of logs to fetch for an upkeep
	logsLimiter *rate.Limiter
}

// logEventProvider manages log filters for upkeeps and enables to read the log events.
type logEventProvider struct {
	lggr logger.Logger

	poller logpoller.LogPoller

	lock   *sync.RWMutex
	active map[string]upkeepFilterEntry
}

func NewLogEventProvider(lggr logger.Logger, poller logpoller.LogPoller) *logEventProvider {
	return &logEventProvider{
		lggr:   lggr.Named("LogEventProvider"),
		poller: poller,
		lock:   &sync.RWMutex{},
		active: make(map[string]upkeepFilterEntry),
	}
}

// Register creates a filter from the given upkeep and calls log poller to register it.
func (lfm *logEventProvider) RegisterFilter(upkeepID *big.Int, cfg LogTriggerConfig) error {
	if err := lfm.validateLogTriggerConfig(cfg); err != nil {
		return errors.Wrap(err, "invalid log trigger config")
	}
	filter := lfm.newLogFilter(upkeepID, cfg)

	// TODO: optimize locking, currently we lock the whole map while registering the filter
	lfm.lock.Lock()
	defer lfm.lock.Unlock()

	uid := upkeepID.String()
	if _, ok := lfm.active[uid]; ok {
		// TODO: check for updates
		return errors.Errorf("filter for upkeep with id %s already registered", uid)
	}
	if err := lfm.poller.RegisterFilter(filter); err != nil {
		return errors.Wrap(err, "failed to register upkeep filter")
	}
	lfm.active[uid] = upkeepFilterEntry{
		id:           upkeepID,
		filter:       filter,
		cfg:          cfg,
		blockLimiter: rate.NewLimiter(blockRateLimit, blockLimitBurst),
		logsLimiter:  rate.NewLimiter(logsRateLimit, logsLimitBurst),
	}

	return nil
}

// Unregister removes the filter for the given upkeepID
func (lfm *logEventProvider) UnregisterFilter(upkeepID *big.Int) error {
	err := lfm.poller.UnregisterFilter(lfm.filterName(upkeepID), nil)
	if err == nil {
		lfm.lock.Lock()
		delete(lfm.active, upkeepID.String())
		lfm.lock.Unlock()
	}
	return errors.Wrap(err, "failed to unregister upkeep filter")
}

// GetLogs returns the logs for the given upkeeps, by reading the logs from the last polled block for each upkeep.
func (lfm *logEventProvider) GetLogs(ctx context.Context, ids ...*big.Int) ([][]logpoller.Log, error) {
	latest, err := lfm.poller.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrHeadNotAvailable, err)
	}
	entries := lfm.getEntries(latest, ids...)

	lfm.lggr.Debugw("polling logs for entries", "latestBlock", latest, "entries", len(entries))

	results, err := lfm.getLogsConcurrently(latest, entries)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs: %w", err)
	}
	// update last poll block
	lfm.lock.Lock()
	defer lfm.lock.Unlock()
	for _, entry := range entries {
		// for successful queries, the last poll block was updated
		orig := lfm.active[entry.id.String()]
		if entry.lastPollBlock == orig.lastPollBlock {
			continue
		}
		orig.lastPollBlock = entry.lastPollBlock
		lfm.active[entry.id.String()] = orig
	}

	return results, nil
}

// getFilters returns the filters for the given upkeepIDs,
// returns empty filter for inactive upkeeps.
//
// TODO: group filters by contract address?
func (lfm *logEventProvider) getEntries(latestBlock int64, ids ...*big.Int) []*upkeepFilterEntry {
	lfm.lock.RLock()
	defer lfm.lock.RUnlock()

	var filters []*upkeepFilterEntry
	for _, id := range ids {
		entry, ok := lfm.active[id.String()]
		if !ok { // entry not found, could be inactive upkeep
			lfm.lggr.Debugw("upkeep filter not found", "upkeep", id.String())
			filters = append(filters, &upkeepFilterEntry{id: id})
			continue
		}
		if entry.lastPollBlock > latestBlock {
			lfm.lggr.Debugw("already polled latest block", "entry.lastPollBlock", entry.lastPollBlock, "latestBlock", latestBlock, "upkeep", id.String())
			filters = append(filters, &upkeepFilterEntry{id: id, lastPollBlock: entry.lastPollBlock})
			continue
		}
		// recreating the struct to be thread safe
		filters = append(filters, &upkeepFilterEntry{
			id:            id,
			filter:        lfm.newLogFilter(id, entry.cfg),
			lastPollBlock: entry.lastPollBlock,
			blockLimiter:  entry.blockLimiter,
			logsLimiter:   entry.logsLimiter,
		})
	}

	return filters
}

func (lfm *logEventProvider) getLogsConcurrently(latest int64, entries []*upkeepFilterEntry) ([][]logpoller.Log, error) {
	results := make([][]logpoller.Log, len(entries))

	var wg sync.WaitGroup
	// using a set of worker goroutines to fetch logs for upkeeps
	for i := 0; i < len(entries); i += queryWorkers {
		end := i + queryWorkers
		if end > len(entries) {
			end = len(entries)
		}
		wg.Add(1)
		go func(i, end int) {
			defer wg.Done()
			localResults, err := lfm.getLogs(latest, entries[i:end]...)
			if err != nil {
				lfm.lggr.Debugw("failed to get logs", "i", i, "end", end, "err", err)
				// TODO: accumulate errors
			}
			n := 0
			for j, res := range localResults {
				n += len(res)
				// each worker writes to a different index in the results slice
				// so we don't need to lock
				results[i+j] = res
			}
			if n > 0 {
				lfm.lggr.Debugw("got logs", "i", i, "end", end, "latestBlock", latest, "n", n, "err", err)
			}
		}(i, end)
	}
	wg.Wait()

	return results, nil
}

// getLogs calls log poller to get the logs for the given upkeep entries.
func (lfm *logEventProvider) getLogs(latest int64, entries ...*upkeepFilterEntry) ([][]logpoller.Log, error) {
	mainLggr := lfm.lggr.With("latestBlock", latest)
	var results [][]logpoller.Log
	for _, entry := range entries {
		if len(entry.filter.Addresses) == 0 {
			results = append(results, nil)
			continue
		}
		lggr := mainLggr.With("upkeep", entry.id.String(), "addrs", entry.filter.Addresses, "sigs", entry.filter.EventSigs)
		start := entry.lastPollBlock
		if start == 0 { // first time polling
			start = latest - logBlocksLookback
			entry.blockLimiter.SetBurst(int(logBlocksLookback + 1))
			entry.logsLimiter.SetBurst(logsLimitBurst * 4)
		}
		start = start - lookbackBuffer // adding a buffer to avoid missing logs
		if start < 0 {
			start = 0
		}
		resv := entry.blockLimiter.ReserveN(time.Now(), int(latest-start))
		if !resv.OK() {
			results = append(results, nil)
			lggr.Warnw("log upkeep block limit exceeded")
			continue
		}
		lggr = lggr.With("startBlock", start)
		// TODO: TBD what function to use to get logs
		logs, err := lfm.poller.LogsWithSigs(start, latest, entry.filter.EventSigs, entry.filter.Addresses[0])
		if err != nil {
			resv.Cancel() // cancels limit reservation as we failed to get logs
			lggr.Warnw("failed to get logs", "err", err)
			results = append(results, nil)
			continue
		}
		// if this limiter's burst was set to the max, we need to reset it
		if entry.blockLimiter.Burst() == int(logBlocksLookback+1) {
			entry.blockLimiter.SetBurst(blockLimitBurst)
			entry.logsLimiter.SetBurst(logsLimitBurst)
		}
		filtered := make([]logpoller.Log, 0)
		for _, log := range logs {
			if entry.lastPollBlock > log.BlockNumber {
				// TODO: check if the log is really known and not a result of some reorg
				continue
			}
			filtered = append(filtered, log)
		}
		resv = entry.logsLimiter.ReserveN(time.Now(), len(filtered))
		if !resv.OK() {
			results = append(results, nil)
			lggr.Warnw("log upkeep log limit exceeded")
			continue
		}
		// lggr.Debugw("fetched logs", "n", len(logs), "filtered", len(filtered))
		results = append(results, filtered)
		entry.lastPollBlock = latest
	}

	return results, nil
}

// newLogFilter creates logpoller.Filter from the given upkeep config
func (lfm *logEventProvider) newLogFilter(upkeepID *big.Int, cfg LogTriggerConfig) logpoller.Filter {
	sigs := lfm.getFiltersBySelector(cfg.FilterSelector, cfg.Topic1[:], cfg.Topic2[:], cfg.Topic3[:])
	sigs = append([]common.Hash{common.BytesToHash(cfg.Topic0[:])}, sigs...)
	return logpoller.Filter{
		Name:      lfm.filterName(upkeepID),
		EventSigs: sigs,
		Addresses: []common.Address{cfg.ContractAddress},
		Retention: logRetention,
	}
}

func (lfm *logEventProvider) validateLogTriggerConfig(cfg LogTriggerConfig) error {
	var zeroAddr common.Address
	var zeroBytes [32]byte
	if bytes.Equal(cfg.ContractAddress[:], zeroAddr[:]) {
		return errors.New("invalid contract address: zeroed")
	}
	if bytes.Equal(cfg.Topic0[:], zeroBytes[:]) {
		return errors.New("invalid topic0: zeroed")
	}
	return nil
}

// getFiltersBySelector the filters based on the filterSelector
func (lfm *logEventProvider) getFiltersBySelector(filterSelector uint8, filters ...[]byte) []common.Hash {
	var sigs []common.Hash
	var zeroBytes [32]byte
	for i, f := range filters {
		// bitwise AND the filterSelector with the index to check if the filter is needed
		mask := uint8(1 << uint8(i))
		a := filterSelector & mask
		if a == uint8(0) {
			continue
		}
		if bytes.Equal(f, zeroBytes[:]) {
			continue
		}
		sigs = append(sigs, common.BytesToHash(common.LeftPadBytes(f, 32)))
	}
	return sigs
}

func (lfm *logEventProvider) filterName(upkeepID *big.Int) string {
	return logpoller.FilterName(upkeepID.String())
}
