package ocrkey

import (
	"crypto/ed25519"

	"github.com/pkg/errors"
)

type offChainPrivateKey ed25519.PrivateKey

// Sign returns the signature on msgHash with k
func (k *offChainPrivateKey) Sign(msg []byte) ([]byte, error) {
	if k == nil {
		return nil, errors.Errorf("attempt to sign with nil key")
	}
	return ed25519.Sign(ed25519.PrivateKey(*k), msg), nil
}

// PublicKey returns the public key which commits to k
func (k *offChainPrivateKey) PublicKey() OffChainPublicKey {
	return OffChainPublicKey(ed25519.PrivateKey(*k).Public().(ed25519.PublicKey))
}
