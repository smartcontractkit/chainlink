package binary

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"slices"

	"github.com/smartcontractkit/libocr/bigbigendian"

	"github.com/smartcontractkit/chainlink-common/pkg/codec/encodings"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type encoder interface {
	binary.AppendByteOrder
	binary.ByteOrder
}

type bigIntEncoder interface {
	serializeSigned(size int, i *big.Int) ([]byte, error)
	deserializeSigned(size int, b []byte) *big.Int
	serializeUnsigned(size int, i *big.Int) ([]byte, error)
	deserializeUnsigned(size int, b []byte) *big.Int
}

type endianEncoder struct {
	encoder
	bigIntEncoder
}

func (e *endianEncoder) Bool() encodings.TypeCodec {
	return Bool{}
}

func (e *endianEncoder) String(maxLen uint) (encodings.TypeCodec, error) {
	return NewString(maxLen, e)
}

func (e *endianEncoder) Float32() encodings.TypeCodec {
	return &Float32{encoder: e.encoder}
}

func (e *endianEncoder) Float64() encodings.TypeCodec {
	return &Float64{encoder: e.encoder}
}

func (e *endianEncoder) OracleID() encodings.TypeCodec {
	return &OracleID{}
}

func (e *endianEncoder) BigInt(bytes uint, signed bool) (encodings.TypeCodec, error) {
	return NewBigInt(bytes, signed, e.bigIntEncoder)
}

func BigEndian() encodings.Builder {
	return bigEndian
}

var bigEndian = &endianEncoder{
	encoder:       binary.BigEndian,
	bigIntEncoder: bigBigInt{},
}

type bigBigInt struct{}

func (bigBigInt) serializeSigned(size int, i *big.Int) ([]byte, error) {
	return bigbigendian.SerializeSigned(size, i)
}

func (bigBigInt) serializeUnsigned(size int, i *big.Int) ([]byte, error) {
	if i.Sign() < 0 {
		return nil, fmt.Errorf("%w: cannot encode %v as unsigned", types.ErrInvalidType, i)
	}

	if i.BitLen() > size*8 {
		return nil, fmt.Errorf("%w: %v doesn't fit into a %v-bytes", types.ErrInvalidType, i, size)
	}

	bytes := make([]byte, size)
	i.FillBytes(bytes)
	return bytes, nil
}

func (bigBigInt) deserializeSigned(size int, b []byte) *big.Int {
	bi, _ := bigbigendian.DeserializeSigned(size, b)
	return bi
}

func (bigBigInt) deserializeUnsigned(_ int, b []byte) *big.Int {
	return new(big.Int).SetBytes(b)
}

func LittleEndian() encodings.Builder {
	return littleEndian
}

var littleEndian = &endianEncoder{
	encoder:       binary.LittleEndian,
	bigIntEncoder: littleBigInt{},
}

type littleBigInt struct{}

func (littleBigInt) serializeSigned(size int, i *big.Int) ([]byte, error) {
	bi, err := bigEndian.serializeSigned(size, i)
	slices.Reverse(bi)
	return bi, err
}

func (littleBigInt) serializeUnsigned(size int, i *big.Int) ([]byte, error) {
	bi, err := bigEndian.serializeUnsigned(size, i)
	slices.Reverse(bi)
	return bi, err
}

func (littleBigInt) deserializeSigned(size int, b []byte) *big.Int {
	slices.Reverse(b)
	bi, _ := bigbigendian.DeserializeSigned(size, b)
	return bi
}

func (littleBigInt) deserializeUnsigned(_ int, b []byte) *big.Int {
	slices.Reverse(b)
	return new(big.Int).SetBytes(b)
}
