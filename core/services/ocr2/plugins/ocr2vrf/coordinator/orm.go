package coordinator

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"

	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//go:generate mockery --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	// HeadsByNumbers fetches the heads with the given numbers from the db, returns nil if none exist
	HeadsByNumbers(ctx context.Context, numbers []uint64) ([]*evmtypes.Head, error)
}

var _ ORM = &orm{}

type orm struct {
	q       pg.Q
	chainID utils.Big
}

func NewORM(db *sqlx.DB, chainID utils.Big, lggr logger.Logger) ORM {
	return &orm{
		q:       pg.NewQ(db, lggr.Named("OCR2VRF_ORM"), sqlConfig{}),
		chainID: chainID,
	}
}

func (o *orm) HeadsByNumbers(ctx context.Context, numbers []uint64) (heads []*evmtypes.Head, err error) {
	q := o.q.WithOpts(pg.WithParentCtx(ctx))
	a := map[string]any{
		"chainid": o.chainID,
		"numbers": numbers,
	}
	query, args, err := sqlx.Named(`SELECT * FROM evm_heads WHERE evm_chain_id = :chainid AND number IN (:numbers)`, a)
	if err != nil {
		return nil, errors.Wrap(err, "sqlx Named")
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "sqlx In")
	}

	query = q.Rebind(query)
	err = q.Select(&heads, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return heads, err
}
