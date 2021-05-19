package pipeline

import (
	"context"
	"sort"

	"github.com/shopspring/decimal"
	"go.uber.org/multierr"
)

type MedianTask struct {
	BaseTask      `mapstructure:",squash"`
	Values        string `json:"values"`
	AllowedFaults string `json:"allowedFaults"`
}

var _ Task = (*MedianTask)(nil)

func (t *MedianTask) Type() TaskType {
	return TaskTypeMedian
}

func (t *MedianTask) SetDefaults(inputValues map[string]string, g TaskDAG, self TaskDAGNode) error {
	return nil
}

func (t *MedianTask) Run(_ context.Context, vars Vars, _ JSONSerializable, inputs []Result) (result Result) {
	var (
		allowedFaults Uint64Param
		values        DecimalSliceParam
	)
	err := multierr.Combine(
		vars.ResolveValue(&allowedFaults, From(NonemptyString(t.AllowedFaults), len(inputs)-1)),
		vars.ResolveValue(&values, From(VariableExpr(t.Values), Inputs(inputs, 1, -1, int(allowedFaults)))),
	)
	if err != nil {
		return Result{Error: err}
	}

	sort.Slice(values, func(i, j int) bool {
		return values[i].LessThan(values[j])
	})
	k := len(values) / 2
	if len(values)%2 == 1 {
		return Result{Value: values[k]}
	}
	median := values[k].Add(values[k-1]).Div(decimal.NewFromInt(2))

	err = vars.Set(t.DotID(), median)
	if err != nil {
		return Result{Error: err}
	}
	return Result{Value: median}
}
