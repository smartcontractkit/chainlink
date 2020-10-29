package ocrkey

import (
	"crypto/ed25519"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type OffChainPublicKey ed25519.PublicKey

func (ocpk OffChainPublicKey) String() string {
	return hexutil.Encode(ocpk[:])
}
