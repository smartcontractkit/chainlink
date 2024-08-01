package utils

import (
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/test-go/testify/require"
)

// TestHexToFelt generates a felt.Felt from a hexadecimal string.
//
// Parameters:
// - t: the testing.TB object for test logging and reporting
// - hex: the hexadecimal string to convert to a felt.Felt
// Returns:
// - *felt.Felt: the generated felt.Felt object
func TestHexToFelt(t testing.TB, hex string) *felt.Felt {
	t.Helper()
	f, err := HexToFelt(hex)
	require.NoError(t, err)
	return f
}

// TestHexArrToFelt generates a slice of *felt.Felt from a slice of strings representing hexadecimal values.
//
// Parameters:
// - t: A testing.TB interface used for test logging and error reporting
// - hexArr: A slice of strings representing hexadecimal values
// Returns:
// - []*felt.Felt: a slice of *felt.Felt
func TestHexArrToFelt(t testing.TB, hexArr []string) []*felt.Felt {
	t.Helper()
	feltArr, err := HexArrToFelt(hexArr)
	require.NoError(t, err)
	return feltArr
}
