package values

import (
	"github.com/smartcontractkit/chainlink-common/pkg/values/pb"
)

type Bool struct {
	Underlying bool
}

func NewBool(b bool) *Bool {
	return &Bool{Underlying: b}
}

func (b *Bool) proto() *pb.Value {
	return pb.NewBoolValue(b.Underlying)
}

func (b *Bool) Unwrap() (any, error) {
	return b.Underlying, nil
}

func (b *Bool) UnwrapTo(to any) error {
	return unwrapTo[bool](b.Underlying, to)
}

func (b *Bool) Copy() Value {
	return &Bool{Underlying: b.Underlying}
}
