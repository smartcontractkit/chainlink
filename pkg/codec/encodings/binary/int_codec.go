package binary

import (
	"fmt"
	"reflect"

	codec2 "github.com/smartcontractkit/chainlink-common/pkg/codec/encodings"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type intCodec struct {
	codec   codec2.TypeCodec
	toInt   func(any) int
	fromInt func(int) any
}

func (i *intCodec) Size(numItems int) (int, error) {
	return i.codec.Size(numItems)
}

func (i *intCodec) FixedSize() (int, error) {
	return i.codec.FixedSize()
}

var _ codec2.TypeCodec = &intCodec{}

func (i *intCodec) Encode(value any, into []byte) ([]byte, error) {
	ival, ok := value.(int)
	if !ok {
		return nil, fmt.Errorf("%w: %T is not an int", types.ErrInvalidType, value)
	}

	return i.codec.Encode(i.fromInt(ival), into)
}

func (i *intCodec) Decode(encoded []byte) (any, []byte, error) {
	value, remaining, err := i.codec.Decode(encoded)
	if err != nil {
		return nil, nil, err
	}

	return i.toInt(value), remaining, nil
}

func (i *intCodec) GetType() reflect.Type {
	return reflect.TypeOf(0)
}

type uintCodec struct {
	codec    codec2.TypeCodec
	toUint   func(any) uint
	fromUint func(uint) any
}

var _ codec2.TypeCodec = &uintCodec{}

func (i *uintCodec) Encode(value any, uinto []byte) ([]byte, error) {
	ival, ok := value.(uint)
	if !ok {
		return nil, fmt.Errorf("%w: %T is not an uint", types.ErrInvalidType, value)
	}

	return i.codec.Encode(i.fromUint(ival), uinto)
}

func (i *uintCodec) Decode(encoded []byte) (any, []byte, error) {
	value, remaining, err := i.codec.Decode(encoded)
	if err != nil {
		return nil, nil, err
	}

	return i.toUint(value), remaining, nil
}

func (i *uintCodec) GetType() reflect.Type {
	return reflect.TypeOf(uint(0))
}

func (i *uintCodec) Size(numItems int) (int, error) {
	return i.codec.Size(numItems)
}

func (i *uintCodec) FixedSize() (int, error) {
	return i.codec.FixedSize()
}
