package pipeline

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

//
// Return types:
//     []byte
//
type ETHABIEncodeTask struct {
	BaseTask `mapstructure:",squash"`
	ABI      string `json:"abi"`
	Data     string `json:"data"`
}

var _ Task = (*ETHABIEncodeTask)(nil)

func (t *ETHABIEncodeTask) Type() TaskType {
	return TaskTypeETHABIEncode
}

func (t *ETHABIEncodeTask) Run(_ context.Context, vars Vars, inputs []Result) (result Result) {
	_, err := CheckInputs(inputs, -1, -1, 0)
	if err != nil {
		return Result{Error: errors.Wrap(err, "task inputs")}
	}

	var (
		inputValues MapParam
		theABI      BytesParam
	)
	err = multierr.Combine(
		errors.Wrap(ResolveParam(&inputValues, From(VarExpr(t.Data, vars), JSONWithVarExprs(t.Data, vars, false), nil)), "data"),
		errors.Wrap(ResolveParam(&theABI, From(NonemptyString(t.ABI))), "abi"),
	)
	if err != nil {
		return Result{Error: err}
	}

	methodName, args, _, err := parseETHABIString([]byte(theABI), false)
	if err != nil {
		return Result{Error: errors.Wrapf(ErrBadInput, "ETHABIEncode: while parsing ABI string: %v", err)}
	}
	method := abi.NewMethod(methodName, methodName, abi.Function, "", false, false, args, nil)

	var vals []interface{}
	for _, arg := range args {
		val, exists := inputValues[arg.Name]
		if !exists {
			return Result{Error: errors.Wrapf(ErrBadInput, "ETHABIEncode: argument '%v' is missing", arg.Name)}
		}
		val, err = convertToETHABIType(val, arg.Type)
		if err != nil {
			return Result{Error: errors.Wrapf(ErrBadInput, "ETHABIEncode: while converting argument '%v' from %T to %v: %v", arg.Name, val, arg.Type, err)}
		}
		vals = append(vals, val)
	}

	argsEncoded, err := method.Inputs.Pack(vals...)
	if err != nil {
		return Result{Error: errors.Wrapf(ErrBadInput, "ETHABIEncode: could not ABI encode values: %v", err)}
	}
	var dataBytes []byte
	if methodName != "" {
		dataBytes = append(method.ID, argsEncoded...)
	} else {
		dataBytes = argsEncoded
	}
	return Result{Value: hexutil.Encode(dataBytes)}
}
