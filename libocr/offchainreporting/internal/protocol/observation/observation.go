package observation

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"
)

type Observation struct{ v *big.Int }

type Observations []Observation

var i = big.NewInt

var MaxObservation = i(0).Sub(i(0).Lsh(i(1), 191), i(1))      var MinObservation = i(0).Sub(i(0).Neg(MaxObservation), i(1)) 
func tooLarge(o *big.Int) error {
	return errors.Errorf("value won't fit in int192: 0x%x", o)
}

func MakeObservation(w types.Observation) (Observation, error) {
	v := (*big.Int)(w)
				if v == nil {
		v = big.NewInt(0)
	}
	if v.Cmp(MaxObservation) > 0 || v.Cmp(MinObservation) < 0 {
		return Observation{}, tooLarge(v)
	}
	return Observation{v}, nil
}

func (o Observation) RawObservation() *big.Int { return o.v }

func (o Observation) Less(o2 Observation) bool { return o.v.Cmp(o2.v) < 0 }

func (o Observation) IsMissingValue() bool { return o.v == nil }

func (o Observation) GoEthereumValue() *big.Int { return o.v }

func (o Observation) Deviates(old Observation, threshold float64) bool {
	if old.v.Cmp(i(0)) == 0 {
		if o.v.Cmp(i(0)) == 0 {
			return false 		}
		return true 	}
		f64, _ := (&big.Rat{}).SetFrac(i(0).Sub(o.v, old.v), old.v).Float64()
	return math.Abs(f64) > threshold
}

func (o Observation) Marshal() []byte {
	if o.v.Cmp(MaxObservation) > 0 || o.v.Cmp(MinObservation) < 0 {
		panic(tooLarge(o.v))
	}
	negative := o.v.Cmp(i(0)) < 0
	val := (&big.Int{}).Set(o.v)
	if negative {
		val.Add(val, big.NewInt(1))
	}
	b := val.Bytes() 	if len(b) > 24 {
		panic("b must fit in 24 bytes, given it's an int192")
	}
	b = bytes.Join([][]byte{bytes.Repeat([]byte{0}, 24-len(b)), b}, []byte{})
	if len(b) != 24 {
		panic("wrong length; there must be an error in the padding of b")
	}
	if negative {
		twosComplement(b)
		b[0] = b[0] | topBit 	}
	return b
}

func UnmarshalObservation(s []byte) (Observation, error) {
	if len(s) != 24 {
		return Observation{}, errors.Errorf("wrong length for serialized "+
			"Observation: length %d 0x%x", len(s), s)
	}
	negative := s[0]&topBit != 0
	if negative {
		t := make([]byte, len(s))
		copy(t, s)
		twosComplement(t)
		s = t
	}
	if s[0]&topBit != 0 {
		panic("two's complement did not cancel top bit")
	}
	rv := (&big.Int{}).SetBytes(s)
	if negative {
		rv.Neg(rv).Sub(rv, big.NewInt(1))
	}
	return MakeObservation(rv)
}

func (o Observation) String() string {
	return fmt.Sprintf("Observation{%d}", o.v)
}

func (o Observation) Equal(o2 Observation) bool {
	return o.v.Cmp(o2.v) == 0
}

var topBit, allBits uint8 = 1 << 7, (1 << 8) - 1

func twosComplement(b []byte) {
	for bi, c := range b {
		b[bi] = allBits ^ c 	}
}

func uInt64sToObservation(w1, w2, w3 uint64) Observation {
	var b [24]byte
	for i, w := range []uint64{w1, w2, w3} {
		start := i * 8
		binary.BigEndian.PutUint64(b[start:start+8], w)
	}
	obs, err := UnmarshalObservation(b[:])
	if err != nil {
		panic(errors.Wrapf(err, "while constructing observation.Observation from 0x%x "+
			"0x%x 0x%x", w1, w2, w3))
	}
	return obs
}

func observationToUInt64s(o Observation) (w1, w2, w3 uint64) {
	b := o.Marshal()
	var uint64s [3]uint64
	for i := 0; i < 24; i += 8 {
		uint64s[i/8] = binary.BigEndian.Uint64(b[i : i+8])
	}
	return uint64s[0], uint64s[1], uint64s[2]
}

func GenObservationValue() gopter.Gen {
	return gopter.DeriveGen(uInt64sToObservation, observationToUInt64s,
		gen.UInt64(), gen.UInt64(), gen.UInt64())
}

func XXXTestingOnlyNewObservation(v *big.Int) Observation {
	return Observation{v: v}
}
