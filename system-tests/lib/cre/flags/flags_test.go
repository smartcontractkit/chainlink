package flags

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

func TestRequiresForwarderContract(t *testing.T) {
	t.Parallel()

	t.Run("returns true for aptos write capability", func(t *testing.T) {
		t.Parallel()
		require.True(t, RequiresForwarderContract([]string{cre.AptosCapability + "-4"}, 4))
	})

	t.Run("returns true for evm, solana, and stellar write paths", func(t *testing.T) {
		t.Parallel()
		require.True(t, RequiresForwarderContract([]string{cre.EVMCapability + "-1337"}, 1337))
		require.True(t, RequiresForwarderContract([]string{cre.SolanaCapability + "-1"}, 1))
		require.True(t, RequiresForwarderContract([]string{cre.StellarCapability}, 0))
	})
}

func TestHasFlagForAnyChain(t *testing.T) {
	t.Parallel()

	values := []string{cre.ConsensusQueueCapability, cre.EVMCapability + "-1337", cre.CronCapability}

	tests := []struct {
		name       string
		capability string
		want       bool
	}{
		{"exact match", cre.CronCapability, true},
		{"exact match with a hyphen in the name", cre.ConsensusQueueCapability, true},
		{"chain ID suffix", cre.EVMCapability, true},
		{"non-numeric suffix is a different flag", cre.ConsensusCapability, false},
		{"absent", cre.VaultCapability, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, HasFlagForAnyChain(values, tt.capability))
			require.Equal(t, tt.want, cre.HasFlag(values, tt.capability))
		})
	}
}
