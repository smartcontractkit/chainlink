package job

import (
	"encoding/json"

	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink/core/utils"
)

type MultiplyTransformer struct {
	BaseTransformer
	Times decimal.Decimal `json:"times" gorm:"type:text;not null"`
}

var (
	_ Transformer   = (*MultiplyTransformer)(nil)
	_ PipelineStage = (*MultiplyTransformer)(nil)
)

func (t *MultiplyTransformer) Transform(input interface{}) (out interface{}, err error) {
	defer func() { t.notifiee.OnEndStage(t, out, err) }()
	t.notifiee.OnBeginStage(t, input)

	value, err := utils.ToDecimal(input)
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
