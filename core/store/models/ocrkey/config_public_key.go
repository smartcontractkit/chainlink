package ocrkey

import (
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"

	"github.com/pkg/errors"
	"golang.org/x/crypto/curve25519"
)

// ConfigPublicKey represents the public key for the config decryption keypair
type ConfigPublicKey [curve25519.PointSize]byte

func (cpk ConfigPublicKey) String() string {
	return hex.EncodeToString(cpk[:])
}

func (cpk ConfigPublicKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(cpk[:]))
}

func (cpk *ConfigPublicKey) UnmarshalJSON(input []byte) error {
	var result [curve25519.PointSize]byte
	var hexString string
	if err := json.Unmarshal(input, &hexString); err != nil {
		return err
	}

	decodedString, err := hex.DecodeString(hexString)
	if err != nil {
		return err
	}
	copy(result[:], decodedString[:curve25519.PointSize])
	*cpk = result
	return nil
}

// Scan reads the database value and returns an instance.
func (cpk *ConfigPublicKey) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.Errorf("unable to convert %v of type %T to ConfigPublicKey", value, value)
	}
	if len(b) != curve25519.PointSize {
		return errors.Errorf("unable to convert blob 0x%x of length %v to ConfigPublicKey", b, len(b))
	}
	copy(cpk[:], b)
	return nil
}

// Value returns this instance serialized for database storage.
func (cpk ConfigPublicKey) Value() (driver.Value, error) {
	return cpk[:], nil
}
