package s4

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
)

const (
	SharedTableName  = "shared"
	s4PostgresSchema = "s4"
)

type orm struct {
	ds        sqlutil.DataSource
	tableName string
	namespace string
}

var _ ORM = (*orm)(nil)

func NewPostgresORM(ds sqlutil.DataSource, tableName, namespace string) ORM {
	return &orm{
		ds:        ds,
		tableName: fmt.Sprintf(`"%s".%s`, s4PostgresSchema, tableName),
		namespace: namespace,
	}
}

func (o *orm) Get(ctx context.Context, address *big.Big, slotId uint) (*Row, error) {
	row := &Row{}

	stmt := fmt.Sprintf(`SELECT address, slot_id, version, expiration, confirmed, payload, signature FROM %s 
WHERE namespace=$1 AND address=$2 AND slot_id=$3;`, o.tableName)
	if err := o.ds.GetContext(ctx, row, stmt, o.namespace, address, slotId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrNotFound
		}
		return nil, err
	}
	return row, nil
}

func (o *orm) Update(ctx context.Context, row *Row) error {
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
	err := o.ds.GetContext(ctx, &id, stmt, o.namespace, row.Address, row.SlotId, row.Version, row.Expiration, row.Confirmed, row.Payload, row.Signature)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrVersionTooLow
	}
	return err
}

func (o *orm) DeleteExpired(ctx context.Context, limit uint, utcNow time.Time) (int64, error) {
	with := fmt.Sprintf(`WITH rows AS (SELECT id FROM %s WHERE namespace = $1 AND expiration < $2 LIMIT $3)`, o.tableName)
	stmt := fmt.Sprintf(`%s DELETE FROM %s WHERE id IN (SELECT id FROM rows);`, with, o.tableName)
	result, err := o.ds.ExecContext(ctx, stmt, o.namespace, utcNow.UnixMilli(), limit)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (o *orm) GetSnapshot(ctx context.Context, addressRange *AddressRange) ([]*SnapshotRow, error) {
	rows := make([]*SnapshotRow, 0)

	stmt := fmt.Sprintf(`SELECT address, slot_id, version, expiration, confirmed, octet_length(payload) AS payload_size FROM %s WHERE namespace = $1 AND address >= $2 AND address <= $3;`, o.tableName)
	if err := o.ds.SelectContext(ctx, &rows, stmt, o.namespace, addressRange.MinAddress, addressRange.MaxAddress); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	return rows, nil
}

func (o *orm) GetUnconfirmedRows(ctx context.Context, limit uint) ([]*Row, error) {
	rows := make([]*Row, 0)

	stmt := fmt.Sprintf(`SELECT address, slot_id, version, expiration, confirmed, payload, signature FROM %s
WHERE namespace = $1 AND confirmed IS FALSE ORDER BY updated_at LIMIT $2;`, o.tableName)
	if err := o.ds.SelectContext(ctx, &rows, stmt, o.namespace, limit); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	return rows, nil
}
