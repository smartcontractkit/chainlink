package binary

import (
	"fmt"
	"math"
	"reflect"

	"github.com/smartcontractkit/chainlink-common/pkg/codec/encodings"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

// Float32 follows IEEE-754 convention for float32, the specification says what the "high" and "low" bits represent
// Leaving endianness to specify the byte ordering.
type Float32 Uint32

func (f *Float32) Encode(value any, into []byte) ([]byte, error) {
	f32, ok := value.(float32)
	if ok {
		return (*Uint32)(f).Encode(math.Float32bits(f32), into)
	}

	return nil, fmt.Errorf("%w expected float32, got %T", types.ErrInvalidType, value)
}

func (f *Float32) Decode(encoded []byte) (any, []byte, error) {
	u32, remaining, err := (*Uint32)(f).Decode(encoded)
	return math.Float32frombits(u32.(uint32)), remaining, err
}

func (f *Float32) GetType() reflect.Type {
	return reflect.TypeOf(float32(0))
}

func (f *Float32) Size(_ int) (int, error) {
	return 4, nil
}

func (f *Float32) FixedSize() (int, error) {
	return 4, nil
}

var _ encodings.TypeCodec = &Float32{}
