package logger

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var _ gormlogger.Interface = &gormWrapper{}

type gormWrapper struct {
	Logger
	logAllQueries bool
	slowThreshold time.Duration
}

// Noop
func (o *gormWrapper) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return o
}

func (o *gormWrapper) Info(ctx context.Context, s string, i ...interface{}) {
	o.Logger.Infow(fmt.Sprintf(s, i...))
}

func (o *gormWrapper) Warn(ctx context.Context, s string, i ...interface{}) {
	o.Logger.Warnw(fmt.Sprintf(s, i...))
}

func (o *gormWrapper) Error(ctx context.Context, s string, i ...interface{}) {
	o.Logger.Errorw(fmt.Sprintf(s, i...))
}

// This is called at the end of every gorm v2 query.
// We always log the sql queries for errors and slow queries (warns).
// Need to set LOG_SQL=true to enable all queries.
func (o *gormWrapper) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	switch {
	case ctx.Err() != nil:
		sql, _ := fc()
		o.Logger.Debugw("Operation cancelled via context", "err", err, "elapsed", float64(elapsed.Nanoseconds())/1e6, "sql", sql)
	case err != nil:
		// NOTE: Silence "record not found" errors since it is the one type of
		// error that we expect/handle and otherwise it fills our logs with
		// noise
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}
		sql, rows := fc()
		// We only log these as debugs as we expect the errors are handled at a higher
		// level.
		if rows == -1 {
			o.Logger.Debugw("Operation failed", "err", err, "elapsed", float64(elapsed.Nanoseconds())/1e6, "sql", sql)
		} else {
			o.Logger.Debugw("Operation failed", "err", err, "elapsed", float64(elapsed.Nanoseconds())/1e6, "rows", rows, "sql", sql)
		}
	case elapsed > o.slowThreshold && o.slowThreshold != 0:
		sql, rows := fc()
		if rows == -1 {
			o.Logger.Warnw(fmt.Sprintf("SQL query took longer than %s", o.slowThreshold), "elapsed", float64(elapsed.Nanoseconds())/1e6, "sql", sql)
		} else {
			o.Logger.Warnw(fmt.Sprintf("SQL query took longer than %s", o.slowThreshold), "elapsed", float64(elapsed.Nanoseconds())/1e6, "rows", rows, "sql", sql)
		}
	case o.logAllQueries:
		sql, rows := fc()
		if rows == -1 {
			o.Logger.Debugw("Query executed", "elapsed", float64(elapsed.Nanoseconds())/1e6, "sql", sql)
		} else {
			o.Logger.Debugw("Query executed", "elapsed", float64(elapsed.Nanoseconds())/1e6, "rows", rows, "sql", sql)
		}
	}
}

func NewGormWrapper(logger Logger, logAllQueries bool, slowThreshold time.Duration) *gormWrapper {
	return &gormWrapper{logger.withCallerSkip(2), logAllQueries, slowThreshold}
}
