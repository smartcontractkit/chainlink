package utils

import (
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func pow2(arg int64) *big.Int {
	return new(big.Int).Exp(big.NewInt(2), big.NewInt(arg), nil)
}

func TestEVMWordUint64(t *testing.T) {
	assert.Equal(t,
		hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000001"),
		EVMWordUint64(1))
	assert.Equal(t,
		hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000100"),
		EVMWordUint64(256))
	assert.Equal(t,
		hexutil.MustDecode("0x000000000000000000000000000000000000000000000000ffffffffffffffff"),
		EVMWordUint64(math.MaxUint64))
}

func TestEVMWordUint128(t *testing.T) {
	tests := []struct {
		name string
		val  *big.Int
		exp  string
	}{
		{
			name: "1",
			val:  big.NewInt(1),
			exp:  "0x0000000000000000000000000000000000000000000000000000000000000001",
		},
		{
			name: "256",
			val:  big.NewInt(256),
			exp:  "0x0000000000000000000000000000000000000000000000000000000000000100",
		},
		{
			name: "Max Uint 128",
			val:  new(big.Int).Sub(pow2(128), big.NewInt(1)),
			exp:  "0x00000000000000000000000000000000ffffffffffffffffffffffffffffffff",
		},
	}
	for _, test := range tests {
		t.Log(test.name)
		ret, err := EVMWordUint128(test.val)
		assert.Equal(t, hexutil.MustDecode(test.exp), ret)
		require.NoError(t, err)
	}
}

func TestEVMWordUint128_Error(t *testing.T) {
	tests := []struct {
		name string
		val  *big.Int
	}{
		{
			name: "Negative number",
			val:  big.NewInt(-1),
		},
		{
			name: "Number too large: 128",
			val:  pow2(128),
		},
	}
	for _, test := range tests {
		t.Log(test.name)
		_, err := EVMWordUint128(test.val)
		assert.Error(t, err)
	}
}

func TestEVMWordSignedBigInt(t *testing.T) {
	val, err := EVMWordSignedBigInt(&big.Int{})
	assert.NoError(t, err)
	assert.Equal(t, hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000000"), val)

	val, err = EVMWordSignedBigInt(new(big.Int).SetInt64(1))
	assert.NoError(t, err)
	assert.Equal(t, hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000001"), val)

	val, err = EVMWordSignedBigInt(new(big.Int).SetInt64(256))
	assert.NoError(t, err)
	assert.Equal(t, hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000100"), val)

	val, err = EVMWordSignedBigInt(new(big.Int).SetInt64(-1))
	assert.NoError(t, err)
	assert.Equal(t, hexutil.MustDecode("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"), val)

	val, err = EVMWordSignedBigInt(MaxInt256)
	assert.NoError(t, err)
	assert.Equal(t, hexutil.MustDecode("0x7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"), val)

	val, err = EVMWordSignedBigInt(new(big.Int).Add(MaxInt256, big.NewInt(1)))
	assert.Error(t, err)
	assert.Nil(t, val)
}

func TestEVMWordBigInt(t *testing.T) {
	val, err := EVMWordBigInt(&big.Int{})
	assert.NoError(t, err)
	assert.Equal(t, hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000000"), val)

	val, err = EVMWordBigInt(new(big.Int).SetInt64(1))
	assert.NoError(t, err)
	assert.Equal(t, hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000001"), val)

	val, err = EVMWordBigInt(new(big.Int).SetInt64(256))
	assert.NoError(t, err)
	assert.Equal(t, hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000100"), val)

	val, err = EVMWordBigInt(new(big.Int).SetInt64(-1))
	assert.Error(t, err)
	assert.Nil(t, val)

	val, err = EVMWordBigInt(MaxUint256)
	assert.NoError(t, err)
	assert.Equal(t, hexutil.MustDecode("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"), val)

	val, err = EVMWordBigInt(new(big.Int).Add(MaxUint256, big.NewInt(1)))
	assert.Error(t, err)
	assert.Nil(t, val)
}

func TestEVMTranscodeBytes(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			"value is string",
			`"hello world"`,
			"0x" +
				"000000000000000000000000000000000000000000000000000000000000000b" +
				"68656c6c6f20776f726c64000000000000000000000000000000000000000000",
		},
		{
			"value is bool true",
			`true`,
			"0x" +
				"0000000000000000000000000000000000000000000000000000000000000020" +
				"0000000000000000000000000000000000000000000000000000000000000001",
		},
		{
			"value is bool false",
			`false`,
			"0x" +
				"0000000000000000000000000000000000000000000000000000000000000020" +
				"0000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			"value is positive integer",
			`19`,
			"0x" +
				"0000000000000000000000000000000000000000000000000000000000000020" +
				"0000000000000000000000000000000000000000000000000000000000000013",
		},
		{
			"value is negative integer",
			`-23`,
			"0x" +
				"0000000000000000000000000000000000000000000000000000000000000020" +
				"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe9",
		},
		// NB: The following is undesirable behavior. For more details, please see
		// https://www.pivotaltracker.com/n/workspaces/755483
		{"value is a number but not an integer",
			`19.99`,
			"0x" +
				"0000000000000000000000000000000000000000000000000000000000000020" +
				"0000000000000000000000000000000000000000000000000000000000000013",
		},
	}
	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			input := gjson.Parse(test.input)
			out, err := EVMTranscodeBytes(input)
			assert.NoError(t, err)
			assert.Equal(t, test.output, hexutil.Encode(out))
		})
	}
}

func TestEVMTranscodeBytes_ErrorsOnOverflow(t *testing.T) {
	input := gjson.Parse("1e+300")
	_, err := EVMTranscodeBytes(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "overflow saving signed big.Int to EVM word")
}

func TestEVMTranscodeBytes_UnsupportedEncoding(t *testing.T) {
	input := gjson.Parse("{}")
	_, err := EVMTranscodeBytes(input)
	assert.Error(t, err)
}

func TestEVMTranscodeBool(t *testing.T) {
	tests := []struct {
		name   string
		input  gjson.Result
		output string
	}{
		{
			"true",
			gjson.Result{Type: gjson.True},
			"0x0000000000000000000000000000000000000000000000000000000000000001",
		},
		{
			"false",
			gjson.Result{Type: gjson.False},
			"0x0000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			"null",
			gjson.Result{Type: gjson.Null},
			"0x0000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			"empty string",
			gjson.Result{Type: gjson.String, Str: ""},
			"0x0000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			"string",
			gjson.Result{Type: gjson.String, Str: "hello world"},
			"0x0000000000000000000000000000000000000000000000000000000000000001",
		},
		{
			"zero",
			gjson.Result{Type: gjson.Number, Num: 0.0},
			"0x0000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			"positive integer",
			gjson.Result{Type: gjson.Number, Num: 1239812},
			"0x0000000000000000000000000000000000000000000000000000000000000001",
		},
		{
			"empty object",
			gjson.Result{Type: gjson.JSON, Raw: "{}"},
			"0x0000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			"object with keys",
			gjson.Result{Type: gjson.JSON, Raw: `{"key": "value"}`},
			"0x0000000000000000000000000000000000000000000000000000000000000001",
		},
		{
			"empty array",
			gjson.Result{Type: gjson.JSON, Raw: "[]"},
			"0x0000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			"array with values",
			gjson.Result{Type: gjson.JSON, Raw: `["value"]`},
			"0x0000000000000000000000000000000000000000000000000000000000000001",
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {

			out, err := EVMTranscodeBool(test.input)
			assert.NoError(t, err)
			assert.Equal(t, test.output, hexutil.Encode(out))
		})
	}
}

func TestEVMTranscodeUint256(t *testing.T) {
	tests := []struct {
		name      string
		input     gjson.Result
		output    string
		wantError bool
	}{
		{
			"true",
			gjson.Result{Type: gjson.True},
			"",
			true,
		},
		{
			"false",
			gjson.Result{Type: gjson.False},
			"",
			true,
		},
		{
			"null",
			gjson.Result{Type: gjson.Null},
			"0x0000000000000000000000000000000000000000000000000000000000000000",
			false,
		},
		{
			"empty string",
			gjson.Result{Type: gjson.String, Str: ""},
			"",
			true,
		},
		{
			"string",
			gjson.Result{Type: gjson.String, Str: "hello world"},
			"",
			true,
		},
		{
			"string decimal",
			gjson.Result{Type: gjson.String, Str: "120"},
			"0x0000000000000000000000000000000000000000000000000000000000000078",
			false,
		},
		{
			"string hex",
			gjson.Result{Type: gjson.String, Str: "0xba"},
			"0x00000000000000000000000000000000000000000000000000000000000000ba",
			false,
		},
		{
			"zero",
			gjson.Result{Type: gjson.Number, Num: 0.0},
			"0x0000000000000000000000000000000000000000000000000000000000000000",
			false,
		},
		{
			"positive integer",
			gjson.Result{Type: gjson.Number, Num: 231},
			"0x00000000000000000000000000000000000000000000000000000000000000e7",
			false,
		},
		{
			"negative integer",
			gjson.Result{Type: gjson.Number, Num: -912},
			"",
			true,
		},
		{
			"unsupported encoding",
			gjson.Result{Type: gjson.JSON},
			"",
			true,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {

			out, err := EVMTranscodeUint256(test.input)
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.output, hexutil.Encode(out))
			}
		})
	}
}

func TestEVMTranscodeInt256(t *testing.T) {
	tests := []struct {
		name      string
		input     gjson.Result
		output    string
		wantError bool
	}{
		{
			"true",
			gjson.Result{Type: gjson.True},
			"",
			true,
		},
		{
			"false",
			gjson.Result{Type: gjson.False},
			"",
			true,
		},
		{
			"null",
			gjson.Result{Type: gjson.Null},
			"0x0000000000000000000000000000000000000000000000000000000000000000",
			false,
		},
		{
			"empty string",
			gjson.Result{Type: gjson.String, Str: ""},
			"",
			true,
		},
		{
			"string",
			gjson.Result{Type: gjson.String, Str: "hello world"},
			"",
			true,
		},
		{
			"string decimal",
			gjson.Result{Type: gjson.String, Str: "120"},
			"0x0000000000000000000000000000000000000000000000000000000000000078",
			false,
		},
		{
			"string hex",
			gjson.Result{Type: gjson.String, Str: "0xba"},
			"0x00000000000000000000000000000000000000000000000000000000000000ba",
			false,
		},
		{
			"zero",
			gjson.Result{Type: gjson.Number, Num: 0.0},
			"0x0000000000000000000000000000000000000000000000000000000000000000",
			false,
		},
		{
			"positive integer",
			gjson.Result{Type: gjson.Number, Num: 231},
			"0x00000000000000000000000000000000000000000000000000000000000000e7",
			false,
		},
		{
			"negative integer",
			gjson.Result{Type: gjson.Number, Num: -912},
			"0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc70",
			false,
		},
		{
			"unsupported encoding",
			gjson.Result{Type: gjson.JSON},
			"",
			true,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {

			out, err := EVMTranscodeInt256(test.input)
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.output, hexutil.Encode(out))
			}
		})
	}
}

func TestEVMTranscodeJSONWithFormat(t *testing.T) {
	tests := []struct {
		name   string
		format string
		input  string
		output string
	}{
		{
			"result is string",
			FormatBytes,
			`{"result": "hello world"}`,
			"0x" +
				"000000000000000000000000000000000000000000000000000000000000000b" +
				"68656c6c6f20776f726c64000000000000000000000000000000000000000000",
		},
		{
			"result is number",
			FormatUint256,
			`{"result": 31223}`,
			"0x" +
				"0000000000000000000000000000000000000000000000000000000000000020" +
				"00000000000000000000000000000000000000000000000000000000000079f7",
		},
		{
			"result is negative number",
			FormatInt256,
			`{"result": -123481273.1}`,
			"0x" +
				"0000000000000000000000000000000000000000000000000000000000000020" +
				"fffffffffffffffffffffffffffffffffffffffffffffffffffffffff8a3d347",
		},
		{
			"result is true",
			FormatBool,
			`{"result": true}`,
			"0x" +
				"0000000000000000000000000000000000000000000000000000000000000020" +
				"0000000000000000000000000000000000000000000000000000000000000001",
		},
		{
			"result is preformatted",
			FormatPreformatted,
			`{"result": "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"}`,
			"0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			input := gjson.GetBytes([]byte(test.input), "result")
			out, err := EVMTranscodeJSONWithFormat(input, test.format)
			require.NoError(t, err)
			assert.Equal(t, test.output, hexutil.Encode(out))
		})
	}
}

func TestEVMTranscodeJSONWithFormat_UnsupportedEncoding(t *testing.T) {
	_, err := EVMTranscodeJSONWithFormat(gjson.Result{}, "burgh")
	assert.Error(t, err)
}

func TestRoundToEVMWordBorder(t *testing.T) {
	assert.Equal(t, 0, roundToEVMWordBorder(0))
	assert.Equal(t, 0, roundToEVMWordBorder(32))
	assert.Equal(t, 31, roundToEVMWordBorder(1))
	assert.Equal(t, 1, roundToEVMWordBorder(31))
}

func TestParseNumericString(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{"0x0", "0"},
		{"0xfffffffffffffffff", "295147905179352825855"},
		{"1.0", "1"},
		{"0", "0"},
		{"1", "1"},
		{"1.0E+0", "1"},
	}

	for _, test := range tests {
		out, err := parseNumericString(test.input)
		assert.NoError(t, err)
		assert.Equal(t, test.output, out.String())
	}
}

func TestParseNumericString_InvalidHex(t *testing.T) {
	_, err := parseNumericString("0xfZ")
	assert.Error(t, err)
}

func TestParseDecimalString(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{"1.0", "1"},
		{"0", "0"},
		{"1", "1"},
		{"1.0E+0", "1"},
		{"1E+0", "1"},
		{"1e+0", "1"},
		{"0.01e+02", "1"},
		{"12072e-4", "1"},
		{"1.2072e+20", "120720000000000000000"},
		{"-1.2072e+20", "-120720000000000000000"},
		{"1.55555555555555555555e+20", "155555555555555555555"},
		{"1.000000000000000001e+18", "1000000000000000001"},
		{"1000000.000000000000000001e+18", "1000000000000000000000001"},
	}

	for _, test := range tests {
		out, err := parseDecimalString(test.input)
		assert.NoError(t, err)
		assert.Equal(t, test.output, out.String())
	}
}

// EVMTranscodeJSONWithFormat given a JSON input and a format specifier, encode the
// value for use by the EVM
func EVMTranscodeJSONWithFormat(value gjson.Result, format string) ([]byte, error) {
	switch format {
	case FormatBytes:
		return EVMTranscodeBytes(value)
	case FormatPreformatted:
		return hex.DecodeString(RemoveHexPrefix(value.Str))
	case FormatUint256:
		data, err := EVMTranscodeUint256(value)
		if err != nil {
			return []byte{}, err
		}
		return EVMEncodeBytes(data), nil

	case FormatInt256:
		data, err := EVMTranscodeInt256(value)
		if err != nil {
			return []byte{}, err
		}
		return EVMEncodeBytes(data), nil

	case FormatBool:
		data, err := EVMTranscodeBool(value)
		if err != nil {
			return []byte{}, err
		}
		return EVMEncodeBytes(data), nil

	default:
		return []byte{}, fmt.Errorf("unsupported format: %s", format)
	}
}

// EVMTranscodeBytes converts a json input to an EVM bytes array
func EVMTranscodeBytes(value gjson.Result) ([]byte, error) {
	switch value.Type {
	case gjson.String:
		return EVMEncodeBytes([]byte(value.Str)), nil

	case gjson.False:
		return EVMEncodeBytes(EVMWordUint64(0)), nil

	case gjson.True:
		return EVMEncodeBytes(EVMWordUint64(1)), nil

	case gjson.Number:
		v := big.NewFloat(value.Num) // precision limited to float64
		vInt, _ := v.Int(nil)
		word, err := EVMWordSignedBigInt(vInt)
		if err != nil {
			return nil, errors.Wrap(err, "while converting float to int256")
		}
		return EVMEncodeBytes(word), nil
	default:
		return []byte{}, fmt.Errorf("unsupported encoding for value: %s", value.Type)
	}
}
