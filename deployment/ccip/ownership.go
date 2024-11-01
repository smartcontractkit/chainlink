package ccipdeployment

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"
)

func TransferAllOwnership(t *testing.T, state CCIPOnChainState, homeCS uint64, e deployment.Environment) {
	for _, source := range e.AllChainSelectors() {
		if state.Chains[source].OnRamp != nil {
			tx, err := state.Chains[source].OnRamp.TransferOwnership(e.Chains[source].DeployerKey, state.Chains[source].Timelock.Address())
			require.NoError(t, err)
			_, err = deployment.ConfirmIfNoError(e.Chains[source], tx, err)
			require.NoError(t, err)
		}
		if state.Chains[source].FeeQuoter != nil {
			tx, err := state.Chains[source].FeeQuoter.TransferOwnership(e.Chains[source].DeployerKey, state.Chains[source].Timelock.Address())
			require.NoError(t, err)
			_, err = deployment.ConfirmIfNoError(e.Chains[source], tx, err)
			require.NoError(t, err)
		}
		// TODO: add offramp and commit stores

	}
	// Transfer CR contract ownership
	tx, err := state.Chains[homeCS].CapabilityRegistry.TransferOwnership(e.Chains[homeCS].DeployerKey, state.Chains[homeCS].Timelock.Address())
	require.NoError(t, err)
	_, err = deployment.ConfirmIfNoError(e.Chains[homeCS], tx, err)
	require.NoError(t, err)
	tx, err = state.Chains[homeCS].CCIPHome.TransferOwnership(e.Chains[homeCS].DeployerKey, state.Chains[homeCS].Timelock.Address())
	require.NoError(t, err)
	_, err = deployment.ConfirmIfNoError(e.Chains[homeCS], tx, err)
	require.NoError(t, err)
}
