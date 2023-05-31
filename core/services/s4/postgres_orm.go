package s4

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"
)

type postgres struct {
	q         pg.Q
	tableName string
}

var _ ORM = (*postgres)(nil)

const fields = "address, slot_id, version, expiration, confirmed, payload, signature, updated_at"

func NewPostgresORM(db *sqlx.DB, lggr logger.Logger, cfg pg.QConfig, tableName string) ORM {
	return &postgres{
		q:         pg.NewQ(db, lggr, cfg),
		tableName: tableName,
	}
}

func (p postgres) Get(address *utils.Big, slotId uint, qopts ...pg.QOpt) (*Row, error) {
	row := &Row{}
	q := p.q.WithOpts(qopts...)

	stmt := fmt.Sprintf(`SELECT %s FROM %s WHERE address=$1 AND slot_id=$2;`, fields, p.tableName)
	if err := q.Get(&row, stmt, address, slotId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrNotFound
		}
		return nil, err
	}
	return row, nil
}

func (p postgres) Update(row *Row, qopts ...pg.QOpt) error {
	q := p.q.WithOpts(qopts...)

	stmt := fmt.Sprintf(`INSERT INTO %s as t (%s)
	                     VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
					     ON CONFLICT (t.address, t.slot_id)
						 DO UPDATE SET t.version = EXCLUDED.version 
									   t.expiration = EXCLUDED.expiration
									   t.confirmed = EXCLUDED.confirmed
									   t.payload = EXCLUDED.payload
									   t.signature = EXCLUDED.signature
									   t.updated_at = NOW()
                         WHERE t.version < EXCLUDED.version;`,
		fields, p.tableName)
	return q.ExecQ(stmt, row.Address, row.SlotId, row.Version)
}

func (p postgres) DeleteExpired(limit uint, utcNow time.Time, qopts ...pg.QOpt) error {
	q := p.q.WithOpts(qopts...)

	with := fmt.Sprintf(`WITH rows AS (SELECT id FROM %s WHERE expiration < $1 LIMIT $2)`, p.tableName)
	stmt := fmt.Sprintf(`%s DELETE FROM %s WHERE id IN (SELECT id FROM rows);`, with, p.tableName)
	return q.ExecQ(stmt, utcNow, limit)
}

func (p postgres) GetSnapshot(addressRange *AddressRange, qopts ...pg.QOpt) ([]*SnapshotRow, error) {
	q := p.q.WithOpts(qopts...)
	rows := make([]*SnapshotRow, 0)

	stmt := fmt.Sprintf(`SELECT address, slot_id, version FROM %s WHERE address >= $1 AND address <= $2;`, p.tableName)
	if err := q.Select(rows, stmt, addressRange.MinAddress, addressRange.MaxAddress); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	return rows, nil
}

func (p postgres) GetUnconfirmedRows(limit uint, qopts ...pg.QOpt) ([]*Row, error) {
	q := p.q.WithOpts(qopts...)
	rows := make([]*Row, 0)

	stmt := fmt.Sprintf(`SELECT %s FROM %s WHERE confirmed IS FALSE ORDER BY updated_at LIMIT $1;`, fields, p.tableName)
	if err := q.Select(rows, stmt, limit); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	return rows, nil
}
