package evm_test

import (
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

// transferOwnershipToTimelock transfers all CCIP contract ownership to the
// MCMS Timelock using the standard test helper. After the transfer, the MCMS
// multisig contracts (ProposerMcm, CancellerMcm, BypasserMcm) are nil-ed out
// because they cannot be made self-governed in the test environment.
func transferOwnershipToTimelock(
	t *testing.T,
	tenv testhelpers.DeployedEnv,
	state stateview.CCIPOnChainState,
	selectors []uint64,
) {
	t.Helper()
	testhelpers.TransferToTimelock(t, tenv, state, selectors, false)
	// MCMS multisig contracts are deployed by the deployer key and cannot
	// easily be made self-governed in the memory test environment.
	// Nil them out so ValidateContractOwnership skips the self-governance checks.
	for _, sel := range selectors {
		cs := state.MustGetEVMChainState(sel)
		cs.ProposerMcm = nil
		cs.CancellerMcm = nil
		cs.BypasserMcm = nil
		state.WriteEVMChainState(sel, cs)
	}
}

// TestValidatePostDeploymentState_HappyPath uses a full memory environment
// to verify that ValidatePostDeploymentState passes on a correctly-wired deployment.
// Contract ownership is transferred to the MCMS Timelock before validation.
func TestValidatePostDeploymentState_HappyPath(t *testing.T) {
	t.Parallel()
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(3))
	state, err := stateview.LoadOnchainState(tenv.Env, stateview.WithLoadLegacyContracts(true))
	require.NoError(t, err)

	evmChains := tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	transferOwnershipToTimelock(t, tenv, state, evmChains)

	err = state.ValidatePostDeploymentState(tenv.Env, true)
	require.NoError(t, err, "expected no errors on a correctly-deployed environment")
}

// TestValidatePostDeploymentState_CollectsMultipleErrors verifies that the
// validation collects all errors rather than returning early on the first one.
func TestValidatePostDeploymentState_CollectsMultipleErrors(t *testing.T) {
	t.Parallel()
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(3))
	state, err := stateview.LoadOnchainState(tenv.Env, stateview.WithLoadLegacyContracts(true))
	require.NoError(t, err)

	evmChains := tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	transferOwnershipToTimelock(t, tenv, state, evmChains)
	require.GreaterOrEqual(t, len(evmChains), 2, "need at least 2 chains for this test")

	// Intentionally break multiple chains' state to force multiple errors:
	// Nil out the RMNProxy on one chain and the FeeQuoter on another.
	chainState0 := state.MustGetEVMChainState(evmChains[0])
	chainState0.RMNProxy = nil
	state.WriteEVMChainState(evmChains[0], chainState0)

	chainState1 := state.MustGetEVMChainState(evmChains[1])
	chainState1.FeeQuoter = nil
	state.WriteEVMChainState(evmChains[1], chainState1)

	err = state.ValidatePostDeploymentState(tenv.Env, false)
	require.Error(t, err, "expected validation errors")

	// The error should contain mentions of both chains' issues.
	errMsg := err.Error()
	assert.True(t, strings.Contains(errMsg, "RMNProxy") || strings.Contains(errMsg, "rmnProxy"),
		"expected error to mention RMNProxy issue, got: %s", errMsg)
	assert.True(t, strings.Contains(errMsg, "fee quoter") || strings.Contains(errMsg, "FeeQuoter"),
		"expected error to mention FeeQuoter issue, got: %s", errMsg)
}

// TestValidateContractOwnership_DetectsWrongOwner verifies that ownership
// validation detects contracts owned by the deployer rather than the timelock.
// In the memory environment, MCMS is deployed but ownership is NOT transferred.
func TestValidateContractOwnership_DetectsWrongOwner(t *testing.T) {
	t.Parallel()
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(2))
	state, err := stateview.LoadOnchainState(tenv.Env)
	require.NoError(t, err)

	evmChains := tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	chainState := state.MustGetEVMChainState(evmChains[0])
	require.NotNil(t, chainState.Timelock, "test expects Timelock to be deployed")

	err = chainState.ValidateContractOwnership(tenv.Env)
	require.Error(t, err, "expected ownership errors since contracts are not owned by timelock")
	assert.Contains(t, err.Error(), "not owned by expected owner")
}

// TestValidateContractOwnership_NoTimelock returns early with an error
// when timelock is nil.
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

// TestValidateRMNProxy_HappyPath validates that the RMNProxy correctly
// points to RMNRemote on a fresh deployment.
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

// TestValidateRMNProxy_MissingContracts returns errors when RMNProxy or RMNRemote is nil.
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

// TestValidateNonceManager_HappyPath validates the NonceManager on a full deployment.
func TestValidateNonceManager_HappyPath(t *testing.T) {
	t.Parallel()
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(3))
	state, err := stateview.LoadOnchainState(tenv.Env, stateview.WithLoadLegacyContracts(true))
	require.NoError(t, err)

	evmChains := tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	for _, sel := range evmChains {
		chainState := state.MustGetEVMChainState(sel)
		// Build connected chains from router
		connectedChains, err := chainState.ValidateRouter(tenv.Env, false)
		require.NoError(t, err, "router validation failed for chain %d", sel)

		err = chainState.ValidateNonceManager(tenv.Env, sel, connectedChains)
		require.NoError(t, err, "NonceManager validation failed for chain %d", sel)
	}
}

// TestValidateNonceManager_NilNonceManager returns error when NonceManager is nil.
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

// TestValidateFeeQuoter_HappyPath validates FeeQuoter chain-level and lane-level
// configurations pass on a correctly-deployed environment.
func TestValidateFeeQuoter_HappyPath(t *testing.T) {
	t.Parallel()
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(3))
	state, err := stateview.LoadOnchainState(tenv.Env, stateview.WithLoadLegacyContracts(true))
	require.NoError(t, err)

	evmChains := tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	for _, sel := range evmChains {
		chainState := state.MustGetEVMChainState(sel)
		connectedChains, err := chainState.ValidateRouter(tenv.Env, false)
		require.NoError(t, err, "router validation failed for chain %d", sel)

		err = chainState.ValidateFeeQuoter(tenv.Env, sel, connectedChains)
		require.NoError(t, err, "FeeQuoter validation failed for chain %d", sel)
	}
}

// TestValidateFeeQuoter_NilFeeQuoter returns error when FeeQuoter is nil.
func TestValidateFeeQuoter_NilFeeQuoter(t *testing.T) {
	t.Parallel()
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(2))
	state, err := stateview.LoadOnchainState(tenv.Env)
	require.NoError(t, err)

	evmChains := tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	chainState := state.MustGetEVMChainState(evmChains[0])
	chainState.FeeQuoter = nil
	err = chainState.ValidateFeeQuoter(tenv.Env, evmChains[0], evmChains[1:])
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no FeeQuoter")
}

// buildHomeChainTestArgs builds the nodes and offRampsByChain arguments needed to call ValidateHomeChain.
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

// TestValidateHomeChain_HappyPath validates home chain + per-chain DON config on a full deployment.
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

// TestValidateHomeChain_MissingContracts returns errors when CCIPHome or CapReg is nil.
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
