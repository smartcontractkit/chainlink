package changeset_test

import (
	"testing"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commonState "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
)

func TestDeployLinkToken(t *testing.T) {
	t.Parallel()

	selector := chain_selectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
	))
	require.NoError(t, err)

	chain := rt.Environment().BlockChains.EVMChains()[selector]

	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(changeset.DeployLinkToken), []uint64{selector}),
	)
	require.NoError(t, err)

	addrs, err := rt.State().AddressBook.AddressesForChain(selector)
	require.NoError(t, err)

	state, err := commonState.MaybeLoadLinkTokenChainState(chain, addrs)
	require.NoError(t, err)

	// View itself already unit tested
	_, err = state.GenerateLinkView()
	require.NoError(t, err)
}

func TestDeployLinkTokenZk(t *testing.T) {
	// Timeouts in CI
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/CCIP-6427")

	t.Parallel()

	selector := chain_selectors.TEST_90000050.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithZKSyncContainer(t, []uint64{selector}),
	))
	require.NoError(t, err)

	chain := rt.Environment().BlockChains.EVMChains()[selector]

	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(changeset.DeployLinkToken), []uint64{selector}),
	)
	require.NoError(t, err)

	addrs, err := rt.State().AddressBook.AddressesForChain(selector)
	require.NoError(t, err)

	state, err := commonState.MaybeLoadLinkTokenChainState(chain, addrs)
	require.NoError(t, err)

	// View itself already unit tested
	_, err = state.GenerateLinkView()
	require.NoError(t, err)
}
