package pipeline

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/logger"
)

//
// Return types:
//     map[string]interface{} with any geth/abigen value type
//
type ETHABIDecodeTask struct {
	BaseTask `mapstructure:",squash"`
	ABI      string `json:"abi"`
	Data     string `json:"data"`
}

var _ Task = (*ETHABIDecodeTask)(nil)

func (t *ETHABIDecodeTask) Type() TaskType {
	return TaskTypeETHABIDecode
}

func (t *ETHABIDecodeTask) Run(_ context.Context, _ logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	_, err := CheckInputs(inputs, 0, 1, 0)
	if err != nil {
		return Result{Error: errors.Wrap(err, "task inputs")}, runInfo
	}

	var (
		data   BytesParam
		theABI BytesParam
	)
	err = multierr.Combine(
		errors.Wrap(ResolveParam(&data, From(VarExpr(t.Data, vars), Input(inputs, 0))), "data"),
		errors.Wrap(ResolveParam(&theABI, From(NonemptyString(t.ABI))), "abi"),
	)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	args, _, err := ParseETHABIArgsString([]byte(theABI), false)
	if err != nil {
		return Result{Error: errors.Wrap(ErrBadInput, err.Error())}, runInfo
	}

	out := make(map[string]interface{})
	if len(data) > 0 {
		if err := args.UnpackIntoMap(out, []byte(data)); err != nil {
			return Result{Error: err}, runInfo
		}
	}
	return Result{Value: out}, runInfo
}
