package pipeline

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/logger"
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

func (t *MultiplyTask) Run(_ context.Context, _ logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	_, err := CheckInputs(inputs, 0, 1, 0)
	if err != nil {
		return Result{Error: errors.Wrap(err, "task inputs")}, runInfo
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
		return Result{Error: err}, runInfo
	}

	value := a.Decimal().Mul(b.Decimal())
	return Result{Value: value}, runInfo
}
