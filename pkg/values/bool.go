package values

import (
	"fmt"

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
	tb, ok := to.(*bool)
	if !ok {
		return fmt.Errorf("cannot unwrap to value of type: %T", to)
	}

	if tb == nil {
		return fmt.Errorf("cannot unwrap to nil pointer")
	}

	*tb = b.Underlying
	return nil
}
