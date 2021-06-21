package pipeline

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

type DivideTask struct {
	BaseTask  `mapstructure:",squash"`
	Input     string `json:"input"`
	Divisor   string `json:"divisor"`
	Precision string `json:"precision"`
}

var _ Task = (*DivideTask)(nil)

func (t *DivideTask) Type() TaskType {
	return TaskTypeDivide
}

func (t *DivideTask) Run(_ context.Context, vars Vars, inputs []Result) (result Result) {
	_, err := CheckInputs(inputs, -1, -1, 0)
	if err != nil {
		return Result{Error: errors.Wrap(err, "task inputs")}
	}

	var (
		a              DecimalParam
		b              DecimalParam
		maybePrecision MaybeInt32Param
	)
	err = multierr.Combine(
		errors.Wrap(ResolveParam(&a, From(VarExpr(t.Input, vars), Input(inputs, 0))), "input"),
		errors.Wrap(ResolveParam(&b, From(VarExpr(t.Divisor, vars), NonemptyString(t.Divisor))), "divisor"),
		errors.Wrap(ResolveParam(&maybePrecision, From(VarExpr(t.Precision, vars), t.Precision)), "precision"),
	)
	if err != nil {
		return Result{Error: err}
	}

	value := a.Decimal().Div(b.Decimal())

	if precision, isSet := maybePrecision.Int32(); isSet {
		value = value.Round(precision)
	}

	return Result{Value: value}
}
