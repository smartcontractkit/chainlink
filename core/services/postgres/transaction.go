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

// NOTE: In an ideal world the timeouts below would be set to something sane in
// the postgres configuration by the user. Since we do not live in an ideal
// world, it is necessary to override them here.
//
// They cannot easily be set at a session level due to how Go's connection
// pooling works.
const (
	// LockTimeout controls the max time we will wait for any kind of database lock.
	// It's good to set this to _something_ because waiting for locks forever is really bad.
	LockTimeout = 15 * time.Second
	// IdleInTxSessionTimeout controls the max time we leave a transaction open and idle.
	// It's good to set this to _something_ because leaving transactions open forever is really bad.
	IdleInTxSessionTimeout = 1 * time.Hour
)

var (
	ErrNoDeadlineSet = errors.New("no deadline set")
)

// WARNING: Only use for nested txes inside ORM methods where you expect db to already have a ctx with a deadline.
func GormTransactionWithoutContext(db *gorm.DB, fc func(tx *gorm.DB) error, txOptss ...sql.TxOptions) (err error) {
	var txOpts sql.TxOptions
	if len(txOptss) > 0 {
		txOpts = txOptss[0]
	} else {
		txOpts = DefaultSqlTxOptions
	}
	return db.Transaction(func(tx *gorm.DB) error {
		err = tx.Exec(fmt.Sprintf(`SET LOCAL lock_timeout = %v; SET LOCAL idle_in_transaction_session_timeout = %v;`, LockTimeout.Milliseconds(), IdleInTxSessionTimeout.Milliseconds())).Error
		if err != nil {
			return errors.Wrap(err, "error setting transaction timeouts")
		}
		return fc(tx)
	}, &txOpts)
}

// DEPRECATED: Use the transaction manager instead.
func GormTransaction(ctx context.Context, db *gorm.DB, fc func(tx *gorm.DB) error, txOptss ...sql.TxOptions) (err error) {
	var txOpts sql.TxOptions
	if len(txOptss) > 0 {
		txOpts = txOptss[0]
	} else {
		txOpts = DefaultSqlTxOptions
	}
	if _, set := ctx.Deadline(); !set {
		return ErrNoDeadlineSet
	}
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err = tx.Exec(fmt.Sprintf(`SET LOCAL lock_timeout = %v; SET LOCAL idle_in_transaction_session_timeout = %v;`, LockTimeout.Milliseconds(), IdleInTxSessionTimeout.Milliseconds())).Error
		if err != nil {
			return errors.Wrap(err, "error setting transaction timeouts")
		}
		return fc(tx)
	}, &txOpts)
}

// DEPRECATED: Use the transaction manager instead.
func GormTransactionWithDefaultContext(db *gorm.DB, fc func(tx *gorm.DB) error, txOptss ...sql.TxOptions) error {
	var txOpts sql.TxOptions
	if len(txOptss) > 0 {
		txOpts = txOptss[0]
	} else {
		txOpts = DefaultSqlTxOptions
	}
	ctx, cancel := DefaultQueryCtx()
	defer cancel()
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Exec(fmt.Sprintf(`SET LOCAL lock_timeout = %v; SET LOCAL idle_in_transaction_session_timeout = %v;`, LockTimeout.Milliseconds(), IdleInTxSessionTimeout.Milliseconds())).Error
		if err != nil {
			return errors.Wrap(err, "error setting transaction timeouts")
		}
		return fc(tx)
	}, &txOpts)
	return err
}

func DBWithDefaultContext(db *gorm.DB, fc func(db *gorm.DB) error) error {
	ctx, cancel := DefaultQueryCtx()
	defer cancel()
	return fc(db.WithContext(ctx))
}

func SqlTransaction(ctx context.Context, rdb *sql.DB, fc func(tx *sqlx.Tx) error, txOpts ...sql.TxOptions) (err error) {
	db := WrapDbWithSqlx(rdb)
	return sqlxTransaction(ctx, db, fc, txOpts...)
}

func sqlxTransaction(ctx context.Context, db *sqlx.DB, fc func(tx *sqlx.Tx) error, txOpts ...sql.TxOptions) (err error) {
	opts := &DefaultSqlTxOptions
	if len(txOpts) > 0 {
		opts = &txOpts[0]
	}
	tx, err := db.BeginTxx(ctx, opts)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}

	defer func() {
		if p := recover(); p != nil {
			// A panic occurred, rollback and repanic
			logger.Errorf("panic in transaction, rolling back: %s", p)
			if rerr := tx.Rollback(); rerr != nil {
				logger.Error("failed to rollback on panic: %s", rerr)
			}
			panic(p)
		} else if err != nil {
			// An error occurred, rollback and return error
			if rerr := tx.Rollback(); rerr != nil {
				err = multierr.Combine(err, errors.WithStack(rerr))
			}
		} else {
			// All good! Time to commit.
			err = errors.WithStack(tx.Commit())
		}
	}()

	_, err = tx.Exec(fmt.Sprintf(`SET LOCAL lock_timeout = %v; SET LOCAL idle_in_transaction_session_timeout = %v;`, LockTimeout.Milliseconds(), IdleInTxSessionTimeout.Milliseconds()))
	if err != nil {
		return errors.Wrap(err, "error setting transaction timeouts")
	}

	err = fc(tx)

	return
}
