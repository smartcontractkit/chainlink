package ocrkey

import (
	"crypto/ed25519"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/common"
)

const offChainPublicKeyPrefix = "ocroff_"

type OffChainPublicKey ed25519.PublicKey

func (ocpk OffChainPublicKey) String() string {
	return fmt.Sprintf("%s%s", offChainPublicKeyPrefix, ocpk.Raw())
}

func (ocpk OffChainPublicKey) Raw() string {
	return hex.EncodeToString(ocpk)
}

func (ocpk OffChainPublicKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(ocpk.String())
}

func (ocpk *OffChainPublicKey) UnmarshalJSON(input []byte) error {
	var hexString string
	if err := json.Unmarshal(input, &hexString); err != nil {
		return err
	}
	return ocpk.UnmarshalText([]byte(hexString))
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

func (ocpk *OffChainPublicKey) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*ocpk = v
		return nil
	default:
		return errors.Errorf("invalid public key bytes got %T wanted []byte", v)
	}
}

func (ocpk OffChainPublicKey) Value() (driver.Value, error) {
	return []byte(ocpk), nil
}
