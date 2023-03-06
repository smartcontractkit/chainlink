// Serializes and deserializes big ints using two's complement big endian representation
package bigbigendian

import (
	"fmt"
	"math/big"
)

// The maximum size (in bytes) of serialized representation we support
const MaxSize = 128

// Serializes a signed big.Int into a byte slice with size bytes. Does not mutate its inputs. Does not panic.
func SerializeSigned(size int, i *big.Int) ([]byte, error) {
	if i == nil {
		return nil, fmt.Errorf("i is nil")
	}

	if !(0 < size && size <= MaxSize) {
		return nil, fmt.Errorf("size is %v, but must be between 1 and %v", size, MaxSize)
	}

	bitSize := size * 8
	negative := i.Sign() < 0

	b := make([]byte, size)
	if negative {
		// To find the two's complement we subtract one from the absolute value, then invert

		// If input is valid, max abs(i) here: 2**(bitSize - 1)

		tmp := big.NewInt(1)
		tmp.Add(i, tmp) // i is negative, so to subtract from its absolute value we need to add

		// If input is valid, max abs(tmp) here: 2**(bitSize - 1)-1 = 2**(bitSize - 2) + ... + 2**0

		if bitSize <= tmp.BitLen() {
			return nil, fmt.Errorf("i doesn't fit into a %v-byte two's complement", size)
		}
		tmp.FillBytes(b) // encodes abs(tmp) into b
		for i := range b {
			b[i] ^= 0xff
		}
	} else {
		if bitSize <= i.BitLen() {
			return nil, fmt.Errorf("i doesn't fit into a %v-byte two's complement", size)
		}
		i.FillBytes(b)
	}
	return b, nil
}

// Deserializes a byte slice with size bytes into a signed big.Int. Does not mutate its inputs. Does not panic.
func DeserializeSigned(size int, b []byte) (*big.Int, error) {
	if !(0 < size && size <= MaxSize) {
		return nil, fmt.Errorf("size is %v, but must be between 1 and %v", size, MaxSize)
	}
	if len(b) != size {
		return &big.Int{}, fmt.Errorf("expected b to have length %v, but got length %v", size, len(b))
	}
	bitSize := size * 8
	val := (&big.Int{}).SetBytes(b)
	negative := b[0]&0x80 != 0
	if negative {
		// In two's complement representation, the msb has a negative sign, e.g.
		// "1011" represents -(2**3) + 2**1 + 2**0.
		// However, SetBytes considered the msb to have a positive sign.
		// We thus compute val - 2**(bitSize-1) - 2**(bitSize-1) = val - 2**bitSize.
		powerOfTwo := (&big.Int{}).SetInt64(1)
		powerOfTwo.Lsh(powerOfTwo, uint(bitSize))
		val.Sub(val, powerOfTwo)
	}
	return val, nil
}
