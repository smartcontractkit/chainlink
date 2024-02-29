package values

import (
	"github.com/smartcontractkit/chainlink-common/pkg/values/pb"
)

type String struct {
	Underlying string
}

func NewString(s string) *String {
	return &String{Underlying: s}
}

func (s *String) Proto() *pb.Value {
	return pb.NewStringValue(s.Underlying)
}

func (s *String) Unwrap() (any, error) {
	return s.Underlying, nil
}
