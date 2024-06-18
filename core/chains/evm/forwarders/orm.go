package forwarders

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"

	"github.com/ethereum/go-ethereum/common"
	pkgerrors "github.com/pkg/errors"
)

//go:generate mockery --quiet --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	CreateForwarder(ctx context.Context, addr common.Address, evmChainId big.Big) (fwd Forwarder, err error)
	//FindForwarders(ctx context.Context, offset, limit int) ([]Forwarder, int, error)
	FindForwardersByChain(ctx context.Context, evmChainId big.Big) ([]Forwarder, error)
	DeleteForwarder(ctx context.Context, id int64, cleanup func(tx sqlutil.DataSource, evmChainId int64, addr common.Address) error) error
	// FindForwardersInListByChain(ctx context.Context, evmChainId big.Big, addrs []common.Address) ([]Forwarder, error)
}

type DSORM struct {
	ds  sqlutil.DataSource
	cid *big.Big
}

var _ ORM = &DSORM{}

func NewScopedORM(ds sqlutil.DataSource, evmChainId *big.Big) *DSORM {
	return &DSORM{ds: ds, cid: evmChainId}
}

func (o *DSORM) Transact(ctx context.Context, fn func(*DSORM) error) (err error) {
	return sqlutil.Transact(ctx, o.new, o.ds, nil, fn)
}

// new returns a NewORM like o, but backed by q.
func (o *DSORM) new(q sqlutil.DataSource) *DSORM { return NewScopedORM(q, o.cid) }

func (o *DSORM) schemaName() string {
	if o.cid != nil {
		return fmt.Sprintf("evm_%s", o.cid.String())
	}
	return "evm"
}

// CreateForwarder creates the Forwarder address associated with the current EVM chain id.
func (o *DSORM) CreateForwarder(ctx context.Context, addr common.Address, evmChainId big.Big) (fwd Forwarder, err error) {
	if o.cid != nil && !o.cid.Equal(&evmChainId) {
		// hacking
		evmChainId = *o.cid
	}
	//	sql := `INSERT INTO evm.forwarders (address, evm_chain_id, created_at, updated_at) VALUES ($1, $2, now(), now()) RETURNING *`
	sql := fmt.Sprintf("INSERT INTO %s.forwarders (address, created_at, updated_at) VALUES ($1, now(), now()) RETURNING *", o.schemaName())
	err = o.ds.GetContext(ctx, &fwd, sql, addr)
	if err != nil {
		return fwd, err
	}
	fwd.EVMChainID = *o.cid
	return fwd, nil
}

// DeleteForwarder removes a forwarder address.
// If cleanup is non-nil, it can be used to perform any chain- or contract-specific cleanup that need to happen atomically
// on forwarder deletion.  If cleanup returns an error, forwarder deletion will be aborted.
func (o *DSORM) DeleteForwarder(ctx context.Context, id int64, cleanup func(tx sqlutil.DataSource, evmChainID int64, addr common.Address) error) (err error) {
	return o.Transact(ctx, func(orm *DSORM) error {
		var dest struct {
			EvmChainId int64
			Address    common.Address
		}
		selectStmt := fmt.Sprintf("SELECT  address FROM %s.forwarders WHERE id = $1", o.schemaName())
		err := orm.ds.GetContext(ctx, &dest, selectStmt, id)
		if err != nil {
			return err
		}
		if cleanup != nil {
			if err = cleanup(orm.ds, dest.EvmChainId, dest.Address); err != nil {
				return err
			}
		}
		deleteStmt := fmt.Sprintf("DELETE FROM %s.forwarders WHERE id = $1", o.schemaName())
		result, err := orm.ds.ExecContext(ctx, deleteStmt, id)
		// If the forwarder wasn't found, we still want to delete the filter.
		// In that case, the transaction must return nil, even though DeleteForwarder
		// will return sql.ErrNoRows
		if err != nil && !pkgerrors.Is(err, sql.ErrNoRows) {
			return err
		}
		rowsAffected, err := result.RowsAffected()
		if err == nil && rowsAffected == 0 {
			err = sql.ErrNoRows
		}
		return err
	})
}

// FindForwarders returns all forwarder addresses from offset up until limit.
func (o *DSORM) FindForwarders(ctx context.Context, offset, limit int) (fwds []Forwarder, count int, err error) {
	//	sql := `SELECT count(*) FROM evm.forwarders`
	sql := fmt.Sprintf("SELECT count(*) FROM %s.forwarders", o.schemaName())

	if err = o.ds.GetContext(ctx, &count, sql); err != nil {
		return
	}

	//	sql = `SELECT * FROM evm.forwarders ORDER BY created_at DESC, id DESC LIMIT $1 OFFSET $2`
	sql = fmt.Sprintf("SELECT * FROM %s.forwarders ORDER BY created_at DESC, id DESC LIMIT $1 OFFSET $2", o.schemaName())
	if err = o.ds.SelectContext(ctx, &fwds, sql, limit, offset); err != nil {
		return
	}
	return
}

// FindForwardersByChain returns all forwarder addresses for a chain.
func (o *DSORM) FindForwardersByChain(ctx context.Context, evmChainId big.Big) (fwds []Forwarder, err error) {
	//	sql := `SELECT * FROM evm.forwarders where evm_chain_id = $1 ORDER BY created_at DESC, id DESC`
	sql := fmt.Sprintf("SELECT * FROM %s.forwarders ORDER BY created_at DESC, id DESC", o.schemaName())
	err = o.ds.SelectContext(ctx, &fwds, sql)
	return
}

/*
func (o *DSORM) FindForwardersInListByChain(ctx context.Context, evmChainId big.Big, addrs []common.Address) ([]Forwarder, error) {
	var fwdrs []Forwarder

	arg := map[string]interface{}{
		"addresses": addrs,
		"chainid":   evmChainId,
	}

	query, args, err := sqlx.Named(`
		SELECT * FROM evm.forwarders
		WHERE evm_chain_id = :chainid
		AND address IN (:addresses)
		ORDER BY created_at DESC, id DESC`,
		arg,
	)

	if err != nil {
		return nil, pkgerrors.Wrap(err, "Failed to format query")
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "Failed to run sqlx.IN on query")
	}

	query = o.ds.Rebind(query)
	err = o.ds.SelectContext(ctx, &fwdrs, query, args...)

	if err != nil {
		return nil, pkgerrors.Wrap(err, "Failed to execute query")
	}

	return fwdrs, nil
}
*/
