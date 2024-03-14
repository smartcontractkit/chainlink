package sqlutil

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

const (
	getDur = 10 * time.Millisecond
	selDur = 200 * time.Millisecond
)

func TestWrapDataSource(t *testing.T) {
	lggr, ol := logger.TestObserved(t, zapcore.InfoLevel)
	var ds DataSource = &dataSource{}
	var sentinelErr = errors.New("intercepted error")
	const fakeError = "fake warning"
	ds = WrapDataSource(ds, lggr, TimeoutHook(selDur/2), noopHook, MonitorHook(func() bool { return true }), noopHook, func(ctx context.Context, lggr logger.Logger, do func(context.Context) error, query string, args ...any) error {
		err := do(ctx)
		if err != nil {
			return err
		}
		lggr.Error(fakeError)
		return sentinelErr
	})
	ctx := tests.Context(t)

	// Error intercepted
	err := ds.GetContext(ctx, "test", "foo", 42, "bar")
	_, file, line, ok := runtime.Caller(0)
	require.True(t, ok)
	expCaller := fmt.Sprintf("%s:%d", file, line-1)
	require.ErrorIs(t, err, sentinelErr)
	logs := ol.FilterMessage(slowMsg).All()
	require.Len(t, logs, 1)
	assert.Equal(t, zapcore.WarnLevel, logs[0].Level)
	assert.Equal(t, expCaller, logs[0].Caller.String())
	logs = ol.FilterMessage(fakeError).All()
	require.Len(t, logs, 1)
	assert.Equal(t, zapcore.ErrorLevel, logs[0].Level)
	assert.Equal(t, expCaller, logs[0].Caller.String())
	_ = ol.TakeAll()

	// Timeout applied
	err = ds.SelectContext(ctx, "test", "foo", 42, "bar")
	require.ErrorIs(t, err, context.DeadlineExceeded)
	logs = ol.FilterMessage(slowMsg).All()
	require.Len(t, logs, 1)
	assert.Equal(t, zapcore.DPanicLevel, logs[0].Level)
	_ = ol.TakeAll()

	// Without default timeout
	err = ds.SelectContext(WithoutDefaultTimeout(ctx), "test", "foo", 42, "bar")
	require.ErrorIs(t, err, sentinelErr)

	// W/o default, but with our own
	ctx2, cancel := context.WithTimeout(WithoutDefaultTimeout(ctx), selDur/100)
	t.Cleanup(cancel)
	err = ds.SelectContext(ctx2, "test", "foo", 42, "bar")
	require.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestWrapDataSource_transactional(t *testing.T) {
	lggr := logger.Test(t)

	txional := (*transactional)(nil)

	var ds DataSource = (*sqlx.DB)(nil)
	assert.Implements(t, txional, ds)
	got := WrapDataSource(ds, lggr)
	assert.Implements(t, txional, got)
	got = WrapDataSource(ds, lggr, noopHook)
	assert.Implements(t, txional, got)
	got = WrapDataSource(ds, lggr, noopHook, noopHook)
	assert.Implements(t, txional, got)

	ds = (*sqlx.Tx)(nil)
	assert.NotImplements(t, txional, ds)
	got = WrapDataSource(ds, lggr)
	assert.NotImplements(t, txional, got)
	got = WrapDataSource(ds, lggr, noopHook)
	assert.NotImplements(t, txional, got)
	got = WrapDataSource(ds, lggr, noopHook, noopHook)
	assert.NotImplements(t, txional, got)
}

var _ DataSource = &dataSource{}

type dataSource struct{}

func (q *dataSource) DriverName() string { return "" }

func (q *dataSource) Rebind(s string) string { return "" }

func (q *dataSource) BindNamed(s string, i interface{}) (string, []interface{}, error) {
	return "", nil, nil
}

func (q *dataSource) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

func (q *dataSource) QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	return nil, nil
}

func (q *dataSource) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	return nil
}

func (q *dataSource) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}

func (q *dataSource) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return nil, nil
}

func (q *dataSource) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(getDur):
	}
	return nil
}

func (q *dataSource) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(selDur):
	}
	return nil
}
