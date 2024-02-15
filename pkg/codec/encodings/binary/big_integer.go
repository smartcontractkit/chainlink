package binary

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/smartcontractkit/libocr/bigbigendian"

	"github.com/smartcontractkit/chainlink-common/pkg/codec/encodings"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

func NewBigInt(numBytes uint, signed bool, intEncoder bigIntEncoder) (encodings.TypeCodec, error) {
	if numBytes > bigbigendian.MaxSize {
		return nil, fmt.Errorf(
			"%w: numBytes is %v, but must be between 1 and %v", types.ErrInvalidConfig, numBytes, bigbigendian.MaxSize)
	}
	return &bigInt{
		NumBytes:   int(numBytes),
		Signed:     signed,
		intEncoder: intEncoder,
	}, nil
}

type bigInt struct {
	NumBytes   int
	Signed     bool
	intEncoder bigIntEncoder
}

var _ encodings.TypeCodec = &bigInt{}

func (i *bigInt) Encode(value any, into []byte) ([]byte, error) {
	bi, ok := value.(*big.Int)
	if !ok {
		return nil, fmt.Errorf("%w: expected big.Int, got %T", types.ErrInvalidType, value)
	}

	if i.Signed {
		bytes, err := i.intEncoder.serializeSigned(i.NumBytes, bi)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", types.ErrInvalidType, err)
		}

		return append(into, bytes...), nil
	}

	bytes, err := i.intEncoder.serializeUnsigned(i.NumBytes, bi)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", types.ErrInvalidType, err)
	}

	return append(into, bytes...), nil
}

func (i *bigInt) Decode(encoded []byte) (any, []byte, error) {
	if i.Signed {
		return encodings.SafeDecode[*big.Int](encoded, i.NumBytes, func(bytes []byte) *big.Int {
			return i.intEncoder.deserializeSigned(i.NumBytes, bytes)
		})
	}

	return encodings.SafeDecode[*big.Int](encoded, i.NumBytes, func(bytes []byte) *big.Int {
		return i.intEncoder.deserializeUnsigned(i.NumBytes, bytes)
	})
}

func (i *bigInt) GetType() reflect.Type {
	return reflect.TypeOf((*big.Int)(nil))
}

func (i *bigInt) Size(_ int) (int, error) {
	return i.NumBytes, nil
}

func (i *bigInt) FixedSize() (int, error) {
	return i.NumBytes, nil
}
