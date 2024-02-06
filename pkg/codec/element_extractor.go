package codec

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

// ElementExtractorLocation is used to determine which element to extract from a slice or array.
// The default is ElementExtractorLocationMiddle, which will extract the middle element.
// valid json values are "first", "middle", and "last".
type ElementExtractorLocation int

const (
	ElementExtractorLocationFirst ElementExtractorLocation = iota
	ElementExtractorLocationMiddle
	ElementExtractorLocationLast
	ElementExtractorLocationDefault = ElementExtractorLocationMiddle
)

func (e ElementExtractorLocation) MarshalJSON() ([]byte, error) {
	switch e {
	case ElementExtractorLocationFirst:
		return json.Marshal("first")
	case ElementExtractorLocationMiddle:
		return json.Marshal("middle")
	case ElementExtractorLocationLast:
		return json.Marshal("last")
	default:
		return nil, fmt.Errorf("%w: %d", types.ErrInvalidType, e)
	}
}

func (e *ElementExtractorLocation) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	switch strings.ToLower(s) {
	case "first":
		*e = ElementExtractorLocationFirst
	case "middle":
		*e = ElementExtractorLocationMiddle
	case "last":
		*e = ElementExtractorLocationLast
	default:
		return fmt.Errorf("%w: %s", types.ErrInvalidType, s)
	}
	return nil
}

// NewElementExtractor creates a modifier that extracts an element from a slice or array.
// fields is used to determine which fields to extract elements from and which element to extract.
// This modifier is lossy, as TransformToOffChain will always return a slice of length 1 with the single element,
// so calling TransformToOnChain, then TransformToOffChain will not return the original value, if it has multiple elements.
func NewElementExtractor(fields map[string]*ElementExtractorLocation) Modifier {
	m := &elementExtractor{
		modifierBase: modifierBase[*ElementExtractorLocation]{
			fields:           fields,
			onToOffChainType: map[reflect.Type]reflect.Type{},
			offToOnChainType: map[reflect.Type]reflect.Type{},
		},
	}
	m.modifyFieldForInput = func(_ string, field *reflect.StructField, _ string, _ *ElementExtractorLocation) error {
		field.Type = reflect.SliceOf(field.Type)
		return nil
	}

	return m
}

type elementExtractor struct {
	modifierBase[*ElementExtractorLocation]
}

func (e *elementExtractor) TransformToOnChain(offChainValue any, _ string) (any, error) {
	return transformWithMaps(offChainValue, e.offToOnChainType, e.fields, extractMap)
}

func (e *elementExtractor) TransformToOffChain(onChainValue any, _ string) (any, error) {
	return transformWithMaps(onChainValue, e.onToOffChainType, e.fields, expandMap)
}

func extractMap(extractMap map[string]any, key string, elementLocation *ElementExtractorLocation) error {
	item, ok := extractMap[key]
	if !ok {
		return fmt.Errorf("%w: cannot find %s", types.ErrInvalidType, key)
	} else if item == nil {
		return nil
	}

	if elementLocation == nil {
		tmp := ElementExtractorLocationDefault
		elementLocation = &tmp
	}

	rItem := reflect.ValueOf(item)
	switch rItem.Kind() {
	case reflect.Array, reflect.Slice:
	default:
		return fmt.Errorf("%w: %s is not a slice or array", types.ErrInvalidType, key)
	}

	if rItem.Len() == 0 {
		extractMap[key] = nil
		return nil
	}

	switch *elementLocation {
	case ElementExtractorLocationFirst:
		extractMap[key] = rItem.Index(0).Interface()
	case ElementExtractorLocationMiddle:
		extractMap[key] = rItem.Index(rItem.Len() / 2).Interface()
	case ElementExtractorLocationLast:
		extractMap[key] = rItem.Index(rItem.Len() - 1).Interface()
	}

	return nil
}

func expandMap(extractMap map[string]any, key string, _ *ElementExtractorLocation) error {
	item, ok := extractMap[key]
	if !ok {
		return fmt.Errorf("%w: cannot find %s", types.ErrInvalidType, key)
	} else if item == nil {
		return nil
	}

	rItem := reflect.ValueOf(item)
	slice := reflect.MakeSlice(reflect.SliceOf(rItem.Type()), 1, 1)
	slice.Index(0).Set(rItem)
	extractMap[key] = slice.Interface()
	return nil
}
