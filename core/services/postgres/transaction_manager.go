package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"go.uber.org/multierr"
	"gorm.io/gorm"
)

//go:generate mockery --name TransactionManager --output ./mocks/ --case=underscore

// A TxFn is a function that will be called with a context which has a transaction
// injected as a value. This can be used for executing statements and queries
// against a database.
type TxFn func(context.Context) error

// TransactionManagerOption configures how we set up the transaction
type TransactionOption func(opts *transactionOptions)

// transactionOptions defines a list of transaction options which configure the
// transaction.
type transactionOptions struct {
	txOpts          sql.TxOptions
	withoutDeadline bool
}

// WithTxOptions returns a TransactionOption which sets the sql.TxOptions on the
// transaction.
func WithTxOptions(txopts sql.TxOptions) TransactionOption {
	return func(opts *transactionOptions) {
		opts.txOpts = txopts
	}
}

func WithoutDeadline() TransactionOption {
	return func(opts *transactionOptions) {
		opts.withoutDeadline = true
	}
}

type TransactionManager interface {
	Transact(TxFn, ...TransactionOption) error
	TransactWithContext(ctx context.Context, fn TxFn, optsFn ...TransactionOption) (err error)
}

type gormTransactionManager struct {
	db *gorm.DB
}

func NewGormTransactionManager(db *gorm.DB) TransactionManager {
	return &gormTransactionManager{db: db}
}

// txKey is the context key for a transaction value.
type txKey struct{}

// Transact creates a new transaction with sane defaults.
func (txm *gormTransactionManager) Transact(fn TxFn, optsFn ...TransactionOption) (err error) {
	ctx, cancel := DefaultQueryCtx()
	defer cancel()

	return txm.TransactWithContext(ctx, fn, optsFn...)
}

// Transact creates a new transaction and injects it into the provided context.
// It handles the rollback/commit based on the error object returned by the `TxFn`.
func (txm *gormTransactionManager) TransactWithContext(ctx context.Context, fn TxFn, optsFn ...TransactionOption) (err error) {
	// Initialize the options with defaults
	opts := &transactionOptions{
		txOpts: DefaultSqlTxOptions,
	}

	// Overwrite any opts with declared option setters
	for _, set := range optsFn {
		set(opts)
	}

	// Start the transaction and insert it into the context.
	tx := txm.db.Begin(&opts.txOpts)
	if err = tx.Error; err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	ctx = context.WithValue(ctx, txKey{}, tx)

	// Ensure that a deadline is set unless disabled by an option.
	if !opts.withoutDeadline {
		if _, ok := ctx.Deadline(); !ok {
			return ErrNoDeadlineSet
		}
	}

	// Handle rollback/commits
	defer func() {
		if p := recover(); p != nil {
			// A panic occurred, rollback and repanic
			logger.Errorf("panic in transaction, rolling back: %s", p)
			done := make(chan struct{})
			go func() {
				if rerr := tx.Rollback().Error; rerr != nil {
					logger.Error("failed to rollback on panic: %s", rerr)
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
			logger.Debugf("error in transaction, rolling back: %s", err)
			// An error occurred, rollback and return error
			if rerr := tx.Rollback().Error; rerr != nil {
				err = multierr.Combine(err, errors.WithStack(rerr))
			}
		} else {
			// All good! Time to commit.
			err = errors.WithStack(tx.Commit().Error)
		}
	}()

	// Set the local lock timeout
	err = tx.Exec(
		fmt.Sprintf(`SET LOCAL lock_timeout = %v; SET LOCAL idle_in_transaction_session_timeout = %v;`,
			LockTimeout.Milliseconds(),
			IdleInTxSessionTimeout.Milliseconds()),
	).Error
	if err != nil {
		return errors.Wrap(err, "error setting transaction timeouts")
	}

	err = fn(ctx)
	return err
}

// TxFromContext extracts the tx from the context. If no transaction value is
// provided in the context, it returns the gorm.DB.
func TxFromContext(ctx context.Context, db *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}

	return db
}
