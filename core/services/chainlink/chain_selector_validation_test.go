package chainlink

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func Test_validateChainSelectorsForConfiguredChains(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)

	t.Run("succeeds with no relayers configured", func(t *testing.T) {
		t.Parallel()

		relayChainInterops := &CoreRelayerChainInteroperators{
			loopRelayers: make(map[commontypes.RelayID]loop.Relayer),
		}

		err := validateChainSelectorsForConfiguredChains(lggr, relayChainInterops)
		require.NoError(t, err)
	})

	t.Run("succeeds with valid chains", func(t *testing.T) {
		t.Parallel()

		relayChainInterops := &CoreRelayerChainInteroperators{
			loopRelayers: map[commontypes.RelayID]loop.Relayer{
				// EVM chains
				{Network: "evm", ChainID: "1"}:        nil, // Ethereum Mainnet
				{Network: "evm", ChainID: "11155111"}: nil, // Sepolia
				{Network: "evm", ChainID: "42161"}:    nil, // Arbitrum One
				// Solana chains
				{Network: "solana", ChainID: "5eykt4UsFv8P8NJdTREpY1vzqKqZKvdpKuc147dw2N9d"}: nil, // Mainnet
				{Network: "solana", ChainID: "EtWTRABZaYq6iMfeYKouRu166VU2xqa1wcaWoxPkrZBG"}: nil, // Devnet
			},
		}

		err := validateChainSelectorsForConfiguredChains(lggr, relayChainInterops)
		require.NoError(t, err)
	})

	t.Run("skips dummy relayers", func(t *testing.T) {
		t.Parallel()

		relayChainInterops := &CoreRelayerChainInteroperators{
			loopRelayers: map[commontypes.RelayID]loop.Relayer{
				// Valid chain
				{Network: "evm", ChainID: "1"}: nil,
				// Dummy relayers should be skipped
				{Network: "dummy", ChainID: "test-chain-1"}: nil,
				{Network: "dummy", ChainID: "test-chain-2"}: nil,
			},
		}

		err := validateChainSelectorsForConfiguredChains(lggr, relayChainInterops)
		require.NoError(t, err)
	})

	t.Run("fails with invalid EVM chain ID", func(t *testing.T) {
		t.Parallel()

		relayChainInterops := &CoreRelayerChainInteroperators{
			loopRelayers: map[commontypes.RelayID]loop.Relayer{
				{Network: "evm", ChainID: "999999999999999"}: nil,
			},
		}

		err := validateChainSelectorsForConfiguredChains(lggr, relayChainInterops)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid chain configuration")
		assert.Contains(t, err.Error(), "999999999999999")
		assert.Contains(t, err.Error(), "not recognized")
	})

	t.Run("fails with invalid Solana chain ID", func(t *testing.T) {
		t.Parallel()

		relayChainInterops := &CoreRelayerChainInteroperators{
			loopRelayers: map[commontypes.RelayID]loop.Relayer{
				{Network: "solana", ChainID: "invalid-solana-genesis-hash"}: nil,
			},
		}

		err := validateChainSelectorsForConfiguredChains(lggr, relayChainInterops)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid chain configuration")
		assert.Contains(t, err.Error(), "solana")
		assert.Contains(t, err.Error(), "invalid-solana-genesis-hash")
	})

	t.Run("fails with typo in chain ID", func(t *testing.T) {
		t.Parallel()

		relayChainInterops := &CoreRelayerChainInteroperators{
			loopRelayers: map[commontypes.RelayID]loop.Relayer{
				{Network: "evm", ChainID: "1115511"}: nil, //Sepolia mistyped
			},
		}

		err := validateChainSelectorsForConfiguredChains(lggr, relayChainInterops)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "1115511")
		assert.Contains(t, err.Error(), "ChainID in TOML is correct")
	})
}
