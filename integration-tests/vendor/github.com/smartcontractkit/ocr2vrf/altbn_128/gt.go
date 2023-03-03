package altbn_128

import (
	"github.com/smartcontractkit/ocr2vrf/altbn_128/scalar"

	"go.dedis.ch/kyber/v3"
)

type GT struct{}

var _ kyber.Group = (*GT)(nil)

func (c *GT) String() string {
	return "AltBN-128 GT"
}

func (c *GT) ScalarLen() int {
	panic("not implemented")
}

func (c *GT) Scalar() kyber.Scalar {

	return scalar.NewScalarInt64(0)
}

func (c *GT) PointLen() int {
	panic("not implemented")
}

func (c *GT) Point() kyber.Point {
	return newGTPoint()
}
