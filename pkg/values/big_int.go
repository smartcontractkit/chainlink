package values

import (
	"fmt"
	"math/big"

	"github.com/smartcontractkit/chainlink-common/pkg/values/pb"
)

type BigInt struct {
	Underlying big.Int
}

func NewBigInt(b big.Int) *BigInt {
	return &BigInt{Underlying: b}
}

func (b *BigInt) proto() *pb.Value {
	return pb.NewBigIntValue(b.Underlying.Bytes())
}

func (b *BigInt) Unwrap() (any, error) {
	return b.Underlying, nil
}

func (b *BigInt) UnwrapTo(to any) error {
	tb, ok := to.(*big.Int)
	if !ok {
		return fmt.Errorf("cannot unwrap to value of type: %T", to)
	}

	if tb == nil {
		return fmt.Errorf("cannot unwrap to nil pointer: %+v", to)
	}

	*tb = b.Underlying
	return nil
}
