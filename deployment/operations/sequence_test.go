package operations

import (
	"context"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func Test_NewSequence(t *testing.T) {
	t.Parallel()

	version := semver.MustParse("1.0.0")

	description := "test operation"
	handler := func(e Bundle, deps OpDeps, input OpInput) (output int, err error) {
		return input.A + input.B, nil
	}
	op := NewOperation("sum", version, description, handler)

	seqDescription := "test sequence"
	sequence := NewSequence("seq-sum", version, seqDescription, func(env Bundle, deps any, input OpInput) (int, error) {
		operationRes, err := ExecuteOperation(env, op, OpDeps{}, input)
		if err != nil {
			return 0, err
		}
		return operationRes.Output, nil
	})

	assert.Equal(t, "seq-sum", sequence.ID())
	assert.Equal(t, version.String(), sequence.Version())
	assert.Equal(t, seqDescription, sequence.Description())
	e := NewBundle(context.Background, logger.Test(t), NewMemoryReporter())
	res, err := sequence.handler(e, OpDeps{}, OpInput{1, 2})
	require.NoError(t, err)
	assert.Equal(t, 3, res)
}
