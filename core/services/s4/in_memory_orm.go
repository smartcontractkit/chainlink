package s4

import (
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

type inMemoryOrm struct {
	entries map[string]*Entry
	mu      sync.RWMutex
}

var _ ORM = (*inMemoryOrm)(nil)

func NewInMemoryORM() ORM {
	return &inMemoryOrm{
		entries: make(map[string]*Entry),
	}
}

func (o *inMemoryOrm) Get(address common.Address, slotId uint, qopts ...pg.QOpt) (*Entry, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	key := fmt.Sprintf("%s_%d", address, slotId)
	entry, ok := o.entries[key]
	if !ok {
		return nil, ErrRecordNotFound
	}
	return entry.Clone(), nil
}

func (o *inMemoryOrm) Upsert(address common.Address, slotId uint, entry *Entry, qopts ...pg.QOpt) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	key := fmt.Sprintf("%s_%d", address, slotId)
	o.entries[key] = entry.Clone()
	return nil
}

func (o *inMemoryOrm) DeleteExpired(qopts ...pg.QOpt) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	queue := make([]string, 0)
	now := time.Now().UnixMilli()
	for k, v := range o.entries {
		if v.HighestExpiration < now {
			queue = append(queue, k)
		}
	}
	for _, k := range queue {
		delete(o.entries, k)
	}

	return nil
}
