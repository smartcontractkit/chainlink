package ocrkey

import (
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/crypto/curve25519"
)

const configPublicKeyPrefix = "ocrcfg_"

// ConfigPublicKey represents the public key for the config decryption keypair
type ConfigPublicKey [curve25519.PointSize]byte

func (cpk ConfigPublicKey) String() string {
	return fmt.Sprintf("%s%s", configPublicKeyPrefix, cpk.Raw())
}

func (cpk ConfigPublicKey) Raw() string {
	return hex.EncodeToString(cpk[:])
}

func (cpk ConfigPublicKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(cpk.String())
}

func (cpk *ConfigPublicKey) UnmarshalJSON(input []byte) error {
	var hexString string
	if err := json.Unmarshal(input, &hexString); err != nil {
		return err
	}

	return cpk.UnmarshalText([]byte(hexString))
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

func (cpk *ConfigPublicKey) UnmarshalText(bs []byte) error {
	input := string(bs)
	if strings.HasPrefix(input, configPublicKeyPrefix) {
		input = string(bs[len(configPublicKeyPrefix):])
	}

	decodedString, err := hex.DecodeString(input)
	if err != nil {
		return err
	}
	var result [curve25519.PointSize]byte
	copy(result[:], decodedString[:curve25519.PointSize])
	*cpk = result
	return nil
}
