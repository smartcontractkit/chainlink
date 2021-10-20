package bridges

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/sqlx"
)

type ORM interface {
	FindBridge(name TaskType) (bt BridgeType, err error)
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
	db *sqlx.DB
}

var _ ORM = (*orm)(nil)

func NewORM(db *sqlx.DB) ORM {
	return &orm{db}
}

// FindBridge looks up a Bridge by its Name.
func (o *orm) FindBridge(name TaskType) (bt BridgeType, err error) {
	sql := "SELECT * FROM bridge_types WHERE name = $1"
	err = o.db.Get(&bt, sql, name.String())
	return
}

// DeleteBridgeType removes the bridge type
func (o *orm) DeleteBridgeType(bt *BridgeType) error {
	query := "DELETE FROM bridge_types WHERE name = $1"
	result, err := o.db.Exec(query, bt.Name)
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
	if err = o.db.Get(&count, "SELECT COUNT(*) FROM bridge_types"); err != nil {
		return
	}

	sql := `SELECT * FROM bridge_types ORDER BY name asc LIMIT $1 OFFSET $2;`
	if err = o.db.Select(&bridges, sql, limit, offset); err != nil {
		return
	}

	return
}

// CreateBridgeType saves the bridge type.
func (o *orm) CreateBridgeType(bt *BridgeType) error {
	sql := `INSERT INTO bridge_types (name, url, confirmations, incoming_token_hash, salt, outgoing_token, minimum_contract_payment, created_at, updated_at)
	VALUES (:name, :url, :confirmations, :incoming_token_hash, :salt, :outgoing_token, :minimum_contract_payment, now(), now())
	RETURNING *;`
	stmt, err := o.db.PrepareNamed(sql)
	if err != nil {
		return err
	}
	return stmt.Get(bt, bt)
}

// UpdateBridgeType updates the bridge type.
func (o *orm) UpdateBridgeType(bt *BridgeType, btr *BridgeTypeRequest) error {
	sql := "UPDATE bridge_types SET url = $1, confirmations = $2, minimum_contract_payment = $3 WHERE name = $4 RETURNING *"
	return o.db.Get(bt, sql, btr.URL, btr.Confirmations, btr.MinimumContractPayment, bt.Name)
}

// --- External Initiator

// ExternalInitiators returns a list of external initiators sorted by name
func (o *orm) ExternalInitiators(offset int, limit int) (exis []ExternalInitiator, count int, err error) {
	if err = o.db.Get(&count, "SELECT COUNT(*) FROM external_initiators"); err != nil {
		return
	}

	sql := `SELECT * FROM external_initiators ORDER BY name asc LIMIT $1 OFFSET $2;`
	if err = o.db.Select(&exis, sql, limit, offset); err != nil {
		return
	}
	return
}

// CreateExternalInitiator inserts a new external initiator
func (o *orm) CreateExternalInitiator(externalInitiator *ExternalInitiator) error {
	sql := `INSERT INTO external_initiators (name, url, access_key, salt, hashed_secret, outgoing_secret, outgoing_token, created_at, updated_at)
	VALUES (:name, :url, :access_key, :salt, :hashed_secret, :outgoing_secret, :outgoing_token, now(), now())
	RETURNING *
	`
	stmt, err := o.db.PrepareNamed(sql)
	if err != nil {
		return errors.Wrap(err, "CreateExternalInitiator failed")
	}
	return stmt.Get(externalInitiator, externalInitiator)
}

// DeleteExternalInitiator removes an external initiator
func (o *orm) DeleteExternalInitiator(name string) error {
	query := "DELETE FROM external_initiators WHERE name = $1"
	result, err := o.db.Exec(query, name)
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
	err := o.db.Get(exi, `SELECT * FROM external_initiators WHERE access_key = $1`, eia.AccessKey)
	return exi, err
}

// FindExternalInitiatorByName finds an external initiator given an authentication request
func (o *orm) FindExternalInitiatorByName(iname string) (exi ExternalInitiator, err error) {
	err = o.db.Get(&exi, `SELECT * FROM external_initiators WHERE lower(name) = lower($1)`, iname)
	return
}
