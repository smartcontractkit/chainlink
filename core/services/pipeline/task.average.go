package pipeline

// import (
// 	"context"

// 	"github.com/pkg/errors"
// 	"github.com/shopspring/decimal"
// 	"go.uber.org/multierr"
// )

// type AverageTask struct {
// 	BaseTask      `mapstructure:",squash"`
// 	Precision     uint8  `json:"precision"`
// 	AllowedFaults uint64 `json:"allowedFaults"`
// }

// var _ Task = (*AverageTask)(nil)

// func (t *AverageTask) Type() TaskType {
// 	return TaskTypeAverage
// }

// func (t *AverageTask) Run(_ context.Context, _ JSONSerializable, inputs []Result) (result Result) {
// 	var (
// 		maybeAllowedFaults MaybeUint64Param
// 		valuesAndErrs      SliceParam
// 		decimalValues      DecimalSliceParam
// 		allowedFaults      int
// 		faults             int
// 	)
// 	err := multierr.Combine(
// 		vars.ResolveValue(&maybeAllowedFaults, From(t.AllowedFaults)),
// 		vars.ResolveValue(&valuesAndErrs, From(VariableExpr(t.Values), Inputs(inputs))),
// 	)
// 	if err != nil {
// 		return Result{Error: err}
// 	}

// 	if allowed, isSet := maybeAllowedFaults.Uint64(); isSet {
// 		allowedFaults = int(allowed)
// 	} else {
// 		allowedFaults = len(valuesAndErrs) - 1
// 	}

// 	values, faults := valuesAndErrs.FilterErrors()
// 	if faults > allowedFaults {
// 		return Result{Error: errors.Wrapf(ErrTooManyErrors, "Number of faulty inputs %v to average task > number allowed faults %v", faults, allowedFaults)}
// 	} else if len(values) == 0 {
// 		return Result{Error: errors.Wrap(ErrWrongInputCardinality, "no values to average")}
// 	}

// 	err = decimalValues.UnmarshalPipelineParam(values, nil)
// 	if err != nil {
// 		return Result{Error: err}
// 	}

// 	total := decimal.NewFromFloat(0)
// 	for _, val := range decimalValues {
// 		total = total.Add(answer)
// 	}
// 	average := total.Div(decimal.NewFromInt(len(decimalValues)))
// 	return Result{Value: average}
// }
