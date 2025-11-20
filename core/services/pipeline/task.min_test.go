package pipeline_test

import (
	"strconv"
	"testing"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

func TestMinTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		inputs        []pipeline.Result
		allowedFaults string
		lax           bool
		want          pipeline.Result
	}{
		{
			"happy",
			[]pipeline.Result{{Value: mustDecimal(t, "1")}, {Value: mustDecimal(t, "2")}, {Value: mustDecimal(t, "3")}},
			"1",
			false,
			pipeline.Result{Value: mustDecimal(t, "1")},
		},
		{
			"happy (one input)",
			[]pipeline.Result{{Value: mustDecimal(t, "1")}},
			"0",
			false,
			pipeline.Result{Value: mustDecimal(t, "1")},
		},
		{
			"happy (with zero)",
			[]pipeline.Result{{Value: mustDecimal(t, "1")}, {Value: mustDecimal(t, "0")}},
			"1",
			false,
			pipeline.Result{Value: mustDecimal(t, "0")},
		},
		{
			"happy (with negative)",
			[]pipeline.Result{{Value: mustDecimal(t, "1")}, {Value: mustDecimal(t, "-1")}},
			"1",
			false,
			pipeline.Result{Value: mustDecimal(t, "-1")},
		},
		{
			"happy (with fractional)",
			[]pipeline.Result{{Value: mustDecimal(t, "0.2")}, {Value: mustDecimal(t, "0.1")}},
			"1",
			false,
			pipeline.Result{Value: mustDecimal(t, "0.1")},
		},
		{
			"nil and non-nil inputs",
			[]pipeline.Result{{Value: mustDecimal(t, "1")}, {}},
			"1",
			false,
			pipeline.Result{Error: pipeline.ErrBadInput},
		},
		{
			"only nil inputs",
			[]pipeline.Result{{}},
			"0",
			false,
			pipeline.Result{Error: pipeline.ErrBadInput},
		},
		{
			"zero inputs",
			[]pipeline.Result{},
			"0",
			false,
			pipeline.Result{Error: pipeline.ErrWrongInputCardinality},
		},
		{
			"fewer errors than threshold",
			[]pipeline.Result{{Error: errors.New("")}, {Value: mustDecimal(t, "2")}, {Value: mustDecimal(t, "3")}, {Value: mustDecimal(t, "4")}},
			"2",
			false,
			pipeline.Result{Value: mustDecimal(t, "2")},
		},
		{
			"exactly threshold of errors",
			[]pipeline.Result{{Error: errors.New("")}, {Error: errors.New("")}, {Value: mustDecimal(t, "3")}, {Value: mustDecimal(t, "4")}},
			"2",
			false,
			pipeline.Result{Value: mustDecimal(t, "3")},
		},
		{
			"more errors than threshold",
			[]pipeline.Result{{Error: errors.New("")}, {Error: errors.New("")}, {Error: errors.New("")}, {Value: mustDecimal(t, "4")}},
			"2",
			false,
			pipeline.Result{Error: pipeline.ErrTooManyErrors},
		},
		{
			"(unspecified AllowedFaults) fewer errors than threshold",
			[]pipeline.Result{{Error: errors.New("")}, {Error: errors.New("")}, {Value: mustDecimal(t, "3")}, {Value: mustDecimal(t, "4")}},
			"",
			false,
			pipeline.Result{Value: mustDecimal(t, "3")},
		},
		{
			"(unspecified AllowedFaults) exactly threshold of errors",
			[]pipeline.Result{{Error: errors.New("")}, {Error: errors.New("")}, {Error: errors.New("")}, {Value: mustDecimal(t, "4")}},
			"",
			false,
			pipeline.Result{Value: mustDecimal(t, "4")},
		},
		{
			"(unspecified AllowedFaults) more errors than threshold",
			[]pipeline.Result{{Error: errors.New("")}, {Error: errors.New("")}, {Error: errors.New("")}},
			"",
			false,
			pipeline.Result{Error: pipeline.ErrTooManyErrors},
		},
		{
			"lax with nil and non-nil inputs",
			[]pipeline.Result{{Value: mustDecimal(t, "1")}, {Value: mustDecimal(t, "2")}, {}},
			"1",
			true,
			pipeline.Result{Value: mustDecimal(t, "1")},
		},
		{
			"lax with more nils than allowed faults",
			[]pipeline.Result{{Value: mustDecimal(t, "1")}, {}, {}, {}},
			"3",
			true,
			pipeline.Result{Value: mustDecimal(t, "1")},
		},
		{
			"lax with nils and errors",
			[]pipeline.Result{{Value: mustDecimal(t, "1")}, {Error: errors.New("1")}, {Error: errors.New("2")}, {}},
			"2",
			true,
			pipeline.Result{Value: mustDecimal(t, "1")},
		},
		{
			"lax with nils and more errors than allowed faults",
			[]pipeline.Result{{Value: mustDecimal(t, "1")}, {Error: errors.New("1")}, {Error: errors.New("2")}, {}},
			"1",
			true,
			pipeline.Result{Error: pipeline.ErrTooManyErrors},
		},
		{
			"lax with numbers and errors and unset allowed faults",
			[]pipeline.Result{{Value: mustDecimal(t, "1")}, {Error: errors.New("1")}, {Error: errors.New("2")}, {}},
			"",
			true,
			pipeline.Result{Value: mustDecimal(t, "1")},
		},
		{
			"lax with only errors and unset allowed faults",
			[]pipeline.Result{{Error: errors.New("1")}, {Error: errors.New("2")}, {}, {}},
			"",
			true,
			pipeline.Result{Error: pipeline.ErrTooManyErrors},
		},
		{
			"lax with only nils",
			[]pipeline.Result{{}, {}, {}, {}},
			"1",
			true,
			pipeline.Result{},
		},
		{
			"lax with only nils and unset allowed faults",
			[]pipeline.Result{{}, {}, {}, {}},
			"",
			true,
			pipeline.Result{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Run("without vars", func(t *testing.T) {
				task := pipeline.MinTask{
					BaseTask:      pipeline.NewBaseTask(0, "task", nil, nil, 0),
					AllowedFaults: test.allowedFaults,
					Lax:           strconv.FormatBool(test.lax),
				}
				output, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), pipeline.NewVarsFrom(nil), test.inputs)
				assert.False(t, runInfo.IsPending)
				assert.False(t, runInfo.IsRetryable)

				switch {
				case test.want.Error != nil:
					require.Equal(t, test.want.Error, errors.Cause(output.Error))
					require.Nil(t, output.Value)
				case test.want.Value != nil:
					require.Equal(t, test.want.Value.(*decimal.Decimal).String(), output.Value.(decimal.Decimal).String())
					require.NoError(t, output.Error)
				default:
					require.Nil(t, output.Value)
					require.NoError(t, output.Error)
				}
			})

			t.Run("with vars", func(t *testing.T) {
				var inputs []any
				for _, input := range test.inputs {
					if input.Error != nil {
						inputs = append(inputs, input.Error)
					} else {
						inputs = append(inputs, input.Value)
					}
				}
				vars := pipeline.NewVarsFrom(map[string]any{
					"foo": map[string]any{"bar": inputs},
				})
				task := pipeline.MinTask{
					BaseTask:      pipeline.NewBaseTask(0, "task", nil, nil, 0),
					Values:        "$(foo.bar)",
					AllowedFaults: test.allowedFaults,
					Lax:           strconv.FormatBool(test.lax),
				}
				output, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), vars, nil)
				assert.False(t, runInfo.IsPending)
				assert.False(t, runInfo.IsRetryable)

				switch {
				case test.want.Error != nil:
					require.Equal(t, test.want.Error, errors.Cause(output.Error))
					require.Nil(t, output.Value)
				case test.want.Value != nil:
					require.Equal(t, test.want.Value.(*decimal.Decimal).String(), output.Value.(decimal.Decimal).String())
					require.NoError(t, output.Error)
				default:
					require.Nil(t, output.Value)
					require.NoError(t, output.Error)
				}
			})

			t.Run("with json vars", func(t *testing.T) {
				var inputs []any
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
					vars = pipeline.NewVarsFrom(map[string]any{"foo": inputs[0]})
				case 2:
					valuesParam = "[ $(foo), $(bar) ]"
					vars = pipeline.NewVarsFrom(map[string]any{"foo": inputs[0], "bar": inputs[1]})
				case 3:
					valuesParam = "[ $(foo), $(bar), $(chain) ]"
					vars = pipeline.NewVarsFrom(map[string]any{"foo": inputs[0], "bar": inputs[1], "chain": inputs[2]})
				case 4:
					valuesParam = "[ $(foo), $(bar), $(chain), $(link) ]"
					vars = pipeline.NewVarsFrom(map[string]any{"foo": inputs[0], "bar": inputs[1], "chain": inputs[2], "link": inputs[3]})
				}

				task := pipeline.MinTask{
					BaseTask:      pipeline.NewBaseTask(0, "task", nil, nil, 0),
					Values:        valuesParam,
					AllowedFaults: test.allowedFaults,
					Lax:           strconv.FormatBool(test.lax),
				}
				output, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), vars, nil)
				assert.False(t, runInfo.IsPending)
				assert.False(t, runInfo.IsRetryable)

				switch {
				case test.want.Error != nil:
					require.Equal(t, test.want.Error, errors.Cause(output.Error))
					require.Nil(t, output.Value)
				case test.want.Value != nil:
					require.Equal(t, test.want.Value.(*decimal.Decimal).String(), output.Value.(decimal.Decimal).String())
					require.NoError(t, output.Error)
				default:
					require.Nil(t, output.Value)
					require.NoError(t, output.Error)
				}
			})
		})
	}
}
