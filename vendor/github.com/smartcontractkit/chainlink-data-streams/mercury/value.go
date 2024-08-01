package mercury

import (
	"math/big"

	"github.com/smartcontractkit/libocr/bigbigendian"
)

var MaxInt192 *big.Int
var MaxInt192Enc []byte

func init() {
	one := big.NewInt(1)
	// Compute the maximum value of int192
	// 1<<191 - 1
	MaxInt192 = new(big.Int).Lsh(one, 191)
	MaxInt192.Sub(MaxInt192, one)

	var err error
	MaxInt192Enc, err = EncodeValueInt192(MaxInt192)
	if err != nil {
		panic(err)
	}
}

// Bounds on an int192
const ByteWidthInt192 = 24

// Encodes a value using 24-byte big endian two's complement representation. This function never panics.
func EncodeValueInt192(i *big.Int) ([]byte, error) {
	return bigbigendian.SerializeSigned(ByteWidthInt192, i)
}

// Decodes a value using 24-byte big endian two's complement representation. This function never panics.
func DecodeValueInt192(s []byte) (*big.Int, error) {
	return bigbigendian.DeserializeSigned(ByteWidthInt192, s)
}

func MustEncodeValueInt192(i *big.Int) []byte {
	val, err := EncodeValueInt192(i)
	if err != nil {
		panic(err)
	}
	return val
}
