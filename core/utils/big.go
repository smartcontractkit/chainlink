package utils

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	bigmath "github.com/smartcontractkit/chainlink/core/utils/big_math"
)

const base10 = 10

// BigFloat accepts both string and float JSON values.
type BigFloat big.Float

// MarshalJSON implements the json.Marshaler interface.
func (b BigFloat) MarshalJSON() ([]byte, error) {
	var j = big.Float(b)
	return json.Marshal(&j)
}

// UnmarshalJSON implements the json.Unmarshal interface.
func (b *BigFloat) UnmarshalJSON(buf []byte) error {
	var n json.Number
	if err := json.Unmarshal(buf, &n); err == nil {
		f, _, err := new(big.Float).Parse(n.String(), 0)
		if err != nil {
			return err
		}
		*b = BigFloat(*f)
		return nil
	}
	var bf big.Float
	if err := json.Unmarshal(buf, &bf); err != nil {
		return err
	}
	*b = BigFloat(bf)
	return nil
}

// Value returns the big.Float value.
func (b *BigFloat) Value() *big.Float {
	return (*big.Float)(b)
}

// Big stores large integers and can deserialize a variety of inputs.
type Big big.Int

// NewBig constructs a Big from *big.Int.
func NewBig(i *big.Int) *Big {
	if i != nil {
		var b big.Int
		b.Set(i)
		return (*Big)(&b)
	}
	return nil
}

// NewBigI constructs a Big from int64.
func NewBigI(i int64) *Big {
	return NewBig(big.NewInt(i))
}

// MarshalText marshals this instance to base 10 number as string.
func (b Big) MarshalText() ([]byte, error) {
	return []byte((*big.Int)(&b).Text(base10)), nil
}

// MarshalJSON marshals this instance to base 10 number as string.
func (b Big) MarshalJSON() ([]byte, error) {
	text, err := b.MarshalText()
	if err != nil {
		return nil, err
	}
	return json.Marshal(string(text))
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (b *Big) UnmarshalText(input []byte) error {
	input = RemoveQuotes(input)
	str := string(input)
	if HasHexPrefix(str) {
		decoded, err := hexutil.DecodeBig(str)
		if err != nil {
			return err
		}
		*b = Big(*decoded)
		return nil
	}

	_, ok := b.setString(str, 10)
	if !ok {
		return fmt.Errorf("unable to convert %s to Big", str)
	}
	return nil
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
			return fmt.Errorf("unable to set string %v of %T to base 10 big.Int for Big", value, value)
		}
		*b = *decoded
	case []uint8:
		// The SQL library returns numeric() types as []uint8 of the string representation
		decoded, ok := b.setString(string(v), 10)
		if !ok {
			return fmt.Errorf("unable to set string %v of %T to base 10 big.Int for Big", value, value)
		}
		*b = *decoded
	default:
		return fmt.Errorf("unable to convert %v of %T to Big", value, value)
	}

	return nil
}

// ToInt converts b to a big.Int.
func (b *Big) ToInt() *big.Int {
	return (*big.Int)(b)
}

// String returns the base 10 encoding of b.
func (b *Big) String() string {
	return b.ToInt().String()
}

// Bytes returns the absolute value of b as a big-endian byte slice.
func (b *Big) Hex() string {
	return hexutil.EncodeBig(b.ToInt())
}

// Bytes returns the
func (b *Big) Bytes() []byte {
	return b.ToInt().Bytes()
}

// Cmp compares b and c as big.Ints.
func (b *Big) Cmp(c *Big) int {
	return b.ToInt().Cmp(c.ToInt())
}

// Equal returns true if c is equal according to Cmp.
func (b *Big) Equal(c *Big) bool {
	return b.Cmp(c) == 0
}

// Int64 casts b as an int64 type
func (b *Big) Int64() int64 {
	return b.ToInt().Int64()
}

// Add returns the sum of b and c
func (b *Big) Add(c interface{}) *Big {
	return NewBig(bigmath.Add(b, c))
}

// Sub returns the differencs between b and c
func (b *Big) Sub(c interface{}) *Big {
	return NewBig(bigmath.Sub(b, c))
}

// Sub returns b % c
func (b *Big) Mod(c interface{}) *Big {
	return NewBig(bigmath.Mod(b, c))
}
