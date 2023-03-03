package starknet

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"

	junotypes "github.com/NethermindEth/juno/pkg/types"
	caigotypes "github.com/dontpanicdao/caigo/types"

	"github.com/pkg/errors"
	"golang.org/x/exp/constraints"
)

const chunkSize = 31

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

func Min[T constraints.Ordered](a, b T) T {
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

		length -= uint64(bytesLen)
	}

	if length != 0 {
		return nil, errors.New("invalid: contained less bytes than the specified length")
	}

	return data, nil
}

// SignedBigToFelt wraps negative values correctly into felt
func SignedBigToFelt(num *big.Int) *big.Int {
	return new(big.Int).Mod(num, caigotypes.MaxFelt.Big())
}

// FeltToSigned unwraps felt into negative values
func FeltToSignedBig(felt *caigotypes.Felt) (num *big.Int) {
	num = felt.Big()
	prime := caigotypes.MaxFelt.Big()
	half := new(big.Int).Div(prime, big.NewInt(2))
	// if num > PRIME/2, then -PRIME to convert to negative value
	if num.Cmp(half) > 0 {
		return new(big.Int).Sub(num, prime)
	}
	return num
}

func HexToSignedBig(str string) (num *big.Int) {
	felt := junotypes.HexToFelt(str)
	return FeltToSignedBig(&caigotypes.Felt{Int: felt.Big()})
}

func FeltsToBig(in []*caigotypes.Felt) (out []*big.Int) {
	for _, f := range in {
		out = append(out, f.Int)
	}

	return out
}

// StringsToFelt maps felts from 'string' (hex) representation to 'caigo.Felt' representation
func StringsToFelt(in []string) (out []*caigotypes.Felt, _ error) {
	if in == nil {
		return nil, errors.New("invalid: input value")
	}

	for _, f := range in {
		felt := caigotypes.StrToFelt(f)
		if felt == nil {
			return nil, errors.New("invalid: string value")
		}

		out = append(out, felt)
	}

	return out, nil
}

func StringsToJunoFelts(in []string) []junotypes.Felt {
	out := make([]junotypes.Felt, len(in))
	for i := 0; i < len(in); i++ {
		out[i] = junotypes.HexToFelt(in[i])
	}
	return out
}

// CompareAddress compares different hex starknet addresses with potentially different 0 padding
func CompareAddress(a, b string) bool {
	aBytes, err := caigotypes.HexToBytes(a)
	if err != nil {
		return false
	}

	bBytes, err := caigotypes.HexToBytes(b)
	if err != nil {
		return false
	}

	return bytes.Equal(PadBytes(aBytes, 32), PadBytes(bBytes, 32))
}

/* Testing utils - do not use (XXX) outside testing context */

func XXXMustHexDecodeString(data string) []byte {
	bytes, err := hex.DecodeString(data)
	if err != nil {
		panic(err)
	}
	return bytes
}
