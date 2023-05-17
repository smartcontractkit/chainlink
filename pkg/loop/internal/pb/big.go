package pb

import (
	"math/big"
)

func NewBigIntFromInt(b *big.Int) *BigInt {
	if b == nil {
		return nil
	}
	return &BigInt{
		Negative: b.Sign() < 0,
		Value:    b.Bytes(),
	}
}

func (b *BigInt) Int() *big.Int {
	if b == nil {
		return nil
	}
	i := new(big.Int)
	i.SetBytes(b.Value)
	if b.Negative {
		i = i.Neg(i)
	}
	return i
}
