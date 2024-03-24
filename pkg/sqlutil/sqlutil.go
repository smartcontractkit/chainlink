package sqlutil

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type Queryer = DataSource

var _ DataSource = (*sqlx.DB)(nil)
var _ DataSource = (*sqlx.Tx)(nil)

// DataSource is implemented by [*sqlx.DB] & [*sqlx.Tx].
type DataSource interface {
	sqlx.ExtContext
	sqlx.PreparerContext
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

type TxOptions struct {
	sql.TxOptions
	OnPanic func(recovered any, rollbackErr error)
}

// TransactDataSource is a helper for executing transactions.
// This useful for executing raw SQL queries, or when using more than one type which is not supported by Transact.
func TransactDataSource(ctx context.Context, ds DataSource, opts *TxOptions, fn func(tx DataSource) error) error {
	return Transact(ctx, func(tx DataSource) DataSource { return tx }, ds, opts, fn)
}

// Transact is a helper for executing transactions with a domain specific type.
// A typical use looks like:
//
//	func (d *MyD) Transact(ctx context.Context, fn func(tx *MyD) error) (err error) {
//	  return sqlutil.Transact(ctx, d.new, d.db, nil, fn)
//	}
//
// If you need to combine multiple types in one transaction, you can declare a new type, or use TransactDataSource.
func Transact[D any](ctx context.Context, newD func(DataSource) D, ds DataSource, opts *TxOptions, fn func(tx D) error) (err error) {
	txds, ok := ds.(transactional)
	if !ok {
		// Unsupported or already inside another transaction.
		return fn(newD(ds))
	}
	if opts == nil {
		opts = &TxOptions{}
	}
	// Begin tx
	tx, err := func() (transaction, error) {
		// Support [DataSource]s wrapped via [WrapDataSource]
		if wrapped, ok := ds.(wrappedTransactional); ok {
			tx, terr := wrapped.BeginWrappedTxx(ctx, &opts.TxOptions)
			if terr != nil {
				return nil, terr
			}
			return tx, nil
		}

		tx, terr := txds.BeginTxx(ctx, &opts.TxOptions)
		if terr != nil {
			return nil, terr
		}
		return tx, nil
	}()
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

var _ transactional = (*sqlx.DB)(nil)

type transactional interface {
	// BeginTxx is implemented by *sqlx.DB but not *sqlx.Tx.
	BeginTxx(context.Context, *sql.TxOptions) (*sqlx.Tx, error)
}

var _ wrappedTransactional = (*wrappedTransactionalDataSource)(nil)

type wrappedTransactional interface {
	BeginWrappedTxx(context.Context, *sql.TxOptions) (transaction, error)
}

var _ transaction = (*wrappedTx)(nil)

type transaction interface {
	DataSource
	Commit() error
	Rollback() error
}
