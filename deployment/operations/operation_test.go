package operations

import (
	"context"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

type OpDeps struct{}

type OpInput struct {
	A int
	B int
}

func Test_NewOperation(t *testing.T) {
	t.Parallel()

	version := semver.MustParse("1.0.0")
	description := "test operation"
	handler := func(b Bundle, deps OpDeps, input OpInput) (output int, err error) {
		return input.A + input.B, nil
	}

	op := NewOperation("sum", version, description, handler)

	assert.Equal(t, "sum", op.ID())
	assert.Equal(t, version.String(), op.Version())
	assert.Equal(t, description, op.Description())
	res, err := op.handler(Bundle{}, OpDeps{}, OpInput{1, 2})
	require.NoError(t, err)
	assert.Equal(t, 3, res)
}

func Test_Operation_Execute(t *testing.T) {
	t.Parallel()

	version := semver.MustParse("1.0.0")
	description := "test operation"
	log, observedLog := logger.TestObserved(t, zapcore.InfoLevel)

	// simulate an addition operation
	handler := func(b Bundle, deps OpDeps, input OpInput) (output int, err error) {
		return input.A + input.B, nil
	}

	op := NewOperation("sum", version, description, handler)
	e := NewBundle(context.Background, log)
	input := OpInput{
		A: 1,
		B: 2,
	}

	output, err := op.execute(e, OpDeps{}, input)

	require.NoError(t, err)
	assert.Equal(t, 3, output)

	require.Equal(t, 1, observedLog.Len())
	entry := observedLog.All()[0]
	assert.Equal(t, "Executing operation", entry.Message)
	assert.Equal(t, "sum", entry.ContextMap()["id"])
	assert.Equal(t, version.String(), entry.ContextMap()["version"])
	assert.Equal(t, description, entry.ContextMap()["description"])
}
