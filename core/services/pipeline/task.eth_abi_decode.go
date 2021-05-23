package pipeline

import (
	"context"
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
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

func (t *ETHABIDecodeTask) Run(_ context.Context, vars Vars, _ JSONSerializable, inputs []Result) Result {
	_, err := CheckInputs(inputs, 0, 1, 0)
	if err != nil {
		return Result{Error: err}
	}

	var (
		data   BytesParam
		theABI StringParam
	)
	err = multierr.Combine(
		vars.ResolveValue(&data, From(VariableExpr(t.Data), Input(inputs, 0))),
		vars.ResolveValue(&theABI, From(NonemptyString(t.ABI))),
	)
	if err != nil {
		return Result{Error: err}
	}

	parts := strings.Split(string(theABI), ",")
	var args abi.Arguments
	for _, part := range parts {
		part = strings.TrimSpace(part)
		argStr := strings.Split(part, " ")
		argStr = trimStrings(argStr)

		var typeStr, name string
		switch len(argStr) {
		case 0, 1:
			return Result{Error: errors.New("bad ABI specification, too few components in argument")}
		case 2:
			typeStr = argStr[0]
			name = argStr[1]
		default:
			return Result{Error: errors.New("bad ABI specification, too many components in argument")}
		}
		typ, err := abi.NewType(typeStr, "", nil)
		if err != nil {
			return Result{Error: err}
		}

		args = append(args, abi.Argument{Type: typ, Name: name})
	}

	out := make(map[string]interface{})
	if len(data) > 0 {
		if err := args.UnpackIntoMap(out, []byte(data)); err != nil {
			return Result{Error: err}
		}
	}

	err = vars.Set(t.DotID(), out)
	if err != nil {
		return Result{Error: err}
	}
	return Result{Value: out}
}
