package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/chainlink/core/logger"

	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	// SignatureLength is the length of the signature in bytes: v = 1, r = 32, s
	// = 32; v + r + s = 65
	SignatureLength = 65
)

// Signature is a byte array fixed to the size of a signature
type Signature [SignatureLength]byte

// NewSignature returns a new Signature
func NewSignature(s string) (Signature, error) {
	bytes := common.FromHex(s)
	return BytesToSignature(bytes), nil
}

// BytesToSignature converts an arbitrary length byte array to a Signature
func BytesToSignature(b []byte) Signature {
	var s Signature
	s.SetBytes(b)
	return s
}

// Bytes returns the raw bytes
func (s Signature) Bytes() []byte { return s[:] }

// Big returns a big.Int representation
func (s Signature) Big() *big.Int { return new(big.Int).SetBytes(s[:]) }

// Hex returns a hexadecimal string
func (s Signature) Hex() string { return hexutil.Encode(s[:]) }

// String implements the stringer interface and is used also by the logger.
func (s Signature) String() string {
	return s.Hex()
}

// Format implements fmt.Formatter
func (s Signature) Format(state fmt.State, c rune) {
	_, err := fmt.Fprintf(state, "%"+string(c), s.String())
	logger.ErrorIf(err, "failed when format signature to state")
}

// SetBytes assigns the byte array to the signature
func (s *Signature) SetBytes(b []byte) {
	if len(b) > len(s) {
		b = b[len(b)-SignatureLength:]
	}

	copy(s[SignatureLength-len(b):], b)
}

// UnmarshalText parses the signature from a hexadecimal representation
func (s *Signature) UnmarshalText(input []byte) error {
	var err error
	*s, err = NewSignature(string(input))
	return err
}

// MarshalText encodes the signature in hexadecimal
func (s Signature) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

// UnmarshalJSON parses a signature from a JSON string
func (s *Signature) UnmarshalJSON(input []byte) error {
	input = utils.RemoveQuotes(input)
	return s.UnmarshalText(input)
}

// MarshalJSON prints the signature as a hexadecimal encoded string
func (s Signature) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// Value returns this instance serialized for database storage.
func (s Signature) Value() (driver.Value, error) {
	return s.String(), nil
}

// Scan reads the database value and returns an instance.
func (s *Signature) Scan(value interface{}) error {
	temp, ok := value.(string)
	if !ok {
		return fmt.Errorf("unable to convert %v of %T to Signature", value, value)
	}

	newSig, err := NewSignature(temp)
	*s = newSig
	return err
}
