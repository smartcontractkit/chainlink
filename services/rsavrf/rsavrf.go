package rsavrf

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/chainlink/utils"
)

// KeySizeBits is the number of bits expected in RSA VRF keys.
// Any change to this forces change to RSAVRF.sol and RSAVRF_test.js, too.
const KeySizeBits = 2048
const keySizeBytes = KeySizeBits / 8
const keySizeWords = KeySizeBits / 256

const wordBytes = 32

// PublicExponent is the exponent used in RSA public keys
// Any change to this forces change to RSAVRF.sol and RSAVRF_test.js, too.
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
func decrypt(k *rsa.PrivateKey, seed *big.Int) (m *big.Int) {
	panicUnless(len(k.Primes) == 2,
		"RSA VRF only works with two-factor moduli.")
	if k.Precomputed.Dp == nil { // Do it the slow way
		return new(big.Int).Exp(seed, k.D, k.N)
	}
	// We have the precalculated values needed for the CRT.
	m = new(big.Int).Exp(seed, k.Precomputed.Dp, k.Primes[0])
	m2 := new(big.Int).Exp(seed, k.Precomputed.Dq, k.Primes[1])
	m.Sub(m, m2)
	if m.Sign() < 0 {
		m.Add(m, k.Primes[0])
	}
	m.Mul(m, k.Precomputed.Qinv)
	m.Mod(m, k.Primes[0])
	m.Mul(m, k.Primes[1])
	m.Add(m, m2)
	return
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
}

// decryptToOutput generates the actual randomness to be output by the VRF, from
// the output of the RSA "decryption"
func decryptToOutput(decryption *big.Int) (rv *big.Int, err error) {
	decryptionBytes := decryption.Bytes()
	deficit := keySizeBytes - len(decryptionBytes)
	if deficit < 0 {
		return &big.Int{}, fmt.Errorf("decrypt string too long: %v",
			decryption)
	}
	decrypt := append(make([]byte, deficit), // Prepad with zeros
		decryptionBytes...)
	panicUnless(len(decrypt) != keySizeBytes, "decrypt has incorrect length")
	output, err := utils.Keccak256(decrypt)
	if err != nil {
		return &big.Int{}, err
	}
	return rv.SetBytes(output), nil
}

// asUint256 returns i represented as an array of packed uint256's, a la solidity
func asUint256Array(i *big.Int) []byte {
	inputBytes := i.Bytes()
	outputBytesDeficit := wordBytes - (len(inputBytes) % wordBytes)
	rv := append(make([]byte, outputBytesDeficit), inputBytes...)
	panicUnless(len(rv)%32 == 0 && i.Cmp(new(big.Int).SetBytes(rv)) == 0,
		"rv is not i as big-endian uint256 array")
	return rv
}

func asKeySizeUint256Array(i *big.Int) []byte {
	a := asUint256Array(i)
	return append(make([]byte, (keySizeBytes-len(a))), a...)
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
// or exploit simple arithmetic relationships which hold over the naturals and
// hence all moduli, like 2³≡8 mod k.N. See also
// https://www.di.ens.fr/~fouque/ens-rennes/coppersmith.pdf
func seedToRingValue(seed *big.Int, k *rsa.PublicKey) (*big.Int, error) {
	seedBytes := asUint256Array(seed)
	expAsUint256 := asUint256Array(new(big.Int).SetUint64(uint64(k.E)))
	if len(seedBytes) != 1 || seed.Cmp(big.NewInt(0)) == -1 {
		return &big.Int{}, fmt.Errorf("Seed must fit in uint256")
	}
	if k.N.BitLen() != KeySizeBits || k.N.Cmp(big.NewInt(0)) == -1 {
		return &big.Int{}, fmt.Errorf("public modulus must fit in key size")
	}
	hash, err := utils.Keccak256(bytes.Join([][]byte{
		seedBytes, asKeySizeUint256Array(k.N), expAsUint256}, []byte{}))
	if err != nil {
		return &big.Int{}, err
	}
	ringValueBytes := make([]byte, wordBytes)
	copy(ringValueBytes, hash)
	// Ensure more uniform distribution in {0,...,k.N-1} by initially
	// generating a much larger value to take the modulus of.
	// (keySizeWords+1, rather than keySizeWords.)
	for len(ringValueBytes) != keySizeWords+1 {
		// Recursively hash last hash and use it to extend ring value.
		hash, err := utils.Keccak256(hash)
		if err != nil {
			return &big.Int{}, err
		}
		panicUnless(len(hash) == wordBytes, "hash is not 256 bits")
		ringValueBytes = append(ringValueBytes, hash...)
	}
	ringValue := new(big.Int).Mod(new(big.Int).SetBytes(ringValueBytes), k.N)
	panicUnless(PublicExponent*ringValue.BitLen() >= 2*k.N.BitLen(),
		`(ring value)^exponent is too short to be secure.
(For a 2048-bit or longer key, this is cryptographically impossible.)`)
	return ringValue, nil
}

// Generate returns VRF output and correctness proof from given key and seed
func Generate(k *rsa.PrivateKey, seed *big.Int) (*Proof, error) {
	if k.E != PublicExponent {
		return &Proof{}, fmt.Errorf("public exponent of key must be PublicExponent")
	}
	if err := k.Validate(); err != nil {
		return &Proof{}, err
	}
	// Prove knowledge of the private key by "decrypting" to seed used to
	// generate Proof.Output. Nothing hidden, here, so not really decryption
	cipherText, err := seedToRingValue(seed, &k.PublicKey)
	if err != nil {
		return &Proof{}, err
	}
	decryption := decrypt(k, cipherText)
	output, err := decryptToOutput(decryption) // Actual VRF "randomness"
	if err != nil {
		return &Proof{}, err
	}
	return &Proof{
		Key:     &k.PublicKey,
		Seed:    seed,
		Decrypt: decryption,
		Output:  output,
	}, nil
}

// Verify returns true iff p is a correct proof for its output.
func (p *Proof) Verify() (bool, error) {
	output, err := decryptToOutput(p.Decrypt)
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

// safePrime returns p which is composite with probability less than 2^{-80},
// and for which (p-1)/2 is composite with probability less than 2^{-40}.
// https://en.wikipedia.org/wiki/Safe_prime
//
// This must use golang version at least 1.10.3. See section 4.15,
// https://eprint.iacr.org/2018/749.pdf#page=19
func safePrime(bits uint32) *big.Int {
	scratch1, scratch2 := new(big.Int), new(big.Int)
	one := big.NewInt(1)
	for {
		p, err := rand.Prime(rand.Reader, int(bits)-1)
		if err != nil {
			panic(err)
		}
		rv := scratch1.Add(scratch2.Lsh(p, 1), one) // 2*p+1
		// XXX: Verify Pocklington as well:
		// https://en.wikipedia.org/wiki/Pocklington_primality_test
		if rv.ProbablyPrime(40) {
			return rv
		}
	}
}

// coprime returns true iff GCD(m,n) = 1
func coprime(m, n *big.Int) bool {
	return new(big.Int).GCD(nil, nil, m, n).Cmp(big.NewInt(1)) == 0
}

// coprimalityChecks verifies that everything's pairwise coprime.
func coprimalityChecks(p, q, pMinusOne, qMinusOne, multOrder, exp *big.Int) {
	for _, tt := range []struct {
		m   *big.Int
		n   *big.Int
		msg string
	}{
		{p, q, "p and q"},
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
// Because this searches for safe primes, it may take a couple of minutes even
// on a modern machine.
func MakeKey(bitsizes ...uint32) (*rsa.PrivateKey, error) {
	if len(bitsizes) > 1 {
		return &rsa.PrivateKey{}, fmt.Errorf("specify at most one bit size")
	}
	bitsize := uint32(KeySizeBits)
	if len(bitsizes) == 1 {
		bitsize = bitsizes[0]
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
	fmt.Println("D", D)
	panicUnless(D.Cmp(zero) != 1, "D not positive") // From prior line
	dExpProd := new(big.Int).Mod(new(big.Int).Mul(D, exp), multOrder)
	panicUnless(dExpProd.Cmp(one) != 0, "(exp * D) ~≡ 1 mod N")
	rv := rsa.PrivateKey{
		PublicKey: rsa.PublicKey{N: N, E: PublicExponent},
		D:         D,
		Primes:    []*big.Int{p, q},
	}
	rv.Precompute()
	if err := rv.Validate(); err != nil {
		return &rsa.PrivateKey{}, err
	}
	return &rv, nil
}
