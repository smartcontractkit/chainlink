package values

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/values/pb"
)

type Bytes struct {
	Underlying []byte
}

func NewBytes(b []byte) *Bytes {
	return &Bytes{Underlying: b}
}

func (b *Bytes) proto() *pb.Value {
	return pb.NewBytesValue(b.Underlying)
}

func (b *Bytes) Unwrap() (any, error) {
	return b.Underlying, nil
}

func (b *Bytes) UnwrapTo(to any) error {
	tb, ok := to.(*[]byte)
	if !ok {
		return fmt.Errorf("can only unwrap to a byte array, got type %T", to)
	}

	if tb == nil {
		return fmt.Errorf("cannot unwrap to nil pointer: %+v", to)
	}

	*tb = b.Underlying
	return nil
}
