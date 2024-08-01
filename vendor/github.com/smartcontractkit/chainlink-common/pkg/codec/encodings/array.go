package encodings

import (
	"fmt"
	"reflect"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

func NewArray(numElements int, underlying TypeCodec) (TypeCodec, error) {
	if underlying == nil {
		return nil, fmt.Errorf("%w: field type cannot be nil", types.ErrInvalidConfig)
	}

	return &array{
		NumElements: numElements,
		Field:       underlying,
	}, nil
}

type array struct {
	NumElements int
	Field       TypeCodec
}

var _ TypeCodec = &array{}

func (a *array) Encode(value any, into []byte) ([]byte, error) {
	rValue := reflect.ValueOf(value)
	if kind := rValue.Kind(); kind != reflect.Array && kind != reflect.Slice {
		return nil, fmt.Errorf("%w: expected array or slice but got %s", types.ErrNotASlice, kind)
	}

	if rValue.Len() != a.NumElements {
		return nil, fmt.Errorf("%w: expected %v elements, got %v", types.ErrSliceWrongLen, a.NumElements, rValue.Len())
	}

	return EncodeEach(rValue, into, a.Field)
}

func (a *array) Decode(encoded []byte) (any, []byte, error) {
	rArray := reflect.New(a.GetType()).Elem()
	return DecodeEach(encoded, rArray, a.NumElements, a.Field)
}

func (a *array) GetType() reflect.Type {
	return reflect.ArrayOf(a.NumElements, a.Field.GetType())
}

func (a *array) Size(_ int) (int, error) {
	return a.FixedSize()
}

func (a *array) FixedSize() (int, error) {
	fs, err := a.Field.FixedSize()
	if err != nil {
		return 0, err
	}
	return fs * a.NumElements, nil
}
