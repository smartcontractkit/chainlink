package pipeline_test

import (
	"context"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func TestMedian(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		inputs        []pipeline.Result
		allowedFaults string
		want          pipeline.Result
	}{
		{
			"odd number of inputs",
			[]pipeline.Result{{Value: mustDecimal(t, "1")}, {Value: mustDecimal(t, "2")}, {Value: mustDecimal(t, "3")}},
			"1",
			pipeline.Result{Value: mustDecimal(t, "2")},
		},
		{
			"even number of inputs",
			[]pipeline.Result{{Value: mustDecimal(t, "1")}, {Value: mustDecimal(t, "2")}, {Value: mustDecimal(t, "3")}, {Value: mustDecimal(t, "4")}},
			"2",
			pipeline.Result{Value: mustDecimal(t, "2.5")},
		},
		{
			"one input",
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
			pipeline.Result{Value: mustDecimal(t, "3")},
		},
		{
			"exactly threshold of errors",
			[]pipeline.Result{{Error: errors.New("")}, {Error: errors.New("")}, {Value: mustDecimal(t, "3")}, {Value: mustDecimal(t, "4")}},
			"2",
			pipeline.Result{Value: mustDecimal(t, "3.5")},
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
			pipeline.Result{Value: mustDecimal(t, "3.5")},
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
		test := test
		t.Run(test.name, func(t *testing.T) {
			vars := pipeline.Vars{}
			task := pipeline.MedianTask{
				BaseTask:      pipeline.NewBaseTask("task", nil, 0, 0),
				AllowedFaults: test.allowedFaults,
			}
			output := task.Run(context.Background(), vars, pipeline.JSONSerializable{}, test.inputs)
			if output.Error != nil {
				require.Equal(t, test.want.Error, errors.Cause(output.Error))
				require.Nil(t, output.Value)
			} else {
				require.Equal(t, test.want.Value.(*decimal.Decimal).String(), output.Value.(decimal.Decimal).String())
				require.NoError(t, output.Error)
			}
		})
	}

	for _, test := range tests {
		test := test
		t.Run(test.name+" (with pipeline.Vars)", func(t *testing.T) {
			t.Parallel()

			var inputs []interface{}
			for _, input := range test.inputs {
				if input.Error != nil {
					inputs = append(inputs, input.Error)
				} else {
					inputs = append(inputs, input.Value)
				}
			}
			vars := pipeline.Vars{
				"foo": map[string]interface{}{"bar": inputs},
			}
			task := pipeline.MedianTask{
				BaseTask:      pipeline.NewBaseTask("task", nil, 0, 0),
				Values:        "$(foo.bar)",
				AllowedFaults: test.allowedFaults,
			}
			output := task.Run(context.Background(), vars, pipeline.JSONSerializable{}, nil)
			if output.Error != nil {
				require.Equal(t, test.want.Error, errors.Cause(output.Error))
				require.Nil(t, output.Value)
			} else {
				require.Equal(t, test.want.Value.(*decimal.Decimal).String(), output.Value.(decimal.Decimal).String())
				require.NoError(t, output.Error)
			}
		})
	}
}

func TestMedian_MedianizeResultsOfJSONParse(t *testing.T) {
	t.Parallel()

	vars := pipeline.Vars{
		"input": map[string]interface{}{"response": `{"data": [1, 3, 5, 7]}`},
	}

	jsonParse := pipeline.JSONParseTask{
		BaseTask: pipeline.NewBaseTask("json", nil, 0, 0),
		Path:     "data",
		Data:     "$(input.response)",
	}
	result := jsonParse.Run(context.Background(), vars, pipeline.JSONSerializable{}, []pipeline.Result{})
	require.NoError(t, result.Error)
	require.NotNil(t, result.Value)

	median := pipeline.MedianTask{
		BaseTask: pipeline.NewBaseTask("median", nil, 0, 0),
		Values:   "$(json)",
	}
	result = median.Run(context.Background(), vars, pipeline.JSONSerializable{}, []pipeline.Result{})
	require.NoError(t, result.Error)
	require.True(t, mustDecimal(t, "4").Equal(result.Value.(decimal.Decimal)))
}

func TestMedian_AllowedFaults_Unmarshal(t *testing.T) {
	t.Parallel()

	var taskDAG pipeline.TaskDAG
	err := taskDAG.UnmarshalText([]byte(`
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

	answer1 [type=median                      index=0 allowedFaults=10];
	answer2 [type=bridge name=election_winner index=1];
`))
	require.NoError(t, err)
	ts, err := taskDAG.TasksInDependencyOrder()
	require.NoError(t, err)
	for _, task := range ts {
		if task.Type() == pipeline.TaskTypeMedian {
			assert.Equal(t, "10", task.(*pipeline.MedianTask).AllowedFaults)
		}
	}
}
