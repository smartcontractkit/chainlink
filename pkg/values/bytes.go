package values

import (
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
	return unwrapTo(b.Underlying, to)
}
