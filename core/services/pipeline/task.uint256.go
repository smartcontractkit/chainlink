package pipeline

import (
	"context"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/utils"
)

type UInt256Task struct {
	BaseTask `mapstructure:",squash"`
}

var _ Task = (*UInt256Task)(nil)

func (t *UInt256Task) Type() TaskType {
	return TaskTypeUInt256
}

func (t *UInt256Task) SetDefaults(inputValues map[string]string, g TaskDAG, self taskDAGNode) error {
	return nil
}

func (t *UInt256Task) Run(_ context.Context, _ JSONSerializable, inputs []Result) (result Result) {
	if len(inputs) != 1 {
		return Result{Error: errors.Wrapf(ErrWrongInputCardinality, "Uint256Task requires a single input")}
	} else if inputs[0].Error != nil {
		return Result{Error: inputs[0].Error}
	}

	value, err := utils.ToDecimal(inputs[0].Value)
	if err != nil {
		return Result{Error: err}
	}

	evmByteArray, err := utils.EVMWordBigInt(value.BigInt())
	if err != nil {
		return Result{Error: err}
	}

	return Result{Value: evmByteArray}
}
