package assets

import (
	"fmt"
	"math/big"
)

// Link contains a field to represent the smallest units of LINK
type Link struct {
	wei *big.Int
}

// NewLink returns a new struct to represent LINK from it's smallest unit
func NewLink(w *big.Int) Link {
	return Link{wei: w}
}

func (l Link) String() string {
	b := &big.Int{}
	b.SetString("1000000000000000000", 10)
	r := &big.Rat{}
	r.SetFrac(l.wei, b)
	return fmt.Sprintf("%v", r.FloatString(18))
}

// MarshalText implements the encoding.TextMarshaler interface.
func (l Link) MarshalText() ([]byte, error) {
	return l.wei.MarshalText()
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (l *Link) UnmarshalText(text []byte) error {
	b := &big.Int{}
	if _, ok := b.SetString(string(text), 10); !ok {
		return fmt.Errorf("assets: cannot unmarshal %q into a *assets.Link", text)
	}
	l.wei = b
	return nil
}
