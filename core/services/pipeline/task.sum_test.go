package pipeline_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func TestSumTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		inputs        []pipeline.Result
		allowedFaults string
		want          pipeline.Result
	}{
		{
			"happy",
			[]pipeline.Result{{Value: mustDecimal(t, "1")}, {Value: mustDecimal(t, "2")}, {Value: mustDecimal(t, "3")}},
			"1",
			pipeline.Result{Value: mustDecimal(t, "6")},
		},
		{
			"happy (one input)",
			[]pipeline.Result{{Value: mustDecimal(t, "1")}},
			"0",
			pipeline.Result{Value: mustDecimal(t, "1")},
		},
		{
			"zero inputs",
			[]pipeline.Result{},
			"0",
			pipeline.Result{Error: pipeline.ErrWrongInputCardinality},
		},
		{
			"fewer errors than threshold",
			[]pipeline.Result{{Error: errors.New("")}, {Value: mustDecimal(t, "2")}, {Value: mustDecimal(t, "3")}, {Value: mustDecimal(t, "4")}},
			"2",
			pipeline.Result{Value: mustDecimal(t, "9")},
		},
		{
			"exactly threshold of errors",
			[]pipeline.Result{{Error: errors.New("")}, {Error: errors.New("")}, {Value: mustDecimal(t, "3")}, {Value: mustDecimal(t, "4")}},
			"2",
			pipeline.Result{Value: mustDecimal(t, "7")},
		},
		{
			"more errors than threshold",
			[]pipeline.Result{{Error: errors.New("")}, {Error: errors.New("")}, {Error: errors.New("")}, {Value: mustDecimal(t, "4")}},
			"2",
			pipeline.Result{Error: pipeline.ErrTooManyErrors},
		},
		{
			"(unspecified AllowedFaults) fewer errors than threshold",
			[]pipeline.Result{{Error: errors.New("")}, {Error: errors.New("")}, {Value: mustDecimal(t, "3")}, {Value: mustDecimal(t, "4")}},
			"",
			pipeline.Result{Value: mustDecimal(t, "7")},
		},
		{
			"(unspecified AllowedFaults) exactly threshold of errors",
			[]pipeline.Result{{Error: errors.New("")}, {Error: errors.New("")}, {Error: errors.New("")}, {Value: mustDecimal(t, "4")}},
			"",
			pipeline.Result{Value: mustDecimal(t, "4")},
		},
		{
			"(unspecified AllowedFaults) more errors than threshold",
			[]pipeline.Result{{Error: errors.New("")}, {Error: errors.New("")}, {Error: errors.New("")}},
			"",
			pipeline.Result{Error: pipeline.ErrTooManyErrors},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Run("without vars", func(t *testing.T) {
				task := pipeline.SumTask{
					BaseTask:      pipeline.NewBaseTask(0, "task", nil, nil, 0),
					AllowedFaults: test.allowedFaults,
				}
				output, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), pipeline.NewVarsFrom(nil), test.inputs)
				assert.False(t, runInfo.IsPending)
				assert.False(t, runInfo.IsRetryable)
				if output.Error != nil {
					require.Equal(t, test.want.Error, errors.Cause(output.Error))
					require.Nil(t, output.Value)
				} else {
					require.Equal(t, test.want.Value.(*decimal.Decimal).String(), output.Value.(decimal.Decimal).String())
					require.NoError(t, output.Error)
				}
			})
			t.Run("with vars", func(t *testing.T) {
				var inputs []interface{}
				for _, input := range test.inputs {
					if input.Error != nil {
						inputs = append(inputs, input.Error)
					} else {
						inputs = append(inputs, input.Value)
					}
				}
				vars := pipeline.NewVarsFrom(map[string]interface{}{
					"foo": map[string]interface{}{"bar": inputs},
				})
				task := pipeline.SumTask{
					BaseTask:      pipeline.NewBaseTask(0, "task", nil, nil, 0),
					Values:        "$(foo.bar)",
					AllowedFaults: test.allowedFaults,
				}
				output, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), vars, nil)
				assert.False(t, runInfo.IsPending)
				assert.False(t, runInfo.IsRetryable)

				if output.Error != nil {
					require.Equal(t, test.want.Error, errors.Cause(output.Error))
					require.Nil(t, output.Value)
				} else {
					require.Equal(t, test.want.Value.(*decimal.Decimal).String(), output.Value.(decimal.Decimal).String())
					require.NoError(t, output.Error)
				}
			})
			t.Run("with json vars", func(t *testing.T) {
				var inputs []interface{}
				for _, input := range test.inputs {
					if input.Error != nil {
						inputs = append(inputs, input.Error)
					} else {
						inputs = append(inputs, input.Value)
					}
				}
				var valuesParam string
				var vars pipeline.Vars
				switch len(inputs) {
				case 0:
					valuesParam = "[]"
					vars = pipeline.NewVarsFrom(nil)
				case 1:
					valuesParam = "[ $(foo) ]"
					vars = pipeline.NewVarsFrom(map[string]interface{}{"foo": inputs[0]})
				case 3:
					valuesParam = "[ $(foo), $(bar), $(chain) ]"
					vars = pipeline.NewVarsFrom(map[string]interface{}{"foo": inputs[0], "bar": inputs[1], "chain": inputs[2]})
				case 4:
					valuesParam = "[ $(foo), $(bar), $(chain), $(link) ]"
					vars = pipeline.NewVarsFrom(map[string]interface{}{"foo": inputs[0], "bar": inputs[1], "chain": inputs[2], "link": inputs[3]})
				}

				task := pipeline.SumTask{
					BaseTask:      pipeline.NewBaseTask(0, "task", nil, nil, 0),
					Values:        valuesParam,
					AllowedFaults: test.allowedFaults,
				}
				output, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), vars, nil)
				assert.False(t, runInfo.IsPending)
				assert.False(t, runInfo.IsRetryable)
				if output.Error != nil {
					require.Equal(t, test.want.Error, errors.Cause(output.Error))
					require.Nil(t, output.Value)
				} else {
					require.Equal(t, test.want.Value.(*decimal.Decimal).String(), output.Value.(decimal.Decimal).String())
					require.NoError(t, output.Error)
				}
			})
		})
	}
}
