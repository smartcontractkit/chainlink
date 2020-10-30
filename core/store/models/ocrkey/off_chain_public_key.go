package ocrkey

import (
	"crypto/ed25519"
	"encoding/hex"
)

type OffChainPublicKey ed25519.PublicKey

func (ocpk OffChainPublicKey) MarshalJSON() ([]byte, error) {
	return []byte(hex.EncodeToString(ocpk)), nil
}

func (ocpk *OffChainPublicKey) UnmarshalJSON(bs []byte) error {
	var err error
	*ocpk, err = hex.DecodeString(string(bs))
	return err
}
