package sqlutil

import (
	"context"
	"database/sql/driver"
	"strconv"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func Test_sprintQ(t *testing.T) {
	for _, tt := range []struct {
		name  string
		query string
		args  []interface{}
		exp   string
	}{
		{"none",
			"SELECT * FROM table;",
			nil,
			"SELECT * FROM table;"},
		{"one",
			"SELECT $1 FROM table;",
			[]interface{}{"foo"},
			"SELECT 'foo' FROM table;"},
		{"two",
			"SELECT $1 FROM table WHERE bar = $2;",
			[]interface{}{"foo", 1},
			"SELECT 'foo' FROM table WHERE bar = 1;"},
		{"limit",
			"SELECT $1 FROM table LIMIT $2;",
			[]interface{}{"foo", limit(10)},
			"SELECT 'foo' FROM table LIMIT 10;"},
		{"limit-all",
			"SELECT $1 FROM table LIMIT $2;",
			[]interface{}{"foo", limit(-1)},
			"SELECT 'foo' FROM table LIMIT NULL;"},
		{"bytea",
			"SELECT $1 FROM table WHERE b = $2;",
			[]interface{}{"foo", []byte{0x0a}},
			"SELECT 'foo' FROM table WHERE b = '\\x0a';"},
		{"bytea[]",
			"SELECT $1 FROM table WHERE b = $2;",
			[]interface{}{"foo", pq.ByteaArray([][]byte{{0xa}, {0xb}})},
			"SELECT 'foo' FROM table WHERE b = ARRAY['\\x0a','\\x0b'];"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := sprintQ(tt.query, tt.args)
			t.Log(tt.query, tt.args)
			t.Log(got)
			require.Equal(t, tt.exp, got)
		})
	}
}

var _ driver.Valuer = limit(-1)

// limit is a helper driver.Valuer for LIMIT queries which uses nil/NULL for negative values.
type limit int

func (l limit) String() string {
	if l < 0 {
		return "NULL"
	}
	return strconv.Itoa(int(l))
}

func (l limit) Value() (driver.Value, error) {
	if l < 0 {
		return nil, nil
	}
	return l, nil
}

func Test_queryLogger_logTiming(t *testing.T) {
	for _, tt := range []struct {
		name       string
		thresholds *LogThresholds
	}{
		{"default", nil},
		{"custom", &LogThresholds{
			Warn: func(timeout time.Duration) time.Duration {
				return timeout / 2 // 50%
			},
			Error: func(timeout time.Duration) time.Duration {
				return timeout - (timeout / 10) // 90%
			},
		}},
		{"partial", &LogThresholds{
			Error: func(timeout time.Duration) time.Duration {
				return timeout - (timeout / 4) // 75%
			},
		}},
	} {
		t.Run(tt.name, func(t *testing.T) {

			t.Run("no-deadline", func(t *testing.T) {
				lggr, ol := logger.TestObservedSugared(t, zap.DebugLevel)
				ql := newQueryLogger(lggr, "TEST QUERY", "foo", "bar")

				start := time.Now().Add(-time.Second)

				ctx, cancel := context.WithCancel(tests.Context(t))
				if tt.thresholds != nil {
					ctx = tt.thresholds.ContextWithValue(ctx)
				}
				ql.logTiming(ctx, start)
				cancel()

				// no logs
				logs := ol.TakeAll()
				if !assert.Empty(t, logs) {
					t.Logf("Unexpected logs: %v", logs)
				}
			})

			t.Run("fast", func(t *testing.T) {
				lggr, ol := logger.TestObservedSugared(t, zap.DebugLevel)
				ql := newQueryLogger(lggr, "TEST QUERY", "foo", "bar")

				start := time.Now().Add(-time.Second)

				ctx, cancel := context.WithTimeout(tests.Context(t), time.Minute)
				ctx = tt.thresholds.ContextWithValue(ctx)
				ql.logTiming(ctx, start)
				cancel()

				// no logs
				logs := ol.TakeAll()
				if !assert.Empty(t, logs) {
					t.Logf("Unexpected logs: %v", logs)
				}
			})

			t.Run("warn", func(t *testing.T) {
				lggr, ol := logger.TestObservedSugared(t, zap.DebugLevel)
				ql := newQueryLogger(lggr, "TEST QUERY", "foo", "bar")

				threshold := time.Second // default is 10%
				if tt.thresholds != nil && tt.thresholds.Warn != nil {
					threshold = tt.thresholds.Warn(10 * time.Second)
				}
				start := time.Now().Add(-threshold)
				deadline := time.Now().Add(10*time.Second - threshold)

				ctx, cancel := context.WithDeadline(tests.Context(t), deadline)
				ctx = tt.thresholds.ContextWithValue(ctx)
				ql.logTiming(ctx, start)
				cancel()

				// warning
				for _, l := range ol.All() {
					assert.LessOrEqual(t, l.Level, zap.WarnLevel, "unexpected log message: %v", l)
				}
				logs := ol.FilterMessageSnippet(slowMsg).TakeAll()
				require.Len(t, logs, 1)
				log := logs[0]
				assert.Equal(t, zap.WarnLevel, log.Level)
				assert.Equal(t, "TEST QUERY", log.ContextMap()["sql"])

			})

			t.Run("error", func(t *testing.T) {
				lggr, ol := logger.TestObservedSugared(t, zap.DebugLevel)
				ql := newQueryLogger(lggr, "TEST QUERY", "foo", "bar")

				threshold := time.Second // default is 20%
				if tt.thresholds != nil && tt.thresholds.Error != nil {
					threshold = tt.thresholds.Error(5 * time.Second)
				}
				start := time.Now().Add(-threshold)
				deadline := time.Now().Add(5*time.Second - threshold)

				ctx, cancel := context.WithDeadline(tests.Context(t), deadline)
				if tt.thresholds != nil {
					ctx = tt.thresholds.ContextWithValue(ctx)
				}
				ql.logTiming(ctx, start)
				cancel()

				// error
				for _, l := range ol.All() {
					assert.LessOrEqual(t, l.Level, zap.ErrorLevel, "unexpected log message: %v", l)
				}
				logs := ol.FilterMessageSnippet(slowMsg).TakeAll()
				require.Len(t, logs, 1)
				log := logs[0]
				assert.Equal(t, zap.ErrorLevel, log.Level)
				assert.Equal(t, "TEST QUERY", log.ContextMap()["sql"])
			})

			t.Run("critical", func(t *testing.T) {
				lggr, ol := logger.TestObservedSugared(t, zap.DebugLevel)
				ql := newQueryLogger(lggr, "TEST QUERY", "foo", "bar")

				// >100%
				start := time.Now().Add(-10 * time.Second)
				deadline := time.Now()

				ctx, cancel := context.WithDeadline(tests.Context(t), deadline)
				if tt.thresholds != nil {
					ctx = tt.thresholds.ContextWithValue(ctx)
				}
				ql.logTiming(ctx, start)
				cancel()

				// critical
				for _, l := range ol.All() {
					assert.LessOrEqual(t, l.Level, zap.DPanicLevel, "unexpected log message: %v", l)
				}
				logs := ol.FilterMessageSnippet(slowMsg).TakeAll()
				require.Len(t, logs, 1)
				log := logs[0]
				assert.Equal(t, zap.DPanicLevel, log.Level)
				assert.Equal(t, "TEST QUERY", log.ContextMap()["sql"])

				logs = ol.FilterMessageSnippet("SQL Deadline Exceeded").TakeAll()
				require.Len(t, logs, 1)
				log = logs[0]
				assert.Equal(t, zap.DebugLevel, log.Level)
				assert.Equal(t, "TEST QUERY", log.ContextMap()["sql"])
			})

			t.Run("cancelled", func(t *testing.T) {
				lggr, ol := logger.TestObservedSugared(t, zap.DebugLevel)
				ql := newQueryLogger(lggr, "TEST QUERY", "foo", "bar")

				// >100%
				start := time.Now()

				ctx, cancel := context.WithCancel(tests.Context(t))
				cancel() // pre-cancel
				if tt.thresholds != nil {
					ctx = tt.thresholds.ContextWithValue(ctx)
				}
				ql.logTiming(ctx, start)

				// debug
				for _, l := range ol.All() {
					assert.LessOrEqual(t, l.Level, zap.DebugLevel, "unexpected log message: %v", l)
				}

				require.Empty(t, ol.FilterMessageSnippet(slowMsg).TakeAll())

				logs := ol.FilterMessageSnippet("SQL Context Canceled").TakeAll()
				require.Len(t, logs, 1)
				log := logs[0]
				assert.Equal(t, zap.DebugLevel, log.Level)
				assert.Equal(t, "TEST QUERY", log.ContextMap()["sql"])
			})

			t.Run("deadline-before", func(t *testing.T) {
				lggr, ol := logger.TestObservedSugared(t, zap.DebugLevel)
				ql := newQueryLogger(lggr, "TEST QUERY", "foo", "bar")

				// >100%
				start := time.Now()

				ctx, cancel := context.WithDeadline(tests.Context(t), start.Add(-time.Second))
				defer cancel()
				if tt.thresholds != nil {
					ctx = tt.thresholds.ContextWithValue(ctx)
				}
				ql.logTiming(ctx, start)

				// debug
				for _, l := range ol.All() {
					assert.LessOrEqual(t, l.Level, zap.DebugLevel, "unexpected log message: %v", l)
				}

				require.Empty(t, ol.FilterMessageSnippet(slowMsg).TakeAll())

				logs := ol.FilterMessageSnippet("SQL Deadline Exceeded").TakeAll()
				require.Len(t, logs, 1)
				log := logs[0]
				assert.Equal(t, zap.DebugLevel, log.Level)
				assert.Equal(t, "TEST QUERY", log.ContextMap()["sql"])
			})
		})
	}
}
