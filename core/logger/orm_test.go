package logger_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/stretchr/testify/assert"
)

func TestLoggerORM(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
	orm := logger.NewORM(db, lggr)

	_, ok := orm.GetServiceLogLevel("foo")
	assert.False(t, ok)

	err := orm.SetServiceLogLevel(testutils.Context(t), "foo", "debug")
	assert.NoError(t, err)

	lvl, ok := orm.GetServiceLogLevel("foo")
	assert.True(t, ok)
	assert.Equal(t, "debug", lvl)

	// overwrite
	err = orm.SetServiceLogLevel(testutils.Context(t), "foo", "info")
	assert.NoError(t, err)

	lvl, ok = orm.GetServiceLogLevel("foo")
	assert.True(t, ok)
	assert.Equal(t, "info", lvl)
}
