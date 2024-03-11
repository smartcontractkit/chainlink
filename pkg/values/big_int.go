package values

import (
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
	return unwrapTo(b.Underlying, to)
}
