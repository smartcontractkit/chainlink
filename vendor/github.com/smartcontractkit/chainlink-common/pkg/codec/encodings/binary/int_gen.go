// DO NOT MODIFY: automatically generated from chainlink-common/pkg/codec/encodings/binary/gen/main.go using the template int_gen.go

package binary

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/smartcontractkit/chainlink-common/pkg/codec/encodings"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type Int8 struct{ encoder }

var _ encodings.TypeCodec = &Int8{}

func (i *Int8) Encode(value any, into []byte) ([]byte, error) {
	v, ok := value.(int8)
	if !ok {
		return nil, fmt.Errorf("%w: %T is not an int8", types.ErrInvalidType, value)
	}
	return append(into, byte(v)), nil
}

func (i *Int8) Decode(encoded []byte) (any, []byte, error) {
	ui, remaining, err := encodings.SafeDecode[uint8](encoded, 1, func(encoded []byte) byte { return encoded[0] })
	return int8(ui), remaining, err
}

func (*Int8) GetType() reflect.Type {
	return reflect.TypeOf(int8(0))
}

func (*Int8) Size(int) (int, error) {
	return 1, nil
}

func (*Int8) FixedSize() (int, error) {
	return 1, nil
}

func (e *endianEncoder) Int8() encodings.TypeCodec {
	return &Int8{encoder: e.encoder}
}

type Uint8 struct{ encoder }

var _ encodings.TypeCodec = &Uint8{}

func (i *Uint8) Encode(value any, into []byte) ([]byte, error) {
	v, ok := value.(uint8)
	if !ok {
		return nil, fmt.Errorf("%w: %T is not an uint8", types.ErrInvalidType, value)
	}
	return append(into, v), nil
}

func (i *Uint8) Decode(encoded []byte) (any, []byte, error) {
	return encodings.SafeDecode[uint8](encoded, 1, func(encoded []byte) byte { return encoded[0] })
}

func (*Uint8) GetType() reflect.Type {
	return reflect.TypeOf(uint8(0))
}

func (*Uint8) Size(int) (int, error) {
	return 1, nil
}

func (*Uint8) FixedSize() (int, error) {
	return 1, nil
}

func (e *endianEncoder) Uint8() encodings.TypeCodec {
	return &Uint8{encoder: e.encoder}
}

type Int16 struct{ encoder }

var _ encodings.TypeCodec = &Int16{}

func (i *Int16) Encode(value any, into []byte) ([]byte, error) {
	v, ok := value.(int16)
	if !ok {
		return nil, fmt.Errorf("%w: %T is not an int16", types.ErrInvalidType, value)
	}
	return i.encoder.AppendUint16(into, uint16(v)), nil
}

func (i *Int16) Decode(encoded []byte) (any, []byte, error) {
	ui, remaining, err := encodings.SafeDecode[uint16](encoded, 2, i.encoder.Uint16)
	return int16(ui), remaining, err
}

func (*Int16) GetType() reflect.Type {
	return reflect.TypeOf(int16(0))
}

func (*Int16) Size(int) (int, error) {
	return 2, nil
}

func (*Int16) FixedSize() (int, error) {
	return 2, nil
}

func (e *endianEncoder) Int16() encodings.TypeCodec {
	return &Int16{encoder: e.encoder}
}

type Uint16 struct{ encoder }

var _ encodings.TypeCodec = &Uint16{}

func (i *Uint16) Encode(value any, into []byte) ([]byte, error) {
	v, ok := value.(uint16)
	if !ok {
		return nil, fmt.Errorf("%w: %T is not an uint16", types.ErrInvalidType, value)
	}
	return i.encoder.AppendUint16(into, v), nil
}

func (i *Uint16) Decode(encoded []byte) (any, []byte, error) {
	return encodings.SafeDecode[uint16](encoded, 2, i.encoder.Uint16)
}

func (*Uint16) GetType() reflect.Type {
	return reflect.TypeOf(uint16(0))
}

func (*Uint16) Size(int) (int, error) {
	return 2, nil
}

func (*Uint16) FixedSize() (int, error) {
	return 2, nil
}

func (e *endianEncoder) Uint16() encodings.TypeCodec {
	return &Uint16{encoder: e.encoder}
}

type Int32 struct{ encoder }

var _ encodings.TypeCodec = &Int32{}

func (i *Int32) Encode(value any, into []byte) ([]byte, error) {
	v, ok := value.(int32)
	if !ok {
		return nil, fmt.Errorf("%w: %T is not an int32", types.ErrInvalidType, value)
	}
	return i.encoder.AppendUint32(into, uint32(v)), nil
}

func (i *Int32) Decode(encoded []byte) (any, []byte, error) {
	ui, remaining, err := encodings.SafeDecode[uint32](encoded, 4, i.encoder.Uint32)
	return int32(ui), remaining, err
}

func (*Int32) GetType() reflect.Type {
	return reflect.TypeOf(int32(0))
}

func (*Int32) Size(int) (int, error) {
	return 4, nil
}

func (*Int32) FixedSize() (int, error) {
	return 4, nil
}

func (e *endianEncoder) Int32() encodings.TypeCodec {
	return &Int32{encoder: e.encoder}
}

type Uint32 struct{ encoder }

var _ encodings.TypeCodec = &Uint32{}

func (i *Uint32) Encode(value any, into []byte) ([]byte, error) {
	v, ok := value.(uint32)
	if !ok {
		return nil, fmt.Errorf("%w: %T is not an uint32", types.ErrInvalidType, value)
	}
	return i.encoder.AppendUint32(into, v), nil
}

func (i *Uint32) Decode(encoded []byte) (any, []byte, error) {
	return encodings.SafeDecode[uint32](encoded, 4, i.encoder.Uint32)
}

func (*Uint32) GetType() reflect.Type {
	return reflect.TypeOf(uint32(0))
}

func (*Uint32) Size(int) (int, error) {
	return 4, nil
}

func (*Uint32) FixedSize() (int, error) {
	return 4, nil
}

func (e *endianEncoder) Uint32() encodings.TypeCodec {
	return &Uint32{encoder: e.encoder}
}

type Int64 struct{ encoder }

var _ encodings.TypeCodec = &Int64{}

func (i *Int64) Encode(value any, into []byte) ([]byte, error) {
	v, ok := value.(int64)
	if !ok {
		return nil, fmt.Errorf("%w: %T is not an int64", types.ErrInvalidType, value)
	}
	return i.encoder.AppendUint64(into, uint64(v)), nil
}

func (i *Int64) Decode(encoded []byte) (any, []byte, error) {
	ui, remaining, err := encodings.SafeDecode[uint64](encoded, 8, i.encoder.Uint64)
	return int64(ui), remaining, err
}

func (*Int64) GetType() reflect.Type {
	return reflect.TypeOf(int64(0))
}

func (*Int64) Size(int) (int, error) {
	return 8, nil
}

func (*Int64) FixedSize() (int, error) {
	return 8, nil
}

func (e *endianEncoder) Int64() encodings.TypeCodec {
	return &Int64{encoder: e.encoder}
}

type Uint64 struct{ encoder }

var _ encodings.TypeCodec = &Uint64{}

func (i *Uint64) Encode(value any, into []byte) ([]byte, error) {
	v, ok := value.(uint64)
	if !ok {
		return nil, fmt.Errorf("%w: %T is not an uint64", types.ErrInvalidType, value)
	}
	return i.encoder.AppendUint64(into, v), nil
}

func (i *Uint64) Decode(encoded []byte) (any, []byte, error) {
	return encodings.SafeDecode[uint64](encoded, 8, i.encoder.Uint64)
}

func (*Uint64) GetType() reflect.Type {
	return reflect.TypeOf(uint64(0))
}

func (*Uint64) Size(int) (int, error) {
	return 8, nil
}

func (*Uint64) FixedSize() (int, error) {
	return 8, nil
}

func (e *endianEncoder) Uint64() encodings.TypeCodec {
	return &Uint64{encoder: e.encoder}
}

func (e *endianEncoder) Int(bytes uint) (encodings.TypeCodec, error) {
	switch bytes {
	case 1:
		return &intCodec{
			codec:   &Int8{encoder: e.encoder},
			toInt:   func(v any) int { return int(v.(int8)) },
			fromInt: func(v int) any { return int8(v) },
		}, nil
	case 2:
		return &intCodec{
			codec:   &Int16{encoder: e.encoder},
			toInt:   func(v any) int { return int(v.(int16)) },
			fromInt: func(v int) any { return int16(v) },
		}, nil
	case 4:
		return &intCodec{
			codec:   &Int32{encoder: e.encoder},
			toInt:   func(v any) int { return int(v.(int32)) },
			fromInt: func(v int) any { return int32(v) },
		}, nil
	case 8:
		return &intCodec{
			codec:   &Int64{encoder: e.encoder},
			toInt:   func(v any) int { return int(v.(int64)) },
			fromInt: func(v int) any { return int64(v) },
		}, nil
	default:
		c, err := NewBigInt(bytes, true, e.bigIntEncoder)
		return &intCodec{
			codec:   c,
			toInt:   func(v any) int { return int(v.(*big.Int).Int64()) },
			fromInt: func(v int) any { return big.NewInt(int64(v)) },
		}, err
	}
}

func (e *endianEncoder) Uint(bytes uint) (encodings.TypeCodec, error) {
	switch bytes {
	case 1:
		return &uintCodec{
			codec:    &Uint8{encoder: e.encoder},
			toUint:   func(v any) uint { return uint(v.(uint8)) },
			fromUint: func(v uint) any { return uint8(v) },
		}, nil
	case 2:
		return &uintCodec{
			codec:    &Uint16{encoder: e.encoder},
			toUint:   func(v any) uint { return uint(v.(uint16)) },
			fromUint: func(v uint) any { return uint16(v) },
		}, nil
	case 4:
		return &uintCodec{
			codec:    &Uint32{encoder: e.encoder},
			toUint:   func(v any) uint { return uint(v.(uint32)) },
			fromUint: func(v uint) any { return uint32(v) },
		}, nil
	case 8:
		return &uintCodec{
			codec:    &Uint64{encoder: e.encoder},
			toUint:   func(v any) uint { return uint(v.(uint64)) },
			fromUint: func(v uint) any { return uint64(v) },
		}, nil
	default:
		c, err := NewBigInt(bytes, false, e.bigIntEncoder)
		return &uintCodec{
			codec:    c,
			toUint:   func(v any) uint { return uint(v.(*big.Int).Uint64()) },
			fromUint: func(v uint) any { return new(big.Int).SetUint64(uint64(v)) },
		}, err
	}
}
