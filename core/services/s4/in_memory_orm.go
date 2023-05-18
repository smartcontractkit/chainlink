package s4

import (
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

type inMemoryOrm struct {
	entires map[string]*Entry
	mu      sync.RWMutex
}

func NewInMemoryORM() ORM {
	return &inMemoryOrm{
		entires: make(map[string]*Entry),
	}
}

func (o *inMemoryOrm) Get(address common.Address, slotId int, qopts ...pg.QOpt) (*Entry, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	key := fmt.Sprintf("%s_%d", address, slotId)
	entry, ok := o.entires[key]
	if !ok {
		return nil, ErrEntryNotFound
	}
	return entry.Clone(), nil
}

func (o *inMemoryOrm) Upsert(address common.Address, slotId int, entry *Entry, qopts ...pg.QOpt) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	key := fmt.Sprintf("%s_%d", address, slotId)
	o.entires[key] = entry.Clone()
	return nil
}

func (o *inMemoryOrm) DeleteExpired(qopts ...pg.QOpt) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	queue := make([]string, 0)
	now := time.Now().UnixMilli()
	for k, v := range o.entires {
		if v.HighestExpiration < now {
			queue = append(queue, k)
		}
	}
	for _, k := range queue {
		delete(o.entires, k)
	}

	return nil
}
