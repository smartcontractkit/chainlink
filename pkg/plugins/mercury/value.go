package mercury

import (
	"math/big"

	"github.com/smartcontractkit/libocr/bigbigendian"
)

// Bounds on an int192
const byteWidthInt192 = 24

// Encodes a value using 24-byte big endian two's complement representation. This function never panics.
func EncodeValueInt192(i *big.Int) ([]byte, error) {
	return bigbigendian.SerializeSigned(byteWidthInt192, i)
}

// Decodes a value using 24-byte big endian two's complement representation. This function never panics.
func DecodeValueInt192(s []byte) (*big.Int, error) {
	return bigbigendian.DeserializeSigned(byteWidthInt192, s)
}
