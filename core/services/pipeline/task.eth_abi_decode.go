package pipeline

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

type ETHABIDecodeTask struct {
	BaseTask `mapstructure:",squash"`
	ABI      string `json:"abi"`
	Data     string `json:"data"`
}

var _ Task = (*ETHABIDecodeTask)(nil)

func (t *ETHABIDecodeTask) Type() TaskType {
	return TaskTypeETHABIDecode
}

func (t *ETHABIDecodeTask) Run(_ context.Context, vars Vars, inputs []Result) Result {
	_, err := CheckInputs(inputs, -1, -1, 0)
	if err != nil {
		return Result{Error: errors.Wrap(err, "task inputs")}
	}

	var (
		data   BytesParam
		theABI BytesParam
	)
	err = multierr.Combine(
		errors.Wrap(ResolveParam(&data, From(VarExpr(t.Data, vars))), "data"),
		errors.Wrap(ResolveParam(&theABI, From(NonemptyString(t.ABI))), "abi"),
	)
	if err != nil {
		return Result{Error: err}
	}

	args, _, err := parseETHABIArgsString([]byte(theABI), false)
	if err != nil {
		return Result{Error: errors.Wrap(ErrBadInput, err.Error())}
	}

	out := make(map[string]interface{})
	if len(data) > 0 {
		if err := args.UnpackIntoMap(out, []byte(data)); err != nil {
			return Result{Error: err}
		}
	}
	return Result{Value: out}
}
