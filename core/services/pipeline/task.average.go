package pipeline

// import (
// 	"context"

// 	"github.com/pkg/errors"
// 	"github.com/shopspring/decimal"
// 	"go.uber.org/multierr"

// 	"github.com/smartcontractkit/chainlink/core/utils"
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

// func (t *AverageTask) SetDefaults(inputValues map[string]string, g TaskDAG, self TaskDAGNode) error {
// 	if _, exists := inputValues["allowedFaults"]; !exists {
// 		if len(self.inputs()) == 0 {
// 			return errors.Wrapf(ErrWrongInputCardinality, "AverageTask requires at least 1 input")
// 		}
// 		t.AllowedFaults = uint64(len(self.inputs()) - 1)
// 	}
// 	return nil
// }

// func (t *AverageTask) Run(_ context.Context, _ JSONSerializable, inputs []Result) (result Result) {
// 	if len(inputs) == 0 {
// 		return Result{Error: errors.Wrapf(ErrWrongInputCardinality, "AverageTask requires at least 1 input")}
// 	}

// 	fetchErrors := []error{}
// 	total := decimal.New(0, int32(t.Precision))
// 	validInputs := 0

// 	for _, input := range inputs {
// 		if input.Error != nil {
// 			fetchErrors = append(fetchErrors, input.Error)
// 			continue
// 		}

// 		answer, err := utils.ToDecimal(input.Value)
// 		if err != nil {
// 			fetchErrors = append(fetchErrors, err)
// 			continue
// 		}
// 		total = total.Add(answer)
// 		validInputs++
// 	}

// 	if uint64(len(fetchErrors)) > t.AllowedFaults {
// 		return Result{Error: errors.Wrapf(ErrBadInput, "Number of faulty inputs %v to average task > number allowed faults %v. Fetch errors: %v", len(fetchErrors), t.AllowedFaults, multierr.Combine(fetchErrors...).Error())}
// 	}

// 	average := total.Div(decimal.NewFromInt(int64(validInputs)))
// 	return Result{Value: average}
// }
