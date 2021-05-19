package pipeline_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func TestMedian(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		inputs        []pipeline.Result
		allowedFaults uint64
		want          pipeline.Result
	}{
		{
			"odd number of inputs",
			[]pipeline.Result{{Value: mustDecimal(t, "1")}, {Value: mustDecimal(t, "2")}, {Value: mustDecimal(t, "3")}},
			1,
			pipeline.Result{Value: mustDecimal(t, "2")},
		},
		{
			"even number of inputs",
			[]pipeline.Result{{Value: mustDecimal(t, "1")}, {Value: mustDecimal(t, "2")}, {Value: mustDecimal(t, "3")}, {Value: mustDecimal(t, "4")}},
			2,
			pipeline.Result{Value: mustDecimal(t, "2.5")},
		},
		{
			"one input",
			[]pipeline.Result{{Value: mustDecimal(t, "1")}},
			0,
			pipeline.Result{Value: mustDecimal(t, "1")},
		},
		{
			"zero inputs",
			[]pipeline.Result{},
			0,
			pipeline.Result{Error: pipeline.ErrWrongInputCardinality},
		},
		{
			"fewer errors than threshold",
			[]pipeline.Result{{Error: errors.New("")}, {Value: mustDecimal(t, "2")}, {Value: mustDecimal(t, "3")}, {Value: mustDecimal(t, "4")}},
			2,
			pipeline.Result{Value: mustDecimal(t, "3")},
		},
		{
			"exactly threshold of errors",
			[]pipeline.Result{{Error: errors.New("")}, {Error: errors.New("")}, {Value: mustDecimal(t, "3")}, {Value: mustDecimal(t, "4")}},
			2,
			pipeline.Result{Value: mustDecimal(t, "3.5")},
		},
		{
			"more errors than threshold",
			[]pipeline.Result{{Error: errors.New("")}, {Error: errors.New("")}, {Error: errors.New("")}, {Value: mustDecimal(t, "4")}},
			2,
			pipeline.Result{Error: pipeline.ErrBadInput},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			task := pipeline.MedianTask{AllowedFaults: fmt.Sprintf("%v", test.allowedFaults)}
			output := task.Run(context.Background(), nil, pipeline.JSONSerializable{}, test.inputs)
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

func TestMedian_Defaults(t *testing.T) {
	t.Parallel()

	var taskDAG pipeline.TaskDAG
	err := taskDAG.UnmarshalText([]byte(pipeline.DotStr))
	require.NoError(t, err)

	tasks, err := taskDAG.TasksInDependencyOrder()
	require.NoError(t, err)

	for _, task := range tasks {
		if asMedian, isMedian := task.(*pipeline.MedianTask); isMedian {
			require.Equal(t, uint64(1), asMedian.AllowedFaults)
			break
		}
	}
}

func TestMedian_Faults_Unmarshal(t *testing.T) {
	t.Parallel()

	var taskDAG pipeline.TaskDAG
	err := taskDAG.UnmarshalText([]byte(`
	// data source 1
	ds1          [type=bridge name=voter_turnout];
	ds1_parse    [type=jsonparse path="one,two"];
	ds1_multiply [type=multiply times=1.23];

	// data source 2
	ds2          [type=http method=GET url="https://chain.link/voter_turnout/USA-2020" requestData="{\"hi\": \"hello\"}"];
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
			assert.Equal(t, uint64(10), task.(*pipeline.MedianTask).AllowedFaults)
		}
	}
}
