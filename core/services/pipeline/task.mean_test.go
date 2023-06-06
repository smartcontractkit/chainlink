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

func TestMeanTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		inputs        []pipeline.Result
		allowedFaults string
		precision     string
		want          pipeline.Result
	}{
		{
			"odd number of inputs",
			[]pipeline.Result{{Value: mustDecimal(t, "1")}, {Value: mustDecimal(t, "2")}, {Value: mustDecimal(t, "3")}},
			"1",
			"",
			pipeline.Result{Value: mustDecimal(t, "2")},
		},
		{
			"even number of inputs",
			[]pipeline.Result{{Value: mustDecimal(t, "1")}, {Value: mustDecimal(t, "2")}, {Value: mustDecimal(t, "3")}, {Value: mustDecimal(t, "4")}},
			"2",
			"",
			pipeline.Result{Value: mustDecimal(t, "2.5")},
		},
		{
			"one input",
			[]pipeline.Result{{Value: mustDecimal(t, "1")}},
			"0",
			"",
			pipeline.Result{Value: mustDecimal(t, "1")},
		},
		{
			"zero inputs",
			[]pipeline.Result{},
			"0",
			"",
			pipeline.Result{Error: pipeline.ErrWrongInputCardinality},
		},
		{
			"fewer errors than threshold",
			[]pipeline.Result{{Error: errors.New("")}, {Value: mustDecimal(t, "2")}, {Value: mustDecimal(t, "3")}, {Value: mustDecimal(t, "4")}},
			"2",
			"",
			pipeline.Result{Value: mustDecimal(t, "3")},
		},
		{
			"exactly threshold of errors",
			[]pipeline.Result{{Error: errors.New("")}, {Error: errors.New("")}, {Value: mustDecimal(t, "3")}, {Value: mustDecimal(t, "4")}},
			"2",
			"",
			pipeline.Result{Value: mustDecimal(t, "3.5")},
		},
		{
			"more errors than threshold",
			[]pipeline.Result{{Error: errors.New("")}, {Error: errors.New("")}, {Error: errors.New("")}, {Value: mustDecimal(t, "4")}},
			"2",
			"",
			pipeline.Result{Error: pipeline.ErrTooManyErrors},
		},
		{
			"(unspecified AllowedFaults) fewer errors than threshold",
			[]pipeline.Result{{Error: errors.New("")}, {Error: errors.New("")}, {Value: mustDecimal(t, "3")}, {Value: mustDecimal(t, "4")}},
			"",
			"",
			pipeline.Result{Value: mustDecimal(t, "3.5")},
		},
		{
			"(unspecified AllowedFaults) exactly threshold of errors",
			[]pipeline.Result{{Error: errors.New("")}, {Error: errors.New("")}, {Error: errors.New("")}, {Value: mustDecimal(t, "4")}},
			"",
			"",
			pipeline.Result{Value: mustDecimal(t, "4")},
		},
		{
			"(unspecified AllowedFaults) more errors than threshold",
			[]pipeline.Result{{Error: errors.New("")}, {Error: errors.New("")}, {Error: errors.New("")}},
			"",
			"",
			pipeline.Result{Error: pipeline.ErrTooManyErrors},
		},
		{
			"precision",
			[]pipeline.Result{{Value: mustDecimal(t, "1.234")}, {Value: mustDecimal(t, "2.345")}, {Value: mustDecimal(t, "3.456")}, {Value: mustDecimal(t, "4.567")}},
			"1",
			"2",
			pipeline.Result{Value: mustDecimal(t, "2.90")},
		},
		{
			"precision (> 16)",
			[]pipeline.Result{{Value: mustDecimal(t, "1.11111111111111111111")}, {Value: mustDecimal(t, "3.33333333333333333333")}},
			"1",
			"18",
			pipeline.Result{Value: mustDecimal(t, "2.222222222222222222")},
		},
		{
			"precision (negative)",
			[]pipeline.Result{{Value: mustDecimal(t, "12.34")}, {Value: mustDecimal(t, "23.45")}, {Value: mustDecimal(t, "34.56")}, {Value: mustDecimal(t, "45.67")}},
			"1",
			"-1",
			pipeline.Result{Value: mustDecimal(t, "30")},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Run("without vars", func(t *testing.T) {
				task := pipeline.MeanTask{
					BaseTask:      pipeline.NewBaseTask(0, "task", nil, nil, 0),
					AllowedFaults: test.allowedFaults,
					Precision:     test.precision,
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
				task := pipeline.MeanTask{
					BaseTask:      pipeline.NewBaseTask(0, "task", nil, nil, 0),
					Values:        "$(foo.bar)",
					AllowedFaults: test.allowedFaults,
					Precision:     test.precision,
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
				case 2:
					valuesParam = "[ $(foo), $(bar) ]"
					vars = pipeline.NewVarsFrom(map[string]interface{}{"foo": inputs[0], "bar": inputs[1]})
				case 3:
					valuesParam = "[ $(foo), $(bar), $(chain) ]"
					vars = pipeline.NewVarsFrom(map[string]interface{}{"foo": inputs[0], "bar": inputs[1], "chain": inputs[2]})
				case 4:
					valuesParam = "[ $(foo), $(bar), $(chain), $(link) ]"
					vars = pipeline.NewVarsFrom(map[string]interface{}{"foo": inputs[0], "bar": inputs[1], "chain": inputs[2], "link": inputs[3]})
				}

				task := pipeline.MeanTask{
					BaseTask:      pipeline.NewBaseTask(0, "task", nil, nil, 0),
					Values:        valuesParam,
					AllowedFaults: test.allowedFaults,
					Precision:     test.precision,
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
