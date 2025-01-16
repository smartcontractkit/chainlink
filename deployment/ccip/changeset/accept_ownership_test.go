package changeset_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

func Test_NewAcceptOwnershipChangeset(t *testing.T) {
	t.Parallel()
	e, _ := testhelpers.NewMemoryEnvironment(t)
	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)

	allChains := maps.Keys(e.Env.Chains)
	source := allChains[0]
	dest := allChains[1]

	timelockContracts := map[uint64]*proposalutils.TimelockExecutionContracts{
		source: {
			Timelock:  state.Chains[source].Timelock,
			CallProxy: state.Chains[source].CallProxy,
		},
		dest: {
			Timelock:  state.Chains[dest].Timelock,
			CallProxy: state.Chains[dest].CallProxy,
		},
	}

	// at this point we have the initial deploys done, now we need to transfer ownership
	// to the timelock contract
	state, err = changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)

	// compose the transfer ownership and accept ownership changesets
	_, err = commonchangeset.ApplyChangesets(t, e.Env, timelockContracts, []commonchangeset.ChangesetApplication{
		// note this doesn't have proposals.
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.TransferToMCMSWithTimelock),
			Config:    testhelpers.GenTestTransferOwnershipConfig(e, allChains, state),
		},
	})
	require.NoError(t, err)

	testhelpers.AssertTimelockOwnership(t, e, allChains, state)
}
