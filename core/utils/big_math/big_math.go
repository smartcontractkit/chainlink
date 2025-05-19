// Package bigmath compensates for awkward big.Int API. Can cause an extra allocation or two.
package bigmath

import (
	"math/big"
)

// I returns a new big.Int.
func I() *big.Int { return new(big.Int) }

// Add performs addition with the given values.
func Add(addend1, addend2 *big.Int) *big.Int { return I().Add(addend1, addend2) }

// Div performs division with the given values.
func Div(dividend, divisor *big.Int) *big.Int { return I().Div(dividend, divisor) }

// Equal compares the given values.
func Equal(left, right *big.Int) bool { return left.Cmp(right) == 0 }

// Exp performs modular eponentiation with the given values.
func Exp(base, exponent, modulus *big.Int) *big.Int {
	return I().Exp(base, exponent, modulus)
}

// Mul performs multiplication with the given values.
func Mul(multiplicand, multiplier *big.Int) *big.Int {
	return I().Mul(multiplicand, multiplier)
}

// Mod performs modulus with the given values.
func Mod(dividend, divisor *big.Int) *big.Int { return I().Mod(dividend, divisor) }

// Sub performs subtraction with the given values.
func Sub(minuend, subtrahend *big.Int) *big.Int { return I().Sub(minuend, subtrahend) }

// Max returns the maximum of the two given values.
func Max(x, y *big.Int) *big.Int {
	if x.Cmp(y) == 1 {
		return x
	}
	return y
}

// Min returns the min of the two given values.
func Min(x, y *big.Int) *big.Int {
	if x.Cmp(y) == -1 {
		return x
	}
	return y
}

// Accumulate returns the sum of the given slice.
func Accumulate(s []*big.Int) (r *big.Int) {
	r = big.NewInt(0)
	for _, e := range s {
		r.Add(r, e)
	}
	return
}

var (
	Zero  = big.NewInt(0)
	One   = big.NewInt(1)
	Two   = big.NewInt(2)
	Three = big.NewInt(3)
	Four  = big.NewInt(4)
	Seven = big.NewInt(7)
)
