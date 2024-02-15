package binary

import (
	"fmt"
	"reflect"

	"github.com/smartcontractkit/chainlink-common/pkg/codec/encodings"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type Bool struct{}

var _ encodings.TypeCodec = Bool{}

func (Bool) Encode(value any, into []byte) ([]byte, error) {
	if b, ok := value.(bool); ok {
		if b {
			return append(into, 1), nil
		}
		return append(into, 0), nil
	}

	return nil, fmt.Errorf("%w: expected bool, got %T", types.ErrInvalidType, value)
}

func (Bool) Decode(encoded []byte) (any, []byte, error) {
	return encodings.SafeDecode[bool](encoded, 1, func(b []byte) bool {
		return b[0] != 0
	})
}

func (Bool) GetType() reflect.Type {
	return reflect.TypeOf(true)
}

func (Bool) Size(numItems int) (int, error) {
	return 1, nil
}

func (Bool) FixedSize() (int, error) {
	return 1, nil
}
