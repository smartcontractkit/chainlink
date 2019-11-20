package bulletin_board

import (
	"chainlink/core/services/signatures/secp256k1"
	"crypto/cipher"

	"go.dedis.ch/kyber/v3"
)

// A SecretKey is a node's secret identity, used to sign all messages
type SecretKey struct{ K kyber.Scalar }

// A publicKey is a node's public identity
type PublicKey struct{ K kyber.Point }

// Set returns true if k has already been set.
func (k SecretKey) Set() bool {
	return !k.K.Equal(k.K.Clone().Zero())
}

// Assign sets k to kPrime
func (k SecretKey) Assign(kPrime kyber.Scalar) {
	k.K.Set(kPrime)

}

var suite = secp256k1.NewBlakeKeccackSecp256k1()
var generator = suite.Point().Base()
var zero = suite.Scalar()

func (k SecretKey) PublicKey() PublicKey {
	p := generator.Clone()
	return PublicKey{K: p.Mul(k.K, generator)}
}

func PickKey(s cipher.Stream) SecretKey {
	return SecretKey{K: zero.Clone().Pick(s)}
}
