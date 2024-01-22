package pg

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commonpg "github.com/smartcontractkit/chainlink-common/pkg/services/pg"
)

type Queryer = commonpg.Queryer

func WrapDbWithSqlx(rdb *sql.DB) *sqlx.DB {
	return commonpg.WrapDbWithSqlx(rdb)
}

func SqlxTransaction(ctx context.Context, q Queryer, lggr logger.Logger, fc func(q Queryer) error, txOpts ...TxOption) (err error) {
	return commonpg.SqlxTransaction(ctx, q, lggr, fc, txOpts...)
}
