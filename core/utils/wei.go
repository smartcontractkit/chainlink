package utils

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
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
	// TODO alternative generic unit suffixes?
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

type Wei Big

// NewWei constructs a Wei from *big.Int.
func NewWei(i *big.Int) *Wei {
	return (*Wei)(i)
}

//TODO accept generic suffixes too?
func (w Wei) Text(suffix string) string {
	switch suffix {
	default: // empty or unknown
		fallthrough
	case wei:
		return (*big.Int)(&w).String()
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

func (w Wei) text(suf string, exp int32) string {
	d := decimal.NewFromBigInt((*big.Int)(&w), -exp)
	return fmt.Sprintf("%s %s", d, suf)

}

const u64Eth = 1_000_000_000_000_000_000

var (
	bigKeth = new(big.Int).Mul(big.NewInt(u64Eth), big.NewInt(1_000))
	bigMeth = new(big.Int).Mul(big.NewInt(u64Eth), big.NewInt(1_000_000))
	bigGeth = new(big.Int).Mul(big.NewInt(u64Eth), big.NewInt(1_000_000_000))
	bigTeth = new(big.Int).Mul(big.NewInt(u64Eth), big.NewInt(1_000_000_000_000))
)

func (w Wei) MarshalText() ([]byte, error) {
	b := (*big.Int)(&w)
	if b.IsUint64() {
		// <= math.MaxUint64 = 9.223_372_036_854_775_808 eth
		u := b.Uint64()
		switch {
		case u >= u64Eth:
			return []byte(w.Text(eth)), nil
		case u >= 1_000_000_000_000_000:
			return []byte(w.Text(milli)), nil
		case u >= 1_000_000_000_000:
			return []byte(w.Text(micro)), nil
		case u >= 1_000_000_000:
			return []byte(w.Text(gwei)), nil
		case u >= 1_000_000:
			return []byte(w.Text(mwei)), nil
		case u >= 1_000:
			return []byte(w.Text(kwei)), nil
		default:
			return []byte(w.Text(wei)), nil
		}
	}
	// > math.MaxUint64 = 9.223_372_036_854_775_808 eth
	if b.Cmp(bigTeth) >= 0 {
		return []byte(w.Text(teth)), nil
	}
	if b.Cmp(bigGeth) >= 0 {
		return []byte(w.Text(geth)), nil
	}
	if b.Cmp(bigMeth) >= 0 {
		return []byte(w.Text(meth)), nil
	}
	if b.Cmp(bigKeth) >= 0 {
		return []byte(w.Text(keth)), nil
	}
	return []byte(w.Text(eth)), nil
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
		if strings.HasSuffix(t, " ") {
			t = t[:len(t)-1]
		}
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
	// no suffix?
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
