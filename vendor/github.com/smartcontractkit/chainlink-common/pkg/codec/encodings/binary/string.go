package binary

import (
	"fmt"
	"math/bits"
	"reflect"

	"github.com/smartcontractkit/chainlink-common/pkg/codec/encodings"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

func NewString(maxLength uint, encoder encodings.Builder) (encodings.TypeCodec, error) {
	headerSize := uint((bits.Len(maxLength) + 7) / 8)
	sizeEncoder, err := encoder.Int(headerSize)
	if err != nil {
		return nil, err
	}

	codec, err := encodings.NewSlice(encoder.Uint8(), sizeEncoder)
	if err != nil {
		return nil, err
	}

	return stringCodec{codec: codec, maxLength: maxLength}, nil
}

type stringCodec struct {
	codec     encodings.TypeCodec
	maxLength uint
}

func (s stringCodec) Encode(value any, into []byte) ([]byte, error) {
	if str, ok := value.(string); ok {
		if uint(len(str)) > s.maxLength {
			return nil, fmt.Errorf("%w string longer than max length %d", types.ErrInvalidType, s.maxLength)
		}
		return s.codec.Encode([]byte(str), into)
	}
	return nil, fmt.Errorf("%w expected string, got %T", types.ErrInvalidType, value)
}

func (s stringCodec) Decode(encoded []byte) (any, []byte, error) {
	bytes, remaining, err := s.codec.Decode(encoded)
	if bytes == nil {
		return nil, remaining, err
	}

	sVal := string(bytes.([]byte))
	if uint(len(sVal)) > s.maxLength {
		return nil, nil, fmt.Errorf("%w string longer than max length %d", types.ErrInvalidEncoding, s.maxLength)
	}

	return sVal, remaining, err
}

func (s stringCodec) GetType() reflect.Type {
	return reflect.TypeOf("")
}

func (s stringCodec) Size(_ int) (int, error) {
	return 0, fmt.Errorf("%w strings are not sided with number of reports", types.ErrInvalidType)
}

func (s stringCodec) FixedSize() (int, error) {
	return 0, fmt.Errorf("%w strings do not have a fixed size", types.ErrInvalidType)
}

var _ encodings.TypeCodec = stringCodec{}
