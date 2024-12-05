package changeset

import (
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
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
	e := NewMemoryEnvironmentWithJobs(t, logger.TestLogger(t), memory.MemoryEnvironmentConfig{
		Chains:             4,
		NumOfUsersPerChain: 1,
		Nodes:              4,
		Bootstraps:         1,
	})
	state, err := LoadOnchainState(e.Env)
	require.NoError(t, err)
	// Take first non-home chain as the new chain.
	newChain := e.Env.AllChainSelectorsExcluding([]uint64{e.HomeChainSel})[0]
	// We deploy to the rest.
	initialDeploy := e.Env.AllChainSelectorsExcluding([]uint64{newChain})
	newAddresses := deployment.NewMemoryAddressBook()
	err = deployPrerequisiteChainContracts(e.Env, newAddresses, initialDeploy, nil)
	require.NoError(t, err)
	require.NoError(t, e.Env.ExistingAddresses.Merge(newAddresses))

	cfg := commontypes.MCMSWithTimelockConfig{
		Canceller:         commonchangeset.SingleGroupMCMS(t),
		Bypasser:          commonchangeset.SingleGroupMCMS(t),
		Proposer:          commonchangeset.SingleGroupMCMS(t),
		TimelockExecutors: e.Env.AllDeployerKeys(),
		TimelockMinDelay:  big.NewInt(0),
	}
	e.Env, err = commonchangeset.ApplyChangesets(t, e.Env, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployLinkToken),
			Config:    initialDeploy,
		},
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployMCMSWithTimelock),
			Config: map[uint64]commontypes.MCMSWithTimelockConfig{
				initialDeploy[0]: cfg,
				initialDeploy[1]: cfg,
				initialDeploy[2]: cfg,
			},
		},
	})
	require.NoError(t, err)
	newAddresses = deployment.NewMemoryAddressBook()
	tokenConfig := NewTestTokenConfig(state.Chains[e.FeedChainSel].USDFeeds)

	chainConfig := make(map[uint64]CCIPOCRParams)
	for _, chain := range initialDeploy {
		chainConfig[chain] = DefaultOCRParams(e.FeedChainSel, nil, nil)
	}
	err = deployCCIPContracts(e.Env, newAddresses, NewChainsConfig{
		HomeChainSel:       e.HomeChainSel,
		FeedChainSel:       e.FeedChainSel,
		ChainConfigByChain: chainConfig,
		OCRSecrets:         deployment.XXXGenerateTestOCRSecrets(),
	})
	require.NoError(t, err)

	state, err = LoadOnchainState(e.Env)
	require.NoError(t, err)

	// Connect all the existing lanes.
	for _, source := range initialDeploy {
		for _, dest := range initialDeploy {
			if source != dest {
				require.NoError(t, AddLaneWithDefaultPricesAndFeeQuoterConfig(e.Env, state, source, dest, false))
			}
		}
	}

	rmnHomeAddress, err := deployment.SearchAddressBook(e.Env.ExistingAddresses, e.HomeChainSel, RMNHome)
	require.NoError(t, err)
	require.True(t, common.IsHexAddress(rmnHomeAddress))
	rmnHome, err := rmn_home.NewRMNHome(common.HexToAddress(rmnHomeAddress), e.Env.Chains[e.HomeChainSel].Client)
	require.NoError(t, err)

	//  Deploy contracts to new chain
	e.Env, err = commonchangeset.ApplyChangesets(t, e.Env, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployLinkToken),
			Config:    []uint64{newChain},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployMCMSWithTimelock),
			Config: map[uint64]commontypes.MCMSWithTimelockConfig{
				newChain: cfg,
			},
		},
	})
	require.NoError(t, err)
	newAddresses = deployment.NewMemoryAddressBook()

	err = deployPrerequisiteChainContracts(e.Env, newAddresses, []uint64{newChain}, nil)
	require.NoError(t, err)
	require.NoError(t, e.Env.ExistingAddresses.Merge(newAddresses))
	newAddresses = deployment.NewMemoryAddressBook()
	err = deployChainContracts(e.Env,
		e.Env.Chains[newChain], newAddresses, rmnHome)
	require.NoError(t, err)
	require.NoError(t, e.Env.ExistingAddresses.Merge(newAddresses))
	state, err = LoadOnchainState(e.Env)
	require.NoError(t, err)

	// configure the testrouter appropriately on each chain
	for _, source := range initialDeploy {
		tx, err := state.Chains[source].TestRouter.ApplyRampUpdates(e.Env.Chains[source].DeployerKey, []router.RouterOnRamp{
			{
				DestChainSelector: newChain,
				OnRamp:            state.Chains[source].OnRamp.Address(),
			},
		}, nil, nil)
		_, err = deployment.ConfirmIfNoError(e.Env.Chains[source], tx, err)
		require.NoError(t, err)
	}

	// transfer ownership to timelock
	_, err = commonchangeset.ApplyChangesets(t, e.Env, map[uint64]*gethwrappers.RBACTimelock{
		initialDeploy[0]: state.Chains[initialDeploy[0]].Timelock,
		initialDeploy[1]: state.Chains[initialDeploy[1]].Timelock,
		initialDeploy[2]: state.Chains[initialDeploy[2]].Timelock,
	}, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.TransferToMCMSWithTimelock),
			Config:    genTestTransferOwnershipConfig(e, initialDeploy, state),
		},
		{
			Changeset: commonchangeset.WrapChangeSet(NewChainInboundChangeset),
			Config: ChainInboundChangesetConfig{
				HomeChainSelector:    e.HomeChainSel,
				NewChainSelector:     newChain,
				SourceChainSelectors: initialDeploy,
			},
		},
	})
	require.NoError(t, err)

	assertTimelockOwnership(t, e, initialDeploy, state)

	nodes, err := deployment.NodeInfo(e.Env.NodeIDs, e.Env.Offchain)
	require.NoError(t, err)

	// TODO This currently is not working - Able to send the request here but request gets stuck in execution
	// Send a new message and expect that this is delivered once the chain is completely set up as inbound
	//TestSendRequest(t, e.Env, state, initialDeploy[0], newChain, true)
	var nodeIDs []string
	for _, node := range nodes {
		nodeIDs = append(nodeIDs, node.NodeID)
	}

	_, err = commonchangeset.ApplyChangesets(t, e.Env, map[uint64]*gethwrappers.RBACTimelock{
		e.HomeChainSel: state.Chains[e.HomeChainSel].Timelock,
		newChain:       state.Chains[newChain].Timelock,
	}, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(AddDonAndSetCandidateChangeset),
			Config: AddDonAndSetCandidateChangesetConfig{
				HomeChainSelector: e.HomeChainSel,
				FeedChainSelector: e.FeedChainSel,
				NewChainSelector:  newChain,
				PluginType:        types.PluginTypeCCIPCommit,
				NodeIDs:           nodeIDs,
				OCRSecrets:        deployment.XXXGenerateTestOCRSecrets(),
				CCIPOCRParams: DefaultOCRParams(
					e.FeedChainSel,
					tokenConfig.GetTokenInfo(logger.TestLogger(t), state.Chains[newChain].LinkToken, state.Chains[newChain].Weth9),
					nil,
				),
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(SetCandidatePluginChangeset),
			Config: AddDonAndSetCandidateChangesetConfig{
				HomeChainSelector: e.HomeChainSel,
				FeedChainSelector: e.FeedChainSel,
				NewChainSelector:  newChain,
				PluginType:        types.PluginTypeCCIPExec,
				NodeIDs:           nodeIDs,
				OCRSecrets:        deployment.XXXGenerateTestOCRSecrets(),
				CCIPOCRParams: DefaultOCRParams(
					e.FeedChainSel,
					tokenConfig.GetTokenInfo(logger.TestLogger(t), state.Chains[newChain].LinkToken, state.Chains[newChain].Weth9),
					nil,
				),
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(PromoteAllCandidatesChangeset),
			Config: PromoteAllCandidatesChangesetConfig{
				HomeChainSelector: e.HomeChainSel,
				NewChainSelector:  newChain,
				NodeIDs:           nodeIDs,
			},
		},
	})

	// verify if the configs are updated
	require.NoError(t, ValidateCCIPHomeConfigSetUp(
		state.Chains[e.HomeChainSel].CapabilityRegistry,
		state.Chains[e.HomeChainSel].CCIPHome,
		newChain,
	))
	replayBlocks, err := LatestBlocksByChain(testcontext.Get(t), e.Env.Chains)
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
	tx, err := state.Chains[newChain].OffRamp.ApplySourceChainConfigUpdates(e.Env.Chains[newChain].DeployerKey, offRampEnables)
	require.NoError(t, err)
	_, err = deployment.ConfirmIfNoError(e.Env.Chains[newChain], tx, err)
	require.NoError(t, err)
	// Set the OCR3 config on new 4th chain to enable the plugin.
	latestDON, err := internal.LatestCCIPDON(state.Chains[e.HomeChainSel].CapabilityRegistry)
	require.NoError(t, err)
	ocrConfigs, err := internal.BuildSetOCR3ConfigArgs(latestDON.Id, state.Chains[e.HomeChainSel].CCIPHome, newChain)
	require.NoError(t, err)
	tx, err = state.Chains[newChain].OffRamp.SetOCR3Configs(e.Env.Chains[newChain].DeployerKey, ocrConfigs)
	require.NoError(t, err)
	_, err = deployment.ConfirmIfNoError(e.Env.Chains[newChain], tx, err)
	require.NoError(t, err)

	// Assert the inbound lanes to the new chain are wired correctly.
	state, err = LoadOnchainState(e.Env)
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
	ReplayLogs(t, e.Env.Offchain, replayBlocks)

	// TODO: Send via all inbound lanes and use parallel helper
	// Now that the proposal has been executed we expect to be able to send traffic to this new 4th chain.
	latesthdr, err := e.Env.Chains[newChain].Client.HeaderByNumber(testcontext.Get(t), nil)
	require.NoError(t, err)
	startBlock := latesthdr.Number.Uint64()
	msgSentEvent := TestSendRequest(t, e.Env, state, initialDeploy[0], newChain, true, router.ClientEVM2AnyMessage{
		Receiver:     common.LeftPadBytes(state.Chains[newChain].Receiver.Address().Bytes(), 32),
		Data:         []byte("hello world"),
		TokenAmounts: nil,
		FeeToken:     common.HexToAddress("0x0"),
		ExtraArgs:    nil,
	})
	require.NoError(t,
		commonutils.JustError(ConfirmCommitWithExpectedSeqNumRange(t, e.Env.Chains[initialDeploy[0]], e.Env.Chains[newChain], state.Chains[newChain].OffRamp, &startBlock, cciptypes.SeqNumRange{
			cciptypes.SeqNum(1),
			cciptypes.SeqNum(msgSentEvent.SequenceNumber),
		}, true)))
	require.NoError(t,
		commonutils.JustError(
			ConfirmExecWithSeqNrs(
				t,
				e.Env.Chains[initialDeploy[0]],
				e.Env.Chains[newChain],
				state.Chains[newChain].OffRamp,
				&startBlock,
				[]uint64{msgSentEvent.SequenceNumber},
			),
		),
	)

	linkAddress := state.Chains[newChain].LinkToken.Address()
	feeQuoter := state.Chains[newChain].FeeQuoter
	timestampedPrice, err := feeQuoter.GetTokenPrice(nil, linkAddress)
	require.NoError(t, err)
	require.Equal(t, MockLinkPrice, timestampedPrice.Value)
}
