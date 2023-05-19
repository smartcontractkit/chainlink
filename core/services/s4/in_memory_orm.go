package s4

import (
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

type inMemoryOrm struct {
	rows map[string]*Row
	mu   sync.RWMutex
}

var _ ORM = (*inMemoryOrm)(nil)

func NewInMemoryORM() ORM {
	return &inMemoryOrm{
		rows: make(map[string]*Row),
	}
}

func (o *inMemoryOrm) Get(address common.Address, slotId uint, qopts ...pg.QOpt) (*Row, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	key := fmt.Sprintf("%s_%d", address, slotId)
	row, ok := o.rows[key]
	if !ok {
		return nil, ErrNotFound
	}
	return row.Clone(), nil
}

func (o *inMemoryOrm) Upsert(address common.Address, slotId uint, row *Row, qopts ...pg.QOpt) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	key := fmt.Sprintf("%s_%d", address, slotId)
	o.rows[key] = row.Clone()
	return nil
}

func (o *inMemoryOrm) DeleteExpired(qopts ...pg.QOpt) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	queue := make([]string, 0)
	now := time.Now().UnixMilli()
	for k, v := range o.rows {
		if v.HighestExpiration < now {
			queue = append(queue, k)
		}
	}
	for _, k := range queue {
		delete(o.rows, k)
	}

	return nil
}
