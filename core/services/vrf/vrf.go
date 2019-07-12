package vrf

// vrf provides a cryptographically secure pseudo-random number generator.
// Numbers are deterministically generated from a seed and a secret key, and are
// statistically indistinguishable from uniform sampling from {0, ..., 2**256},
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
// The public key associated with it can be calculated from it by XXX
//
// To generate random output from a big.Int seed, pass sk and the seed to
// GenerateProof, and use the Output field of the returned Proof object.
//
// To verify a Proof object p, run p.Verify(), or pass its fields to the
// corresponding arguments of isValidVRFOutput on the VRF solidity contract, to
// verify it on-chain.

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/core/utils"
	"go.dedis.ch/kyber"
)

func bigFromHex(s string) *big.Int {
	n, ok := new(big.Int).SetString(s, 16)
	if !ok {
		panic(fmt.Errorf(`failed to convert "%s" as hex to big.Int`, s))
	}
	return n
}

// P is number of elements in the Galois field over which secp256k1 is defined
var P = bigFromHex(
	"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F")

// Order is the number of rational points on the curve in GF(P) (group size)
var Order = bigFromHex(
	"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141")

// Compensate for awkward big.Int API.
var bi = big.NewInt
var zero, one, two, three, four, seven = bi(0), bi(1), bi(2), bi(3), bi(4), bi(7)

func i() *big.Int                                    { return new(big.Int) }
func add(addend1, addend2 *big.Int) *big.Int         { return i().Add(addend1, addend2) }
func div(dividend, divisor *big.Int) *big.Int        { return i().Div(dividend, divisor) }
func equal(left, right *big.Int) bool                { return left.Cmp(right) == 0 }
func exp(base, exponent, modulus *big.Int) *big.Int  { return i().Exp(base, exponent, modulus) }
func lsh(num *big.Int, bits uint) *big.Int           { return i().Lsh(num, bits) }
func mul(multiplicand, multiplier *big.Int) *big.Int { return i().Mul(multiplicand, multiplier) }
func mod(dividend, divisor *big.Int) *big.Int        { return i().Mod(dividend, divisor) }
func sub(minuend, subtrahend *big.Int) *big.Int      { return i().Sub(minuend, subtrahend) }

var (
	// (P-1)/2: Half Fermat's Little Theorem exponent
	eulersCriterionPower = div(sub(P, one), two)
	// (P+1)/4: As long as P%4==3 and n=x^2 in GF(P), n^((P+1)/4)=±x
	sqrtPower = div(add(P, one), four)
)

// IsSquare returns true iff x = y^2 for some y in GF(p)
func IsSquare(x *big.Int) bool {
	return equal(one, exp(x, eulersCriterionPower, P))
}

// SquareRoot returns a s.t. a^2=x. Assumes x is a square
func SquareRoot(x *big.Int) *big.Int {
	return exp(x, sqrtPower, P)
}

// YSquared returns x^3+7 mod P
func YSquared(x *big.Int) *big.Int {
	return mod(add(exp(x, three, P), seven), P)
}

// IsCurveXOrdinate returns true iff there is y s.t. y^2=x^3+7
func IsCurveXOrdinate(x *big.Int) bool {
	return IsSquare(YSquared(x))
}

func packUint256s(xs ...*big.Int) ([]byte, error) {
	mem := []byte{}
	for _, x := range xs {
		word, err := utils.EVMWordBigInt(x)
		if err != nil {
			return []byte{}, errors.Wrap(err, "vrf.packUint256s#EVMWordBigInt")
		}
		mem = append(mem, word...)
	}
	return mem, nil
}

type curveT *secp256k1.Secp256k1

var curve = secp256k1.Secp256k1{}
var rcurve = &curve

// Generator is a specific generator of the curve group. Any non-zero point will
// do, since the group order is prime. But one must be specified as part of the
// protocol.
var Generator = rcurve.Point().Base()

// CoordsFromPoint returns the (x, y) coordinates of p
func CoordsFromPoint(p kyber.Point) (*big.Int, *big.Int) {
	return secp256k1.Coordinates(p)
}

// HashUint256s returns a uint256 representing the hash of the concatenated byte
// representations of the inputs
func HashUint256s(xs ...*big.Int) (*big.Int, error) {
	packed, err := packUint256s(xs...)
	if err != nil {
		return &big.Int{}, err
	}
	hash, err := utils.Keccak256(packed)
	if err != nil {
		return &big.Int{}, errors.Wrap(err, "vrf.HashUint256s#Keccak256")
	}
	return i().SetBytes(hash), nil
}

func asUint256(x *big.Int) []byte {
	if x.BitLen() > 256 {
		panic("vrf.asUint256: too big to marshal to uint256")
	}
	return common.LeftPadBytes(x.Bytes(), 32)
}

var numWords = lsh(two, 256)
var mask = sub(numWords, one)

// ZqHash hashes xs uniformly into {0, ..., q-1}. q must be 256 bits long, and
// msg is assumed to already be a 256-bit hash
func ZqHash(q *big.Int, msg []byte) (*big.Int, error) {
	if q.BitLen() != 256 || len(msg) > 256 {
		panic(fmt.Errorf(
			"will only work for moduli 256 bits long, need %v",
			q.BitLen()))
	}
	rv := i().SetBytes(msg)
	// Hash recursively until rv < q. P(success per iteration) >= 0.5, so
	// number of extra hashes is geometrically distributed, with mean < 1.
	for rv.Cmp(q) != -1 {
		hash, err := utils.Keccak256(asUint256(rv))
		if err != nil {
			return nil, errors.Wrap(err, "vrf.ZqHash#Keccak256.loop")
		}
		rv.SetBytes(hash)
	}
	return rv, nil
}

// HashToCurve is a one-way hash function onto the curve
func HashToCurve(p kyber.Point, input *big.Int) (kyber.Point, error) {
	if !(secp256k1.ValidPublicKey(p) && input.BitLen() <= 32) {
		return nil, fmt.Errorf("bad input to vrf.HashToCurve")
	}
	iHash, err := utils.Keccak256(
		append(secp256k1.LongMarshal(p), asUint256(input)...))
	if err != nil {
		panic(errors.Wrap(err, "while attempting initial hash"))
	}
	x, err := ZqHash(P, iHash)
	if err != nil {
		return nil, errors.Wrap(err, "vrf.HashToCurve#ZqHash")
	}
	count := 0
	for !IsCurveXOrdinate(x) { // Hash recursively until x^3+7 is a square
		count += 1
		if count >= 10 {
			panic("done")
		}
		nHash, err := utils.Keccak256(asUint256(x))
		if err != nil {
			panic(errors.Wrap(err, "while attempting to rehash x"))
		}
		nx, err := ZqHash(P, nHash)
		if err != nil {
			panic(err)
		}
		x.Set(nx)
	}
	rv := secp256k1.SetCoordinates(x, SquareRoot(YSquared(x)))
	// Two possible y ordinates for this x ordinate; pick one "randomly"
	nhash, err := HashUint256s(x, input) // nhash is the random value
	if err != nil {
		return nil, errors.Wrap(err, "vrf.HashToCurve#HashUint256s")
	}
	if i().Mod(nhash, two).Cmp(zero) == 0 { // Negate response if nhash even
		rv = rv.Neg(rv)
	}
	return rv, nil
}

// ScalarFromCurvePoints returns a hash for the curve points. Corresponds to the
// hash computed in Curve.sol#scalarFromCurve
func ScalarFromCurvePoints(
	hash, pk, gamma kyber.Point, uWitness [20]byte, v kyber.Point) *big.Int {
	if !(secp256k1.ValidPublicKey(hash) && secp256k1.ValidPublicKey(pk) &&
		secp256k1.ValidPublicKey(gamma) && secp256k1.ValidPublicKey(v)) {
		panic("bad arguments to vrf.ScalarFromCurvePoints")
	}
	msg := secp256k1.LongMarshal(hash)
	msg = append(msg, secp256k1.LongMarshal(pk)...)
	msg = append(msg, secp256k1.LongMarshal(gamma)...)
	msg = append(msg, secp256k1.LongMarshal(v)...)
	msg = append(msg, uWitness[:]...)
	preHash, err := utils.Keccak256(msg)
	if err != nil {
		panic(err)
	}
	h, err := ZqHash(Order, preHash)
	if err != nil {
		panic(err)
	}
	return h
}

// linearComination returns c*p1+s*p2
func linearCombination(c *big.Int, p1 kyber.Point,
	s *big.Int, p2 kyber.Point) kyber.Point {
	return rcurve.Point().Add(
		rcurve.Point().Mul(secp256k1.IntToScalar(c), p1),
		rcurve.Point().Mul(secp256k1.IntToScalar(s), p2))
}

// Proof represents a proof that Gamma was constructed from the Seed
// according to the process mandated by the PublicKey.
//
// N.B.: The kyber.Point fields must contain secp256k1.secp256k1Point values
type Proof struct {
	PublicKey, Gamma   kyber.Point
	C, S, Seed, Output *big.Int
}

func (p *Proof) WellFormed() bool {
	return (secp256k1.ValidPublicKey(p.PublicKey) && secp256k1.ValidPublicKey(p.Gamma) &&
		secp256k1.RepresentsScalar(p.C) && secp256k1.RepresentsScalar(p.S) &&
		p.Output.BitLen() <= 256)
}

// VerifyProof is true iff gamma was generated in the mandated way from the
// given publicKey and seed
func (proof *Proof) Verify() (bool, error) {
	if !proof.WellFormed() {
		return false, fmt.Errorf("badly-formatted proof")
	}
	h, err := HashToCurve(proof.PublicKey, proof.Seed)
	if err != nil {
		return false, err
	}
	// publicKey = secretKey*Generator. See GenerateProof for u, v, m, s
	// c*secretKey*Generator + (m - c*secretKey)*Generator = m*Generator = u
	uPrime := linearCombination(proof.C, proof.PublicKey, proof.S, Generator)
	// c*secretKey*h + (m - c*secretKey)*h = m*h = v
	vPrime := linearCombination(proof.C, proof.Gamma, proof.S, h)
	uWitness, err := secp256k1.EthereumAddress(uPrime)
	if err != nil {
		return false, errors.Wrap(err, "vrf.VerifyProof#EthereumAddress")
	}
	cPrime := ScalarFromCurvePoints(h, proof.PublicKey, proof.Gamma, uWitness, vPrime)
	output, err := utils.Keccak256(secp256k1.LongMarshal(proof.Gamma))
	if err != nil {
		panic(errors.Wrap(err, "while hashing to compute proof output"))
	}
	return (proof.C.Cmp(cPrime) == 0) &&
			(proof.Output.Cmp(i().SetBytes(output)) == 0),
		nil
}

// GenerateProof returns gamma, plus proof that gamma was constructed from seed
// as mandated from the given secretKey, with public key secretKey*Generator
// Proof is constructed using nonce as the ephemeral key. If provided, it must
// be treated as key material (cryptographically-securely randomly generated,
// kept confidential or just forgotten.) If it's nil, it will be generated here.
func GenerateProof(secretKey, seed, nonce *big.Int) (*Proof, error) {
	if !(secp256k1.RepresentsScalar(secretKey) && seed.BitLen() <= 256) {
		return nil, fmt.Errorf("badly-formatted key or seed")
	}
	publicKey := rcurve.Point().Mul(secp256k1.IntToScalar(secretKey), nil)
	h, err := HashToCurve(publicKey, seed)
	if err != nil {
		return &Proof{}, errors.Wrap(err, "vrf.makeProof#HashToCurve")
	}
	gamma := rcurve.Point().Mul(secp256k1.IntToScalar(secretKey), h)
	if nonce == nil {
		nonce, err = rand.Int(rand.Reader, Order)
		if err != nil {
			return &Proof{}, errors.Wrap(err, "vrf.makeProof#rand.Int")
		}
	}
	sm := secp256k1.IntToScalar(nonce)
	u := rcurve.Point().Mul(sm, Generator)
	uWitness, err := secp256k1.EthereumAddress(u)
	if err != nil {
		panic(errors.Wrap(err, "while computing Ethereum Address for proof"))
	}
	v := rcurve.Point().Mul(sm, h)
	c := ScalarFromCurvePoints(h, publicKey, gamma, uWitness, v)
	// s = (m - c*secretKey) % Order
	s := mod(sub(nonce, mul(c, secretKey)), Order)
	outputHash, err := utils.Keccak256(secp256k1.LongMarshal(gamma))
	if err != nil {
		panic("failed to hash gamma")
	}
	rv := Proof{
		PublicKey: publicKey,
		Gamma:     gamma,
		C:         c,
		S:         s,
		Seed:      seed,
		Output:    i().SetBytes(outputHash),
	}
	valid, err := rv.Verify()
	if !valid || err != nil {
		panic("constructed invalid proof")
	}
	return &rv, nil
}
