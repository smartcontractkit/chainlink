package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
)

const (
	// FormatBytes encodes the output as bytes
	FormatBytes = "bytes"
	// FormatPreformatted encodes the output, assumed to be hex, as bytes.
	FormatPreformatted = "preformatted"
	// FormatUint256 encodes the output as bytes containing a uint256
	FormatUint256 = "uint256"
	// FormatInt256 encodes the output as bytes containing an int256
	FormatInt256 = "int256"
	// FormatBool encodes the output as bytes containing a bool
	FormatBool = "bool"
)

// ABIEncode is the equivalent of abi.encode.
// See a full set of examples https://github.com/ethereum/go-ethereum/blob/420b78659bef661a83c5c442121b13f13288c09f/accounts/abi/packing_test.go#L31
func ABIEncode(abiStr string, values ...interface{}) ([]byte, error) {
	// Create a dummy method with arguments
	inDef := fmt.Sprintf(`[{ "name" : "method", "type": "function", "inputs": %s}]`, abiStr)
	inAbi, err := abi.JSON(strings.NewReader(inDef))
	if err != nil {
		return nil, err
	}
	res, err := inAbi.Pack("method", values...)
	if err != nil {
		return nil, err
	}
	return res[4:], nil
}

// ABIEncode is the equivalent of abi.decode.
// See a full set of examples https://github.com/ethereum/go-ethereum/blob/420b78659bef661a83c5c442121b13f13288c09f/accounts/abi/packing_test.go#L31
func ABIDecode(abiStr string, data []byte) ([]interface{}, error) {
	inDef := fmt.Sprintf(`[{ "name" : "method", "type": "function", "outputs": %s}]`, abiStr)
	inAbi, err := abi.JSON(strings.NewReader(inDef))
	if err != nil {
		return nil, err
	}
	return inAbi.Unpack("method", data)
}

// ConcatBytes appends a bunch of byte arrays into a single byte array
func ConcatBytes(bufs ...[]byte) []byte {
	return bytes.Join(bufs, []byte{})
}

func roundToEVMWordBorder(length int) int {
	mod := length % EVMWordByteLen
	if mod == 0 {
		return 0
	}
	return EVMWordByteLen - mod
}

// EVMEncodeBytes encodes arbitrary bytes as bytes expected by the EVM
func EVMEncodeBytes(input []byte) []byte {
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

func parseDecimalString(input string) (*big.Int, error) {
	d, err := decimal.NewFromString(input)
	return d.BigInt(), err
}

func parseNumericString(input string) (*big.Int, error) {
	if HasHexPrefix(input) {
		output, ok := big.NewInt(0).SetString(RemoveHexPrefix(input), 16)
		if !ok {
			return nil, fmt.Errorf("error parsing hex %s", input)
		}
		return output, nil
	}

	output, ok := big.NewInt(0).SetString(input, 10)
	if !ok {
		return parseDecimalString(input)
	}
	return output, nil
}

func parseJSONAsEVMWord(value gjson.Result) (*big.Int, error) {
	output := new(big.Int)

	switch value.Type {
	case gjson.String:
		var err error
		output, err = parseNumericString(value.Str)
		if err != nil {
			return nil, err
		}

	case gjson.Number:
		output.SetInt64(int64(value.Num))

	case gjson.Null:

	default:
		return nil, fmt.Errorf("unsupported encoding for value: %s", value.Type)
	}

	return output, nil
}

// EVMTranscodeUint256 converts a json input to an EVM uint256
func EVMTranscodeUint256(value gjson.Result) ([]byte, error) {
	output, err := parseJSONAsEVMWord(value)
	if err != nil {
		return nil, err
	}

	if output.Cmp(big.NewInt(0)) < 0 {
		return nil, fmt.Errorf("%v cannot be represented as uint256", output)
	}

	return EVMWordBigInt(output)
}

// EVMTranscodeInt256 converts a json input to an EVM int256
func EVMTranscodeInt256(value gjson.Result) ([]byte, error) {
	output, err := parseJSONAsEVMWord(value)
	if err != nil {
		return nil, err
	}

	return EVMWordSignedBigInt(output)
}

// EVMWordUint64 returns a uint64 as an EVM word byte array.
func EVMWordUint64(val uint64) []byte {
	word := make([]byte, EVMWordByteLen)
	binary.BigEndian.PutUint64(word[EVMWordByteLen-8:], val)
	return word
}

// EVMWordUint32 returns a uint32 as an EVM word byte array.
func EVMWordUint32(val uint32) []byte {
	word := make([]byte, EVMWordByteLen)
	binary.BigEndian.PutUint32(word[EVMWordByteLen-4:], val)
	return word
}

// EVMWordUint128 returns a uint128 as an EVM word byte array.
func EVMWordUint128(val *big.Int) ([]byte, error) {
	bytes := val.Bytes()
	if val.BitLen() > 128 {
		return nil, fmt.Errorf("overflow saving uint128 to EVM word: %v", val)
	} else if val.Sign() == -1 {
		return nil, fmt.Errorf("invalid attempt to save negative value as uint128 to EVM word: %v", val)
	}
	return common.LeftPadBytes(bytes, EVMWordByteLen), nil
}

// EVMWordSignedBigInt returns a big.Int as an EVM word byte array, with
// support for a signed representation. Returns error on overflow.
func EVMWordSignedBigInt(val *big.Int) ([]byte, error) {
	bytes := val.Bytes()
	if val.BitLen() > (8*EVMWordByteLen - 1) {
		return nil, fmt.Errorf("overflow saving signed big.Int to EVM word: %v", val)
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
		return nil, fmt.Errorf("overflow saving big.Int to EVM word: %v", val)
	}
	return common.LeftPadBytes(bytes, EVMWordByteLen), nil
}

// Bytes32FromString returns a 32 byte array filled from the given string, which may be of any length.
func Bytes32FromString(s string) [32]byte {
	var b32 [32]byte
	copy(b32[:], s)
	return b32
}

// Bytes4FromString returns a 4 byte array filled from the given string, which may be of any length.
func Bytes4FromString(s string) [4]byte {
	var b4 [4]byte
	copy(b4[:], s)
	return b4
}

func MustAbiType(ts string, components []abi.ArgumentMarshaling) abi.Type {
	ty, err := abi.NewType(ts, "", components)
	if err != nil {
		panic(err)
	}
	return ty
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
