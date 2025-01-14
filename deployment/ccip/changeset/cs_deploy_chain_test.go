package changeset

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-ccip/chainconfig"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestDeployChainContractsChangeset(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     2,
		SolChains:  1,
		Nodes:      4,
	})
	evmSelectors := e.AllChainSelectors()
	homeChainSel := evmSelectors[0]
	solChainSelectors := e.AllChainSelectorsSolana()
	selectors := make([]uint64, 0, len(evmSelectors)+len(solChainSelectors))
	selectors = append(selectors, evmSelectors...)
	selectors = append(selectors, solChainSelectors...)
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

	SavePreloadedSolAddresses(t, e, solChainSelectors[0])
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
	for _, sel := range evmSelectors {
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

	solState, err := LoadOnchainStateSolana(e)
	require.NoError(t, err)
	for _, sel := range solChainSelectors {
		require.NotNil(t, solState.SolChains[sel].LinkToken)
		require.NotNil(t, solState.SolChains[sel].SolCcipRouter)
	}

}

func TestDeployCCIPContracts(t *testing.T) {
	t.Parallel()
	e := NewMemoryEnvironment(t)
	// Deploy all the CCIP contracts.
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

func TestHomeChainChangesetSolana(t *testing.T) {
	t.Parallel()
	e := NewMemoryEnvironment(t, WithSolChains(1))
	evmSelectors := e.Env.AllChainSelectors()
	homeChainSel := evmSelectors[0]
	solChainSelectors := e.Env.AllChainSelectorsSolana()
	nodes, err := deployment.NodeInfo(e.Env.NodeIDs, e.Env.Offchain)
	require.NoError(t, err)
	cfg := make(map[uint64]commontypes.MCMSWithTimelockConfig)
	for _, chain := range e.Env.AllChainSelectors() {
		cfg[chain] = proposalutils.SingleGroupTimelockConfig(t)
	}
	SavePreloadedSolAddresses(t, e.Env, solChainSelectors[0])
	e.Env, err = commonchangeset.ApplyChangesets(t, e.Env, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployLinkToken),
			Config:    solChainSelectors,
		},
		{
			Changeset: commonchangeset.WrapChangeSet(DeployChainContracts),
			Config: DeployChainContractsConfig{
				ChainSelectors:    solChainSelectors,
				HomeChainSelector: homeChainSel,
			},
		},
	})
	require.NoError(t, err)
	solState, err := LoadOnchainStateSolana(e.Env)
	require.NoError(t, err)
	for _, sel := range solChainSelectors {
		require.NotNil(t, solState.SolChains[sel].LinkToken)
		require.NotNil(t, solState.SolChains[sel].SolCcipRouter)
	}

	// Build the per chain config.
	ocrConfigs := make(map[uint64]CCIPOCRParams)
	chainConfigs := make(map[uint64]ChainConfig)
	for _, chain := range solChainSelectors {
		var tokenDataProviders []pluginconfig.TokenDataObserverConfig
		ocrParams := DefaultOCRParams(chain, map[cciptypes.UnknownEncodedAddress]pluginconfig.TokenInfo{}, tokenDataProviders)
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
					FeedChainSelector:               solChainSelectors[0],
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
					FeedChainSelector:               solChainSelectors[0],
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
				RemoteChainSelectors: solChainSelectors,
				PluginType:           types.PluginTypeCCIPCommit,
			},
		},
		{
			// Promote everything
			Changeset: commonchangeset.WrapChangeSet(PromoteAllCandidatesChangeset),
			Config: PromoteCandidatesChangesetConfig{
				HomeChainSelector:    homeChainSel,
				RemoteChainSelectors: solChainSelectors,
				PluginType:           types.PluginTypeCCIPExec,
			},
		},
		{
			// Enable the OCR config on the remote chains.
			Changeset: commonchangeset.WrapChangeSet(SetOCR3ConfigSolana),
			Config: SetOCR3OffRampConfig{
				HomeChainSel:    homeChainSel,
				RemoteChainSels: solChainSelectors,
			},
		},
	})
	require.NoError(t, err)

}
