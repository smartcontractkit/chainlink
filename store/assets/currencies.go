package assets

import (
	"fmt"
	"math/big"
)

// Link contains a field to represent the smallest units of LINK
type Link big.Int

// NewLink returns a new struct to represent LINK from it's smallest unit
func NewLink(w int64) *Link {
	return (*Link)(big.NewInt(w))
}

func (l *Link) String() string {
	b := &big.Int{}
	b.SetString("1000000000000000000", 10)
	r := &big.Rat{}
	r.SetFrac((*big.Int)(l), b)
	return fmt.Sprintf("%v", r.FloatString(18))
}

// SetInt64 delegates to *big.Int.SetInt64
func (l *Link) SetInt64(w int64) *Link {
	return (*Link)((*big.Int)(l).SetInt64(w))
}

// SetString delegates to *big.Int.SetString
func (l *Link) SetString(s string, base int) (*Link, bool) {
	w, ok := (*big.Int)(l).SetString(s, base)
	return (*Link)(w), ok
}

// MarshalText implements the encoding.TextMarshaler interface.
func (l *Link) MarshalText() ([]byte, error) {
	return (*big.Int)(l).MarshalText()
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (l *Link) UnmarshalText(text []byte) error {
	if _, ok := l.SetString(string(text), 10); !ok {
		return fmt.Errorf("assets: cannot unmarshal %q into a *assets.Link", text)
	}
	return nil
}
