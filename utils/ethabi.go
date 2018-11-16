package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tidwall/gjson"
)

const (
	// FormatBytes encodes the output as bytes
	FormatBytes = "bytes"
	// FormatUint256 encodes the output as bytes containing a uint256
	FormatUint256 = "uint256"
	// FormatInt256 encodes the output as bytes containing an int256
	FormatInt256 = "int256"
	// FormatBool encodes the output as bytes containing a bool
	FormatBool = "bool"
)

// ConcatBytes appends a bunch of byte arrays into a single byte array
func ConcatBytes(bufs ...[]byte) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	for _, b := range bufs {
		_, err := buffer.Write(b)
		if err != nil {
			return nil, err
		}
	}
	return buffer.Bytes(), nil
}

// EVMTranscodeBytes converts a json input to an EVM bytes array
func EVMTranscodeBytes(value gjson.Result) ([]byte, error) {
	switch value.Type {
	case gjson.String:
		return EVMEncodeBytes([]byte(value.Str))

	case gjson.False:
		return EVMEncodeBytes(EVMWordUint64(0))

	case gjson.True:
		return EVMEncodeBytes(EVMWordUint64(1))

	case gjson.Number:
		word, err := EVMWordSignedBigInt(big.NewInt(int64(value.Num)))
		if err != nil {
			return []byte{}, nil
		}

		return EVMEncodeBytes(word)

	default:
		return []byte{}, fmt.Errorf("unsupported encoding for value: %s", value.Type)
	}
}

func roundToEVMWordBorder(length int) int {
	mod := length % EVMWordByteLen
	if mod == 0 {
		return 0
	}
	return EVMWordByteLen - mod
}

// EVMEncodeBytes encodes arbitrary bytes as bytes expected by the EVM
func EVMEncodeBytes(input []byte) ([]byte, error) {
	length := len(input)
	return ConcatBytes(
		EVMWordUint64(uint64(length)),
		input,
		make([]byte, roundToEVMWordBorder(length)))
}

// EVMTranscodeBool converts a json input to an EVM bool
func EVMTranscodeBool(value gjson.Result) ([]byte, error) {
	var output uint64

	switch value.Type {
	case gjson.Number:
		if value.Num != 0 {
			output = 1
		}

	case gjson.String:
		if len(value.Str) > 0 {
			output = 1
		}

	case gjson.True:
		output = 1

	case gjson.JSON:
		value.ForEach(func(key, value gjson.Result) bool {
			output = 1
			return false
		})

	case gjson.False, gjson.Null:

	default:
		panic(fmt.Errorf("unreachable/unsupported encoding for value: %s", value.Type))
	}

	return EVMWordUint64(output), nil
}

// EVMTranscodeUint256 converts a json input to an EVM uint256
func EVMTranscodeUint256(value gjson.Result) ([]byte, error) {
	output := new(big.Int)

	switch value.Type {
	case gjson.String:
		var ok bool
		if HasHexPrefix(value.Str) {
			output, ok = output.SetString(RemoveHexPrefix(value.Str), 16)
		} else {
			output, ok = output.SetString(value.Str, 10)
		}
		if !ok {
			return []byte{}, fmt.Errorf("error parsing %s", value.Str)
		}

	case gjson.Number:
		if value.Num < 0 {
			return []byte{}, errors.New("cannot format negative number as uint256")
		}
		output.SetUint64(uint64(value.Num))

	case gjson.Null:

	default:
		return []byte{}, fmt.Errorf("unsupported encoding for value: %s", value.Type)
	}

	return EVMWordBigInt(output)
}

// EVMTranscodeInt256 converts a json input to an EVM int256
func EVMTranscodeInt256(value gjson.Result) ([]byte, error) {
	output := new(big.Int)

	switch value.Type {
	case gjson.String:
		var ok bool
		if HasHexPrefix(value.Str) {
			output, ok = output.SetString(RemoveHexPrefix(value.Str), 16)
		} else {
			output, ok = output.SetString(value.Str, 10)
		}
		if !ok {
			return []byte{}, fmt.Errorf("error parsing %s", value.Str)
		}

	case gjson.Number:
		output.SetInt64(int64(value.Num))

	case gjson.Null:

	default:
		return []byte{}, fmt.Errorf("unsupported encoding for value: %s", value.Type)
	}

	return EVMWordSignedBigInt(output)
}

// EVMTranscodeJSONWithFormat given a JSON input and a format specifier, encode the
// value for use by the EVM
func EVMTranscodeJSONWithFormat(value gjson.Result, format string) ([]byte, error) {
	switch format {
	case FormatBytes:
		return EVMTranscodeBytes(value)

	case FormatUint256:
		data, err := EVMTranscodeUint256(value)
		if err != nil {
			return []byte{}, err
		}
		return EVMEncodeBytes(data)

	case FormatInt256:
		data, err := EVMTranscodeInt256(value)
		if err != nil {
			return []byte{}, err
		}
		return EVMEncodeBytes(data)

	case FormatBool:
		data, err := EVMTranscodeBool(value)
		if err != nil {
			return []byte{}, err
		}
		return EVMEncodeBytes(data)

	default:
		return []byte{}, fmt.Errorf("unsupported format: %s", format)
	}
}

// EVMWordUint64 returns a uint64 as an EVM word byte array.
func EVMWordUint64(val uint64) []byte {
	word := make([]byte, EVMWordByteLen)
	binary.BigEndian.PutUint64(word[EVMWordByteLen-8:], val)
	return word
}

// EVMWordSignedBigInt returns a big.Int as an EVM word byte array, with
// support for a signed representation. Returns error on overflow.
func EVMWordSignedBigInt(val *big.Int) ([]byte, error) {
	bytes := val.Bytes()
	if val.BitLen() > (8*EVMWordByteLen - 1) {
		return nil, fmt.Errorf("Overflow saving signed big.Int to EVM word: %v", val)
	}
	if val.Sign() == -1 {
		twosComplement := new(big.Int).Add(val, MaxUint256)
		bytes = new(big.Int).Add(twosComplement, big.NewInt(1)).Bytes()
	}
	return common.LeftPadBytes(bytes, EVMWordByteLen), nil
}

// EVMWordBigInt returns a big.Int as an EVM word byte array, with support for
// a signed representation. Returns error on overflow.
func EVMWordBigInt(val *big.Int) ([]byte, error) {
	if val.Sign() == -1 {
		return nil, errors.New("Uint256 cannot be negative")
	}
	bytes := val.Bytes()
	if len(bytes) > EVMWordByteLen {
		return nil, fmt.Errorf("Overflow saving big.Int to EVM word: %v", val)
	}
	return common.LeftPadBytes(bytes, EVMWordByteLen), nil
}

// "Constants" used by EVM words
var (
	maxUint257 = &big.Int{}
	// MaxUint256 represents the largest number represented by an EVM word
	MaxUint256 = &big.Int{}
	// MaxInt256 represents the largest number represented by an EVM word using
	// signed encoding.
	MaxInt256 = &big.Int{}
	// MinInt256 represents the smallest number represented by an EVM word using
	// signed encoding.
	MinInt256 = &big.Int{}
)

func init() {
	maxUint257 = new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
	MaxUint256 = new(big.Int).Sub(maxUint257, big.NewInt(1))
	MaxInt256 = new(big.Int).Div(MaxUint256, big.NewInt(2))
	MinInt256 = new(big.Int).Neg(MaxInt256)
}
