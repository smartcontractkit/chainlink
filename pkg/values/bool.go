package values

import (
	"github.com/smartcontractkit/chainlink-common/pkg/values/pb"
)

type Bool struct {
	Underlying bool
}

func NewBool(b bool) (*Bool, error) {
	return &Bool{Underlying: b}, nil
}

func (b *Bool) Proto() (*pb.Value, error) {
	return pb.NewBoolValue(b.Underlying)
}

func (b *Bool) Unwrap() (any, error) {
	return b.Underlying, nil
}
