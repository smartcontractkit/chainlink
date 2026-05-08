package state_test

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

// TestMaybeLoadMCMSWithTimelockChainState_LabeledManyChainMultisig_AddressMap verifies ManyChainMultiSig @ v1.0.0
// with standard MCMS role labels maps to ManyChainMultiSig bindings the same way as typed Proposer/Bypasser/Canceller MCMS rows.
func TestMaybeLoadMCMSWithTimelockChainState_LabeledManyChainMultisig_AddressMap(t *testing.T) {
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

	labeledAddrs := map[string]cldf.TypeAndVersion{
		canonical.Timelock.Address().Hex():     cldf.NewTypeAndVersion(commontypes.RBACTimelock, deployment.Version1_0_0),
		canonical.CallProxy.Address().Hex():    cldf.NewTypeAndVersion(commontypes.CallProxy, deployment.Version1_0_0),
		canonical.ProposerMcm.Address().Hex():  manyChainMultisigWithRole(t, commontypes.ProposerRole),
		canonical.BypasserMcm.Address().Hex():  manyChainMultisigWithRole(t, commontypes.BypasserRole),
		canonical.CancellerMcm.Address().Hex(): manyChainMultisigWithRole(t, commontypes.CancellerRole),
	}

	loaded, err := state.MaybeLoadMCMSWithTimelockChainState(chain, labeledAddrs)
	require.NoError(t, err)
	require.NoError(t, loaded.Validate())
	require.Equal(t, canonical.ProposerMcm.Address(), loaded.ProposerMcm.Address())
	require.Equal(t, canonical.BypasserMcm.Address(), loaded.BypasserMcm.Address())
	require.Equal(t, canonical.CancellerMcm.Address(), loaded.CancellerMcm.Address())
	require.Equal(t, canonical.Timelock.Address(), loaded.Timelock.Address())
	require.Equal(t, canonical.CallProxy.Address(), loaded.CallProxy.Address())
}

// TestMaybeLoadMCMSWithTimelockChainStateFromRefs_LabeledManyChainMultisig exercises DataStore-style refs with ManyChainMultiSig @ v1.0.0 and role labels.
func TestMaybeLoadMCMSWithTimelockChainStateFromRefs_LabeledManyChainMultisig(t *testing.T) {
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
			Type: datastore.ContractType(commontypes.ManyChainMultisig), Version: v10,
			Labels: datastore.NewLabelSet(commontypes.ProposerRole.String()),
		},
		{
			Address: canonical.BypasserMcm.Address().Hex(), ChainSelector: sel,
			Type: datastore.ContractType(commontypes.ManyChainMultisig), Version: v10,
			Labels: datastore.NewLabelSet(commontypes.BypasserRole.String()),
		},
		{
			Address: canonical.CancellerMcm.Address().Hex(), ChainSelector: sel,
			Type: datastore.ContractType(commontypes.ManyChainMultisig), Version: v10,
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

func manyChainMultisigWithRole(t *testing.T, role commontypes.MCMSRole) cldf.TypeAndVersion {
	t.Helper()
	tv := cldf.NewTypeAndVersion(commontypes.ManyChainMultisig, deployment.Version1_0_0)
	tv.Labels.Add(role.String())
	return tv
}
