package s4

import (
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type key struct {
	address string
	slot    uint
}

type inMemoryOrm struct {
	rows map[key]*Row
	mu   sync.RWMutex
}

var _ ORM = (*inMemoryOrm)(nil)

func NewInMemoryORM() ORM {
	return &inMemoryOrm{
		rows: make(map[key]*Row),
	}
}

func (o *inMemoryOrm) Get(address common.Address, slotId uint, qopts ...pg.QOpt) (*Row, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	mkey := key{
		address: address.String(),
		slot:    slotId,
	}
	row, ok := o.rows[mkey]
	if !ok {
		return nil, ErrNotFound
	}
	return row.Clone(), nil
}

func (o *inMemoryOrm) Update(row *Row, qopts ...pg.QOpt) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	mkey := key{
		address: row.Address,
		slot:    row.SlotId,
	}
	existing, ok := o.rows[mkey]
	if ok && existing.Version >= row.Version {
		return ErrVersionTooLow
	}

	clone := row.Clone()
	clone.UpdatedAt = time.Now().UnixMilli()
	o.rows[mkey] = clone
	return nil
}

func (o *inMemoryOrm) DeleteExpired(qopts ...pg.QOpt) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	queue := make([]key, 0)
	now := time.Now().UnixMilli()
	for k, v := range o.rows {
		if v.Expiration < now {
			queue = append(queue, k)
		}
	}
	for _, k := range queue {
		delete(o.rows, k)
	}

	return nil
}

func (o *inMemoryOrm) GetSnapshot(addressRange *AddressRange, qopts ...pg.QOpt) ([]*Row, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	now := time.Now().UnixMilli()
	var selection []key
	for k, v := range o.rows {
		bigAddress := utils.NewBig(common.HexToAddress(v.Address).Big())
		if v.Expiration > now && addressRange.Contains(bigAddress) {
			selection = append(selection, k)
		}
	}

	rows := make([]*Row, len(selection))
	for i, s := range selection {
		rows[i] = o.rows[s].Clone()
	}

	return rows, nil
}
