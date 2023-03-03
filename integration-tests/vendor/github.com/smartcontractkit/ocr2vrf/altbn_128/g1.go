package altbn_128

import (
	"crypto/cipher"
	"math/big"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/sign/anon"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"

	"github.com/smartcontractkit/ocr2vrf/altbn_128/scalar"
)

type G1 struct{ r cipher.Stream }

var _ kyber.Group = (*G1)(nil)
var _ anon.Suite = (*G1)(nil)

func newG1() *G1 {
	return &G1{}
}

func (g *G1) String() string {
	return "AltBN-128 G₁"
}

func (g *G1) ScalarLen() int {
	return g1ScalarLength
}

func (g *G1) Scalar() kyber.Scalar {
	return scalar.NewScalarInt64(0)
}

func (g *G1) PointLen() int {
	return g1PointLength
}

func (g *G1) Point() kyber.Point {
	return newG1Point()
}

var g1ScalarLength, g1PointLength int
var zero, one *big.Int
var null, g1Base *g1Point

func init() {
	zero, one = big.NewInt(0), big.NewInt(1)

	b, err := new(G1).Scalar().Zero().MarshalBinary()
	if err != nil {
		panic(err)
	}
	g1ScalarLength = len(b)

	rawG1Base := new(bn256.G1).ScalarBaseMult(one)
	g1Base = &g1Point{rawG1Base}
	negG1Base := new(bn256.G1).Neg(rawG1Base)
	g1Null := new(bn256.G1).Add(rawG1Base, negG1Base)
	null = &g1Point{g1Null}
	b, err = null.MarshalBinary()
	if err != nil {
		panic(err)
	}
	g1PointLength = len(b)
}
