package orm

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var _ gormlogger.Interface = &ormLogWrapper{}

type ormLogWrapper struct {
	logger.Logger
	logAllQueries bool
	slowThreshold time.Duration
}

// Noop
func (o *ormLogWrapper) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return o
}

func (o *ormLogWrapper) Info(ctx context.Context, s string, i ...interface{}) {
	o.Infow(fmt.Sprintf(s, i...))
}

func (o *ormLogWrapper) Warn(ctx context.Context, s string, i ...interface{}) {
	o.Warnw(fmt.Sprintf(s, i...))
}

func (o *ormLogWrapper) Error(ctx context.Context, s string, i ...interface{}) {
	o.Errorw(fmt.Sprintf(s, i...))
}

// This is called at the end of every gorm v2 query.
// We always log the sql queries for errors and slow queries (warns).
// Need to set LOG_SQL=true to enable all queries.
func (o *ormLogWrapper) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	switch {
	case ctx.Err() != nil:
		sql, _ := fc()
		o.Debugw("Operation cancelled via context", "err", err, "elapsed", float64(elapsed.Nanoseconds())/1e6, "sql", sql)
	case err != nil:
		// NOTE: Silence "record not found" errors since it is the one type of
		// error that we expect/handle and otherwise it fills our logs with
		// noise
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}
		sql, rows := fc()
		if rows == -1 {
			o.Errorw("Operation failed", "err", err, "elapsed", float64(elapsed.Nanoseconds())/1e6, "sql", sql)
		} else {
			o.Errorw("Operation failed", "err", err, "elapsed", float64(elapsed.Nanoseconds())/1e6, "rows", rows, "sql", sql)
		}
	case elapsed > o.slowThreshold && o.slowThreshold != 0:
		sql, rows := fc()
		if rows == -1 {
			o.Warnw(fmt.Sprintf("SQL query took longer than %s", o.slowThreshold), "elapsed", float64(elapsed.Nanoseconds())/1e6, "sql", sql)
		} else {
			o.Warnw(fmt.Sprintf("SQL query took longer than %s", o.slowThreshold), "elapsed", float64(elapsed.Nanoseconds())/1e6, "rows", rows, "sql", sql)
		}
	case o.logAllQueries:
		sql, rows := fc()
		if rows == -1 {
			o.Debugw("Query executed", "elapsed", float64(elapsed.Nanoseconds())/1e6, "sql", sql)
		} else {
			o.Debugw("Query executed", "elapsed", float64(elapsed.Nanoseconds())/1e6, "rows", rows, "sql", sql)
		}
	}
}

func newOrmLogWrapper(logger logger.Logger, logAllQueries bool, slowThreshold time.Duration) *ormLogWrapper {
	// newLogger := logger.
	// 	SugaredLogger.
	// 	Desugar().
	// 	WithOptions(zap.AddCallerSkip(2)).
	// 	Sugar()
	return &ormLogWrapper{logger, logAllQueries, slowThreshold}
}
