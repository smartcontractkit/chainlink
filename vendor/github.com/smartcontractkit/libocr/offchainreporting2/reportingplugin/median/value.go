package median

import (
	"math/big"

	"github.com/smartcontractkit/libocr/bigbigendian"
)

var i = big.NewInt

// Bounds on an int192
const byteWidth = 24
const bitWidth = byteWidth * 8

var one *big.Int = big.NewInt(1)

// 2**191-1
func MaxValue() *big.Int {
	result := MinValue()
	result.Abs(result)
	result.Sub(result, one)
	return result
}

// -2**191
func MinValue() *big.Int {
	result := &big.Int{}
	result.Lsh(one, bitWidth-1)
	result.Neg(result)
	return result
}

// Encodes a value using 24-byte big endian two's complement representation. This function never panics.
func EncodeValue(i *big.Int) ([]byte, error) {
	return bigbigendian.SerializeSigned(byteWidth, i)
}

// Decodes a value using 24-byte big endian two's complement representation. This function never panics.
func DecodeValue(s []byte) (*big.Int, error) {
	return bigbigendian.DeserializeSigned(byteWidth, s)
}
