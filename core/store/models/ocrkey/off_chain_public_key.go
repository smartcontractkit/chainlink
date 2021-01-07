package ocrkey

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

const offChainPublicKeyPrefix = "ocroff"

type OffChainPublicKey ed25519.PublicKey

func (ocpk OffChainPublicKey) String() string {
	return fmt.Sprintf("%s%s", offChainPublicKeyPrefix, ocpk.Raw())
}

func (ocpk OffChainPublicKey) Raw() string {
	return hex.EncodeToString(ocpk)
}

func (ocpk OffChainPublicKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(ocpk.Raw())
}

func (ocpk *OffChainPublicKey) UnmarshalText(bs []byte) error {
	input := string(bs)
	if strings.HasPrefix(input, offChainPublicKeyPrefix) {
		input = string(bs[len(offChainPublicKeyPrefix):])
	}
	result, err := hex.DecodeString(input)
	if err != nil {
		return err
	}
	copy(result[:], result[:common.AddressLength])
	*ocpk = result
	return nil
}

func (ocpk *OffChainPublicKey) UnmarshalJSON(input []byte) error {
	var hexString string
	if err := json.Unmarshal(input, &hexString); err != nil {
		return err
	}
	return ocpk.UnmarshalText([]byte(hexString))
}
