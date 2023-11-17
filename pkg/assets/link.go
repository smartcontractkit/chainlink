package assets

import (
	"database/sql/driver"
	"fmt"
	"math/big"
	"strings"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/bytes"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

var ErrNoQuotesForCurrency = errors.New("cannot unmarshal json.Number into currency")

// getDenominator returns 10**precision.
func getDenominator(precision int) *big.Int {
	x := big.NewInt(10)
	return new(big.Int).Exp(x, big.NewInt(int64(precision)), nil)
}

func Format(i *big.Int, precision int) string {
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
	return Format((*big.Int)(l), 18)
}

// SetInt64 delegates to *big.Int.SetInt64
func (l *Link) SetInt64(w int64) *Link {
	return (*Link)((*big.Int)(l).SetInt64(w))
}

// ToInt returns the Link value as a *big.Int.
func (l *Link) ToInt() *big.Int {
	return (*big.Int)(l)
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

var linkFmtThreshold = (*Link)(new(big.Int).Exp(big.NewInt(10), big.NewInt(12), nil))

// MarshalText implements the encoding.TextMarshaler interface.
func (l *Link) MarshalText() ([]byte, error) {
	if l.Cmp(linkFmtThreshold) >= 0 {
		return []byte(fmt.Sprintf("%s link", decimal.NewFromBigInt(l.ToInt(), -18))), nil
	}
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
	if bytes.HasQuotes(data) {
		return l.UnmarshalText(bytes.TrimQuotes(data))
	}
	return ErrNoQuotesForCurrency
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (l *Link) UnmarshalText(text []byte) error {
	s := string(text)
	if strings.HasSuffix(s, "link") {
		s = strings.TrimSuffix(s, "link")
		s = strings.TrimSuffix(s, " ")
		d, err := decimal.NewFromString(s)
		if err != nil {
			return errors.Wrapf(err, "assets: cannot unmarshal %q into a *assets.Link", text)
		}
		d = d.Mul(decimal.New(1, 18))
		if !d.IsInteger() {
			err := errors.New("maximum precision is juels")
			return errors.Wrapf(err, "assets: cannot unmarshal %q into a *assets.Link", text)
		}
		l.Set((*Link)(d.Rat().Num()))
		return nil
	}
	if strings.HasSuffix(s, "juels") {
		s = strings.TrimSuffix(s, "juels")
		s = strings.TrimSuffix(s, " ")
	}
	if _, ok := l.SetString(s, 10); !ok {
		return errors.Errorf("assets: cannot unmarshal %q into a *assets.Link", text)
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
