package altbn_128

import (
	"bytes"
	"crypto/cipher"
	"fmt"
	"io"
	"math/big"
	"sync"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"

	"github.com/smartcontractkit/ocr2vrf/altbn_128/scalar"
	"github.com/smartcontractkit/ocr2vrf/internal/util"

	"go.dedis.ch/kyber/v3"
)

type g2Point struct{ G2 g2Interface }

var _ kyber.Point = (*g2Point)(nil)

func bn256G2Null() *bn256.G2 {
	g := new(bn256.G2).ScalarBaseMult(big.NewInt(1))
	ng := new(bn256.G2).Neg(g)
	return new(bn256.G2).Add(g, ng)
}

func newG2Point() *g2Point {
	return &g2Point{bn256G2Null()}
}

func (p *g2Point) ensurePG2() {

	if p == nil {
		panic("attempt to ensure allocation of nil *g2Point")
	}
	if p.G2 == nil {
		p.G2 = bn256G2Null()
	}
}

func (p *g2Point) mustBeValidPoint() {
	if _, err := new(bn256.G2).Unmarshal(p.G2.Marshal()); err != nil {
		panic(util.WrapErrorf(err, "invalid G₂ point: 0x%x", p.G2.Marshal()))
	}
}

func (p *g2Point) MarshalSize() int {
	panic("not implemented")
}

func (p *g2Point) MarshalBinary() (data []byte, err error) {
	p.mustBeValidPoint()
	return p.G2.Marshal(), nil
}

func (p *g2Point) UnmarshalBinary(data []byte) error {
	if p == nil {
		return fmt.Errorf("can't assign to nil pointer")
	}

	p.G2 = new(bn256.G2)
	rem, err := p.G2.Unmarshal(data)
	if err != nil {
		return util.WrapErrorf(err, "while unmarshalling to G₂ point: 0x%x", data)
	}
	if len(rem) > 0 {
		errMsg := "overage of %d bytes in representation of AltBN-128 G2 point"
		return fmt.Errorf(errMsg, len(rem))
	}
	return nil
}

func (p *g2Point) Mul(s kyber.Scalar, p2 kyber.Point) kyber.Point {
	sc := s.(*scalar.Scalar)
	if p2 == nil {
		p2 = newG2Point().Base()
	}
	pP := p2.(*g2Point)
	p.ensurePG2()
	_ = p.G2.ScalarMult(pP.G2.(*bn256.G2), sc.Big())
	return p
}

func (p *g2Point) Add(a kyber.Point, b kyber.Point) kyber.Point {
	aG2 := a.(*g2Point)
	bG2 := b.(*g2Point)
	aG2.mustBeValidPoint()
	bG2.mustBeValidPoint()
	p.ensurePG2()
	_ = p.G2.Add(aG2.G2.(*bn256.G2), bG2.G2.(*bn256.G2))
	return p
}

var g2BasePoint *g2Point
var g2BasePointLock = sync.RWMutex{}

func (p *g2Point) Base() kyber.Point {
	g2BasePointLock.Lock()
	defer g2BasePointLock.Unlock()
	if g2BasePoint == nil {
		g2BasePoint = &g2Point{new(bn256.G2).ScalarBaseMult(one)}
	}
	p.Set(g2BasePoint)
	return p
}

func (p *g2Point) Null() kyber.Point {
	p.ensurePG2()
	p.G2 = bn256G2Null()
	return p
}

func (p *g2Point) Sub(a kyber.Point, b kyber.Point) kyber.Point {
	aG2 := a.(*g2Point)
	bG2 := b.Clone().Neg(b).(*g2Point)
	aG2.mustBeValidPoint()
	bG2.mustBeValidPoint()
	p.ensurePG2()
	_ = p.G2.Add(aG2.G2.(*bn256.G2), bG2.G2.(*bn256.G2))
	return p
}

func (p *g2Point) Neg(a kyber.Point) kyber.Point {
	_ = p.G2.Neg(a.(*g2Point).G2.(*bn256.G2))
	return p
}

func (p *g2Point) Set(p2 kyber.Point) kyber.Point {
	p.ensurePG2()
	p.G2 = new(bn256.G2)
	_ = p.G2.Set(p2.(*g2Point).G2.(*bn256.G2))
	return p
}

func (p *g2Point) Equal(p2 kyber.Point) bool {
	p2g2, ok := p2.(*g2Point)
	return ok && bytes.Equal(p.G2.Marshal(), p2g2.G2.Marshal())
}

func (p *g2Point) Clone() kyber.Point {
	p.ensurePG2()
	rv := newG2Point()
	rv.G2.Set(p.G2.(*bn256.G2))
	return rv
}

func (p *g2Point) Pick(rand cipher.Stream) kyber.Point {
	_ = p.G2.ScalarBaseMult(i(0).Pick(rand).(*scalar.Scalar).Big())
	return p
}

func (p *g2Point) String() string { return fmt.Sprintf("&g2Point{%s}", p.G2.String()) }

func (p *g2Point) MarshalTo(w io.Writer) (int, error)             { panic("not implemented") }
func (p *g2Point) UnmarshalFrom(r io.Reader) (int, error)         { panic("not implemented") }
func (p *g2Point) EmbedLen() int                                  { panic("not implemented") }
func (p *g2Point) Embed(data []byte, r cipher.Stream) kyber.Point { panic("not implemented") }
func (p *g2Point) Data() ([]byte, error)                          { panic("not implemented") }
