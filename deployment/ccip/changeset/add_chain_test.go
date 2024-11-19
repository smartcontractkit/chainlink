package changeset

import (
	"testing"
	"time"

	ccipdeployment "github.com/smartcontractkit/chainlink/deployment/ccip"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_home"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	commonutils "github.com/smartcontractkit/chainlink-common/pkg/utils"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/deployment"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestAddChainInbound(t *testing.T) {
	// 4 chains where the 4th is added after initial deployment.
	e := ccipdeployment.NewMemoryEnvironmentWithJobs(t, logger.TestLogger(t), 4, 4)
	state, err := ccipdeployment.LoadOnchainState(e.Env)
	require.NoError(t, err)
	// Take first non-home chain as the new chain.
	newChain := e.Env.AllChainSelectorsExcluding([]uint64{e.HomeChainSel})[0]
	// We deploy to the rest.
	initialDeploy := e.Env.AllChainSelectorsExcluding([]uint64{newChain})
	newAddresses := deployment.NewMemoryAddressBook()
	err = ccipdeployment.DeployPrerequisiteChainContracts(e.Env, newAddresses, initialDeploy)
	require.NoError(t, err)
	require.NoError(t, e.Env.ExistingAddresses.Merge(newAddresses))

	tokenConfig := ccipdeployment.NewTestTokenConfig(state.Chains[e.FeedChainSel].USDFeeds)
	newAddresses = deployment.NewMemoryAddressBook()
	err = ccipdeployment.DeployCCIPContracts(e.Env, newAddresses, ccipdeployment.DeployCCIPContractConfig{
		HomeChainSel:   e.HomeChainSel,
		FeedChainSel:   e.FeedChainSel,
		ChainsToDeploy: initialDeploy,
		TokenConfig:    tokenConfig,
		MCMSConfig:     ccipdeployment.NewTestMCMSConfig(t, e.Env),
		OCRSecrets:     deployment.XXXGenerateTestOCRSecrets(),
	})
	require.NoError(t, err)
	require.NoError(t, e.Env.ExistingAddresses.Merge(newAddresses))
	state, err = ccipdeployment.LoadOnchainState(e.Env)
	require.NoError(t, err)

	// Connect all the existing lanes.
	for _, source := range initialDeploy {
		for _, dest := range initialDeploy {
			if source != dest {
				require.NoError(t, ccipdeployment.AddLaneWithDefaultPrices(e.Env, state, source, dest))
			}
		}
	}

	rmnHomeAddress, err := deployment.SearchAddressBook(e.Env.ExistingAddresses, e.HomeChainSel, ccipdeployment.RMNHome)
	require.NoError(t, err)
	require.True(t, common.IsHexAddress(rmnHomeAddress))
	rmnHome, err := rmn_home.NewRMNHome(common.HexToAddress(rmnHomeAddress), e.Env.Chains[e.HomeChainSel].Client)
	require.NoError(t, err)

	//  Deploy contracts to new chain
	newAddresses = deployment.NewMemoryAddressBook()
	err = ccipdeployment.DeployPrerequisiteChainContracts(e.Env, newAddresses, []uint64{newChain})
	require.NoError(t, err)
	require.NoError(t, e.Env.ExistingAddresses.Merge(newAddresses))
	newAddresses = deployment.NewMemoryAddressBook()
	err = ccipdeployment.DeployChainContracts(e.Env,
		e.Env.Chains[newChain], newAddresses,
		ccipdeployment.NewTestMCMSConfig(t, e.Env), rmnHome)
	require.NoError(t, err)
	require.NoError(t, e.Env.ExistingAddresses.Merge(newAddresses))
	state, err = ccipdeployment.LoadOnchainState(e.Env)
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
	tx, err = state.Chains[e.HomeChainSel].CCIPHome.TransferOwnership(e.Env.Chains[e.HomeChainSel].DeployerKey, state.Chains[e.HomeChainSel].Timelock.Address())
	require.NoError(t, err)
	_, err = deployment.ConfirmIfNoError(e.Env.Chains[e.HomeChainSel], tx, err)
	require.NoError(t, err)

	acceptOwnershipProposal, err := ccipdeployment.GenerateAcceptOwnershipProposal(state, e.HomeChainSel, initialDeploy)
	require.NoError(t, err)
	acceptOwnershipExec := ccipdeployment.SignProposal(t, e.Env, acceptOwnershipProposal)
	// Apply the accept ownership proposal to all the chains.
	for _, sel := range initialDeploy {
		ccipdeployment.ExecuteProposal(t, e.Env, acceptOwnershipExec, state, sel)
	}
	for _, chain := range initialDeploy {
		owner, err2 := state.Chains[chain].OnRamp.Owner(nil)
		require.NoError(t, err2)
		require.Equal(t, state.Chains[chain].Timelock.Address(), owner)
	}
	cfgOwner, err := state.Chains[e.HomeChainSel].CCIPHome.Owner(nil)
	require.NoError(t, err)
	crOwner, err := state.Chains[e.HomeChainSel].CapabilityRegistry.Owner(nil)
	require.NoError(t, err)
	require.Equal(t, state.Chains[e.HomeChainSel].Timelock.Address(), cfgOwner)
	require.Equal(t, state.Chains[e.HomeChainSel].Timelock.Address(), crOwner)

	nodes, err := deployment.NodeInfo(e.Env.NodeIDs, e.Env.Offchain)
	require.NoError(t, err)

	// Generate and sign inbound proposal to new 4th chain.
	chainInboundChangeset, err := NewChainInboundChangeset(e.Env, state, e.HomeChainSel, newChain, initialDeploy)
	require.NoError(t, err)
	ccipdeployment.ProcessChangeset(t, e.Env, chainInboundChangeset)

	// TODO This currently is not working - Able to send the request here but request gets stuck in execution
	// Send a new message and expect that this is delivered once the chain is completely set up as inbound
	//TestSendRequest(t, e.Env, state, initialDeploy[0], newChain, true)

	t.Logf("Executing add don and set candidate proposal for commit plugin on chain %d", newChain)
	addDonChangeset, err := AddDonAndSetCandidateChangeset(state, e.Env, nodes, deployment.XXXGenerateTestOCRSecrets(), e.HomeChainSel, e.FeedChainSel, newChain, tokenConfig, types.PluginTypeCCIPCommit)
	require.NoError(t, err)
	ccipdeployment.ProcessChangeset(t, e.Env, addDonChangeset)

	t.Logf("Executing promote candidate proposal for exec plugin on chain %d", newChain)
	setCandidateForExecChangeset, err := SetCandidatePluginChangeset(state, e.Env, nodes, deployment.XXXGenerateTestOCRSecrets(), e.HomeChainSel, e.FeedChainSel, newChain, tokenConfig, types.PluginTypeCCIPExec)
	require.NoError(t, err)
	ccipdeployment.ProcessChangeset(t, e.Env, setCandidateForExecChangeset)

	t.Logf("Executing promote candidate proposal for both commit and exec plugins on chain %d", newChain)
	donPromoteChangeset, err := PromoteAllCandidatesChangeset(state, e.HomeChainSel, newChain, nodes)
	require.NoError(t, err)
	ccipdeployment.ProcessChangeset(t, e.Env, donPromoteChangeset)

	// verify if the configs are updated
	require.NoError(t, ccipdeployment.ValidateCCIPHomeConfigSetUp(
		state.Chains[e.HomeChainSel].CapabilityRegistry,
		state.Chains[e.HomeChainSel].CCIPHome,
		newChain,
	))
	replayBlocks, err := ccipdeployment.LatestBlocksByChain(testcontext.Get(t), e.Env.Chains)
	require.NoError(t, err)

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
	latestDON, err := ccipdeployment.LatestCCIPDON(state.Chains[e.HomeChainSel].CapabilityRegistry)
	require.NoError(t, err)
	ocrConfigs, err := ccipdeployment.BuildSetOCR3ConfigArgs(latestDON.Id, state.Chains[e.HomeChainSel].CCIPHome, newChain)
	require.NoError(t, err)
	tx, err = state.Chains[newChain].OffRamp.SetOCR3Configs(e.Env.Chains[newChain].DeployerKey, ocrConfigs)
	require.NoError(t, err)
	_, err = deployment.ConfirmIfNoError(e.Env.Chains[newChain], tx, err)
	require.NoError(t, err)

	// Assert the inbound lanes to the new chain are wired correctly.
	state, err = ccipdeployment.LoadOnchainState(e.Env)
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
	ccipdeployment.ReplayLogs(t, e.Env.Offchain, replayBlocks)

	// TODO: Send via all inbound lanes and use parallel helper
	// Now that the proposal has been executed we expect to be able to send traffic to this new 4th chain.
	latesthdr, err := e.Env.Chains[newChain].Client.HeaderByNumber(testcontext.Get(t), nil)
	require.NoError(t, err)
	startBlock := latesthdr.Number.Uint64()
	msgSentEvent := ccipdeployment.TestSendRequest(t, e.Env, state, initialDeploy[0], newChain, true, router.ClientEVM2AnyMessage{
		Receiver:     common.LeftPadBytes(state.Chains[newChain].Receiver.Address().Bytes(), 32),
		Data:         []byte("hello world"),
		TokenAmounts: nil,
		FeeToken:     common.HexToAddress("0x0"),
		ExtraArgs:    nil,
	})
	require.NoError(t,
		ccipdeployment.ConfirmCommitWithExpectedSeqNumRange(t, e.Env.Chains[initialDeploy[0]], e.Env.Chains[newChain], state.Chains[newChain].OffRamp, &startBlock, cciptypes.SeqNumRange{
			cciptypes.SeqNum(1),
			cciptypes.SeqNum(msgSentEvent.SequenceNumber),
		}))
	require.NoError(t,
		commonutils.JustError(ccipdeployment.ConfirmExecWithSeqNr(t, e.Env.Chains[initialDeploy[0]], e.Env.Chains[newChain], state.Chains[newChain].OffRamp, &startBlock, msgSentEvent.SequenceNumber)))

	linkAddress := state.Chains[newChain].LinkToken.Address()
	feeQuoter := state.Chains[newChain].FeeQuoter
	timestampedPrice, err := feeQuoter.GetTokenPrice(nil, linkAddress)
	require.NoError(t, err)
	require.Equal(t, ccipdeployment.MockLinkPrice, timestampedPrice.Value)
}
