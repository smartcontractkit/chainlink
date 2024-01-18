package pipeline

import (
	"context"
	"errors"

	pkgerrors "github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// Return types:
//
//	*decimal.Decimal
type MeanTask struct {
	BaseTask      `mapstructure:",squash"`
	Values        string `json:"values"`
	AllowedFaults string `json:"allowedFaults"`
	Precision     string `json:"precision"`
}

var _ Task = (*MeanTask)(nil)

func (t *MeanTask) Type() TaskType {
	return TaskTypeMean
}

func (t *MeanTask) Run(ctx context.Context, lggr logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	var (
		maybeAllowedFaults MaybeUint64Param
		maybePrecision     MaybeInt32Param
		valuesAndErrs      SliceParam
		decimalValues      DecimalSliceParam
		allowedFaults      int
		faults             int
	)
	err := errors.Join(
		pkgerrors.Wrap(ResolveParam(&maybeAllowedFaults, From(t.AllowedFaults)), "allowedFaults"),
		pkgerrors.Wrap(ResolveParam(&maybePrecision, From(VarExpr(t.Precision, vars), t.Precision)), "precision"),
		pkgerrors.Wrap(ResolveParam(&valuesAndErrs, From(VarExpr(t.Values, vars), JSONWithVarExprs(t.Values, vars, true), Inputs(inputs))), "values"),
	)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	if allowed, isSet := maybeAllowedFaults.Uint64(); isSet {
		allowedFaults = int(allowed)
	} else {
		allowedFaults = len(valuesAndErrs) - 1
	}

	values, faults := valuesAndErrs.FilterErrors()
	if faults > allowedFaults {
		return Result{Error: pkgerrors.Wrapf(ErrTooManyErrors, "Number of faulty inputs %v to mean task > number allowed faults %v", faults, allowedFaults)}, runInfo
	} else if len(values) == 0 {
		return Result{Error: pkgerrors.Wrap(ErrWrongInputCardinality, "values")}, runInfo
	}

	err = decimalValues.UnmarshalPipelineParam(values)
	if err != nil {
		return Result{Error: pkgerrors.Wrapf(ErrBadInput, "values: %v", err)}, runInfo
	}

	total := decimal.NewFromInt(0)
	for _, val := range decimalValues {
		total = total.Add(val)
	}

	numValues := decimal.NewFromInt(int64(len(decimalValues)))

	if precision, isSet := maybePrecision.Int32(); isSet {
		return Result{Value: total.DivRound(numValues, precision)}, runInfo
	}
	// Note that decimal library defaults to rounding to 16 precision
	//https://github.com/shopspring/decimal/blob/2568a29459476f824f35433dfbef158d6ad8618c/decimal.go#L44
	return Result{Value: total.Div(numValues)}, runInfo
}
