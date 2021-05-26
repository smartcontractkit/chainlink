package pipeline

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

type MultiplyTask struct {
	BaseTask `mapstructure:",squash"`
	Input    string `json:"input" pipeline:"@expand_vars"`
	Times    string `json:"times" pipeline:"@expand_vars"`
}

var _ Task = (*MultiplyTask)(nil)

func (t *MultiplyTask) Type() TaskType {
	return TaskTypeMultiply
}

func (t *MultiplyTask) Run(_ context.Context, _ JSONSerializable, inputs []Result) (result Result) {
	_, err := CheckInputs(inputs, 0, 1, 0)
	if err != nil {
		return Result{Error: err}
	}

	var (
		a DecimalParam
		b DecimalParam
	)
	err = multierr.Combine(
		errors.Wrap(ResolveParam(&a, From(t.Input, Input(inputs, 0))), "input"),
		errors.Wrap(ResolveParam(&b, From(NonemptyString(t.Times))), "times"),
	)
	if err != nil {
		return Result{Error: err}
	}

	value := a.Decimal().Mul(b.Decimal())
	return Result{Value: value}
}
