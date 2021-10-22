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

func (t *DivideTask) Run(_ context.Context, _ logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	_, err := CheckInputs(inputs, -1, -1, 0)
	if err != nil {
		return Result{Error: errors.Wrap(err, "task inputs")}, runInfo
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
		return Result{Error: err}, runInfo
	}

	if precision, isSet := maybePrecision.Int32(); isSet {
		return Result{Value: a.Decimal().DivRound(b.Decimal(), precision)}, runInfo
	}
	// Note that decimal library defaults to rounding to 16 precision
	// https://github.com/shopspring/decimal/blob/2568a29459476f824f35433dfbef158d6ad8618c/decimal.go#L44
	return Result{Value: a.Decimal().Div(b.Decimal())}, runInfo
}
