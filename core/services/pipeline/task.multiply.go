package pipeline

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

type MultiplyTask struct {
	BaseTask `mapstructure:",squash"`
	Input    string `json:"input"`
	Times    string `json:"times"`
}

var _ Task = (*MultiplyTask)(nil)

func (t *MultiplyTask) Type() TaskType {
	return TaskTypeMultiply
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
		errors.Wrap(vars.ResolveValue(&a, From(VariableExpr(t.Input), Input(inputs, 0))), "input"),
		errors.Wrap(vars.ResolveValue(&b, From(VariableExpr(t.Times), NonemptyString(t.Times))), "times"),
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
