package pipeline

import (
	"context"
	"encoding/json"
	"math/big"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

type JSONParseTask struct {
	BaseTask `mapstructure:",squash"`
	Path     string `json:"path"`
	Data     string `json:"data"`
	// Lax when disabled will return an error if the path does not exist
	// Lax when enabled will return nil with no error if the path does not exist
	Lax string
}

var _ Task = (*JSONParseTask)(nil)

func (t *JSONParseTask) Type() TaskType {
	return TaskTypeJSONParse
}

func (t *JSONParseTask) SetDefaults(inputValues map[string]string, g TaskDAG, self TaskDAGNode) error {
	return nil
}

func (t *JSONParseTask) Run(_ context.Context, vars Vars, _ JSONSerializable, inputs []Result) (result Result) {
	_, err := CheckInputs(inputs, 0, 1, 0)
	if err != nil {
		return Result{Error: err}
	}

	var (
		path StringSliceParam
		data StringParam
		lax  BoolParam
	)
	err = multierr.Combine(
		vars.ResolveValue(&path, From(NonemptyString(t.Path))),
		vars.ResolveValue(&data, From(VariableExpr(t.Data), NonemptyString(t.Data), Input(inputs, 0))),
		vars.ResolveValue(&lax, From(NonemptyString(t.Lax))),
	)
	if err != nil {
		return Result{Error: err}
	}

	var decoded interface{}
	err = json.Unmarshal([]byte(data), &decoded)
	if err != nil {
		return Result{Error: err}
	}

	for _, part := range path {
		switch d := decoded.(type) {
		case map[string]interface{}:
			var exists bool
			decoded, exists = d[part]
			if !exists && bool(lax) {
				return Result{Value: nil}
			} else if !exists {
				return Result{Error: errors.Errorf(`could not resolve path ["%v"] in %s`, strings.Join(path, `","`), data)}
			}

		case []interface{}:
			bigindex, ok := big.NewInt(0).SetString(part, 10)
			if !ok {
				return Result{Error: errors.Errorf("JSONParse task error: %v is not a valid array index", part)}
			} else if !bigindex.IsInt64() {
				if bool(lax) {
					return Result{Value: nil}
				}
				return Result{Error: errors.Errorf(`could not resolve path ["%v"] in %s`, strings.Join(path, `","`), data)}
			}
			index := int(bigindex.Int64())
			if index < 0 {
				index = len(d) + index
			}

			exists := index >= 0 && index < len(d)
			if !exists && bool(lax) {
				return Result{Value: nil}
			} else if !exists {
				return Result{Error: errors.Errorf(`could not resolve path ["%v"] in %s`, strings.Join(path, `","`), data)}
			}
			decoded = d[index]

		default:
			return Result{Error: errors.Errorf(`could not resolve path ["%v"] in %s`, strings.Join(path, `","`), data)}
		}
	}
	return Result{Value: decoded}
}
