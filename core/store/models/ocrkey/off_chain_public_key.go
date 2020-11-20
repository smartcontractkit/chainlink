package ocrkey

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
)

type OffChainPublicKey ed25519.PublicKey

func (ocpk OffChainPublicKey) String() string {
	return hex.EncodeToString(ocpk)
}

func (ocpk OffChainPublicKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(ocpk.String())
}

func (ocpk *OffChainPublicKey) UnmarshalJSON(input []byte) error {
	var hexString string
	var result OffChainPublicKey

	if err := json.Unmarshal(input, &hexString); err != nil {
		return err
	}
	result, err := hex.DecodeString(hexString)
	if err != nil {
		return err
	}
	copy(result[:], result[:common.AddressLength])
	*ocpk = result
	return nil
}
