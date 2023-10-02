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

const (
	SharedTableName  = "shared"
	s4PostgresSchema = "s4"
)

type orm struct {
	q         pg.Q
	tableName string
	namespace string
}

var _ ORM = (*orm)(nil)

func NewPostgresORM(db *sqlx.DB, lggr logger.Logger, cfg pg.QConfig, tableName, namespace string) ORM {
	return &orm{
		q:         pg.NewQ(db, lggr, cfg),
		tableName: fmt.Sprintf(`"%s".%s`, s4PostgresSchema, tableName),
		namespace: namespace,
	}
}

func (o orm) Get(address *utils.Big, slotId uint, qopts ...pg.QOpt) (*Row, error) {
	row := &Row{}
	q := o.q.WithOpts(qopts...)

	stmt := fmt.Sprintf(`SELECT address, slot_id, version, expiration, confirmed, payload, signature FROM %s 
WHERE namespace=$1 AND address=$2 AND slot_id=$3;`, o.tableName)
	if err := q.Get(row, stmt, o.namespace, address, slotId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrNotFound
		}
		return nil, err
	}
	return row, nil
}

func (o orm) Update(row *Row, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)

	// This query inserts or updates a row, depending on whether the version is higher than the existing one.
	// We only allow the same version when the row is confirmed.
	// We never transition back from unconfirmed to confirmed state.
	stmt := fmt.Sprintf(`INSERT INTO %s as t (namespace, address, slot_id, version, expiration, confirmed, payload, signature, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
ON CONFLICT (namespace, address, slot_id)
DO UPDATE SET version = EXCLUDED.version,
expiration = EXCLUDED.expiration,
confirmed = EXCLUDED.confirmed,
payload = EXCLUDED.payload,
signature = EXCLUDED.signature,
updated_at = NOW()
WHERE (t.version < EXCLUDED.version) OR (t.version <= EXCLUDED.version AND EXCLUDED.confirmed IS TRUE)
RETURNING id;`, o.tableName)
	var id uint64
	err := q.Get(&id, stmt, o.namespace, row.Address, row.SlotId, row.Version, row.Expiration, row.Confirmed, row.Payload, row.Signature)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrVersionTooLow
	}
	return err
}

func (o orm) DeleteExpired(limit uint, utcNow time.Time, qopts ...pg.QOpt) (int64, error) {
	q := o.q.WithOpts(qopts...)

	with := fmt.Sprintf(`WITH rows AS (SELECT id FROM %s WHERE namespace = $1 AND expiration < $2 LIMIT $3)`, o.tableName)
	stmt := fmt.Sprintf(`%s DELETE FROM %s WHERE id IN (SELECT id FROM rows);`, with, o.tableName)
	result, err := q.Exec(stmt, o.namespace, utcNow.UnixMilli(), limit)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (o orm) GetSnapshot(addressRange *AddressRange, qopts ...pg.QOpt) ([]*SnapshotRow, error) {
	q := o.q.WithOpts(qopts...)
	rows := make([]*SnapshotRow, 0)

	stmt := fmt.Sprintf(`SELECT address, slot_id, version, expiration, confirmed FROM %s WHERE namespace = $1 AND address >= $2 AND address <= $3;`, o.tableName)
	if err := q.Select(&rows, stmt, o.namespace, addressRange.MinAddress, addressRange.MaxAddress); err != nil {
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
WHERE namespace = $1 AND confirmed IS FALSE ORDER BY updated_at LIMIT $2;`, o.tableName)
	if err := q.Select(&rows, stmt, o.namespace, limit); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	return rows, nil
}
