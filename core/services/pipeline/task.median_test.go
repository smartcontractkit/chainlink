package pipeline_test

import (
	"github.com/pkg/errors"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func TestMedian(t *testing.T) {
	tests := []struct {
		name   string
		inputs []pipeline.Result
		want   pipeline.Result
	}{
		{
			"odd number of inputs",
			[]pipeline.Result{{Value: mustDecimal(t, "1")}, {Value: mustDecimal(t, "2")}, {Value: mustDecimal(t, "3")}},
			pipeline.Result{Value: mustDecimal(t, "2")},
		},
		{
			"even number of inputs",
			[]pipeline.Result{{Value: mustDecimal(t, "1")}, {Value: mustDecimal(t, "2")}, {Value: mustDecimal(t, "3")}, {Value: mustDecimal(t, "4")}},
			pipeline.Result{Value: mustDecimal(t, "2.5")},
		},
		{
			"one input",
			[]pipeline.Result{{Value: mustDecimal(t, "1")}},
			pipeline.Result{Value: mustDecimal(t, "1")},
		},
		{
			"zero inputs",
			[]pipeline.Result{},
			pipeline.Result{Error: pipeline.ErrWrongInputCardinality},
		},
		{
			"< 50% errors",
			[]pipeline.Result{{Error: errors.New("")}, {Value: mustDecimal(t, "2")}, {Value: mustDecimal(t, "3")}, {Value: mustDecimal(t, "4")}},
			pipeline.Result{Value: mustDecimal(t, "3")},
		},
		{
			"50% errors",
			[]pipeline.Result{{Error: errors.New("")}, {Error: errors.New("")}, {Value: mustDecimal(t, "3")}, {Value: mustDecimal(t, "4")}},
			pipeline.Result{Error: pipeline.ErrBadInput},
		},
		{
			"> 50% errors",
			[]pipeline.Result{{Error: errors.New("")}, {Error: errors.New("")}, {Error: errors.New("")}, {Value: mustDecimal(t, "4")}},
			pipeline.Result{Error: pipeline.ErrBadInput},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			task := pipeline.MedianTask{}
			output := task.Run(pipeline.TaskRun{}, test.inputs)
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
