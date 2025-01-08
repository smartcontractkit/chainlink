package changeset

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/smartcontractkit/chainlink-ccip/chainconfig"
	solConfig "github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/config"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solCommomUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
)

// TODO: Solana re-write

func TestDeployChainContractsChangeset(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     2,
		Nodes:      4,
	})
	fmt.Println("Created Env")
	selectors := e.AllChainSelectors()
	homeChainSel := selectors[0]
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	require.NoError(t, err)
	p2pIds := nodes.NonBootstraps().PeerIDs()
	cfg := make(map[uint64]commontypes.MCMSWithTimelockConfig)
	for _, chain := range e.AllChainSelectors() {
		cfg[chain] = proposalutils.SingleGroupTimelockConfig(t)
	}
	var prereqCfg []DeployPrerequisiteConfigPerChain
	for _, chain := range e.AllChainSelectors() {
		prereqCfg = append(prereqCfg, DeployPrerequisiteConfigPerChain{
			ChainSelector: chain,
		})
	}
	e, err = commonchangeset.ApplyChangesets(t, e, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(DeployHomeChain),
			Config: DeployHomeChainConfig{
				HomeChainSel:     homeChainSel,
				RMNStaticConfig:  NewTestRMNStaticConfig(),
				RMNDynamicConfig: NewTestRMNDynamicConfig(),
				NodeOperators:    NewTestNodeOperator(e.Chains[homeChainSel].DeployerKey.From),
				NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
					"NodeOperator": p2pIds,
				},
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployLinkToken),
			Config:    selectors,
		},
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployMCMSWithTimelock),
			Config:    cfg,
		},
		{
			Changeset: commonchangeset.WrapChangeSet(DeployPrerequisites),
			Config: DeployPrerequisiteConfig{
				Configs: prereqCfg,
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(DeployChainContracts),
			Config: DeployChainContractsConfig{
				ChainSelectors:    selectors,
				HomeChainSelector: homeChainSel,
			},
		},
	})
	require.NoError(t, err)

	// load onchain state
	state, err := LoadOnchainState(e)
	require.NoError(t, err)

	// verify all contracts populated
	require.NotNil(t, state.Chains[homeChainSel].CapabilityRegistry)
	require.NotNil(t, state.Chains[homeChainSel].CCIPHome)
	require.NotNil(t, state.Chains[homeChainSel].RMNHome)
	for _, sel := range selectors {
		require.NotNil(t, state.Chains[sel].LinkToken)
		require.NotNil(t, state.Chains[sel].Weth9)
		require.NotNil(t, state.Chains[sel].TokenAdminRegistry)
		require.NotNil(t, state.Chains[sel].RegistryModule)
		require.NotNil(t, state.Chains[sel].Router)
		require.NotNil(t, state.Chains[sel].RMNRemote)
		require.NotNil(t, state.Chains[sel].TestRouter)
		require.NotNil(t, state.Chains[sel].NonceManager)
		require.NotNil(t, state.Chains[sel].FeeQuoter)
		require.NotNil(t, state.Chains[sel].OffRamp)
		require.NotNil(t, state.Chains[sel].OnRamp)
	}
}

func TestDeployChainContractsChangesetSolana(t *testing.T) {
	t.Parallel()
	e := NewMemoryEnvironment(t)
	fmt.Println("Created Env")
	selectors := e.Env.AllChainSelectors()
	homeChainSel := selectors[0]
	allChains := []uint64{deployment.SolanaChainSelector}
	nodes, err := deployment.NodeInfo(e.Env.NodeIDs, e.Env.Offchain)
	require.NoError(t, err)
	// p2pIds := nodes.NonBootstraps().PeerIDs()
	cfg := make(map[uint64]commontypes.MCMSWithTimelockConfig)
	for _, chain := range e.Env.AllChainSelectors() {
		cfg[chain] = proposalutils.SingleGroupTimelockConfig(t)
	}
	// var prereqCfg []DeployPrerequisiteConfigPerChain
	// for _, chain := range e.AllChainSelectors() {
	// 	prereqCfg = append(prereqCfg, DeployPrerequisiteConfigPerChain{
	// 		ChainSelector: chain,
	// 	})
	// }
	fmt.Println(e.Env.SolChains)
	e.Env, err = commonchangeset.ApplyChangesets(t, e.Env, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployLinkTokenSolana),
			Config:    []uint64{deployment.SolanaChainSelector},
		},
		// {
		// 	Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployMCMSWithTimelock),
		// 	Config:    cfg,
		// },
		{
			Changeset: commonchangeset.WrapChangeSet(DeployChainContractsSolana),
			Config: DeployChainContractsConfig{
				ChainSelectors:    []uint64{deployment.SolanaChainSelector},
				HomeChainSelector: homeChainSel,
			},
		},
	})
	require.NoError(t, err)
	state, err := LoadOnchainStateSolana(e.Env)
	require.NoError(t, err)
	// Assert link present
	require.NotNil(t, state.SolChains[deployment.SolanaChainSelector].LinkToken)
	require.NotNil(t, state.SolChains[deployment.SolanaChainSelector].Weth9)

	tokenConfig := NewTestTokenConfig(
		state.SolChains[deployment.SolanaChainSelector].LinkToken.String(),
		state.SolChains[deployment.SolanaChainSelector].Weth9.String(),
		deployment.SolanaChainSelector)
	var tokenDataProviders []pluginconfig.TokenDataObserverConfig
	// Build the per chain config.
	ocrConfigs := make(map[uint64]CCIPOCRParams)
	chainConfigs := make(map[uint64]ChainConfig)
	for _, chain := range allChains {
		tokenInfo := tokenConfig.GetTokenInfo(e.Env.Logger, state.SolChains[chain].LinkToken.String(), state.SolChains[chain].Weth9.String())
		ocrParams := DefaultOCRParams(deployment.SolanaChainSelector, tokenInfo, tokenDataProviders)
		ocrConfigs[chain] = ocrParams
		chainConfigs[chain] = ChainConfig{
			Readers: nodes.NonBootstraps().PeerIDs(),
			FChain:  uint8(len(nodes.NonBootstraps().PeerIDs()) / 3),
			EncodableChainConfig: chainconfig.ChainConfig{
				GasPriceDeviationPPB:    cciptypes.BigInt{Int: big.NewInt(internal.GasPriceDeviationPPB)},
				DAGasPriceDeviationPPB:  cciptypes.BigInt{Int: big.NewInt(internal.DAGasPriceDeviationPPB)},
				OptimisticConfirmations: internal.OptimisticConfirmations,
			},
		}
	}
	// Deploy second set of changesets to deploy and configure the CCIP contracts.
	e.Env, err = commonchangeset.ApplyChangesets(t, e.Env, nil, []commonchangeset.ChangesetApplication{
		{
			// Add the chain configs for the new chains.
			Changeset: commonchangeset.WrapChangeSet(UpdateChainConfig),
			Config: UpdateChainConfigConfig{
				HomeChainSelector: homeChainSel,
				RemoteChainAdds:   chainConfigs,
			},
		},
		// For everything below, we need node spinup to support Solana OCR
		{
			// Add the DONs and candidate commit OCR instances for the chain.
			Changeset: commonchangeset.WrapChangeSet(AddDonAndSetCandidateChangeset),
			Config: AddDonAndSetCandidateChangesetConfig{
				SetCandidateConfigBase{
					HomeChainSelector:               homeChainSel,
					FeedChainSelector:               deployment.SolanaChainSelector,
					OCRConfigPerRemoteChainSelector: ocrConfigs,
					PluginType:                      types.PluginTypeCCIPCommit,
				},
			},
		},
		{
			// Add the exec OCR instances for the new chains.
			Changeset: commonchangeset.WrapChangeSet(SetCandidateChangeset),
			Config: SetCandidateChangesetConfig{
				SetCandidateConfigBase{
					HomeChainSelector:               homeChainSel,
					FeedChainSelector:               deployment.SolanaChainSelector,
					OCRConfigPerRemoteChainSelector: ocrConfigs,
					PluginType:                      types.PluginTypeCCIPExec,
				},
			},
		},
		{
			// Promote everything
			Changeset: commonchangeset.WrapChangeSet(PromoteAllCandidatesChangeset),
			Config: PromoteCandidatesChangesetConfig{
				HomeChainSelector:    homeChainSel,
				RemoteChainSelectors: allChains,
				PluginType:           types.PluginTypeCCIPCommit,
			},
		},
		{
			// Promote everything
			Changeset: commonchangeset.WrapChangeSet(PromoteAllCandidatesChangeset),
			Config: PromoteCandidatesChangesetConfig{
				HomeChainSelector:    homeChainSel,
				RemoteChainSelectors: allChains,
				PluginType:           types.PluginTypeCCIPExec,
			},
		},
		{
			// Enable the OCR config on the remote chains.
			Changeset: commonchangeset.WrapChangeSet(SetOCR3ConfigSolana),
			Config: SetOCR3OffRampConfig{
				HomeChainSel:    homeChainSel,
				RemoteChainSels: allChains,
			},
		},
	})
	require.NoError(t, err)

	// load onchain state
	state, err = LoadOnchainStateSolana(e.Env)
	require.NoError(t, err)

	// // verify all contracts populated
	// require.NotNil(t, state.Chains[homeChainSel].CapabilityRegistry)
	// require.NotNil(t, state.Chains[homeChainSel].CCIPHome)
	// require.NotNil(t, state.Chains[homeChainSel].RMNHome)
	// for _, sel := range selectors {
	require.NotNil(t, state.SolChains[deployment.SolanaChainSelector].LinkToken)
	// 	require.NotNil(t, state.Chains[sel].Weth9)
	// 	require.NotNil(t, state.Chains[sel].TokenAdminRegistry)
	// 	require.NotNil(t, state.Chains[sel].RegistryModule)
	require.NotNil(t, state.SolChains[deployment.SolanaChainSelector].CcipRouter)
	// 	require.NotNil(t, state.Chains[sel].RMNRemote)
	// 	require.NotNil(t, state.Chains[sel].TestRouter)
	// 	require.NotNil(t, state.Chains[sel].NonceManager)
	// 	require.NotNil(t, state.Chains[sel].FeeQuoter)
	// 	require.NotNil(t, state.Chains[sel].OffRamp)
	// 	require.NotNil(t, state.Chains[sel].OnRamp)
	// }

	var configAccount ccip_router.Config
	err = solCommomUtil.GetAccountDataBorshInto(e.Env.GetContext(), e.Env.SolChains[deployment.SolanaChainSelector].Client, GetRouterConfigPDA(state.SolChains[deployment.SolanaChainSelector].CcipRouter), solConfig.DefaultCommitment, &configAccount)
	require.NoError(t, err)
	require.Equal(t, deployment.SolanaChainSelector, configAccount.SolanaChainSelector)

}

func TestDeployCCIPContracts(t *testing.T) {
	t.Parallel()
	e := NewMemoryEnvironment(t)
	// Deploy all the CCIP contracts.
	// TODO: not sure how this ends up deploying contracts ?
	state, err := LoadOnchainState(e.Env)
	require.NoError(t, err)
	snap, err := state.View(e.Env.AllChainSelectors())
	require.NoError(t, err)

	// Assert expect every deployed address to be in the address book.
	// TODO (CCIP-3047): Add the rest of CCIPv2 representation
	b, err := json.MarshalIndent(snap, "", "	")
	require.NoError(t, err)
	fmt.Println(string(b))
}
