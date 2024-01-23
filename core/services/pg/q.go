package pg

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commonpg "github.com/smartcontractkit/chainlink-common/pkg/services/pg"
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
// A `Q` struct that wraps a `sqlx.DB` or `sqlx.Tx` and implements the `pg.Queryer` interface.
//
// This struct is initialised with `QOpts` which define how the queryer should behave. `QOpts` can define a parent context, an open transaction or other options to configure the Queryer.
//
// A sample ORM method looks like this:
//
//	func (o *orm) GetFoo(id int64, qopts ...pg.QOpt) (Foo, error) {
//		q := pg.NewQ(q, qopts...)
//		return q.Exec(...)
//	}
//
// Now you can call it like so:
//
//	orm.GetFoo(1) // will automatically have default query timeout context set
//	orm.GetFoo(1, pg.WithParentCtx(ctx)) // will wrap the supplied parent context with the default query context
//	orm.GetFoo(1, pg.WithQueryer(tx)) // allows to pass in a running transaction or anything else that implements Queryer
//	orm.GetFoo(q, pg.WithQueryer(tx), pg.WithParentCtx(ctx)) // options can be combined
type QOpt = commonpg.QOpt //func(*Q)

// WithQueryer sets the queryer
func WithQueryer(queryer Queryer) QOpt {
	return commonpg.WithQueryer(queryer)
}

// WithParentCtx sets or overwrites the parent ctx
func WithParentCtx(ctx context.Context) QOpt {
	return commonpg.WithParentCtx(ctx)
}

// If the parent has a timeout, just use that instead of DefaultTimeout
func WithParentCtxInheritTimeout(ctx context.Context) QOpt {
	return commonpg.WithParentCtxInheritTimeout(ctx)
}

// WithLongQueryTimeout prevents the usage of the `DefaultQueryTimeout` duration and uses `OneMinuteQueryTimeout` instead
// Some queries need to take longer when operating over big chunks of data, like deleting jobs, but we need to keep some upper bound timeout
func WithLongQueryTimeout() QOpt {
	return func(q *Q) {
		q.QueryTimeout = longQueryTimeout
	}
}

var _ Queryer = Q{}

type QConfig interface {
	LogSQL() bool
	DefaultQueryTimeout() time.Duration
}

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

type Q = commonpg.Q

func NewQ(db *sqlx.DB, lggr logger.Logger, config QConfig, qopts ...QOpt) (q Q) {
	return commonpg.NewQ(db, lggr, config, qopts...)
}

func PrepareQueryRowx(q Queryer, sql string, dest interface{}, arg interface{}) error {
	return commonpg.PrepareQueryRowx(q, sql, dest, arg)
}

// sprintQ formats the query with the given args and returns the resulting string.
func sprintQ(query string, args []interface{}) string {
	return commonpg.SprintQ(query, args)
}
