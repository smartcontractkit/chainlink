package sqlutil

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

// Queryer is implemented by [*sqlx.DB], [*sqlx.Tx], & [*sqlx.Conn].
type Queryer interface {
	sqlx.ExtContext
	sqlx.PreparerContext
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type TxOptions struct {
	sql.TxOptions
	OnPanic func(recovered any, rollbackErr error)
}

// Transact is a helper for executing transactions with a domain specific type.
// A typical use looks like:
//
//	func (d *MyD) Transaction(ctx context.Context, fn func(*MyD) error) (err error) {
//	  return sqlutil.Transact(ctx, d.new, d.q, nil, fn)
//	}
func Transact[D any](ctx context.Context, newD func(Queryer) D, q Queryer, opts *TxOptions, fn func(D) error) (err error) {
	db, ok := q.(interface {
		// BeginTxx is implemented by *sqlx.DB & *sqlx.Conn, but not *sqlx.Tx.
		BeginTxx(context.Context, *sql.TxOptions) (*sqlx.Tx, error)
	})
	if !ok {
		// Unsupported or already inside another transaction.
		return fn(newD(q))
	}
	if opts == nil {
		opts = &TxOptions{}
	}
	tx, err := db.BeginTxx(ctx, &opts.TxOptions)
	if err != nil {
		return err
	}
	defer func() {
		if recovered := recover(); recovered != nil {
			rbErr := tx.Rollback()
			if onPanic := opts.OnPanic; onPanic != nil {
				onPanic(recovered, rbErr)
			}
			panic(recovered)
		} else if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil &&
				!errors.Is(rbErr, sql.ErrTxDone) { // already committed or rolled back
				err = errors.Join(err, rbErr)
			}
			return
		}
		err = tx.Commit()
	}()
	err = fn(newD(tx))
	return
}
