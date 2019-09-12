package adapters

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/tidwall/gjson"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var isSupportedABITypeTests = []struct {
	typeString string
	valid      bool
}{
	// primitive types
	{"address", true},
	{"bool", true},
	{"bytes", true},
	// abi.NewType will interpret this as `bytes`
	// {"bytes0", false},
	{"bytes1", true},
	{"bytes32", true},
	{"bytes33", false},
	{"int7", false},
	{"int8", true},
	{"int33", false},
	{"int256", true},
	{"int264", false},
	{"string", true},
	{"uint0", false},
	{"uint8", true},
	{"uint256", true},
	{"uint257", false},
	// arrays
	{"address[3]", true},
	{"address[3][3]", true},
	{"bytes[2]", false},
	{"bytes32[3][3]", true},
	{"uint256[2][3]", true},
	// slices
	{"bytes[]", false},
	{"int128[]", true},
	{"string[]", false},
	{"uint256[2][3][]", true},
	{"uint256[][]", false},
}

func TestEthTxABIEncodeAdapter_isSupportedABIType(t *testing.T) {
	for _, test := range isSupportedABITypeTests {
		typ, err := abi.NewType(test.typeString, []abi.ArgumentMarshaling{})
		assert.NoError(t, err)

		require.Equal(t, test.valid, isSupportedABIType(&typ), "failed for %s", test.typeString)
	}
}

var encodeTests = []struct {
	desc       string
	abiJSON    string
	resultJSON string
	hexEncoded string // leave empty to signal that encode is expected to fail
}{
	{
		"empty result should fail",
		`[{"inputs":[{"name":"a","type":"uint8"}],"name":"foo","type":"function"}]`,
		`{}`,
		``,
	},
	{
		"result with wrong key should fail",
		`[{"inputs":[{"name":"a","type":"uint8"}],"name":"foo","type":"function"}]`,
		`{"b": "0xf"}`,
		``,
	},
	{
		"result with extra key should fail",
		`[{"inputs":[{"name":"a","type":"uint8"}],"name":"foo","type":"function"}]`,
		`{"a": "0xf", "b": "0xf"}`,
		``,
	},
	// Testvectors
	{
		"testvec 1 from https://solidity.readthedocs.io/en/v0.5.11/abi-spec.html#examples",
		`[{"inputs":[{"name":"x","type":"uint32"},{"name":"y","type":"bool"}],"name":"baz","type":"function"}]`,
		`{"x": 69, "y": true}`,
		`cdcd77c000000000000000000000000000000000000000000000000000000000000000450000000000000000000000000000000000000000000000000000000000000001`,
	},
	{
		"testvec 2 from https://solidity.readthedocs.io/en/v0.5.11/abi-spec.html#examples",
		`[{"inputs":[{"name":"x","type":"bytes3[2]"}],"name":"bar","type":"function"}]`,
		`{"x": [[97, 98, 99], [100, 101, 102]]}`,
		`fce353f661626300000000000000000000000000000000000000000000000000000000006465660000000000000000000000000000000000000000000000000000000000`,
	},
	{
		"testvec 3 from https://solidity.readthedocs.io/en/v0.5.11/abi-spec.html#examples",
		`[{"inputs":[{"name":"x","type":"bytes"},{"name":"y","type":"bool"},{"name":"z","type":"uint256[]"}],"name":"sam","type":"function"}]`,
		`{"x": "0x64617665", "y": true, "z": ["1","2","3"]}`,
		`a5643bf20000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000000464617665000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000003`,
	},
	{
		"testvec 1 from https://solidity.readthedocs.io/en/v0.5.11/abi-spec.html#use-of-dynamic-types",
		`[{"inputs": [
			{"name":"a","type":"uint256"},
			{"name":"b","type":"uint32[]"},
			{"name":"c","type":"bytes10"},
			{"name":"d","type":"bytes"}
		  ],
		  "name":"f", "type":"function"}]`,
		`{
			"a": "0x123", 
			"b": ["0x456", "0x789"], 
			"c": "0x31323334353637383930", 
			"d": [72, 101, 108, 108, 111, 44, 32, 119, 111, 114, 108, 100, 33]
		}`,
		`8be6524600000000000000000000000000000000000000000000000000000000000001230000000000000000000000000000000000000000000000000000000000000080313233343536373839300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e0000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000004560000000000000000000000000000000000000000000000000000000000000789000000000000000000000000000000000000000000000000000000000000000d48656c6c6f2c20776f726c642100000000000000000000000000000000000000`,
	},
	// Address
	{
		"valid address",
		`[{"inputs":[{"name":"a","type":"address"}],"name":"foo","type":"function"}]`,
		`{"a": "0x98d60255f917e3eb94eae199d827dad837fac4cb"}`,
		`fdf80bda00000000000000000000000098d60255f917e3eb94eae199d827dad837fac4cb`,
	},
	{
		"invalid address",
		`[{"inputs":[{"name":"a","type":"address"}],"name":"foo","type":"function"}]`,
		`{"a": "98d60255f917e3eb94eae199d827dad837fac4cb"}`,
		``,
	},
	{
		"short address",
		`[{"inputs":[{"name":"a","type":"address"}],"name":"foo","type":"function"}]`,
		`{"a": "0xf917e3eb94eae199d827dad837fac4cb"}`,
		`fdf80bda00000000000000000000000000000000f917e3eb94eae199d827dad837fac4cb`,
	},
	{
		"address too long",
		`[{"inputs":[{"name":"a","type":"address"}],"name":"foo","type":"function"}]`,
		`{"a": "0xffffffffffffffffffffffffffffffffffffffffff"}`,
		``,
	},
	// Array
	{
		"valid array",
		`[{"inputs":[{"name":"a","type":"address[2]"}],"name":"foo","type":"function"}]`,
		`{"a": ["0x98d60255f917e3eb94eae199d827dad837fac4cb", "0x98d60255f917e3eb94eae199d827dad837fac4cd"]}`,
		`8d833f3000000000000000000000000098d60255f917e3eb94eae199d827dad837fac4cb00000000000000000000000098d60255f917e3eb94eae199d827dad837fac4cd`,
	},
	{
		"nested array",
		`[{"inputs":[{"name":"a","type":"uint128[2][4]"}],"name":"foo","type":"function"}]`,
		`{"a": [["1", "2"], ["0x3", "0x4"], ["5", "6"], ["7", "8"]]}`,
		`9c78d0ac00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000005000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000070000000000000000000000000000000000000000000000000000000000000008`,
	},
	// Bool
	{
		"valid bool",
		`[{"inputs":[{"name":"a","type":"bool"},{"name":"b","type":"bool"}],"name":"foo","type":"function"}]`,
		`{"a": true, "b": false}`,
		`b3cedfcf00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000`,
	},
	{
		"invalid bool",
		`[{"inputs":[{"name":"a","type":"bool"},{"name":"b","type":"bool"}],"name":"foo","type":"function"}]`,
		`{"a": 1, "b": null}`,
		``,
	},
	// Bytes (fixed size)
	{
		"bytes1 and bytes32",
		`[{"inputs":[{"name":"a","type":"bytes1"},{"name":"b","type":"bytes32"}],"name":"foo","type":"function"}]`,
		`{"a": "0x12","b": "0xffbb22aaaccaaaa00aaaa13aaa88d60255f91243eb94eae1943827dad837fac4"}`,
		`296874791200000000000000000000000000000000000000000000000000000000000000ffbb22aaaccaaaa00aaaa13aaa88d60255f91243eb94eae1943827dad837fac4`,
	},
	{
		"bytes32 too short",
		`[{"inputs":[{"name":"a","type":"bytes32"}],"name":"foo","type":"function"}]`,
		`{"a": "0xaaaccaaaa00aaaa13aaa88d60255f91243eb94eae1943827dad837fac4"}`,
		``,
	},
	{
		"bytes32 too long",
		`[{"inputs":[{"name":"a","type":"bytes32"}],"name":"foo","type":"function"}]`,
		`{"a": "0xffffbb22aaaccaaaa00aaaa13aaa88d60255f91243eb94eae1943827dad837fac4"}`,
		``,
	},
	// Bytes (variable size)
	{
		"valid bytes",
		`[{"inputs":[{"name":"a","type":"bytes"}],"name":"foo","type":"function"}]`,
		`{"a": "0x12"}`,
		`30c8d1da000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000011200000000000000000000000000000000000000000000000000000000000000`,
	},
	{
		"valid bytes",
		`[{"inputs":[{"name":"a","type":"bytes"}],"name":"foo","type":"function"}]`,
		`{"a": "0xaaaccaaaa0aaaaaaaa0aaaa13aaa88d60255f91243eb94eae1943827dad837fac4"}`,
		`30c8d1da00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000021aaaccaaaa0aaaaaaaa0aaaa13aaa88d60255f91243eb94eae1943827dad837fac400000000000000000000000000000000000000000000000000000000000000`,
	},
	{
		"valid bytes",
		`[{"inputs":[{"name":"a","type":"bytes"}],"name":"foo","type":"function"}]`,
		`{"a": [1,2,3]}`,
		`30c8d1da000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000030102030000000000000000000000000000000000000000000000000000000000`,
	},
	{
		"invalid bytes",
		`[{"inputs":[{"name":"a","type":"bytes"}],"name":"foo","type":"function"}]`,
		`{"a": "aabbcc"}`,
		``,
	},
	{
		"invalid bytes",
		`[{"inputs":[{"name":"a","type":"bytes"}],"name":"foo","type":"function"}]`,
		`{"a": [0,1,2,3,"p"]}`,
		``,
	},
	// Int
	{
		"value too small for int8",
		`[{"inputs":[{"name":"a","type":"int8"}],"name":"foo","type":"function"}]`,
		`{"b": "-129"}`,
		``,
	},
	{
		"value too large for int8",
		`[{"inputs":[{"name":"a","type":"int8"}],"name":"foo","type":"function"}]`,
		`{"b": "128"}`,
		``,
	},
	{
		"min and max int8",
		`[{"inputs":[{"name":"a","type":"int8"},{"name":"b","type":"int8"}],"name":"foo","type":"function"}]`,
		`{"a": "-128", "b": "0x7f"}`,
		`0affedbfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80000000000000000000000000000000000000000000000000000000000000007f`,
	},
	{
		"int32 and int128",
		`[{"inputs":[{"name":"a","type":"int32"},{"name":"b","type":"int128"}],"name":"foo","type":"function"}]`,
		`{"a": "0xfff", "b": "-170141183460469231731687303715884105728"}`,
		`a1f401870000000000000000000000000000000000000000000000000000000000000fffffffffffffffffffffffffffffffffff80000000000000000000000000000000`,
	},
	{
		"value too large for uint256",
		`[{"inputs":[{"name":"a","type":"uint256"}],"name":"foo","type":"function"}]`,
		`{"a": "0xffaaaaaaaaaaaaaaaaaaaaaaaa88d60255f917e3eb94eae199d827dad837fac4cb"}`,
		``,
	},
	// Slice
	{
		"simple slice",
		`[{"inputs":[{"name":"a","type":"bool[]"}],"name":"foo","type":"function"}]`,
		`{"a": [true, false, true, true]}`,
		`78a4a116000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001`,
	},
	{
		"complex slice",
		`[{"inputs":[{"name":"a","type":"int32[2][]"}],"name":"foo","type":"function"}]`,
		`{"a": [[-12, 12], [17, "0xabc"], ["-2", "-2147483648"]]}`,
		`c291cfa600000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000003fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff4000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000110000000000000000000000000000000000000000000000000000000000000abcfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffeffffffffffffffffffffffffffffffffffffffffffffffffffffffff80000000`,
	},
	// String
	{
		"valid string",
		`[{"inputs":[{"name":"a","type":"string"}],"name":"foo","type":"function"}]`,
		`{"a": "asdf"}`,
		`f31a6969000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000046173646600000000000000000000000000000000000000000000000000000000`,
	},
	{
		"valid empty string",
		`[{"inputs":[{"name":"a","type":"string"}],"name":"foo","type":"function"}]`,
		`{"a": ""}`,
		`f31a696900000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000000`,
	},
	// Uint
	{
		"value too large for uint",
		`[{"inputs":[{"name":"a","type":"uint8"}],"name":"foo","type":"function"}]`,
		`{"a": "0x100"}`,
		``,
	},
	{
		"hex without 0x prefix for uint",
		`[{"inputs":[{"name":"a","type":"uint8"}],"name":"foo","type":"function"}]`,
		`{"a": "ff"}`,
		``,
	},
	{
		"value too small for uint",
		`[{"inputs":[{"name":"a","type":"uint8"}],"name":"foo","type":"function"}]`,
		`{"a": "-1"}`,
		``,
	},
	{
		"different encodings for uint",
		`[{"inputs":[{"name":"a","type":"uint8"}, {"name":"b","type":"uint8"}],"name":"foo","type":"function"}]`,
		`{"a": "0x00", "b": "255"}`,
		`7ce60e3b000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ff`,
	},
	{
		"different encodings for uint",
		`[{"inputs":[{"name":"a","type":"uint8"}, {"name":"b","type":"uint8"}],"name":"foo","type":"function"}]`,
		`{"a": "0xa", "b": "0"}`,
		`7ce60e3b000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000000`,
	},
	{
		"different encodings of max uint256",
		`[{"inputs":[{"name":"a","type":"uint256"}, {"name":"b","type":"uint256"}],"name":"foo","type":"function"}]`,
		`{
			"a": "115792089237316195423570985008687907853269984665640564039457584007913129639935",
			"b": "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
		}`,
		`04bc52f8ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff`,
	},
	{
		"uint48 as number and string",
		`[{"inputs":[{"name":"a","type":"uint48"},{"name":"b","type":"uint48"}],"name":"foo","type":"function"}]`,
		`{"a": 140737488355328, "b": "0x800000000001"}`,
		`183cb15b00000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000800000000001`,
	},
	{
		"uint56 doesn't accept numbers",
		`[{"inputs":[{"name":"a","type":"uint56"}],"name":"foo","type":"function"}]`,
		`{"a": 1"}`,
		``,
	},
}

func TestEthTxABIEncodeAdapter_encode(t *testing.T) {
	for i, test := range encodeTests {
		testABI, err := abi.JSON(strings.NewReader(test.abiJSON))
		assert.NoError(t, err)

		// there should be a single method, get its abi
		assert.Equal(t, 1, len(testABI.Methods))
		var fnABI abi.Method
		for _, fnABI = range testABI.Methods {
		}

		args, ok := gjson.Parse(test.resultJSON).Value().(map[string]interface{})
		assert.True(t, ok, "Failed to parse resultJSON in test %v: %s", i, test.desc)

		encoded, err := abiEncode(&fnABI, args)
		_ = encoded
		if test.hexEncoded != "" {
			assert.NoError(t, err, "in test %v: %s", i, test.desc)
			if test.hexEncoded != "?" {
				require.Equal(t, test.hexEncoded, hex.EncodeToString(encoded), "in test %v: %s", i, test.desc)
			}
		} else {
			require.Error(t, err, "in test %v: %s", i, test.desc)
		}
	}
}
