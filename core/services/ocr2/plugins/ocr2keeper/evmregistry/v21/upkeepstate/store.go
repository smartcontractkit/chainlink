package upkeepstate

import (
	"context"
	"fmt"
	"io"
	"math/big"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/core"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	UpkeepStateStoreServiceName = "UpkeepStateStore"
	// CacheExpiration is the amount of time that we keep a record in the cache.
	CacheExpiration = 24 * time.Hour
	// GCInterval is the amount of time between cache cleanups.
	GCInterval = 2 * time.Hour
	// flushCadence is the amount of time between flushes to the DB.
	flushCadence         = 30 * time.Second
	concurrentBatchCalls = 10
)

type ORM interface {
	BatchInsertRecords(context.Context, []persistedStateRecord) error
	SelectStatesByWorkIDs(context.Context, []string) ([]persistedStateRecord, error)
	DeleteExpired(context.Context, time.Time) error
}

// UpkeepStateStore is the interface for managing upkeeps final state in a local store.
type UpkeepStateStore interface {
	ocr2keepers.UpkeepStateUpdater
	core.UpkeepStateReader
	Start(context.Context) error
	io.Closer
}

var (
	_           UpkeepStateStore = &upkeepStateStore{}
	newTickerFn                  = time.NewTicker
	batchSize                    = 1000
)

// upkeepStateRecord is a record that we save in a local cache.
type upkeepStateRecord struct {
	workID string
	state  ocr2keepers.UpkeepState

	addedAt time.Time
}

// upkeepStateStore implements UpkeepStateStore.
// It stores the state of ineligible upkeeps in a local, in-memory cache.
// In addition, performed events are fetched by the scanner on demand.
type upkeepStateStore struct {
	services.StateMachine
	threadCtrl utils.ThreadControl

	orm     ORM
	lggr    logger.Logger
	scanner PerformedLogsScanner

	retention    time.Duration
	cleanCadence time.Duration

	mu    sync.RWMutex
	cache map[string]*upkeepStateRecord

	pendingRecords []persistedStateRecord
	sem            chan struct{}
	batchSize      int
}

// NewUpkeepStateStore creates a new state store
func NewUpkeepStateStore(orm ORM, lggr logger.Logger, scanner PerformedLogsScanner) *upkeepStateStore {
	return &upkeepStateStore{
		orm:            orm,
		lggr:           lggr.Named(UpkeepStateStoreServiceName),
		cache:          map[string]*upkeepStateRecord{},
		scanner:        scanner,
		retention:      CacheExpiration,
		cleanCadence:   GCInterval,
		threadCtrl:     utils.NewThreadControl(),
		pendingRecords: []persistedStateRecord{},
		sem:            make(chan struct{}, concurrentBatchCalls),
		batchSize:      batchSize,
	}
}

// Start starts the upkeep state store.
// it does background cleanup of the cache every GCInterval,
// and flush records to DB every flushCadence.
func (u *upkeepStateStore) Start(pctx context.Context) error {
	return u.StartOnce(UpkeepStateStoreServiceName, func() error {
		if err := u.scanner.Start(pctx); err != nil {
			return fmt.Errorf("failed to start scanner")
		}

		u.lggr.Debug("Starting upkeep state store")

		u.threadCtrl.Go(func(ctx context.Context) {
			ticker := time.NewTicker(utils.WithJitter(u.cleanCadence))
			defer ticker.Stop()

			flushTicker := newTickerFn(utils.WithJitter(flushCadence))
			defer flushTicker.Stop()

			for {
				select {
				case <-ticker.C:
					if err := u.cleanup(ctx); err != nil {
						u.lggr.Errorw("unable to clean old state values", "err", err)
					}
					ticker.Reset(utils.WithJitter(u.cleanCadence))
				case <-flushTicker.C:
					u.flush(ctx)
					flushTicker.Reset(utils.WithJitter(flushCadence))
				case <-ctx.Done():
					u.flush(ctx)
					return
				}
			}
		})
		return nil
	})
}

func (u *upkeepStateStore) flush(ctx context.Context) {
	u.mu.Lock()
	cloneRecords := make([]persistedStateRecord, len(u.pendingRecords))
	copy(cloneRecords, u.pendingRecords)
	u.pendingRecords = []persistedStateRecord{}
	u.mu.Unlock()

	for i := 0; i < len(cloneRecords); i += u.batchSize {
		end := i + u.batchSize
		if end > len(cloneRecords) {
			end = len(cloneRecords)
		}

		batch := cloneRecords[i:end]

		u.sem <- struct{}{}

		go func() {
			if err := u.orm.BatchInsertRecords(ctx, batch); err != nil {
				u.lggr.Errorw("error inserting records", "err", err)
			}
			<-u.sem
		}()
	}
}

// Close stops the service of pruning stale data; implements io.Closer
func (u *upkeepStateStore) Close() error {
	return u.StopOnce(UpkeepStateStoreServiceName, func() error {
		u.threadCtrl.Close()
		return nil
	})
}

func (u *upkeepStateStore) HealthReport() map[string]error {
	return map[string]error{UpkeepStateStoreServiceName: u.Healthy()}
}

// SelectByWorkIDs returns the current state of the upkeep for the provided ids.
// If an id is not found, the state is returned as StateUnknown.
// We first check the cache, and if any ids are missing, we fetch them from the scanner and DB.
func (u *upkeepStateStore) SelectByWorkIDs(ctx context.Context, workIDs ...string) ([]ocr2keepers.UpkeepState, error) {
	states, missing := u.selectFromCache(workIDs...)
	if len(missing) == 0 {
		// all ids were found in the cache
		return states, nil
	}
	if err := u.fetchPerformed(ctx, missing...); err != nil {
		return nil, err
	}
	if err := u.fetchFromDB(ctx, missing...); err != nil {
		return nil, err
	}

	// at this point all values should be in the cache. if values are missing
	// their state is indicated as unknown
	states, _ = u.selectFromCache(workIDs...)

	return states, nil
}

// SetUpkeepState updates the state of the upkeep.
// Currently we only store the state if the upkeep is ineligible.
// Performed events will be fetched on demand.
func (u *upkeepStateStore) SetUpkeepState(ctx context.Context, result ocr2keepers.CheckResult, _ ocr2keepers.UpkeepState) error {
	if result.Eligible {
		return nil
	}

	return u.upsertStateRecord(ctx, result.WorkID, ocr2keepers.Ineligible, uint64(result.Trigger.BlockNumber), result.UpkeepID.BigInt(), result.IneligibilityReason)
}

// upsertStateRecord inserts or updates a record for the provided
// check result. If an item already exists in the data store, the state and
// block are updated.
func (u *upkeepStateStore) upsertStateRecord(ctx context.Context, workID string, s ocr2keepers.UpkeepState, b uint64, upkeepID *big.Int, reason uint8) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	record, ok := u.cache[workID]
	if !ok {
		record = &upkeepStateRecord{
			workID:  workID,
			addedAt: time.Now(),
		}
	}

	record.state = s

	u.cache[workID] = record

	u.pendingRecords = append(u.pendingRecords, persistedStateRecord{
		UpkeepID:            ubig.New(upkeepID),
		WorkID:              record.workID,
		CompletionState:     uint8(record.state),
		IneligibilityReason: reason,
		InsertedAt:          record.addedAt,
	})

	return nil
}

// fetchPerformed fetches all performed logs from the scanner to populate the cache.
func (u *upkeepStateStore) fetchPerformed(ctx context.Context, workIDs ...string) error {
	performed, err := u.scanner.ScanWorkIDs(ctx, workIDs...)
	if err != nil {
		return err
	}

	if len(performed) > 0 {
		u.lggr.Debugw("Fetched performed logs", "performed", len(performed))
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	for _, workID := range performed {
		if _, ok := u.cache[workID]; !ok {
			s := &upkeepStateRecord{
				workID:  workID,
				state:   ocr2keepers.Performed,
				addedAt: time.Now(),
			}

			u.cache[workID] = s
		}
	}

	return nil
}

// fetchFromDB fetches all upkeeps indicated as ineligible from the db to
// populate the cache.
func (u *upkeepStateStore) fetchFromDB(ctx context.Context, workIDs ...string) error {
	states, err := u.orm.SelectStatesByWorkIDs(ctx, workIDs)
	if err != nil {
		return err
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	for _, state := range states {
		if _, ok := u.cache[state.WorkID]; !ok {
			u.cache[state.WorkID] = &upkeepStateRecord{
				workID:  state.WorkID,
				state:   ocr2keepers.UpkeepState(state.CompletionState),
				addedAt: state.InsertedAt,
			}
		}
	}

	return nil
}

// selectFromCache returns all saved state values for the provided ids,
// returning stateNotFound for any ids that are not found.
// the second return value is true if all ids were found in the cache.
func (u *upkeepStateStore) selectFromCache(workIDs ...string) ([]ocr2keepers.UpkeepState, []string) {
	u.mu.RLock()
	defer u.mu.RUnlock()

	var missing []string
	states := make([]ocr2keepers.UpkeepState, len(workIDs))
	for i, workID := range workIDs {
		if state, ok := u.cache[workID]; ok {
			states[i] = state.state
		} else {
			missing = append(missing, workID)
		}
	}

	return states, missing
}

// cleanup removes any records that are older than the TTL from both cache and DB.
func (u *upkeepStateStore) cleanup(ctx context.Context) error {
	u.cleanCache()

	return u.cleanDB(ctx)
}

// cleanDB cleans up records in the DB that are older than the TTL.
func (u *upkeepStateStore) cleanDB(ctx context.Context) error {
	tm := time.Now().Add(-1 * u.retention)

	ctx, cancel := context.WithTimeout(sqlutil.WithoutDefaultTimeout(ctx), time.Minute)
	defer cancel()
	return u.orm.DeleteExpired(ctx, tm)
}

// cleanupCache removes any records from the cache that are older than the TTL.
func (u *upkeepStateStore) cleanCache() {
	u.mu.Lock()
	defer u.mu.Unlock()

	for id, state := range u.cache {
		if time.Since(state.addedAt) > u.retention {
			delete(u.cache, id)
		}
	}
}
