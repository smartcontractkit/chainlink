package crypto

import (
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// PublicKey defines a type which can be used for JSON and SQL.
type PublicKey []byte

// PublicKeyFromHex generates a public key from a hex string
func PublicKeyFromHex(hexStr string) (*PublicKey, error) {
	result, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}

	pubKey := PublicKey(result)

	return &pubKey, err
}

func (k PublicKey) String() string {
	return hex.EncodeToString(k)
}

func (k PublicKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(k))
}

func (k *PublicKey) UnmarshalJSON(in []byte) error {
	var hexStr string
	if err := json.Unmarshal(in, &hexStr); err != nil {
		return err
	}

	result, err := hex.DecodeString(hexStr)
	if err != nil {
		return err
	}

	*k = PublicKey(result)
	return nil
}

func (k *PublicKey) Scan(value interface{}) error {
	switch v := value.(type) {
	case nil:
		*k = nil
		return nil
	case []byte:
		*k = v
		return nil
	default:
		return fmt.Errorf("invalid public key bytes got %T wanted []byte", v)
	}
}

func (k PublicKey) Value() (driver.Value, error) {
	return []byte(k), nil
}
