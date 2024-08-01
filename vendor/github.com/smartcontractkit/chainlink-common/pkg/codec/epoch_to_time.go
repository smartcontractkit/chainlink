package codec

import (
	"fmt"
	"math/big"
	"reflect"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

// NewEpochToTimeModifier converts all fields from time.Time off-chain to int64.
func NewEpochToTimeModifier(fields []string) Modifier {
	fieldMap := map[string]bool{}
	for _, field := range fields {
		fieldMap[field] = true
	}

	m := &timeToUnixModifier{
		modifierBase: modifierBase[bool]{
			fields:           fieldMap,
			onToOffChainType: map[reflect.Type]reflect.Type{},
			offToOnChainType: map[reflect.Type]reflect.Type{},
		},
	}

	m.modifyFieldForInput = func(_ string, field *reflect.StructField, _ string, _ bool) error {
		t, err := convertInt64InTypeToTime(field.Type, field.Name)
		if err != nil {
			return err
		}
		field.Type = t
		return nil
	}

	return m
}

type timeToUnixModifier struct {
	modifierBase[bool]
}

func (t *timeToUnixModifier) TransformToOnChain(offChainValue any, itemType string) (any, error) {
	// since the hook will convert time.Time to epoch, we don't need to worry about converting them in the maps
	return transformWithMaps(offChainValue, t.offToOnChainType, t.fields, noop, EpochToTimeHook, BigIntHook)
}

func (t *timeToUnixModifier) TransformToOffChain(onChainValue any, itemType string) (any, error) {
	// since the hook will convert epoch to time.Time, we don't need to worry about converting them in the maps
	return transformWithMaps(onChainValue, t.onToOffChainType, t.fields, noop, EpochToTimeHook, BigIntHook)
}

func noop(_ map[string]any, _ string, _ bool) error {
	return nil
}

func convertInt64InTypeToTime(t reflect.Type, field string) (reflect.Type, error) {
	converter := func(t reflect.Type) reflect.Type { return t }
	for {
		if t.ConvertibleTo(i64Type) {
			return converter(reflect.TypeOf(&time.Time{})), nil
		}

		switch t.Kind() {
		case reflect.Pointer:
			if t.ConvertibleTo(reflect.TypeOf(&big.Int{})) {
				return converter(reflect.TypeOf(&time.Time{})), nil
			}
		case reflect.Slice, reflect.Array:
			tmp := converter
			// works for array decoding too, [SliceToArrayVerifySizeHook]
			// is used if the on-chain type requires array size checking
			converter = func(t reflect.Type) reflect.Type { return reflect.SliceOf(tmp(t)) }
		default:
			return nil, fmt.Errorf("%w: cannot convert time for field %s", types.ErrInvalidType, field)
		}
		t = t.Elem()
	}
}
