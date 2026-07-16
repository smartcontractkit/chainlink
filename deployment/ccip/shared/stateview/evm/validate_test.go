package evm_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/offramp"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

func buildV16ActiveChains(
	t *testing.T,
	tenv testhelpers.DeployedEnv,
	state stateview.CCIPOnChainState,
) map[uint64]bool {
	t.Helper()
	homeChainState := state.MustGetEVMChainState(tenv.HomeChainSel)
	v16Active, err := homeChainState.V16ActiveChainSelectors(tenv.Env.GetContext())
	require.NoError(t, err)
	return v16Active
}

// transferOwnershipToTimelock transfers ownership and nils out MCMS multisig
// contracts that can't be self-governed in test environments.
func transferOwnershipToTimelock(
	t *testing.T,
	tenv testhelpers.DeployedEnv,
	state stateview.CCIPOnChainState,
	selectors []uint64,
) {
	t.Helper()
	testhelpers.TransferToTimelock(t, tenv, state, selectors, false)
	for _, sel := range selectors {
		cs := state.MustGetEVMChainState(sel)
		cs.ProposerMcm = nil
		cs.CancellerMcm = nil
		cs.BypasserMcm = nil
		state.WriteEVMChainState(sel, cs)
	}
}

func TestValidatePostDeploymentState_HappyPath(t *testing.T) {
	t.Parallel()
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(3))
	state, err := stateview.LoadOnchainState(tenv.Env, stateview.WithLoadLegacyContracts(true))
	require.NoError(t, err)

	evmChains := tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	transferOwnershipToTimelock(t, tenv, state, evmChains)

	chainErrs := state.ValidatePostDeploymentState(tenv.Env, true, nil)
	require.Empty(t, chainErrs, "expected no errors on a correctly-deployed environment")
}

func TestValidatePostDeploymentState_CollectsMultipleErrors(t *testing.T) {
	t.Parallel()
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(3))
	state, err := stateview.LoadOnchainState(tenv.Env, stateview.WithLoadLegacyContracts(true))
	require.NoError(t, err)

	evmChains := tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	transferOwnershipToTimelock(t, tenv, state, evmChains)
	require.GreaterOrEqual(t, len(evmChains), 2, "need at least 2 chains for this test")

	chainState0 := state.MustGetEVMChainState(evmChains[0])
	chainState0.RMNProxy = nil
	state.WriteEVMChainState(evmChains[0], chainState0)

	chainState1 := state.MustGetEVMChainState(evmChains[1])
	chainState1.FeeQuoter = nil
	state.WriteEVMChainState(evmChains[1], chainState1)

	chainErrs := state.ValidatePostDeploymentState(tenv.Env, false, nil)
	require.NotEmpty(t, chainErrs, "expected validation errors")

	var allErrs []string
	for _, errs := range chainErrs {
		for _, e := range errs {
			allErrs = append(allErrs, e.Error())
		}
	}
	errMsg := strings.Join(allErrs, "; ")
	assert.True(t, strings.Contains(errMsg, "RMNProxy") || strings.Contains(errMsg, "rmnProxy"),
		"expected error to mention RMNProxy issue, got: %s", errMsg)
	assert.True(t, strings.Contains(errMsg, "fee quoter") || strings.Contains(errMsg, "FeeQuoter"),
		"expected error to mention FeeQuoter issue, got: %s", errMsg)
}

func TestValidateContractOwnership_DetectsWrongOwner(t *testing.T) {
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(2))
	state, err := stateview.LoadOnchainState(tenv.Env)
	require.NoError(t, err)

	evmChains := tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	slices.Sort(evmChains)
	chainState := state.MustGetEVMChainState(evmChains[0])
	require.NotNil(t, chainState.Timelock, "test expects Timelock to be deployed")

	err = chainState.ValidateContractOwnership(tenv.Env)
	require.Error(t, err, "expected ownership errors since contracts are not owned by timelock")
	assert.Contains(t, err.Error(), "not owned by expected owner")
}

func TestValidateContractOwnership_NoTimelock(t *testing.T) {
	t.Parallel()
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(2))
	state, err := stateview.LoadOnchainState(tenv.Env)
	require.NoError(t, err)

	evmChains := tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	chainState := state.MustGetEVMChainState(evmChains[0])
	chainState.Timelock = nil
	err = chainState.ValidateContractOwnership(tenv.Env)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "timelock not found")
}

func TestValidateRMNProxy_HappyPath(t *testing.T) {
	t.Parallel()
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(2))
	state, err := stateview.LoadOnchainState(tenv.Env)
	require.NoError(t, err)

	for _, sel := range tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM)) {
		chainState := state.MustGetEVMChainState(sel)
		err := chainState.ValidateRMNProxy(tenv.Env)
		require.NoError(t, err, "RMNProxy validation failed for chain %d", sel)
	}
}

func TestValidateRMNProxy_MissingContracts(t *testing.T) {
	t.Parallel()
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(2))
	state, err := stateview.LoadOnchainState(tenv.Env)
	require.NoError(t, err)

	evmChains := tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))

	t.Run("nil RMNProxy", func(t *testing.T) {
		chainState := state.MustGetEVMChainState(evmChains[0])
		chainState.RMNProxy = nil
		err := chainState.ValidateRMNProxy(tenv.Env)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no RMNProxy")
	})

	t.Run("nil RMNRemote", func(t *testing.T) {
		chainState := state.MustGetEVMChainState(evmChains[0])
		chainState.RMNRemote = nil
		err := chainState.ValidateRMNProxy(tenv.Env)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no RMNRemote")
	})
}

func TestValidateNonceManager_HappyPath(t *testing.T) {
	t.Parallel()
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(3))
	state, err := stateview.LoadOnchainState(tenv.Env, stateview.WithLoadLegacyContracts(true))
	require.NoError(t, err)

	evmChains := tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	for _, sel := range evmChains {
		chainState := state.MustGetEVMChainState(sel)
		v16Active := buildV16ActiveChains(t, tenv, state)
		connectedChains, err := chainState.ValidateRouter(tenv.Env, false, v16Active)
		require.NoError(t, err, "router validation failed for chain %d", sel)

		err = chainState.ValidateNonceManager(tenv.Env, sel, connectedChains)
		require.NoError(t, err, "NonceManager validation failed for chain %d", sel)
	}
}

func TestValidateNonceManager_NilNonceManager(t *testing.T) {
	t.Parallel()
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(2))
	state, err := stateview.LoadOnchainState(tenv.Env)
	require.NoError(t, err)

	evmChains := tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	chainState := state.MustGetEVMChainState(evmChains[0])
	chainState.NonceManager = nil
	err = chainState.ValidateNonceManager(tenv.Env, evmChains[0], evmChains[1:])
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no NonceManager")
}

func buildHomeChainTestArgs(
	t *testing.T,
	tenv testhelpers.DeployedEnv,
	state stateview.CCIPOnChainState,
) (deployment.Nodes, map[uint64]offramp.OffRampInterface) {
	t.Helper()
	nodes, err := deployment.NodeInfo(tenv.Env.NodeIDs, tenv.Env.Offchain)
	require.NoError(t, err)
	offRamps := make(map[uint64]offramp.OffRampInterface)
	for _, sel := range state.EVMChains() {
		cs := state.MustGetEVMChainState(sel)
		offRamps[sel] = cs.OffRamp
	}
	return nodes, offRamps
}

func TestValidateHomeChain_HappyPath(t *testing.T) {
	t.Parallel()
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(3))
	state, err := stateview.LoadOnchainState(tenv.Env)
	require.NoError(t, err)

	homeChainState := state.MustGetEVMChainState(tenv.HomeChainSel)
	nodes, offRamps := buildHomeChainTestArgs(t, tenv, state)
	err = homeChainState.ValidateHomeChain(tenv.Env, nodes, offRamps)
	require.NoError(t, err, "home chain validation failed")
}

func TestValidateHomeChain_MissingContracts(t *testing.T) {
	t.Parallel()
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(2))
	state, err := stateview.LoadOnchainState(tenv.Env)
	require.NoError(t, err)

	nodes, offRamps := buildHomeChainTestArgs(t, tenv, state)

	t.Run("nil CCIPHome", func(t *testing.T) {
		homeChainState := state.MustGetEVMChainState(tenv.HomeChainSel)
		homeChainState.CCIPHome = nil
		err := homeChainState.ValidateHomeChain(tenv.Env, nodes, offRamps)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no CCIPHome")
	})

	t.Run("nil CapabilityRegistry", func(t *testing.T) {
		homeChainState := state.MustGetEVMChainState(tenv.HomeChainSel)
		homeChainState.CapabilityRegistry = nil
		err := homeChainState.ValidateHomeChain(tenv.Env, nodes, offRamps)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no CapabilityRegistry")
	})
}
