package values

import (
	"github.com/smartcontractkit/chainlink-common/pkg/values/pb"
)

type Nil struct {
}

func NewNil() (*Nil, error) {
	return &Nil{}, nil
}

func (n *Nil) Proto() (*pb.Value, error) {
	return pb.NewNilValue()
}

func (n *Nil) Unwrap() (any, error) {
	return nil, nil
}
