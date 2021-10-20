package vrfkey

import (
	"math/big"
)

func MustNewPrivateKey(rawKey *big.Int) *PrivateKey {
	k, err := newPrivateKey(rawKey)
	if err != nil {
		panic(err)
	}
	return k
}
