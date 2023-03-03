package altbn_128

import (
	"crypto/cipher"

	"github.com/smartcontractkit/ocr2vrf/altbn_128/scalar"

	"go.dedis.ch/kyber/v3"
)

type G2 struct{ r cipher.Stream }

var _ kyber.Group = (*G2)(nil)

func (g *G2) String() string {
	return "AltBN-128 G₂"
}

func (g *G2) ScalarLen() int {
	panic("not implemented")
}

func (g *G2) Scalar() kyber.Scalar {
	return scalar.NewScalarInt64(0)
}

func (g *G2) PointLen() int {
	return len(g.Point().(*g2Point).G2.Marshal())
}

func (g *G2) Point() kyber.Point {
	return newG2Point()
}
