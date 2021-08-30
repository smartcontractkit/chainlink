package pipeline

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/store/models"
)

//
// Return types:
//     map[string]interface{} with potential value types:
//         float64
//         string
//         bool
//         map[string]interface{}
//         []interface{}
//         nil
//
type CBORParseTask struct {
	BaseTask `mapstructure:",squash"`
	Data     string `json:"data"`
}

var _ Task = (*CBORParseTask)(nil)

func (t *CBORParseTask) Type() TaskType {
	return TaskTypeCBORParse
}

func (t *CBORParseTask) Run(_ context.Context, vars Vars, inputs []Result) (result Result) {
	_, err := CheckInputs(inputs, -1, -1, 0)
	if err != nil {
		return Result{Error: errors.Wrap(err, "task inputs")}
	}

	var (
		data BytesParam
	)
	err = multierr.Combine(
		errors.Wrap(ResolveParam(&data, From(VarExpr(t.Data, vars))), "data"),
	)
	if err != nil {
		return Result{Error: err}
	}

	parsed, err := models.ParseCBOR([]byte(data))
	if err != nil {
		return Result{Error: errors.Wrapf(ErrBadInput, "CBORParse: data: %v", err)}
	}
	m, ok := parsed.Result.Value().(map[string]interface{})
	if !ok {
		return Result{Error: errors.Wrapf(ErrBadInput, "CBORParse: data: expected map[string]interface{}, got %T", parsed.Result.Value())}
	}
	return Result{Value: m}
}
