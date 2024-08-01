package codec

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

// NewHardCoder creates a modifier that will hard-code values for on-chain and off-chain types
// The modifier will override any values of the same name, if you need an overwritten value to be used in a different field,
// NewRenamer must be used before NewHardCoder.
func NewHardCoder(onChain map[string]any, offChain map[string]any, hooks ...mapstructure.DecodeHookFunc) (Modifier, error) {
	if err := verifyHardCodeKeys(onChain); err != nil {
		return nil, err
	} else if err = verifyHardCodeKeys(offChain); err != nil {
		return nil, err
	}

	myHooks := make([]mapstructure.DecodeHookFunc, len(hooks)+1)
	copy(myHooks, hooks)
	myHooks[len(hooks)] = hardCodeManyHook

	m := &onChainHardCoder{
		modifierBase: modifierBase[any]{
			fields:           offChain,
			onToOffChainType: map[reflect.Type]reflect.Type{},
			offToOnChainType: map[reflect.Type]reflect.Type{},
		},
		onChain: onChain,
		hooks:   myHooks,
	}
	m.modifyFieldForInput = func(_ string, field *reflect.StructField, key string, v any) error {
		// if we are typing it differently, we need to make sure it's hard-coded the other way
		newType := reflect.TypeOf(v)
		if _, ok := m.onChain[key]; !ok && field.Type != newType {
			return fmt.Errorf(
				"%w: cannot change field type without hard-coding its onchain value for key %s",
				types.ErrInvalidType,
				key)
		}
		field.Type = newType
		return nil
	}
	m.addFieldForInput = func(_, key string, value any) reflect.StructField {
		return reflect.StructField{
			Name: key,
			Type: reflect.TypeOf(value),
		}
	}
	return m, nil
}

type onChainHardCoder struct {
	modifierBase[any]
	onChain map[string]any
	hooks   []mapstructure.DecodeHookFunc
}

// verifyHardCodeKeys checks that no key is a prefix of another key
// This is important because if you hard code "A" : {"B" : 10}, and "A.C" : 20
// A key will override all A values and the A.C will add to existing values, which is inconsistent.
// instead the user should do "A" : {"B" : 10, "C" : 20} if they want to override or
// "A.B" : 10, "A.C" : 20 if they want to add
func verifyHardCodeKeys(values map[string]any) error {
	seen := map[string]bool{}
	for _, k := range subkeysLast(values) {
		parts := strings.Split(k, ".")
		on := ""
		for _, part := range parts {
			on += part
			if seen[on] {
				return fmt.Errorf("%w: key %s and %s cannot both be present", types.ErrInvalidConfig, on, k)
			}
		}
		seen[k] = true
	}
	return nil
}

func (o *onChainHardCoder) TransformToOnChain(offChainValue any, _ string) (any, error) {
	return transformWithMaps(offChainValue, o.offToOnChainType, o.onChain, hardCode, o.hooks...)
}

func (o *onChainHardCoder) TransformToOffChain(onChainValue any, _ string) (any, error) {
	allHooks := make([]mapstructure.DecodeHookFunc, len(o.hooks)+1)
	copy(allHooks, o.hooks)
	allHooks[len(o.hooks)] = hardCodeManyHook
	return transformWithMaps(onChainValue, o.onToOffChainType, o.fields, hardCode, allHooks...)
}

func hardCode(extractMap map[string]any, key string, item any) error {
	extractMap[key] = item
	return nil
}

// hardCodeManyHook allows a user to specify a single value for a slice or array
// This is useful because users may not know how many values are in an array ahead of time (e.g. number of reports)
// Instead, a user can specify A.C = 10 and if A is an array, all A.C values will be set to 10
func hardCodeManyHook(from reflect.Value, to reflect.Value) (any, error) {
	// A slice or array could be behind pointers. mapstructure could add an extra pointer level too.
	for to.Kind() == reflect.Pointer {
		to = to.Elem()
	}

	for from.Kind() == reflect.Pointer {
		from = from.Elem()
	}

	switch to.Kind() {
	case reflect.Slice, reflect.Array:
		switch from.Kind() {
		case reflect.Slice, reflect.Array:
			return from.Interface(), nil
		default:
		}
	default:
		return from.Interface(), nil
	}

	length := to.Len()
	array := reflect.MakeSlice(reflect.SliceOf(from.Type()), length, length)
	for i := 0; i < length; i++ {
		array.Index(i).Set(from)
	}
	return array.Interface(), nil
}
