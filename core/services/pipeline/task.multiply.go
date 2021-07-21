package pipeline

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

//
// Return types:
//    *decimal.Decimal
//
type MultiplyTask struct {
	BaseTask `mapstructure:",squash"`
	Input    string `json:"input"`
	Times    string `json:"times"`
}

var _ Task = (*MultiplyTask)(nil)

func (t *MultiplyTask) Type() TaskType {
	return TaskTypeMultiply
}

func (t *MultiplyTask) Run(_ context.Context, vars Vars, inputs []Result) (result Result) {
	_, err := CheckInputs(inputs, 0, 1, 0)
	if err != nil {
		return Result{Error: errors.Wrap(err, "task inputs")}
	}

	var (
		a DecimalParam
		b DecimalParam
	)
	err = multierr.Combine(
		errors.Wrap(ResolveParam(&a, From(VarExpr(t.Input, vars), Input(inputs, 0))), "input"),
		errors.Wrap(ResolveParam(&b, From(VarExpr(t.Times, vars), NonemptyString(t.Times))), "times"),
	)
	if err != nil {
		return Result{Error: err}
	}

	value := a.Decimal().Mul(b.Decimal())
	return Result{Value: value}
}
