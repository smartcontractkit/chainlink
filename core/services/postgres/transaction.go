package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/sqlx"
	"go.uber.org/multierr"
	"gorm.io/gorm"
)

type TxOptions struct {
	sql.TxOptions
	LockTimeout            time.Duration
	IdleInTxSessionTimeout time.Duration
}

// NOTE: In an ideal world the timeouts below would be set to something sane in
// the postgres configuration by the user. Since we do not live in an ideal
// world, it is necessary to override them here.
//
// They cannot easily be set at a session level due to how Go's connection
// pooling works.
const (
	// DefaultLockTimeout controls the max time we will wait for any kind of database lock.
	// It's good to set this to _something_ because waiting for locks forever is really bad.
	DefaultLockTimeout = 15 * time.Second
	// DefaultIdleInTxSessionTimeout controls the max time we leave a transaction open and idle.
	// It's good to set this to _something_ because leaving transactions open forever is really bad.
	DefaultIdleInTxSessionTimeout = 1 * time.Hour
	// NOTE: This is the default level in Postgres anyway, we just make it
	// explicit here
	DefaultIsolation = sql.LevelReadCommitted
)

var (
	ErrNoDeadlineSet = errors.New("no deadline set")
)

func applyDefaults(optss []TxOptions) (lockTimeout, idleInTxSessionTimeout time.Duration, txOpts sql.TxOptions) {
	lockTimeout = DefaultLockTimeout
	idleInTxSessionTimeout = DefaultIdleInTxSessionTimeout
	txIsolation := DefaultIsolation
	readOnly := false
	if len(optss) > 0 {
		opts := optss[0]
		if opts.LockTimeout != 0 {
			lockTimeout = opts.LockTimeout
		}
		if opts.IdleInTxSessionTimeout != 0 {
			idleInTxSessionTimeout = opts.IdleInTxSessionTimeout
		}
		if opts.Isolation != 0 {
			txIsolation = opts.Isolation
		}
		readOnly = opts.ReadOnly
	}
	txOpts = sql.TxOptions{
		Isolation: txIsolation,
		ReadOnly:  readOnly,
	}
	return
}

func GormTransactionWithoutContext(db *gorm.DB, fn func(tx *gorm.DB) error, optss ...TxOptions) (err error) {
	lockTimeout, idleInTxSessionTimeout, txOpts := applyDefaults(optss)
	return db.Transaction(func(tx *gorm.DB) error {
		err = tx.Exec(fmt.Sprintf(`SET LOCAL lock_timeout = %v; SET LOCAL idle_in_transaction_session_timeout = %v;`, lockTimeout.Milliseconds(), idleInTxSessionTimeout.Milliseconds())).Error
		if err != nil {
			return errors.Wrap(err, "error setting transaction timeouts")
		}
		return fn(tx)
	}, &txOpts)
}

func GormTransaction(ctx context.Context, db *gorm.DB, fn func(tx *gorm.DB) error, optss ...TxOptions) (err error) {
	lockTimeout, idleInTxSessionTimeout, txOpts := applyDefaults(optss)
	if _, set := ctx.Deadline(); !set {
		return ErrNoDeadlineSet
	}
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err = tx.Exec(fmt.Sprintf(`SET LOCAL lock_timeout = %v; SET LOCAL idle_in_transaction_session_timeout = %v;`, lockTimeout.Milliseconds(), idleInTxSessionTimeout.Milliseconds())).Error
		if err != nil {
			return errors.Wrap(err, "error setting transaction timeouts")
		}
		return fn(tx)
	}, &txOpts)
}

func GormTransactionWithDefaultContext(db *gorm.DB, fn func(tx *gorm.DB) error, optss ...TxOptions) error {
	lockTimeout, idleInTxSessionTimeout, txOpts := applyDefaults(optss)
	ctx, cancel := DefaultQueryCtx()
	defer cancel()
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Exec(fmt.Sprintf(`SET LOCAL lock_timeout = %v; SET LOCAL idle_in_transaction_session_timeout = %v;`, lockTimeout.Milliseconds(), idleInTxSessionTimeout.Milliseconds())).Error
		if err != nil {
			return errors.Wrap(err, "error setting transaction timeouts")
		}
		return fn(tx)
	}, &txOpts)
	return err
}

func DBWithDefaultContext(db *gorm.DB, fn func(db *gorm.DB) error) error {
	ctx, cancel := DefaultQueryCtx()
	defer cancel()
	return fn(db.WithContext(ctx))
}

func SqlTransaction(ctx context.Context, rdb *sql.DB, lggr logger.Logger, fn func(tx *sqlx.Tx) error, optss ...TxOptions) (err error) {
	db := WrapDbWithSqlx(rdb)
	return sqlxTransaction(ctx, db, lggr, fn, optss...)
}

func sqlxTransaction(ctx context.Context, db *sqlx.DB, lggr logger.Logger, fn func(tx *sqlx.Tx) error, optss ...TxOptions) (err error) {
	wrapFn := func(q Queryer) error {
		tx, ok := q.(*sqlx.Tx)
		if !ok {
			panic(fmt.Sprintf("expected q to be %T but got %T", tx, q))
		}
		return fn(tx)
	}
	return sqlxTransactionQ(ctx, db, lggr, wrapFn, optss...)
}

func sqlxTransactionQ(ctx context.Context, db *sqlx.DB, lggr logger.Logger, fn func(q Queryer) error, optss ...TxOptions) (err error) {
	lockTimeout, idleInTxSessionTimeout, txOpts := applyDefaults(optss)
	fmt.Println(txOpts)

	tx, err := db.BeginTxx(ctx, &txOpts)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}

	defer func() {
		if p := recover(); p != nil {
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
			lggr.Debugf("Error in transaction, rolling back: %s", err)
			// An error occurred, rollback and return error
			if rerr := tx.Rollback(); rerr != nil {
				err = multierr.Combine(err, errors.WithStack(rerr))
			}
		} else {
			// All good! Time to commit.
			err = errors.WithStack(tx.Commit())
		}
	}()

	_, err = tx.Exec(fmt.Sprintf(`SET LOCAL lock_timeout = %v; SET LOCAL idle_in_transaction_session_timeout = %v;`, lockTimeout.Milliseconds(), idleInTxSessionTimeout.Milliseconds()))
	if err != nil {
		return errors.Wrap(err, "error setting transaction timeouts")
	}

	err = fn(tx)

	return
}
