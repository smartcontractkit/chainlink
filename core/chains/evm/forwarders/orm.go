package forwarders

import (
	"database/sql"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//go:generate mockery --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	CreateForwarder(addr common.Address, evmChainId utils.Big) (fwd Forwarder, err error)
	FindForwarders(offset, limit int) ([]Forwarder, int, error)
	DeleteForwarder(id int32) error
}

type orm struct {
	q pg.Q
}

var _ ORM = (*orm)(nil)

func NewORM(db *sqlx.DB, lggr logger.Logger, cfg pg.LogConfig) *orm {
	return &orm{pg.NewQ(db, lggr, cfg)}
}

// CreateForwarder creates the Forwarder address associated with the current EVM chain id.
func (o *orm) CreateForwarder(addr common.Address, evmChainId utils.Big) (fwd Forwarder, err error) {
	sql := `INSERT INTO evm_forwarders (address, evm_chain_id, created_at, updated_at) VALUES ($1, $2, now(), now()) RETURNING *`
	err = o.q.Get(&fwd, sql, addr, evmChainId)
	return fwd, err
}

// DeleteForwarder removes a forwarder address.
func (o *orm) DeleteForwarder(id int32) error {
	q := `DELETE FROM evm_forwarders WHERE id = $1`
	result, err := o.q.Exec(q, id)
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

// FindForwarders returns all forwarder addresses from offset up until limit.
func (o *orm) FindForwarders(offset, limit int) (fwds []Forwarder, count int, err error) {
	sql := `SELECT count(*) FROM evm_forwarders`
	if err = o.q.Get(&count, sql); err != nil {
		return
	}

	sql = `SELECT * FROM evm_forwarders ORDER BY created_at DESC, id DESC LIMIT $1 OFFSET $2`
	if err = o.q.Select(&fwds, sql, limit, offset); err != nil {
		return
	}
	return
}
