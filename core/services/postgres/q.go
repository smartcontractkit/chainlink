package postgres

// TODO: Rename this package to "pg"
// https://app.shortcut.com/chainlinklabs/story/20021/rename-postgres-to-pg

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/logger"
)

// QOpt pattern for ORM methods aims to clarify usage and remove some common footguns, notably:
//
// 1. It should be easy and obvious how to pass a parent context or a transaction into an ORM method
// 2. Simple queries should not be cluttered
// 3. It should have compile-time safety and be explicit
// 4. It should enforce some sort of context deadline on all queries by default
// 5. It should optimise for clarity and readability
// 6. It should mandate using sqlx everywhere, gorm is forbidden in new code
// 7. It should make using sqlx a little more convenient by wrapping certain methods
// 8. It allows easier mocking of DB calls (Queryer is an interface)
//
// The two main concepts introduced are:
//
// A `Q` struct that wraps a `sqlx.DB` or `sqlx.Tx` and implements the `postgres.Queryer` interface.
//
// This struct is initialised with `QOpts` which define how the queryer should behave. `QOpts` can define a parent context, an open transaction or other options to configure the Queryer.
//
// A sample ORM method looks like this:
//
// 	func (o *orm) GetFoo(id int64, qopts ...postgres.QOpt) (Foo, error) {
// 		q := postgres.NewQ(q, qopts...)
// 		return q.Exec(...)
// 	}
//
// Now you can call it like so:
//
// 	orm.GetFoo(1) // will automatically have default query timeout context set
// 	orm.GetFoo(1, postgres.WithParentCtx(ctx)) // will wrap the supplied parent context with the default query context
// 	orm.GetFoo(1, postgres.WithQueryer(tx)) // allows to pass in a running transaction or anything else that implements Queryer
// 	orm.GetFoo(q, postgres.WithQueryer(tx), postgres.WithParentCtx(ctx)) // options can be combined
type QOpt func(*Q)

// WithQueryer sets the queryer
func WithQueryer(queryer Queryer) func(q *Q) {
	return func(q *Q) {
		if q.Queryer != nil {
			panic("queryer already set")
		}
		q.Queryer = queryer
	}
}

// WithParentCtx sets the parent ctx
func WithParentCtx(ctx context.Context) func(q *Q) {
	return func(q *Q) {
		q.ParentCtx = ctx
	}
}

var _ Queryer = Q{}

// Q wraps an underlying queryer (either a *sqlx.DB or a *sqlx.Tx)
//
// It is designed to make handling *sqlx.Tx or *sqlx.DB a little bit safer by
// preventing footguns such as having no deadline on contexts.
//
// It also handles nesting transactions.
//
// It automatically adds the default context deadline to all non-context
// queries (if you _really_ want to issue a query without a context, use the
// underlying Queryer)
//
// This is not the prettiest construct but without macros its about the best we
// can do.
type Q struct {
	Queryer
	lggr      logger.Logger
	ParentCtx context.Context
}

// NewQFromOpts is intended to be used in ORMs where the caller may wish to use
// either the default DB or pass an explicit Tx
func NewQFromOpts(qopts []QOpt) (q Q) {
	for _, opt := range qopts {
		opt(&q)
	}
	return q
}

func NewQ(queryer Queryer, qopts ...QOpt) (q Q) {
	q = NewQFromOpts(qopts)
	if q.Queryer == nil {
		q.Queryer = queryer
	}
	return
}

func PrepareQueryRowx(q Queryer, sql string, dest interface{}, arg interface{}) error {
	stmt, err := q.PrepareNamed(sql)
	if err != nil {
		return errors.Wrap(err, "error preparing named statement")
	}
	return errors.Wrap(stmt.QueryRowx(arg).Scan(dest), "error querying row")
}

func (q Q) Context() (context.Context, context.CancelFunc) {
	if q.ParentCtx == nil {
		return DefaultQueryCtx()
	}
	return DefaultQueryCtxWithParent(q.ParentCtx)
}

func (q Q) Transaction(lggr logger.Logger, fc func(q Queryer) error) error {
	ctx, cancel := q.Context()
	defer cancel()
	return SqlxTransaction(ctx, q.Queryer, lggr, fc)
}

// CAUTION: A subtle problem lurks here, because the following code is buggy:
//
//     ctx, cancel := context.WithCancel(context.Background())
//     rows, err := db.QueryContext(ctx, "SELECT foo")
//     cancel() // canceling here "poisons" the scan below
//     for rows.Next() {
//       rows.Scan(...)
//     }
//
// We must cancel the context only after we have completely finished using the
// returned rows or result from the query/exec
//
// For this reasons, the following functions return a context.CancelFunc and it
// is up to the caller to ensure that cancel is called after it has finished
//
// Generally speaking, it makes more sense to use Get/Select in most cases,
// which avoids this problem
func (q Q) ExecQIter(query string, args ...interface{}) (sql.Result, context.CancelFunc, error) {
	ctx, cancel := q.Context()
	res, err := q.Queryer.ExecContext(ctx, query, args...)
	return res, cancel, err
}
func (q Q) ExecQ(query string, args ...interface{}) error {
	_, cancel, err := q.ExecQIter(query, args...)
	cancel()
	return err
}

// Select and Get are safe to wrap the context cancellation because the rows
// are entirely consumed within the call
func (q Q) Select(dest interface{}, query string, args ...interface{}) error {
	ctx, cancel := q.Context()
	defer cancel()
	return q.Queryer.SelectContext(ctx, dest, query, args...)
}
func (q Q) Get(dest interface{}, query string, args ...interface{}) error {
	ctx, cancel := q.Context()
	defer cancel()
	return q.Queryer.GetContext(ctx, dest, query, args...)
}
func (q Q) GetNamed(sql string, dest interface{}, arg interface{}) error {
	query, args, err := q.BindNamed(sql, arg)
	if err != nil {
		return errors.Wrap(err, "error binding arg")
	}
	ctx, cancel := q.Context()
	defer cancel()
	return errors.Wrap(q.GetContext(ctx, dest, query, args...), "error in get query")
}
