package codec

import (
	"fmt"
	"reflect"
	"unicode"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

func NewRenamer(fields map[string]string) Modifier {
	m := &renamer{
		modifierBase: modifierBase[string]{
			fields:           fields,
			onToOffChainType: map[reflect.Type]reflect.Type{},
			offToOnChainType: map[reflect.Type]reflect.Type{},
		},
	}
	m.modifyFieldForInput = func(pkgPath string, field *reflect.StructField, _, newName string) error {
		field.Name = newName
		if unicode.IsLower(rune(field.Name[0])) {
			field.PkgPath = pkgPath
		}
		return nil
	}
	return m
}

type renamer struct {
	modifierBase[string]
}

func (r *renamer) TransformToOffChain(onChainValue any, _ string) (any, error) {
	rOutput, err := renameTransform(r.onToOffChainType, reflect.ValueOf(onChainValue))
	if err != nil {
		return nil, err
	}
	return rOutput.Interface(), nil
}

func (r *renamer) TransformToOnChain(offChainValue any, _ string) (any, error) {
	rOutput, err := renameTransform(r.offToOnChainType, reflect.ValueOf(offChainValue))
	if err != nil {
		return nil, err
	}
	return rOutput.Interface(), nil
}

func renameTransform(typeMap map[reflect.Type]reflect.Type, rInput reflect.Value) (reflect.Value, error) {
	toType, ok := typeMap[rInput.Type()]
	if !ok {
		return reflect.Value{}, fmt.Errorf("%w: cannot rename unknown type %v", types.ErrInvalidType, toType)
	}

	if toType == rInput.Type() {
		return rInput, nil
	}

	switch rInput.Kind() {
	case reflect.Pointer:
		return reflect.NewAt(toType.Elem(), rInput.UnsafePointer()), nil
	case reflect.Struct, reflect.Slice, reflect.Array:
		return transformNonPointer(toType, rInput)
	default:
		return reflect.Value{}, fmt.Errorf("%w: cannot rename kind %v", types.ErrInvalidType, rInput.Kind())
	}
}

func transformNonPointer(toType reflect.Type, rInput reflect.Value) (reflect.Value, error) {
	// make sure the input is addressable
	ptr := reflect.New(rInput.Type())
	reflect.Indirect(ptr).Set(rInput)
	changed := reflect.NewAt(toType, ptr.UnsafePointer()).Elem()
	return changed, nil
}
