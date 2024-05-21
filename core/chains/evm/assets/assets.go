package assets

import (
	"database/sql/driver"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/chainlink-common/pkg/assets"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/bytes"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"

	"github.com/shopspring/decimal"
)

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
	if e == nil {
		return "<nil>"
	}
	return assets.Format(e.ToInt(), 18)
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
	if bytes.HasQuotes(data) {
		return e.UnmarshalText(bytes.TrimQuotes(data))
	}
	return assets.ErrNoQuotesForCurrency
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
	return (*ubig.Big)(e).Scan(value)
}

// Value returns the Eth value for serialization to database.
func (e Eth) Value() (driver.Value, error) {
	return (ubig.Big)(e).Value()
}
