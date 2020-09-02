package job

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type JSONParseTransformer struct {
	BaseTask
	Path JSONPath `json:"path" gorm:"type:jsonb"`
}

var _ Task = (*JSONParseTransformer)(nil)

func (t *JSONParseTransformer) Run(inputs []Result) (out interface{}, err error) {
	if len(inputs) != 1 {
		return nil, errors.Wrapf(ErrWrongInputCardinality, "JSONParseTransformer requires a single input")
	} else if inputs[0].Error != nil {
		return nil, inputs[0].Error
	}

	var bs []byte
	switch v := inputs[0].Value.(type) {
	case []byte:
		bs = v
	case string:
		bs = []byte(v)
	default:
		return nil, errors.Errorf("JSONParseTransformer does not accept inputs of type %T", inputs[0].Value)
	}

	var decoded interface{}
	err = json.Unmarshal(bs, &decoded)
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

type JSONPath []string

func (p *JSONPath) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), p)
}
func (p JSONPath) Value() (driver.Value, error) {
	return json.Marshal(p)
}
