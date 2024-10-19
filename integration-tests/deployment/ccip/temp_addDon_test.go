package ccipdeployment

import (
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_home"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func Test_addDon_test(t *testing.T) {
	// 4 chains where the 4th is added after initial deployment.
	e := NewMemoryEnvironmentWithJobs(t, logger.TestLogger(t), 4)
	state, err := LoadOnchainState(e.Env, e.Ab)
	require.NoError(t, err)
	// Take first non-home chain as the new chain.
	newChain := e.Env.AllChainSelectorsExcluding([]uint64{e.HomeChainSel})[0]
	// We deploy to the rest.
	initialDeploy := e.Env.AllChainSelectorsExcluding([]uint64{newChain})

	feeds := state.Chains[e.FeedChainSel].USDFeeds
	tokenConfig := NewTokenConfig()
	tokenConfig.UpsertTokenInfo(LinkSymbol,
		pluginconfig.TokenInfo{
			AggregatorAddress: feeds[LinkSymbol].Address().String(),
			Decimals:          LinkDecimals,
			DeviationPPB:      cciptypes.NewBigIntFromInt64(1e9),
		},
	)
	err = DeployCCIPContracts(e.Env, e.Ab, DeployCCIPContractConfig{
		HomeChainSel:       e.HomeChainSel,
		FeedChainSel:       e.FeedChainSel,
		ChainsToDeploy:     initialDeploy,
		TokenConfig:        tokenConfig,
		MCMSConfig:         NewTestMCMSConfig(t, e.Env),
		FeeTokenContracts:  e.FeeTokenContracts,
		CapabilityRegistry: state.Chains[e.HomeChainSel].CapabilityRegistry.Address(),
	})
	require.NoError(t, err)
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

	rmnHomeAddress, err := deployment.SearchAddressBook(e.Ab, e.HomeChainSel, RMNHome)
	require.NoError(t, err)
	require.True(t, common.IsHexAddress(rmnHomeAddress))
	rmnHome, err := rmn_home.NewRMNHome(common.HexToAddress(rmnHomeAddress), e.Env.Chains[e.HomeChainSel].Client)
	require.NoError(t, err)

	//  Deploy contracts to new chain
	err = DeployChainContracts(e.Env, e.Env.Chains[newChain], e.Ab, e.FeeTokenContracts[newChain], NewTestMCMSConfig(t, e.Env), rmnHome)
	require.NoError(t, err)
	state, err = LoadOnchainState(e.Env, e.Ab)
	require.NoError(t, err)

	nodes, err := deployment.NodeInfo(e.Env.NodeIDs, e.Env.Offchain)
	require.NoError(t, err)

	newDONArgs, err := BuildOCR3ConfigForCCIPHome(
		e.Env.Logger,
		state.Chains[newChain].OffRamp,
		e.Env.Chains[newChain],
		e.FeedChainSel,
		tokenConfig.GetTokenInfo(e.Env.Logger, state.Chains[newChain].LinkToken),
		nodes.NonBootstraps(),
		common.HexToAddress(rmnHomeAddress),
	)
	require.NoError(t, err)

	allDons, err := state.Chains[e.HomeChainSel].CapabilityRegistry.GetDONs(nil)

	_, err = state.Chains[e.HomeChainSel].CapabilityRegistry.AddDON(
		e.Env.Chains[e.HomeChainSel].DeployerKey,
		nodes.PeerIDs(),
		[]capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{},
		false,
		false, nodes.DefaultF(),
	)
	require.NoError(t, err)

	newDons, err := state.Chains[e.HomeChainSel].CapabilityRegistry.GetDONs(nil)
	require.Equal(t, len(newDons), len(allDons)+1)

	// fetch DON ID for the chain
	donID, err := DonIDForChain(
		state.Chains[e.HomeChainSel].CapabilityRegistry,
		state.Chains[e.HomeChainSel].CCIPHome,
		newChain)
	require.NoError(t, err)
	// donID was provisioned?
	require.NotNil(t, donID)

	encodedSetCandidateCall, err := CCIPHomeABI.Pack(
		"setCandidate",
		donID,
		cctypes.PluginTypeCCIPExec,
		newDONArgs[cctypes.PluginTypeCCIPExec],
		[32]byte{},
	)
	require.NoError(t, err)

	_, err = state.Chains[e.HomeChainSel].CapabilityRegistry.UpdateDON(
		e.Env.Chains[e.HomeChainSel].DeployerKey,
		donID,
		nodes.PeerIDs(),
		[]capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
			{
				CapabilityId: CCIPCapabilityID,
				Config:       encodedSetCandidateCall,
			},
		},
		false,
		nodes.DefaultF())
	require.NoError(t, err)
}
