package hashlib

import (
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// BytesOfBytesKeccak will compute a keccak256 hash of the provided bytes of bytes slice
func BytesOfBytesKeccak(b [][]byte) ([32]byte, error) {
	if len(b) == 0 {
		return [32]byte{}, nil
	}

	h := utils.Keccak256Fixed(b[0])
	for _, v := range b[1:] {
		h = utils.Keccak256Fixed(append(h[:], v...))
	}
	return h, nil
}
