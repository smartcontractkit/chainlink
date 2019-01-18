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
// A secret key sk should be securely sampled uniformly from {0, ..., Order-1}.
// The public key associated with it can be calculated from it by
// bn256.ScalarMult(Generator,sk). A convenience function for this is provided
// by GenerateKeyPair().
//
// To generate a random output, pass the keypair from GenerateKeyPair and a
// *big.Int seed to GenerateProof, and use the Output field of the returned
// Proof object.
//
// To verify a Proof object p, run p.Verify(), or pass its fields to the
// corresponding arguments of isValidVRFOutput on the VRF solidity contract, to
// verify it on-chain.

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"math/big"

	// NB: this curve is actually alt-bn128, not bn256!
	"github.com/ethereum/go-ethereum/crypto/bn256"

	"github.com/smartcontractkit/chainlink/utils"
)

// Int compensates for awkward big.Int API. This package is fast; we don't need
// to bash things in-place.
type integer struct{ i *big.Int }

func bi(i uint) integer { return integer{big.NewInt(int64(i))} }

var zero, one, two, three, four = bi(0), bi(1), bi(2), bi(3), bi(4)

func i() *big.Int                              { return new(big.Int) }
func (x integer) add(s integer) integer        { return integer{i().Add(x.i, s.i)} }
func (x integer) div(divisor integer) integer  { return integer{i().Div(x.i, divisor.i)} }
func (x integer) equal(y integer) bool         { return x.i.Cmp(y.i) == 0 }
func (x integer) lessThan(y integer) bool      { return x.i.Cmp(y.i) == -1 }
func (x integer) greaterThan(y integer) bool   { return x.i.Cmp(y.i) == 1 }
func (x integer) exp(exponent integer) integer { return integer{i().Exp(x.i, exponent.i, P.i)} }
func (x integer) lsh(bits uint) integer        { return integer{i().Lsh(x.i, bits)} }
func (x integer) rsh(bits uint) integer        { return integer{i().Rsh(x.i, bits)} }
func (x integer) mul(y integer) integer        { return integer{i().Mul(x.i, y.i)} }
func (x integer) mod(y integer) integer        { return integer{i().Mod(x.i, y.i)} }
func (x integer) sub(y integer) integer        { return integer{i().Sub(x.i, y.i)} }
func (x integer) modSqrt() integer             { return integer{i().ModSqrt(x.i, P.i)} }
func (x integer) and(mask integer) integer     { return integer{i().And(x.i, mask.i)} }
func (x integer) bitLen() uint                 { return uint(x.i.BitLen()) }
func (x integer) copy(y integer)               { x.i.Set(y.i) }

func integerFromBytes(b []byte) integer { return integer{i().SetBytes(b)} }

func bigFromBase10(s string) integer {
	n, ok := new(big.Int).SetString(s, 10)
	if !ok {
		panic("failed to allocate big.Int")
	}
	return integer{n}
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

// (P-1)/2: Half Fermat's Little Theorem exponent
var eulersCriterionPower = P.sub(one).rsh(1)

// isSquare returns true iff x = y^2 for some y in GF(p)
func isSquare(x integer) bool {
	return x.exp(eulersCriterionPower).equal(one)
}

// ySquared returns x^3+3 mod P
func ySquared(x integer) integer {
	return x.exp(three).add(three).mod(P)
}

// IsCurveXOrdinate returns true iff there is y s.t. y^2=x^3+3
func isCurveXOrdinate(x integer) bool {
	return isSquare(ySquared(x))
}

func packUint256s(xs ...integer) ([]byte, error) {
	mem := bytes.Buffer{}
	for _, x := range xs {
		word, err := utils.EVMWordBigInt(x.i)
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
func newCurvePoint(x, y integer) (*bn256.G1, error) {
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
func hashUint256s(xs ...integer) (integer, error) {
	packed, err := packUint256s(xs...)
	if err != nil {
		return zero, err
	}
	hash, err := utils.Keccak256(packed)
	if err != nil {
		return zero, err
	}
	return integerFromBytes(hash), nil
}

// maskHash returns hashUint256s(xs...) & mask
func maskHash(mask integer, xs ...integer) (integer, error) {
	x, err := hashUint256s(xs...)
	if err != nil {
		return zero, err
	}
	return x.and(mask), nil
}

// zqHash hashes xs uniformly into {0, ..., q-1}
func zqHash(q integer, xs ...integer) (integer, error) {
	if len(xs) < 1 {
		panic("can't take hash of empty list. You might have forgotten argument q")
	}
	if q.bitLen() > 256 {
		panic(fmt.Errorf(
			"will only generate 256 bits of entropy, need %v",
			q.bitLen()))
	}
	// Bits which can be used in representation of a number less than q.
	// 2^(q.BitLen)-1
	orderMask := one.lsh(uint(q.bitLen())).sub(one)
	rv, err := maskHash(orderMask, xs...)
	if err != nil {
		return zero, err
	}
	// Hash recursively until rv < q. P(success per iteration) >= 0.5, so
	// number of extra hashes is geometrically distributed, with mean < 1.
	for rv.greaterThan(q) {
		nrv, err := maskHash(orderMask, rv)
		if err != nil {
			return zero, err
		}
		rv.copy(nrv)
	}
	return rv, nil
}

// hashToCurve is a one-way hash function onto the curve
func hashToCurve(px, py, input integer) (*bn256.G1, error) {
	x, err := zqHash(P, px, py, input)
	if err != nil {
		return &bn256.G1{}, err
	}
	for !isCurveXOrdinate(x) { // Hash recursively until x^3+3 is a square
		nx, err := zqHash(P, x)
		if err != nil {
			return &bn256.G1{}, err
		}
		x.copy(nx)
	}
	return newCurvePoint(x, ySquared(x).modSqrt())
}

// scalarFromCurve returns a hash for the curve points. Corresponds to the hash
// computed in Curve.sol#scalarFromCurve
func scalarFromCurve(ps ...*bn256.G1) (integer, error) {
	coordinates := make([]integer, (len(ps)+1)*2)
	gx, gy := CoordsFromPoint(Generator)
	coordinates[0] = integer{gx}
	coordinates[1] = integer{gy}
	for ordidx, p := range ps {
		x, y := CoordsFromPoint(p)
		coordinates[2*ordidx+2] = integer{x}
		coordinates[2*ordidx+3] = integer{y}
	}
	return zqHash(Order, coordinates...)
}

// linearComination returns c*p1+s*p2
func linearComination(c integer, p1 *bn256.G1, s integer, p2 *bn256.G1) *bn256.G1 {
	return new(bn256.G1).Add(
		new(bn256.G1).ScalarMult(p1, c.i),
		new(bn256.G1).ScalarMult(p2, s.i))
}

// KeyPair represents a public/private keypair
type KeyPair struct {
	Public *bn256.G1
	secret *big.Int
}

// GenerateKeyPair returns a public/private keypair generated from
// cryptographically strong pseudorandomness, or an error on failure.
func GenerateKeyPair() (*KeyPair, error) {
	rv := new(KeyPair)
	var err error
	rv.secret, err = rand.Int(rand.Reader, Order.sub(one).i)
	if err != nil {
		return rv, err
	}
	rv.Public = new(bn256.G1).ScalarMult(Generator, rv.secret)
	return rv, nil
}

// Marshal should not be used on a KeyPair
func (keypair *KeyPair) Marshal(out []byte) {
	panic("don't use this, use MarshalEncryptedKeyPair instead")
}

// Unmarshal should not be used on a KeyPair
func (keypair *KeyPair) Unmarshal(in []byte) {
	panic("don't use this, use UnmarshalEncryptedKeyPair or UnmarshalPublicKey instead")
}

// dummyCipher is a blank cipher to be returned on an error condition
type dummyCipher struct{}

func (*dummyCipher) BlockSize() int          { panic("never use this function") }
func (*dummyCipher) Encrypt(dst, src []byte) { panic("never use this function") }
func (*dummyCipher) Decrypt(dst, src []byte) { panic("never use this function") }

func createKeyCipher(salt, passPhrase []byte) (cipher.Block, error) {
	// XXX: Append is over-writing the encrypted key!!! Take a copy, or
	// send one in.
	key, err := utils.Keccak256(append(salt, passPhrase...))
	if err != nil {
		return &dummyCipher{}, nil
	}
	cipherRV, err := aes.NewCipher(key)
	if err != nil {
		return &dummyCipher{}, nil
	}
	return cipherRV, nil
}

// MarshalEncryptedKeyPair marshals keypair with Secret encrypted by passPhrase,
// or returns an error on failure.
//
// The format is (Public.x || Public.y || encrypted secret), where all fields
// are 256 bits, || is concatenation, and Public.{x,y} are big-endian affine
// coordinates of KeyPair.Public.
//
// The encryption key is keccak256(Public.x || Public.y || passPhrase), and the
// plaintext of the encrypted secret is KeyPair.secret as big-endian byte array.
func (keypair *KeyPair) MarshalEncryptedKeyPair(passPhrase []byte) ([]byte, error) {
	if (integer{keypair.secret}).lessThan(zero) ||
		(integer{keypair.secret}).greaterThan(Order.sub(one)) {
		return []byte{}, fmt.Errorf(
			"Secret key %v is not in {1,...Order-1}", keypair.secret)
	}
	rv := keypair.Public.Marshal()                 // Marshal public key in plaintext
	cipher, err := createKeyCipher(rv, passPhrase) // Salt with encryption key with public key.
	if err != nil {
		return []byte{}, err
	}
	cipherText := make([]byte, 32)
	byteDeficit := (256 - keypair.secret.BitLen()) / 8
	// Bytes returns big-endian (higher-order bytes first), so pre-pad with
	// 0's to bring input length up to 256 bits/32 bytes.
	prefix := make([]byte, byteDeficit)
	for i := range prefix {
		prefix[i] = 0
	}
	secretInput := append(prefix, keypair.secret.Bytes()...)
	if (len(secretInput) != 32) ||
		!integerFromBytes(secretInput).equal(integer{keypair.secret}) {
		panic("failed to encode secret with correct length")
	}
	cipher.Encrypt(cipherText[:16], secretInput[:16])
	testOutput := make([]byte, 16)
	cipher.Decrypt(testOutput, cipherText[:16])
	cipher.Encrypt(cipherText[16:], secretInput[16:])
	result := append(rv, cipherText...)
	return result, nil
}

// UnmarshalEncryptedKeyPair returns a KeyPair after decrypting the secret key with passPhrase
func UnmarshalEncryptedKeyPair(encKeyPair, passPhrase []byte) (*KeyPair, error) {
	if len(encKeyPair) != 96 {
		return &KeyPair{}, fmt.Errorf("serialized key pair should be 3 uint256's")
	}
	rv := new(KeyPair)
	publicKey, err := UnmarshalPublicKey(encKeyPair)
	encSecretKey := make([]byte, 32) // Copy because createKeyCipher
	copy(encSecretKey, encKeyPair[64:])
	if err != nil {
		return &KeyPair{}, err
	}
	rv.Public = publicKey
	marshaledPublicKeyAsSalt := encKeyPair[:len(encKeyPair)-len(encSecretKey)]
	cipher, err := createKeyCipher(marshaledPublicKeyAsSalt, passPhrase)
	if err != nil {
		return &KeyPair{}, err
	}
	secretOutput := make([]byte, 32)
	cipher.Decrypt(secretOutput[:16], encSecretKey[:16])
	cipher.Decrypt(secretOutput[16:], encSecretKey[16:])
	rv.secret = integerFromBytes(secretOutput).i

	return rv, nil
}

// UnmarshalPublicKey returns the unencrypted public key of a marshaled KeyPair
func UnmarshalPublicKey(encKeyPair []byte) (*bn256.G1, error) {
	rv := new(bn256.G1)
	_, err := rv.Unmarshal(encKeyPair)
	if err != nil {
		return &bn256.G1{}, err
	}
	return rv, nil
}

// String generates a human-readable representation of the KeyPair, blinding the secret key
func String(k *KeyPair) string {
	return fmt.Sprintf("KeyPair{Public:%v, secret:<redacted>}", k.Public)
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
	h, err := hashToCurve(integer{px}, integer{py}, integer{proof.Seed})
	if err != nil {
		return false, err
	}
	// publicKey = secretKey*Generator. See GenerateProof for u, v, m, s
	// c*secretKey*Generator + (m - c*secretKey)*Generator = m*Generator = u
	uPrime := linearComination(integer{proof.C}, proof.PublicKey, integer{proof.S}, Generator)
	// c*secretKey*h + (m - c*secretKey)*h = m*h = v
	vPrime := linearComination(integer{proof.C}, proof.Gamma, integer{proof.S}, h)
	cPrime, _ := scalarFromCurve(
		h, proof.PublicKey, proof.Gamma, uPrime, vPrime)
	if err != nil {
		return false, err
	}
	output, err := utils.Keccak256(proof.Gamma.Marshal())
	if err != nil {
		return false, err
	}
	return integer{proof.C}.equal(cPrime) &&
			integer{proof.Output}.equal(integerFromBytes(output)),
		nil
}

// makeProof proof generates the actual proof, modulo the actual random output
func makeProof(keypair *KeyPair, seed integer) (*Proof, error) {
	px, py := CoordsFromPoint(keypair.Public)
	h, err := hashToCurve(integer{px}, integer{py}, seed)
	if err != nil {
		return &Proof{}, err
	}
	gamma := new(bn256.G1).ScalarMult(h, keypair.secret)
	baseM, err := rand.Int(rand.Reader, Order.i)
	if err != nil {
		return &Proof{}, err
	}
	m := integer{baseM}
	u := new(bn256.G1).ScalarMult(Generator, m.i)
	v := new(bn256.G1).ScalarMult(h, m.i)
	c, err := scalarFromCurve(h, keypair.Public, gamma, u, v)
	if err != nil {
		return &Proof{}, err
	}
	// s = (m - c*secretKey1) % Order
	s := m.sub(c.mul(integer{keypair.secret})).mod(Order)
	return &Proof{
		PublicKey: keypair.Public,
		Gamma:     gamma,
		C:         c.i,
		S:         s.i,
		Seed:      seed.i,
	}, nil
}

// GenerateProof returns gamma, plus proof that gamma was constructed from seed
// as mandated from the given secretKey, with public key secretKey*Generator
func GenerateProof(keypair *KeyPair, seed *big.Int) (*Proof, error) {
	proof, err := makeProof(keypair, integer{seed})
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
