package bn256

import (
	"crypto/cipher"
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"io"
	"math/big"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/mod"
)

var marshalPointID1 = [8]byte{'b', 'n', '2', '5', '6', '.', 'g', '1'}
var marshalPointID2 = [8]byte{'b', 'n', '2', '5', '6', '.', 'g', '2'}
var marshalPointIDT = [8]byte{'b', 'n', '2', '5', '6', '.', 'g', 't'}

type pointG1 struct {
	g *curvePoint
}

func newPointG1() *pointG1 {
	p := &pointG1{g: &curvePoint{}}
	return p
}

func (p *pointG1) Equal(q kyber.Point) bool {
	x, _ := p.MarshalBinary()
	y, _ := q.MarshalBinary()
	return subtle.ConstantTimeCompare(x, y) == 1
}

func (p *pointG1) Null() kyber.Point {
	p.g.SetInfinity()
	return p
}

func (p *pointG1) Base() kyber.Point {
	p.g.Set(curveGen)
	return p
}

func (p *pointG1) Pick(rand cipher.Stream) kyber.Point {
	s := mod.NewInt64(0, Order).Pick(rand)
	p.Base()
	p.g.Mul(p.g, &s.(*mod.Int).V)
	return p
}

func (p *pointG1) Set(q kyber.Point) kyber.Point {
	x := q.(*pointG1).g
	p.g.Set(x)
	return p
}

// Clone makes a hard copy of the point
func (p *pointG1) Clone() kyber.Point {
	q := newPointG1()
	q.g = p.g.Clone()
	return q
}

func (p *pointG1) EmbedLen() int {
	panic("bn256.G1: unsupported operation")
}

func (p *pointG1) Embed(data []byte, rand cipher.Stream) kyber.Point {
	// XXX: An approach to implement this is:
	// - Encode data as the x-coordinate of a point on y²=x³+3 where len(data)
	//   is stored in the least significant byte of x and the rest is being
	//   filled with random values, i.e., x = rand || data || len(data).
	// - Use the Tonelli-Shanks algorithm to compute the y-coordinate.
	// - Convert the new point to Jacobian coordinates and set it as p.
	panic("bn256.G1: unsupported operation")
}

func (p *pointG1) Data() ([]byte, error) {
	panic("bn256.G1: unsupported operation")
}

func (p *pointG1) Add(a, b kyber.Point) kyber.Point {
	x := a.(*pointG1).g
	y := b.(*pointG1).g
	p.g.Add(x, y) // p = a + b
	return p
}

func (p *pointG1) Sub(a, b kyber.Point) kyber.Point {
	q := newPointG1()
	return p.Add(a, q.Neg(b))
}

func (p *pointG1) Neg(q kyber.Point) kyber.Point {
	x := q.(*pointG1).g
	p.g.Neg(x)
	return p
}

func (p *pointG1) Mul(s kyber.Scalar, q kyber.Point) kyber.Point {
	if q == nil {
		q = newPointG1().Base()
	}
	t := s.(*mod.Int).V
	r := q.(*pointG1).g
	p.g.Mul(r, &t)
	return p
}

func (p *pointG1) MarshalBinary() ([]byte, error) {
	// Clone is required as we change the point
	p = p.Clone().(*pointG1)

	n := p.ElementSize()
	// Take a copy so that p is not written to, so calls to MarshalBinary
	// are threadsafe.
	pgtemp := *p.g
	pgtemp.MakeAffine()
	ret := make([]byte, p.MarshalSize())
	if pgtemp.IsInfinity() {
		return ret, nil
	}
	tmp := &gfP{}
	montDecode(tmp, &pgtemp.x)
	tmp.Marshal(ret)
	montDecode(tmp, &pgtemp.y)
	tmp.Marshal(ret[n:])
	return ret, nil
}

func (p *pointG1) MarshalID() [8]byte {
	return marshalPointID1
}

func (p *pointG1) MarshalTo(w io.Writer) (int, error) {
	buf, err := p.MarshalBinary()
	if err != nil {
		return 0, err
	}
	return w.Write(buf)
}

func (p *pointG1) UnmarshalBinary(buf []byte) error {
	n := p.ElementSize()
	if len(buf) < p.MarshalSize() {
		return errors.New("bn256.G1: not enough data")
	}
	if p.g == nil {
		p.g = &curvePoint{}
	} else {
		p.g.x, p.g.y = gfP{0}, gfP{0}
	}

	p.g.x.Unmarshal(buf)
	p.g.y.Unmarshal(buf[n:])
	montEncode(&p.g.x, &p.g.x)
	montEncode(&p.g.y, &p.g.y)

	zero := gfP{0}
	if p.g.x == zero && p.g.y == zero {
		// This is the point at infinity
		p.g.y = *newGFp(1)
		p.g.z = gfP{0}
		p.g.t = gfP{0}
	} else {
		p.g.z = *newGFp(1)
		p.g.t = *newGFp(1)
	}

	if !p.g.IsOnCurve() {
		return errors.New("bn256.G1: malformed point")
	}

	return nil
}

func (p *pointG1) UnmarshalFrom(r io.Reader) (int, error) {
	buf := make([]byte, p.MarshalSize())
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return n, err
	}
	return n, p.UnmarshalBinary(buf)
}

func (p *pointG1) MarshalSize() int {
	return 2 * p.ElementSize()
}

func (p *pointG1) ElementSize() int {
	return 256 / 8
}

func (p *pointG1) String() string {
	return "bn256.G1" + p.g.String()
}

func (p *pointG1) Hash(m []byte) kyber.Point {
	leftPad32 := func(in []byte) []byte {
		if len(in) > 32 {
			panic("input cannot be more than 32 bytes")
		}

		out := make([]byte, 32)
		copy(out[32-len(in):], in)
		return out
	}

	bigX, bigY := hashToPoint(m)
	if p.g == nil {
		p.g = new(curvePoint)
	}

	x, y := new(gfP), new(gfP)
	x.Unmarshal(leftPad32(bigX.Bytes()))
	y.Unmarshal(leftPad32(bigY.Bytes()))
	montEncode(x, x)
	montEncode(y, y)

	p.g.Set(&curvePoint{*x, *y, *newGFp(1), *newGFp(1)})
	return p
}

// hashes a byte slice into two points on a curve represented by big.Int
// ideally we want to do this using gfP, but gfP doesn't have a ModSqrt function
func hashToPoint(m []byte) (*big.Int, *big.Int) {
	// we need to convert curveB into a bigInt for our computation
	intCurveB := new(big.Int)
	{
		decodedCurveB := new(gfP)
		montDecode(decodedCurveB, curveB)
		bufCurveB := make([]byte, 32)
		decodedCurveB.Marshal(bufCurveB)
		intCurveB.SetBytes(bufCurveB)
	}

	h := sha256.Sum256(m)
	x := new(big.Int).SetBytes(h[:])
	x.Mod(x, p)

	for {
		xxx := new(big.Int).Mul(x, x)
		xxx.Mul(xxx, x)
		xxx.Mod(xxx, p)

		t := new(big.Int).Add(xxx, intCurveB)
		y := new(big.Int).ModSqrt(t, p)
		if y != nil {
			return x, y
		}

		x.Add(x, big.NewInt(1))
	}
}

type pointG2 struct {
	g *twistPoint
}

func newPointG2() *pointG2 {
	p := &pointG2{g: &twistPoint{}}
	return p
}

func (p *pointG2) Equal(q kyber.Point) bool {
	x, _ := p.MarshalBinary()
	y, _ := q.MarshalBinary()
	return subtle.ConstantTimeCompare(x, y) == 1
}

func (p *pointG2) Null() kyber.Point {
	p.g.SetInfinity()
	return p
}

func (p *pointG2) Base() kyber.Point {
	p.g.Set(twistGen)
	return p
}

func (p *pointG2) Pick(rand cipher.Stream) kyber.Point {
	s := mod.NewInt64(0, Order).Pick(rand)
	p.Base()
	p.g.Mul(p.g, &s.(*mod.Int).V)
	return p
}

func (p *pointG2) Set(q kyber.Point) kyber.Point {
	x := q.(*pointG2).g
	p.g.Set(x)
	return p
}

// Clone makes a hard copy of the field
func (p *pointG2) Clone() kyber.Point {
	q := newPointG2()
	q.g = p.g.Clone()
	return q
}

func (p *pointG2) EmbedLen() int {
	panic("bn256.G2: unsupported operation")
}

func (p *pointG2) Embed(data []byte, rand cipher.Stream) kyber.Point {
	panic("bn256.G2: unsupported operation")
}

func (p *pointG2) Data() ([]byte, error) {
	panic("bn256.G2: unsupported operation")
}

func (p *pointG2) Add(a, b kyber.Point) kyber.Point {
	x := a.(*pointG2).g
	y := b.(*pointG2).g
	p.g.Add(x, y) // p = a + b
	return p
}

func (p *pointG2) Sub(a, b kyber.Point) kyber.Point {
	q := newPointG2()
	return p.Add(a, q.Neg(b))
}

func (p *pointG2) Neg(q kyber.Point) kyber.Point {
	x := q.(*pointG2).g
	p.g.Neg(x)
	return p
}

func (p *pointG2) Mul(s kyber.Scalar, q kyber.Point) kyber.Point {
	if q == nil {
		q = newPointG2().Base()
	}
	t := s.(*mod.Int).V
	r := q.(*pointG2).g
	p.g.Mul(r, &t)
	return p
}

func (p *pointG2) MarshalBinary() ([]byte, error) {
	// Clone is required as we change the point during the operation
	p = p.Clone().(*pointG2)

	n := p.ElementSize()
	if p.g == nil {
		p.g = &twistPoint{}
	}

	p.g.MakeAffine()

	ret := make([]byte, p.MarshalSize())
	if p.g.IsInfinity() {
		return ret, nil
	}

	temp := &gfP{}
	montDecode(temp, &p.g.x.x)
	temp.Marshal(ret[0*n:])
	montDecode(temp, &p.g.x.y)
	temp.Marshal(ret[1*n:])
	montDecode(temp, &p.g.y.x)
	temp.Marshal(ret[2*n:])
	montDecode(temp, &p.g.y.y)
	temp.Marshal(ret[3*n:])

	return ret, nil
}

func (p *pointG2) MarshalID() [8]byte {
	return marshalPointID2
}

func (p *pointG2) MarshalTo(w io.Writer) (int, error) {
	buf, err := p.MarshalBinary()
	if err != nil {
		return 0, err
	}
	return w.Write(buf)
}

func (p *pointG2) UnmarshalBinary(buf []byte) error {
	n := p.ElementSize()
	if p.g == nil {
		p.g = &twistPoint{}
	}

	if len(buf) < p.MarshalSize() {
		return errors.New("bn256.G2: not enough data")
	}

	p.g.x.x.Unmarshal(buf[0*n:])
	p.g.x.y.Unmarshal(buf[1*n:])
	p.g.y.x.Unmarshal(buf[2*n:])
	p.g.y.y.Unmarshal(buf[3*n:])
	montEncode(&p.g.x.x, &p.g.x.x)
	montEncode(&p.g.x.y, &p.g.x.y)
	montEncode(&p.g.y.x, &p.g.y.x)
	montEncode(&p.g.y.y, &p.g.y.y)

	if p.g.x.IsZero() && p.g.y.IsZero() {
		// This is the point at infinity.
		p.g.y.SetOne()
		p.g.z.SetZero()
		p.g.t.SetZero()
	} else {
		p.g.z.SetOne()
		p.g.t.SetOne()

		if !p.g.IsOnCurve() {
			return errors.New("bn256.G2: malformed point")
		}
	}
	return nil
}

func (p *pointG2) UnmarshalFrom(r io.Reader) (int, error) {
	buf := make([]byte, p.MarshalSize())
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return n, err
	}
	return n, p.UnmarshalBinary(buf)
}

func (p *pointG2) MarshalSize() int {
	return 4 * p.ElementSize()
}

func (p *pointG2) ElementSize() int {
	return 256 / 8
}

func (p *pointG2) String() string {
	return "bn256.G2" + p.g.String()
}

type pointGT struct {
	g *gfP12
}

func newPointGT() *pointGT {
	p := &pointGT{g: &gfP12{}}
	return p
}

func (p *pointGT) Equal(q kyber.Point) bool {
	x, _ := p.MarshalBinary()
	y, _ := q.MarshalBinary()
	return subtle.ConstantTimeCompare(x, y) == 1
}

func (p *pointGT) Null() kyber.Point {
	p.g.Set(gfP12Inf)
	return p
}

func (p *pointGT) Base() kyber.Point {
	p.g.Set(gfP12Gen)
	return p
}

func (p *pointGT) Pick(rand cipher.Stream) kyber.Point {
	s := mod.NewInt64(0, Order).Pick(rand)
	p.Base()
	p.g.Exp(p.g, &s.(*mod.Int).V)
	return p
}

func (p *pointGT) Set(q kyber.Point) kyber.Point {
	x := q.(*pointGT).g
	p.g.Set(x)
	return p
}

// Clone makes a hard copy of the point
func (p *pointGT) Clone() kyber.Point {
	q := newPointGT()
	q.g = p.g.Clone()
	return q
}

func (p *pointGT) EmbedLen() int {
	panic("bn256.GT: unsupported operation")
}

func (p *pointGT) Embed(data []byte, rand cipher.Stream) kyber.Point {
	panic("bn256.GT: unsupported operation")
}

func (p *pointGT) Data() ([]byte, error) {
	panic("bn256.GT: unsupported operation")
}

func (p *pointGT) Add(a, b kyber.Point) kyber.Point {
	x := a.(*pointGT).g
	y := b.(*pointGT).g
	p.g.Mul(x, y)
	return p
}

func (p *pointGT) Sub(a, b kyber.Point) kyber.Point {
	q := newPointGT()
	return p.Add(a, q.Neg(b))
}

func (p *pointGT) Neg(q kyber.Point) kyber.Point {
	x := q.(*pointGT).g
	p.g.Conjugate(x)
	return p
}

func (p *pointGT) Mul(s kyber.Scalar, q kyber.Point) kyber.Point {
	if q == nil {
		q = newPointGT().Base()
	}
	t := s.(*mod.Int).V
	r := q.(*pointGT).g
	p.g.Exp(r, &t)
	return p
}

func (p *pointGT) MarshalBinary() ([]byte, error) {
	n := p.ElementSize()
	ret := make([]byte, p.MarshalSize())
	temp := &gfP{}

	montDecode(temp, &p.g.x.x.x)
	temp.Marshal(ret[0*n:])
	montDecode(temp, &p.g.x.x.y)
	temp.Marshal(ret[1*n:])
	montDecode(temp, &p.g.x.y.x)
	temp.Marshal(ret[2*n:])
	montDecode(temp, &p.g.x.y.y)
	temp.Marshal(ret[3*n:])
	montDecode(temp, &p.g.x.z.x)
	temp.Marshal(ret[4*n:])
	montDecode(temp, &p.g.x.z.y)
	temp.Marshal(ret[5*n:])
	montDecode(temp, &p.g.y.x.x)
	temp.Marshal(ret[6*n:])
	montDecode(temp, &p.g.y.x.y)
	temp.Marshal(ret[7*n:])
	montDecode(temp, &p.g.y.y.x)
	temp.Marshal(ret[8*n:])
	montDecode(temp, &p.g.y.y.y)
	temp.Marshal(ret[9*n:])
	montDecode(temp, &p.g.y.z.x)
	temp.Marshal(ret[10*n:])
	montDecode(temp, &p.g.y.z.y)
	temp.Marshal(ret[11*n:])

	return ret, nil
}

func (p *pointGT) MarshalID() [8]byte {
	return marshalPointIDT
}

func (p *pointGT) MarshalTo(w io.Writer) (int, error) {
	buf, err := p.MarshalBinary()
	if err != nil {
		return 0, err
	}
	return w.Write(buf)
}

func (p *pointGT) UnmarshalBinary(buf []byte) error {
	n := p.ElementSize()
	if len(buf) < p.MarshalSize() {
		return errors.New("bn256.GT: not enough data")
	}

	if p.g == nil {
		p.g = &gfP12{}
	}

	p.g.x.x.x.Unmarshal(buf[0*n:])
	p.g.x.x.y.Unmarshal(buf[1*n:])
	p.g.x.y.x.Unmarshal(buf[2*n:])
	p.g.x.y.y.Unmarshal(buf[3*n:])
	p.g.x.z.x.Unmarshal(buf[4*n:])
	p.g.x.z.y.Unmarshal(buf[5*n:])
	p.g.y.x.x.Unmarshal(buf[6*n:])
	p.g.y.x.y.Unmarshal(buf[7*n:])
	p.g.y.y.x.Unmarshal(buf[8*n:])
	p.g.y.y.y.Unmarshal(buf[9*n:])
	p.g.y.z.x.Unmarshal(buf[10*n:])
	p.g.y.z.y.Unmarshal(buf[11*n:])
	montEncode(&p.g.x.x.x, &p.g.x.x.x)
	montEncode(&p.g.x.x.y, &p.g.x.x.y)
	montEncode(&p.g.x.y.x, &p.g.x.y.x)
	montEncode(&p.g.x.y.y, &p.g.x.y.y)
	montEncode(&p.g.x.z.x, &p.g.x.z.x)
	montEncode(&p.g.x.z.y, &p.g.x.z.y)
	montEncode(&p.g.y.x.x, &p.g.y.x.x)
	montEncode(&p.g.y.x.y, &p.g.y.x.y)
	montEncode(&p.g.y.y.x, &p.g.y.y.x)
	montEncode(&p.g.y.y.y, &p.g.y.y.y)
	montEncode(&p.g.y.z.x, &p.g.y.z.x)
	montEncode(&p.g.y.z.y, &p.g.y.z.y)

	// TODO: check if point is on curve

	return nil
}

func (p *pointGT) UnmarshalFrom(r io.Reader) (int, error) {
	buf := make([]byte, p.MarshalSize())
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return n, err
	}
	return n, p.UnmarshalBinary(buf)
}

func (p *pointGT) MarshalSize() int {
	return 12 * p.ElementSize()
}

func (p *pointGT) ElementSize() int {
	return 256 / 8
}

func (p *pointGT) String() string {
	return "bn256.GT" + p.g.String()
}

func (p *pointGT) Finalize() kyber.Point {
	buf := finalExponentiation(p.g)
	p.g.Set(buf)
	return p
}

func (p *pointGT) Miller(p1, p2 kyber.Point) kyber.Point {
	a := p1.(*pointG1).g
	b := p2.(*pointG2).g
	p.g.Set(miller(b, a))
	return p
}

func (p *pointGT) Pair(p1, p2 kyber.Point) kyber.Point {
	a := p1.(*pointG1).g
	b := p2.(*pointG2).g
	p.g.Set(optimalAte(b, a))
	return p
}
