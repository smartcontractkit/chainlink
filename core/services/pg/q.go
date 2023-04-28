package pg

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var promSQLQueryTime = promauto.NewHistogram(prometheus.HistogramOpts{
	Name:    "sql_query_timeout_percent",
	Help:    "SQL query time as a pecentage of timeout.",
	Buckets: []float64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 110, 120},
})

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
type QOpt func(*Q)

// WithQueryer sets the queryer
func WithQueryer(queryer Queryer) QOpt {
	return func(q *Q) {
		if q.Queryer != nil {
			panic("queryer already set")
		}
		q.Queryer = queryer
	}
}

// WithParentCtx sets or overwrites the parent ctx
func WithParentCtx(ctx context.Context) QOpt {
	return func(q *Q) {
		q.ParentCtx = ctx
	}
}

// If the parent has a timeout, just use that instead of DefaultTimeout
func WithParentCtxInheritTimeout(ctx context.Context) QOpt {
	return func(q *Q) {
		q.ParentCtx = ctx
		deadline, ok := q.ParentCtx.Deadline()
		if ok {
			q.QueryTimeout = time.Until(deadline)
		}
	}
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
	DatabaseDefaultQueryTimeout() time.Duration
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
	ParentCtx    context.Context
	db           *sqlx.DB
	logger       logger.Logger
	config       QConfig
	QueryTimeout time.Duration
}

func NewQ(db *sqlx.DB, logger logger.Logger, config QConfig, qopts ...QOpt) (q Q) {
	for _, opt := range qopts {
		opt(&q)
	}

	q.db = db
	q.logger = logger.Helper(2)
	q.config = config

	if q.Queryer == nil {
		q.Queryer = db
	}
	if q.ParentCtx == nil {
		q.ParentCtx = context.Background()
	}
	if q.QueryTimeout <= 0 {
		q.QueryTimeout = q.config.DatabaseDefaultQueryTimeout()
	}
	return
}

func (q Q) originalLogger() logger.Logger {
	return q.logger.Helper(-2)
}

func PrepareQueryRowx(q Queryer, sql string, dest interface{}, arg interface{}) error {
	stmt, err := q.PrepareNamed(sql)
	if err != nil {
		return errors.Wrap(err, "error preparing named statement")
	}
	return errors.Wrap(stmt.QueryRowx(arg).Scan(dest), "error querying row")
}

func (q Q) WithOpts(qopts ...QOpt) Q {
	return NewQ(q.db, q.originalLogger(), q.config, qopts...)
}

func (q Q) Context() (context.Context, context.CancelFunc) {
	return context.WithTimeout(q.ParentCtx, q.QueryTimeout)
}

func (q Q) Transaction(fc func(q Queryer) error, txOpts ...TxOptions) error {
	ctx, cancel := q.Context()
	defer cancel()
	return SqlxTransaction(ctx, q.Queryer, q.originalLogger(), fc, txOpts...)
}

// CAUTION: A subtle problem lurks here, because the following code is buggy:
//
//	ctx, cancel := context.WithCancel(context.Background())
//	rows, err := db.QueryContext(ctx, "SELECT foo")
//	cancel() // canceling here "poisons" the scan below
//	for rows.Next() {
//	  rows.Scan(...)
//	}
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

	ql := q.newQueryLogger(query, args)
	ql.logSqlQuery()
	defer ql.postSqlLog(ctx, time.Now())

	res, err := q.Queryer.ExecContext(ctx, query, args...)
	return res, cancel, ql.withLogError(err)
}
func (q Q) ExecQ(query string, args ...interface{}) error {
	ctx, cancel := q.Context()
	defer cancel()

	ql := q.newQueryLogger(query, args)
	ql.logSqlQuery()
	defer ql.postSqlLog(ctx, time.Now())

	_, err := q.Queryer.ExecContext(ctx, query, args...)
	return ql.withLogError(err)
}
func (q Q) ExecQNamed(query string, arg interface{}) (err error) {
	query, args, err := q.BindNamed(query, arg)
	if err != nil {
		return errors.Wrap(err, "error binding arg")
	}
	ctx, cancel := q.Context()
	defer cancel()

	ql := q.newQueryLogger(query, args)
	ql.logSqlQuery()
	defer ql.postSqlLog(ctx, time.Now())

	_, err = q.Queryer.ExecContext(ctx, query, args...)
	return ql.withLogError(err)
}

// Select and Get are safe to wrap the context cancellation because the rows
// are entirely consumed within the call
func (q Q) Select(dest interface{}, query string, args ...interface{}) error {
	ctx, cancel := q.Context()
	defer cancel()

	ql := q.newQueryLogger(query, args)
	ql.logSqlQuery()
	defer ql.postSqlLog(ctx, time.Now())

	return ql.withLogError(q.Queryer.SelectContext(ctx, dest, query, args...))
}
func (q Q) Get(dest interface{}, query string, args ...interface{}) error {
	ctx, cancel := q.Context()
	defer cancel()

	ql := q.newQueryLogger(query, args)
	ql.logSqlQuery()
	defer ql.postSqlLog(ctx, time.Now())

	return ql.withLogError(q.Queryer.GetContext(ctx, dest, query, args...))
}
func (q Q) GetNamed(sql string, dest interface{}, arg interface{}) error {
	query, args, err := q.BindNamed(sql, arg)
	if err != nil {
		return errors.Wrap(err, "error binding arg")
	}
	ctx, cancel := q.Context()
	defer cancel()

	ql := q.newQueryLogger(query, args)
	ql.logSqlQuery()
	defer ql.postSqlLog(ctx, time.Now())

	return ql.withLogError(errors.Wrap(q.GetContext(ctx, dest, query, args...), "error in get query"))
}

func (q Q) newQueryLogger(query string, args []interface{}) *queryLogger {
	return &queryLogger{Q: q, query: query, args: args}
}

// sprintQ formats the query with the given args and returns the resulting string.
func sprintQ(query string, args []interface{}) string {
	if args == nil {
		return query
	}
	var pairs []string
	for i, arg := range args {
		// We print by type so one can directly take the logged query string and execute it manually in pg.
		// Annoyingly it seems as though the logger itself will add an extra \, so you still have to remove that.
		switch v := arg.(type) {
		case []byte:
			pairs = append(pairs, fmt.Sprintf("$%d", i+1), fmt.Sprintf("'\\x%x'", v))
		case common.Address:
			pairs = append(pairs, fmt.Sprintf("$%d", i+1), fmt.Sprintf("'\\x%x'", v.Bytes()))
		case common.Hash:
			pairs = append(pairs, fmt.Sprintf("$%d", i+1), fmt.Sprintf("'\\x%x'", v.Bytes()))
		case pq.ByteaArray:
			var s strings.Builder
			fmt.Fprintf(&s, "('\\x%x'", v[0])
			for j := 1; j < len(v); j++ {
				fmt.Fprintf(&s, ",'\\x%x'", v[j])
			}
			pairs = append(pairs, fmt.Sprintf("$%d", i+1), fmt.Sprintf("%s)", s.String()))
		default:
			pairs = append(pairs, fmt.Sprintf("$%d", i+1), fmt.Sprintf("%v", arg))
		}
	}
	replacer := strings.NewReplacer(pairs...)
	queryWithVals := replacer.Replace(query)
	return strings.ReplaceAll(strings.ReplaceAll(queryWithVals, "\n", " "), "\t", " ")
}

// queryLogger extends Q with logging helpers for a particular query w/ args.
type queryLogger struct {
	Q

	query string
	args  []interface{}

	str     string
	strOnce sync.Once
}

func (q *queryLogger) String() string {
	q.strOnce.Do(func() {
		q.str = sprintQ(q.query, q.args)
	})
	return q.str
}

func (q *queryLogger) logSqlQuery() {
	if q.config != nil && q.config.LogSQL() {
		q.logger.Debugw("SQL QUERY", "sql", q)
	}
}

func (q *queryLogger) withLogError(err error) error {
	if err != nil && !errors.Is(err, sql.ErrNoRows) && q.config != nil && q.config.LogSQL() {
		q.logger.Errorw("SQL ERROR", "err", err, "sql", q)
	}
	return err
}

// postSqlLog logs about context cancellation and timing after a query returns.
// Queries which use their full timeout log critical level. More than 50% log error, and 10% warn.
func (q *queryLogger) postSqlLog(ctx context.Context, begin time.Time) {
	elapsed := time.Since(begin)
	if ctx.Err() != nil {
		q.logger.Debugw("SQL CONTEXT CANCELLED", "ms", elapsed.Milliseconds(), "err", ctx.Err(), "sql", q)
	}

	timeout := q.QueryTimeout
	if timeout <= 0 {
		timeout = DefaultQueryTimeout
	}

	pct := float64(elapsed) / float64(timeout)
	pct *= 100

	kvs := []any{"ms", elapsed.Milliseconds(), "timeout", timeout.Milliseconds(), "percent", strconv.FormatFloat(pct, 'f', 1, 64), "sql", q}

	if elapsed >= timeout {
		q.logger.Criticalw("SLOW SQL QUERY", kvs...)
	} else if errThreshold := timeout / 5; errThreshold > 0 && elapsed > errThreshold {
		q.logger.Errorw("SLOW SQL QUERY", kvs...)
	} else if warnThreshold := timeout / 10; warnThreshold > 0 && elapsed > warnThreshold {
		q.logger.Warnw("SLOW SQL QUERY", kvs...)
	}

	promSQLQueryTime.Observe(pct)
}
