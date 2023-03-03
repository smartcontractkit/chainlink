package types

import (
	"fmt"
	"math/big"
)

type Uint256 struct {
	Low  *Felt
	High *Felt
}

func NewUint256(low, high *Felt) (*Uint256, error) {
	if low.Cmp(new(big.Int).Lsh(big.NewInt(1), 128)) >= 0 {
		return nil, fmt.Errorf("invalid low felt value")
	}
	if high.Cmp(new(big.Int).Lsh(big.NewInt(1), 128)) >= 0 {
		return nil, fmt.Errorf("invalid high felt value")
	}
	return &Uint256{
		Low:  low,
		High: high,
	}, nil
}

func (u *Uint256) Big() *big.Int {
	return new(big.Int).Add(new(big.Int).Lsh(u.High.Int, 128), u.Low.Int)
}

func (u *Uint256) String() string {
	return u.Big().String()
}

func Uint256FromBig(b *big.Int) (*Uint256, error) {
	if b.Cmp(new(big.Int).Lsh(big.NewInt(1), 256)) >= 0 {
		return nil, fmt.Errorf("invalid uint256 value")
	}
	return &Uint256{
		Low:  &Felt{Int: new(big.Int).Mod(b, new(big.Int).Exp(big.NewInt(2), big.NewInt(128), nil))},
		High: &Felt{Int: new(big.Int).Rsh(b, 128)},
	}, nil
}
