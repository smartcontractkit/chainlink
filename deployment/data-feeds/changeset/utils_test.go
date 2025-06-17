package changeset

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDecimalsFromFeedIDValidFeedID(t *testing.T) {
	decimals, err := GetDecimalsFromFeedID("0x01d0fd1ef80003320000000000000000") // decimal 18
	require.NoError(t, err)
	require.Equal(t, uint8(0x12), decimals)

	decimals, err = GetDecimalsFromFeedID("0x01d0fd1ef80003600000000000000000") // decimal 64
	require.NoError(t, err)
	require.Equal(t, uint8(0x40), decimals)

	_, err = GetDecimalsFromFeedID("0x01d0fd1ef80003ff0000000000000000") // invalid bytes for decimal
	require.Error(t, err)

	_, err = GetDecimalsFromFeedID("0x00") // invalid feed id
	require.Error(t, err)
}
