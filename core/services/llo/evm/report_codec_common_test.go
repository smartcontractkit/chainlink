package evm

import (
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-data-streams/llo"
	ubig "github.com/smartcontractkit/chainlink-integrations/evm/utils/big"
)

func TestEncodePackedBigInt(t *testing.T) {
	testCases := []struct {
		value       *big.Int
		typeStr     string
		expectedHex string // expected output as lowercase hex string without 0x prefix
		shouldErr   bool
	}{
		// Valid cases:
		{big.NewInt(100), "uint8", "64", false},
		{big.NewInt(100), "uint256", "0000000000000000000000000000000000000000000000000000000000000064", false},
		{big.NewInt(-1), "int8", "ff", false},          // -1 mod 256 = 0xff
		{big.NewInt(-100), "int32", "ffffff9c", false}, // -100 mod 2^32
		{big.NewInt(123456789), "uint32", "075bcd15", false},
		{big.NewInt(123456789), "int160", "00000000000000000000000000000000075bcd15", false},
		// For a 192-bit unsigned integer, 24 bytes (48 hex digits)
		{big.NewInt(100), "uint192", "000000000000000000000000000000000000000000000064", false},
		// For a 192-bit signed integer; -100 mod 2^192 = 2^192 - 100.
		// The expected value is (2^192 - 100) in 24 bytes, which is 23 bytes of "ff" followed by "9c".
		{big.NewInt(-100), "int192", "ffffffffffffffffffffffffffffffffffffffffffffff9c", false},
		// For a 256-bit signed integer; -1 mod 2^256 = 2^256 - 1 (all bytes "ff")
		{big.NewInt(-1), "int256", "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", false},
		// For a 256-bit signed integer with a positive value.
		// 123456789 in hex is 075bcd15; padded to 32 bytes (64 hex digits)
		{big.NewInt(123456789), "int256", "00000000000000000000000000000000000000000000000000000000075bcd15", false},

		// Error cases:
		{big.NewInt(256), "uint8", "", true}, // 256 does not fit in 8 bits
		{big.NewInt(-1), "uint8", "", true},  // negative value for unsigned type
		{big.NewInt(100), "int7", "", true},  // invalid Solidity type (bitwidth not supported)
	}

	for _, tc := range testCases {
		result, err := EncodePackedBigInt(tc.value, tc.typeStr)
		if tc.shouldErr {
			require.Error(t, err, "expected error for value %s with type %s", tc.value, tc.typeStr)
		} else {
			require.NoError(t, err, "unexpected error for value %s with type %s", tc.value, tc.typeStr)
			hexResult := hex.EncodeToString(result)
			assert.Equal(t, tc.expectedHex, hexResult, "For value %s and type %s", tc.value, tc.typeStr)
		}
	}
}

func Test_ABIEncoder_EncodePacked(t *testing.T) {
	t.Run("encodes decimals", func(t *testing.T) {
		enc := ABIEncoder{
			Type:       "uint192",
			Multiplier: ubig.NewI(10000),
		}
		encoded, err := enc.EncodePadded(tests.Context(t), llo.ToDecimal(decimal.NewFromFloat32(123456.789123)))
		require.NoError(t, err)
		assert.Equal(t, "00000000000000000000000000000000000000000000000000000000499602dc", hex.EncodeToString(encoded))
	})
	t.Run("errors on unsupported type (e.g. Quote)", func(t *testing.T) {
		enc := ABIEncoder{
			Type: "Quote",
		}
		_, err := enc.EncodePacked(tests.Context(t), &llo.Quote{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "EncodePacked only currently supports StreamValue type of *llo.Decimal")
	})
}

func Test_ABIEncoder_EncodePadded(t *testing.T) {
	t.Run("encodes decimals", func(t *testing.T) {
		tcs := []struct {
			name       string
			sv         llo.StreamValue
			abiType    string
			multiplier *big.Int
			errStr     string
			expected   string
		}{
			{
				name:    "overflow int8",
				sv:      llo.ToDecimal(decimal.NewFromFloat32(123456789.123456789)),
				abiType: "int8",
				errStr:  "invalid type: cannot fit 123456790 into int8",
			},
			{
				name:     "successful int8",
				sv:       llo.ToDecimal(decimal.NewFromFloat32(123.456)),
				abiType:  "int8",
				expected: padLeft32Byte(fmt.Sprintf("%x", 123)),
			},
			{
				name:       "negative multiplied int8",
				sv:         llo.ToDecimal(decimal.NewFromFloat32(1.11)),
				multiplier: big.NewInt(-100),
				abiType:    "int8",
				expected:   "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff91",
			},
			{
				name:    "negative uint32",
				sv:      llo.ToDecimal(decimal.NewFromFloat32(-123.456)),
				abiType: "uint32",
				errStr:  "invalid type: cannot fit -123 into uint32",
			},
			{
				name:     "successful uint32",
				sv:       llo.ToDecimal(decimal.NewFromFloat32(123456.456)),
				abiType:  "uint32",
				expected: padLeft32Byte(fmt.Sprintf("%x", 123456)),
			},
			{
				name:       "multiplied uint32",
				sv:         llo.ToDecimal(decimal.NewFromFloat32(123.456)),
				multiplier: big.NewInt(100),
				abiType:    "uint32",
				expected:   padLeft32Byte(fmt.Sprintf("%x", 12345)),
			},
			{
				name:       "negative multiplied int32",
				sv:         llo.ToDecimal(decimal.NewFromFloat32(123.456)),
				multiplier: big.NewInt(-100),
				abiType:    "int32",
				expected:   "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffcfc7",
			},
			{
				name:       "overflowing multiplied int32",
				sv:         llo.ToDecimal(decimal.NewFromInt(math.MaxInt32)),
				multiplier: big.NewInt(2),
				abiType:    "int32",
				errStr:     "invalid type: cannot fit 4294967294 into int32",
			},
			{
				name:       "successful int192",
				sv:         llo.ToDecimal(decimal.NewFromFloat32(123456.789123)),
				abiType:    "int192",
				multiplier: big.NewInt(1e18),
				expected:   "000000000000000000000000000000000000000000001a249b2292e49d8f0000",
			},
		}
		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				enc := ABIEncoder{
					Type:       tc.abiType,
					Multiplier: (*ubig.Big)(tc.multiplier),
				}
				encoded, err := enc.EncodePadded(tests.Context(t), tc.sv)
				if tc.errStr != "" {
					require.Error(t, err)
					assert.Contains(t, err.Error(), tc.errStr)
				} else {
					require.NoError(t, err)
					require.Equal(t, tc.expected, hex.EncodeToString(encoded))
				}
			})
		}
	})
}
