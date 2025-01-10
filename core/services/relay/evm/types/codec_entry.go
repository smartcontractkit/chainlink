package types

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
)

// MaxTopicFields is three because the EVM has a max of four topics, but the first topic is always the event signature.
const MaxTopicFields = 3

type CodecEntry interface {
	Init() error
	Args() abi.Arguments
	EncodingPrefix() []byte
	GetMaxSize(n int) (int, error)
	Modifier() codec.Modifier

	// CheckedType provides a type that can be used to decode into with type-safety around sizes of integers etc.
	CheckedType() reflect.Type

	// ToNative converts a pointer to checked value into a pointer of a type to use with the go-ethereum ABI encoder
	// Note that modification of the returned value will modify the original checked value and vice versa.
	ToNative(checked reflect.Value) (reflect.Value, error)

	// IsNativePointer returns if the type is a pointer to the native type
	IsNativePointer(item reflect.Type) bool
}

func NewCodecEntry(args abi.Arguments, encodingPrefix []byte, mod codec.Modifier) CodecEntry {
	if mod == nil {
		mod = codec.MultiModifier{}
	}
	return &codecEntry{args: args, encodingPrefix: encodingPrefix, mod: mod}
}

type codecEntry struct {
	args           abi.Arguments
	encodingPrefix []byte
	checkedType    reflect.Type
	nativeType     reflect.Type
	mod            codec.Modifier
}

func (entry *codecEntry) CheckedType() reflect.Type {
	return entry.checkedType
}

func (entry *codecEntry) NativeType() reflect.Type {
	return entry.nativeType
}

func (entry *codecEntry) ToNative(checked reflect.Value) (val reflect.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			val = reflect.Value{}
			err = fmt.Errorf("invalid checked value: %v", r)
		}
	}()

	// some checked types are expected to be pointers already for e.g. big numbers, so this is fine
	checkedTypeIsPtr := entry.checkedType == checked.Type()
	if checked.Type() != reflect.PointerTo(entry.checkedType) && !checkedTypeIsPtr {
		return reflect.Value{}, fmt.Errorf("%w: checked type %v does not match expected type %v", commontypes.ErrInvalidType, checked.Type(), entry.checkedType)
	}

	if checkedTypeIsPtr {
		return reflect.NewAt(entry.nativeType.Elem(), checked.UnsafePointer()), nil
	}

	return reflect.Indirect(reflect.NewAt(entry.nativeType, checked.UnsafePointer())), nil
}

func (entry *codecEntry) IsNativePointer(item reflect.Type) bool {
	return item == reflect.PointerTo(entry.nativeType)
}

func (entry *codecEntry) Modifier() codec.Modifier {
	return entry.mod
}

func (entry *codecEntry) Args() abi.Arguments {
	tmp := make(abi.Arguments, len(entry.args))
	copy(tmp, entry.args)
	return tmp
}

func (entry *codecEntry) EncodingPrefix() []byte {
	tmp := make([]byte, len(entry.encodingPrefix))
	copy(tmp, entry.encodingPrefix)
	return tmp
}

func (entry *codecEntry) Init() (err error) {
	// Since reflection panics if errors occur, best to recover in case of any unknown errors
	defer func() {
		if r := recover(); r != nil {
			entry.checkedType = nil
			entry.nativeType = nil
			err = fmt.Errorf("%w: %v", commontypes.ErrInvalidConfig, r)
		}
	}()
	if entry.checkedType != nil {
		return nil
	}

	args := unwrapArgs(entry.args)
	argLen := len(args)
	native := make([]reflect.StructField, argLen)
	checked := make([]reflect.StructField, argLen)

	// Single returns that aren't named will return that type
	// whereas named parameters will return a struct with the fields
	// Eg: function foo() returns (int256) ... will return a *big.Int for the native type
	// function foo() returns (int256 i) ... will return a struct { I *big.Int } for the native type
	// function foo() returns (int256 i1, int256 i2) ... will return a struct { I1 *big.Int, I2 *big.Int } for the native type
	if len(args) == 1 && args[0].Name == "" {
		nativeArg, checkedArg, err := getNativeAndCheckedTypesForArg(&args[0])
		if err != nil {
			return err
		}
		entry.nativeType = nativeArg
		entry.checkedType = checkedArg
		return nil
	}

	numIndices := 0
	seenNames := map[string]bool{}
	for i, arg := range args {
		if arg.Indexed {
			if numIndices == MaxTopicFields {
				return fmt.Errorf("%w: too many indexed arguments", commontypes.ErrInvalidConfig)
			}
			numIndices++
		}

		tmp := arg
		nativeArg, checkedArg, err := getNativeAndCheckedTypesForArg(&tmp)
		if err != nil {
			return err
		}
		allowRename := false
		if len(arg.Name) == 0 {
			arg.Name = fmt.Sprintf("F%d", i)
			allowRename = true
		}

		name := strings.ToUpper(arg.Name[:1]) + arg.Name[1:]
		if seenNames[name] {
			if !allowRename {
				return fmt.Errorf("%w: duplicate field name %s, after ToCamelCase", commontypes.ErrInvalidConfig, name)
			}
			for {
				name = name + "_X"
				arg.Name = name
				if !seenNames[name] {
					break
				}
			}
		}
		args[i] = arg
		seenNames[name] = true
		native[i] = reflect.StructField{Name: name, Type: nativeArg}
		checked[i] = reflect.StructField{Name: name, Type: checkedArg}
	}

	entry.nativeType = structOfPointers(native)
	entry.checkedType = structOfPointers(checked)
	return nil
}

func (entry *codecEntry) GetMaxSize(n int) (int, error) {
	return GetMaxSize(n, entry.args)
}

func unwrapArgs(args abi.Arguments) abi.Arguments {
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

func getNativeAndCheckedTypesForArg(arg *abi.Argument) (reflect.Type, reflect.Type, error) {
	tmp := arg.Type
	if arg.Indexed {
		switch arg.Type.T {
		case abi.StringTy:
			return reflect.TypeOf(common.Hash{}), reflect.TypeOf(common.Hash{}), nil
		case abi.ArrayTy:
			u8, _ := GetAbiEncodingType("uint8")
			if arg.Type.Elem.GetType() == u8.native {
				return reflect.TypeOf(common.Hash{}), reflect.TypeOf(common.Hash{}), nil
			}
			fallthrough
		case abi.SliceTy, abi.TupleTy, abi.FixedPointTy, abi.FunctionTy:
			// https://github.com/ethereum/go-ethereum/blob/release/1.12/accounts/abi/topics.go#L78
			return nil, nil, fmt.Errorf("%w: unsupported indexed type: %v", commontypes.ErrInvalidConfig, arg.Type)
		default:
		}
	}

	return getNativeAndCheckedTypes(&tmp)
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
			return nil, nil, fmt.Errorf(
				"%w: cannot create type for kind %v", commontypes.ErrInvalidType, curType.GetType().Kind())
		}
	}
	base, ok := GetAbiEncodingType(curType.String())
	if ok {
		return converter(base.native), converter(base.checked), nil
	}

	return createTupleType(curType, converter)
}

func createTupleType(curType *abi.Type, converter func(reflect.Type) reflect.Type) (reflect.Type, reflect.Type, error) {
	if len(curType.TupleElems) == 0 {
		if curType.TupleType == nil {
			return nil, nil, fmt.Errorf("%w: unsupported solidity type: %v", commontypes.ErrInvalidType, curType.String())
		}
		return curType.TupleType, curType.TupleType, nil
	}

	// Our naive types always have the same layout as the checked ones.
	// This differs intentionally from the type.GetType() in abi as fields on structs are pointers in ours to
	// verify that fields are intentionally set.
	nativeFields := make([]reflect.StructField, len(curType.TupleElems))
	checkedFields := make([]reflect.StructField, len(curType.TupleElems))
	for i, elm := range curType.TupleElems {
		name := curType.TupleRawNames[i]
		name = strings.ToUpper(name[:1]) + name[1:]
		nativeFields[i].Name = name
		checkedFields[i].Name = name
		nativeArgType, checkedArgType, err := getNativeAndCheckedTypes(elm)
		if err != nil {
			return nil, nil, err
		}
		nativeFields[i].Type = nativeArgType
		checkedFields[i].Type = checkedArgType
	}
	return converter(structOfPointers(nativeFields)), converter(structOfPointers(checkedFields)), nil
}

func structOfPointers(fields []reflect.StructField) reflect.Type {
	for i := range fields {
		if fields[i].Type.Kind() != reflect.Pointer {
			fields[i].Type = reflect.PointerTo(fields[i].Type)
		}
	}
	return reflect.StructOf(fields)
}
