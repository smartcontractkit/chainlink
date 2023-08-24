package logprovider

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/time/rate"
)

type UpkeepFilterStore interface {
	GetIDs(selector func(upkeepFilter) bool) []*big.Int
	UpdateFilters(updater func(upkeepFilter, upkeepFilter) upkeepFilter, filters ...upkeepFilter)
	Has(id *big.Int) bool
	Get(id *big.Int) *upkeepFilter
	RangeFiltersByIDs(iterator func(int, upkeepFilter), ids ...*big.Int)
	GetFilters(selector func(upkeepFilter) bool) []upkeepFilter
	AddActiveUpkeeps(filters ...upkeepFilter)
	RemoveActiveUpkeeps(filters ...upkeepFilter)
	Size() int
}

var _ UpkeepFilterStore = &upkeepFilterStore{}

type upkeepFilter struct {
	addr     []byte
	topics   []common.Hash
	upkeepID *big.Int
	// configUpdateBlock is the block number the filter was last updated at
	configUpdateBlock uint64
	// lastPollBlock is the last block number the logs were fetched for this upkeep
	// used by log event provider.
	lastPollBlock int64
	// blockLimiter is used to limit the number of blocks to fetch logs for an upkeep.
	// used by log event provider.
	blockLimiter *rate.Limiter
	// lastRePollBlock is the last block number the logs were recovered for this upkeep
	// used by log recoverer.
	lastRePollBlock int64
}

func (f upkeepFilter) Clone() upkeepFilter {
	topics := make([]common.Hash, len(f.topics))
	copy(topics, f.topics)
	addr := make([]byte, len(f.addr))
	copy(addr, f.addr)
	return upkeepFilter{
		upkeepID:          f.upkeepID,
		topics:            topics,
		addr:              addr,
		configUpdateBlock: f.configUpdateBlock,
		lastPollBlock:     f.lastPollBlock,
		lastRePollBlock:   f.lastRePollBlock,
		blockLimiter:      f.blockLimiter,
	}
}

type upkeepFilterStore struct {
	lock *sync.RWMutex
	// filters is a map of upkeepID to upkeepFilter
	filters map[string]upkeepFilter
}

func NewUpkeepFilterStore() *upkeepFilterStore {
	return &upkeepFilterStore{
		lock:    &sync.RWMutex{},
		filters: make(map[string]upkeepFilter),
	}
}

func (s *upkeepFilterStore) GetIDs(selector func(upkeepFilter) bool) []*big.Int {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if selector == nil {
		// noop selector returns true for all filters
		selector = func(upkeepFilter) bool { return true }
	}

	var ids []*big.Int
	for _, f := range s.filters {
		if selector(f) {
			ids = append(ids, f.upkeepID)
		}
	}

	return ids
}

func (s *upkeepFilterStore) UpdateFilters(resolveUpdated func(upkeepFilter, upkeepFilter) upkeepFilter, filters ...upkeepFilter) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if resolveUpdated == nil {
		// noop resolveUpdated will use the newer filter
		resolveUpdated = func(_ upkeepFilter, f upkeepFilter) upkeepFilter { return f }
	}

	for _, f := range filters {
		uid := f.upkeepID.String()
		orig, ok := s.filters[uid]
		if !ok {
			// not found, turned inactive probably
			continue
		}
		updated := resolveUpdated(orig, f)
		s.filters[uid] = updated
	}
}

func (s *upkeepFilterStore) Has(id *big.Int) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	_, ok := s.filters[id.String()]
	return ok
}

func (s *upkeepFilterStore) Get(id *big.Int) *upkeepFilter {
	s.lock.RLock()
	defer s.lock.RUnlock()

	f, ok := s.filters[id.String()]
	if !ok {
		return nil
	}
	fp := f.Clone()
	return &fp
}

func (s *upkeepFilterStore) RangeFiltersByIDs(iterator func(int, upkeepFilter), ids ...*big.Int) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if iterator == nil {
		// noop iterator does nothing
		iterator = func(int, upkeepFilter) {}
	}

	for i, id := range ids {
		f, ok := s.filters[id.String()]
		if !ok {
			// in case the filter is not found, we still want to call the iterator
			// with an empty filter, so
			iterator(i, upkeepFilter{upkeepID: id})
		} else {
			iterator(i, f)
		}
	}
}

func (s *upkeepFilterStore) GetFilters(selector func(upkeepFilter) bool) []upkeepFilter {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if selector == nil {
		// noop selector returns true for all filters
		selector = func(upkeepFilter) bool { return true }
	}

	var filters []upkeepFilter
	for _, f := range s.filters {
		if selector(f) {
			filters = append(filters, f.Clone())
		}
	}
	return filters
}

func (s *upkeepFilterStore) AddActiveUpkeeps(filters ...upkeepFilter) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, f := range filters {
		s.filters[f.upkeepID.String()] = f
	}
}

func (s *upkeepFilterStore) RemoveActiveUpkeeps(filters ...upkeepFilter) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, f := range filters {
		uid := f.upkeepID.String()
		delete(s.filters, uid)
	}
}

func (s *upkeepFilterStore) Size() int {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return len(s.filters)
}
