package sqlutil

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

var _ DB = &WrappedDB{}

// WrappedDB is a [DB] which invokes a [QueryHook] on each call.
type WrappedDB struct {
	db   DB
	lggr logger.Logger
	hook QueryHook
}

// QueryHook is a func that is executed for each query, providing an opportunity to measure, log, inspect/modify errors, etc.
// The do func *must* be called.
// Logs emitted through the provided logger.Logger will have the caller's line info.
//
// See [MonitorHook] and [TimeoutHook] for examples.
type QueryHook func(ctx context.Context, lggr logger.Logger, do func(context.Context) error, query string, args ...any) error

// NewWrappedDB returns a new [WrappedDB] that calls each [QueryHook] in the provided order.
func NewWrappedDB(db DB, l logger.Logger, hs ...QueryHook) *WrappedDB {
	iq := WrappedDB{db: db,
		lggr: logger.Helper(logger.Named(l, "WrappedDB"), 2), // skip our own wrapper and one interceptor
		hook: noopHook,
	}
	switch len(hs) {
	case 0:
		return &iq
	case 1:
		iq.hook = hs[0]
		return &iq
	}

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
	return &iq
}

func noopHook(ctx context.Context, lggr logger.Logger, do func(context.Context) error, query string, args ...any) error {
	return do(ctx)
}

func (w *WrappedDB) DriverName() string {
	return w.db.DriverName()
}

func (w *WrappedDB) Rebind(s string) string {
	return w.db.Rebind(s)
}

func (w *WrappedDB) BindNamed(s string, i interface{}) (string, []any, error) {
	return w.db.BindNamed(s, i)
}

func (w *WrappedDB) QueryContext(ctx context.Context, query string, args ...any) (rows *sql.Rows, err error) {
	err = w.hook(ctx, w.lggr, func(ctx context.Context) (err error) {
		rows, err = w.db.QueryContext(ctx, query, args...) //nolint
		return
	}, query, args...)
	return
}

func (w *WrappedDB) QueryxContext(ctx context.Context, query string, args ...any) (rows *sqlx.Rows, err error) {
	err = w.hook(ctx, w.lggr, func(ctx context.Context) (err error) {
		rows, err = w.db.QueryxContext(ctx, query, args...) //nolint:sqlclosecheck
		return
	}, query, args...)
	return
}

func (w *WrappedDB) QueryRowxContext(ctx context.Context, query string, args ...any) (row *sqlx.Row) {
	_ = w.hook(ctx, w.lggr, func(ctx context.Context) error {
		row = w.db.QueryRowxContext(ctx, query, args...)
		return nil
	}, query, args...)
	return
}

func (w *WrappedDB) ExecContext(ctx context.Context, query string, args ...any) (res sql.Result, err error) {
	err = w.hook(ctx, w.lggr, func(ctx context.Context) (err error) {
		res, err = w.db.ExecContext(ctx, query, args...)
		return
	}, query, args...)
	return
}

func (w *WrappedDB) PrepareContext(ctx context.Context, query string) (stmt *sql.Stmt, err error) {
	err = w.hook(ctx, w.lggr, func(ctx context.Context) (err error) {
		stmt, err = w.db.PrepareContext(ctx, query) //nolint:sqlclosecheck
		return
	}, query, nil)
	return
}

func (w *WrappedDB) GetContext(ctx context.Context, dest interface{}, query string, args ...any) error {
	return w.hook(ctx, w.lggr, func(ctx context.Context) error {
		return w.db.GetContext(ctx, dest, query, args...)
	}, query, args...)
}

func (w *WrappedDB) SelectContext(ctx context.Context, dest interface{}, query string, args ...any) error {
	return w.hook(ctx, w.lggr, func(ctx context.Context) error {
		return w.db.SelectContext(ctx, dest, query, args...)
	}, query, args...)
}
