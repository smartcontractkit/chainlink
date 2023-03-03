// Package key creates asymmetric key pairs.
package key

import (
	"crypto/cipher"

	"go.dedis.ch/kyber/v3"
)

// Generator is a type that needs to implement a special case in order
// to correctly choose a key.
type Generator interface {
	NewKey(random cipher.Stream) kyber.Scalar
}

// Suite represents the list of functionalities needed by this package.
type Suite interface {
	kyber.Group
	kyber.Random
}

// Pair represents a public/private keypair together with the
// ciphersuite the key was generated from.
type Pair struct {
	Public  kyber.Point  // Public key
	Private kyber.Scalar // Private key
}

// NewKeyPair directly creates a secret/public key pair
func NewKeyPair(suite Suite) *Pair {
	kp := new(Pair)
	kp.Gen(suite)
	return kp
}

// Gen creates a fresh public/private keypair with the given
// ciphersuite, using a given source of cryptographic randomness. If
// suite implements key.Generator, then suite.NewKey is called
// to generate the private key, otherwise the normal technique
// of choosing a random scalar from the group is used.
func (p *Pair) Gen(suite Suite) {
	random := suite.RandomStream()
	if g, ok := suite.(Generator); ok {
		p.Private = g.NewKey(random)
	} else {
		p.Private = suite.Scalar().Pick(random)
	}
	p.Public = suite.Point().Mul(p.Private, nil)
}
