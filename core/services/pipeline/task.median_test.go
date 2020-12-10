package pipeline_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func TestMedian(t *testing.T) {
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
			pipeline.Result{Error: pipeline.ErrBadInput},
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
			task := pipeline.MedianTask{AllowedFaults: test.allowedFaults}
			output := task.Run(context.Background(), pipeline.TaskRun{}, test.inputs)
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
	var taskDAG pipeline.TaskDAG
	err := taskDAG.UnmarshalText([]byte(dotStr))
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
