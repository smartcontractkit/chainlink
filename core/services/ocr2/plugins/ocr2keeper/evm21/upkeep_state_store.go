package evm

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type upkeepState struct {
	UpkeepID        []byte
	WorkID          string
	CompletionState uint8
	BlockNumber     uint64
}

type UpkeepStateReader interface {
	SelectByWorkIDs(workIDs ...string) (map[string]ocr2keepers.UpkeepState, error)
}

var (
	_ ocr2keepers.UpkeepStateUpdater = &UpkeepStateStore{}
	_ UpkeepStateReader              = &UpkeepStateStore{}
)

type UpkeepStateStore struct {
	// dependencies
	orm  *ORM
	lggr logger.Logger

	// configuration
	pruneDepth   time.Duration
	pruneCadence time.Duration

	// service values
	utils.StartStopOnce
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewUpkeepStateStore creates a new state store
func NewUpkeepStateStore(orm *ORM, lggr logger.Logger) *UpkeepStateStore {
	return &UpkeepStateStore{
		orm:          orm,
		lggr:         lggr.Named("UpkeepStateStore"),
		pruneDepth:   24 * time.Hour,
		pruneCadence: 10 * time.Second,
	}
}

// Start will start the store as a service and automatically prune stale data.
// Calling this function again is a noop
func (u *UpkeepStateStore) Start(parentCtx context.Context) error {
	if u.pruneDepth == 0 {
		return errors.Errorf("pruneDepth %d must be greater than zero", u.pruneDepth)
	}

	return u.StartOnce("UpkeepStateStore", func() error {
		ctx, cancel := context.WithCancel(parentCtx)

		u.ctx = ctx
		u.cancel = cancel

		u.wg.Add(1)
		go u.run()

		return nil
	})
}

// Close stops the service of pruning stale data
func (u *UpkeepStateStore) Close() error {
	return u.StopOnce("UpkeepStateStore", func() error {
		u.cancel()
		u.wg.Wait()

		return nil
	})
}

// Name provides the service name
func (u *UpkeepStateStore) Name() string {
	return u.lggr.Name()
}

// HealthReport provides a health report on the service
func (u *UpkeepStateStore) HealthReport() map[string]error {
	return map[string]error{u.Name(): u.StartStopOnce.Healthy()}
}

// SelectByWorkIDs returns all saved state values for the provided ids
func (u *UpkeepStateStore) SelectByWorkIDs(workIDs ...string) (map[string]ocr2keepers.UpkeepState, error) {
	states, err := u.orm.SelectStatesByWorkIDs(workIDs)
	if err != nil {
		return nil, err
	}

	//

	statesMap := make(map[string]ocr2keepers.UpkeepState)
	for _, state := range states {
		statesMap[state.WorkID] = ocr2keepers.UpkeepState(state.CompletionState)
	}

	return statesMap, nil
}

// SetUpkeepState applies the provided state to the data store for the provided
// check result. If an item already exists in the data store, the state and
// block are updated. Otherwise, the new state is added.
func (u *UpkeepStateStore) SetUpkeepState(_ context.Context, result ocr2keepers.CheckResult, _ ocr2keepers.UpkeepState) error {
	// only if the result is ineligible is the upkeep state updated
	if !result.Eligible {
		storedState := upkeepState{
			UpkeepID:        result.UpkeepID.BigInt().Bytes(),
			WorkID:          result.WorkID,
			CompletionState: uint8(ocr2keepers.Ineligible),
			BlockNumber:     uint64(result.Trigger.BlockNumber),
		}

		return u.orm.InsertUpkeepState(storedState)
	}

	return nil
}

func (u *UpkeepStateStore) run() {
	pruneTick := time.After(u.pruneCadence)

	for {
		select {
		case <-u.ctx.Done():
			return
		case <-pruneTick:
			if err := u.pruneStoredValues(); err != nil {
				u.lggr.Errorw("unable to prune old state values", "err", err)
			}

			pruneTick = time.After(utils.WithJitter(u.pruneCadence))
		}
	}
}

func (u *UpkeepStateStore) pruneStoredValues() error {
	tm := time.Now().Add(-1 * u.pruneDepth)

	return u.orm.DeleteBeforeTime(tm)
}
