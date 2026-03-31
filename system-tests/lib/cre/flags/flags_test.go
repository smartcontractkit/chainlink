package flags

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

func TestRequiresForwarderContract(t *testing.T) {
	t.Run("returns true for aptos write capability", func(t *testing.T) {
		require.True(t, RequiresForwarderContract([]string{cre.AptosCapability + "-4"}, 4))
	})

	t.Run("returns true for evm and solana write paths", func(t *testing.T) {
		require.True(t, RequiresForwarderContract([]string{cre.EVMCapability + "-1337"}, 1337))
		require.True(t, RequiresForwarderContract([]string{cre.WriteEVMCapability + "-1337"}, 1337))
		require.True(t, RequiresForwarderContract([]string{cre.SolanaCapability + "-1"}, 1))
	})

	t.Run("returns false for read-only aptos capability set", func(t *testing.T) {
		require.False(t, RequiresForwarderContract([]string{cre.ReadContractCapability + "-4"}, 4))
	})
}
