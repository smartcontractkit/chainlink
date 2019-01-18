package vrf

// vrf provides a cryptographically secure pseudo-random number generator.
// Numbers are deterministically generated from a seed and a secret key, and are
// statistically indistinguishable from uniform sampling from {0, ..., 2**256-1},
// to observers who don't know the key. But each number comes with a proof that
// it was generated according to the procedure mandated by a public key
// associated with that private key.
//
// See VRF.sol for design notes.
//
// Usage
// -----
//
// A secret key sk should be securely sampled uniformly from {0, ..., Order}.
// The public key associated with it can be calculated from it by
// bn256.ScalarMult(Generator,sk).
//
// To generate random output from a big.Int seed, pass sk and the seed to
// GenerateProof, and use the Output field of the returned Proof object.
//
// To verify a Proof object p, run p.Verify(), or pass its fields to the
// corresponding arguments of isValidVRFOutput on the VRF solidity contract, to
// verify it on-chain.

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math/big"

	// NB: this curve is actually alt-bn128, not bn256!
	"github.com/ethereum/go-ethereum/crypto/bn256"

	"github.com/smartcontractkit/chainlink/utils"
)

func bigFromBase10(s string) *big.Int {
	n, _ := new(big.Int).SetString(s, 10)
	return n
}

// The following variables are not exported, despite being capitalized, from
// github.com/ethereum/go-ethereum/crypto/bn256/cloudfare/constants.go

// These values don't match the published BN256 values, because the bn256
// package does not implement the BN256 curve, but the Alt-BN128 curve.

// P is the number of elements in the Galois field over which Alt-BN128 is defined
var P = bigFromBase10(
	"21888242871839275222246405745257275088696311157297823662689037894645226208583")

// Order is the number of rational points on the curve in GF(P) (group size)
var Order = bigFromBase10(
	"21888242871839275222246405745257275088548364400416034343698204186575808495617")

// Compensate for awkward big.Int API.
var bi = big.NewInt
var zero, one, two, three, four = bi(0), bi(1), bi(2), bi(3), bi(4)

func i() *big.Int                                    { return new(big.Int) }
func add(addend1, addend2 *big.Int) *big.Int         { return i().Add(addend1, addend2) }
func div(dividend, divisor *big.Int) *big.Int        { return i().Div(dividend, divisor) }
func equal(left, right *big.Int) bool                { return left.Cmp(right) == 0 }
func exp(base, exponent, modulus *big.Int) *big.Int  { return i().Exp(base, exponent, modulus) }
func lsh(num *big.Int, bits uint) *big.Int           { return i().Lsh(num, bits) }
func mul(multiplicand, multiplier *big.Int) *big.Int { return i().Mul(multiplicand, multiplier) }
func mod(dividend, divisor *big.Int) *big.Int        { return i().Mod(dividend, divisor) }
func sub(minuend, subtrahend *big.Int) *big.Int      { return i().Sub(minuend, subtrahend) }

// Sadly, the hidden interfaces in bn256 require us to duplicate some
// functionality.

var (
	// (P-1)/2: Half Fermat's Little Theorem exponent
	eulersCriterionPower = div(sub(P, one), two)
	// (P+1)/4: As long as P%4==3 and n=x^2 in GF(P), n^((P+1)/4)=±x
	sqrtPower = div(add(P, one), four)
)

// isSquare returns true iff x = y^2 for some y in GF(p)
func isSquare(x *big.Int) bool {
	return equal(one, exp(x, eulersCriterionPower, P))
}

// squareRoot returns a s.t. a^2=x. Assumes x is a square
func squareRoot(x *big.Int) *big.Int {
	return exp(x, sqrtPower, P)
}

// ySquared returns x^3+3 mod P
func ySquared(x *big.Int) *big.Int {
	return mod(add(exp(x, three, P), three), P)
}

// IsCurveXOrdinate returns true iff there is y s.t. y^2=x^3+3
func isCurveXOrdinate(x *big.Int) bool {
	return isSquare(ySquared(x))
}

func packUint256s(xs ...*big.Int) ([]byte, error) {
	mem := bytes.Buffer{}
	for _, x := range xs {
		word, err := utils.EVMWordBigInt(x)
		if err != nil {
			return []byte{}, err
		}
		n, err := mem.Write(word)
		if n != 32 {
			return []byte{}, fmt.Errorf(
				"Failed to write as uint256: %v", x)
		}
		if err != nil {
			return []byte{}, err
		}
	}
	if mem.Len() != len(xs)*32 {
		panic(fmt.Errorf("Package of %v uint256s unexpected length, %v",
			len(xs), mem.Len()))
	}
	return mem.Bytes(), nil
}

// newCurvePoint returns the bn256.G1 point corresponding to (x, y)
func newCurvePoint(x, y *big.Int) (*bn256.G1, error) {
	p := new(bn256.G1)
	packed, err := packUint256s(x, y)
	if err != nil {
		return &bn256.G1{}, err
	}
	b, err := p.Unmarshal(packed) // Unmarshal's comment lies about return type
	if len(b) != 0 {
		panic(fmt.Errorf(
			"Did not consume all of packed (%v, %v)", x, y))
	}
	if err != nil {
		return &bn256.G1{}, err
	}
	return p, nil
}

// Generator is a specific generator of the curve group. Any non-zero point will
// do, since the group order is prime. But one must be specified as part of the
// protocol.
var Generator, _ = newCurvePoint(one, two)

// CoordsFromPoint returns the (x, y) coordinates of p
func CoordsFromPoint(p *bn256.G1) (*big.Int, *big.Int) {
	b := p.Marshal()
	if len(b) != 64 {
		panic(fmt.Errorf("did not get 512 bits from %v", p))
	}
	return i().SetBytes(b[:32]), i().SetBytes(b[32:])
}

// hashUint256s returns a uint256 representing the hash of the concatenated byte
// representations of the inputs
func hashUint256s(xs ...*big.Int) (*big.Int, error) {
	packed, err := packUint256s(xs...)
	if err != nil {
		return &big.Int{}, err
	}
	hash, err := utils.Keccak256(packed)
	if err != nil {
		return &big.Int{}, err
	}
	return i().SetBytes(hash), nil
}

// maskHash returns hashUint256s(xs...) & mask
func maskHash(mask *big.Int, xs ...*big.Int) (*big.Int, error) {
	x, err := hashUint256s(xs...)
	if err != nil {
		return &big.Int{}, err
	}
	return x.And(x, mask), nil
}

// zqHash hashes xs uniformly into {0, ..., q-1}
func zqHash(q *big.Int, xs ...*big.Int) (*big.Int, error) {
	if len(xs) < 1 {
		panic("can't take hash of empty list. You might have forgotten argument q")
	}
	if q.BitLen() > 256 {
		panic(fmt.Errorf(
			"will only generate 256 bits of entropy, need %v",
			q.BitLen()))
	}
	// Bits which can be used in representation of a number less than q.
	// 2^(q.BitLen)-1
	orderMask := sub(lsh(one, uint(q.BitLen())), one)
	rv, err := maskHash(orderMask, xs...)
	if err != nil {
		return &big.Int{}, err
	}
	// Hash recursively until rv < q. P(success per iteration) >= 0.5, so
	// number of extra hashes is geometrically distributed, with mean < 1.
	for rv.Cmp(q) != -1 {
		nrv, err := maskHash(orderMask, rv)
		if err != nil {
			return &big.Int{}, err
		}
		rv.Set(nrv)
	}
	return rv, nil
}

// hashToCurve is a one-way hash function onto the curve
func hashToCurve(px, py, input *big.Int) (*bn256.G1, error) {
	x, err := zqHash(P, px, py, input)
	if err != nil {
		return &bn256.G1{}, err
	}
	for !isCurveXOrdinate(x) { // Hash recursively until x^3+3 is a square
		nx, err := zqHash(P, x)
		if err != nil {
			return &bn256.G1{}, err
		}
		x.Set(nx)
	}
	return newCurvePoint(x, squareRoot(ySquared(x)))
}

// scalarFromCurve returns a hash for the curve points. Corresponds to the hash
// computed in Curve.sol#scalarFromCurve
func scalarFromCurve(ps ...*bn256.G1) (*big.Int, error) {
	coordinates := make([]*big.Int, (len(ps)+1)*2)
	gx, gy := CoordsFromPoint(Generator)
	coordinates[0] = gx
	coordinates[1] = gy
	for ordidx, p := range ps {
		x, y := CoordsFromPoint(p)
		coordinates[2*ordidx+2] = x
		coordinates[2*ordidx+3] = y
	}
	return zqHash(Order, coordinates...)
}

// linearComination returns c*p1+s*p2
func linearComination(c *big.Int, p1 *bn256.G1, s *big.Int, p2 *bn256.G1) *bn256.G1 {
	return new(bn256.G1).Add(
		new(bn256.G1).ScalarMult(p1, c),
		new(bn256.G1).ScalarMult(p2, s))
}

// Proof represents a proof that Gamma was constructed from the Seed
// according to the process mandated by the PublicKey
type Proof struct {
	PublicKey, Gamma   *bn256.G1
	C, S, Seed, Output *big.Int
}

// VerifyProof is true iff gamma was generated in the mandated way from the
// given publicKey and seed
func (proof *Proof) VerifyProof() (bool, error) {
	px, py := CoordsFromPoint(proof.PublicKey)
	h, err := hashToCurve(px, py, proof.Seed)
	if err != nil {
		return false, err
	}
	// publicKey = secretKey*Generator. See GenerateProof for u, v, m, s
	// c*secretKey*Generator + (m - c*secretKey)*Generator = m*Generator = u
	uPrime := linearComination(proof.C, proof.PublicKey, proof.S, Generator)
	// c*secretKey*h + (m - c*secretKey)*h = m*h = v
	vPrime := linearComination(proof.C, proof.Gamma, proof.S, h)
	cPrime, _ := scalarFromCurve(
		h, proof.PublicKey, proof.Gamma, uPrime, vPrime)
	if err != nil {
		return false, err
	}
	output, err := utils.Keccak256(proof.Gamma.Marshal())
	if err != nil {
		return false, err
	}
	return (proof.C.Cmp(cPrime) == 0) &&
			(proof.Output.Cmp(i().SetBytes(output)) == 0),
		nil
}

// makeProof proof generates the actual proof, modulo the actual random output
func makeProof(secretKey, seed *big.Int) (*Proof, error) {
	publicKey := new(bn256.G1).ScalarMult(Generator, secretKey)
	px, py := CoordsFromPoint(publicKey)
	h, err := hashToCurve(px, py, seed)
	if err != nil {
		return &Proof{}, err
	}
	gamma := new(bn256.G1).ScalarMult(h, secretKey)
	m, err := rand.Int(rand.Reader, Order)
	if err != nil {
		return &Proof{}, err
	}
	u := new(bn256.G1).ScalarMult(Generator, m)
	v := new(bn256.G1).ScalarMult(h, m)
	c, err := scalarFromCurve(h, publicKey, gamma, u, v)
	if err != nil {
		return &Proof{}, err
	}
	// s = (m - c*secretKey) % Order
	s := mod(sub(m, mul(c, secretKey)), Order)
	return &Proof{
		PublicKey: publicKey,
		Gamma:     gamma,
		C:         c,
		S:         s,
		Seed:      seed,
	}, nil
}

// GenerateProof returns gamma, plus proof that gamma was constructed from seed
// as mandated from the given secretKey, with public key secretKey*Generator
func GenerateProof(secretKey, seed *big.Int) (*Proof, error) {
	proof, err := makeProof(secretKey, seed)
	if err != nil {
		return &Proof{}, err
	}
	output, err := utils.Keccak256(proof.Gamma.Marshal())
	if err != nil {
		return &Proof{}, err
	}
	proof.Output = i().SetBytes(output)
	return proof, nil
}
