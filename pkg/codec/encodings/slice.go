package encodings

import (
	"fmt"
	"reflect"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

func NewSlice(field, size TypeCodec) (TypeCodec, error) {
	if field == nil || size == nil {
		return nil, fmt.Errorf("%w: field and size must be non-nil", types.ErrInvalidConfig)
	}

	if size.GetType() != reflect.TypeOf(0) {
		return nil, fmt.Errorf("%w: size must be an int", types.ErrInvalidConfig)
	}

	return &slice{
		Field:     field,
		SizeCodec: size,
	}, nil
}

type slice struct {
	Field     TypeCodec
	SizeCodec TypeCodec
}

var _ TypeCodec = &slice{}

func (s *slice) Encode(value any, into []byte) ([]byte, error) {
	rValue := reflect.ValueOf(value)
	if rValue.Kind() != reflect.Array && rValue.Kind() != reflect.Slice {
		return nil, types.ErrNotASlice
	}

	fs, err := s.SizeCodec.FixedSize()
	if err != nil {
		return nil, err
	}

	numElements := rValue.Len()
	if numElements > 1<<(fs*8) {
		return nil, fmt.Errorf("%w: %v is too big to encode into a %v-bytes slice", types.ErrSliceWrongLen, numElements, fs)
	}

	toEncode := reflect.ValueOf(rValue.Len()).Convert(s.SizeCodec.GetType()).Interface()
	into, err = s.SizeCodec.Encode(toEncode, into)
	if err != nil {
		return nil, err
	}

	return EncodeEach(rValue, into, s.Field)
}

func (s *slice) Decode(encoded []byte) (any, []byte, error) {
	size, remaining, err := s.SizeCodec.Decode(encoded)
	if err != nil {
		return nil, nil, err
	}

	intSize, ok := size.(int)
	if !ok {
		return nil, nil, fmt.Errorf("%w: %T is not an int indicating the size", types.ErrInternal, size)
	}

	if intSize < 0 {
		return nil, nil, fmt.Errorf("%w: negative slice size %v", types.ErrInvalidEncoding, size)
	}

	rSlice := reflect.MakeSlice(s.GetType(), intSize, intSize)
	return DecodeEach(remaining, rSlice, intSize, s.Field)
}

func (s *slice) GetType() reflect.Type {
	return reflect.SliceOf(s.Field.GetType())
}

func (s *slice) Size(numItems int) (int, error) {
	sizeSize, err := s.SizeCodec.FixedSize()
	if err != nil {
		return 0, err
	}

	elemSize, err := s.Field.FixedSize()
	if err != nil {
		return 0, err
	}

	return sizeSize + elemSize*numItems, nil
}

func (s *slice) FixedSize() (int, error) {
	return 0, fmt.Errorf("%w: slices are not fixed size", types.ErrInvalidType)
}
