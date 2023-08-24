package evm

import (
	"math/big"
	"sync"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
)

// ActiveUpkeepList is a list to manage active upkeep IDs
type ActiveUpkeepList interface {
	// Reset resets the list to the given IDs
	Reset(ids ...*big.Int)
	// Add adds new entries to the list
	Add(id ...*big.Int) int
	// Remove removes entries from the list
	Remove(id ...*big.Int) int
	// View returns the list of IDs of the given type
	View(...ocr2keepers.UpkeepType) []*big.Int
	// IsActive returns true if the given ID is of an active upkeep
	IsActive(id *big.Int) bool
	Size() int
}

type activeList struct {
	items map[string]bool
	lock  sync.RWMutex
}

var _ ActiveUpkeepList = &activeList{}

// NewActiveList creates a new ActiveList
func NewActiveUpkeepList() ActiveUpkeepList {
	return &activeList{
		items: make(map[string]bool),
	}
}

// Reset resets the list to the given IDs
func (al *activeList) Reset(ids ...*big.Int) {
	al.lock.Lock()
	defer al.lock.Unlock()

	al.items = make(map[string]bool)
	for _, id := range ids {
		al.items[id.String()] = true
	}
}

// Add adds new entries to the list
func (al *activeList) Add(ids ...*big.Int) int {
	al.lock.Lock()
	defer al.lock.Unlock()

	count := 0
	for _, id := range ids {
		if key := id.String(); !al.items[key] {
			count++
			al.items[key] = true
		}
	}
	return count
}

// Remove removes entries from the list
func (al *activeList) Remove(ids ...*big.Int) int {
	al.lock.Lock()
	defer al.lock.Unlock()

	count := 0
	for _, id := range ids {
		key := id.String()
		if al.items[key] {
			count++
			delete(al.items, key)
		}
	}
	return count
}

// View returns the list of IDs of the given type
func (al *activeList) View(upkeepTypes ...ocr2keepers.UpkeepType) []*big.Int {
	al.lock.RLock()
	defer al.lock.RUnlock()

	var keys []*big.Int
	for key := range al.items {
		id := &ocr2keepers.UpkeepIdentifier{}
		bint, ok := big.NewInt(0).SetString(key, 10)
		if !ok {
			continue
		}
		if !id.FromBigInt(bint) {
			continue
		}
		currentType := core.GetUpkeepType(*id)
		for _, t := range upkeepTypes {
			if currentType == t {
				keys = append(keys, bint)
				break
			}
		}
	}
	return keys
}

func (al *activeList) IsActive(id *big.Int) bool {
	al.lock.RLock()
	defer al.lock.RUnlock()

	return al.items[id.String()]
}

func (al *activeList) Size() int {
	al.lock.RLock()
	defer al.lock.RUnlock()

	return len(al.items)
}
