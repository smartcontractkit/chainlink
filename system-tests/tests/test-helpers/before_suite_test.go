package helpers

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"

	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
)

func TestLocalEnvironmentSatisfiesRequestedConfig(t *testing.T) {
	t.Run("returns true when saved state contains requested chain families and ids", func(t *testing.T) {
		ok, err := savedEnvironmentSatisfiesRequestedConfig(
			&envconfig.Config{
				Blockchains: []*blockchain.Input{
					{Type: blockchain.TypeAnvil, ChainID: "1337"},
					{Type: blockchain.TypeAptos, ChainID: "4"},
				},
			},
			&envconfig.Config{
				Blockchains: []*blockchain.Input{
					{Type: blockchain.TypeAnvil, ChainID: "1337", Out: &blockchain.Output{Family: blockchain.FamilyEVM, ChainID: "1337"}},
					{Type: blockchain.TypeAptos, ChainID: "4", Out: &blockchain.Output{Family: blockchain.FamilyAptos, ChainID: "4"}},
				},
			},
		)
		require.NoError(t, err)
		require.True(t, ok)
	})

	t.Run("returns false when saved state is missing a requested chain", func(t *testing.T) {
		ok, err := savedEnvironmentSatisfiesRequestedConfig(
			&envconfig.Config{
				Blockchains: []*blockchain.Input{
					{Type: blockchain.TypeAnvil, ChainID: "1337"},
					{Type: blockchain.TypeAptos, ChainID: "4"},
				},
			},
			&envconfig.Config{
				Blockchains: []*blockchain.Input{
					{Type: blockchain.TypeAnvil, ChainID: "1337", Out: &blockchain.Output{Family: blockchain.FamilyEVM, ChainID: "1337"}},
					{Type: blockchain.TypeAnvil, ChainID: "2337", Out: &blockchain.Output{Family: blockchain.FamilyEVM, ChainID: "2337"}},
				},
			},
		)
		require.NoError(t, err)
		require.False(t, ok)
	})
}

func TestChainKey(t *testing.T) {
	t.Run("derives family from type when output is absent", func(t *testing.T) {
		family, chainID, err := chainKey(&blockchain.Input{Type: blockchain.TypeAptos, ChainID: "4"})
		require.NoError(t, err)
		require.Equal(t, blockchain.FamilyAptos, family)
		require.Equal(t, "4", chainID)
	})

	t.Run("uses output family and chain id when available", func(t *testing.T) {
		family, chainID, err := chainKey(&blockchain.Input{
			Type:    blockchain.TypeAnvil,
			ChainID: "",
			Out: &blockchain.Output{
				Family:  blockchain.FamilyEVM,
				ChainID: "1337",
			},
		})
		require.NoError(t, err)
		require.Equal(t, blockchain.FamilyEVM, family)
		require.Equal(t, "1337", chainID)
	})
}
