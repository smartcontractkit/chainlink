package hashlib

import (
	"bytes"
	"encoding/binary"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
)

// BytesOfBytesKeccak will compute a keccak256 hash of the provided bytes of bytes slice
func BytesOfBytesKeccak(b [][]byte) ([32]byte, error) {
	if len(b) == 0 {
		return [32]byte{}, nil
	}

	encodedArr, err := encodeBytesOfBytes(b)
	if err != nil {
		return [32]byte{}, err
	}

	return utils.Keccak256Fixed(encodedArr), nil
}

// encodeBytesOfBytes encodes the nested byte arrays into a single byte array as follows
//  1. total number of nested arrays is encoded into fix-size 8 bytes at the front of the result
//  2. for each nested array
//     encode the array length into fixed-size 8 bytes, append to result
//     append the array contents to result
func encodeBytesOfBytes(b [][]byte) ([]byte, error) {
	var buffer bytes.Buffer
	if err := binary.Write(&buffer, binary.BigEndian, uint64(len(b))); err != nil {
		return nil, err
	}
	for _, arr := range b {
		if err := binary.Write(&buffer, binary.BigEndian, uint64(len(arr))); err != nil {
			return nil, err
		}
		if _, err := buffer.Write(arr); err != nil {
			return nil, err
		}
	}
	return buffer.Bytes(), nil
}
