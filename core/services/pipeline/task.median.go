package pipeline

import (
	"context"
	"sort"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/logger"
)

//
// Return types:
//    *decimal.Decimal
//
type MedianTask struct {
	BaseTask      `mapstructure:",squash"`
	Values        string `json:"values"`
	AllowedFaults string `json:"allowedFaults"`
}

var _ Task = (*MedianTask)(nil)

func (t *MedianTask) Type() TaskType {
	return TaskTypeMedian
}

func (t *MedianTask) Run(_ context.Context, _ logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	var (
		maybeAllowedFaults MaybeUint64Param
		valuesAndErrs      SliceParam
		decimalValues      DecimalSliceParam
		allowedFaults      int
		faults             int
	)
	err := multierr.Combine(
		errors.Wrap(ResolveParam(&maybeAllowedFaults, From(t.AllowedFaults)), "allowedFaults"),
		errors.Wrap(ResolveParam(&valuesAndErrs, From(VarExpr(t.Values, vars), JSONWithVarExprs(t.Values, vars, true), Inputs(inputs))), "values"),
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
		return Result{Error: errors.Wrapf(ErrTooManyErrors, "Number of faulty inputs %v to median task > number allowed faults %v", faults, allowedFaults)}, runInfo
	} else if len(values) == 0 {
		return Result{Error: errors.Wrap(ErrWrongInputCardinality, "no values to medianize")}, runInfo
	}

	err = decimalValues.UnmarshalPipelineParam(values)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	sort.Slice(decimalValues, func(i, j int) bool {
		return decimalValues[i].LessThan(decimalValues[j])
	})
	k := len(decimalValues) / 2
	if len(decimalValues)%2 == 1 {
		return Result{Value: decimalValues[k]}, runInfo
	}
	median := decimalValues[k].Add(decimalValues[k-1]).Div(decimal.NewFromInt(2))
	return Result{Value: median}, runInfo
}
