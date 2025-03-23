package evm

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
		{big.NewInt(256), "int8", "", true},  // 256 does not fit in 8 bits
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

func Test_ABIEncoder_UnmarshalJSON(t *testing.T) {
	j := `[{"type":"uint192","multiplier":"10000"},[{"type":"uint64"},{"type":"int192","multiplier":"100"}]]`
	encs := []ABIEncoder{}
	err := json.Unmarshal([]byte(j), &encs)
	require.NoError(t, err)
	require.Len(t, encs, 2)

	t.Run("padded encoding", func(t *testing.T) {
		// simple decimal case
		enc := encs[0]
		assert.Equal(t, "uint192", enc.encoders[0].Type)
		assert.Equal(t, "10000", enc.encoders[0].Multiplier.String())

		res, err := enc.EncodePadded(llo.ToDecimal(decimal.NewFromFloat32(123456.789123)))
		require.NoError(t, err)
		assert.Len(t, res, 32)
		assert.Equal(t, "00000000000000000000000000000000000000000000000000000000499602dc", hex.EncodeToString(res))

		// nested array timestamped stream value case
		enc = encs[1]
		assert.Equal(t, "uint64", enc.encoders[0].Type)
		assert.Equal(t, "<nil>", enc.encoders[0].Multiplier.String())
		assert.Equal(t, "int192", enc.encoders[1].Type)
		assert.Equal(t, "100", enc.encoders[1].Multiplier.String())

		res, err = enc.EncodePadded(&llo.TimestampedStreamValue{
			ObservedAtNanoseconds: 0x123456789,
			StreamValue:           llo.ToDecimal(decimal.NewFromFloat32(123456.789123)),
		})
		require.NoError(t, err)
		assert.Len(t, res, 64)
		assert.Equal(t, "00000000000000000000000000000000000000000000000000000001234567890000000000000000000000000000000000000000000000000000000000bc614f", hex.EncodeToString(res))
	})
	t.Run("packed encoding", func(t *testing.T) {
		// simple decimal case
		enc := encs[0]
		assert.Equal(t, "uint192", enc.encoders[0].Type)
		assert.Equal(t, "10000", enc.encoders[0].Multiplier.String())

		res, err := enc.EncodePacked(llo.ToDecimal(decimal.NewFromFloat32(123456.789123)))
		require.NoError(t, err)
		assert.Len(t, res, 24)
		assert.Equal(t, "0000000000000000000000000000000000000000499602dc", hex.EncodeToString(res))

		// nested array timestamped stream value case
		enc = encs[1]
		assert.Equal(t, "uint64", enc.encoders[0].Type)
		assert.Equal(t, "<nil>", enc.encoders[0].Multiplier.String())
		assert.Equal(t, "int192", enc.encoders[1].Type)
		assert.Equal(t, "100", enc.encoders[1].Multiplier.String())

		res, err = enc.EncodePacked(&llo.TimestampedStreamValue{
			ObservedAtNanoseconds: 0x123456789,
			StreamValue:           llo.ToDecimal(decimal.NewFromFloat32(123456.789123)),
		})
		require.NoError(t, err)
		assert.Len(t, res, 32)
		assert.Equal(t, "0000000123456789000000000000000000000000000000000000000000bc614f", hex.EncodeToString(res))
	})
}

func Test_ABIEncoder_EncodePacked(t *testing.T) {
	t.Run("encodes decimals", func(t *testing.T) {
		enc := ABIEncoder{
			encoders: []singleABIEncoder{{
				Type:       "uint192",
				Multiplier: ubig.NewI(10000),
			}},
		}
		encoded, err := enc.EncodePadded(llo.ToDecimal(decimal.NewFromFloat32(123456.789123)))
		require.NoError(t, err)
		assert.Equal(t, "00000000000000000000000000000000000000000000000000000000499602dc", hex.EncodeToString(encoded))
	})
	t.Run("errors on unsupported type (e.g. Quote)", func(t *testing.T) {
		enc := ABIEncoder{
			encoders: []singleABIEncoder{{
				Type: "Quote",
			}},
		}
		_, err := enc.EncodePacked(&llo.Quote{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unhandled type; supported types are: *llo.Decimal or *llo.TimestampedStreamValue; got: *llo.Quote")
	})
}

func Test_ABIEncoder_EncodePadded_EncodePacked(t *testing.T) {
	t.Run("encodes decimals", func(t *testing.T) {
		tcs := []struct {
			name           string
			sv             llo.StreamValue
			abiType        string
			multiplier     *big.Int
			errStr         string
			expectedPadded string
			expectedPacked string
		}{
			{
				name:    "overflow int8",
				sv:      llo.ToDecimal(decimal.NewFromFloat32(123456789.123456789)),
				abiType: "int8",
				errStr:  "value 123456790 out of range for type int8",
			},
			{
				name:           "successful int8",
				sv:             llo.ToDecimal(decimal.NewFromFloat32(123.456)),
				abiType:        "int8",
				expectedPadded: padLeft32Byte(fmt.Sprintf("%x", 123)),
				expectedPacked: fmt.Sprintf("%x", 123),
			},
			{
				name:           "negative multiplied int8",
				sv:             llo.ToDecimal(decimal.NewFromFloat32(1.11)),
				multiplier:     big.NewInt(-100),
				abiType:        "int8",
				expectedPadded: "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff91",
				expectedPacked: "91",
			},
			{
				name:           "negative int192",
				sv:             llo.ToDecimal(decimal.NewFromFloat32(1.11)),
				multiplier:     big.NewInt(-100),
				abiType:        "int192",
				expectedPadded: "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff91",
				expectedPacked: "ffffffffffffffffffffffffffffffffffffffffffffff91",
			},
			{
				name:    "negative uint32",
				sv:      llo.ToDecimal(decimal.NewFromFloat32(-123.456)),
				abiType: "uint32",
				errStr:  "negative value provided for unsigned type uint32",
			},
			{
				name:           "successful uint32",
				sv:             llo.ToDecimal(decimal.NewFromFloat32(123456.456)),
				abiType:        "uint32",
				expectedPadded: padLeft32Byte(fmt.Sprintf("%x", 123456)),
				expectedPacked: "0001e240",
			},
			{
				name:           "multiplied uint32",
				sv:             llo.ToDecimal(decimal.NewFromFloat32(123.456)),
				multiplier:     big.NewInt(100),
				abiType:        "uint32",
				expectedPadded: padLeft32Byte(fmt.Sprintf("%x", 12345)),
				expectedPacked: "00003039",
			},
			{
				name:           "negative multiplied int32",
				sv:             llo.ToDecimal(decimal.NewFromFloat32(123.456)),
				multiplier:     big.NewInt(-100),
				abiType:        "int32",
				expectedPadded: "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffcfc7",
				expectedPacked: "ffffcfc7",
			},
			{
				name:       "overflowing multiplied int32",
				sv:         llo.ToDecimal(decimal.NewFromInt(math.MaxInt32)),
				multiplier: big.NewInt(2),
				abiType:    "int32",
				errStr:     "value 4294967294 out of range for type int32",
			},
			{
				name:           "successful int192",
				sv:             llo.ToDecimal(decimal.NewFromFloat32(123456.789123)),
				abiType:        "int192",
				multiplier:     big.NewInt(1e18),
				expectedPadded: "000000000000000000000000000000000000000000001a249b2292e49d8f0000",
				expectedPacked: "00000000000000000000000000001a249b2292e49d8f0000",
			},
			{
				name:    "invalid type",
				sv:      llo.ToDecimal(decimal.NewFromFloat32(123.456)),
				abiType: "blah",
				errStr:  "invalid Solidity type: blah",
			},
			{
				name:    "invalid type",
				sv:      llo.ToDecimal(decimal.NewFromFloat32(123.456)),
				abiType: "int",
				errStr:  "invalid Solidity type: int",
			},
		}
		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				enc := ABIEncoder{
					[]singleABIEncoder{{
						Type:       tc.abiType,
						Multiplier: (*ubig.Big)(tc.multiplier),
					}},
				}
				padded, err := enc.EncodePadded(tc.sv)
				if tc.errStr != "" {
					require.Error(t, err)
					assert.Contains(t, err.Error(), tc.errStr)
				} else {
					assert.Len(t, padded, 32)
					require.NoError(t, err)
					require.Equal(t, tc.expectedPadded, hex.EncodeToString(padded))
				}
				packed, err := enc.EncodePacked(tc.sv)
				if tc.errStr != "" {
					require.Error(t, err)
					assert.Contains(t, err.Error(), tc.errStr)
				} else {
					require.NoError(t, err)
					require.Equal(t, tc.expectedPacked, hex.EncodeToString(packed))
				}
			})
		}
	})
	t.Run("encodes TimestampedStreamValue", func(t *testing.T) {
		tcs := []struct {
			name           string
			sv             llo.StreamValue
			encoder        ABIEncoder
			paddedErrStr   string
			expectedPadded string
			packedErrStr   string
			expectedPacked string
		}{
			{
				name: "with valid ABI types",
				sv: &llo.TimestampedStreamValue{
					ObservedAtNanoseconds: 0x123456789,
					StreamValue:           llo.ToDecimal(decimal.NewFromFloat32(123456.789123)),
				},
				encoder: ABIEncoder{
					[]singleABIEncoder{{
						Type:       "uint64",
						Multiplier: nil,
					}, {
						Type:       "int192",
						Multiplier: ubig.NewI(1e18),
					}},
				},
				expectedPadded: "0000000000000000000000000000000000000000000000000000000123456789000000000000000000000000000000000000000000001a249b2292e49d8f0000",
				expectedPacked: "000000012345678900000000000000000000000000001a249b2292e49d8f0000",
			},
			{
				name: "passing bytes0 as the first type serializes only the nested value",
				sv: &llo.TimestampedStreamValue{
					ObservedAtNanoseconds: 0x123456789,
					StreamValue:           llo.ToDecimal(decimal.NewFromFloat32(123456.789123)),
				},
				encoder: ABIEncoder{
					[]singleABIEncoder{{
						Type:       "bytes0",
						Multiplier: nil,
					}, {
						Type:       "int192",
						Multiplier: ubig.NewI(1e18),
					}},
				},
				paddedErrStr:   "invalid Solidity type: bytes0",
				expectedPacked: "00000000000000000000000000001a249b2292e49d8f0000",
			},
			{
				name: "passing bytes0 as the second type serializes only the timestamp",
				sv: &llo.TimestampedStreamValue{
					ObservedAtNanoseconds: 0x123456789,
					StreamValue:           llo.ToDecimal(decimal.NewFromFloat32(123456.789123)),
				},
				encoder: ABIEncoder{
					[]singleABIEncoder{{
						Type: "uint64",
					}, {
						Type: "bytes0",
					}},
				},
				paddedErrStr:   "invalid Solidity type: bytes0",
				expectedPacked: "0000000123456789",
			},
			{
				name: "multiplier applied to both types",
				sv: &llo.TimestampedStreamValue{
					ObservedAtNanoseconds: 0x123456789,
					StreamValue:           llo.ToDecimal(decimal.NewFromFloat32(123456.789123)),
				},
				encoder: ABIEncoder{
					[]singleABIEncoder{{
						Type:       "uint64",
						Multiplier: ubig.NewI(100),
					}, {
						Type:       "int192",
						Multiplier: ubig.NewI(100),
					}},
				},
				expectedPadded: "00000000000000000000000000000000000000000000000000000071c71c71840000000000000000000000000000000000000000000000000000000000bc614f",
				expectedPacked: "00000071c71c7184000000000000000000000000000000000000000000bc614f",
			},
			{
				name: "errors if timestamp would overflow",
				sv: &llo.TimestampedStreamValue{
					ObservedAtNanoseconds: 0x123456789,
					StreamValue:           llo.ToDecimal(decimal.NewFromFloat(float64(200))),
				},
				encoder: ABIEncoder{
					[]singleABIEncoder{{
						Type: "uint8",
					}, {
						Type: "int192",
					}},
				},
				paddedErrStr: "value 4886718345 out of range for type uint8",
				packedErrStr: "value 4886718345 out of range for type uint8",
			},
			{
				name: "errors if stream value would overflow",
				sv: &llo.TimestampedStreamValue{
					ObservedAtNanoseconds: 0x123456789,
					StreamValue:           llo.ToDecimal(decimal.NewFromFloat(float64(200))),
				},
				encoder: ABIEncoder{
					[]singleABIEncoder{{
						Type: "uint64",
					}, {
						Type:       "int8",
						Multiplier: ubig.NewI(1e18),
					}},
				},
				paddedErrStr: "value 200000000000000000000 out of range for type int8",
				packedErrStr: "value 200000000000000000000 out of range for type int8",
			},
			{
				name: "unsupported nested type",
				sv: &llo.TimestampedStreamValue{
					ObservedAtNanoseconds: 0x123456789,
					StreamValue:           &llo.Quote{},
				},
				encoder: ABIEncoder{
					[]singleABIEncoder{{
						Type: "uint64",
					}, {
						Type: "int192",
					}},
				},
				paddedErrStr: "unhandled type; supported nested types for *llo.TimestampedStreamValue are: *llo.Decimal; got: *llo.Quote",
				packedErrStr: "failed to encode nested stream value; EncodePacked only currently supports StreamValue type of *llo.Decimal, got: *llo.Quote",
			},
			{
				name: "successfully encodes timestamped stream value",
				sv: &llo.TimestampedStreamValue{
					ObservedAtNanoseconds: 0x123456789,
					StreamValue:           llo.ToDecimal(decimal.NewFromFloat32(123456.789123)),
				},
				encoder: ABIEncoder{
					[]singleABIEncoder{{
						Type: "uint64",
					}, {
						Type:       "int192",
						Multiplier: ubig.NewI(1e18),
					}},
				},
				expectedPadded: "0000000000000000000000000000000000000000000000000000000123456789000000000000000000000000000000000000000000001a249b2292e49d8f0000",
				expectedPacked: "000000012345678900000000000000000000000000001a249b2292e49d8f0000",
			},
			{
				name: "wrong abi encoder",
				sv: &llo.TimestampedStreamValue{
					ObservedAtNanoseconds: 0x123456789,
					StreamValue:           llo.ToDecimal(decimal.NewFromFloat32(123456.789123)),
				},
				encoder: ABIEncoder{
					[]singleABIEncoder{{
						Type: "uint64",
					}},
				},
				paddedErrStr: "expected exactly two encoders for *llo.TimestampedStreamValue; got: 1",
				packedErrStr: "expected exactly two encoders for *llo.TimestampedStreamValue; got: 1",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				enc := tc.encoder
				padded, err := enc.EncodePadded(tc.sv)
				if tc.paddedErrStr != "" {
					require.Error(t, err)
					assert.Contains(t, err.Error(), tc.paddedErrStr)
				} else {
					assert.Len(t, padded, 64)
					require.NoError(t, err)
					assert.Equal(t, tc.expectedPadded, hex.EncodeToString(padded))
				}
				packed, err := enc.EncodePacked(tc.sv)
				if tc.packedErrStr != "" {
					require.Error(t, err)
					assert.Contains(t, err.Error(), tc.packedErrStr)
				} else {
					require.NoError(t, err)
					assert.Equal(t, tc.expectedPacked, hex.EncodeToString(packed))
				}
			})
		}
	})
}
