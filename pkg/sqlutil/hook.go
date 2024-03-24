package sqlutil

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

var _ DataSource = &wrappedDataSource{}

// wrappedDataSource is a [DataSource] which invokes a [QueryHook] on each call.
type wrappedDataSource struct {
	db   DataSource
	lggr logger.Logger
	hook QueryHook
}

// QueryHook is a func that is executed for each query, providing an opportunity to measure, log, inspect/modify errors, etc.
// The do func *must* be called.
// Logs emitted through the provided logger.Logger will have the caller's line info.
//
// See [MonitorHook] and [TimeoutHook] for examples.
type QueryHook func(ctx context.Context, lggr logger.Logger, do func(context.Context) error, query string, args ...any) error

// WrapDataSource returns a new [DataSource] that calls each [QueryHook] in the provided order.
// If db implements [sqlx.BeginTxx], then the returned DataSource will also.
func WrapDataSource(db DataSource, l logger.Logger, hs ...QueryHook) DataSource {
	iq := wrappedDataSource{db: db,
		lggr: logger.Helper(logger.Named(l, "WrappedDB"), 2), // skip our own wrapper and one interceptor
		hook: noopHook,
	}
	switch len(hs) {
	case 0:
	case 1:
		iq.hook = hs[0]
	default:
		// Nest the QueryHook calls so that they are wrapped from first to last.
		// Example:
		// 	[A, B, C] => A(B(C(do())))
		for i := len(hs) - 1; i >= 0; i-- {
			next := hs[i]
			prev := iq.hook
			iq.hook = func(ctx context.Context, lggr logger.Logger, do func(context.Context) error, query string, args ...any) error {
				// opt: cache the construction of these loggers
				lggr = logger.Helper(lggr, 1) // skip one more for this wrapper
				return next(ctx, lggr, func(ctx context.Context) error {
					lggr = logger.Helper(lggr, 2) // skip two more for do() and this extra wrapper
					return prev(ctx, lggr, do, query, args...)
				}, query, args...)
			}
		}
	}

	if txdb, ok := db.(transactional); ok {
		// extra wrapper to make BeginTxx available
		return &wrappedTransactionalDataSource{
			wrappedDataSource: iq,
			txdb:              txdb,
		}
	}

	return &iq
}

func noopHook(ctx context.Context, lggr logger.Logger, do func(context.Context) error, query string, args ...any) error {
	return do(ctx)
}

func (w *wrappedDataSource) DriverName() string {
	return w.db.DriverName()
}

func (w *wrappedDataSource) Rebind(s string) string {
	return w.db.Rebind(s)
}

func (w *wrappedDataSource) BindNamed(s string, i interface{}) (string, []any, error) {
	return w.db.BindNamed(s, i)
}

func (w *wrappedDataSource) QueryContext(ctx context.Context, query string, args ...any) (rows *sql.Rows, err error) {
	err = w.hook(ctx, w.lggr, func(ctx context.Context) (err error) {
		rows, err = w.db.QueryContext(ctx, query, args...) //nolint
		return
	}, query, args...)
	return
}

func (w *wrappedDataSource) QueryxContext(ctx context.Context, query string, args ...any) (rows *sqlx.Rows, err error) {
	err = w.hook(ctx, w.lggr, func(ctx context.Context) (err error) {
		rows, err = w.db.QueryxContext(ctx, query, args...) //nolint:sqlclosecheck
		return
	}, query, args...)
	return
}

func (w *wrappedDataSource) QueryRowxContext(ctx context.Context, query string, args ...any) (row *sqlx.Row) {
	_ = w.hook(ctx, w.lggr, func(ctx context.Context) error {
		row = w.db.QueryRowxContext(ctx, query, args...)
		return nil
	}, query, args...)
	return
}

func (w *wrappedDataSource) ExecContext(ctx context.Context, query string, args ...any) (res sql.Result, err error) {
	err = w.hook(ctx, w.lggr, func(ctx context.Context) (err error) {
		res, err = w.db.ExecContext(ctx, query, args...)
		return
	}, query, args...)
	return
}

func (w *wrappedDataSource) NamedExecContext(ctx context.Context, query string, arg interface{}) (res sql.Result, err error) {
	err = w.hook(ctx, w.lggr, func(ctx context.Context) (err error) {
		res, err = w.db.NamedExecContext(ctx, query, arg)
		return
	}, query, arg)
	return
}

func (w *wrappedDataSource) PrepareContext(ctx context.Context, query string) (stmt *sql.Stmt, err error) {
	err = w.hook(ctx, w.lggr, func(ctx context.Context) (err error) {
		stmt, err = w.db.PrepareContext(ctx, query) //nolint:sqlclosecheck
		return
	}, query, nil)
	return
}

func (w *wrappedDataSource) PrepareNamedContext(ctx context.Context, query string) (stmt *sqlx.NamedStmt, err error) {
	err = w.hook(ctx, w.lggr, func(ctx context.Context) (err error) {
		stmt, err = w.db.PrepareNamedContext(ctx, query) //nolint:sqlclosecheck
		return
	}, query, nil)
	return
}

func (w *wrappedDataSource) GetContext(ctx context.Context, dest interface{}, query string, args ...any) error {
	return w.hook(ctx, w.lggr, func(ctx context.Context) error {
		return w.db.GetContext(ctx, dest, query, args...)
	}, query, args...)
}

func (w *wrappedDataSource) SelectContext(ctx context.Context, dest interface{}, query string, args ...any) error {
	return w.hook(ctx, w.lggr, func(ctx context.Context) error {
		return w.db.SelectContext(ctx, dest, query, args...)
	}, query, args...)
}

// wrappedTransactionalDataSource extends [wrappedDataSource] with BeginTxx and BeginWrappedTxx for initiating transactions.
type wrappedTransactionalDataSource struct {
	wrappedDataSource
	txdb transactional
}

func (w *wrappedTransactionalDataSource) BeginTxx(ctx context.Context, opts *sql.TxOptions) (tx *sqlx.Tx, err error) {
	err = w.hook(ctx, w.lggr, func(ctx context.Context) (err error) {
		tx, err = w.txdb.BeginTxx(ctx, opts)
		return
	}, "START TRANSACTION", nil)
	return
}

// BeginWrappedTxx is like BeginTxx, but wraps the returned tx with the same hook.
func (w *wrappedTransactionalDataSource) BeginWrappedTxx(ctx context.Context, opts *sql.TxOptions) (tx transaction, err error) {
	tx, err = w.BeginTxx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &wrappedTx{
		wrappedDataSource: wrappedDataSource{
			db:   tx,
			lggr: w.lggr,
			hook: w.hook,
		},
		tx: tx,
	}, nil
}

// wrappedTx extends [wrappedDataSource] with Commit and Rollback for completing a transaction.
type wrappedTx struct {
	wrappedDataSource
	tx transaction
}

func (w *wrappedTx) Commit() error { return w.tx.Commit() }

func (w *wrappedTx) Rollback() error { return w.tx.Rollback() }
