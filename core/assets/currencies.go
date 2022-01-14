package assets

import (
	"database/sql/driver"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

var ErrNoQuotesForCurrency = errors.New("cannot unmarshal json.Number into currency")

// getDenominator returns 10**precision.
func getDenominator(precision int) *big.Int {
	x := big.NewInt(10)
	return new(big.Int).Exp(x, big.NewInt(int64(precision)), nil)
}

func format(i *big.Int, precision int) string {
	r := big.NewRat(1, 1).SetFrac(i, getDenominator(precision))
	return fmt.Sprintf("%v", r.FloatString(precision))
}

// Link contains a field to represent the smallest units of LINK
type Link big.Int

// NewLinkFromJuels returns a new struct to represent LINK from it's smallest unit
func NewLinkFromJuels(w int64) *Link {
	return (*Link)(big.NewInt(w))
}

// String returns Link formatted as a string.
func (l *Link) String() string {
	if l == nil {
		return "0"
	}
	return fmt.Sprintf("%v", (*big.Int)(l))
}

// Link returns Link formatted as a string, in LINK units
func (l *Link) Link() string {
	if l == nil {
		return "0"
	}
	return format((*big.Int)(l), 18)
}

// SetInt64 delegates to *big.Int.SetInt64
func (l *Link) SetInt64(w int64) *Link {
	return (*Link)((*big.Int)(l).SetInt64(w))
}

// ToInt returns the Link value as a *big.Int.
func (l *Link) ToInt() *big.Int {
	return (*big.Int)(l)
}

// ToHash returns a 32 byte representation of this value
func (l *Link) ToHash() common.Hash {
	return common.BigToHash((*big.Int)(l))
}

// Set delegates to *big.Int.Set
func (l *Link) Set(x *Link) *Link {
	il := (*big.Int)(l)
	ix := (*big.Int)(x)

	w := il.Set(ix)
	return (*Link)(w)
}

// SetString delegates to *big.Int.SetString
func (l *Link) SetString(s string, base int) (*Link, bool) {
	w, ok := (*big.Int)(l).SetString(s, base)
	return (*Link)(w), ok
}

// Cmp defers to big.Int Cmp
func (l *Link) Cmp(y *Link) int {
	return (*big.Int)(l).Cmp((*big.Int)(y))
}

// Add defers to big.Int Add
func (l *Link) Add(x, y *Link) *Link {
	il := (*big.Int)(l)
	ix := (*big.Int)(x)
	iy := (*big.Int)(y)

	return (*Link)(il.Add(ix, iy))
}

// Text defers to big.Int Text
func (l *Link) Text(base int) string {
	return (*big.Int)(l).Text(base)
}

// MarshalText implements the encoding.TextMarshaler interface.
func (l *Link) MarshalText() ([]byte, error) {
	return (*big.Int)(l).MarshalText()
}

// MarshalJSON implements the json.Marshaler interface.
func (l Link) MarshalJSON() ([]byte, error) {
	value, err := l.MarshalText()
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf(`"%s"`, value)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (l *Link) UnmarshalJSON(data []byte) error {
	if utils.IsQuoted(data) {
		return l.UnmarshalText(utils.RemoveQuotes(data))
	}
	return ErrNoQuotesForCurrency
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (l *Link) UnmarshalText(text []byte) error {
	if _, ok := l.SetString(string(text), 10); !ok {
		return fmt.Errorf("assets: cannot unmarshal %q into a *assets.Link", text)
	}
	return nil
}

// IsZero returns true when the value is 0 and false otherwise
func (l *Link) IsZero() bool {
	zero := big.NewInt(0)
	return (*big.Int)(l).Cmp(zero) == 0
}

// Symbol returns LINK
func (*Link) Symbol() string {
	return "LINK"
}

// Value returns the Link value for serialization to database.
func (l Link) Value() (driver.Value, error) {
	b := (big.Int)(l)
	return b.String(), nil
}

// Scan reads the database value and returns an instance.
func (l *Link) Scan(value interface{}) error {
	switch v := value.(type) {
	case string:
		decoded, ok := l.SetString(v, 10)
		if !ok {
			return fmt.Errorf("unable to set string %v of %T to base 10 big.Int for Link", value, value)
		}
		*l = *decoded
	case []uint8:
		// The SQL library returns numeric() types as []uint8 of the string representation
		decoded, ok := l.SetString(string(v), 10)
		if !ok {
			return fmt.Errorf("unable to set string %v of %T to base 10 big.Int for Link", value, value)
		}
		*l = *decoded
	case int64:
		return fmt.Errorf("unable to convert %v of %T to Link, is the sql type set to varchar?", value, value)
	default:
		return fmt.Errorf("unable to convert %v of %T to Link", value, value)
	}

	return nil
}

// Eth contains a field to represent the smallest units of ETH
type Eth big.Int

// NewEth returns a new struct to represent ETH from it's smallest unit (wei)
func NewEth(w int64) *Eth {
	return (*Eth)(big.NewInt(w))
}

// NewEthValue returns a new struct to represent ETH from it's smallest unit (wei)
func NewEthValue(w int64) Eth {
	eth := NewEth(w)
	return *eth
}

// NewEthValueS returns a new struct to represent ETH from a string value of Eth (not wei)
// the underlying value is still wei
func NewEthValueS(s string) (Eth, error) {
	e, err := decimal.NewFromString(s)
	if err != nil {
		return Eth{}, err
	}
	w := e.Mul(decimal.RequireFromString("10").Pow(decimal.RequireFromString("18")))
	return *(*Eth)(w.BigInt()), nil
}

// Cmp delegates to *big.Int.Cmp
func (e *Eth) Cmp(y *Eth) int {
	return e.ToInt().Cmp(y.ToInt())
}

func (e *Eth) String() string {
	return format(e.ToInt(), 18)
}

// SetInt64 delegates to *big.Int.SetInt64
func (e *Eth) SetInt64(w int64) *Eth {
	return (*Eth)(e.ToInt().SetInt64(w))
}

// SetString delegates to *big.Int.SetString
func (e *Eth) SetString(s string, base int) (*Eth, bool) {
	w, ok := e.ToInt().SetString(s, base)
	return (*Eth)(w), ok
}

// MarshalJSON implements the json.Marshaler interface.
func (e Eth) MarshalJSON() ([]byte, error) {
	value, err := e.MarshalText()
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf(`"%s"`, value)), nil
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e *Eth) MarshalText() ([]byte, error) {
	return e.ToInt().MarshalText()
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (e *Eth) UnmarshalJSON(data []byte) error {
	if utils.IsQuoted(data) {
		return e.UnmarshalText(utils.RemoveQuotes(data))
	}
	return ErrNoQuotesForCurrency
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *Eth) UnmarshalText(text []byte) error {
	if _, ok := e.SetString(string(text), 10); !ok {
		return fmt.Errorf("assets: cannot unmarshal %q into a *assets.Eth", text)
	}
	return nil
}

// IsZero returns true when the value is 0 and false otherwise
func (e *Eth) IsZero() bool {
	zero := big.NewInt(0)
	return e.ToInt().Cmp(zero) == 0
}

// Symbol returns ETH
func (*Eth) Symbol() string {
	return "ETH"
}

// ToInt returns the Eth value as a *big.Int.
func (e *Eth) ToInt() *big.Int {
	return (*big.Int)(e)
}

// Scan reads the database value and returns an instance.
func (e *Eth) Scan(value interface{}) error {
	return (*utils.Big)(e).Scan(value)
}

// Value returns the Eth value for serialization to database.
func (e Eth) Value() (driver.Value, error) {
	return (utils.Big)(e).Value()
}
