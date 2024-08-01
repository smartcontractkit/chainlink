package utils

import (
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/stretchr/testify/require"
)

func HexToFelt(t *testing.T, hex string) *felt.Felt {
	t.Helper()
	f, err := new(felt.Felt).SetString(hex)
	require.NoError(t, err)
	return f
}
