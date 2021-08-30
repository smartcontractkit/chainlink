package pipeline

import (
	"context"
	"sort"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"go.uber.org/multierr"
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

func (t *MedianTask) Run(_ context.Context, vars Vars, inputs []Result) (result Result) {
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
		return Result{Error: err}
	}

	if allowed, isSet := maybeAllowedFaults.Uint64(); isSet {
		allowedFaults = int(allowed)
	} else {
		allowedFaults = len(valuesAndErrs) - 1
	}

	values, faults := valuesAndErrs.FilterErrors()
	if faults > allowedFaults {
		return Result{Error: errors.Wrapf(ErrTooManyErrors, "Number of faulty inputs %v to median task > number allowed faults %v", faults, allowedFaults)}
	} else if len(values) == 0 {
		return Result{Error: errors.Wrap(ErrWrongInputCardinality, "no values to medianize")}
	}

	err = decimalValues.UnmarshalPipelineParam(values)
	if err != nil {
		return Result{Error: err}
	}

	sort.Slice(decimalValues, func(i, j int) bool {
		return decimalValues[i].LessThan(decimalValues[j])
	})
	k := len(decimalValues) / 2
	if len(decimalValues)%2 == 1 {
		return Result{Value: decimalValues[k]}
	}
	median := decimalValues[k].Add(decimalValues[k-1]).Div(decimal.NewFromInt(2))
	return Result{Value: median}
}
