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

func TestCoalesceTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		inputs []pipeline.Result
		want   pipeline.Result
	}{
		{
			"zero inputs",
			[]pipeline.Result{},
			pipeline.Result{Value: nil},
		},
		{
			"one non-errored decimal input",
			[]pipeline.Result{{Value: mustDecimal(t, "42")}},
			pipeline.Result{Value: mustDecimal(t, "42")},
		},
		{
			"one errored decimal input",
			[]pipeline.Result{{Value: mustDecimal(t, "42"), Error: errors.New("foo")}},
			pipeline.Result{Value: nil},
		},
		{
			"one non-errored string input",
			[]pipeline.Result{{Value: "42"}},
			pipeline.Result{Value: "42"},
		},
		{
			"one errored input and one non-errored input",
			[]pipeline.Result{{Error: errors.New("foo"), Value: "1"}, {Value: "42"}},
			pipeline.Result{Value: "42"},
		},
		{
			"one non-errored input and one errored input",
			[]pipeline.Result{{Value: "42"}, {Error: errors.New("foo"), Value: "1"}},
			pipeline.Result{Value: "42"},
		},
		{
			"two errored inputs",
			[]pipeline.Result{{Value: "42", Error: errors.New("bar")}, {Error: errors.New("foo"), Value: "1"}},
			pipeline.Result{Value: nil},
		},
		{
			"one errored input and two non-errored inputs",
			[]pipeline.Result{{Error: errors.New("foo")}, {Value: "42"}, {Value: "44"}},
			pipeline.Result{Value: "42"},
		},
		{
			"one nil input",
			[]pipeline.Result{{Value: nil}},
			pipeline.Result{Value: nil},
		},
		{
			"two errored or nil inputs",
			[]pipeline.Result{{Error: errors.New("foo")}, {Value: nil}},
			pipeline.Result{Value: nil},
		},
		{
			"one nil input and two non-errored inputs",
			[]pipeline.Result{{Value: nil}, {Value: "42"}, {Value: "44"}},
			pipeline.Result{Value: "42"},
		},
		{
			"two non-errored inputs and one nil input",
			[]pipeline.Result{{Value: "42"}, {Value: nil}, {Value: "44"}},
			pipeline.Result{Value: "42"},
		},
		{
			"three errored or nil inputs and one decimal input",
			[]pipeline.Result{{Error: errors.New("foo")}, {Value: nil}, {Error: errors.New("bar")}, {Value: mustDecimal(t, "42")}},
			pipeline.Result{Value: mustDecimal(t, "42")},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			task := pipeline.CoalesceTask{}
			output, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), pipeline.NewVarsFrom(nil), test.inputs)
			assert.False(t, runInfo.IsPending)
			assert.False(t, runInfo.IsRetryable)
			if output.Error != nil {
				require.Equal(t, test.want.Error, errors.Cause(output.Error))
				require.Nil(t, output.Value)
			} else {
				switch val := test.want.Value.(type) {
				case *decimal.Decimal:
					require.Equal(t, val.String(), output.Value.(*decimal.Decimal).String())
				default:
					require.Equal(t, val, output.Value)
				}
				require.NoError(t, output.Error)
			}
		})
	}
}
