package codec

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

// NewPropertyExtractor creates a modifier that will extract a single property from a struct.
// This modifier is lossy, as TransformToOffchain will discard unwanted struct properties and
// return a single element. Calling TransformToOnchain will result in unset properties.
func NewPropertyExtractor(fieldName string) Modifier {
	m := &propertyExtractor{
		onToOffChainType: map[reflect.Type]reflect.Type{},
		offToOnChainType: map[reflect.Type]reflect.Type{},
		fieldName:        fieldName,
	}

	return m
}

type propertyExtractor struct {
	onToOffChainType map[reflect.Type]reflect.Type
	offToOnChainType map[reflect.Type]reflect.Type
	fieldName        string
}

func (e *propertyExtractor) RetypeToOffChain(onChainType reflect.Type, itemType string) (reflect.Type, error) {
	if e.fieldName == "" {
		return nil, fmt.Errorf("%w: field name required for extraction", types.ErrInvalidConfig)
	}

	if cached, ok := e.onToOffChainType[onChainType]; ok {
		return cached, nil
	}

	switch onChainType.Kind() {
	case reflect.Pointer:
		elm, err := e.RetypeToOffChain(onChainType.Elem(), "")
		if err != nil {
			return nil, err
		}

		ptr := reflect.PointerTo(elm)
		e.onToOffChainType[onChainType] = ptr
		e.offToOnChainType[ptr] = onChainType

		return ptr, nil
	case reflect.Slice:
		elm, err := e.RetypeToOffChain(onChainType.Elem(), "")
		if err != nil {
			return nil, err
		}

		sliceType := reflect.SliceOf(elm)
		e.onToOffChainType[onChainType] = sliceType
		e.offToOnChainType[sliceType] = onChainType

		return sliceType, nil
	case reflect.Array:
		elm, err := e.RetypeToOffChain(onChainType.Elem(), "")
		if err != nil {
			return nil, err
		}

		arrayType := reflect.ArrayOf(onChainType.Len(), elm)
		e.onToOffChainType[onChainType] = arrayType
		e.offToOnChainType[arrayType] = onChainType

		return arrayType, nil
	case reflect.Struct:
		return e.getPropTypeFromStruct(onChainType)
	default:
		return nil, fmt.Errorf("%w: cannot retype the kind %v", types.ErrInvalidType, onChainType.Kind())
	}
}

func (e *propertyExtractor) TransformToOnChain(offChainValue any, _ string) (any, error) {
	return extractOrExpandWithMaps(offChainValue, e.offToOnChainType, e.fieldName, expandWithMapsHelper)
}

func (e *propertyExtractor) TransformToOffChain(onChainValue any, _ string) (any, error) {
	return extractOrExpandWithMaps(onChainValue, e.onToOffChainType, e.fieldName, extractWithMapsHelper)
}

func (e *propertyExtractor) getPropTypeFromStruct(onChainType reflect.Type) (reflect.Type, error) {
	filedLocations, err := getFieldIndices(onChainType)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(e.fieldName, ".")
	fieldName := parts[len(parts)-1]
	parts = parts[:len(parts)-1]

	curLocations := filedLocations
	for _, part := range parts {
		if curLocations, err = curLocations.populateSubFields(part); err != nil {
			return nil, err
		}
	}

	curLocations.updateTypeFromSubkeyMods(fieldName)
	field, ok := curLocations.fieldByName(fieldName)
	if !ok {
		return nil, fmt.Errorf("%w: field not found in on-chain type %s", types.ErrInvalidType, e.fieldName)
	}

	e.onToOffChainType[onChainType] = field.Type
	e.offToOnChainType[field.Type] = onChainType

	return field.Type, nil
}

type transformHelperFunc func(reflect.Value, reflect.Type, string) (reflect.Value, error)

func extractOrExpandWithMaps(input any, typeMap map[reflect.Type]reflect.Type, field string, fn transformHelperFunc) (any, error) {
	rItem := reflect.ValueOf(input)

	toType, ok := typeMap[rItem.Type()]
	if !ok {
		return reflect.Value{}, fmt.Errorf("%w: cannot retype %v", types.ErrInvalidType, rItem.Type())
	}

	output, err := fn(rItem, toType, field)
	if err != nil {
		return reflect.Value{}, err
	}

	return output.Interface(), err
}

func expandWithMapsHelper(rItem reflect.Value, toType reflect.Type, field string) (reflect.Value, error) {
	switch toType.Kind() {
	case reflect.Pointer:
		if rItem.Kind() != reflect.Pointer {
			return reflect.Value{}, fmt.Errorf("%w: value to expand should be pointer", types.ErrInvalidType)
		}

		if toType.Elem().Kind() == reflect.Struct {
			into := reflect.New(toType.Elem())
			err := setFieldValue(rItem.Elem().Interface(), into.Interface(), field)
			return into, err
		}

		tmp, err := expandWithMapsHelper(rItem.Elem(), toType.Elem(), field)
		result := reflect.New(toType.Elem())
		reflect.Indirect(result).Set(tmp)

		return result, err
	case reflect.Struct:
		into := reflect.New(toType)
		err := setFieldValue(rItem.Interface(), into.Interface(), field)

		return into.Elem(), err
	case reflect.Slice:
		if rItem.Kind() != reflect.Slice {
			return reflect.Value{}, fmt.Errorf("%w: value to expand should be slice", types.ErrInvalidType)
		}

		length := rItem.Len()
		into := reflect.MakeSlice(toType, length, length)
		err := extractOrExpandMany(rItem, into, field, expandWithMapsHelper)

		return into, err
	case reflect.Array:
		if rItem.Kind() != reflect.Array {
			return reflect.Value{}, fmt.Errorf("%w: value to expand should be array", types.ErrInvalidType)
		}

		into := reflect.New(toType).Elem()
		err := extractOrExpandMany(rItem, into, field, expandWithMapsHelper)

		return into, err
	default:
		return reflect.Value{}, fmt.Errorf("%w: cannot retype", types.ErrInvalidType)
	}
}

func extractWithMapsHelper(rItem reflect.Value, toType reflect.Type, field string) (reflect.Value, error) {
	switch rItem.Kind() {
	case reflect.Pointer:
		elm := rItem.Elem()
		if elm.Kind() == reflect.Struct {
			tmp, err := extractElement(rItem.Interface(), field)
			result := reflect.New(toType.Elem())
			reflect.Indirect(result).Set(tmp)

			return result, err
		}

		tmp, err := extractWithMapsHelper(elm, toType.Elem(), field)
		result := reflect.New(toType.Elem())
		reflect.Indirect(result).Set(tmp)

		return result, err
	case reflect.Struct:
		return extractElement(rItem.Interface(), field)
	case reflect.Slice:
		length := rItem.Len()
		into := reflect.MakeSlice(toType, length, length)
		err := extractOrExpandMany(rItem, into, field, extractWithMapsHelper)
		return into, err
	case reflect.Array:
		into := reflect.New(toType).Elem()
		err := extractOrExpandMany(rItem, into, field, extractWithMapsHelper)
		return into, err
	default:
		return reflect.Value{}, fmt.Errorf("%w: cannot retype %v", types.ErrInvalidType, rItem.Type())
	}
}

type extractOrExpandHelperFunc func(reflect.Value, reflect.Type, string) (reflect.Value, error)

func extractOrExpandMany(rInput, rOutput reflect.Value, field string, fn extractOrExpandHelperFunc) error {
	length := rInput.Len()

	for i := 0; i < length; i++ {
		inTmp := rInput.Index(i)
		outTmp := rOutput.Index(i)

		output, err := fn(inTmp, outTmp.Type(), field)
		if err != nil {
			return err
		}

		outTmp.Set(output)
	}

	return nil
}

func extractElement(src any, field string) (reflect.Value, error) {
	valueMapping := map[string]any{}
	if err := mapstructure.Decode(src, &valueMapping); err != nil {
		return reflect.Value{}, err
	}

	path, name := pathAndName(field)

	extractMaps, err := getMapsFromPath(valueMapping, path)
	if err != nil {
		return reflect.Value{}, err
	}

	if len(extractMaps) != 1 {
		return reflect.Value{}, fmt.Errorf("%w: cannot find %s", types.ErrInvalidType, field)
	}

	em := extractMaps[0]

	item, ok := em[name]
	if !ok {
		return reflect.Value{}, fmt.Errorf("%w: cannot find %s", types.ErrInvalidType, field)
	}

	return reflect.ValueOf(item), nil
}

func setFieldValue(src, dest any, field string) error {
	valueMapping := map[string]any{}
	if err := mapstructure.Decode(dest, &valueMapping); err != nil {
		return fmt.Errorf("%w: %w", types.ErrInvalidType, err)
	}

	path, name := pathAndName(field)

	extractMaps, err := getMapsFromPath(valueMapping, path)
	if err != nil {
		return err
	}

	if len(extractMaps) != 1 {
		return fmt.Errorf("%w: only 1 extract map expected", types.ErrInvalidType)
	}

	extractMaps[0][name] = src

	conf := &mapstructure.DecoderConfig{Result: &dest, Squash: true}
	decoder, err := mapstructure.NewDecoder(conf)
	if err != nil {
		return fmt.Errorf("%w: %w", types.ErrInvalidType, err)
	}

	if err = decoder.Decode(valueMapping); err != nil {
		return fmt.Errorf("%w: %w", types.ErrInvalidType, err)
	}

	return nil
}

func pathAndName(field string) ([]string, string) {
	path := strings.Split(field, ".")
	name := path[len(path)-1]
	path = path[:len(path)-1]

	return path, name
}
