package pipeline

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

//
// Return types:
//     map[string]interface{} with any geth/abigen value type
//
type ETHABIDecodeLogTask struct {
	BaseTask `mapstructure:",squash"`
	ABI      string `json:"abi"`
	Data     string `json:"data"`
	Topics   string `json:"topics"`
}

var _ Task = (*ETHABIDecodeLogTask)(nil)

func (t *ETHABIDecodeLogTask) Type() TaskType {
	return TaskTypeETHABIDecodeLog
}

func (t *ETHABIDecodeLogTask) Run(_ context.Context, vars Vars, inputs []Result) (result Result) {
	_, err := CheckInputs(inputs, -1, -1, 0)
	if err != nil {
		return Result{Error: errors.Wrap(err, "task inputs")}
	}

	var (
		theABI BytesParam
		data   BytesParam
		topics HashSliceParam
	)
	err = multierr.Combine(
		errors.Wrap(ResolveParam(&data, From(VarExpr(t.Data, vars), nil)), "data"),
		errors.Wrap(ResolveParam(&topics, From(VarExpr(t.Topics, vars))), "topics"),
		errors.Wrap(ResolveParam(&theABI, From(NonemptyString(t.ABI))), "abi"),
	)
	if err != nil {
		return Result{Error: err}
	}

	_, args, indexedArgs, err := parseETHABIString([]byte(theABI), true)
	if err != nil {
		return Result{Error: errors.Wrap(ErrBadInput, err.Error())}
	}

	out := make(map[string]interface{})
	if len(data) > 0 {
		if err2 := args.UnpackIntoMap(out, []byte(data)); err2 != nil {
			return Result{Error: errors.Wrap(ErrBadInput, err2.Error())}
		}
	}
	if len(indexedArgs) > 0 {
		if len(topics) != len(indexedArgs)+1 {
			return Result{Error: errors.Wrap(ErrBadInput, "topic/field count mismatch")}
		}
		err = abi.ParseTopicsIntoMap(out, indexedArgs, topics[1:])
		if err != nil {
			return Result{Error: errors.Wrap(ErrBadInput, err.Error())}
		}
	}
	return Result{Value: out}
}
