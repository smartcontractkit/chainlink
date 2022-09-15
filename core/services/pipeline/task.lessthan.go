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
type LessThan struct {
	BaseTask `mapstructure:",squash"`
	Input    string `json:"input"`
	Limit    string `json:"limit"`
}

var (
	_ Task = (*LessThan)(nil)
)

func (t *LessThan) Type() TaskType {
	return TaskTypeLessThan
}

func (t *LessThan) Run(_ context.Context, _ logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	_, err := CheckInputs(inputs, 0, 1, 0)
	if err != nil {
		return Result{Error: errors.Wrap(err, "task inputs")}, runInfo
	}

	var (
		a DecimalParam
		b DecimalParam
	)

	err = multierr.Combine(
		errors.Wrap(ResolveParam(&a, From(VarExpr(t.Input, vars), NonemptyString(t.Input), Input(inputs, 0))), "input"),
		errors.Wrap(ResolveParam(&b, From(VarExpr(t.Limit, vars), NonemptyString(t.Limit))), "limit"),
	)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	value := a.Decimal().LessThan(b.Decimal())
	return Result{Value: value}, runInfo
}
