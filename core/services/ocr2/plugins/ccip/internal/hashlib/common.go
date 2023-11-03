package hashlib

import (
	"strconv"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// BytesOfBytesKeccak will compute a keccak256 hash of the provided bytes of bytes slice
func BytesOfBytesKeccak(b [][]byte) ([32]byte, error) {
	if len(b) == 0 {
		return [32]byte{}, nil
	}

	joinedBytes := make([]byte, 0)
	joinedBytes = append(joinedBytes, intToBytes(int64(len(b)))...)
	for i := range b {
		joinedBytes = append(joinedBytes, intToBytes(int64(len(b[i])))...)
		joinedBytes = append(joinedBytes, b[i]...)
	}

	return utils.Keccak256Fixed(joinedBytes), nil
}

func intToBytes(v int64) []byte {
	return []byte(strconv.FormatInt(v, 10))
}
