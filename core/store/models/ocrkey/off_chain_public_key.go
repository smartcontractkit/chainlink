package ocrkey

import (
	"crypto/ed25519"
	"encoding/hex"
)

type OffChainPublicKey ed25519.PublicKey

func (ocpk OffChainPublicKey) String() string {
	return hex.EncodeToString(ocpk[:])
}
