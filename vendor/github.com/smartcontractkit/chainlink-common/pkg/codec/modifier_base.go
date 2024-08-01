package codec

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/mitchellh/mapstructure"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type modifierBase[T any] struct {
	fields              map[string]T
	onToOffChainType    map[reflect.Type]reflect.Type
	offToOnChainType    map[reflect.Type]reflect.Type
	modifyFieldForInput func(pkgPath string, outputField *reflect.StructField, fullPath string, change T) error
	addFieldForInput    func(pkgPath, name string, change T) reflect.StructField
}

func (m *modifierBase[T]) RetypeToOffChain(onChainType reflect.Type, itemType string) (tpe reflect.Type, err error) {
	defer func() {
		// StructOf can panic if the fields are not valid
		if r := recover(); r != nil {
			tpe = nil
			err = fmt.Errorf("%w: %v", types.ErrInvalidType, r)
		}
	}()
	if m.fields == nil || len(m.fields) == 0 {
		m.offToOnChainType[onChainType] = onChainType
		m.onToOffChainType[onChainType] = onChainType
		return onChainType, nil
	}

	if cached, ok := m.onToOffChainType[onChainType]; ok {
		return cached, nil
	}

	switch onChainType.Kind() {
	case reflect.Pointer:
		elm, err := m.RetypeToOffChain(onChainType.Elem(), "")
		if err != nil {
			return nil, err
		}

		ptr := reflect.PointerTo(elm)
		m.onToOffChainType[onChainType] = ptr
		m.offToOnChainType[ptr] = onChainType
		return ptr, nil
	case reflect.Slice:
		elm, err := m.RetypeToOffChain(onChainType.Elem(), "")
		if err != nil {
			return nil, err
		}

		sliceType := reflect.SliceOf(elm)
		m.onToOffChainType[onChainType] = sliceType
		m.offToOnChainType[sliceType] = onChainType
		return sliceType, nil
	case reflect.Array:
		elm, err := m.RetypeToOffChain(onChainType.Elem(), "")
		if err != nil {
			return nil, err
		}

		arrayType := reflect.ArrayOf(onChainType.Len(), elm)
		m.onToOffChainType[onChainType] = arrayType
		m.offToOnChainType[arrayType] = onChainType
		return arrayType, nil
	case reflect.Struct:
		return m.getStructType(onChainType)
	default:
		return nil, fmt.Errorf("%w: cannot retype the kind %v", types.ErrInvalidType, onChainType.Kind())
	}
}

func (m *modifierBase[T]) getStructType(outputType reflect.Type) (reflect.Type, error) {
	filedLocations, err := getFieldIndices(outputType)
	if err != nil {
		return nil, err
	}

	for _, key := range m.subkeysFirst() {
		parts := strings.Split(key, ".")
		fieldName := parts[len(parts)-1]
		parts = parts[:len(parts)-1]
		curLocations := filedLocations
		for _, part := range parts {
			if curLocations, err = curLocations.populateSubFields(part); err != nil {
				return nil, err
			}
		}

		// If a subkey has been modified, update the underlying types first
		curLocations.updateTypeFromSubkeyMods(fieldName)
		if field, ok := curLocations.fieldByName(fieldName); ok {
			if err = m.modifyFieldForInput(curLocations.pkgPath, field, key, m.fields[key]); err != nil {
				return nil, err
			}
		} else {
			if m.addFieldForInput == nil {
				return nil, fmt.Errorf("%w: cannot find %s", types.ErrInvalidType, key)
			}
			curLocations.addNewField(m.addFieldForInput(curLocations.pkgPath, fieldName, m.fields[key]))
		}
	}

	newStruct := filedLocations.makeNewType()
	m.onToOffChainType[outputType] = newStruct
	m.offToOnChainType[newStruct] = outputType
	return newStruct, nil
}

// subkeysFirst returns a list of keys that will always have a sub-key before the key if both are present
func (m *modifierBase[T]) subkeysFirst() []string {
	orderedKeys := make([]string, 0, len(m.fields))
	for k := range m.fields {
		orderedKeys = append(orderedKeys, k)
	}

	sort.Slice(orderedKeys, func(i, j int) bool {
		return orderedKeys[i] > orderedKeys[j]
	})

	return orderedKeys
}

// subkeysLast returns a list of keys that will always have a sub-key after the key if both are present
func subkeysLast[T any](fields map[string]T) []string {
	orderedKeys := make([]string, 0, len(fields))
	for k := range fields {
		orderedKeys = append(orderedKeys, k)
	}

	sort.Strings(orderedKeys)
	return orderedKeys
}

type mapAction[T any] func(extractMap map[string]any, key string, element T) error

func transformWithMaps[T any](
	item any,
	typeMap map[reflect.Type]reflect.Type,
	fields map[string]T,
	fn mapAction[T],
	hooks ...mapstructure.DecodeHookFunc) (any, error) {
	rItem := reflect.ValueOf(item)

	toType, ok := typeMap[rItem.Type()]
	if !ok {
		return reflect.Value{}, fmt.Errorf("%w: cannot retype %v", types.ErrInvalidType, rItem.Type())
	}

	rOutput, err := transformWithMapsHelper(rItem, toType, fields, fn, hooks)
	if err != nil {
		return reflect.Value{}, err
	}

	return rOutput.Interface(), nil
}

func transformWithMapsHelper[T any](
	rItem reflect.Value,
	toType reflect.Type,
	fields map[string]T,
	fn mapAction[T],
	hooks []mapstructure.DecodeHookFunc) (reflect.Value, error) {
	switch rItem.Kind() {
	case reflect.Pointer:
		elm := rItem.Elem()
		if elm.Kind() == reflect.Struct {
			into := reflect.New(toType.Elem())
			err := changeElements(rItem.Interface(), into.Interface(), fields, fn, hooks)

			return into, err
		}

		tmp, err := transformWithMapsHelper(elm, toType.Elem(), fields, fn, hooks)
		result := reflect.New(toType.Elem())
		reflect.Indirect(result).Set(tmp)

		return result, err
	case reflect.Struct:
		into := reflect.New(toType)
		err := changeElements(rItem.Interface(), into.Interface(), fields, fn, hooks)

		return into.Elem(), err
	case reflect.Slice:
		length := rItem.Len()
		into := reflect.MakeSlice(toType, length, length)
		err := doMany(rItem, into, fields, fn, hooks)

		return into, err
	case reflect.Array:
		into := reflect.New(toType).Elem()
		err := doMany(rItem, into, fields, fn, hooks)

		return into, err
	default:
		return reflect.Value{}, fmt.Errorf("%w: cannot retype %v", types.ErrInvalidType, rItem.Type())
	}
}

func doMany[T any](rInput, rOutput reflect.Value, fields map[string]T, fn mapAction[T], hooks []mapstructure.DecodeHookFunc) error {
	length := rInput.Len()
	for i := 0; i < length; i++ {
		// Make sure the items are addressable
		inTmp := rInput.Index(i)
		outTmp := rOutput.Index(i)

		output, err := transformWithMapsHelper(inTmp, outTmp.Type(), fields, fn, hooks)
		if err != nil {
			return err
		}

		outTmp.Set(output)
	}

	return nil
}

func changeElements[T any](src, dest any, fields map[string]T, fn mapAction[T], hooks []mapstructure.DecodeHookFunc) error {
	valueMapping := map[string]any{}
	if err := mapstructure.Decode(src, &valueMapping); err != nil {
		return fmt.Errorf("%w: failed to decode source type: %w", types.ErrInvalidType, err)
	}

	if err := doForMapElements(valueMapping, fields, fn); err != nil {
		return err
	}

	conf := &mapstructure.DecoderConfig{Result: &dest, Squash: true}
	if len(hooks) != 0 {
		conf.DecodeHook = mapstructure.ComposeDecodeHookFunc(hooks...)
	}

	hookedDecoder, err := mapstructure.NewDecoder(conf)
	if err != nil {
		return fmt.Errorf("%w: failed to create configured decoder: %w", types.ErrInvalidType, err)
	}

	if err = hookedDecoder.Decode(valueMapping); err != nil {
		return fmt.Errorf("%w: failed to decode destination type: %w", types.ErrInvalidType, err)
	}

	return nil
}

func doForMapElements[T any](valueMapping map[string]any, fields map[string]T, fn mapAction[T]) error {
	for key, value := range fields {
		path := strings.Split(key, ".")
		name := path[len(path)-1]
		path = path[:len(path)-1]

		extractMaps, err := getMapsFromPath(valueMapping, path)
		if err != nil {
			return PathMappingError{Err: err, Path: key}
		}

		for _, em := range extractMaps {
			if err = fn(em, name, value); err != nil {
				return PathMappingError{Err: err, Path: key}
			}
		}
	}

	return nil
}

type PathMappingError struct {
	Err  error
	Path string
}

func (e PathMappingError) Error() string {
	return fmt.Sprintf("mapping error for path (%s): %s", e.Path, e.Err)
}

func (e PathMappingError) Cause() error {
	return e.Err
}
