package ccipdeployment

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestAddChainInbound(t *testing.T) {
	// 4 chains where the 4th is added after initial deployment.
	e := NewEnvironmentWithCRAndJobs(t, logger.TestLogger(t), 4)
	require.Equal(t, len(e.Nodes), 5)
	state, err := LoadOnchainState(e.Env, e.Ab)
	require.NoError(t, err)
	// Take first non-home chain as the new chain.
	newChain := e.Env.AllChainSelectorsExcluding([]uint64{e.HomeChainSel})[0]
	// We deploy to the rest.
	initialDeploy := e.Env.AllChainSelectorsExcluding([]uint64{newChain})

	ab, err := DeployCCIPContracts(e.Env, DeployCCIPContractConfig{
		HomeChainSel:     e.HomeChainSel,
		ChainsToDeploy:   initialDeploy,
		CCIPOnChainState: state,
	})
	require.NoError(t, err)
	require.NoError(t, e.Ab.Merge(ab))
	state, err = LoadOnchainState(e.Env, e.Ab)
	require.NoError(t, err)

	// Connect all the existing lanes.
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

	// Transfer onramp/fq ownership to timelock.
	// Enable the new dest on the test router.
	for _, source := range initialDeploy {
		tx, err := state.Chains[source].OnRamp.TransferOwnership(e.Env.Chains[source].DeployerKey, state.Chains[source].Timelock.Address())
		require.NoError(t, err)
		_, err = deployment.ConfirmIfNoError(e.Env.Chains[source], tx, err)
		require.NoError(t, err)
		tx, err = state.Chains[source].FeeQuoter.TransferOwnership(e.Env.Chains[source].DeployerKey, state.Chains[source].Timelock.Address())
		require.NoError(t, err)
		_, err = deployment.ConfirmIfNoError(e.Env.Chains[source], tx, err)
		require.NoError(t, err)
		tx, err = state.Chains[source].TestRouter.ApplyRampUpdates(e.Env.Chains[source].DeployerKey, []router.RouterOnRamp{
			{
				DestChainSelector: newChain,
				OnRamp:            state.Chains[source].OnRamp.Address(),
			},
		}, nil, nil)
		_, err = deployment.ConfirmIfNoError(e.Env.Chains[source], tx, err)
		require.NoError(t, err)
	}
	// Transfer CR contract ownership
	tx, err := state.Chains[e.HomeChainSel].CapabilityRegistry.TransferOwnership(e.Env.Chains[e.HomeChainSel].DeployerKey, state.Chains[e.HomeChainSel].Timelock.Address())
	require.NoError(t, err)
	_, err = deployment.ConfirmIfNoError(e.Env.Chains[e.HomeChainSel], tx, err)
	require.NoError(t, err)
	tx, err = state.Chains[e.HomeChainSel].CCIPConfig.TransferOwnership(e.Env.Chains[e.HomeChainSel].DeployerKey, state.Chains[e.HomeChainSel].Timelock.Address())
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
		require.Equal(t, state.Chains[chain].Timelock.Address(), owner)
	}
	cfgOwner, err := state.Chains[e.HomeChainSel].CCIPConfig.Owner(nil)
	require.NoError(t, err)
	crOwner, err := state.Chains[e.HomeChainSel].CapabilityRegistry.Owner(nil)
	require.NoError(t, err)
	require.Equal(t, state.Chains[e.HomeChainSel].Timelock.Address(), cfgOwner)
	require.Equal(t, state.Chains[e.HomeChainSel].Timelock.Address(), crOwner)

	// Generate and sign inbound proposal to new 4th chain.
	chainInboundProposal, err := NewChainInboundProposal(e.Env, state, e.HomeChainSel, newChain, initialDeploy)
	require.NoError(t, err)
	chainInboundExec := SignProposal(t, e.Env, chainInboundProposal)
	for _, sel := range initialDeploy {
		ExecuteProposal(t, e.Env, chainInboundExec, state, sel)
	}

	// Now configure the new chain using deployer key (not transferred to timelock yet).
	var offRampEnables []offramp.OffRampSourceChainConfigArgs
	for _, source := range initialDeploy {
		offRampEnables = append(offRampEnables, offramp.OffRampSourceChainConfigArgs{
			Router:              state.Chains[newChain].Router.Address(),
			SourceChainSelector: source,
			IsEnabled:           true,
			OnRamp:              common.LeftPadBytes(state.Chains[source].OnRamp.Address().Bytes(), 32),
		})
	}
	tx, err = state.Chains[newChain].OffRamp.ApplySourceChainConfigUpdates(e.Env.Chains[newChain].DeployerKey, offRampEnables)
	require.NoError(t, err)
	_, err = deployment.ConfirmIfNoError(e.Env.Chains[newChain], tx, err)
	require.NoError(t, err)
	// Set the OCR3 config on new 4th chain to enable the plugin.
	latestDON, err := LatestCCIPDON(state.Chains[e.HomeChainSel].CapabilityRegistry)
	require.NoError(t, err)
	ocrConfigs, err := BuildSetOCR3ConfigArgs(latestDON.Id, state.Chains[e.HomeChainSel].CCIPConfig)
	require.NoError(t, err)
	tx, err = state.Chains[newChain].OffRamp.SetOCR3Configs(e.Env.Chains[newChain].DeployerKey, ocrConfigs)
	require.NoError(t, err)
	_, err = deployment.ConfirmIfNoError(e.Env.Chains[newChain], tx, err)
	require.NoError(t, err)

	// Assert the inbound lanes to the new chain are wired correctly.
	state, err = LoadOnchainState(e.Env, e.Ab)
	require.NoError(t, err)
	for _, chain := range initialDeploy {
		cfg, err2 := state.Chains[chain].OnRamp.GetDestChainConfig(nil, newChain)
		require.NoError(t, err2)
		assert.Equal(t, cfg.Router, state.Chains[chain].TestRouter.Address())
		fqCfg, err2 := state.Chains[chain].FeeQuoter.GetDestChainConfig(nil, newChain)
		require.NoError(t, err2)
		assert.True(t, fqCfg.IsEnabled)
		s, err2 := state.Chains[newChain].OffRamp.GetSourceChainConfig(nil, chain)
		require.NoError(t, err2)
		assert.Equal(t, common.LeftPadBytes(state.Chains[chain].OnRamp.Address().Bytes(), 32), s.OnRamp)
	}
	// Ensure job related logs are up to date.
	time.Sleep(30 * time.Second)
	require.NoError(t, ReplayAllLogs(e.Nodes, e.Env.Chains))

	// TODO: Send via all inbound lanes and use parallel helper
	// Now that the proposal has been executed we expect to be able to send traffic to this new 4th chain.
	seqNr := SendRequest(t, e.Env, state, initialDeploy[0], newChain, true)
	ConfirmExecution(t, e.Env.Chains[initialDeploy[0]], e.Env.Chains[newChain], state.Chains[newChain].OffRamp, seqNr)
}
