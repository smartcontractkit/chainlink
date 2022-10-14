package assets

import (
	"database/sql/driver"
	"fmt"
	"math/big"
	"strings"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"golang.org/x/exp/constraints"

	"github.com/smartcontractkit/chainlink/core/utils"
	bigmath "github.com/smartcontractkit/chainlink/core/utils/big_math"
)

const (
	// canonical unit suffixes
	wei   = "wei"
	kwei  = "kwei"
	mwei  = "mwei"
	gwei  = "gwei"
	micro = "micro"
	milli = "milli"
	eth   = "ether"
	keth  = "kether"
	meth  = "mether"
	geth  = "gether"
	teth  = "tether"
)

func suffixExp(suf string) int32 {
	switch suf {
	default:
		panic("unrecognized suffix: " + suf)
	case wei:
		return 0
	case kwei:
		return 3
	case mwei:
		return 6
	case gwei:
		return 9
	case micro:
		return 12
	case milli:
		return 15
	case eth:
		return 18
	case keth:
		return 21
	case meth:
		return 24
	case geth:
		return 27
	case teth:
		return 30
	}
}

// Wei extends utils.Big to implement encoding.TextMarshaler and
// encoding.TextUnmarshaler with support for unit suffixes, as well as
// additional functions
type Wei utils.Big

func MaxWei(w, x *Wei) *Wei {
	return NewWei(bigmath.Max(w.ToInt(), x.ToInt()))
}

// NewWei constructs a Wei from *big.Int.
func NewWei(i *big.Int) *Wei {
	return (*Wei)(i)
}

func NewWeiI[T constraints.Signed](i T) *Wei {
	return NewWei(big.NewInt(int64(i)))
}

func (w *Wei) Text(suffix string) string {
	switch suffix {
	default: // empty or unknown
		fallthrough
	case wei:
		return w.text(wei, 0)
	case kwei:
		return w.text(kwei, 3)
	case mwei:
		return w.text(mwei, 6)
	case gwei:
		return w.text(gwei, 9)
	case micro:
		return w.text(micro, 12)
	case milli:
		return w.text(milli, 15)
	case eth:
		return w.text(eth, 18)
	case keth:
		return w.text(keth, 21)
	case meth:
		return w.text(meth, 24)
	case geth:
		return w.text(geth, 27)
	case teth:
		return w.text(teth, 30)
	}
}

// text formats w with the given suffix and exponent. As a special case, the suffix is omitted for `0`.
func (w *Wei) text(suf string, exp int32) string {
	d := decimal.NewFromBigInt((*big.Int)(w), -exp)
	if d.IsZero() {
		return "0"
	}
	return fmt.Sprintf("%s %s", d, suf)

}

const u64Eth = 1_000_000_000_000_000_000

var (
	bigKeth = new(big.Int).Mul(big.NewInt(u64Eth), big.NewInt(1_000))
	bigMeth = new(big.Int).Mul(big.NewInt(u64Eth), big.NewInt(1_000_000))
	bigGeth = new(big.Int).Mul(big.NewInt(u64Eth), big.NewInt(1_000_000_000))
	bigTeth = new(big.Int).Mul(big.NewInt(u64Eth), big.NewInt(1_000_000_000_000))
)

func (w *Wei) MarshalText() ([]byte, error) {
	return []byte(w.String()), nil
}

func (w *Wei) String() string {
	b := (*big.Int)(w)
	if b.IsUint64() {
		// <= math.MaxUint64 = 9.223_372_036_854_775_808 eth
		u := b.Uint64()
		switch {
		case u >= u64Eth:
			return w.Text(eth)
		case u >= 1_000_000_000_000_000:
			return w.Text(milli)
		case u >= 1_000_000_000_000:
			return w.Text(micro)
		case u >= 1_000_000_000:
			return w.Text(gwei)
		case u >= 1_000_000:
			return w.Text(mwei)
		case u >= 1_000:
			return w.Text(kwei)
		default:
			return w.Text(wei)
		}
	}
	// > math.MaxUint64 = 9.223_372_036_854_775_808 eth
	if b.Cmp(bigTeth) >= 0 {
		return w.Text(teth)
	}
	if b.Cmp(bigGeth) >= 0 {
		return w.Text(geth)
	}
	if b.Cmp(bigMeth) >= 0 {
		return w.Text(meth)
	}
	if b.Cmp(bigKeth) >= 0 {
		return w.Text(keth)
	}
	return w.Text(eth)
}

func (w *Wei) UnmarshalText(b []byte) error {
	s := string(b)
	for _, suf := range []string{
		teth, geth, meth, keth, eth,
		milli, micro,
		gwei, mwei, kwei, wei,
	} {
		if !strings.HasSuffix(s, suf) {
			continue
		}
		t := strings.TrimSuffix(s, suf)
		t = strings.TrimSuffix(t, " ")
		d, err := decimal.NewFromString(t)
		if err != nil {
			return errors.Wrapf(err, "unable to parse %q", s)
		}
		se := suffixExp(suf)
		if d.IsInteger() {
			m := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(se)), nil)
			*w = (Wei)(*new(big.Int).Mul(d.BigInt(), m))
			return nil
		}

		d = d.Mul(decimal.New(1, se))
		if !d.IsInteger() {
			err := errors.New("maximum precision is wei")
			return errors.Wrapf(err, "unable to parse %q", s)
		}
		*w = (Wei)(*d.BigInt())
		return nil

	}
	// unrecognized or missing suffix
	d, err := decimal.NewFromString(s)
	if err != nil {
		return errors.Wrapf(err, "unable to parse %q", s)
	}
	if d.IsInteger() {
		*w = (Wei)(*d.BigInt())
		return nil
	}
	return errors.Errorf("unable to parse %q", s)
}

func (w *Wei) ToInt() *big.Int {
	return (*big.Int)(w)
}

func (w *Wei) Int64() int64 {
	return w.ToInt().Int64()
}

func (w *Wei) Cmp(y *Wei) int {
	return w.ToInt().Cmp(y.ToInt())
}

func (w *Wei) IsNegative() bool {
	return w.Cmp(NewWeiI(0)) < 0
}

func (w *Wei) IsZero() bool {
	return w.Cmp(NewWeiI(0)) == 0
}

func (w *Wei) Equal(y *Wei) bool {
	return w.Cmp(y) == 0
}

func WeiMax(x, y *Wei) *Wei {
	return NewWei(bigmath.Max(x.ToInt(), y.ToInt()))
}

func WeiMin(x, y *Wei) *Wei {
	return NewWei(bigmath.Min(x.ToInt(), y.ToInt()))
}

// NOTE: Maths functions always return newly allocated number and do not mutate

func (w *Wei) Sub(y *Wei) *Wei {
	result := big.NewInt(0).Sub(w.ToInt(), y.ToInt())
	return NewWei(result)
}

func (w *Wei) Add(y *Wei) *Wei {
	return NewWei(big.NewInt(0).Add(w.ToInt(), y.ToInt()))
}

func (w *Wei) Mul(y *big.Int) *Wei {
	return NewWei(big.NewInt(0).Mul(w.ToInt(), y))
}

func (w *Wei) AddPercentage(percentage uint16) *Wei {
	bumped := new(big.Int)
	bumped.Mul(w.ToInt(), big.NewInt(int64(100+percentage)))
	bumped.Div(bumped, big.NewInt(100))
	return NewWei(bumped)
}

// Scan reads the database value and returns an instance.
func (w *Wei) Scan(value interface{}) error {
	return (*utils.Big)(w).Scan(value)
}

// Value returns this instance serialized for database storage.
func (w Wei) Value() (driver.Value, error) {
	return (utils.Big)(w).Value()
}
