package evm

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/ocr2keepers/pkg/types"
)

var (
	zeroAddr = common.Address{}
)

// UpkeepFetcher is a function to get a single upkeep info
type UpkeepFetcher func(context.Context, *big.Int, *big.Int) (upkeepInfoEntry, error)

// UpkeepFetcher is a function to get a batch of upkeeps info
type UpkeepBatchFetcher func(context.Context, *big.Int, []*big.Int) ([]upkeepInfoEntry, error)

// UpkeepFilter is a function that can be used to filter upkeeps
type UpkeepFilter func(types.UpkeepIdentifier, upkeepInfoEntry) bool

// upkeepIndex holds available upkeeps and enable other components to query them.
// TODO: limits, seperate cache for info? config? state?
type upkeepIndex struct {
	lggr logger.Logger

	// info *cache.Cache // TBD
	upkeeps map[string]upkeepInfoEntry
	lock    *sync.RWMutex

	fetcher      UpkeepFetcher
	batchFetcher UpkeepBatchFetcher
}

func newUpkeepIndex(lggr logger.Logger, fetcher UpkeepFetcher, batchFetcher UpkeepBatchFetcher) *upkeepIndex {
	return &upkeepIndex{
		lggr: lggr,

		upkeeps: make(map[string]upkeepInfoEntry),
		lock:    new(sync.RWMutex),

		fetcher:      fetcher,
		batchFetcher: batchFetcher,
	}
}

// Initialize enables to initialize the index with the upkeeps from the registry
func (index *upkeepIndex) Initialize(ctx context.Context, ids ...*big.Int) error {
	infos, err := index.batchFetcher(ctx, nil, ids)
	if err != nil {
		return errors.Wrap(err, "failed to get configs to initialize the upkeep index")
	}
	upkeeps := make(map[string]upkeepInfoEntry)
	for _, upkeepInfo := range infos {
		upkeeps[upkeepInfo.id.String()] = upkeepInfo
	}

	index.lock.Lock()
	defer index.lock.Unlock()
	index.upkeeps = upkeeps

	return nil
}

// AddActiveUpkeep adds an active upkeep to the index
func (index *upkeepIndex) AddActiveUpkeep(ctx context.Context, id *big.Int, force bool) error {
	_, err := index.GetUpkeepInfo(ctx, nil, id, force)
	return err
}

// GetActiveUpkeepIDs returns the IDs of filtered active upkeeps.
func (index *upkeepIndex) GetActiveUpkeepIDs(filters ...UpkeepFilter) []types.UpkeepIdentifier {
	filters = append([]UpkeepFilter{ActiveUpkeepsFilter()}, filters...)
	return index.GetUpkeepIDs(filters...)
}

// GetActiveUpkeepIDs returns the IDs of filtered upkeeps.
func (index *upkeepIndex) GetUpkeepIDs(filters ...UpkeepFilter) []types.UpkeepIdentifier {
	index.lock.RLock()
	defer index.lock.RUnlock()

	f := upkeepFilters(filters)
	var ids []types.UpkeepIdentifier
	for id, entry := range index.upkeeps {
		uid := types.UpkeepIdentifier(id)
		if f.Apply(uid, entry) {
			ids = append(ids, uid)
		}
	}

	return ids
}

// GetUpkeepConfig returns the offchain config of the upkeep
func (index *upkeepIndex) GetUpkeepConfig(ctx context.Context, upkeepID *big.Int) ([]byte, error) {
	upkeepInfo, err := index.GetUpkeepInfo(ctx, nil, upkeepID, false)
	return upkeepInfo.offchainConfig, err
}

// GetUpkeepInfo returns the onchain info of the upkeep
func (index *upkeepIndex) GetUpkeepInfo(ctx context.Context, block *big.Int, upkeepID *big.Int, force bool) (upkeepInfoEntry, error) {
	upkeepInfo, found := index.getLocalUpkeepInfo(upkeepID)
	if found && !force {
		if !upkeepInfo.updatedAt.Add(DefaultUpkeepExpiration).Before(time.Now()) {
			index.lggr.Debugf("UpkeepInfo upkeep %s block %d cache hit UpkeepInfo: %+v", upkeepID.String(), block, upkeepInfo)
			return upkeepInfo, nil
		}
		// otherwise the entry is expired, fetching...
	}

	upkeepInfo, err := index.fetcher(ctx, block, upkeepID)
	if err != nil {
		return upkeepInfo, err
	}
	index.lggr.Debugf("UpkeepInfo upkeep %s block %d cache miss UpkeepInfo: %+v", upkeepID.String(), block, upkeepInfo)
	index.setLocalUpkeepInfo(upkeepID, upkeepInfo)
	return upkeepInfo, nil
}

func (index *upkeepIndex) getLocalUpkeepInfo(upkeepID *big.Int) (upkeepInfoEntry, bool) {
	index.lock.RLock()
	defer index.lock.RUnlock()

	entry, found := index.upkeeps[upkeepID.String()]
	if !found {
		return upkeepInfoEntry{}, false
	}
	// TODO: check if expired
	return entry, true
}

func (index *upkeepIndex) setLocalUpkeepInfo(upkeepID *big.Int, upkeepEntry upkeepInfoEntry) {
	index.lock.Lock()
	defer index.lock.Unlock()

	index.upkeeps[upkeepID.String()] = upkeepEntry
}

func (index *upkeepIndex) setUpkeepConfig(upkeepID *big.Int, cfg []byte) {
	index.lock.Lock()
	defer index.lock.Unlock()

	entry, found := index.upkeeps[upkeepID.String()]
	if !found {
		entry = upkeepInfoEntry{
			id: upkeepID,
		}
	}
	entry.offchainConfig = cfg
	entry.updatedAt = time.Now()
	index.upkeeps[upkeepID.String()] = entry
}

func (index *upkeepIndex) setUpkeepState(upkeepID *big.Int, state upkeepState) {
	index.lock.Lock()
	defer index.lock.Unlock()

	entry, found := index.upkeeps[upkeepID.String()]
	if !found {
		entry = upkeepInfoEntry{}
	}
	entry.state = state
	entry.updatedAt = time.Now()
	index.upkeeps[upkeepID.String()] = entry
}

func getUpkeepType(id types.UpkeepIdentifier) upkeepType {
	// TODO: implement
	return blockTrigger
}

type upkeepFilters []UpkeepFilter

func (uf upkeepFilters) Apply(id types.UpkeepIdentifier, upkeep upkeepInfoEntry) bool {
	for _, f := range uf {
		if !f(id, upkeep) {
			return false
		}
	}
	return true
}

func ActiveUpkeepsFilter() UpkeepFilter {
	return func(_ types.UpkeepIdentifier, entry upkeepInfoEntry) bool {
		return entry.state == stateActive
	}
}

func LogUpkeepsFilter() UpkeepFilter {
	return func(id types.UpkeepIdentifier, _ upkeepInfoEntry) bool {
		return getUpkeepType(id) == logTrigger
	}
}

func BlockUpkeepsFilter() UpkeepFilter {
	return func(id types.UpkeepIdentifier, _ upkeepInfoEntry) bool {
		return getUpkeepType(id) == blockTrigger
	}
}
