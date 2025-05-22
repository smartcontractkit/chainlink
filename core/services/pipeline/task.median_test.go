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

func TestMedianTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		inputs        []pipeline.Result
		allowedFaults string
		lax           string
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
			"(unspecified AllowedFaults) zero inputs",
			[]pipeline.Result{},
			"",
			"",
			pipeline.Result{Error: pipeline.ErrWrongInputCardinality},
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
			"(unspecified Lax) error on parsing nil inputs",
			[]pipeline.Result{{}, {Value: mustDecimal(t, "2")}, {Value: mustDecimal(t, "3")}},
			"",
			"",
			pipeline.Result{Error: pipeline.ErrBadInput},
		},
		{
			"nil inputs with Lax enabled",
			[]pipeline.Result{{}, {Value: errors.New("")}, {Value: mustDecimal(t, "2")}, {Value: mustDecimal(t, "3")}},
			"",
			"true",
			pipeline.Result{Value: mustDecimal(t, "2.5")},
		},
		{
			"zero inputs with Lax enabled",
			[]pipeline.Result{},
			"",
			"true",
			pipeline.Result{},
		},
		{
			"zero non-nil inputs with Lax enabled",
			[]pipeline.Result{{}, {}, {}},
			"",
			"true",
			pipeline.Result{},
		},
		{
			"nil inputs and exact threshold of errors with Lax enabled",
			[]pipeline.Result{{}, {}, {Value: errors.New("")}, {Value: errors.New("")}},
			"2",
			"true",
			pipeline.Result{},
		},
		{
			"nil inputs and more errors than threshold with Lax enabled",
			[]pipeline.Result{{}, {}, {Value: errors.New("")}, {Value: errors.New("")}},
			"1",
			"true",
			pipeline.Result{Error: pipeline.ErrTooManyErrors},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Run("without vars", func(t *testing.T) {
				task := pipeline.MedianTask{
					BaseTask:      pipeline.NewBaseTask(0, "task", nil, nil, 0),
					AllowedFaults: test.allowedFaults,
					Lax:           test.lax,
				}
				output, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), pipeline.NewVarsFrom(nil), test.inputs)
				assert.False(t, runInfo.IsPending)
				assert.False(t, runInfo.IsRetryable)
				if output.Error != nil {
					require.Equal(t, test.want.Error, errors.Cause(output.Error))
					require.Nil(t, output.Value)
				} else {
					if test.want.Value == nil {
						require.Nil(t, output.Value)
					} else {
						require.Equal(t, test.want.Value.(*decimal.Decimal).String(), output.Value.(decimal.Decimal).String())
					}
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
				task := pipeline.MedianTask{
					BaseTask:      pipeline.NewBaseTask(0, "task", nil, nil, 0),
					Values:        "$(foo.bar)",
					AllowedFaults: test.allowedFaults,
					Lax:           test.lax,
				}
				output, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), vars, nil)
				assert.False(t, runInfo.IsPending)
				assert.False(t, runInfo.IsRetryable)
				if output.Error != nil {
					require.Equal(t, test.want.Error, errors.Cause(output.Error))
					require.Nil(t, output.Value)
				} else {
					if test.want.Value == nil {
						require.Nil(t, output.Value)
					} else {
						require.Equal(t, test.want.Value.(*decimal.Decimal).String(), output.Value.(decimal.Decimal).String())
					}
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

				task := pipeline.MedianTask{
					BaseTask:      pipeline.NewBaseTask(0, "task", nil, nil, 0),
					Values:        valuesParam,
					AllowedFaults: test.allowedFaults,
					Lax:           test.lax,
				}
				output, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), vars, nil)
				assert.False(t, runInfo.IsPending)
				assert.False(t, runInfo.IsRetryable)
				if output.Error != nil {
					require.Equal(t, test.want.Error, errors.Cause(output.Error))
					require.Nil(t, output.Value)
				} else {
					if test.want.Value == nil {
						require.Nil(t, output.Value)
					} else {
						require.Equal(t, test.want.Value.(*decimal.Decimal).String(), output.Value.(decimal.Decimal).String())
					}
					require.NoError(t, output.Error)
				}
			})
		})
	}
}

func TestMedianTask_AllowedFaultsAndLax_Unmarshal(t *testing.T) {
	t.Parallel()

	p, err := pipeline.Parse(`
	// data source 1
	ds1          [type=bridge name=voter_turnout];
	ds1_parse    [type=jsonparse path="one,two"];
	ds1_multiply [type=multiply times=1.23];

	// data source 2
	ds2          [type=http method=GET url="https://chain.link/voter_turnout/USA-2020" requestData=<{"hi": "hello"}>];
	ds2_parse    [type=jsonparse path="three,four"];
	ds2_multiply [type=multiply times=4.56];

	ds1 -> ds1_parse -> ds1_multiply -> answer1;
	ds2 -> ds2_parse -> ds2_multiply -> answer1;

	answer1 [type=median                      index=0 allowedFaults=10 lax=true];
	answer2 [type=bridge name=election_winner index=1];
`)
	require.NoError(t, err)
	for _, task := range p.Tasks {
		if task.Type() == pipeline.TaskTypeMedian {
			require.Equal(t, "10", task.(*pipeline.MedianTask).AllowedFaults)
			require.Equal(t, "true", task.(*pipeline.MedianTask).Lax)
		}
	}
}
