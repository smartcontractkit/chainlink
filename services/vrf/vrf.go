// Package vrf implements an RSA-based Verifiable Random Function, roughly as
// specified in https://tools.ietf.org/html/draft-irtf-cfrg-vrf-03#section-4.1
// (with some modifications to make on-chain verification cheaper), and as
// described in https://eprint.iacr.org/2017/099.pdf , figure 1 and section 4.1.
//
// The key size and public exponent used for this protocol are set in the
// constants KeySizeBits and PublicExponent. Recompile with a different key
// size, if necessary. Note that this will require a corresponding change to the
// on-chain contract, VRF.sol, and its tests, VRF_test.js. We don't recommend
// changing PublicExponent: any change will at least double the gas cost for
// on-chain verification. Precautions have been taken (in seedToRingValue) to
// mitigate the risk of using a small public exponent.
//
// Private keys with the required PublicExponent and key size can be generated
// with code like
//
//   key, err := MakeKey()
//   if err != nil { panic(err) }
//
// The prime factors p and q in key are safe primes
// (https://en.wikipedia.org/wiki/Safe_prime), meaning (p-1)/2, (q-1)/2 are
// intended to be prime. Searching for these is a bit slow, so allow MakeKey a
// couple of minutes.
//
// Generate(key, seed) will generate a Proof object. Its Output field contains
// the VRF output for the given seed and key. The rest of the fields provide the
// proof that Output was generated as mandated by the key and seed. To those who
// don't know the secret key, Output values should be statistically
// indistinguishable from uniform random samples from {0, ..., 2**256-1}.
//
// Given a Proof object p, use p.Verify() to check that p is correct, or pass
// its fields to the corresponding arguments of VRF.sol's isValidVRFOutput
// method for on-chain verification.
package vrf

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/chainlink/utils"
)

// KeySizeBits is the number of bits expected in RSA VRF keys.
// Any change to this forces change to VRF.sol and VRF_test.js, too.
const KeySizeBits = 2048
const keySizeBytes = KeySizeBits / 8
const keySizeWords = KeySizeBits / 256

const wordBytes = 32

// PublicExponent is the exponent used in RSA public keys
// Any change to this forces change to VRF.sol and VRF_test.js, too.
//
// A public exponent of 3 is very cheap to calculate, in terms of ethereum gas,
// and no risk for this application because the seed is uniformly randomly
// sampled from the ring.
const PublicExponent = 3

// panicUnless is a basic golang "assert"
func panicUnless(prop bool, message string) {
	if !prop {
		panic(message)
	}
}

// Textbook RSA "decryption", copied from crypto/rsa.go/decypt function. Returns
// seed raised to the private exponent of k, k's modulus. Uses faster CRT method
// if enabled on k.
func decrypt(k *rsa.PrivateKey, seed *big.Int) *big.Int {
	panicUnless(len(k.Primes) == 2,
		"RSA VRF only works with two-factor moduli.")
	if k.Precomputed.Dp == nil { // Do it the slow way
		return new(big.Int).Exp(seed, k.D, k.N)
	}
	// We have the precalculated values needed for the CRT.
	m := new(big.Int).Exp(seed, k.Precomputed.Dp, k.Primes[0])
	m2 := new(big.Int).Exp(seed, k.Precomputed.Dq, k.Primes[1])
	m.Sub(m, m2)
	if m.Sign() < 0 {
		m.Add(m, k.Primes[0])
	}
	m.Mul(m, k.Precomputed.Qinv)
	m.Mod(m, k.Primes[0])
	m.Mul(m, k.Primes[1])
	m.Add(m, m2)
	return m
}

// encrypt "encrypts" m under the publicKey, with textbook RSA.
//
// This is inadequate for actual encryption. Use rsa.EncryptOAEP for that. The
// point here is not to encrypt sensitive data, but to use an operation which
// can only be reversed with knowledge of the private key.
func encrypt(publicKey *rsa.PublicKey, m *big.Int) *big.Int {
	exponentBig := new(big.Int).SetUint64(uint64(publicKey.E))
	return new(big.Int).Exp(m, exponentBig, publicKey.N)
}

// Proof represents a proof that Output was generated as mandated by the Key
// from the given Seed, via Decrypt.
type Proof struct {
	Key                   *rsa.PublicKey
	Seed, Decrypt, Output *big.Int
	bitSize               uint32
}

// asUint256 returns i represented as an array of packed uint256's, a la solidity
func asUint256Array(i *big.Int) []byte {
	inputBytes := i.Bytes()
	outputBytesDeficit := wordBytes - (len(inputBytes) % wordBytes)
	if outputBytesDeficit == wordBytes {
		outputBytesDeficit = 0
	}
	rv := append(make([]byte, outputBytesDeficit), inputBytes...)
	panicUnless(len(rv)%32 == 0 && i.Cmp(new(big.Int).SetBytes(rv)) == 0,
		"rv is not i as big-endian uint256 array")
	return rv
}

func asKeySizeUint256Array(i *big.Int) []byte {
	a := asUint256Array(i)
	rv := append(make([]byte, (keySizeBytes-len(a))), a...)
	panicUnless(len(rv) == keySizeBytes, "i must fit in keySizeBytes")
	panicUnless(len(rv)%wordBytes == 0, "must generate packed uint256 array")
	panicUnless(i.Cmp(new(big.Int).SetBytes(rv)) == 0, "must represent i")
	return rv
}

// decryptionToOutput generates the actual randomness to be output by the VRF, from
// the output of the RSA "decryption"
func decryptionToOutput(decryption *big.Int) (*big.Int, error) {
	decrypt := asKeySizeUint256Array(decryption)
	output, err := utils.Keccak256(decrypt)
	if err != nil {
		return nil, err
	}
	return new(big.Int).SetBytes(output), nil
}

// makeLongHash returns a numWords*256-bits hash of iv, or an error
func makeLongHash(iv []byte, numWords uint32) ([]byte, error) {
	hash, err := utils.Keccak256(iv)
	if err != nil {
		return nil, err
	}
	panicUnless(len(hash) == wordBytes, "hash is not 256 bits")
	rv := make([]byte, wordBytes)
	copy(rv, hash)
	for uint32(len(rv)) != numWords*wordBytes {
		// Recursively hash last hash and use it to extend ring value.
		hash, err := utils.Keccak256(hash)
		if err != nil {
			return nil, err
		}
		panicUnless(len(hash) == wordBytes, "hash is not 256 bits")
		rv = append(rv, hash...)
	}
	return rv, nil
}

// seedToRingValue hashes seed roughly uniformly to {0, ..., k.N - 1}
//
// This plays the same role as the Mask Generation Function, MGF1, in
// https://tools.ietf.org/html/draft-irtf-cfrg-vrf-03#section-4.1 , and the same
// role as padding, in regular RSA.
//
// Forcing the VRF prover to "decrypt" this, rather than the initial seed
// itself, prevents an adversary from submitting random "proofs" that, for
// instance, back out a seed from encrypting a "decrypt" using the public key,
// or exploiting simple arithmetic relationships which hold over the naturals
// and hence all moduli, like 2³≡8 mod k.N. See also
// https://www.di.ens.fr/~fouque/ens-rennes/coppersmith.pdf
func seedToRingValue(seed *big.Int, k *rsa.PublicKey) (*big.Int, error) {
	seedBytes := asUint256Array(seed)
	if len(seedBytes) != wordBytes || seed.Cmp(big.NewInt(0)) == -1 {
		return nil, fmt.Errorf("Seed must fit in uint256")
	}
	if k.N.BitLen() > KeySizeBits || k.N.BitLen() < 0 {
		return nil, fmt.Errorf("public modulus must fit in key size")
	}
	expAsUint256 := asUint256Array(new(big.Int).SetUint64(uint64(k.E)))
	ringValueBytes, err := makeLongHash(bytes.Join([][]byte{
		seedBytes, asKeySizeUint256Array(k.N), expAsUint256}, nil),
		// Generate something much longer than we need (keySizeWords+1),
		// to ensure more uniform sampling from {0, ..., k.N-1} on
		// reduction mod k.N.
		keySizeWords+1)
	if err != nil {
		return nil, err
	}
	ringValue := new(big.Int).Mod(new(big.Int).SetBytes(ringValueBytes), k.N)
	panicUnless(PublicExponent*ringValue.BitLen() >= 2*k.N.BitLen(),
		`(ring value)^exponent is too short to be secure.
(For a 2048-bit or longer key, this is cryptographically impossible.)`)
	return ringValue, nil
}

// checkKey returns an error describing any problems with k.
func checkKey(k *rsa.PrivateKey) error {
	if k.E != PublicExponent {
		return fmt.Errorf("public exponent of key must be PublicExponent")
	}
	return k.Validate()
}

// Generate returns VRF output and correctness proof from given key and seed
func Generate(k *rsa.PrivateKey, seed *big.Int) (*Proof, error) {
	if err := checkKey(k); err != nil {
		return nil, err
	}
	// Prove knowledge of the private key by "decrypting" to seed used to
	// generate Proof.Output. Nothing hidden, here, so not really decryption
	cipherText, err := seedToRingValue(seed, &k.PublicKey)
	if err != nil {
		return nil, err
	}
	decryption := decrypt(k, cipherText)
	output, err := decryptionToOutput(decryption) // Actual VRF "randomness"
	if err != nil {
		return nil, err
	}
	rv := &Proof{
		Key:     &k.PublicKey,
		Seed:    seed,
		Decrypt: decryption,
		Output:  output,
	}
	ok, err := rv.Verify()
	if err != nil {
		panic(err)
	}
	panicUnless(ok, "couldn't verify proof we just generated")
	return rv, nil
}

// Verify returns true iff p is a correct proof for its output.
func (p *Proof) Verify() (bool, error) {
	output, err := decryptionToOutput(p.Decrypt)
	if err != nil {
		return false, err
	}
	if output.Cmp(p.Output) != 0 { // Verify Output is hash of Decrypt
		return false, nil
	}
	// Get the value from the seed which prover should have "decrypted"
	expected, err := seedToRingValue(p.Seed, p.Key)
	if err != nil {
		return false, err
	}
	return encrypt(p.Key, p.Decrypt).Cmp(expected) == 0, nil
}

// safePrime returns 2p+1 which is composite with probability less than 2^{-80},
// and for which (p-1)/2 is composite with probability less than 2^{-40}.
// https://en.wikipedia.org/wiki/Safe_prime
//
// This must use golang version at least 1.10.3. See section 4.15,
// https://eprint.iacr.org/2018/749.pdf#page=19
func safePrime(bits uint32) *big.Int {
	one, two, three := big.NewInt(1), big.NewInt(2), big.NewInt(3)
	scratch1, scratch2, scratch3 := new(big.Int), new(big.Int), new(big.Int)
	for {
		p, err := rand.Prime(rand.Reader, int(bits)-1)
		if err != nil {
			panic(err)
		}
		twoP := scratch2.Lsh(p, 1)
		rv := scratch1.Add(twoP, one) // 2*p+1
		// See https://en.wikipedia.org/wiki/Pocklington_primality_test
		// equations 1-3, N:=2p+1, p:=p and a:=2
		if coprime(three, rv) && // Eq. 3: a^{(N-1)/p}-1=2^2-1=4-1=3
			scratch3.Exp(two, twoP, rv).Cmp(one) == 0 { // Eq1
			// This extra primality test, combined with the
			// Pocklington test and the primality test on p, gives
			// (2^{-40})^2 probability that the returned value is
			// actually composite
			if rv.ProbablyPrime(20) {
				return rv
			}
		}
	}
}

// coprime returns true iff GCD(m,n) = 1
func coprime(m, n *big.Int) bool {
	return new(big.Int).GCD(nil, nil, m, n).Cmp(big.NewInt(1)) == 0
}

// coprimalityChecks panics if any expected coprimality does not hold.
func coprimalityChecks(p, q, pMinusOne, qMinusOne, multOrder, exp *big.Int) {
	for _, tt := range []struct {
		m   *big.Int
		n   *big.Int
		msg string
	}{
		{p, q, "p and q"},
		{new(big.Int).Mul(p, q), multOrder, "pq and (p-1)(q-1)"},
		{new(big.Int).Rsh(pMinusOne, 1), qMinusOne, "(p-1)/2, q-1"},
		{multOrder, exp, "(p-1)(q-1), exponent"},
	} {
		panicUnless(coprime(tt.m, tt.n), tt.msg+" not coprime")
	}
}

// MakeKey securely randomly samples a pair of large primes for an RSA modulus,
// and sets the public exponent to PublicExponent.
//
// Without an argument, defaults to KeySizeBits-sized key. Otherwise, makes a
// key of the requested size. The argument form should only be used for testing.
//
// Because this searches for safe primes, it may take a couple of minutes, even
// on a modern machine.
func MakeKey(bitsizes ...uint32) (*rsa.PrivateKey, error) {
	if len(bitsizes) > 1 {
		return nil, fmt.Errorf("specify at most one bit size")
	}
	bitsize := uint32(KeySizeBits)
	if len(bitsizes) == 1 {
		bitsize = bitsizes[0]
		fmt.Printf("Warning, generating a key of length %d. %d bits demanded by protocol\n",
			bitsize, KeySizeBits)
	}
	zero, one := big.NewInt(0), big.NewInt(1)
	exp := new(big.Int).SetUint64(uint64(PublicExponent))
	p := safePrime(bitsize / 2)
	pMinusOne := new(big.Int).Sub(p, one)
	q := safePrime(bitsize / 2)
	qMinusOne := new(big.Int).Sub(q, one)
	N := new(big.Int).Mul(p, q)
	panicUnless(uint32(N.BitLen()) == bitsize, "modulus doesn't match key size")
	multOrder := new(big.Int).Mul(pMinusOne, qMinusOne)
	coprimalityChecks(p, q, pMinusOne, qMinusOne, multOrder, exp)
	D := new(big.Int) // Will receive "exp^{-1} mod multOrder" from GCD
	_ = new(big.Int).GCD(D, nil, exp, multOrder)
	D = D.Mod(D, multOrder)
	panicUnless(D.Cmp(zero) == 1, "D not positive") // From prior line
	dExpProd := new(big.Int).Mod(new(big.Int).Mul(D, exp), multOrder)
	panicUnless(dExpProd.Cmp(one) == 0, "(exp * D) ~≡ 1 mod N")
	rv := rsa.PrivateKey{
		PublicKey: rsa.PublicKey{N: N, E: PublicExponent},
		D:         D,
		Primes:    []*big.Int{p, q},
	}
	rv.Precompute()
	if err := rv.Validate(); err != nil {
		return nil, err
	}
	return &rv, nil
}
