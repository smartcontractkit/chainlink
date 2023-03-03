package altbn_128

import (
	"math/big"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

type g1Interface interface {
	String() string
	ScalarBaseMult(k *big.Int) *bn256.G1
	ScalarMult(a *bn256.G1, k *big.Int) *bn256.G1
	Add(a, b *bn256.G1) *bn256.G1
	Neg(a *bn256.G1) *bn256.G1
	Set(a *bn256.G1) *bn256.G1
	Marshal() []byte
	Unmarshal(m []byte) ([]byte, error)
}

var _ g1Interface = (*bn256.G1)(nil)

type g2Interface interface {
	String() string
	ScalarBaseMult(k *big.Int) *bn256.G2
	ScalarMult(a *bn256.G2, k *big.Int) *bn256.G2
	Add(a, b *bn256.G2) *bn256.G2
	Neg(a *bn256.G2) *bn256.G2
	Set(a *bn256.G2) *bn256.G2
	Marshal() []byte
	Unmarshal(m []byte) ([]byte, error)
}

var _ g2Interface = (*bn256.G2)(nil)
