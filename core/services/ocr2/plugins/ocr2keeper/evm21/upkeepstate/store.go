package upkeepstate

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

const (
	// CacheExpiration is the amount of time that we keep a record in the cache.
	CacheExpiration = 24 * time.Hour
	// GCInterval is the amount of time between cache cleanups.
	GCInterval = 2 * time.Hour
	// TODO: use sentinel value from ocr2keepers
	StateUnknown = ocr2keepers.UpkeepState(100)
)

// UpkeepStateReader is the interface for reading the current state of upkeeps.
type UpkeepStateReader interface {
	SelectByWorkIDsInRange(ctx context.Context, start, end int64, workIDs ...string) ([]ocr2keepers.UpkeepState, error)
}

type ORM interface {
	InsertUpkeepState(PersistedStateRecord, ...pg.QOpt) error
	SelectStatesByWorkIDs([]string, ...pg.QOpt) ([]PersistedStateRecord, error)
	DeleteExpired(time.Time, ...pg.QOpt) error
}

// UpkeepStateStore is the interface for managing upkeeps final state in a local store.
type UpkeepStateStore interface {
	ocr2keepers.UpkeepStateUpdater
	UpkeepStateReader
}

var (
	_ UpkeepStateStore = &upkeepStateStore{}
)

// upkeepStateRecord is a record that we save in a local cache.
type upkeepStateRecord struct {
	WorkID          string
	CompletionState ocr2keepers.UpkeepState
	BlockNumber     uint64

	AddedAt time.Time
}

// upkeepStateStore implements UpkeepStateStore.
// It stores the state of ineligible upkeeps in a local, in-memory cache (TODO: save in DB).
// In addition, performed events are fetched by the scanner on demand.
type upkeepStateStore struct {
	// dependencies
	orm     ORM
	lggr    logger.Logger
	scanner PerformedLogsScanner

	// configuration
	retention    time.Duration
	cleanCadence time.Duration

	mu    sync.RWMutex
	cache map[string]*upkeepStateRecord

	// service values
	cancel context.CancelFunc
}

// NewUpkeepStateStore creates a new state store
func NewUpkeepStateStore(orm ORM, lggr logger.Logger, scanner PerformedLogsScanner) *upkeepStateStore {
	return &upkeepStateStore{
		orm:          orm,
		lggr:         lggr.Named("UpkeepStateStore"),
		cache:        map[string]*upkeepStateRecord{},
		scanner:      scanner,
		retention:    CacheExpiration,
		cleanCadence: GCInterval,
	}
}

// Start starts the upkeep state store.
// it does background cleanup of the cache.
func (u *upkeepStateStore) Start(pctx context.Context) error {
	if u.retention == 0 {
		return fmt.Errorf("pruneDepth %d must be greater than zero", u.retention)
	}

	ctx, cancel := context.WithCancel(pctx)
	defer cancel()

	u.mu.Lock()
	u.cancel = cancel
	u.mu.Unlock()

	u.lggr.Debug("Starting upkeep state store")

	cleanTick := time.NewTicker(u.cleanCadence)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-cleanTick.C:
			if err := u.cleanup(); err != nil {
				u.lggr.Errorw("unable to clean old state values", "err", err)
			}
		}
	}
}

// Close stops the service of pruning stale data; implements io.Closer
func (u *upkeepStateStore) Close() error {
	u.mu.Lock()
	cancel := u.cancel
	u.mu.Unlock()

	if cancel != nil {
		cancel()
	}

	return nil
}

// SelectByWorkIDs returns the current state of the upkeep for the provided ids.
// If an id is not found, the state is returned as StateUnknown.
// We first check the cache, and if any ids are missing, we fetch them from the scanner.
func (u *upkeepStateStore) SelectByWorkIDsInRange(ctx context.Context, start, end int64, workIDs ...string) ([]ocr2keepers.UpkeepState, error) {
	states, ok := u.selectFromCache(workIDs...)
	if ok {
		// all ids were found in the cache
		return states, nil
	}

	// fetch values from chain to populate the cache with missing values
	if err := u.fetchPerformed(ctx, start, end); err != nil {
		return nil, err
	}

	idsWithUnknownState := []string{}
	for i, id := range states {
		if id == StateUnknown {
			idsWithUnknownState = append(idsWithUnknownState, workIDs[i])
		}
	}

	if len(idsWithUnknownState) > 0 {
		// fetch values from the db to populate the cache with missing values here
		if err := u.fetchFromDB(ctx, idsWithUnknownState...); err != nil {
			return nil, err
		}
	}

	// at this point all values should be in the cache. if values are missing
	// their state is indicated as unknown
	states, _ = u.selectFromCache(workIDs...)

	return states, nil
}

// SetUpkeepState updates the state of the upkeep.
// Currently we only store the state if the upkeep is ineligible.
// Performed events will be fetched on demand.
func (u *upkeepStateStore) SetUpkeepState(_ context.Context, result ocr2keepers.CheckResult, _ ocr2keepers.UpkeepState) error {
	if result.Eligible {
		return nil
	}

	return u.upsertStateRecord(result.WorkID, ocr2keepers.Ineligible, uint64(result.Trigger.BlockNumber), result.UpkeepID.BigInt(), result.IneligibilityReason)
}

// upsertStateRecord inserts or updates a record for the provided
// check result. If an item already exists in the data store, the state and
// block are updated.
func (u *upkeepStateStore) upsertStateRecord(workID string, s ocr2keepers.UpkeepState, b uint64, upkeepID *big.Int, reason uint8) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	record, ok := u.cache[workID]
	if !ok {
		record = &upkeepStateRecord{
			WorkID:  workID,
			AddedAt: time.Now(),
		}
	}

	record.CompletionState = s
	record.BlockNumber = b

	u.cache[workID] = record

	return u.orm.InsertUpkeepState(PersistedStateRecord{
		UpkeepID:            upkeepID,
		WorkID:              record.WorkID,
		CompletionState:     uint8(record.CompletionState),
		BlockNumber:         record.BlockNumber,
		IneligibilityReason: reason,
		AddedAt:             record.AddedAt,
	})
}

// fetchPerformed fetches all performed logs from the scanner to populate the cache.
func (u *upkeepStateStore) fetchPerformed(ctx context.Context, start, end int64, workIDs ...string) error {
	performed, err := u.scanner.WorkIDsInRange(ctx, start, end)
	if err != nil {
		return err
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	for _, workID := range performed {
		if _, ok := u.cache[workID]; !ok {
			s := &upkeepStateRecord{
				WorkID:          workID,
				CompletionState: ocr2keepers.Performed,
				AddedAt:         time.Now(),
				BlockNumber:     uint64(end), // TODO: use block number from log
			}

			u.cache[workID] = s
		}
	}

	return nil
}

// fetchFromDB fetches all upkeeps indicated as ineligible from the db to
// populate the cache
func (u *upkeepStateStore) fetchFromDB(ctx context.Context, workIDs ...string) error {
	states, err := u.orm.SelectStatesByWorkIDs(workIDs)
	if err != nil {
		return err
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	for _, state := range states {
		if _, ok := u.cache[state.WorkID]; !ok {
			u.cache[state.WorkID] = &upkeepStateRecord{
				WorkID:          state.WorkID,
				CompletionState: ocr2keepers.UpkeepState(state.CompletionState),
				BlockNumber:     state.BlockNumber,
				AddedAt:         state.AddedAt,
			}
		}
	}

	return nil
}

// selectFromCache returns all saved state values for the provided ids,
// returning stateNotFound for any ids that are not found.
// the second return value is true if all ids were found in the cache.
func (u *upkeepStateStore) selectFromCache(workIDs ...string) ([]ocr2keepers.UpkeepState, bool) {
	u.mu.RLock()
	defer u.mu.RUnlock()

	var hasMisses bool
	states := make([]ocr2keepers.UpkeepState, len(workIDs))
	for i, workID := range workIDs {
		if state, ok := u.cache[workID]; ok {
			states[i] = state.CompletionState
		} else {
			hasMisses = true
			states[i] = StateUnknown
		}
	}

	return states, !hasMisses
}

// cleanup removes any records that are older than the TTL from both cache and DB.
func (u *upkeepStateStore) cleanup() error {
	u.cleanCache()

	return u.cleanDB()
}

// cleanDB cleans up records in the DB that are older than the TTL.
func (u *upkeepStateStore) cleanDB() error {
	tm := time.Now().Add(-1 * u.retention)

	return u.orm.DeleteExpired(tm)
}

// cleanupCache removes any records from the cache that are older than the TTL.
func (u *upkeepStateStore) cleanCache() {
	u.mu.Lock()
	defer u.mu.Unlock()

	for id, state := range u.cache {
		if time.Since(state.AddedAt) > u.retention {
			delete(u.cache, id)
		}
	}
}
