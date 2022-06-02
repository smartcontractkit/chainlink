package operators

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
	CreateOperator(addr common.Address, evmChainId utils.Big) (fwd Operator, err error)
	FindOperators(offset, limit int) ([]Operator, int, error)
	FindOperatorsByChain(evmChainId utils.Big) ([]Operator, error)
	DeleteOperator(id int32) error
}

type orm struct {
	q pg.Q
}

var _ ORM = (*orm)(nil)

func NewORM(db *sqlx.DB, lggr logger.Logger, cfg pg.LogConfig) *orm {
	return &orm{pg.NewQ(db, lggr, cfg)}
}

// CreateOperator creates the Operator address associated with the current chain id.
func (o *orm) CreateOperator(addr common.Address, chainId utils.Big) (opr Operator, err error) {
	sql := `INSERT INTO operators (address, chain_id, created_at, updated_at) VALUES ($1, $2, now(), now()) RETURNING *`
	err = o.q.Get(&opr, sql, addr, chainId)
	return opr, err
}

// DeleteOperator removes a Operator address.
func (o *orm) DeleteOperator(id int32) error {
	q := `DELETE FROM operators WHERE id = $1`
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

// FindOperators returns all Operator addresses from offset up until limit.
func (o *orm) FindOperators(offset, limit int) (oprs []Operator, count int, err error) {
	sql := `SELECT count(*) FROM operators`
	if err = o.q.Get(&count, sql); err != nil {
		return
	}

	sql = `SELECT * FROM operators ORDER BY created_at DESC, id DESC LIMIT $1 OFFSET $2`
	if err = o.q.Select(&oprs, sql, limit, offset); err != nil {
		return
	}
	return
}

// FindOperatorsByChain returns all Operator addresses for a chain.
func (o *orm) FindOperatorsByChain(evmChainId utils.Big) (oprs []Operator, err error) {
	sql := `SELECT * FROM operators where chain_id = $1 ORDER BY created_at DESC, id DESC`
	err = o.q.Select(&oprs, sql, evmChainId)
	return
}
