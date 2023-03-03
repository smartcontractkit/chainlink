// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package edwards25519

// Group elements are members of the elliptic curve -x^2 + y^2 = 1 + d * x^2 *
// y^2 where d = -121665/121666.
//
// Several representations are used:
//   projectiveGroupElement: (X:Y:Z) satisfying x=X/Z, y=Y/Z
//   extendedGroupElement: (X:Y:Z:T) satisfying x=X/Z, y=Y/Z, XY=ZT
//   completedGroupElement: ((X:Z),(Y:T)) satisfying x=X/Z, y=Y/T
//   preComputedGroupElement: (y+x,y-x,2dxy)

type projectiveGroupElement struct {
	X, Y, Z fieldElement
}

type extendedGroupElement struct {
	X, Y, Z, T fieldElement
}

type completedGroupElement struct {
	X, Y, Z, T fieldElement
}

type preComputedGroupElement struct {
	yPlusX, yMinusX, xy2d fieldElement
}

type cachedGroupElement struct {
	yPlusX, yMinusX, Z, T2d fieldElement
}

func (p *projectiveGroupElement) Zero() {
	feZero(&p.X)
	feOne(&p.Y)
	feOne(&p.Z)
}

func (p *projectiveGroupElement) Double(r *completedGroupElement) {
	var t0 fieldElement

	feSquare(&r.X, &p.X)
	feSquare(&r.Z, &p.Y)
	feSquare2(&r.T, &p.Z)
	feAdd(&r.Y, &p.X, &p.Y)
	feSquare(&t0, &r.Y)
	feAdd(&r.Y, &r.Z, &r.X)
	feSub(&r.Z, &r.Z, &r.X)
	feSub(&r.X, &t0, &r.Y)
	feSub(&r.T, &r.T, &r.Z)
}

func (p *projectiveGroupElement) ToBytes(s *[32]byte) {
	var recip, x, y fieldElement

	feInvert(&recip, &p.Z)
	feMul(&x, &p.X, &recip)
	feMul(&y, &p.Y, &recip)
	feToBytes(s, &y)
	s[31] ^= feIsNegative(&x) << 7
}

func (p *extendedGroupElement) Zero() {
	feZero(&p.X)
	feOne(&p.Y)
	feOne(&p.Z)
	feZero(&p.T)
}

func (p *extendedGroupElement) Neg(s *extendedGroupElement) {
	feNeg(&p.X, &s.X)
	feCopy(&p.Y, &s.Y)
	feCopy(&p.Z, &s.Z)
	feNeg(&p.T, &s.T)
}

func (p *extendedGroupElement) Double(r *completedGroupElement) {
	var q projectiveGroupElement
	p.ToProjective(&q)
	q.Double(r)
}

func (p *extendedGroupElement) ToCached(r *cachedGroupElement) {
	feAdd(&r.yPlusX, &p.Y, &p.X)
	feSub(&r.yMinusX, &p.Y, &p.X)
	feCopy(&r.Z, &p.Z)
	feMul(&r.T2d, &p.T, &d2)
}

func (p *extendedGroupElement) ToProjective(r *projectiveGroupElement) {
	feCopy(&r.X, &p.X)
	feCopy(&r.Y, &p.Y)
	feCopy(&r.Z, &p.Z)
}

func (p *extendedGroupElement) ToBytes(s *[32]byte) {
	var recip, x, y fieldElement

	feInvert(&recip, &p.Z)
	feMul(&x, &p.X, &recip)
	feMul(&y, &p.Y, &recip)
	feToBytes(s, &y)
	s[31] ^= feIsNegative(&x) << 7
}

func (p *extendedGroupElement) FromBytes(s []byte) bool {
	var u, v, v3, vxx, check fieldElement

	if len(s) != 32 {
		return false
	}
	feFromBytes(&p.Y, s)
	feOne(&p.Z)
	feSquare(&u, &p.Y)
	feMul(&v, &u, &d)
	feSub(&u, &u, &p.Z) // y = y^2-1
	feAdd(&v, &v, &p.Z) // v = dy^2+1

	feSquare(&v3, &v)
	feMul(&v3, &v3, &v) // v3 = v^3
	feSquare(&p.X, &v3)
	feMul(&p.X, &p.X, &v)
	feMul(&p.X, &p.X, &u) // x = uv^7

	fePow22523(&p.X, &p.X) // x = (uv^7)^((q-5)/8)
	feMul(&p.X, &p.X, &v3)
	feMul(&p.X, &p.X, &u) // x = uv^3(uv^7)^((q-5)/8)

	feSquare(&vxx, &p.X)
	feMul(&vxx, &vxx, &v)
	feSub(&check, &vxx, &u) // vx^2-u
	if feIsNonZero(&check) == 1 {
		feAdd(&check, &vxx, &u) // vx^2+u
		if feIsNonZero(&check) == 1 {
			return false
		}
		feMul(&p.X, &p.X, &sqrtM1)
	}

	if feIsNegative(&p.X) != (s[31] >> 7) {
		feNeg(&p.X, &p.X)
	}

	feMul(&p.T, &p.X, &p.Y)
	return true
}

func (p *extendedGroupElement) String() string {
	return "extendedGroupElement{\n\t" +
		p.X.String() + ",\n\t" +
		p.Y.String() + ",\n\t" +
		p.Z.String() + ",\n\t" +
		p.T.String() + ",\n}"
}

// completedGroupElement methods

func (c *completedGroupElement) ToProjective(r *projectiveGroupElement) {
	feMul(&r.X, &c.X, &c.T)
	feMul(&r.Y, &c.Y, &c.Z)
	feMul(&r.Z, &c.Z, &c.T)
}

func (c *completedGroupElement) ToExtended(r *extendedGroupElement) {
	feMul(&r.X, &c.X, &c.T)
	feMul(&r.Y, &c.Y, &c.Z)
	feMul(&r.Z, &c.Z, &c.T)
	feMul(&r.T, &c.X, &c.Y)
}

func (p *preComputedGroupElement) Zero() {
	feOne(&p.yPlusX)
	feOne(&p.yMinusX)
	feZero(&p.xy2d)
}

func (c *completedGroupElement) Add(p *extendedGroupElement, q *cachedGroupElement) {
	var t0 fieldElement

	feAdd(&c.X, &p.Y, &p.X)
	feSub(&c.Y, &p.Y, &p.X)
	feMul(&c.Z, &c.X, &q.yPlusX)
	feMul(&c.Y, &c.Y, &q.yMinusX)
	feMul(&c.T, &q.T2d, &p.T)
	feMul(&c.X, &p.Z, &q.Z)
	feAdd(&t0, &c.X, &c.X)
	feSub(&c.X, &c.Z, &c.Y)
	feAdd(&c.Y, &c.Z, &c.Y)
	feAdd(&c.Z, &t0, &c.T)
	feSub(&c.T, &t0, &c.T)
}

func (c *completedGroupElement) Sub(p *extendedGroupElement, q *cachedGroupElement) {
	var t0 fieldElement

	feAdd(&c.X, &p.Y, &p.X)
	feSub(&c.Y, &p.Y, &p.X)
	feMul(&c.Z, &c.X, &q.yMinusX)
	feMul(&c.Y, &c.Y, &q.yPlusX)
	feMul(&c.T, &q.T2d, &p.T)
	feMul(&c.X, &p.Z, &q.Z)
	feAdd(&t0, &c.X, &c.X)
	feSub(&c.X, &c.Z, &c.Y)
	feAdd(&c.Y, &c.Z, &c.Y)
	feSub(&c.Z, &t0, &c.T)
	feAdd(&c.T, &t0, &c.T)
}

func (c *completedGroupElement) MixedAdd(p *extendedGroupElement, q *preComputedGroupElement) {
	var t0 fieldElement

	feAdd(&c.X, &p.Y, &p.X)
	feSub(&c.Y, &p.Y, &p.X)
	feMul(&c.Z, &c.X, &q.yPlusX)
	feMul(&c.Y, &c.Y, &q.yMinusX)
	feMul(&c.T, &q.xy2d, &p.T)
	feAdd(&t0, &p.Z, &p.Z)
	feSub(&c.X, &c.Z, &c.Y)
	feAdd(&c.Y, &c.Z, &c.Y)
	feAdd(&c.Z, &t0, &c.T)
	feSub(&c.T, &t0, &c.T)
}

func (c *completedGroupElement) MixedSub(p *extendedGroupElement, q *preComputedGroupElement) {
	var t0 fieldElement

	feAdd(&c.X, &p.Y, &p.X)
	feSub(&c.Y, &p.Y, &p.X)
	feMul(&c.Z, &c.X, &q.yMinusX)
	feMul(&c.Y, &c.Y, &q.yPlusX)
	feMul(&c.T, &q.xy2d, &p.T)
	feAdd(&t0, &p.Z, &p.Z)
	feSub(&c.X, &c.Z, &c.Y)
	feAdd(&c.Y, &c.Z, &c.Y)
	feSub(&c.Z, &t0, &c.T)
	feAdd(&c.T, &t0, &c.T)
}

// preComputedGroupElement methods

// Set to u conditionally based on b
func (p *preComputedGroupElement) CMove(u *preComputedGroupElement, b int32) {
	feCMove(&p.yPlusX, &u.yPlusX, b)
	feCMove(&p.yMinusX, &u.yMinusX, b)
	feCMove(&p.xy2d, &u.xy2d, b)
}

// Set to negative of t
func (p *preComputedGroupElement) Neg(t *preComputedGroupElement) {
	feCopy(&p.yPlusX, &t.yMinusX)
	feCopy(&p.yMinusX, &t.yPlusX)
	feNeg(&p.xy2d, &t.xy2d)
}

// cachedGroupElement methods

func (r *cachedGroupElement) Zero() {
	feOne(&r.yPlusX)
	feOne(&r.yMinusX)
	feOne(&r.Z)
	feZero(&r.T2d)
}

// Set to u conditionally based on b
func (r *cachedGroupElement) CMove(u *cachedGroupElement, b int32) {
	feCMove(&r.yPlusX, &u.yPlusX, b)
	feCMove(&r.yMinusX, &u.yMinusX, b)
	feCMove(&r.Z, &u.Z, b)
	feCMove(&r.T2d, &u.T2d, b)
}

// Set to negative of t
func (r *cachedGroupElement) Neg(t *cachedGroupElement) {
	feCopy(&r.yPlusX, &t.yMinusX)
	feCopy(&r.yMinusX, &t.yPlusX)
	feCopy(&r.Z, &t.Z)
	feNeg(&r.T2d, &t.T2d)
}

// Expand the 32-byte (256-bit) exponent in slice a into
// a sequence of 256 multipliers, one per exponent bit position.
// Clumps nearby 1 bits into multi-bit multipliers to reduce
// the total number of add/sub operations in a point multiply;
// each multiplier is either zero or an odd number between -15 and 15.
// Assumes the target array r has been preinitialized with zeros
// in case the input slice a is less than 32 bytes.
func slide(r *[256]int8, a *[32]byte) {

	// Explode the exponent a into a little-endian array, one bit per byte
	for i := range a {
		ai := int8(a[i])
		for j := 0; j < 8; j++ {
			r[i*8+j] = ai & 1
			ai >>= 1
		}
	}

	// Go through and clump sequences of 1-bits together wherever possible,
	// while keeping r[i] in the range -15 through 15.
	// Note that each nonzero r[i] in the result will always be odd,
	// because clumping is triggered by the first, least-significant,
	// 1-bit encountered in a clump, and that first bit always remains 1.
	for i := range r {
		if r[i] != 0 {
			for b := 1; b <= 6 && i+b < 256; b++ {
				if r[i+b] != 0 {
					if r[i]+(r[i+b]<<uint(b)) <= 15 {
						r[i] += r[i+b] << uint(b)
						r[i+b] = 0
					} else if r[i]-(r[i+b]<<uint(b)) >= -15 {
						r[i] -= r[i+b] << uint(b)
						for k := i + b; k < 256; k++ {
							if r[k] == 0 {
								r[k] = 1
								break
							}
							r[k] = 0
						}
					} else {
						break
					}
				}
			}
		}
	}
}

// equal returns 1 if b == c and 0 otherwise.
func equal(b, c int32) int32 {
	x := uint32(b ^ c)
	x--
	return int32(x >> 31)
}

// negative returns 1 if b < 0 and 0 otherwise.
func negative(b int32) int32 {
	return (b >> 31) & 1
}

func selectPreComputed(t *preComputedGroupElement, pos int32, b int32) {
	var minusT preComputedGroupElement
	bNegative := negative(b)
	bAbs := b - (((-bNegative) & b) << 1)

	t.Zero()
	for i := int32(0); i < 8; i++ {
		t.CMove(&base[pos][i], equal(bAbs, i+1))
	}
	minusT.Neg(t)
	t.CMove(&minusT, bNegative)
}

// geScalarMultBase computes h = a*B, where
//   a = a[0]+256*a[1]+...+256^31 a[31]
//   B is the Ed25519 base point (x,4/5) with x positive.
//
// Preconditions:
//   a[31] <= 127
func geScalarMultBase(h *extendedGroupElement, a *[32]byte) {
	var e [64]int8

	for i, v := range a {
		e[2*i] = int8(v & 15)
		e[2*i+1] = int8((v >> 4) & 15)
	}

	// each e[i] is between 0 and 15 and e[63] is between 0 and 7.

	carry := int8(0)
	for i := 0; i < 63; i++ {
		e[i] += carry
		carry = (e[i] + 8) >> 4
		e[i] -= carry << 4
	}
	e[63] += carry
	// each e[i] is between -8 and 8.

	h.Zero()
	var t preComputedGroupElement
	var r completedGroupElement
	for i := int32(1); i < 64; i += 2 {
		selectPreComputed(&t, i/2, int32(e[i]))
		r.MixedAdd(h, &t)
		r.ToExtended(h)
	}

	var s projectiveGroupElement

	h.Double(&r)
	r.ToProjective(&s)
	s.Double(&r)
	r.ToProjective(&s)
	s.Double(&r)
	r.ToProjective(&s)
	s.Double(&r)
	r.ToExtended(h)

	for i := int32(0); i < 64; i += 2 {
		selectPreComputed(&t, i/2, int32(e[i]))
		r.MixedAdd(h, &t)
		r.ToExtended(h)
	}
}

func selectCached(c *cachedGroupElement, Ai *[8]cachedGroupElement, b int32) {
	bNegative := negative(b)
	bAbs := b - (((-bNegative) & b) << 1)

	// in constant-time pick cached multiplier for exponent 0 through 8
	c.Zero()
	for i := int32(0); i < 8; i++ {
		c.CMove(&Ai[i], equal(bAbs, i+1))
	}

	// in constant-time compute negated version, conditionally use it
	var minusC cachedGroupElement
	minusC.Neg(c)
	c.CMove(&minusC, bNegative)
}

// geScalarMult computes h = a*B, where
//   a = a[0]+256*a[1]+...+256^31 a[31]
//   B is the Ed25519 base point (x,4/5) with x positive.
//
// Preconditions:
//   a[31] <= 127
func geScalarMult(h *extendedGroupElement, a *[32]byte,
	A *extendedGroupElement) {

	var t completedGroupElement
	var u extendedGroupElement
	var r projectiveGroupElement
	var c cachedGroupElement
	var i int

	// Break the exponent into 4-bit nybbles.
	var e [64]int8
	for i, v := range a {
		e[2*i] = int8(v & 15)
		e[2*i+1] = int8((v >> 4) & 15)
	}
	// each e[i] is between 0 and 15 and e[63] is between 0 and 7.

	carry := int8(0)
	for i := 0; i < 63; i++ {
		e[i] += carry
		carry = (e[i] + 8) >> 4
		e[i] -= carry << 4
	}
	e[63] += carry
	// each e[i] is between -8 and 8.

	// compute cached array of multiples of A from 1A through 8A
	var Ai [8]cachedGroupElement // A,1A,2A,3A,4A,5A,6A,7A
	A.ToCached(&Ai[0])
	for i := 0; i < 7; i++ {
		t.Add(A, &Ai[i])
		t.ToExtended(&u)
		u.ToCached(&Ai[i+1])
	}

	// special case for exponent nybble i == 63
	u.Zero()
	selectCached(&c, &Ai, int32(e[63]))
	t.Add(&u, &c)

	for i = 62; i >= 0; i-- {

		// t <<= 4
		t.ToProjective(&r)
		r.Double(&t)
		t.ToProjective(&r)
		r.Double(&t)
		t.ToProjective(&r)
		r.Double(&t)
		t.ToProjective(&r)
		r.Double(&t)

		// Add next nybble
		t.ToExtended(&u)
		selectCached(&c, &Ai, int32(e[i]))
		t.Add(&u, &c)
	}

	t.ToExtended(h)
}
