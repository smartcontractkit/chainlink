package models

import (
	"database/sql/driver"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// Big stores large integers and can deserialize a variety of inputs.
type Big big.Int

// NewBig constructs a Big from *big.Int.
func NewBig(i *big.Int) *Big {
	if i != nil {
		b := Big(*i)
		return &b
	}
	return nil
}

// MarshalText marshals this instance to base 10 number as string.
func (b *Big) MarshalText() ([]byte, error) {
	return []byte((*big.Int)(b).Text(10)), nil
}

// MarshalJSON marshals this instance to base 10 number as string.
func (b *Big) MarshalJSON() ([]byte, error) {
	return b.MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (b *Big) UnmarshalText(input []byte) error {
	input = utils.RemoveQuotes(input)
	str := string(input)
	if utils.HasHexPrefix(str) {
		decoded, err := hexutil.DecodeBig(str)
		if err != nil {
			return err
		}
		*b = Big(*decoded)
		return nil
	}

	_, ok := b.setString(str, 10)
	if !ok {
		return fmt.Errorf("Unable to convert %s to Big", str)
	}
	return nil
}

func (b *Big) setBytes(value []uint8) *Big {
	w := (*big.Int)(b).SetBytes(value)
	return (*Big)(w)
}

func (b *Big) setString(s string, base int) (*Big, bool) {
	w, ok := (*big.Int)(b).SetString(s, base)
	return (*Big)(w), ok
}

// UnmarshalJSON implements encoding.JSONUnmarshaler.
func (b *Big) UnmarshalJSON(input []byte) error {
	return b.UnmarshalText(input)
}

// Value returns this instance serialized for database storage.
func (b Big) Value() (driver.Value, error) {
	return b.String(), nil
}

// Scan reads the database value and returns an instance.
func (b *Big) Scan(value interface{}) error {
	switch v := value.(type) {
	case string:
		decoded, ok := b.setString(v, 10)
		if !ok {
			return fmt.Errorf("Unable to set string %v of %T to base 10 big.Int for Big", value, value)
		}
		*b = *decoded
	case []uint8:
		// The SQL library returns numeric() types as []uint8 of the string representation
		decoded, ok := b.setString(string(v), 10)
		if !ok {
			return fmt.Errorf("Unable to set string %v of %T to base 10 big.Int for Big", value, value)
		}
		*b = *decoded
	default:
		return fmt.Errorf("Unable to convert %v of %T to Big", value, value)
	}

	return nil
}

// ToInt converts b to a big.Int.
func (b *Big) ToInt() *big.Int {
	return (*big.Int)(b)
}

// String returns the base 10 encoding of b.
func (b *Big) String() string {
	return b.ToInt().Text(10)
}

// Hex returns the hex encoding of b.
func (b *Big) Hex() string {
	return hexutil.EncodeBig(b.ToInt())
}
