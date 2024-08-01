package felt

import (
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/consensys/gnark-crypto/ecc/stark-curve/fp"
	"github.com/fxamacker/cbor/v2"
)

type Felt struct {
	val fp.Element
}

func NewFelt(element *fp.Element) *Felt {
	return &Felt{
		val: *element,
	}
}

const (
	Base16 = 16
	Base10 = 10
)

const (
	Limbs = fp.Limbs // number of 64 bits words needed to represent a Element
	Bits  = fp.Bits  // number of bits needed to represent a Element
	Bytes = fp.Bytes // number of bytes needed to represent a Element
)

// Zero felt constant
var Zero = Felt{}

var bigIntPool = sync.Pool{
	New: func() interface{} {
		return new(big.Int)
	},
}

// Impl returns the underlying field element type
func (z *Felt) Impl() *fp.Element {
	return &z.val
}

// UnmarshalJSON accepts numbers and strings as input.
// See Element.SetString for valid prefixes (0x, 0b, ...).
// If there is an error, we try to explicitly unmarshal from hex before
// returning an error. This implementation is based on [gnark-crypto]'s UnmarshalJSON.
//
// [gnark-crypto]: https://github.com/ConsenSys/gnark-crypto/blob/master/ecc/stark-curve/fp/element.go
func (z *Felt) UnmarshalJSON(data []byte) error {
	s := string(data)
	if len(s) > fp.Bits*3 {
		return errors.New("value too large (max = Element.Bits * 3)")
	}

	// we accept numbers and strings, remove leading and trailing quotes if any
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}

	_, err := z.SetString(s)
	return err
}

// MarshalJSON forwards the call to underlying field element implementation
func (z *Felt) MarshalJSON() ([]byte, error) {
	return []byte("\"" + z.String() + "\""), nil
}

// SetBytes forwards the call to underlying field element implementation
func (z *Felt) SetBytes(e []byte) *Felt {
	z.val.SetBytes(e)
	return z
}

// SetString forwards the call to underlying field element implementation
func (z *Felt) SetString(number string) (*Felt, error) {
	// get temporary big int from the pool
	vv := bigIntPool.Get().(*big.Int)
	// release object into pool
	defer bigIntPool.Put(vv)

	if _, ok := vv.SetString(number, 0); !ok {
		if _, ok := vv.SetString(number, Base16); !ok {
			return z, errors.New("can't parse into a big.Int: " + number)
		}
	}

	if vv.BitLen() > fp.Bits {
		return z, errors.New("can't fit in felt: " + number)
	}

	var bytes [32]byte
	vv.FillBytes(bytes[:])
	return z, z.val.SetBytesCanonical(bytes[:])
}

// SetUint64 forwards the call to underlying field element implementation
func (z *Felt) SetUint64(v uint64) *Felt {
	z.val.SetUint64(v)
	return z
}

// SetRandom forwards the call to underlying field element implementation
func (z *Felt) SetRandom() (*Felt, error) {
	_, err := z.val.SetRandom()
	return z, err
}

// String forwards the call to underlying field element implementation
func (z *Felt) String() string {
	return "0x" + z.val.Text(Base16)
}

// ShortString prints the felt to a string in a shortened format
func (z *Felt) ShortString() string {
	shortFelt := 8
	hex := z.val.Text(Base16)

	if len(hex) <= shortFelt {
		return fmt.Sprintf("0x%s", hex)
	}
	return fmt.Sprintf("0x%s...%s", hex[:4], hex[len(hex)-4:])
}

// Text forwards the call to underlying field element implementation
func (z *Felt) Text(base int) string {
	return z.val.Text(base)
}

// Equal forwards the call to underlying field element implementation
func (z *Felt) Equal(x *Felt) bool {
	return z.val.Equal(&x.val)
}

// Marshal forwards the call to underlying field element implementation
func (z *Felt) Marshal() []byte {
	return z.val.Marshal()
}

// Bytes forwards the call to underlying field element implementation
func (z *Felt) Bytes() [32]byte {
	return z.val.Bytes()
}

// IsOne forwards the call to underlying field element implementation
func (z *Felt) IsOne() bool {
	return z.val.IsOne()
}

// IsZero forwards the call to underlying field element implementation
func (z *Felt) IsZero() bool {
	return z.val.IsZero()
}

// Add forwards the call to underlying field element implementation
func (z *Felt) Add(x, y *Felt) *Felt {
	z.val.Add(&x.val, &y.val)
	return z
}

// Halve forwards the call to underlying field element implementation
func (z *Felt) Halve() {
	z.val.Halve()
}

// MarshalCBOR lets Felt be encoded in CBOR format with private `val`
func (z *Felt) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(z.val)
}

// UnmarshalCBOR lets Felt be decoded from CBOR format with private `val`
func (z *Felt) UnmarshalCBOR(data []byte) error {
	return cbor.Unmarshal(data, &z.val)
}

// Bits forwards the call to underlying field element implementation
func (z *Felt) Bits() [4]uint64 {
	return z.val.Bits()
}

// BigInt forwards the call to underlying field element implementation
func (z *Felt) BigInt(res *big.Int) *big.Int {
	return z.val.BigInt(res)
}

// Set forwards the call to underlying field element implementation
func (z *Felt) Set(x *Felt) *Felt {
	z.val.Set(&x.val)
	return z
}

// Double forwards the call to underlying field element implementation
func (z *Felt) Double(x *Felt) *Felt {
	z.val.Double(&x.val)
	return z
}

// Sub forwards the call to underlying field element implementation
func (z *Felt) Sub(x, y *Felt) *Felt {
	z.val.Sub(&x.val, &y.val)
	return z
}

// Exp forwards the call to underlying field element implementation
func (z *Felt) Exp(x *Felt, y *big.Int) *Felt {
	z.val.Exp(x.val, y)
	return z
}

// Mul forwards the call to underlying field element implementation
func (z *Felt) Mul(x, y *Felt) *Felt {
	z.val.Mul(&x.val, &y.val)
	return z
}

// Cmp forwards the call to underlying field element implementation
func (z *Felt) Cmp(x *Felt) int {
	return z.val.Cmp(&x.val)
}
