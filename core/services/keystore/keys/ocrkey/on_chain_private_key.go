package ocrkey

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
)

type onChainPrivateKey struct {
	pk func() *ecdsa.PrivateKey
}

// Sign returns the signature on msgHash with k
func (k *onChainPrivateKey) Sign(msg []byte) (signature []byte, err error) {
	sig, err := crypto.Sign(onChainHash(msg), k.pk())
	return sig, err
}

func (k onChainPrivateKey) Address() OnChainSigningAddress {
	return OnChainSigningAddress(crypto.PubkeyToAddress(k.pk().PublicKey))
}

func onChainHash(msg []byte) []byte {
	return crypto.Keccak256(msg)
}
