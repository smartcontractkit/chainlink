package scalar

import (
	"bytes"
	"crypto/cipher"
	"fmt"
	"io"
	"math/big"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/pkg/errors"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/mod"
)

type Scalar struct{ int *mod.Int }

var _ kyber.Scalar = (*Scalar)(nil)

var order = bn256.Order

func NewScalar() *Scalar { return &Scalar{mod.NewInt64(0, order)} }

func NewScalarInt64(i int64) *Scalar {
	rv := NewScalar()
	_ = rv.int.SetInt64(i)
	return rv
}

func (s *Scalar) MarshalBinary() (data []byte, err error) {
	s.ensureAllocation()
	if s.int.V.Cmp(s.int.M) >= 0 {

		return nil, errors.Errorf("0x%x too large for AltBN128Scalar", s)
	}
	data = s.int.V.Bytes()
	if len(data) > s.MarshalSize() {
		panic(fmt.Sprintf("0x%x too large for AltBN128Scalar", s))
	}

	return append(bytes.Repeat([]byte{0}, s.MarshalSize()-len(data)), data...), nil
}

func (s *Scalar) UnmarshalBinary(data []byte) error {
	if s == nil {
		return errors.Errorf("can't set the value of a nil *Scalar")
	}
	s.ensureAllocation()

	r := big.NewInt(0).SetBytes(data)
	if r.Cmp(order) >= 0 {
		return errors.Errorf("0x%x too large for AltBN128Scalar", s)
	}
	s.int.V = *r
	return nil
}

func (s *Scalar) String() string {
	s.ensureAllocation()
	return fmt.Sprintf("AltBN128Scalar{0x%s}", s.int)
}

var marshalSize = len(order.Bytes())

func (s Scalar) MarshalSize() int {
	return marshalSize
}

func (s *Scalar) MarshalTo(w io.Writer) (int, error) {
	data, err := s.MarshalBinary()
	if err != nil {
		return 0, errors.Wrapf(err, "while marshalling for writing")
	}
	n, err := w.Write(data)
	return n, errors.Wrapf(err, "while writing marshaled value 0x%x", data)
}

func (s *Scalar) UnmarshalFrom(r io.Reader) (int, error) {
	buf := make([]byte, s.MarshalSize())
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return n, errors.Wrap(err, "while reading for unmarshalling")
	}
	return n, errors.Wrapf(s.UnmarshalBinary(buf), "while unmarshalling 0x%x", buf)
}

var zero = big.NewInt(0)

func (s *Scalar) Equal(s2 kyber.Scalar) bool {
	_, ok := s2.(*Scalar)
	if !ok || s == nil || s2 == nil {
		return false
	}
	diff := &NewScalar().Sub(s, s2).(*Scalar).int.V

	return diff.Mod(diff, s.int.M).Cmp(zero) == 0
}

func mustScalar(s, a kyber.Scalar, verb string) *Scalar {
	if aScalar, ok := a.(*Scalar); ok {
		return aScalar
	}
	panic(fmt.Sprintf(
		"attempt to combine %s non AltBN-128 Scalar %s with AltBN-128 Scalar "+
			"operation %s", a, s, verb),
	)
}

func (s *Scalar) Set(a kyber.Scalar) kyber.Scalar {
	aScalar := mustScalar(s, a, "set")
	_ = s.int.Set(aScalar.int)
	return s
}

func (s *Scalar) Clone() kyber.Scalar {
	return &Scalar{s.int.Clone().(*mod.Int)}
}

func (s *Scalar) ensureAllocation() {
	if s == nil {
		panic("attempt to ensure allocation on nil *Scalar")
	}
	if s.int == nil {
		s.int = mod.NewInt64(0, order)
	}
}

func (s *Scalar) SetInt64(v int64) kyber.Scalar {
	s.ensureAllocation()
	_ = s.int.SetInt64(v)
	return s
}

func (s *Scalar) Zero() kyber.Scalar {
	return s.SetInt64(0)
}

func (s *Scalar) Add(a kyber.Scalar, b kyber.Scalar) kyber.Scalar {
	aScalar := mustScalar(s, a, "sum")
	bScalar := mustScalar(s, b, "sum")
	s.ensureAllocation()
	s.int.Add(aScalar.int, bScalar.int)
	return s
}

func (s *Scalar) Sub(a kyber.Scalar, b kyber.Scalar) kyber.Scalar {
	aScalar := mustScalar(s, a, "subtract")
	bScalar := mustScalar(s, b, "subtract")
	s.ensureAllocation()
	_ = s.int.Sub(aScalar.int, bScalar.int)
	return s
}

func (s *Scalar) Neg(a kyber.Scalar) kyber.Scalar {
	panic("not implemented")
}

func (s *Scalar) One() kyber.Scalar {
	s.ensureAllocation()
	_ = s.int.SetInt64(1)
	return s
}

func (s *Scalar) Mul(a kyber.Scalar, b kyber.Scalar) kyber.Scalar {
	aScalar := mustScalar(s, a, "multiply")
	bScalar := mustScalar(s, b, "multiply")
	s.ensureAllocation()
	_ = s.int.Mul(aScalar.int, bScalar.int)
	return s
}

func (s *Scalar) Div(a kyber.Scalar, b kyber.Scalar) kyber.Scalar {
	aScalar := mustScalar(s, a, "multiply")
	bScalar := mustScalar(s, b, "multiply")
	return s.Mul(aScalar, NewScalar().Inv(bScalar))
}

func (s *Scalar) Inv(a kyber.Scalar) kyber.Scalar {
	aScalar := mustScalar(s, a, "invert")
	s.ensureAllocation()
	_ = s.int.Inv(aScalar.int)
	return s
}

func (s *Scalar) Pick(rand cipher.Stream) kyber.Scalar {
	s.ensureAllocation()
	_ = s.int.Pick(rand)
	return s
}

func (s *Scalar) SetBytes(data []byte) kyber.Scalar {
	s.int.V.SetBytes(data)
	return s
}

func (s *Scalar) Big() *big.Int {
	return big.NewInt(0).Set(&s.int.V)
}
