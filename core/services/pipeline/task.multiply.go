package pipeline

import (
	"context"

	"go.uber.org/multierr"
)

type MultiplyTask struct {
	BaseTask `mapstructure:",squash"`
	Times    string `json:"times"`
}

var _ Task = (*MultiplyTask)(nil)

func (t *MultiplyTask) Type() TaskType {
	return TaskTypeMultiply
}

func (t *MultiplyTask) SetDefaults(inputValues map[string]string, g TaskDAG, self TaskDAGNode) error {
	return nil
}

func (t *MultiplyTask) Run(_ context.Context, vars Vars, _ JSONSerializable, inputs []Result) (result Result) {
	_, err := CheckInputs(inputs, 0, 1, 0)
	if err != nil {
		return Result{Error: err}
	}

	var (
		a DecimalParam
		b DecimalParam
	)
	err = multierr.Combine(
		vars.ResolveValue(&a, From(Input(inputs, 0))),
		vars.ResolveValue(&b, From(VariableExpr(t.Times), NonemptyString(t.Times))),
	)
	if err != nil {
		return Result{Error: err}
	}

	value := a.Decimal().Mul(b.Decimal())

	err = vars.Set(t.DotID(), value)
	if err != nil {
		return Result{Error: err}
	}
	return Result{Value: value}
}
