package pipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

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
		return Result{Error: errors.Wrap(ErrBadInput, err.Error())}
	}
	method := abi.NewMethod(methodName, methodName, abi.Function, "", false, false, args, nil)

	var vals []interface{}
	for _, arg := range args {
		val, exists := inputValues[arg.Name]
		if !exists {
			return Result{Error: errors.Wrapf(ErrBadInput, "ETHABIEncode: argument '%v' is missing", arg.Name)}
		}
		val, err = convertToABIType(val, arg.Type)
		if err != nil {
			return Result{Error: errors.Wrapf(ErrBadInput, "ETHABIEncode: %v", err)}
		}
		vals = append(vals, val)
	}

	argsEncoded, err := method.Inputs.Pack(vals...)
	if err != nil {
		return Result{Error: errors.Wrapf(ErrBadInput, "ETHABIEncode: could not ABI encode values: %v", err)}
	}
	dataBytes := append(method.ID, argsEncoded...)
	return Result{Value: dataBytes}
}

func convertToABIType(val interface{}, abiType abi.Type) (interface{}, error) {
	err := checkArrayLengths(reflect.ValueOf(val), abiType)
	if err != nil {
		return nil, err
	}

	bs, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}

	goType := abiType.GetType()
	converted := reflect.New(goType).Interface()
	err = json.Unmarshal(bs, converted)
	return converted, err
}

// JSON marshaling will not fail when a longer array is unmarshaled into a shorter
// one, so we manually check for this case and error if values would otherwise
// be truncated.
func checkArrayLengths(rval reflect.Value, abiType abi.Type) error {
	if rval.Kind() == reflect.Interface {
		rval = rval.Elem()
	}

	switch abiType.T {
	case abi.SliceTy:
		switch rval.Kind() {
		case reflect.Slice, reflect.Array:
			for i := 0; i < rval.Len(); i++ {
				err := checkArrayLengths(rval.Index(i), *abiType.Elem)
				if err != nil {
					return err
				}
			}
		default:
			panic(fmt.Sprintf("invariant violation: abi specifies slice, got %T", rval.Interface()))
		}
	case abi.ArrayTy:
		switch rval.Kind() {
		case reflect.Slice, reflect.Array:
			if abiType.Size != rval.Len() {
				return errors.Wrapf(ErrBadInput, "ETHABIEncode: input array length does not match ABI type")
			}
			for i := 0; i < rval.Len(); i++ {
				err := checkArrayLengths(rval.Index(i), *abiType.Elem)
				if err != nil {
					return err
				}
			}
		default:
			panic(fmt.Sprintf("invariant violation: abi specifies array, got %T (%v)", rval.Interface(), rval.Kind()))
		}
	}
	return nil
}
