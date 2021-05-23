package pipeline

import (
	"context"
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/multierr"
)

type ETHABIDecodeLogTask struct {
	BaseTask `mapstructure:",squash"`
	ABI      string `json:"abi"`
	Log      string `json:"log"`
}

var _ Task = (*ETHABIDecodeLogTask)(nil)

func (t *ETHABIDecodeLogTask) Type() TaskType {
	return TaskTypeETHABIDecodeLog
}

func (t *ETHABIDecodeLogTask) Run(_ context.Context, vars Vars, _ JSONSerializable, inputs []Result) (result Result) {
	_, err := CheckInputs(inputs, 0, 1, 0)
	if err != nil {
		return Result{Error: err}
	}

	var (
		log    MapParam
		theABI StringParam
	)
	err = multierr.Combine(
		vars.ResolveValue(&log, From(VariableExpr(t.Log), Input(inputs, 0))),
		vars.ResolveValue(&theABI, From(NonemptyString(t.ABI))),
	)
	if err != nil {
		return Result{Error: err}
	}

	parts := strings.Split(string(theABI), ",")
	var args, indexedArgs abi.Arguments
	for _, part := range parts {
		part = strings.TrimSpace(part)
		argStr := strings.Split(part, " ")
		argStr = trimStrings(argStr)

		var typeStr, name string
		var indexed bool
		switch len(argStr) {
		case 0, 1:
			return Result{Error: errors.New("bad ABI specification, too few components in argument")}
		case 2:
			if argStr[1] == "indexed" {
				return Result{Error: errors.New("bad ABI specification, too few components in argument")}
			}
			typeStr = argStr[0]
			name = argStr[1]
		case 3:
			if argStr[1] != "indexed" {
				return Result{Error: errors.New("bad ABI specification, unknown component in argument")}
			}
			typeStr = argStr[0]
			name = argStr[2]
			indexed = true
		default:
			return Result{Error: errors.New("bad ABI specification, too many components in argument")}
		}
		typ, err := abi.NewType(typeStr, "", nil)
		if err != nil {
			return Result{Error: err}
		}

		args = append(args, abi.Argument{Type: typ, Name: name, Indexed: indexed})
		if indexed {
			indexedArgs = append(indexedArgs, abi.Argument{Type: typ, Name: name, Indexed: indexed})
		}
	}

	data, ok := log["data"].([]byte)
	if !ok {
		return Result{Error: errors.New("log data: expected []byte")}
	}

	topics, ok := log["topics"].([]common.Hash)
	if !ok {
		return Result{Error: errors.New("log topics: expected []common.Hash")}
	}

	out := make(map[string]interface{})
	if len(data) > 0 {
		if err := args.UnpackIntoMap(out, data); err != nil {
			return Result{Error: err}
		}
	}
	err = abi.ParseTopicsIntoMap(out, indexedArgs, topics[1:])
	if err != nil {
		return Result{Error: err}
	}
	return Result{Value: out}
}
