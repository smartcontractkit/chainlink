package job

import (
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink/core/utils"
)

type MultiplyTask struct {
	BaseTask
	Times decimal.Decimal `json:"times" gorm:"type:text;not null"`
}

var _ Task = (*MultiplyTask)(nil)

func (t *MultiplyTask) Run(inputs []Result) (out interface{}, err error) {
	if len(inputs) != 1 {
		return nil, errors.Wrapf(ErrWrongInputCardinality, "MultiplyTask requires a single input")
	} else if inputs[0].Error != nil {
		return nil, inputs[0].Error
	}

	value, err := utils.ToDecimal(inputs[0].Value)
	if err != nil {
		return nil, err
	}
	return value.Mul(t.Times), nil
}
