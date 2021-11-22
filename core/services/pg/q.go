package pg

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/sqlx"
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
// 	func (o *orm) GetFoo(id int64, qopts ...pg.QOpt) (Foo, error) {
// 		q := pg.NewQ(q, qopts...)
// 		return q.Exec(...)
// 	}
//
// Now you can call it like so:
//
// 	orm.GetFoo(1) // will automatically have default query timeout context set
// 	orm.GetFoo(1, pg.WithParentCtx(ctx)) // will wrap the supplied parent context with the default query context
// 	orm.GetFoo(1, pg.WithQueryer(tx)) // allows to pass in a running transaction or anything else that implements Queryer
// 	orm.GetFoo(q, pg.WithQueryer(tx), pg.WithParentCtx(ctx)) // options can be combined
type QOpt func(*Q)

type LogConfig interface {
	LogSQL() bool
}

// WithQueryer sets the queryer
func WithQueryer(queryer Queryer) func(q *Q) {
	return func(q *Q) {
		if q.Queryer != nil {
			panic("queryer already set")
		}
		q.Queryer = queryer
	}
}

// WithParentCtx sets or overwrites the parent ctx
func WithParentCtx(ctx context.Context) func(q *Q) {
	return func(q *Q) {
		q.ParentCtx = ctx
	}
}

// MergeCtx allows callers to combine a ctx with a previously set parent context
// Responsibility for cancelling the passed context lies with caller
func MergeCtx(fn func(parentCtx context.Context) context.Context) func(q *Q) {
	return func(q *Q) {
		q.ParentCtx = fn(q.ParentCtx)
	}
}

var _ Queryer = Q{}
var slowSqlThreshold = time.Second

func init() {
	slowSqlThresholdStr := os.Getenv("SLOW_SQL_THRESHOLD")
	if len(slowSqlThresholdStr) > 0 {
		d, err := time.ParseDuration(slowSqlThresholdStr)
		if err != nil {
			log.Fatalf("failed to parse SLOW_SQL_THRESHOLD: %s", err)
		}
		slowSqlThreshold = d
	}
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
type Q struct {
	Queryer
	ParentCtx context.Context
	db        *sqlx.DB
	logger    logger.Logger
	config    LogConfig
}

// newQFromOpts is intended to be used in ORMs where the caller may wish to use
// either the default DB or pass an explicit Tx
func newQFromOpts(qopts []QOpt) (q Q) {
	for _, opt := range qopts {
		opt(&q)
	}
	return q
}

func NewQ(queryer Queryer, qopts ...QOpt) (q Q) {
	q = newQFromOpts(qopts)
	if q.Queryer == nil {
		q.Queryer = queryer
	}
	return
}

// TODO: this has to become new NewQ after all usages are fixed
func NewNewQ(db *sqlx.DB, logger logger.Logger, config LogConfig, qopts ...QOpt) (q Q) {
	q = newQFromOpts(qopts)
	if q.Queryer == nil {
		q.Queryer = db
	}
	q.db = db
	q.logger = logger
	q.config = config
	return
}

func PrepareQueryRowx(q Queryer, sql string, dest interface{}, arg interface{}) error {
	stmt, err := q.PrepareNamed(sql)
	if err != nil {
		return errors.Wrap(err, "error preparing named statement")
	}
	return errors.Wrap(stmt.QueryRowx(arg).Scan(dest), "error querying row")
}

func (q Q) WithOpts(qopts ...QOpt) (nq Q) {
	return NewNewQ(q.db, q.logger, q.config, qopts...)
}

func (q Q) Context() (context.Context, context.CancelFunc) {
	if q.ParentCtx == nil {
		return DefaultQueryCtx()
	}
	return DefaultQueryCtxWithParent(q.ParentCtx)
}

func (q Q) Transaction(lggr logger.Logger, fc func(q Queryer) error, txOpts ...TxOptions) error {
	ctx, cancel := q.Context()
	defer cancel()
	return SqlxTransaction(ctx, q.Queryer, lggr, fc, txOpts...)
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

	q.logSqlQuery(query, args...)
	begin := time.Now()
	defer q.postSqlLog(ctx, begin)

	res, err := q.Queryer.ExecContext(ctx, query, args...)
	return res, cancel, q.withLogError(err)
}
func (q Q) ExecQ(query string, args ...interface{}) error {
	_, cancel, err := q.ExecQIter(query, args...)
	cancel()
	return err
}
func (q Q) ExecQNamed(query string, arg interface{}) (err error) {
	query, args, err := q.BindNamed(query, arg)
	if err != nil {
		return errors.Wrap(err, "error binding arg")
	}
	return q.ExecQ(query, args...)
}

// Select and Get are safe to wrap the context cancellation because the rows
// are entirely consumed within the call
func (q Q) Select(dest interface{}, query string, args ...interface{}) error {
	ctx, cancel := q.Context()
	defer cancel()

	q.logSqlQuery(query, args...)
	begin := time.Now()
	defer q.postSqlLog(ctx, begin)

	return q.withLogError(q.Queryer.SelectContext(ctx, dest, query, args...))
}
func (q Q) Get(dest interface{}, query string, args ...interface{}) error {
	ctx, cancel := q.Context()
	defer cancel()

	q.logSqlQuery(query, args...)
	begin := time.Now()
	defer q.postSqlLog(ctx, begin)

	return q.withLogError(q.Queryer.GetContext(ctx, dest, query, args...))
}
func (q Q) GetNamed(sql string, dest interface{}, arg interface{}) error {
	query, args, err := q.BindNamed(sql, arg)
	if err != nil {
		return errors.Wrap(err, "error binding arg")
	}
	ctx, cancel := q.Context()
	defer cancel()

	q.logSqlQuery(query, args...)
	begin := time.Now()
	defer q.postSqlLog(ctx, begin)

	return q.withLogError(errors.Wrap(q.GetContext(ctx, dest, query, args...), "error in get query"))
}

type queryFmt struct {
	query string
	args  []interface{}
}

func (q queryFmt) String() string {
	if q.args == nil {
		return q.query
	}
	var pairs []string
	for i, arg := range q.args {
		pairs = append(pairs, fmt.Sprintf("$%d", i+1), fmt.Sprintf("%v", arg))
	}
	replacer := strings.NewReplacer(pairs...)
	return replacer.Replace(q.query)
}

func (q Q) logSqlQuery(query string, args ...interface{}) {
	if q.config == nil || q.logger == nil {
		return
	}
	if q.config.LogSQL() {
		q.logger.Debugf("SQL: %s", queryFmt{query, args})
	}
}

func (q Q) withLogError(err error) error {
	if err != nil && err != sql.ErrNoRows && q.logger != nil && q.config.LogSQL() {
		q.logger.Errorf("SQL ERROR: %v", err)
	}
	return err
}

func (q Q) postSqlLog(ctx context.Context, begin time.Time) {
	if q.logger == nil {
		return
	}
	elapsed := time.Since(begin)
	if ctx.Err() != nil {
		q.logger.Debugf("SQL CONTEXT CANCELLED: %d ms, err=%v", elapsed.Milliseconds(), ctx.Err())
	}
	if slowSqlThreshold > 0 && elapsed > slowSqlThreshold {
		q.logger.Warnf("SLOW SQL QUERY: %d ms", elapsed.Milliseconds())
	}
}
