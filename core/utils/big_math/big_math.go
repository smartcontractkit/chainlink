// Package bigmath compensates for awkward big.Int API. Can cause an extra allocation or two.
package bigmath

import (
	"fmt"
	"math/big"
	"strings"
)

// ToIntable represents a type that is convertable to a big.Int, ex utils.Big
type ToIntable interface {
	ToInt() *big.Int
}

// I returns a new big.Int.
func I() *big.Int { return new(big.Int) }

// Add performs addition with the given values after coercing them to big.Int, or panics if it cannot.
func Add(addend1, addend2 interface{}) *big.Int { return I().Add(bnIfy(addend1), bnIfy(addend2)) }

// Div performs division with the given values after coercing them to big.Int, or panics if it cannot.
func Div(dividend, divisor interface{}) *big.Int { return I().Div(bnIfy(dividend), bnIfy(divisor)) }

// Equal compares the given values after coercing them to big.Int, or panics if it cannot.
func Equal(left, right interface{}) bool { return bnIfy(left).Cmp(bnIfy(right)) == 0 }

// Exp performs modular eponentiation with the given values after coercing them to big.Int, or panics if it cannot.
func Exp(base, exponent, modulus interface{}) *big.Int {
	return I().Exp(bnIfy(base), bnIfy(exponent), bnIfy(modulus))
}

// Mul performs multiplication with the given values after coercing them to big.Int, or panics if it cannot.
func Mul(multiplicand, multiplier interface{}) *big.Int {
	return I().Mul(bnIfy(multiplicand), bnIfy(multiplier))
}

// Mod performs modulus with the given values after coercing them to big.Int, or panics if it cannot.
func Mod(dividend, divisor interface{}) *big.Int { return I().Mod(bnIfy(dividend), bnIfy(divisor)) }

// Sub performs subtraction with the given values after coercing them to big.Int, or panics if it cannot.
func Sub(minuend, subtrahend interface{}) *big.Int { return I().Sub(bnIfy(minuend), bnIfy(subtrahend)) }

// Max returns the maximum of the two given values after coercing them to big.Int,
// or panics if it cannot.
func Max(x, y interface{}) *big.Int {
	xBig := bnIfy(x)
	yBig := bnIfy(y)
	if xBig.Cmp(yBig) == 1 {
		return xBig
	}
	return yBig
}

// Accumulate returns the sum of the given slice after coercing all elements
// to a big.Int, or panics if it cannot.
func Accumulate(s []interface{}) (r *big.Int) {
	r = big.NewInt(0)
	for _, e := range s {
		r.Add(r, bnIfy(e))
	}
	return
}

func bnIfy(val interface{}) *big.Int {
	if toIntable, ok := val.(ToIntable); ok {
		return toIntable.ToInt()
	}
	switch v := val.(type) {
	case uint:
		return big.NewInt(0).SetUint64(uint64(v))
	case uint8:
		return big.NewInt(0).SetUint64(uint64(v))
	case uint16:
		return big.NewInt(0).SetUint64(uint64(v))
	case uint32:
		return big.NewInt(0).SetUint64(uint64(v))
	case uint64:
		return big.NewInt(0).SetUint64(v)
	case int:
		return big.NewInt(int64(v))
	case int8:
		return big.NewInt(int64(v))
	case int16:
		return big.NewInt(int64(v))
	case int32:
		return big.NewInt(int64(v))
	case int64:
		return big.NewInt(int64(v))
	case float64: // when decoding from db: JSON numbers are floats
		return big.NewInt(0).SetUint64(uint64(v))
	case string:
		if strings.TrimSpace(v) == "" {
			panic("invalid big int string")
		}
		n, ok := big.NewInt(0).SetString(v, 10)
		if !ok {
			panic(fmt.Sprintf("unable to convert %s to big.Int", v))
		}
		return n
	case *big.Int:
		return v
	default:
		panic(fmt.Sprintf("invalid type for big num conversion: %v", v))
	}
}

//nolint
var (
	Zero  = big.NewInt(0)
	One   = big.NewInt(1)
	Two   = big.NewInt(2)
	Three = big.NewInt(3)
	Four  = big.NewInt(4)
	Seven = big.NewInt(7)
)
