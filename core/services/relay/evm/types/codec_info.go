package types

import (
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
)

type CodecInfo struct {
	Definitions map[string]*CodecEntry
}

func (c *CodecInfo) init() error {
	for _, v := range c.Definitions {
		if err := v.Init(); err != nil {
			return err
		}
	}
	return nil
}

type CodecEntry struct {
	Args           abi.Arguments
	EncodingPrefix []byte

	CheckedType      reflect.Type
	CheckedArrayType reflect.Type
	ArraySize        int
	NativeType       reflect.Type
}

func (info *CodecEntry) Init() error {
	if info.CheckedType != nil {
		return nil
	}

	args := info.Args
	argLen := len(args)
	native := make([]reflect.StructField, argLen)
	checked := make([]reflect.StructField, argLen)
	for i, arg := range args {
		nativeArg, checkedArg, err := getNativeAndCheckedTypes(&arg.Type)
		if err != nil {
			return err
		}

		native[i] = reflect.StructField{Name: arg.Name, Type: nativeArg}
		checked[i] = reflect.StructField{Name: arg.Name, Type: checkedArg}
	}

	info.NativeType = reflect.StructOf(native)
	info.CheckedType = reflect.StructOf(checked)
	info.CheckedArrayType, info.ArraySize = getArrayType(checked)
	return nil
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
			return nil, nil, types.InvalidTypeError{}
		}
	}
	base, ok := typeMap[curType.String()]
	if ok {
		return converter(base.Native), converter(base.Checked), nil
	}

	return createTupleType(curType, converter)
}

func createTupleType(curType *abi.Type, converter func(reflect.Type) reflect.Type) (reflect.Type, reflect.Type, error) {
	if len(curType.TupleElems) == 0 {
		return nil, nil, types.InvalidTypeError{}
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
