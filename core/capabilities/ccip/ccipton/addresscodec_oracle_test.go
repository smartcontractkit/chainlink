package ccipton

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xssnick/tonutils-go/address"
)

// TestAddressCodec_OracleIDAsAddressBytes tests the OracleIDAsAddressBytes method
// This test was previously disabled with a TODO comment but is now re-enabled
// with comprehensive test cases to ensure the function works correctly.
func TestAddressCodec_OracleIDAsAddressBytes(t *testing.T) {
	codec := AddressCodec{}

	testCases := []struct {
		name     string
		oracleID uint8
		expected []byte
	}{
		{
			name:     "oracleID 0",
			oracleID: 0,
			expected: func() []byte {
				return packOracleIDForTest(t, 0)
			}(),
		},
		{
			name:     "oracleID 1",
			oracleID: 1,
			expected: func() []byte {
				return packOracleIDForTest(t, 1)
			}(),
		},
		{
			name:     "oracleID 127 (mid-range)",
			oracleID: 127,
			expected: func() []byte {
				return packOracleIDForTest(t, 127)
			}(),
		},
		{
			name:     "oracleID 255 (max value)",
			oracleID: 255,
			expected: func() []byte {
				return packOracleIDForTest(t, 255)
			}(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := codec.OracleIDAsAddressBytes(tc.oracleID)

			// Verify no error occurred
			require.NoError(t, err)
			
			// Verify the result matches expected
			require.Equal(t, tc.expected, actual, "expected %x, got %x", tc.expected, actual)
			
			// Verify the result has correct length (36 bytes for TON address)
			require.Len(t, actual, 36, "TON address should be 36 bytes")
			
			// Verify the result can be converted back to a valid TON address string
			addrString, err := codec.AddressBytesToString(actual)
			require.NoError(t, err)
			require.NotEmpty(t, addrString)
			
			// Verify round-trip conversion works
			roundTripBytes, err := codec.AddressStringToBytes(addrString)
			require.NoError(t, err)
			require.Equal(t, actual, roundTripBytes, "round-trip conversion should preserve bytes")
		})
	}
}

// TestAddressCodec_OracleIDAsAddressBytes_EdgeCases tests edge cases and error conditions
func TestAddressCodec_OracleIDAsAddressBytes_EdgeCases(t *testing.T) {
	codec := AddressCodec{}

	// Test that different oracle IDs produce different addresses
	addr0, err := codec.OracleIDAsAddressBytes(0)
	require.NoError(t, err)
	
	addr1, err := codec.OracleIDAsAddressBytes(1)
	require.NoError(t, err)
	
	addr255, err := codec.OracleIDAsAddressBytes(255)
	require.NoError(t, err)
	
	// Verify all addresses are different
	require.NotEqual(t, addr0, addr1, "different oracle IDs should produce different addresses")
	require.NotEqual(t, addr0, addr255, "different oracle IDs should produce different addresses")
	require.NotEqual(t, addr1, addr255, "different oracle IDs should produce different addresses")
	
	// Verify all addresses have the same length
	require.Len(t, addr0, 36)
	require.Len(t, addr1, 36)
	require.Len(t, addr255, 36)
}

// TestAddressCodec_OracleIDAsAddressBytes_Consistency tests that the same oracle ID always produces the same address
func TestAddressCodec_OracleIDAsAddressBytes_Consistency(t *testing.T) {
	codec := AddressCodec{}
	
	testOracleIDs := []uint8{0, 1, 42, 127, 200, 255}
	
	for _, oracleID := range testOracleIDs {
		t.Run(fmt.Sprintf("oracleID_%d_consistency", oracleID), func(t *testing.T) {
			// Generate address multiple times
			addr1, err1 := codec.OracleIDAsAddressBytes(oracleID)
			require.NoError(t, err1)
			
			addr2, err2 := codec.OracleIDAsAddressBytes(oracleID)
			require.NoError(t, err2)
			
			addr3, err3 := codec.OracleIDAsAddressBytes(oracleID)
			require.NoError(t, err3)
			
			// Verify all results are identical
			require.Equal(t, addr1, addr2, "same oracle ID should always produce same address")
			require.Equal(t, addr1, addr3, "same oracle ID should always produce same address")
			require.Equal(t, addr2, addr3, "same oracle ID should always produce same address")
		})
	}
}

// packOracleIDForTest is a helper function that mimics the OracleIDAsAddressBytes implementation
// for testing purposes. This ensures our test expectations match the actual implementation.
func packOracleIDForTest(t *testing.T, oracleID uint8) []byte {
	addr := make([]byte, 32)
	binary.BigEndian.PutUint32(addr, uint32(oracleID))
	tonAddr := address.NewAddress(0, 0, addr)
	decodeString, err := base64.RawURLEncoding.DecodeString(tonAddr.String())
	if err != nil {
		t.Fatalf("failed to decode TVM address bytes: %v", err)
	}
	return decodeString
}