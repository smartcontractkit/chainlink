package bigmath

import (
	"fmt"
	"math/big"
	"strings"
)

// Compensate for awkward big.Int API. Can cause an extra allocation or two.
func I() *big.Int                                { return new(big.Int) }
func Add(addend1, addend2 interface{}) *big.Int  { return I().Add(bnIfy(addend1), bnIfy(addend2)) }
func Div(dividend, divisor interface{}) *big.Int { return I().Div(bnIfy(dividend), bnIfy(divisor)) }
func Equal(left, right interface{}) bool         { return bnIfy(left).Cmp(bnIfy(right)) == 0 }
func Exp(base, exponent, modulus interface{}) *big.Int {
	return I().Exp(bnIfy(base), bnIfy(exponent), bnIfy(modulus))
}
func Mul(multiplicand, multiplier interface{}) *big.Int {
	return I().Mul(bnIfy(multiplicand), bnIfy(multiplier))
}
func Mod(dividend, divisor interface{}) *big.Int   { return I().Mod(bnIfy(dividend), bnIfy(divisor)) }
func Sub(minuend, subtrahend interface{}) *big.Int { return I().Sub(bnIfy(minuend), bnIfy(subtrahend)) }

func bnIfy(val interface{}) *big.Int {
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

var Zero = big.NewInt(0)
var One = big.NewInt(1)
var Two = big.NewInt(2)
var Three = big.NewInt(3)
var Four = big.NewInt(4)
var Seven = big.NewInt(7)
