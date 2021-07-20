package pipeline

import (
	"context"
	"encoding/json"
	"math/big"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

//
// Return types:
//     float64
//     string
//     bool
//     map[string]interface{}
//     []interface{}
//     nil
//
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

func (t *JSONParseTask) Run(_ context.Context, vars Vars, inputs []Result) (result Result) {
	_, err := CheckInputs(inputs, 0, 1, 0)
	if err != nil {
		return Result{Error: errors.Wrap(err, "task inputs")}
	}

	var (
		path JSONPathParam
		data StringParam
		lax  BoolParam
	)
	err = multierr.Combine(
		errors.Wrap(ResolveParam(&path, From(VarExpr(t.Path, vars), t.Path)), "path"),
		errors.Wrap(ResolveParam(&data, From(VarExpr(t.Data, vars), Input(inputs, 0))), "data"),
		errors.Wrap(ResolveParam(&lax, From(NonemptyString(t.Lax), false)), "lax"),
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
				decoded = nil
				break
			} else if !exists {
				return Result{Error: errors.Wrapf(ErrKeypathNotFound, `could not resolve path ["%v"] in %s`, strings.Join(path, `","`), data)}
			}

		case []interface{}:
			bigindex, ok := big.NewInt(0).SetString(part, 10)
			if !ok {
				return Result{Error: errors.Wrapf(ErrKeypathNotFound, "JSONParse task error: %v is not a valid array index", part)}
			} else if !bigindex.IsInt64() {
				if bool(lax) {
					decoded = nil
					break
				}
				return Result{Error: errors.Wrapf(ErrKeypathNotFound, `could not resolve path ["%v"] in %s`, strings.Join(path, `","`), data)}
			}
			index := int(bigindex.Int64())
			if index < 0 {
				index = len(d) + index
			}

			exists := index >= 0 && index < len(d)
			if !exists && bool(lax) {
				decoded = nil
				break
			} else if !exists {
				return Result{Error: errors.Wrapf(ErrKeypathNotFound, `could not resolve path ["%v"] in %s`, strings.Join(path, `","`), data)}
			}
			decoded = d[index]

		default:
			return Result{Error: errors.Wrapf(ErrKeypathNotFound, `could not resolve path ["%v"] in %s`, strings.Join(path, `","`), data)}
		}
	}
	return Result{Value: decoded}
}
