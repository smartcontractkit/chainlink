package state_test

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

// TestMaybeLoadMCMSWithTimelockChainState_ProdTestnetMCM_AddressMap verifies legacy RDD datastore type
// prodTestnetMCM @ v1.0.0 maps to ManyChainMultiSig bindings when entries carry standard MCMS role labels.
func TestMaybeLoadMCMSWithTimelockChainState_ProdTestnetMCM_AddressMap(t *testing.T) {
	t.Parallel()

	sel := chain_selectors.TEST_90000001.Selector
	env, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{sel}),
		environment.WithLogger(logger.Test(t)),
	)
	require.NoError(t, err)

	cfg := proposalutils.SingleGroupTimelockConfigV2(t)
	updatedEnv, err := commonchangeset.Apply(t, *env, commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
		map[uint64]commontypes.MCMSWithTimelockConfigV2{
			sel: cfg,
		},
	))
	require.NoError(t, err)

	chain := updatedEnv.BlockChains.EVMChains()[sel]
	canonicalAddrs, err := updatedEnv.ExistingAddresses.AddressesForChain(sel)
	require.NoError(t, err)

	canonical, err := state.MaybeLoadMCMSWithTimelockChainState(chain, canonicalAddrs)
	require.NoError(t, err)
	require.NoError(t, canonical.Validate())

	prodAddrs := map[string]cldf.TypeAndVersion{
		canonical.Timelock.Address().Hex():    cldf.NewTypeAndVersion(commontypes.RBACTimelock, deployment.Version1_0_0),
		canonical.CallProxy.Address().Hex():   cldf.NewTypeAndVersion(commontypes.CallProxy, deployment.Version1_0_0),
		canonical.ProposerMcm.Address().Hex(): prodTestnetMCMWithRole(t, commontypes.ProposerRole),
		canonical.BypasserMcm.Address().Hex(): prodTestnetMCMWithRole(t, commontypes.BypasserRole),
		canonical.CancellerMcm.Address().Hex(): prodTestnetMCMWithRole(t, commontypes.CancellerRole),
	}

	loaded, err := state.MaybeLoadMCMSWithTimelockChainState(chain, prodAddrs)
	require.NoError(t, err)
	require.NoError(t, loaded.Validate())
	require.Equal(t, canonical.ProposerMcm.Address(), loaded.ProposerMcm.Address())
	require.Equal(t, canonical.BypasserMcm.Address(), loaded.BypasserMcm.Address())
	require.Equal(t, canonical.CancellerMcm.Address(), loaded.CancellerMcm.Address())
	require.Equal(t, canonical.Timelock.Address(), loaded.Timelock.Address())
	require.Equal(t, canonical.CallProxy.Address(), loaded.CallProxy.Address())
}

// TestMaybeLoadMCMSWithTimelockChainStateFromRefs_ProdTestnetMCM exercises DataStore-style refs with prodTestnetMCM.
func TestMaybeLoadMCMSWithTimelockChainStateFromRefs_ProdTestnetMCM(t *testing.T) {
	t.Parallel()

	sel := chain_selectors.TEST_90000001.Selector
	env, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{sel}),
		environment.WithLogger(logger.Test(t)),
	)
	require.NoError(t, err)

	cfg := proposalutils.SingleGroupTimelockConfigV2(t)
	updatedEnv, err := commonchangeset.Apply(t, *env, commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
		map[uint64]commontypes.MCMSWithTimelockConfigV2{
			sel: cfg,
		},
	))
	require.NoError(t, err)

	chain := updatedEnv.BlockChains.EVMChains()[sel]
	canonicalAddrs, err := updatedEnv.ExistingAddresses.AddressesForChain(sel)
	require.NoError(t, err)

	canonical, err := state.MaybeLoadMCMSWithTimelockChainState(chain, canonicalAddrs)
	require.NoError(t, err)
	require.NoError(t, canonical.Validate())

	v10 := semver.MustParse("1.0.0")
	refs := []datastore.AddressRef{
		{
			Address: canonical.Timelock.Address().Hex(), ChainSelector: sel,
			Type: datastore.ContractType(commontypes.RBACTimelock), Version: v10,
			Labels: datastore.NewLabelSet(),
		},
		{
			Address: canonical.CallProxy.Address().Hex(), ChainSelector: sel,
			Type: datastore.ContractType(commontypes.CallProxy), Version: v10,
			Labels: datastore.NewLabelSet(),
		},
		{
			Address: canonical.ProposerMcm.Address().Hex(), ChainSelector: sel,
			Type: datastore.ContractType(commontypes.ProdTestnetMCM), Version: v10,
			Labels: datastore.NewLabelSet(commontypes.ProposerRole.String()),
		},
		{
			Address: canonical.BypasserMcm.Address().Hex(), ChainSelector: sel,
			Type: datastore.ContractType(commontypes.ProdTestnetMCM), Version: v10,
			Labels: datastore.NewLabelSet(commontypes.BypasserRole.String()),
		},
		{
			Address: canonical.CancellerMcm.Address().Hex(), ChainSelector: sel,
			Type: datastore.ContractType(commontypes.ProdTestnetMCM), Version: v10,
			Labels: datastore.NewLabelSet(commontypes.CancellerRole.String()),
		},
	}

	loaded, err := state.MaybeLoadMCMSWithTimelockChainStateFromRefs(chain, refs)
	require.NoError(t, err)
	require.NoError(t, loaded.Validate())
	require.Equal(t, canonical.ProposerMcm.Address(), loaded.ProposerMcm.Address())
	require.Equal(t, canonical.BypasserMcm.Address(), loaded.BypasserMcm.Address())
	require.Equal(t, canonical.CancellerMcm.Address(), loaded.CancellerMcm.Address())
}

func prodTestnetMCMWithRole(t *testing.T, role commontypes.MCMSRole) cldf.TypeAndVersion {
	t.Helper()
	tv := cldf.NewTypeAndVersion(commontypes.ProdTestnetMCM, deployment.Version1_0_0)
	tv.Labels.Add(role.String())
	return tv
}
