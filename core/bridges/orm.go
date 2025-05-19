package bridges

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	pkgerrors "github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/auth"
)

type ORM interface {
	FindBridge(ctx context.Context, name BridgeName) (bt BridgeType, err error)
	FindBridges(ctx context.Context, name []BridgeName) (bts []BridgeType, err error)
	DeleteBridgeType(ctx context.Context, bt *BridgeType) error
	BridgeTypes(ctx context.Context, offset int, limit int) ([]BridgeType, int, error)
	CreateBridgeType(ctx context.Context, bt *BridgeType) error
	UpdateBridgeType(ctx context.Context, bt *BridgeType, btr *BridgeTypeRequest) error

	GetCachedResponse(ctx context.Context, dotId string, specId int32, maxElapsed time.Duration) ([]byte, error)
	UpsertBridgeResponse(ctx context.Context, dotId string, specId int32, response []byte) error

	ExternalInitiators(ctx context.Context, offset int, limit int) ([]ExternalInitiator, int, error)
	CreateExternalInitiator(ctx context.Context, externalInitiator *ExternalInitiator) error
	DeleteExternalInitiator(ctx context.Context, name string) error
	FindExternalInitiator(ctx context.Context, eia *auth.Token) (*ExternalInitiator, error)
	FindExternalInitiatorByName(ctx context.Context, iname string) (exi ExternalInitiator, err error)

	GetCachedResponseWithFinished(ctx context.Context, dotId string, specId int32, maxElapsed time.Duration) ([]byte, time.Time, error)
	BulkUpsertBridgeResponse(ctx context.Context, responses []BridgeResponse) error

	WithDataSource(sqlutil.DataSource) ORM
}

type orm struct {
	ds sqlutil.DataSource
}

var _ ORM = (*orm)(nil)

func NewORM(ds sqlutil.DataSource) ORM {
	return &orm{ds: ds}
}

func (o *orm) WithDataSource(ds sqlutil.DataSource) ORM { return NewORM(ds) }

func (o *orm) transact(ctx context.Context, readOnly bool, fn func(tx *orm) error) error {
	opts := sqlutil.TxOptions{TxOptions: sql.TxOptions{ReadOnly: readOnly}}
	return sqlutil.Transact(ctx, func(ds sqlutil.DataSource) *orm { return &orm{ds: ds} }, o.ds, &opts, fn)
}

// FindBridge looks up a Bridge by its Name.
// Returns sql.ErrNoRows if name not present
func (o *orm) FindBridge(ctx context.Context, name BridgeName) (bt BridgeType, err error) {
	stmt := "SELECT * FROM bridge_types WHERE name = $1"
	err = o.ds.GetContext(ctx, &bt, stmt, name.String())

	return
}

// FindBridges looks up multiple bridges in a single query.
// Errors unless all bridges successfully found. Requires at least one bridge.
// Expects all bridges to be unique
func (o *orm) FindBridges(ctx context.Context, names []BridgeName) ([]BridgeType, error) {
	stmt := "SELECT * FROM bridge_types WHERE name IN (?)"
	query, args, err := sqlx.In(stmt, names)
	if err != nil {
		return nil, err
	}

	var bts []BridgeType

	if err = o.ds.SelectContext(ctx, &bts, o.ds.Rebind(query), args...); err != nil {
		return nil, err
	}

	if len(bts) != len(names) {
		return nil, pkgerrors.Errorf("not all bridges exist, asked for %v, exists %v", names, bts)
	}

	return bts, nil
}

// DeleteBridgeType removes the bridge type
func (o *orm) DeleteBridgeType(ctx context.Context, bt *BridgeType) error {
	query := "DELETE FROM bridge_types WHERE name = $1"
	result, err := o.ds.ExecContext(ctx, query, bt.Name)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// BridgeTypes returns bridge types ordered by name filtered limited by the
// passed params.
func (o *orm) BridgeTypes(ctx context.Context, offset int, limit int) (bridges []BridgeType, count int, err error) {
	err = o.transact(ctx, true, func(tx *orm) error {
		if err = tx.ds.GetContext(ctx, &count, "SELECT COUNT(*) FROM bridge_types"); err != nil {
			return pkgerrors.Wrap(err, "BridgeTypes failed to get count")
		}
		sql := `SELECT * FROM bridge_types ORDER BY name asc LIMIT $1 OFFSET $2;`
		if err = tx.ds.SelectContext(ctx, &bridges, sql, limit, offset); err != nil {
			return pkgerrors.Wrap(err, "BridgeTypes failed to load bridge_types")
		}
		return nil
	})

	return
}

// CreateBridgeType saves the bridge type.
func (o *orm) CreateBridgeType(ctx context.Context, bt *BridgeType) error {
	stmt := `INSERT INTO bridge_types (name, url, confirmations, incoming_token_hash, salt, outgoing_token, minimum_contract_payment, created_at, updated_at)
	VALUES (:name, :url, :confirmations, :incoming_token_hash, :salt, :outgoing_token, :minimum_contract_payment, now(), now())
	RETURNING *;`
	err := o.transact(ctx, false, func(tx *orm) error {
		stmt, err := tx.ds.PrepareNamedContext(ctx, stmt)
		if err != nil {
			return err
		}
		defer stmt.Close()
		return stmt.GetContext(ctx, bt, bt)
	})

	return pkgerrors.Wrap(err, "CreateBridgeType failed")
}

// UpdateBridgeType updates the bridge type.
func (o *orm) UpdateBridgeType(ctx context.Context, bt *BridgeType, btr *BridgeTypeRequest) error {
	stmt := "UPDATE bridge_types SET url = $1, confirmations = $2, minimum_contract_payment = $3 WHERE name = $4 RETURNING *"
	err := o.ds.GetContext(ctx, bt, stmt, btr.URL, btr.Confirmations, btr.MinimumContractPayment, bt.Name)

	return err
}

func (o *orm) GetCachedResponse(ctx context.Context, dotId string, specId int32, maxElapsed time.Duration) ([]byte, error) {
	response, _, err := o.GetCachedResponseWithFinished(ctx, dotId, specId, maxElapsed)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (o *orm) GetCachedResponseWithFinished(ctx context.Context, dotId string, specId int32, maxElapsed time.Duration) ([]byte, time.Time, error) {
	stalenessThreshold := time.Now().Add(-maxElapsed)
	sql := `SELECT value, finished_at FROM bridge_last_value WHERE
				dot_id = $1 AND 
				spec_id = $2 AND 
				finished_at > ($3)	
				ORDER BY finished_at 
				DESC LIMIT 1;`

	type responseType struct {
		Value      []byte
		FinishedAt time.Time
	}

	var result responseType

	if err := pkgerrors.Wrap(
		o.ds.GetContext(ctx, &result, sql, dotId, specId, stalenessThreshold),
		fmt.Sprintf("failed to fetch last good value for task %s spec %d", dotId, specId),
	); err != nil {
		return nil, time.Now(), err
	}

	return result.Value, result.FinishedAt, nil
}

func (o *orm) UpsertBridgeResponse(ctx context.Context, dotId string, specId int32, response []byte) error {
	sql := `INSERT INTO bridge_last_value(dot_id, spec_id, value, finished_at) 
				VALUES($1, $2, $3, $4)
			ON CONFLICT ON CONSTRAINT bridge_last_value_pkey
				DO UPDATE SET value = $3, finished_at = $4;`

	_, err := o.ds.ExecContext(ctx, sql, dotId, specId, response, time.Now())

	return err
}

func (o *orm) BulkUpsertBridgeResponse(ctx context.Context, responses []BridgeResponse) error {
	sql := `INSERT INTO bridge_last_value(dot_id, spec_id, value, finished_at) 
			VALUES (:dot_id, :spec_id, :value, :finished_at)
			ON CONFLICT ON CONSTRAINT bridge_last_value_pkey
				DO UPDATE SET value = excluded.value, finished_at = excluded.finished_at;`

	if _, err := o.ds.NamedExecContext(ctx, sql, responses); err != nil {
		return err
	}

	return nil
}

// --- External Initiator

// ExternalInitiators returns a list of external initiators sorted by name
func (o *orm) ExternalInitiators(ctx context.Context, offset int, limit int) (exis []ExternalInitiator, count int, err error) {
	err = o.transact(ctx, true, func(tx *orm) error {
		if err = tx.ds.GetContext(ctx, &count, "SELECT COUNT(*) FROM external_initiators"); err != nil {
			return pkgerrors.Wrap(err, "ExternalInitiators failed to get count")
		}

		sql := `SELECT * FROM external_initiators ORDER BY name asc LIMIT $1 OFFSET $2;`
		if err = tx.ds.SelectContext(ctx, &exis, sql, limit, offset); err != nil {
			return pkgerrors.Wrap(err, "ExternalInitiators failed to load external_initiators")
		}
		return nil
	})
	return
}

// CreateExternalInitiator inserts a new external initiator
func (o *orm) CreateExternalInitiator(ctx context.Context, externalInitiator *ExternalInitiator) (err error) {
	query := `INSERT INTO external_initiators (name, url, access_key, salt, hashed_secret, outgoing_secret, outgoing_token, created_at, updated_at)
	VALUES (:name, :url, :access_key, :salt, :hashed_secret, :outgoing_secret, :outgoing_token, now(), now())
	RETURNING *
	`
	err = o.transact(ctx, false, func(tx *orm) error {
		var stmt *sqlx.NamedStmt
		stmt, err = tx.ds.PrepareNamedContext(ctx, query)
		if err != nil {
			return pkgerrors.Wrap(err, "failed to prepare named stmt")
		}
		defer stmt.Close()
		return pkgerrors.Wrap(stmt.GetContext(ctx, externalInitiator, externalInitiator), "failed to load external_initiator")
	})
	return pkgerrors.Wrap(err, "CreateExternalInitiator failed")
}

// DeleteExternalInitiator removes an external initiator
func (o *orm) DeleteExternalInitiator(ctx context.Context, name string) error {
	query := "DELETE FROM external_initiators WHERE name = $1"
	result, err := o.ds.ExecContext(ctx, query, name)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return err
}

// FindExternalInitiator finds an external initiator given an authentication request
func (o *orm) FindExternalInitiator(ctx context.Context, eia *auth.Token) (*ExternalInitiator, error) {
	exi := &ExternalInitiator{}
	err := o.ds.GetContext(ctx, exi, `SELECT * FROM external_initiators WHERE access_key = $1`, eia.AccessKey)
	return exi, err
}

// FindExternalInitiatorByName finds an external initiator given an authentication request
func (o *orm) FindExternalInitiatorByName(ctx context.Context, iname string) (exi ExternalInitiator, err error) {
	err = o.ds.GetContext(ctx, &exi, `SELECT * FROM external_initiators WHERE lower(name) = lower($1)`, iname)
	return
}
