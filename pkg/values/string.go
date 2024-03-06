package values

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/values/pb"
)

type String struct {
	Underlying string
}

func NewString(s string) *String {
	return &String{Underlying: s}
}

func (s *String) proto() *pb.Value {
	return pb.NewStringValue(s.Underlying)
}

func (s *String) Unwrap() (any, error) {
	return s.Underlying, nil
}

func (s *String) UnwrapTo(to any) error {
	tv, ok := to.(*string)
	if !ok {
		return fmt.Errorf("cannot unwrap to type %T", to)
	}

	if tv == nil {
		return fmt.Errorf("cannot unwrap to nil pointer: %+v", to)
	}

	*tv = s.Underlying
	return nil
}
