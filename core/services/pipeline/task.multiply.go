package pipeline

import (
	"context"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink/core/utils"
)

type MultiplyTask struct {
	BaseTask `mapstructure:",squash"`
	Times    decimal.Decimal `json:"times"`
}

var _ Task = (*MultiplyTask)(nil)

func (t *MultiplyTask) Type() TaskType {
	return TaskTypeMultiply
}

func (t *MultiplyTask) SetDefaults(inputValues map[string]string, g TaskDAG, self TaskDAGNode) error {
	return nil
}

func (t *MultiplyTask) Run(_ context.Context, _ JSONSerializable, inputs []Result) (result Result) {
	if len(inputs) != 1 {
		return Result{Error: errors.Wrapf(ErrWrongInputCardinality, "MultiplyTask requires a single input")}
	} else if inputs[0].Error != nil {
		return Result{Error: inputs[0].Error}
	}

	value, err := utils.ToDecimal(inputs[0].Value)
	if err != nil {
		return Result{Error: err}
	}
	return Result{Value: value.Mul(t.Times)}
}
