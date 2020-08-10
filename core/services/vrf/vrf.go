// Package vrf provides a cryptographically secure pseudo-random number generator.

// Numbers are deterministically generated from seeds and a secret key, and are
// statistically indistinguishable from uniform sampling from {0,...,2**256-1},
// to computationally-bounded observers who know the seeds, don't know the key,
// and only see the generated numbers. But each number also comes with a proof
// that it was generated according to the procedure mandated by a public key
// associated with that secret key.
//
// See VRF.sol for design notes.
//
// Usage
// -----
//
// You should probably not be using this directly.
// chainlink/store/core/models/vrfkey.PrivateKey provides a simple, more
// misuse-resistant interface to the same functionality, via the CreateKey and
// MarshaledProof methods.
//
// Nonetheless, a secret key sk should be securely sampled uniformly from
// {0,...,Order-1}. Its public key can be calculated from it by
//
//   secp256k1.Secp256k1{}.Point().Mul(secretKey, Generator)
//
// To generate random output from a big.Int seed, pass sk and the seed to
// GenerateProof, and use the Output field of the returned Proof object.
//
// To verify a Proof object p, run p.Verify(); or to verify it on-chain pass
// p.MarshalForSolidityVerifier() to randomValueFromVRFProof on the VRF solidity
// contract.
package vrf

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/core/utils"

	"go.dedis.ch/kyber/v3"
)

func bigFromHex(s string) *big.Int {
	n, ok := new(big.Int).SetString(s, 16)
	if !ok {
		panic(fmt.Errorf(`failed to convert "%s" as hex to big.Int`, s))
	}
	return n
}

// FieldSize is number of elements in secp256k1's base field, i.e. GF(FieldSize)
var FieldSize = bigFromHex(
	"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F")

var bi = big.NewInt
var zero, one, two, three, four, seven = bi(0), bi(1), bi(2), bi(3), bi(4), bi(7)

// Compensate for awkward big.Int API. Can cause an extra allocation or two.
func i() *big.Int                                    { return new(big.Int) }
func add(addend1, addend2 *big.Int) *big.Int         { return i().Add(addend1, addend2) }
func div(dividend, divisor *big.Int) *big.Int        { return i().Div(dividend, divisor) }
func equal(left, right *big.Int) bool                { return left.Cmp(right) == 0 }
func exp(base, exponent, modulus *big.Int) *big.Int  { return i().Exp(base, exponent, modulus) }
func mul(multiplicand, multiplier *big.Int) *big.Int { return i().Mul(multiplicand, multiplier) }
func mod(dividend, divisor *big.Int) *big.Int        { return i().Mod(dividend, divisor) }
func sub(minuend, subtrahend *big.Int) *big.Int      { return i().Sub(minuend, subtrahend) }

var (
	// (fieldSize-1)/2: Half Fermat's Little Theorem exponent
	eulersCriterionPower = div(sub(FieldSize, one), two)
	// (fieldSize+1)/4: As long as P%4==3 and n=x^2 in GF(fieldSize), n^sqrtPower=±x
	sqrtPower = div(add(FieldSize, one), four)
)

// IsSquare returns true iff x = y^2 for some y in GF(p)
func IsSquare(x *big.Int) bool {
	return equal(one, exp(x, eulersCriterionPower, FieldSize))
}

// SquareRoot returns a s.t. a^2=x, as long as x is a square
func SquareRoot(x *big.Int) *big.Int {
	return exp(x, sqrtPower, FieldSize)
}

// YSquared returns x^3+7 mod fieldSize, the right-hand side of the secp256k1
// curve equation.
func YSquared(x *big.Int) *big.Int {
	return mod(add(exp(x, three, FieldSize), seven), FieldSize)
}

// IsCurveXOrdinate returns true iff there is y s.t. y^2=x^3+7
func IsCurveXOrdinate(x *big.Int) bool {
	return IsSquare(YSquared(x))
}

// packUint256s returns xs serialized as concatenated uint256s, or an error
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

var secp256k1Curve = &secp256k1.Secp256k1{}

// Generator is the generator point of secp256k1
var Generator = secp256k1Curve.Point().Base()

// HashUint256s returns a uint256 representing the hash of the concatenated byte
// representations of the inputs
func HashUint256s(xs ...*big.Int) (*big.Int, error) {
	packed, err := packUint256s(xs...)
	if err != nil {
		return &big.Int{}, err
	}
	return utils.MustHash(string(packed)).Big(), nil
}

func uint256ToBytes32(x *big.Int) []byte {
	if x.BitLen() > 256 {
		panic("vrf.uint256ToBytes32: too big to marshal to uint256")
	}
	return common.LeftPadBytes(x.Bytes(), 32)
}

// FieldHash hashes xs uniformly into {0, ..., fieldSize-1}. msg is assumed to
// already be a 256-bit hash
func FieldHash(msg []byte) *big.Int {
	rv := utils.MustHash(string(msg)).Big()
	// Hash recursively until rv < q. P(success per iteration) >= 0.5, so
	// number of extra hashes is geometrically distributed, with mean < 1.
	for rv.Cmp(FieldSize) >= 0 {
		rv = utils.MustHash(string(common.BigToHash(rv).Bytes())).Big()
	}
	return rv
}

// hashToCurveHashPrefix is domain-separation tag for initial HashToCurve hash.
// Corresponds to HASH_TO_CURVE_HASH_PREFIX in VRF.sol.
var hashToCurveHashPrefix = common.BigToHash(one).Bytes()

// HashToCurve is a cryptographic hash function which outputs a secp256k1 point,
// or an error. It passes each candidate x ordinate to ordinates function.
func HashToCurve(p kyber.Point, input *big.Int, ordinates func(x *big.Int),
) (kyber.Point, error) {
	if !(secp256k1.ValidPublicKey(p) && input.BitLen() <= 256 && input.Cmp(zero) >= 0) {
		return nil, fmt.Errorf("bad input to vrf.HashToCurve")
	}
	x := FieldHash(append(hashToCurveHashPrefix, append(secp256k1.LongMarshal(p),
		uint256ToBytes32(input)...)...))
	ordinates(x)
	for !IsCurveXOrdinate(x) { // Hash recursively until x^3+7 is a square
		x.Set(FieldHash(common.BigToHash(x).Bytes()))
		ordinates(x)
	}
	y := SquareRoot(YSquared(x))
	rv := secp256k1.SetCoordinates(x, y)
	if equal(i().Mod(y, two), one) { // Negate response if y odd
		rv = rv.Neg(rv)
	}
	return rv, nil
}

// scalarFromCurveHashPrefix is a domain-separation tag for the hash taken in
// ScalarFromCurve. Corresponds to SCALAR_FROM_CURVE_POINTS_HASH_PREFIX in
// VRF.sol.
var scalarFromCurveHashPrefix = common.BigToHash(two).Bytes()

// ScalarFromCurve returns a hash for the curve points. Corresponds to the
// hash computed in VRF.sol#ScalarFromCurvePoints
func ScalarFromCurvePoints(
	hash, pk, gamma kyber.Point, uWitness [20]byte, v kyber.Point) *big.Int {
	if !(secp256k1.ValidPublicKey(hash) && secp256k1.ValidPublicKey(pk) &&
		secp256k1.ValidPublicKey(gamma) && secp256k1.ValidPublicKey(v)) {
		panic("bad arguments to vrf.ScalarFromCurvePoints")
	}
	// msg will contain abi.encodePacked(hash, pk, gamma, v, uWitness)
	msg := scalarFromCurveHashPrefix
	for _, p := range []kyber.Point{hash, pk, gamma, v} {
		msg = append(msg, secp256k1.LongMarshal(p)...)
	}
	msg = append(msg, uWitness[:]...)
	return i().SetBytes(utils.MustHash(string(msg)).Bytes())
}

// linearComination returns c*p1+s*p2
func linearCombination(c *big.Int, p1 kyber.Point,
	s *big.Int, p2 kyber.Point) kyber.Point {
	return secp256k1Curve.Point().Add(
		secp256k1Curve.Point().Mul(secp256k1.IntToScalar(c), p1),
		secp256k1Curve.Point().Mul(secp256k1.IntToScalar(s), p2))
}

// Proof represents a proof that Gamma was constructed from the Seed
// according to the process mandated by the PublicKey.
//
// N.B.: The kyber.Point fields must contain secp256k1.secp256k1Point values, C,
// S and Seed must be secp256k1Point, and Output must be at
// most 256 bits. See Proof.WellFormed.
type Proof struct {
	PublicKey kyber.Point // secp256k1 public key of private key used in proof
	Gamma     kyber.Point
	C         *big.Int
	S         *big.Int
	Seed      *big.Int // Seed input to verifiable random function
	Output    *big.Int // verifiable random function output;, uniform uint256 sample
}

func (p *Proof) String() string {
	return fmt.Sprintf(
		"vrf.Proof{PublicKey: %s, Gamma: %s, C: %x, S: %x, Seed: %x, Output: %x}",
		p.PublicKey, p.Gamma, p.C, p.S, p.Seed, p.Output)
}

// WellFormed is true iff p's attributes satisfy basic domain checks
func (p *Proof) WellFormed() bool {
	return (secp256k1.ValidPublicKey(p.PublicKey) &&
		secp256k1.ValidPublicKey(p.Gamma) && secp256k1.RepresentsScalar(p.C) &&
		secp256k1.RepresentsScalar(p.S) && p.Output.BitLen() <= 256)
}

var ErrCGammaEqualsSHash = fmt.Errorf(
	"pick a different nonce; c*gamma = s*hash, with this one")

// checkCGammaNotEqualToSHash checks c*gamma ≠ s*hash, as required by solidity
// verifier
func checkCGammaNotEqualToSHash(c *big.Int, gamma kyber.Point, s *big.Int,
	hash kyber.Point) error {
	cGamma := secp256k1Curve.Point().Mul(secp256k1.IntToScalar(c), gamma)
	sHash := secp256k1Curve.Point().Mul(secp256k1.IntToScalar(s), hash)
	if cGamma.Equal(sHash) {
		return ErrCGammaEqualsSHash
	}
	return nil
}

// vrfRandomOutputHashPrefix is a domain-separation tag for the hash used to
// compute the final VRF random output
var vrfRandomOutputHashPrefix = common.BigToHash(three).Bytes()

// VerifyProof is true iff gamma was generated in the mandated way from the
// given publicKey and seed, and no error was encountered
func (p *Proof) VerifyVRFProof() (bool, error) {
	if !p.WellFormed() {
		return false, fmt.Errorf("badly-formatted proof")
	}
	h, err := HashToCurve(p.PublicKey, p.Seed, func(*big.Int) {})
	if err != nil {
		return false, err
	}
	err = checkCGammaNotEqualToSHash(p.C, p.Gamma, p.S, h)
	if err != nil {
		return false, fmt.Errorf("c*γ = s*hash (disallowed in solidity verifier)")
	}
	// publicKey = secretKey*Generator. See GenerateProof for u, v, m, s
	// c*secretKey*Generator + (m - c*secretKey)*Generator = m*Generator = u
	uPrime := linearCombination(p.C, p.PublicKey, p.S, Generator)
	// c*secretKey*h + (m - c*secretKey)*h = m*h = v
	vPrime := linearCombination(p.C, p.Gamma, p.S, h)
	uWitness := secp256k1.EthereumAddress(uPrime)
	cPrime := ScalarFromCurvePoints(h, p.PublicKey, p.Gamma, uWitness, vPrime)
	output := utils.MustHash(string(append(
		vrfRandomOutputHashPrefix, secp256k1.LongMarshal(p.Gamma)...)))
	return equal(p.C, cPrime) && equal(p.Output, output.Big()), nil
}

// generateProofWithNonce allows external nonce generation for testing purposes
//
// As with signatures, using nonces which are in any way predictable to an
// adversary will leak your secret key! Most people should use GenerateProof
// instead.
func generateProofWithNonce(secretKey, seed, nonce *big.Int) (Proof, error) {
	if !(secp256k1.RepresentsScalar(secretKey) && seed.BitLen() <= 256) {
		return Proof{}, fmt.Errorf("badly-formatted key or seed")
	}
	skAsScalar := secp256k1.IntToScalar(secretKey)
	publicKey := secp256k1Curve.Point().Mul(skAsScalar, nil)
	h, err := HashToCurve(publicKey, seed, func(*big.Int) {})
	if err != nil {
		return Proof{}, errors.Wrap(err, "vrf.makeProof#HashToCurve")
	}
	gamma := secp256k1Curve.Point().Mul(skAsScalar, h)
	sm := secp256k1.IntToScalar(nonce)
	u := secp256k1Curve.Point().Mul(sm, Generator)
	uWitness := secp256k1.EthereumAddress(u)
	v := secp256k1Curve.Point().Mul(sm, h)
	c := ScalarFromCurvePoints(h, publicKey, gamma, uWitness, v)
	// (m - c*secretKey) % GroupOrder
	s := mod(sub(nonce, mul(c, secretKey)), secp256k1.GroupOrder)
	if e := checkCGammaNotEqualToSHash(c, gamma, s, h); e != nil {
		return Proof{}, e
	}
	outputHash := utils.MustHash(string(append(vrfRandomOutputHashPrefix,
		secp256k1.LongMarshal(gamma)...)))
	rv := Proof{
		PublicKey: publicKey,
		Gamma:     gamma,
		C:         c,
		S:         s,
		Seed:      seed,
		Output:    outputHash.Big(),
	}
	valid, err := rv.VerifyVRFProof()
	if !valid || err != nil {
		panic("constructed invalid proof")
	}
	return rv, nil
}

// GenerateProof returns gamma, plus proof that gamma was constructed from seed
// as mandated from the given secretKey, with public key secretKey*Generator
//
// secretKey and seed must be less than secp256k1 group order. (Without this
// constraint on the seed, the samples and the possible public keys would
// deviate very slightly from uniform distribution.)
func GenerateProof(secretKey, seed common.Hash) (Proof, error) {
	for {
		nonce, err := rand.Int(rand.Reader, secp256k1.GroupOrder)
		if err != nil {
			return Proof{}, err
		}
		proof, err := generateProofWithNonce(secretKey.Big(), seed.Big(), nonce)
		switch {
		case err == ErrCGammaEqualsSHash:
			// This is cryptographically impossible, but if it were ever to happen, we
			// should try again with a different nonce.
			continue
		case err != nil: // Any other error indicates failure
			return Proof{}, err
		default:
			return proof, err // err should be nil
		}
	}
}
