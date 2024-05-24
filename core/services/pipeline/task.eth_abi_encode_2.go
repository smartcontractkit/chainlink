package pipeline

import (
	"context"
	"encoding/json"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// Return types:
//
//	[]byte
type ETHABIEncodeTask2 struct {
	BaseTask `mapstructure:",squash"`
	ABI      string `json:"abi"`
	Data     string `json:"data"`
}

var _ Task = (*ETHABIEncodeTask2)(nil)

func (t *ETHABIEncodeTask2) Type() TaskType {
	return TaskTypeETHABIEncode
}

func (t *ETHABIEncodeTask2) Run(_ context.Context, _ logger.Logger, vars Vars, inputs []Result) (Result, RunInfo) {
	_, err := CheckInputs(inputs, -1, -1, 0)
	if err != nil {
		return Result{Error: errors.Wrap(err, "task inputs")}, RunInfo{}
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
		return Result{Error: err}, RunInfo{}
	}

	inputMethod := Method{}
	err = json.Unmarshal(theABI, &inputMethod)
	if err != nil {
		return Result{Error: errors.Wrapf(ErrBadInput, "ETHABIEncode: while parsing ABI string: %v", err)}, RunInfo{}
	}

	method := abi.NewMethod(inputMethod.Name, inputMethod.Name, abi.Function, "", false, false, inputMethod.Inputs, nil)

	var vals []interface{}
	for _, arg := range method.Inputs {
		if len(arg.Name) == 0 {
			return Result{Error: errors.Wrapf(ErrBadInput, "ETHABIEncode: bad ABI specification, missing argument name")}, RunInfo{}
		}
		val, exists := inputValues[arg.Name]
		if !exists {
			return Result{Error: errors.Wrapf(ErrBadInput, "ETHABIEncode: argument '%v' is missing", arg.Name)}, RunInfo{}
		}
		val, err = convertToETHABIType(val, arg.Type)
		if err != nil {
			return Result{Error: errors.Wrapf(ErrBadInput, "ETHABIEncode: while converting argument '%v' from %T to %v: %v", arg.Name, val, arg.Type, err)}, RunInfo{}
		}
		vals = append(vals, val)
	}

	argsEncoded, err := method.Inputs.Pack(vals...)
	if err != nil {
		return Result{Error: errors.Wrapf(ErrBadInput, "ETHABIEncode: could not ABI encode values: %v", err)}, RunInfo{}
	}
	var dataBytes []byte
	if method.Name != "" {
		dataBytes = append(method.ID, argsEncoded...)
	} else {
		dataBytes = argsEncoded
	}
	return Result{Value: hexutil.Encode(dataBytes)}, RunInfo{}
}

// go-ethereum's abi.Method doesn't implement json.Marshal for Type, but
// otherwise would have worked fine, in any case we only care about these...
type Method struct {
	Name   string
	Inputs abi.Arguments
}
