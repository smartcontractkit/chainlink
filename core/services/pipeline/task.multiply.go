package pipeline

import (
	"context"
	"math"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// Return types:
//
//	*decimal.Decimal
type MultiplyTask struct {
	BaseTask `mapstructure:",squash"`
	Input    string `json:"input"`
	Times    string `json:"times"`
}

var (
	_                  Task = (*MultiplyTask)(nil)
	ErrMultiplyOverlow      = errors.New("multiply overflow")
)

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
		errors.Wrap(ResolveParam(&a, From(VarExpr(t.Input, vars), NonemptyString(t.Input), Input(inputs, 0))), "input"),
		errors.Wrap(ResolveParam(&b, From(VarExpr(t.Times, vars), NonemptyString(t.Times))), "times"),
	)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	newExp := int64(a.Decimal().Exponent()) + int64(b.Decimal().Exponent())
	if newExp > math.MaxInt32 || newExp < math.MinInt32 {
		return Result{Error: ErrMultiplyOverlow}, runInfo
	}

	value := a.Decimal().Mul(b.Decimal())
	return Result{Value: value}, runInfo
}
