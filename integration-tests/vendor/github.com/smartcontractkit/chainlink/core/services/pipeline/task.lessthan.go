package pipeline

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/logger"
)

// Return types:
//
//	bool
type LessThanTask struct {
	BaseTask `mapstructure:",squash"`
	Left     string `json:"input"`
	Right    string `json:"limit"`
}

var (
	_ Task = (*LessThanTask)(nil)
)

func (t *LessThanTask) Type() TaskType {
	return TaskTypeLessThan
}

func (t *LessThanTask) Run(_ context.Context, _ logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	_, err := CheckInputs(inputs, 0, 1, 0)
	if err != nil {
		return Result{Error: errors.Wrap(err, "task inputs")}, runInfo
	}

	var (
		a DecimalParam
		b DecimalParam
	)

	err = multierr.Combine(
		errors.Wrap(ResolveParam(&a, From(VarExpr(t.Left, vars), NonemptyString(t.Left), Input(inputs, 0))), "left"),
		errors.Wrap(ResolveParam(&b, From(VarExpr(t.Right, vars), NonemptyString(t.Right))), "right"),
	)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	value := a.Decimal().LessThan(b.Decimal())
	return Result{Value: value}, runInfo
}
