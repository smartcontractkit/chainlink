package pipeline

import (
	"context"
	stderrors "errors"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// Return types:
//
//	*decimal.Decimal
type MaxTask struct {
	BaseTask      `mapstructure:",squash"`
	Values        string `json:"values"`
	AllowedFaults string `json:"allowedFaults"`
	// Lax when disabled (default) will return an error if there are no input values or if the input includes nil values.
	// Lax when enabled will return nil with no error if there are no valid input values. If the input includes nil values, they will be excluded from the calculation and do not count as a fault.
	Lax string
}

var _ Task = (*MaxTask)(nil)

func (t *MaxTask) Type() TaskType {
	return TaskTypeMax
}

func (t *MaxTask) Run(_ context.Context, _ logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	var (
		maybeAllowedFaults MaybeUint64Param
		valuesAndErrs      SliceParam
		decimalValues      DecimalSliceParam
		allowedFaults      int
		lax                BoolParam
	)
	err := stderrors.Join(
		errors.Wrap(ResolveParam(&maybeAllowedFaults, From(t.AllowedFaults)), "allowedFaults"),
		errors.Wrap(ResolveParam(&valuesAndErrs, From(VarExpr(t.Values, vars), JSONWithVarExprs(t.Values, vars, true), Inputs(inputs))), "values"),
		errors.Wrap(ResolveParam(&lax, From(NonemptyString(t.Lax), false)), "lax"),
	)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	// if lax is enabled, filter out nil values
	// nil values are not included in the fault calculations
	if lax {
		valuesAndErrs, _ = valuesAndErrs.FilterNils()
	}

	if allowed, isSet := maybeAllowedFaults.Uint64(); isSet {
		allowedFaults = int(allowed) //nolint:gosec // G115: it will not exceed int64
	} else {
		allowedFaults = max(len(valuesAndErrs)-1, 0)
	}

	values, faults := valuesAndErrs.FilterErrors()
	if faults > allowedFaults {
		return Result{Error: errors.Wrapf(ErrTooManyErrors, "Number of faulty inputs %v to max task > number allowed faults %v", faults, allowedFaults)}, runInfo
	}
	if len(values) == 0 {
		if lax {
			return Result{}, runInfo // if lax is enabled, return nil result with no error
		}
		return Result{Error: errors.Wrap(ErrWrongInputCardinality, "no values to maxize")}, runInfo
	}

	err = decimalValues.UnmarshalPipelineParam(values)
	if err != nil {
		return Result{Error: errors.Wrapf(ErrBadInput, "values: %v", err)}, runInfo
	}

	maxVal := decimal.Max(decimalValues[0], decimalValues[1:]...)

	return Result{Value: maxVal}, runInfo
}
