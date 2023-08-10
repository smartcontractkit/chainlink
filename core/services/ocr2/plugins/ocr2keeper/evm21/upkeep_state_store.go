package evm

import (
	"context"
	"sync"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type upkeepState struct {
	upkeepID ocr2keepers.UpkeepIdentifier
	workID   string
	state    ocr2keepers.UpkeepState
	block    uint64
}

type UpkeepStateReader interface {
	SelectByWorkIDs(workIDs ...string) (map[string]ocr2keepers.UpkeepState, error)
}

var (
	_ ocr2keepers.UpkeepStateUpdater = &UpkeepStateStore{}
	_ UpkeepStateReader              = &UpkeepStateStore{}
)

type UpkeepStateStore struct {
	mu        sync.RWMutex
	workIDIdx map[string]*upkeepState
	states    []*upkeepState
	lggr      logger.Logger
}

// NewUpkeepStateStore creates a new state store. This is an initial version of this store. More improvements to come:
// TODO: AUTO-4027
func NewUpkeepStateStore(lggr logger.Logger) *UpkeepStateStore {
	return &UpkeepStateStore{
		states:    []*upkeepState{},
		workIDIdx: map[string]*upkeepState{},
		lggr:      lggr.Named("UpkeepStateStore"),
	}
}

// SelectByWorkIDs returns all saved state values for the provided ids
func (u *UpkeepStateStore) SelectByWorkIDs(workIDs ...string) (map[string]ocr2keepers.UpkeepState, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()

	states := make(map[string]ocr2keepers.UpkeepState)

	for _, workID := range workIDs {
		if state, ok := u.workIDIdx[workID]; ok {
			states[workID] = state.state
		}
	}

	return states, nil
}

// SetUpkeepState applies the provided state to the data store for the provided
// check result. If an item already exists in the data store, the state and
// block are updated. Otherwise, the new state is added.
func (u *UpkeepStateStore) SetUpkeepState(_ context.Context, result ocr2keepers.CheckResult, _ ocr2keepers.UpkeepState) error {
	// only if the result is ineligible is the upkeep state updated
	if !result.Eligible {
		if existing, ok := u.workIDIdx[result.WorkID]; ok {
			u.mu.Lock()

			existing.state = ocr2keepers.Ineligible
			existing.block = uint64(result.Trigger.BlockNumber)

			u.mu.Unlock()

			u.lggr.Infof("upkeep %s is overridden, workID is %s, block is %d", existing.upkeepID, existing.workID, existing.block)
		} else {
			storedState := &upkeepState{
				upkeepID: result.UpkeepID,
				workID:   result.WorkID,
				state:    ocr2keepers.Ineligible,
				block:    uint64(result.Trigger.BlockNumber),
			}

			u.mu.Lock()

			u.workIDIdx[result.WorkID] = storedState
			u.states = append(u.states, storedState)

			u.mu.Unlock()

			u.lggr.Infof("added new state with upkeep %s payload ID %s block %d", storedState.upkeepID, storedState.workID, storedState.block)
		}
	}

	return nil
}
