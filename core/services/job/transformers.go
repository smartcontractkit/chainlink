package job

import (
	"encoding/json"
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

func withStack(err *error) {
	*err = errors.WithStack(*err)
}

func (t *Transformers) UnmarshalJSON(bs []byte) (err error) {
	defer withStack(&err)

	var rawJSON []json.RawMessage
	err = json.Unmarshal(bs, &rawJSON)
	if err != nil {
		return err
	}

	for _, bs := range rawJSON {
		var header struct {
			Type TransformerType `json:"type"`
		}
		err := json.Unmarshal(bs, &header)
		if err != nil {
			return err
		}
		var transformer Transformer
		switch header.Type {
		case TransformerTypeJSONParse:
			jsonTransformer := JSONParseTransformer{}
			err = json.Unmarshal(bs, &jsonTransformer)
			if err != nil {
				return err
			}
			transformer = jsonTransformer

		case TransformerTypeMultiply:
			multiplyTransformer := MultiplyTransformer{}
			err = json.Unmarshal(bs, &multiplyTransformer)
			if err != nil {
				return err
			}
			transformer = multiplyTransformer

		default:
			return errors.Errorf("invalid transformer type '%v'", header.Type)
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
	ID   uint64   `json:"-" gorm:"primary_key;auto_increment"`
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
		return nil, errors.Errorf("JSONParseTransformer does not accept inputs of type %T", input)
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
			if index < 0 {
				index = len(d) + index
			}

			exists := index >= 0 && index < len(d)
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

func (t JSONParseTransformer) MarshalJSON() ([]byte, error) {
	type preventInfiniteRecursion JSONParseTransformer
	type transformerWithType struct {
		Type TransformerType `json:"type"`
		preventInfiniteRecursion
	}
	return json.Marshal(transformerWithType{
		TransformerTypeJSONParse,
		preventInfiniteRecursion(t),
	})
}

type MultiplyTransformer struct {
	ID    uint64          `json:"-" gorm:"primary_key;auto_increment"`
	Times decimal.Decimal `json:"times"`
}

func (t MultiplyTransformer) Transform(input interface{}) (interface{}, error) {
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
