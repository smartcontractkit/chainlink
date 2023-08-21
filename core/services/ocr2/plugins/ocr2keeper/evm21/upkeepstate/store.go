package upkeepstate

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
)

var (
	// CacheExpiration is the amount of time that we keep a record in the cache.
	CacheExpiration = 24 * time.Hour
	// GCInterval is the amount of time between cache cleanups.
	GCInterval = 2 * time.Hour
)

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
	workID string
	state  ocr2keepers.UpkeepState

	addedAt time.Time
}

// upkeepStateStore implements UpkeepStateStore.
// It stores the state of ineligible upkeeps in a local, in-memory cache.
// In addition, performed events are fetched by the scanner on demand.
// TODO: Add DB persistence
type upkeepStateStore struct {
	lggr logger.Logger

	cancel context.CancelFunc

	mu    sync.RWMutex
	cache map[string]*upkeepStateRecord

	scanner PerformedLogsScanner
}

// NewUpkeepStateStore creates a new state store
func NewUpkeepStateStore(lggr logger.Logger, scanner PerformedLogsScanner) *upkeepStateStore {
	return &upkeepStateStore{
		lggr:    lggr.Named("upkeepStateStore"),
		cache:   map[string]*upkeepStateRecord{},
		scanner: scanner,
	}
}

// Start starts the upkeep state store.
// it does background cleanup of the cache.
func (u *upkeepStateStore) Start(context.Context) error {
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
			ticker := time.NewTicker(GCInterval)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					u.cleanup()
				case <-ctx.Done():

				}
			}
		}(ctx)
	}

	return nil
}

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
	if err := u.fetchPerformed(ctx, start, end); err != nil {
		return nil, err
	}
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

	u.upsertStateRecord(result.WorkID, ocr2keepers.Ineligible)

	return nil
}

// upsertStateRecord inserts or updates a record for the provided
// check result. If an item already exists in the data store, the state and
// block are updated.
func (u *upkeepStateStore) upsertStateRecord(workID string, s ocr2keepers.UpkeepState) {
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
				workID:  workID,
				state:   ocr2keepers.Performed,
				addedAt: time.Now(),
			}
			u.cache[workID] = s
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
func (u *upkeepStateStore) cleanup() {
	u.cleanCache()
	u.cleanDB()
}

// cleanDB cleans up records in the DB that are older than the TTL.
func (u *upkeepStateStore) cleanDB() {
}

// cleanupCache removes any records from the cache that are older than the TTL.
func (u *upkeepStateStore) cleanCache() {
	u.mu.Lock()
	defer u.mu.Unlock()

	for id, state := range u.cache {
		if time.Since(state.addedAt) > CacheExpiration {
			delete(u.cache, id)
		}
	}
}
