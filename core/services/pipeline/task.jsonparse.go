package pipeline

import (
	"database/sql/driver"
	"encoding/json"
	"math/big"
	"strings"

	"github.com/pkg/errors"
)

type JSONParseTask struct {
	BaseTask `mapstructure:",squash"`
	Path     JSONPath `json:"path"`
}

var _ Task = (*JSONParseTask)(nil)

func (t *JSONParseTask) Type() TaskType {
	return TaskTypeJSONParse
}

func (t *JSONParseTask) Run(inputs []Result) Result {
	if len(inputs) != 1 {
		return Result{Error: errors.Wrapf(ErrWrongInputCardinality, "JSONParseTask requires a single input")}
	} else if inputs[0].Error != nil {
		return Result{Error: inputs[0].Error}
	}

	var bs []byte
	switch v := inputs[0].Value.(type) {
	case []byte:
		bs = v
	case string:
		bs = []byte(v)
	default:
		return Result{Error: errors.Errorf("JSONParseTask does not accept inputs of type %T", inputs[0].Value)}
	}

	var decoded interface{}
	err := json.Unmarshal(bs, &decoded)
	if err != nil {
		return Result{Error: err}
	}

	for i, part := range t.Path {
		switch d := decoded.(type) {
		case map[string]interface{}:
			var exists bool
			decoded, exists = d[part]
			if !exists && i == len(t.Path)-1 {
				return Result{Value: nil}
			} else if !exists {
				return Result{Error: errors.Errorf(`could not resolve path ["%v"]`, strings.Join(t.Path, `","`))}
			}

		case []interface{}:
			bigindex, ok := big.NewInt(0).SetString(part, 10)
			if !ok {
				return Result{Error: errors.Errorf("JSONParse task error: %v is not a valid array index", part)}
			} else if !bigindex.IsInt64() {
				if i == len(t.Path)-1 {
					return Result{Value: nil}
				} else {
					return Result{Error: errors.Errorf(`could not resolve path ["%v"]`, strings.Join(t.Path, `","`))}
				}
			}
			index := int(bigindex.Int64())
			if index < 0 {
				index = len(d) + index
			}

			exists := index >= 0 && index < len(d)
			if !exists && i == len(t.Path)-1 {
				return Result{Value: nil}
			} else if !exists {
				return Result{Error: errors.Errorf(`could not resolve path ["%v"]`, strings.Join(t.Path, `","`))}
			}
			decoded = d[index]

		default:
			return Result{Error: errors.Errorf(`could not resolve path ["%v"]`, strings.Join(t.Path, `","`))}
		}
	}
	return Result{Value: decoded}
}

type JSONPath []string

func (p *JSONPath) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), p)
}
func (p JSONPath) Value() (driver.Value, error) {
	return json.Marshal(p)
}
