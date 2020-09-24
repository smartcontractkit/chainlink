package ocrkey

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
)

type OnChainPublicKey ecdsa.PublicKey

func (k OnChainPublicKey) Address() OnChainSigningAddress {
	return OnChainSigningAddress(crypto.PubkeyToAddress(ecdsa.PublicKey(k)))
}
