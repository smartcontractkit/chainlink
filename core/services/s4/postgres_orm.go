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

const postgresSchemaName = "s4"

type orm struct {
	q         pg.Q
	tableName string
}

var _ ORM = (*orm)(nil)

func NewPostgresORM(db *sqlx.DB, lggr logger.Logger, cfg pg.QConfig, tableName string) ORM {
	return &orm{
		q:         pg.NewQ(db, lggr, cfg),
		tableName: fmt.Sprintf(`"%s".%s`, postgresSchemaName, tableName),
	}
}

func (o orm) Get(address *utils.Big, slotId uint, qopts ...pg.QOpt) (*Row, error) {
	row := &Row{}
	q := o.q.WithOpts(qopts...)

	stmt := fmt.Sprintf(`SELECT address, slot_id, version, expiration, confirmed, payload, signature FROM %s 
WHERE address=$1 AND slot_id=$2;`, o.tableName)
	if err := q.Get(row, stmt, address, slotId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrNotFound
		}
		return nil, err
	}
	return row, nil
}

func (o orm) Update(row *Row, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)

	stmt := fmt.Sprintf(`INSERT INTO %s as t (address, slot_id, version, expiration, confirmed, payload, signature, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
ON CONFLICT (address, slot_id)
DO UPDATE SET version = EXCLUDED.version,
expiration = EXCLUDED.expiration,
confirmed = EXCLUDED.confirmed,
payload = EXCLUDED.payload,
signature = EXCLUDED.signature,
updated_at = NOW()
WHERE t.version < EXCLUDED.version
RETURNING id;`, o.tableName)
	var id uint64
	err := q.Get(&id, stmt, row.Address, row.SlotId, row.Version, row.Expiration, row.Confirmed, row.Payload, row.Signature)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrVersionTooLow
	}
	return nil
}

func (o orm) DeleteExpired(limit uint, utcNow time.Time, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)

	with := fmt.Sprintf(`WITH rows AS (SELECT id FROM %s WHERE expiration < $1 LIMIT $2)`, o.tableName)
	stmt := fmt.Sprintf(`%s DELETE FROM %s WHERE id IN (SELECT id FROM rows);`, with, o.tableName)
	return q.ExecQ(stmt, utcNow.UnixMilli(), limit)
}

func (o orm) GetSnapshot(addressRange *AddressRange, qopts ...pg.QOpt) ([]*SnapshotRow, error) {
	q := o.q.WithOpts(qopts...)
	rows := make([]*SnapshotRow, 0)

	stmt := fmt.Sprintf(`SELECT address, slot_id, version FROM %s WHERE address >= $1 AND address <= $2;`, o.tableName)
	if err := q.Select(&rows, stmt, addressRange.MinAddress, addressRange.MaxAddress); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	return rows, nil
}

func (o orm) GetUnconfirmedRows(limit uint, qopts ...pg.QOpt) ([]*Row, error) {
	q := o.q.WithOpts(qopts...)
	rows := make([]*Row, 0)

	stmt := fmt.Sprintf(`SELECT address, slot_id, version, expiration, confirmed, payload, signature FROM %s
WHERE confirmed IS FALSE ORDER BY updated_at LIMIT $1;`, o.tableName)
	if err := q.Select(&rows, stmt, limit); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	return rows, nil
}
