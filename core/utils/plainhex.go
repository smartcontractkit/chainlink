package utils

import (
	"encoding/hex"
	"encoding/json"
	"reflect"

	"github.com/pkg/errors"
)

// Similar to go-ethereum's hexutil.Bytes but does not assume a 0x prefix.

// PlainHexBytes marshals/unmarshals as a JSON string without a 0x prefix.
// The empty slice marshals as "".
type PlainHexBytes []byte

// MarshalText implements encoding.TextMarshaler
func (b PlainHexBytes) MarshalText() ([]byte, error) {
	result := make([]byte, len(b)*2)
	hex.Encode(result, b)
	return result, nil
}

func (b PlainHexBytes) String() string {
	return hex.EncodeToString(b)
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *PlainHexBytes) UnmarshalJSON(input []byte) (err error) {
	if !isString(input) {
		return &json.UnmarshalTypeError{Value: "non-string", Type: reflect.TypeOf((PlainHexBytes)(nil))}
	}
	err = b.UnmarshalText(input[1 : len(input)-1])
	return errors.Wrap(err, "UnmarshalJSON failed")
}

func isString(input []byte) bool {
	return len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"'
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (b *PlainHexBytes) UnmarshalText(input []byte) error {
	raw, err := checkText(input, true)
	if err != nil {
		return errors.Wrap(err, "UnmarshalText failed")
	}
	dec := make([]byte, len(raw)/2)
	if _, err = hex.Decode(dec, raw); err != nil {
		return errors.Wrap(err, "UnmarshalText failed")
	}
	*b = dec
	return nil
}

func checkText(input []byte, wantPrefix bool) ([]byte, error) {
	if len(input) == 0 {
		return nil, nil // empty strings are allowed
	}
	if len(input)%2 != 0 {
		return nil, errors.New("odd length")
	}
	return input, nil
}
