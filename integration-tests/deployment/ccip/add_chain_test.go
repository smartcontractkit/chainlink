package ccipdeployment

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestAddChain(t *testing.T) {
	// 4 chains where the 4th is added after initial deployment.
	e := NewEnvironmentWithCRAndJobs(t, logger.TestLogger(t), 4)
	state, err := LoadOnchainState(e.Env, e.Ab)
	require.NoError(t, err)
	// Take first non-home chain as the new chain.
	newChain := e.Env.AllChainSelectorsExcluding([]uint64{e.HomeChainSel})[0]
	// We deploy to the rest.
	initialDeploy := e.Env.AllChainSelectorsExcluding([]uint64{newChain})
	t.Logf("Home %d new %d initial %d\n", e.HomeChainSel, newChain, initialDeploy)

	ab, err := DeployCCIPContracts(e.Env, DeployCCIPContractConfig{
		HomeChainSel:     e.HomeChainSel,
		ChainsToDeploy:   initialDeploy,
		CCIPOnChainState: state,
	})
	require.NoError(t, err)
	require.NoError(t, e.Ab.Merge(ab))
	state, err = LoadOnchainState(e.Env, e.Ab)
	require.NoError(t, err)

	// Contracts deployed and initial DONs set up.
	// Connect all the lanes
	for _, source := range initialDeploy {
		for _, dest := range initialDeploy {
			if source != dest {
				require.NoError(t, AddLane(e.Env, state, source, dest))
			}
		}
	}

	//  Deploy contracts to new chain
	newAddresses, err := DeployChainContracts(e.Env, e.Env.Chains[newChain], deployment.NewMemoryAddressBook())
	require.NoError(t, err)
	require.NoError(t, e.Ab.Merge(newAddresses))
	state, err = LoadOnchainState(e.Env, e.Ab)
	require.NoError(t, err)

	// We can directly enable the sources on the new chain with deployer key.
	var offRampEnables []offramp.OffRampSourceChainConfigArgs
	for _, source := range initialDeploy {
		offRampEnables = append(offRampEnables, offramp.OffRampSourceChainConfigArgs{
			Router:              state.Chains[newChain].Router.Address(),
			SourceChainSelector: source,
			IsEnabled:           true,
			OnRamp:              common.LeftPadBytes(state.Chains[source].OnRamp.Address().Bytes(), 32),
		})
	}
	tx, err := state.Chains[newChain].OffRamp.ApplySourceChainConfigUpdates(e.Env.Chains[newChain].DeployerKey, offRampEnables)
	require.NoError(t, err)
	_, err = deployment.ConfirmIfNoError(e.Env.Chains[newChain], tx, err)
	require.NoError(t, err)

	// Transfer onramp/fq ownership to timelock.
	for _, source := range initialDeploy {
		tx, err := state.Chains[source].OnRamp.TransferOwnership(e.Env.Chains[source].DeployerKey, state.Chains[source].TimelockAddr)
		require.NoError(t, err)
		_, err = deployment.ConfirmIfNoError(e.Env.Chains[source], tx, err)
		require.NoError(t, err)
		tx, err = state.Chains[source].FeeQuoter.TransferOwnership(e.Env.Chains[source].DeployerKey, state.Chains[source].TimelockAddr)
		require.NoError(t, err)
		_, err = deployment.ConfirmIfNoError(e.Env.Chains[source], tx, err)
		require.NoError(t, err)
	}
	// Transfer CR contract ownership
	tx, err = state.Chains[e.HomeChainSel].CapabilityRegistry.TransferOwnership(e.Env.Chains[e.HomeChainSel].DeployerKey, state.Chains[e.HomeChainSel].TimelockAddr)
	require.NoError(t, err)
	_, err = deployment.ConfirmIfNoError(e.Env.Chains[e.HomeChainSel], tx, err)
	require.NoError(t, err)
	tx, err = state.Chains[e.HomeChainSel].CCIPConfig.TransferOwnership(e.Env.Chains[e.HomeChainSel].DeployerKey, state.Chains[e.HomeChainSel].TimelockAddr)
	require.NoError(t, err)
	_, err = deployment.ConfirmIfNoError(e.Env.Chains[e.HomeChainSel], tx, err)
	require.NoError(t, err)

	acceptOwnershipProposal, err := GenerateAcceptOwnershipProposal(state, e.HomeChainSel, initialDeploy)
	require.NoError(t, err)
	acceptOwnershipExec := SignProposal(t, e.Env, acceptOwnershipProposal)
	// Apply the accept ownership proposal to all the chains.
	for _, sel := range initialDeploy {
		ExecuteProposal(t, e.Env, acceptOwnershipExec, state, sel)
	}
	for _, chain := range initialDeploy {
		owner, err2 := state.Chains[chain].OnRamp.Owner(nil)
		require.NoError(t, err2)
		assert.Equal(t, state.Chains[chain].TimelockAddr, owner)
	}
	cfgOwner, err := state.Chains[e.HomeChainSel].CCIPConfig.Owner(nil)
	require.NoError(t, err)
	crOwner, err := state.Chains[e.HomeChainSel].CapabilityRegistry.Owner(nil)
	require.NoError(t, err)
	assert.Equal(t, state.Chains[e.HomeChainSel].TimelockAddr, cfgOwner)
	assert.Equal(t, state.Chains[e.HomeChainSel].TimelockAddr, crOwner)

	// Generate and sign inbound proposal to new 4th chain.
	chainInboundProposal, err := NewChainInboundProposal(e.Env, state, e.HomeChainSel, newChain, initialDeploy)
	require.NoError(t, err)
	chainInboundExec := SignProposal(t, e.Env, chainInboundProposal)
	for _, sel := range initialDeploy {
		ExecuteProposal(t, e.Env, chainInboundExec, state, sel)
	}

	state, err = LoadOnchainState(e.Env, e.Ab)
	require.NoError(t, err)
	for _, chain := range initialDeploy {
		cfg, err2 := state.Chains[chain].OnRamp.GetDestChainConfig(nil, newChain)
		require.NoError(t, err2)
		t.Log("config", cfg)
		s, err2 := state.Chains[newChain].OffRamp.GetSourceChainConfig(nil, chain)
		require.NoError(t, err2)
		t.Log("config", s)
	}

	// Now that the proposal has been executed we expect to be able to send traffic to this new 4th chain.
	//SendRequest(t, e.Env, state, initialDeploy[0], newChain)
	//e.Env.initialDeploy[0]

}
