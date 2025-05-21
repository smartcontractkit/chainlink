package s4

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
)

type key struct {
	address string
	slot    uint
}

type mrow struct {
	Row       *Row
	UpdatedAt time.Time
}

type inMemoryOrm struct {
	rows map[key]*mrow
	mu   sync.RWMutex
}

var _ ORM = (*inMemoryOrm)(nil)

func NewInMemoryORM() ORM {
	return &inMemoryOrm{
		rows: make(map[key]*mrow),
	}
}

func (o *inMemoryOrm) Get(ctx context.Context, address *big.Big, slotId uint) (*Row, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	mkey := key{
		address: address.Hex(),
		slot:    slotId,
	}
	mrow, ok := o.rows[mkey]
	if !ok {
		return nil, ErrNotFound
	}
	return mrow.Row.Clone(), nil
}

func (o *inMemoryOrm) Update(ctx context.Context, row *Row) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	mkey := key{
		address: row.Address.Hex(),
		slot:    row.SlotId,
	}
	existing, ok := o.rows[mkey]
	versionOk := false
	if ok && row.Confirmed {
		versionOk = existing.Row.Version <= row.Version
	}
	if ok && !row.Confirmed {
		versionOk = existing.Row.Version < row.Version
	}
	if ok && !versionOk {
		return ErrVersionTooLow
	}

	o.rows[mkey] = &mrow{
		Row:       row.Clone(),
		UpdatedAt: time.Now().UTC(),
	}
	return nil
}

func (o *inMemoryOrm) DeleteExpired(ctx context.Context, limit uint, now time.Time) (int64, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	queue := make([]key, 0)
	for k, v := range o.rows {
		if v.Row.Expiration < now.UnixMilli() {
			queue = append(queue, k)
			if len(queue) >= int(limit) {
				break
			}
		}
	}
	for _, k := range queue {
		delete(o.rows, k)
	}

	return int64(len(queue)), nil
}

func (o *inMemoryOrm) GetSnapshot(ctx context.Context, _ *AddressRange) ([]*SnapshotRow, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	now := time.Now().UnixMilli()
	var rows []*SnapshotRow
	for _, mrow := range o.rows {
		if mrow.Row.Expiration > now {
			rows = append(rows, &SnapshotRow{
				Address:    big.New(mrow.Row.Address.ToInt()),
				SlotId:     mrow.Row.SlotId,
				Version:    mrow.Row.Version,
				Expiration: mrow.Row.Expiration,
				Confirmed:  mrow.Row.Confirmed,
			})
		}
	}

	return rows, nil
}

func (o *inMemoryOrm) GetUnconfirmedRows(ctx context.Context, limit uint) ([]*Row, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	now := time.Now().UnixMilli()
	var mrows []*mrow
	for _, mrow := range o.rows {
		if mrow.Row.Expiration > now && !mrow.Row.Confirmed {
			mrows = append(mrows, mrow)
		}
	}

	sort.Slice(mrows, func(i, j int) bool {
		return mrows[i].UpdatedAt.Before(mrows[j].UpdatedAt)
	})

	if uint(len(mrows)) > limit {
		mrows = mrows[:limit]
	}

	rows := make([]*Row, len(mrows))
	for i, mrow := range mrows {
		rows[i] = mrow.Row.Clone()
	}

	return rows, nil
}
