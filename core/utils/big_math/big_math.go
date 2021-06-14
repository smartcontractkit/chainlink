package bigmath

import "math/big"

// Compensate for awkward big.Int API. Can cause an extra allocation or two.
func I() *big.Int                                    { return new(big.Int) }
func Add(addend1, addend2 *big.Int) *big.Int         { return I().Add(addend1, addend2) }
func Div(dividend, divisor *big.Int) *big.Int        { return I().Div(dividend, divisor) }
func Equal(left, right *big.Int) bool                { return left.Cmp(right) == 0 }
func Exp(base, exponent, modulus *big.Int) *big.Int  { return I().Exp(base, exponent, modulus) }
func Mul(multiplicand, multiplier *big.Int) *big.Int { return I().Mul(multiplicand, multiplier) }
func Mod(dividend, divisor *big.Int) *big.Int        { return I().Mod(dividend, divisor) }
func Sub(minuend, subtrahend *big.Int) *big.Int      { return I().Sub(minuend, subtrahend) }

var Zero = big.NewInt(0)
var One = big.NewInt(1)
var Two = big.NewInt(2)
var Three = big.NewInt(3)
var Four = big.NewInt(4)
var Seven = big.NewInt(7)
