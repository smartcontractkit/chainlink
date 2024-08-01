package encodings

import (
	"fmt"
	"reflect"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type NamedTypeCodec struct {
	Name  string
	Codec TypeCodec
}

// NewStructCodec creates a codec that encodes fields with the given names and codecs in-order.
// Note: To verify fields are not defaulted,
// Codecs with non-pointer types in fields will be wrapped with encodings.NotNilPointer
func NewStructCodec(fields []NamedTypeCodec) (c TopLevelCodec, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%w: %v", types.ErrInvalidConfig, r)
		}
	}()

	sfs := make([]reflect.StructField, len(fields))
	codecFields := make([]TypeCodec, len(fields))
	for i, field := range fields {
		ft := field.Codec.GetType()
		if ft.Kind() != reflect.Pointer {
			field.Codec = &NotNilPointer{Elm: field.Codec}
			ft = reflect.PointerTo(ft)
		}

		sfs[i] = reflect.StructField{
			Name: field.Name,
			Type: ft,
		}
		codecFields[i] = field.Codec
	}

	return &structCodec{
		fields: codecFields,
		tpe:    reflect.PointerTo(reflect.StructOf(sfs)),
	}, nil
}

type structCodec struct {
	fields []TypeCodec
	tpe    reflect.Type
}

func (s *structCodec) Encode(value any, into []byte) ([]byte, error) {
	rVal := reflect.ValueOf(value)
	if rVal.Type() != s.tpe {
		return nil, fmt.Errorf("%w: expected %v, got %T", types.ErrInvalidType, s.tpe, value)
	}

	rVal = reflect.Indirect(rVal)

	for i, field := range s.fields {
		var err error
		if into, err = field.Encode(rVal.Field(i).Interface(), into); err != nil {
			return nil, err
		}
	}

	return into, nil
}

func (s *structCodec) Decode(encoded []byte) (any, []byte, error) {
	rVal := reflect.New(s.tpe.Elem())
	iVal := reflect.Indirect(rVal)
	for i, field := range s.fields {
		var fieldValue any
		var err error
		if fieldValue, encoded, err = field.Decode(encoded); err != nil {
			return nil, nil, err
		}
		iVal.Field(i).Set(reflect.ValueOf(fieldValue))
	}

	return rVal.Interface(), encoded, nil
}

func (s *structCodec) GetType() reflect.Type {
	return s.tpe
}

func (s *structCodec) Size(_ int) (int, error) {
	return s.FixedSize()
}

func (s *structCodec) FixedSize() (int, error) {
	size := 0
	for _, field := range s.fields {
		fieldSize, err := field.FixedSize()
		if err != nil {
			return 0, err
		}
		size += fieldSize
	}
	return size, nil
}

func (s *structCodec) SizeAtTopLevel(numItems int) (int, error) {
	size := 0
	for _, field := range s.fields {
		fieldSize, err := field.Size(numItems)
		if err != nil {
			return 0, err
		}
		size += fieldSize
	}
	return size, nil
}
