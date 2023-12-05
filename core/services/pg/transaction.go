package pg

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	corelogger "github.com/smartcontractkit/chainlink/v2/core/logger"
)

// NOTE: This is the default level in Postgres anyway, we just make it
// explicit here
const defaultIsolation = sql.LevelReadCommitted

// TxOption is a functional option for SQL transactions.
type TxOption func(*sql.TxOptions)

func OptReadOnlyTx() TxOption {
	return func(opts *sql.TxOptions) {
		opts.ReadOnly = true
	}
}

func SqlTransaction(ctx context.Context, rdb *sql.DB, lggr logger.Logger, fn func(tx *sqlx.Tx) error, opts ...TxOption) (err error) {
	db := WrapDbWithSqlx(rdb)
	wrapFn := func(q Queryer) error {
		tx, ok := q.(*sqlx.Tx)
		if !ok {
			panic(fmt.Sprintf("expected q to be %T but got %T", tx, q))
		}
		return fn(tx)
	}
	return sqlxTransactionQ(ctx, db, lggr, wrapFn, opts...)
}

// txBeginner can be a db or a conn, anything that implements BeginTxx
type txBeginner interface {
	BeginTxx(context.Context, *sql.TxOptions) (*sqlx.Tx, error)
}

func sqlxTransactionQ(ctx context.Context, db txBeginner, lggr logger.Logger, fn func(q Queryer) error, opts ...TxOption) (err error) {
	var txOpts sql.TxOptions
	for _, o := range opts {
		o(&txOpts)
	}

	var tx *sqlx.Tx
	tx, err = db.BeginTxx(ctx, &txOpts)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}

	defer func() {
		if p := recover(); p != nil {
			sentry.CurrentHub().Recover(p)
			sentry.Flush(corelogger.SentryFlushDeadline)

			// A panic occurred, rollback and repanic
			lggr.Errorf("Panic in transaction, rolling back: %s", p)
			done := make(chan struct{})
			go func() {
				if rerr := tx.Rollback(); rerr != nil {
					lggr.Errorf("Failed to rollback on panic: %s", rerr)
				}
				close(done)
			}()
			select {
			case <-done:
				panic(p)
			case <-time.After(10 * time.Second):
				panic(fmt.Sprintf("panic in transaction; aborting rollback that took longer than 10s: %s", p))
			}
		} else if err != nil {
			lggr.Errorf("Error in transaction, rolling back: %s", err)
			// An error occurred, rollback and return error
			if rerr := tx.Rollback(); rerr != nil {
				err = multierr.Combine(err, errors.WithStack(rerr))
			}
		} else {
			// All good! Time to commit.
			err = errors.WithStack(tx.Commit())
		}
	}()

	err = fn(tx)

	return
}
