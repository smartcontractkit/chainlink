package pipeline_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

func TestAnyTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		inputs []pipeline.Result
		want   pipeline.Result
	}{
		{
			"zero inputs",
			[]pipeline.Result{},
			pipeline.Result{Error: pipeline.ErrWrongInputCardinality},
		},
		{
			"one non-errored decimal input",
			[]pipeline.Result{{Value: mustDecimal(t, "42")}},
			pipeline.Result{Value: mustDecimal(t, "42")},
		},
		{
			"one errored decimal input",
			[]pipeline.Result{{Value: mustDecimal(t, "42"), Error: errors.New("foo")}},
			pipeline.Result{Error: pipeline.ErrBadInput},
		},
		{
			"one non-errored string input",
			[]pipeline.Result{{Value: "42"}},
			pipeline.Result{Value: "42"},
		},
		{
			"one errored input and one non-errored input",
			[]pipeline.Result{{Value: "42"}, {Error: errors.New("foo"), Value: "1"}},
			pipeline.Result{Value: "42"},
		},
		{
			"two errored inputs",
			[]pipeline.Result{{Value: "42", Error: errors.New("bar")}, {Error: errors.New("foo"), Value: "1"}},
			pipeline.Result{Error: pipeline.ErrBadInput},
		},
		{
			"two non-errored inputs with one errored input",
			[]pipeline.Result{{Value: "42"}, {Value: "42"}, {Error: errors.New("foo")}},
			pipeline.Result{Value: "42"},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			task := pipeline.AnyTask{}
			output, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), pipeline.NewVarsFrom(nil), test.inputs)
			assert.False(t, runInfo.IsPending)
			assert.False(t, runInfo.IsRetryable)
			if output.Error != nil {
				require.Equal(t, test.want.Error, errors.Cause(output.Error))
				require.Nil(t, output.Value)
			} else {
				switch test.want.Value.(type) {
				case *decimal.Decimal:
					require.Equal(t, test.want.Value.(*decimal.Decimal).String(), output.Value.(*decimal.Decimal).String())
				default:
					require.Equal(t, test.want.Value, output.Value)
				}
				require.NoError(t, output.Error)
			}
		})
	}
}
