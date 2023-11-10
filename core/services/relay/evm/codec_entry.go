package evm

import (
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type CodecEntry struct {
	Args             abi.Arguments
	unwrappedArgs    abi.Arguments
	encodingPrefix   []byte
	checkedType      reflect.Type
	checkedArrayType reflect.Type
	arraySize        int
	nativeType       reflect.Type
}

func (info *CodecEntry) Init() error {
	if info.checkedType != nil {
		return nil
	}

	args := UnwrapArgs(info.Args)
	info.unwrappedArgs = args
	argLen := len(args)
	native := make([]reflect.StructField, argLen)
	checked := make([]reflect.StructField, argLen)

	for i, arg := range args {
		if len(arg.Name) == 0 {
			return relaytypes.InvalidTypeError{}
		}
		nativeArg, checkedArg, err := getNativeAndCheckedTypes(&arg.Type)
		if err != nil {
			return err
		}
		tag := reflect.StructTag(`json:"` + arg.Name + `"`)
		name := strings.ToUpper(arg.Name[:1]) + arg.Name[1:]
		native[i] = reflect.StructField{Name: name, Type: nativeArg, Tag: tag}
		checked[i] = reflect.StructField{Name: name, Type: checkedArg, Tag: tag}
	}

	info.nativeType = reflect.StructOf(native)
	info.checkedType = reflect.StructOf(checked)
	info.checkedArrayType, info.arraySize = getArrayType(checked)
	return nil
}

func UnwrapArgs(args abi.Arguments) abi.Arguments {
	// Unwrap an unnamed tuple so that callers don't need to wrap it
	// Eg: If you have struct Foo { ... } and return an unnamed Foo, you should be able ot decode to a go Foo{} directly
	if len(args) != 1 || args[0].Name != "" {
		return args
	}

	elms := args[0].Type.TupleElems
	if len(elms) != 0 {
		names := args[0].Type.TupleRawNames
		args = make(abi.Arguments, len(elms))
		for i, elm := range elms {
			args[i] = abi.Argument{
				Name: names[i],
				Type: *elm,
			}
		}
	}
	return args
}

func getNativeAndCheckedTypes(curType *abi.Type) (reflect.Type, reflect.Type, error) {
	converter := func(t reflect.Type) reflect.Type { return t }
	for curType.Elem != nil {
		prior := converter
		switch curType.GetType().Kind() {
		case reflect.Slice:
			converter = func(t reflect.Type) reflect.Type {
				return prior(reflect.SliceOf(t))
			}
			curType = curType.Elem
		case reflect.Array:
			tmp := curType
			converter = func(t reflect.Type) reflect.Type {
				return prior(reflect.ArrayOf(tmp.Size, t))
			}
			curType = curType.Elem
		default:
			return nil, nil, relaytypes.InvalidTypeError{}
		}
	}
	base, ok := types.GetType(curType.String())
	if ok {
		return converter(base.Native), converter(base.Checked), nil
	}

	return createTupleType(curType, converter)
}

func createTupleType(curType *abi.Type, converter func(reflect.Type) reflect.Type) (reflect.Type, reflect.Type, error) {
	if len(curType.TupleElems) == 0 {
		return nil, nil, relaytypes.InvalidTypeError{}
	}

	nativeFields := make([]reflect.StructField, len(curType.TupleElems))
	checkedFields := make([]reflect.StructField, len(curType.TupleElems))
	for i, elm := range curType.TupleElems {
		name := curType.TupleRawNames[i]
		nativeFields[i].Name = name
		checkedFields[i].Name = name
		nativeArgType, checkedArgType, err := getNativeAndCheckedTypes(elm)
		if err != nil {
			return nil, nil, err
		}
		nativeFields[i].Type = nativeArgType
		checkedFields[i].Type = checkedArgType
	}
	return converter(reflect.StructOf(nativeFields)), converter(reflect.StructOf(checkedFields)), nil
}

func getArrayType(checked []reflect.StructField) (reflect.Type, int) {
	checkedArray := make([]reflect.StructField, len(checked))
	length := 0
	for i, f := range checked {
		kind := f.Type.Kind()
		if kind == reflect.Slice {
			if i == 0 {
				length = 0
			} else if length != 0 {
				return nil, 0
			}
		} else if kind == reflect.Array {
			if i == 0 {
				length = f.Type.Len()
			} else {
				if f.Type.Len() != length {
					return nil, 0
				}
			}
		} else {
			return nil, 0
		}

		checkedArray[i] = reflect.StructField{Name: f.Name, Type: f.Type.Elem()}
	}
	return reflect.SliceOf(reflect.StructOf(checkedArray)), length
}
