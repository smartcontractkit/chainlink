package forwarders

import (
	"database/sql"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

//go:generate mockery --quiet --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	CreateForwarder(addr common.Address, evmChainId utils.Big) (fwd Forwarder, err error)
	FindForwarders(offset, limit int) ([]Forwarder, int, error)
	FindForwardersByChain(evmChainId utils.Big) ([]Forwarder, error)
	DeleteForwarder(id int64, cleanup func(tx pg.Queryer, evmChainId int64, addr common.Address) error) error
	FindForwardersInListByChain(evmChainId utils.Big, addrs []common.Address) ([]Forwarder, error)
}

type orm struct {
	q pg.Q
}

var _ ORM = (*orm)(nil)

func NewORM(db *sqlx.DB, lggr logger.Logger, cfg pg.QConfig) *orm {
	return &orm{pg.NewQ(db, lggr, cfg)}
}

// CreateForwarder creates the Forwarder address associated with the current EVM chain id.
func (o *orm) CreateForwarder(addr common.Address, evmChainId utils.Big) (fwd Forwarder, err error) {
	sql := `INSERT INTO evm_forwarders (address, evm_chain_id, created_at, updated_at) VALUES ($1, $2, now(), now()) RETURNING *`
	err = o.q.Get(&fwd, sql, addr, evmChainId)
	return fwd, err
}

// DeleteForwarder removes a forwarder address.
// If cleanup is non-nil, it can be used to perform any chain- or contract-specific cleanup that need to happen atomically
// on forwarder deletion.  If cleanup returns an error, forwarder deletion will be aborted.
func (o *orm) DeleteForwarder(id int64, cleanup func(tx pg.Queryer, evmChainID int64, addr common.Address) error) (err error) {
	var dest struct {
		EvmChainId int64
		Address    common.Address
	}

	var rowsAffected int64
	err = o.q.Transaction(func(tx pg.Queryer) error {
		err = tx.Get(&dest, `SELECT evm_chain_id, address FROM evm_forwarders WHERE id = $1`, id)
		if err != nil {
			return err
		}
		if cleanup != nil {
			if err = cleanup(tx, dest.EvmChainId, dest.Address); err != nil {
				return err
			}
		}

		result, err2 := o.q.Exec(`DELETE FROM evm_forwarders WHERE id = $1`, id)
		// If the forwarder wasn't found, we still want to delete the filter.
		// In that case, the transaction must return nil, even though DeleteForwarder
		// will return sql.ErrNoRows
		if err2 != nil && !errors.Is(err2, sql.ErrNoRows) {
			return err2
		}
		rowsAffected, err2 = result.RowsAffected()

		return err2
	})

	if err == nil && rowsAffected == 0 {
		err = sql.ErrNoRows
	}
	return err
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

// FindForwardersByChain returns all forwarder addresses for a chain.
func (o *orm) FindForwardersByChain(evmChainId utils.Big) (fwds []Forwarder, err error) {
	sql := `SELECT * FROM evm_forwarders where evm_chain_id = $1 ORDER BY created_at DESC, id DESC`
	err = o.q.Select(&fwds, sql, evmChainId)
	return
}

func (o *orm) FindForwardersInListByChain(evmChainId utils.Big, addrs []common.Address) ([]Forwarder, error) {
	var fwdrs []Forwarder

	arg := map[string]interface{}{
		"addresses": addrs,
		"chainid":   evmChainId,
	}

	query, args, err := sqlx.Named(`
		SELECT * FROM evm_forwarders 
		WHERE evm_chain_id = :chainid
		AND address IN (:addresses)
		ORDER BY created_at DESC, id DESC`,
		arg,
	)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to format query")
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to run sqlx.IN on query")
	}

	query = o.q.Rebind(query)
	err = o.q.Select(&fwdrs, query, args...)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to execute query")
	}

	return fwdrs, nil
}
