package operations

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChainCapJobNameVariants(t *testing.T) {
	t.Parallel()

	t.Run("evm current and legacy", func(t *testing.T) {
		t.Parallel()
		current, legacy, ok := chainCapJobNameVariants("evm-cap-v2-ethereum-testnet-sepolia-zone-a")
		require.True(t, ok)
		assert.Equal(t, "evm-cap-v2-ethereum-testnet-sepolia-zone-a", current)
		assert.Equal(t, "evm-capabilities-v2-ethereum-testnet-sepolia-zone-a", legacy)
	})

	t.Run("solana has no legacy alias", func(t *testing.T) {
		t.Parallel()
		current, legacy, ok := chainCapJobNameVariants("solana-cap-v2-solana-devnet-zone-b")
		require.True(t, ok)
		assert.Equal(t, "solana-cap-v2-solana-devnet-zone-b", current)
		assert.Empty(t, legacy)
	})

	t.Run("unknown prefix", func(t *testing.T) {
		t.Parallel()
		_, _, ok := chainCapJobNameVariants("http-trigger-zone-a")
		assert.False(t, ok)
	})
}
