package values

import (
	"fmt"
	"math"

	"github.com/smartcontractkit/chainlink-common/pkg/values/pb"
)

type Int64 struct {
	Underlying int64
}

func NewInt64(i int64) *Int64 {
	return &Int64{Underlying: i}
}

func (i *Int64) proto() *pb.Value {
	return pb.NewInt64Value(i.Underlying)
}

func (i *Int64) Unwrap() (any, error) {
	return i.Underlying, nil
}

func (i *Int64) UnwrapTo(to any) error {
	if to == nil {
		return fmt.Errorf("cannot unwrap to nil pointer: %+v", to)
	}

	switch tv := to.(type) {
	case *int64:
		*tv = i.Underlying
		return nil
	case *int:
		if i.Underlying > math.MaxInt {
			return fmt.Errorf("cannot unwrap int64 to int: number would overlflow %d", i)
		}

		*tv = int(i.Underlying)
		return nil
	}

	return fmt.Errorf("cannot unwrap to type %T", to)
}
