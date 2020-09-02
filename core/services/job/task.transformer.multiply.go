package job

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink/core/utils"
)

type MultiplyTransformer struct {
	BaseTask
	Times decimal.Decimal `json:"times" gorm:"type:text;not null"`
}

var _ Task = (*MultiplyTransformer)(nil)

func (t *MultiplyTransformer) Run(inputs []Result) (out interface{}, err error) {
	if len(inputs) != 1 {
		return nil, errors.Wrapf(ErrWrongInputCardinality, "MultiplyTransformer requires a single input")
	} else if inputs[0].Error != nil {
		return nil, inputs[0].Error
	}

	value, err := utils.ToDecimal(inputs[0].Value)
	if err != nil {
		return nil, err
	}
	return value.Mul(t.Times), nil
}

func (t MultiplyTransformer) MarshalJSON() ([]byte, error) {
	type preventInfiniteRecursion MultiplyTransformer
	type transformerWithType struct {
		Type TransformerType `json:"type"`
		preventInfiniteRecursion
	}
	return json.Marshal(transformerWithType{
		TransformerTypeMultiply,
		preventInfiniteRecursion(t),
	})
}
