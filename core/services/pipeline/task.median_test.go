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
		{
			// A bridge returning HTTP 200 with null numeric values produces a nil result.
			// The nil passes through FilterErrors (it's not an error), then hits
			// DecimalSliceParam.UnmarshalPipelineParam which fails hard on nil.
			// The entire median task errors out even though 3 valid values were available.
			"single nil from bridge returning null kills entire median",
			[]pipeline.Result{
				{},                           // bridge returned HTTP 200 with null value
				{Value: mustDecimal(t, "1")}, // 3 valid bridges — enough for a median
				{Value: mustDecimal(t, "2")},
				{Value: mustDecimal(t, "3")},
			},
			"3",
			"",
			pipeline.Result{Error: pipeline.ErrBadInput}, // nil poisons the entire batch
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
				case 3:
					valuesParam = "[ $(foo), $(bar), $(chain) ]"
					vars = pipeline.NewVarsFrom(map[string]any{"foo": inputs[0], "bar": inputs[1], "chain": inputs[2]})
				case 4:
					valuesParam = "[ $(foo), $(bar), $(chain), $(link) ]"
					vars = pipeline.NewVarsFrom(map[string]any{"foo": inputs[0], "bar": inputs[1], "chain": inputs[2], "link": inputs[3]})
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

func TestMedianTask_CountNilsAsFaults(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		inputs            []pipeline.Result
		allowedFaults     string
		countNilsAsFaults string
		want              pipeline.Result
	}{
		{
			name: "1 nil + 7 errors exceeds allowedFaults",
			inputs: []pipeline.Result{
				{}, // nil
				{Value: errors.New("e1")}, {Value: errors.New("e2")}, {Value: errors.New("e3")},
				{Value: errors.New("e4")}, {Value: errors.New("e5")}, {Value: errors.New("e6")},
				{Value: errors.New("e7")},
			},
			allowedFaults:     "7",
			countNilsAsFaults: "true",
			want:              pipeline.Result{Error: pipeline.ErrTooManyErrors},
		},
		{
			// With countNilsAsFaults enabled, the nil is counted as a fault (1 fault total)
			// which is within allowedFaults=3, then filtered out before decimal parsing.
			// Median proceeds on the 3 valid values instead of crashing with ErrBadInput.
			name: "countNilsAsFaults prevents nil from killing median",
			inputs: []pipeline.Result{
				{},                           // bridge returned HTTP 200 with null value
				{Value: mustDecimal(t, "1")}, // 3 valid bridges
				{Value: mustDecimal(t, "2")},
				{Value: mustDecimal(t, "3")},
			},
			allowedFaults:     "3",
			countNilsAsFaults: "true",
			want:              pipeline.Result{Value: mustDecimal(t, "2")},
		},
		{
			// Nils and errors are both counted as faults. Together they exceed allowedFaults,
			// so the task correctly fails even though valid values exist.
			name: "combined nils and errors exceed allowedFaults",
			inputs: []pipeline.Result{
				{},                           // nil (1 fault)
				{},                           // nil (2 faults)
				{Value: errors.New("err1")},  // error (3 faults)
				{Value: mustDecimal(t, "5")}, // valid
			},
			allowedFaults:     "2",
			countNilsAsFaults: "true",
			want:              pipeline.Result{Error: pipeline.ErrTooManyErrors}, // 3 faults > 2
		},
		{
			name: "nils and errors within threshold returns valid median",
			inputs: []pipeline.Result{
				{},                           // nil (counted as fault)
				{Value: errors.New("err")},   // error (counted as fault)
				{Value: mustDecimal(t, "2")}, // valid
				{Value: mustDecimal(t, "4")}, // valid
			},
			allowedFaults:     "2",
			countNilsAsFaults: "true",
			want:              pipeline.Result{Value: mustDecimal(t, "3")},
		},
		{
			name:              "all nils exceed threshold",
			inputs:            []pipeline.Result{{}, {}, {}},
			allowedFaults:     "2",
			countNilsAsFaults: "true",
			want:              pipeline.Result{Error: pipeline.ErrTooManyErrors},
		},
		{
			name:              "all nils within threshold returns empty result",
			inputs:            []pipeline.Result{{}, {}},
			allowedFaults:     "2",
			countNilsAsFaults: "true",
			want:              pipeline.Result{},
		},
		{
			name: "no nils no errors returns normal median",
			inputs: []pipeline.Result{
				{Value: mustDecimal(t, "1")},
				{Value: mustDecimal(t, "2")},
				{Value: mustDecimal(t, "3")},
			},
			allowedFaults:     "1",
			countNilsAsFaults: "true",
			want:              pipeline.Result{Value: mustDecimal(t, "2")},
		},
		{
			name: "only errors no nils behaves like default",
			inputs: []pipeline.Result{
				{Value: errors.New("e1")}, {Value: errors.New("e2")}, {Value: errors.New("e3")},
				{Value: mustDecimal(t, "4")},
			},
			allowedFaults:     "2",
			countNilsAsFaults: "true",
			want:              pipeline.Result{Error: pipeline.ErrTooManyErrors},
		},
		{
			name: "single nil with one valid value within threshold",
			inputs: []pipeline.Result{
				{},
				{Value: mustDecimal(t, "42")},
			},
			allowedFaults:     "1",
			countNilsAsFaults: "true",
			want:              pipeline.Result{Value: mustDecimal(t, "42")},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			task := pipeline.MedianTask{
				BaseTask:          pipeline.NewBaseTask(0, "task", nil, nil, 0),
				AllowedFaults:     test.allowedFaults,
				CountNilsAsFaults: test.countNilsAsFaults,
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
	}

	t.Run("mutual exclusion: lax and countNilsAsFaults cannot both be enabled", func(t *testing.T) {
		task := pipeline.MedianTask{
			BaseTask:          pipeline.NewBaseTask(0, "task", nil, nil, 0),
			AllowedFaults:     "1",
			Lax:               "true",
			CountNilsAsFaults: "true",
		}
		output, _ := task.Run(testutils.Context(t), logger.TestLogger(t), pipeline.NewVarsFrom(nil), []pipeline.Result{{Value: mustDecimal(t, "1")}})
		require.Error(t, output.Error)
		require.Contains(t, output.Error.Error(), "lax and countNilsAsFaults cannot both be enabled")
	})
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

func TestMedianTask_CountNilsAsFaults_Unmarshal(t *testing.T) {
	t.Parallel()

	p, err := pipeline.Parse(`
	ds1          [type=bridge name=voter_turnout];
	ds1_parse    [type=jsonparse path="one,two"];

	ds2          [type=bridge name=voter_turnout2];
	ds2_parse    [type=jsonparse path="three,four"];

	ds1 -> ds1_parse -> answer1;
	ds2 -> ds2_parse -> answer1;

	answer1 [type=median index=0 allowedFaults=5 countNilsAsFaults=true];
`)
	require.NoError(t, err)
	for _, task := range p.Tasks {
		if task.Type() == pipeline.TaskTypeMedian {
			require.Equal(t, "5", task.(*pipeline.MedianTask).AllowedFaults)
			require.Equal(t, "true", task.(*pipeline.MedianTask).CountNilsAsFaults)
			require.Empty(t, task.(*pipeline.MedianTask).Lax)
		}
	}
}
