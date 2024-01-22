package pg

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commonpg "github.com/smartcontractkit/chainlink-common/pkg/services/pg"
)

// NOTE: This is the default level in Postgres anyway, we just make it
// explicit here
const defaultIsolation = sql.LevelReadCommitted

// TxOption is a functional option for SQL transactions.
type TxOption = commonpg.TxOption

func OptReadOnlyTx() TxOption {
	return commonpg.OptReadOnlyTx()
}

func SqlTransaction(ctx context.Context, rdb *sql.DB, lggr logger.Logger, fn func(tx *sqlx.Tx) error, opts ...TxOption) (err error) {
	return commonpg.SqlTransaction(ctx, rdb, lggr, fn, opts...)
}

// TxBeginner can be a db or a conn, anything that implements BeginTxx
type TxBeginner = commonpg.TxBeginner

func sqlxTransactionQ(ctx context.Context, db TxBeginner, lggr logger.Logger, fn func(q Queryer) error, opts ...TxOption) (err error) {
	return commonpg.SqlxTransactionQ(ctx, db, lggr, fn, opts...)
}
