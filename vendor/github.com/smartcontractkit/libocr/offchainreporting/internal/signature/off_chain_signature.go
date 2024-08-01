package signature

import (
	"bytes"
	"crypto/ed25519"

	"github.com/pkg/errors"
)

// OffchainPublicKey is the public key used to cryptographically identify an
// oracle in inter-oracle communications.
type OffchainPublicKey ed25519.PublicKey

// Equal returns true iff k and k2 represent the same public key
func (k OffchainPublicKey) Equal(k2 OffchainPublicKey) bool {
	return bytes.Equal([]byte(ed25519.PublicKey(k)), []byte(ed25519.PublicKey(k2)))
}

// Verify returns true iff signature is a valid signature by k on msg
func (k OffchainPublicKey) Verify(msg, signature []byte) bool {
	return ed25519.Verify(ed25519.PublicKey(k), msg, signature)
}

// OffchainPrivateKey is the secret key oracles use to sign messages destined
// for off-chain verification by other oracles
type OffchainPrivateKey ed25519.PrivateKey

// Sign returns the signature on msgHash with k
func (k *OffchainPrivateKey) Sign(msg []byte) ([]byte, error) {
	if k == nil {
		return nil, errors.Errorf("attempt to sign with nil key")
	}
	return ed25519.Sign(ed25519.PrivateKey(*k), msg), nil
}

// PublicKey returns the public key which commits to k
func (k *OffchainPrivateKey) PublicKey() OffchainPublicKey {
	return OffchainPublicKey(ed25519.PrivateKey(*k).Public().(ed25519.PublicKey))
}
