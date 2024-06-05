package sqlutil

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

const slowMsg = "SLOW SQL QUERY"

// PromSQLQueryTime is exported temporarily while transitioning the core ORMs.
var PromSQLQueryTime = promauto.NewHistogram(prometheus.HistogramOpts{
	Name:    "sql_query_timeout_percent",
	Help:    "SQL query time as a percentage of timeout.",
	Buckets: []float64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 110, 120},
})

// MonitorHook returns a [QueryHook] that measures the timing of each query and logs about slow queries at increasing levels of severity.
// When logAll returns true, every query and error will be logged.
func MonitorHook(logAll func() bool) QueryHook {
	return func(ctx context.Context, lggr logger.Logger, do func(context.Context) error, query string, args ...any) error {
		shouldLog := logAll()
		ql := newQueryLogger(lggr, query, args...)
		if shouldLog {
			ql.logQuery()
		}
		defer ql.logTiming(ctx, time.Now())
		err := do(ctx)
		if shouldLog && err != nil && !errors.Is(err, sql.ErrNoRows) {
			ql.logError(err)
		}
		return err
	}
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
		case pq.ByteaArray:
			pairs = append(pairs, fmt.Sprintf("$%d", i+1))
			if v == nil {
				pairs = append(pairs, "NULL")
				continue
			}
			if len(v) == 0 {
				pairs = append(pairs, "ARRAY[]")
				continue
			}
			var s strings.Builder
			fmt.Fprintf(&s, "ARRAY['\\x%x'", v[0])
			for j := 1; j < len(v); j++ {
				fmt.Fprintf(&s, ",'\\x%x'", v[j])
			}
			pairs = append(pairs, fmt.Sprintf("%s]", s.String()))
		case string:
			pairs = append(pairs, fmt.Sprintf("$%d", i+1), fmt.Sprintf("'%s'", v))
		default:
			pairs = append(pairs, fmt.Sprintf("$%d", i+1), fmt.Sprintf("%v", v))
		}
	}
	replacer := strings.NewReplacer(pairs...)
	queryWithVals := replacer.Replace(query)
	return strings.ReplaceAll(strings.ReplaceAll(queryWithVals, "\n", " "), "\t", " ")
}

// queryLogger extends Q with logging helpers for a particular query w/ args.
type queryLogger struct {
	lggr logger.SugaredLogger

	query string
	args  []interface{}

	str func() string
}

func newQueryLogger(lggr logger.Logger, query string, args ...any) *queryLogger {
	return &queryLogger{
		// skip another level since we use internal helpers
		lggr:  logger.Sugared(logger.Helper(lggr, 1)),
		query: query, args: args,
		str: sync.OnceValue(func() string {
			return sprintQ(query, args)
		}),
	}
}

func (q *queryLogger) String() string {
	return q.str()
}

func (q *queryLogger) logQuery() {
	q.lggr.Debugw("SQL QUERY", "sql", q)
}

func (q *queryLogger) logError(err error) {
	q.lggr.Errorw("SQL ERROR", "err", err, "sql", q)
}

// logTiming logs about context errors and timing.
// To be deferred to run after a query returns.
func (q *queryLogger) logTiming(ctx context.Context, start time.Time) {
	elapsed := time.Since(start)

	if err := ctx.Err(); err != nil {
		if errors.Is(err, context.Canceled) {
			q.lggr.Debugw("SQL Context Canceled", "ms", elapsed.Milliseconds(), "err", context.Cause(ctx), "sql", q)
			return
		} else if !errors.Is(err, context.DeadlineExceeded) {
			q.lggr.Debugw("SQL Context Error", "ms", elapsed.Milliseconds(), "err", err, "sql", q)
			return
		}
		q.lggr.Debugw("SQL Deadline Exceeded", "ms", elapsed.Milliseconds(), "err", err, "sql", q)
	}

	// Success or Deadline Exceeded, so calculate how much of the timeout was used.
	deadline, ok := ctx.Deadline()
	if !ok {
		return
	}
	timeout := deadline.Sub(start)

	var thresholds LogThresholds
	thresholds.FromContextValue(ctx)

	var pct float64 // percent of timeout used
	if timeout > 0 {
		pct = float64(elapsed) / float64(timeout)
		pct *= 100
	} else {
		timeout = 0
		pct = 100
	}

	kvs := []any{"ms", elapsed.Milliseconds(), "timeout", timeout.Milliseconds(), "percent", strconv.FormatFloat(pct, 'f', 1, 64), "sql", q}

	if elapsed >= timeout {
		q.lggr.Criticalw(slowMsg, kvs...)
	} else if errThreshold := thresholds.Error(timeout); errThreshold > 0 && elapsed > errThreshold {
		q.lggr.Errorw(slowMsg, kvs...)
	} else if warnThreshold := thresholds.Warn(timeout); warnThreshold > 0 && elapsed > warnThreshold {
		q.lggr.Warnw(slowMsg, kvs...)
	}

	PromSQLQueryTime.Observe(pct)
}

// LogThresholds holds funcs for computing thresholds for timeout usage.
type LogThresholds struct {
	Warn  func(timeout time.Duration) (threshold time.Duration)
	Error func(timeout time.Duration) (threshold time.Duration)
}

type ctxKeyLogThresholds struct{}

func (t *LogThresholds) ContextWithValue(ctx context.Context) context.Context {
	if t == nil {
		return ctx
	}
	return context.WithValue(ctx, ctxKeyLogThresholds{}, *t)
}

// FromContextValue sets LogThresholds from the context value, if one is present.
// nil fields are set to default functions. 10% for Warn, and 20% for Error.
func (t *LogThresholds) FromContextValue(ctx context.Context) {
	v := ctx.Value(ctxKeyLogThresholds{})
	if v != nil {
		*t = v.(LogThresholds)
	}
	if t.Warn == nil {
		t.Warn = func(timeout time.Duration) time.Duration {
			return timeout / 10
		}
	}
	if t.Error == nil {
		t.Error = func(timeout time.Duration) time.Duration {
			return timeout / 5
		}
	}
}
