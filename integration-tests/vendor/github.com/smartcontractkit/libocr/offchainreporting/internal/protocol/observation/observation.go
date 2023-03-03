// Package observation contains the data structures and logic for handling
// observations provided by the client DataSource. Its role is to encapsulate
// the Observation type so that it can be changed relatively easily.
package observation

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/offchainreporting/types"
)

type Observation struct{ v *big.Int }

type Observations []Observation

var i = big.NewInt

// Bounds on an ethereum int192
const byteWidth = 24
const bitWidth = byteWidth * 8

var MaxObservation = i(0).Sub(i(0).Lsh(i(1), bitWidth-1), i(1)) // 2**191 - 1
var MinObservation = i(0).Sub(i(0).Neg(MaxObservation), i(1))   // -2**191

func tooLarge(o *big.Int) error {
	return errors.Errorf("value won't fit in int%v: 0x%x", bitWidth, o)
}

// MakeObservation returns v as an ethereum int192, if it fits, errors otherwise.
func MakeObservation(w types.Observation) (Observation, error) {
	v := (*big.Int)(w)
	// nil can sometimes occur here because it's the zero value for a pointer in a
	// struct, and w comes from a zero struct with a *big.Int field. We always
	// want the corresponding value for the "zero observation" to be zero.
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

func (o Observation) Deviates(old Observation, thresholdPPB uint64) bool {
	if old.v.Cmp(i(0)) == 0 {
		//nolint:gosimple
		if o.v.Cmp(i(0)) == 0 {
			return false // Both values are zero; no deviation
		}
		return true // Any deviation from 0 is significant
	}
	// ||o.v - old.v|| / ||old.v||, approximated by a float
	change := &big.Rat{}
	change.SetFrac(i(0).Sub(o.v, old.v), old.v)
	change.Abs(change)
	threshold := &big.Rat{}
	threshold.SetFrac(
		(&big.Int{}).SetUint64(thresholdPPB),
		(&big.Int{}).SetUint64(1000000000),
	)
	return change.Cmp(threshold) >= 0
}

// Bytes returns the twos-complement representation of o
//
// This panics on OOB values, because MakeObservation and UnmarshalObservation
// are the only external ways to create an Observation, and that already checks
// the bounds
func (o Observation) Marshal() []byte {
	if o.v.Cmp(MaxObservation) > 0 || o.v.Cmp(MinObservation) < 0 {
		panic(tooLarge(o.v))
	}
	negative := o.v.Sign() < 0
	val := (&big.Int{})
	if negative {
		// compute two's complement as 2**192 - abs(o.v) = 2**192 + o.v
		val.SetInt64(1)
		val.Lsh(val, bitWidth)
		val.Add(val, o.v)
	} else {
		val.Set(o.v)
	}
	b := val.Bytes() // big-endian representation of abs(val)
	if len(b) > byteWidth {
		panic(fmt.Sprintf("b must fit in %v bytes", byteWidth))
	}
	b = bytes.Join([][]byte{bytes.Repeat([]byte{0}, byteWidth-len(b)), b}, []byte{})
	if len(b) != byteWidth {
		panic("wrong length; there must be an error in the padding of b")
	}
	return b
}

func UnmarshalObservation(s []byte) (Observation, error) {
	if len(s) != byteWidth {
		return Observation{}, errors.Errorf("wrong length for serialized "+
			"Observation: length %d 0x%x", len(s), s)
	}
	val := (&big.Int{}).SetBytes(s)
	negative := val.Cmp(MaxObservation) > 0
	if negative {
		maxUint := (&big.Int{}).SetInt64(1)
		maxUint.Lsh(maxUint, bitWidth)
		val.Sub(maxUint, val)
		val.Neg(val)
	}
	return MakeObservation(val)
}

func (o Observation) String() string {
	return fmt.Sprintf("Observation{%d}", o.v)
}

func (o Observation) Equal(o2 Observation) bool {
	return o.v.Cmp(o2.v) == 0
}

var _ encoding.TextMarshaler = Observation{}

func (o Observation) MarshalText() (text []byte, err error) {
	if o.v == nil {
		return []byte{}, nil
	}

	return o.v.MarshalText()
}

func uInt64sToObservation(w1, w2, w3 uint64) Observation {
	var b [byteWidth]byte
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
	for i := 0; i < byteWidth; i += 8 {
		uint64s[i/8] = binary.BigEndian.Uint64(b[i : i+8])
	}
	return uint64s[0], uint64s[1], uint64s[2]
}

func GenObservation() gopter.Gen {
	return gopter.DeriveGen(uInt64sToObservation, observationToUInt64s,
		gen.UInt64(), gen.UInt64(), gen.UInt64())
}

// XXXTestingOnlyNewObservation returns a new observation with no bounds
// checking on v.
func XXXTestingOnlyNewObservation(v *big.Int) Observation {
	return Observation{v: v}
}
