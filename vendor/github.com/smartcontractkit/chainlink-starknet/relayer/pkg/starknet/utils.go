package starknet

import (
	"cmp"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/NethermindEth/juno/core/felt"
)

const (
	FeltLength = 32
	chunkSize  = 31
)

// padd bytes to specific length
func PadBytes(a []byte, length int) []byte {
	if len(a) < length {
		pad := make([]byte, length-len(a))
		return append(pad, a...)
	}

	// return original if length is >= to specified length
	return a
}

func NilResultError(funcName string) error {
	return fmt.Errorf("nil result in %s", funcName)
}

func Min[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// EncodeFelts takes a byte slice and splits as bunch of felts. First felt indicates the total byte size.
func EncodeFelts(data []byte) (felts []*big.Int) {
	// prefix with len
	length := big.NewInt(int64(len(data)))
	felts = append(felts, length)

	// chunk every 31 bytes
	for i := 0; i < len(data); i += chunkSize {
		chunk := data[i:Min(i+chunkSize, len(data))]
		// cast to int
		felt := new(big.Int).SetBytes(chunk)
		felts = append(felts, felt)
	}

	return felts
}

// DecodeFelts is the reverse of EncodeFelts
func DecodeFelts(felts []*big.Int) ([]byte, error) {
	if len(felts) == 0 {
		return []byte{}, nil
	}

	data := []byte{}
	buf := make([]byte, chunkSize)
	length := felts[0].Uint64()

	for _, felt := range felts[1:] {
		bytesLen := Min(chunkSize, length)
		bytesBuffer := buf[:bytesLen]

		// TODO: this is inefficient because Bytes() and FillBytes() duplicate work
		// reuse felt.Bytes()
		if bytesLen < uint64(len(felt.Bytes())) {
			return nil, errors.New("invalid: felt array can't be decoded")
		}

		felt.FillBytes(bytesBuffer)
		data = append(data, bytesBuffer...)

		length -= bytesLen
	}

	if length != 0 {
		return nil, errors.New("invalid: contained less bytes than the specified length")
	}

	return data, nil
}

func FeltsToBig(in []*felt.Felt) (out []*big.Int) {
	for _, f := range in {
		out = append(out, f.BigInt(big.NewInt(0)))
	}

	return out
}

/* Testing utils - do not use (XXX) outside testing context */

func XXXMustHexDecodeString(data string) []byte {
	bytes, err := hex.DecodeString(data)
	if err != nil {
		panic(err)
	}
	return bytes
}
