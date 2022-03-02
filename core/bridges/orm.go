package bridges

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/sqlx"
)

//go:generate mockery --name ORM --output ./mocks --case=underscore

type ORM interface {
	FindBridge(name BridgeName) (bt BridgeType, err error)
	FindBridges(name []BridgeName) (bts []BridgeType, err error)
	DeleteBridgeType(bt *BridgeType) error
	BridgeTypes(offset int, limit int) ([]BridgeType, int, error)
	CreateBridgeType(bt *BridgeType) error
	UpdateBridgeType(bt *BridgeType, btr *BridgeTypeRequest) error

	ExternalInitiators(offset int, limit int) ([]ExternalInitiator, int, error)
	CreateExternalInitiator(externalInitiator *ExternalInitiator) error
	DeleteExternalInitiator(name string) error
	FindExternalInitiator(eia *auth.Token) (*ExternalInitiator, error)
	FindExternalInitiatorByName(iname string) (exi ExternalInitiator, err error)
}

type orm struct {
	q pg.Q
}

var _ ORM = (*orm)(nil)

func NewORM(db *sqlx.DB, lggr logger.Logger, cfg pg.LogConfig) ORM {
	namedLogger := lggr.Named("BridgeORM")
	return &orm{pg.NewQ(db, namedLogger, cfg)}
}

// FindBridge looks up a Bridge by its Name.
// Returns sql.ErrNoRows if name not present
func (o *orm) FindBridge(name BridgeName) (bt BridgeType, err error) {
	sql := "SELECT * FROM bridge_types WHERE name = $1"
	err = o.q.Get(&bt, sql, name.String())
	return
}

// FindBridges looks up multiple bridges in a single query.
// Errors unless all bridges successfully found. Requires at least one bridge.
// Expects all bridges to be unique
func (o *orm) FindBridges(names []BridgeName) (bts []BridgeType, err error) {
	sql := "SELECT * FROM bridge_types WHERE name IN (?)"
	query, args, err := sqlx.In(sql, names)
	if err != nil {
		return nil, err
	}
	err = o.q.Select(&bts, o.q.Rebind(query), args...)
	if err != nil {
		return nil, err
	}
	if len(bts) != len(names) {
		return nil, errors.Errorf("not all bridges exist, asked for %v, exists %v", names, bts)
	}
	return
}

// DeleteBridgeType removes the bridge type
func (o *orm) DeleteBridgeType(bt *BridgeType) error {
	query := "DELETE FROM bridge_types WHERE name = $1"
	result, err := o.q.Exec(query, bt.Name)
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

// BridgeTypes returns bridge types ordered by name filtered limited by the
// passed params.
func (o *orm) BridgeTypes(offset int, limit int) (bridges []BridgeType, count int, err error) {
	err = o.q.Transaction(func(tx pg.Queryer) error {
		if err = tx.Get(&count, "SELECT COUNT(*) FROM bridge_types"); err != nil {
			return errors.Wrap(err, "BridgeTypes failed to get count")
		}
		sql := `SELECT * FROM bridge_types ORDER BY name asc LIMIT $1 OFFSET $2;`
		if err = tx.Select(&bridges, sql, limit, offset); err != nil {
			return errors.Wrap(err, "BridgeTypes failed to load bridge_types")
		}
		return nil
	}, pg.OptReadOnlyTx())

	return
}

// CreateBridgeType saves the bridge type.
func (o *orm) CreateBridgeType(bt *BridgeType) error {
	stmt := `INSERT INTO bridge_types (name, url, confirmations, incoming_token_hash, salt, outgoing_token, minimum_contract_payment, created_at, updated_at)
	VALUES (:name, :url, :confirmations, :incoming_token_hash, :salt, :outgoing_token, :minimum_contract_payment, now(), now())
	RETURNING *;`
	err := o.q.Transaction(func(tx pg.Queryer) error {
		stmt, err := tx.PrepareNamed(stmt)
		if err != nil {
			return err
		}
		return stmt.Get(bt, bt)
	})
	return errors.Wrap(err, "CreateBridgeType failed")
}

// UpdateBridgeType updates the bridge type.
func (o *orm) UpdateBridgeType(bt *BridgeType,
	btr *BridgeTypeRequest) error {
	sql := "UPDATE bridge_types SET url = $1, confirmations = $2, minimum_contract_payment = $3 WHERE name = $4 RETURNING *"
	return o.q.Get(bt, sql, btr.URL, btr.Confirmations, btr.MinimumContractPayment, bt.Name)
}

// --- External Initiator

// ExternalInitiators returns a list of external initiators sorted by name
func (o *orm) ExternalInitiators(offset int, limit int) (exis []ExternalInitiator, count int, err error) {
	err = o.q.Transaction(func(tx pg.Queryer) error {
		if err = tx.Get(&count, "SELECT COUNT(*) FROM external_initiators"); err != nil {
			return errors.Wrap(err, "ExternalInitiators failed to get count")
		}

		sql := `SELECT * FROM external_initiators ORDER BY name asc LIMIT $1 OFFSET $2;`
		if err = tx.Select(&exis, sql, limit, offset); err != nil {
			return errors.Wrap(err, "ExternalInitiators failed to load external_initiators")
		}
		return nil
	}, pg.OptReadOnlyTx())
	return
}

// CreateExternalInitiator inserts a new external initiator
func (o *orm) CreateExternalInitiator(externalInitiator *ExternalInitiator) (err error) {
	query := `INSERT INTO external_initiators (name, url, access_key, salt, hashed_secret, outgoing_secret, outgoing_token, created_at, updated_at)
	VALUES (:name, :url, :access_key, :salt, :hashed_secret, :outgoing_secret, :outgoing_token, now(), now())
	RETURNING *
	`
	err = o.q.Transaction(func(tx pg.Queryer) error {
		var stmt *sqlx.NamedStmt
		stmt, err = tx.PrepareNamed(query)
		if err != nil {
			return errors.Wrap(err, "failed to prepare named stmt")
		}
		return errors.Wrap(stmt.Get(externalInitiator, externalInitiator), "failed to load external_initiator")
	})
	return errors.Wrap(err, "CreateExternalInitiator failed")
}

// DeleteExternalInitiator removes an external initiator
func (o *orm) DeleteExternalInitiator(name string) error {
	query := "DELETE FROM external_initiators WHERE name = $1"
	ctx, cancel := o.q.Context()
	defer cancel()
	result, err := o.q.ExecContext(ctx, query, name)
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
func (o *orm) FindExternalInitiator(
	eia *auth.Token,
) (*ExternalInitiator, error) {
	exi := &ExternalInitiator{}
	err := o.q.Get(exi, `SELECT * FROM external_initiators WHERE access_key = $1`, eia.AccessKey)
	return exi, err
}

// FindExternalInitiatorByName finds an external initiator given an authentication request
func (o *orm) FindExternalInitiatorByName(iname string) (exi ExternalInitiator, err error) {
	err = o.q.Get(&exi, `SELECT * FROM external_initiators WHERE lower(name) = lower($1)`, iname)
	return
}
