package job

import (
	"encoding/json"
	"math/big"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink/core/utils"
)

type Transformers []Transformer

func (t Transformers) Run(input interface{}) (interface{}, error) {
	var err error
	for _, transformer := range t {
		input, err = transformer.Transform(input)
		if err != nil {
			return nil, err
		}
	}
	return input, nil
}

func (t *Transformers) UnmarshalJSON(bs []byte) error {
	var rawJSON []json.RawMessage
	err := json.Unmarshal(bs, &rawJSON)
	if err != nil {
		return err
	}

	for _, x := range rawJSON {
		var y struct {
			Type string `json:"type"`
		}
		err := json.Unmarshal(x, &y)
		if err != nil {
			return err
		}
		var transformer Transformer
		switch y.Type {
		case TransformerTypeMultiply:
			transformer = MultiplyTransformer{}
		default:
			return errors.Errorf("invalid transformer type '%v'", y.Type)
		}
		*t = append(*t, transformer)
	}
	return nil
}

type Transformer interface {
	Transform(input interface{}) (interface{}, error)
}

type TransformerType string

var (
	TransformerTypeJSONParse TransformerType = "jsonparse"
	TransformerTypeMultiply  TransformerType = "multiply"
)

type JSONParseTransformer struct {
	Path []string `json:"path"`
}

func (t JSONParseTransformer) Transform(input interface{}) (interface{}, error) {
	var bs []byte
	switch v := input.(type) {
	case []byte:
		bs = v
	case string:
		bs = []byte(v)
	default:
		return errors.Errorf("JSONParseTransformer does not accept inputs of type %T", input)
	}

	var decoded interface{}
	err := json.Unmarshal(bs, &decoded)
	if err != nil {
		return nil, err
	}

	for i, part := range t.Path {
		switch d := decoded.(type) {
		case map[string]interface{}:
			var exists bool
			decoded, exists = d[part]
			if !exists && i == len(t.Path)-1 {
				return nil, nil
			} else if !exists {
				return nil, errors.Errorf(`could not resolve path ["%v"]`, strings.Join(t.Path, `","`))
			}

		case []interface{}:
			index, err := strconv.Atoi(part)
			if err != nil {
				return nil, err
			}
			exists := index < len(d)
			if !exists && i == len(t.Path)-1 {
				return nil, nil
			} else if !exists {
				return nil, errors.Errorf(`could not resolve path ["%v"]`, strings.Join(t.Path, `","`))
			}
			decoded = d[index]

		default:
			return nil, errors.Errorf(`could not resolve path ["%v"]`, strings.Join(t.Path, `","`))
		}
	}
	return decoded, nil
}

type MultiplyTransformer struct {
	Multiplier decimal.Decimal `json:"times"`
}

func (t MultiplyTransformer) Transform(input interface{}) (interface{}, error) {
	value, err := utils.ToDecimal(input)
	if err != nil {
		return nil, err
	}
	return value.Mul(t.Multiplier), nil
}
