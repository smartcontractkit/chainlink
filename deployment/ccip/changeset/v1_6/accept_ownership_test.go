package v1_6_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

func Test_NewAcceptOwnershipChangeset(t *testing.T) {
	t.Parallel()
	e, _ := testhelpers.NewMemoryEnvironment(t)
	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	allChains := maps.Keys(e.Env.BlockChains.EVMChains())

	// at this point we have the initial deploys done, now we need to transfer ownership
	// to the timelock contract
	state, err = stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	// compose the transfer ownership and accept ownership changesets
	_, err = commonchangeset.Apply(t, e.Env,
		// note this doesn't have proposals.
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(commonchangeset.TransferToMCMSWithTimelockV2),
			testhelpers.GenTestTransferOwnershipConfig(e, allChains, state, true),
		),
	)
	require.NoError(t, err)

	testhelpers.AssertTimelockOwnership(t, e, allChains, state, true)
}
