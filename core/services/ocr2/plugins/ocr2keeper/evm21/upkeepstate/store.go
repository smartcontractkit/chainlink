package upkeepstate

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/big"
	"sync"
	"time"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	// CacheExpiration is the amount of time that we keep a record in the cache.
	CacheExpiration = 24 * time.Hour
	// GCInterval is the amount of time between cache cleanups.
	GCInterval = 2 * time.Hour
)

type ORM interface {
	InsertUpkeepState(persistedStateRecord, ...pg.QOpt) error
	SelectStatesByWorkIDs([]string, ...pg.QOpt) ([]persistedStateRecord, error)
	DeleteExpired(time.Time, ...pg.QOpt) error
}

// UpkeepStateStore is the interface for managing upkeeps final state in a local store.
type UpkeepStateStore interface {
	ocr2keepers.UpkeepStateUpdater
	core.UpkeepStateReader
	Start(context.Context) error
	io.Closer
}

var (
	_ UpkeepStateStore = &upkeepStateStore{}
)

// upkeepStateRecord is a record that we save in a local cache.
type upkeepStateRecord struct {
	workID      string
	state       ocr2keepers.UpkeepState
	blockNumber uint64

	addedAt time.Time
}

// upkeepStateStore implements UpkeepStateStore.
// It stores the state of ineligible upkeeps in a local, in-memory cache.
// In addition, performed events are fetched by the scanner on demand.
// TODO: Add DB persistence
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
		return errors.New("pruneDepth must be greater than zero")
	}

	u.mu.Lock()
	if u.cancel != nil {
		u.mu.Unlock()
		return fmt.Errorf("already started")
	}

	ctx, cancel := context.WithCancel(context.Background())

	u.cancel = cancel
	u.mu.Unlock()

	if err := u.scanner.Start(ctx); err != nil {
		return fmt.Errorf("failed to start scanner")
	}

	u.lggr.Debug("Starting upkeep state store")

	{
		go func(ctx context.Context) {
			ticker := time.NewTicker(u.cleanCadence)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					if err := u.cleanup(ctx); err != nil {
						u.lggr.Errorw("unable to clean old state values", "err", err)
					}

					ticker.Reset(utils.WithJitter(u.cleanCadence))
				case <-ctx.Done():

				}
			}
		}(ctx)
	}

	return nil
}

// Close stops the service of pruning stale data; implements io.Closer
func (u *upkeepStateStore) Close() error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if cancel := u.cancel; cancel != nil {
		u.cancel = nil
		cancel()
	} else {
		return fmt.Errorf("already stopped")
	}
	if err := u.scanner.Close(); err != nil {
		return fmt.Errorf("failed to start scanner")
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

	// fetch values from the db to populate the cache with missing values here
	if err := u.fetchFromDB(ctx, workIDs, states); err != nil {
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

	record.blockNumber = b
	record.state = s

	u.cache[workID] = record

	return u.orm.InsertUpkeepState(persistedStateRecord{
		UpkeepID:            utils.NewBig(upkeepID),
		WorkID:              record.workID,
		CompletionState:     uint8(record.state),
		BlockNumber:         int64(record.blockNumber),
		IneligibilityReason: reason,
		InsertedAt:          record.addedAt,
	}, pg.WithParentCtx(ctx))
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
				workID:      workID,
				state:       ocr2keepers.Performed,
				addedAt:     time.Now(),
				blockNumber: uint64(end),
			}

			u.cache[workID] = s
		}
	}

	return nil
}

// fetchFromDB fetches all upkeeps indicated as ineligible from the db to
// populate the cache
func (u *upkeepStateStore) fetchFromDB(ctx context.Context, workIDs []string, states []ocr2keepers.UpkeepState) error {
	if len(workIDs) == 0 {
		return nil
	}

	idsWithUnknownState := []string{}
	for i, state := range states {
		if state == ocr2keepers.UnknownState {
			idsWithUnknownState = append(idsWithUnknownState, workIDs[i])
		}
	}

	dbStates, err := u.orm.SelectStatesByWorkIDs(idsWithUnknownState, pg.WithParentCtx(ctx))
	if err != nil {
		return err
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	for _, state := range dbStates {
		if _, ok := u.cache[state.WorkID]; !ok {
			u.cache[state.WorkID] = &upkeepStateRecord{
				workID:      state.WorkID,
				state:       ocr2keepers.UpkeepState(state.CompletionState),
				blockNumber: uint64(state.BlockNumber),
				addedAt:     state.InsertedAt,
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
			states[i] = state.state
		} else {
			hasMisses = true
			states[i] = ocr2keepers.UnknownState
		}
	}

	return states, !hasMisses
}

// cleanup removes any records that are older than the TTL from both cache and DB.
func (u *upkeepStateStore) cleanup(ctx context.Context) error {
	u.cleanCache()

	return u.cleanDB(ctx)
}

// cleanDB cleans up records in the DB that are older than the TTL.
func (u *upkeepStateStore) cleanDB(ctx context.Context) error {
	tm := time.Now().Add(-1 * u.retention)

	return u.orm.DeleteExpired(tm, pg.WithParentCtx(ctx), pg.WithLongQueryTimeout())
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
