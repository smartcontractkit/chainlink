package altbn_128

import (
	"bytes"
	"crypto/cipher"
	"fmt"
	"io"
	"math/big"

	"github.com/pkg/errors"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/mod"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"

	"github.com/smartcontractkit/ocr2vrf/altbn_128/scalar"
)

type g1Point struct{ G1 g1Interface }

var _ kyber.Point = (*g1Point)(nil)

func newG1Point() (rv *g1Point) {
	rv = &g1Point{new(bn256.G1)}
	rv.G1.Set(null.G1.(*bn256.G1))
	return
}

var Order = bn256.Order

func (p *g1Point) ensureP() *g1Point {
	if p == nil {
		panic("cannot assign to nil point")
	}
	if p.G1 == nil {
		p.G1 = new(bn256.G1)
		p.G1.Set(null.G1.(*bn256.G1))
	}
	return p
}

func bytesZeroThroughPMinusOne(b []byte) bool {
	return big.NewInt(0).SetBytes(b).Cmp(bn256.P) < 0
}

func (p *g1Point) MarshalBinary() (data []byte, err error) {
	rawData := p.G1.Marshal()
	if len(rawData) != 64 {
		return nil, errors.Errorf("wrong format from bn256 marshalling logic")
	}
	rawX, rawY := rawData[:32], rawData[32:]
	if !bytesZeroThroughPMinusOne(rawX) {
		return nil, errors.Errorf("x ordinate 0x%x too large", rawX)
	}
	if !bytesZeroThroughPMinusOne(rawY) {
		return nil, errors.Errorf("y ordinate 0x%x too large", rawY)
	}
	rawX[0] |= (rawY[31] & 1) << 7
	return rawX, nil
}

func i(x int64) *scalar.Scalar { return scalar.NewScalarInt64(x) }
func m(x int64) *mod.Int       { return mod.NewInt64(x, bn256.P) }

var three = m(3)

func (p *g1Point) UnmarshalBinary(data []byte) error {
	if len(data) != 32 {
		return errors.Errorf("attempt to unmarshal g1Point data of wrong length")
	}
	if bytes.Equal(data, uncompressedZero[:32]) {

		p.Null()
		return nil
	}

	xData := make([]byte, len(data))
	copy(xData, data)
	xData[0] &= 0x7F
	if !bytesZeroThroughPMinusOne(xData) {
		return errors.Errorf("x ordinate 0x%x too large", xData)
	}
	x := mod.NewIntBytes(xData, bn256.P, mod.BigEndian)
	tmp := m(0).Mul(x, x)
	_ = tmp.Mul(tmp, x)
	ySq := tmp.Add(tmp, three)
	y := m(0)
	if !y.Sqrt(ySq) {
		return errors.Errorf("no point on curve with given x ordinate 0x%s", x)
	}
	yParity := (data[0] & 0x80) == 0x80
	yData, err := y.MarshalBinary()
	if err != nil {
		return errors.Wrapf(err, "while marshalling y data")
	}
	currentParity := (yData[31] & 1) == 1
	if yParity != currentParity {

		_ = y.Neg(y)
	}
	yData, err = y.MarshalBinary()
	if err != nil {
		return errors.Wrap(err, "could not re-marshal y after possible negation")
	}
	if (yData[31]&1 == 1) != yParity {
		panic("failed to set correct parity for y")
	}
	p.ensureP()
	if _, err := p.G1.Unmarshal(append(xData, yData...)); err != nil {

		return errors.Wrap(err, "while unmarshalling to bn256 point")
	}
	return nil
}

var uncompressedZero = make([]byte, 64)

func (p *g1Point) String() string {
	if p == nil {
		return "(*g1Point)(nil)"
	}
	b := p.G1.Marshal()
	if bytes.Equal(b, uncompressedZero) {
		return "g1Point{∞}"
	}
	return fmt.Sprintf("g1Point{0x%x,0x%x}", b[:32], b[32:])
}

func (p *g1Point) Equal(p2 kyber.Point) bool {
	if p == nil || p2 == nil {
		return false
	}
	p2G1, ok := p2.(*g1Point)
	return p != nil && ok && bytes.Equal(p.G1.Marshal(), p2G1.G1.Marshal())
}

func (p *g1Point) Null() kyber.Point {
	p.ensureP()
	p.G1.ScalarBaseMult(big.NewInt(0))
	return p
}

func (p *g1Point) Base() kyber.Point {
	p.ensureP()
	p.G1.ScalarBaseMult(big.NewInt(1))
	return p
}

func (p *g1Point) Pick(rand cipher.Stream) kyber.Point {
	p.ensureP()

	_ = p.G1.ScalarBaseMult(i(0).Pick(rand).(*scalar.Scalar).Big())
	return p
}

func (p *g1Point) Set(p2 kyber.Point) kyber.Point {
	p.ensureP()
	p.G1.Set(p2.(*g1Point).G1.(*bn256.G1))
	return p
}

func (p *g1Point) Clone() kyber.Point {
	rv := newG1Point()
	rv.Set(p)
	return rv
}

func (p *g1Point) EmbedLen() int                                  { panic("not implemented") }
func (p *g1Point) Embed(data []byte, r cipher.Stream) kyber.Point { panic("not implemented") }
func (p *g1Point) Data() ([]byte, error)                          { panic("not implemented") }

func (p *g1Point) Add(a kyber.Point, b kyber.Point) kyber.Point {
	p.ensureP()
	ap := a.(*g1Point).G1
	bp := b.(*g1Point).G1
	p.G1.Add(ap.(*bn256.G1), bp.(*bn256.G1))
	return p
}

func (p *g1Point) Sub(a kyber.Point, b kyber.Point) kyber.Point {
	return p.Add(a, newG1Point().Neg(b))
}

func (p *g1Point) Neg(a kyber.Point) kyber.Point {
	aG1 := a.(*g1Point)
	p.ensureP()
	p.G1.Neg(aG1.G1.(*bn256.G1))
	return p
}

func (p *g1Point) Mul(s kyber.Scalar, p2 kyber.Point) kyber.Point {
	sm := s.(*scalar.Scalar)
	p.ensureP()
	if p2 == nil {
		p2 = newG1Point().Base()
	}
	p.G1.ScalarMult(p2.(*g1Point).G1.(*bn256.G1), sm.Big())
	return p
}

func (p *g1Point) MarshalSize() int { return g1PointLength }

func (p *g1Point) MarshalTo(w io.Writer) (numBytesWritten int, err error) {
	data, err := p.MarshalBinary()
	if err != nil {
		return 0, errors.Wrapf(err, "while marshalling for writing")
	}
	n, err := w.Write(data)
	return n, errors.Wrapf(err, "while writing marshaled value 0x%x", data)
}

func (p *g1Point) UnmarshalFrom(r io.Reader) (numBytesRead int, err error) {
	if strm, ok := r.(cipher.Stream); ok {
		p.Pick(strm)
		return -1, nil
	}
	buf := make([]byte, p.MarshalSize())
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return n, errors.Wrap(err, "while reading for unmarshalling")
	}
	return n, errors.Wrapf(p.UnmarshalBinary(buf), "while unmarshalling 0x%x", buf)
}

func LongMarshal(p kyber.Point) (rv [64]byte) {
	m := p.(*g1Point).G1.Marshal()
	if len(m) != 64 {
		panic(fmt.Errorf("wrong length for serialized G1 point 0x%x from %s", m, p))
	}
	copy(rv[:], m)
	return
}

func IsAltBN128G1Point(p kyber.Point) bool {
	_, ok := p.(*g1Point)
	return ok
}
