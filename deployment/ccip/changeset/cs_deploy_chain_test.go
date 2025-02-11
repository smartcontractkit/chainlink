package changeset_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestDeployChainContractsChangeset(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     2,
		Nodes:      4,
	})
	evmSelectors := e.AllChainSelectors()
	homeChainSel := evmSelectors[0]
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	require.NoError(t, err)
	p2pIds := nodes.NonBootstraps().PeerIDs()
	cfg := make(map[uint64]commontypes.MCMSWithTimelockConfig)
	contractParams := make(map[uint64]changeset.ChainContractParams)
	for _, chain := range e.AllChainSelectors() {
		cfg[chain] = proposalutils.SingleGroupTimelockConfig(t)
		contractParams[chain] = changeset.ChainContractParams{
			FeeQuoterParams: changeset.DefaultFeeQuoterParams(),
			OffRampParams:   changeset.DefaultOffRampParams(),
		}
	}
	prereqCfg := make([]changeset.DeployPrerequisiteConfigPerChain, 0)
	for _, chain := range e.AllChainSelectors() {
		prereqCfg = append(prereqCfg, changeset.DeployPrerequisiteConfigPerChain{
			ChainSelector: chain,
		})
	}

	e, err = commonchangeset.ApplyChangesets(t, e, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.DeployHomeChainChangeset),
			Config: changeset.DeployHomeChainConfig{
				HomeChainSel:     homeChainSel,
				RMNStaticConfig:  testhelpers.NewTestRMNStaticConfig(),
				RMNDynamicConfig: testhelpers.NewTestRMNDynamicConfig(),
				NodeOperators:    testhelpers.NewTestNodeOperator(e.Chains[homeChainSel].DeployerKey.From),
				NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
					"NodeOperator": p2pIds,
				},
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployLinkToken),
			Config:    evmSelectors,
		},
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployMCMSWithTimelock),
			Config:    cfg,
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.DeployPrerequisitesChangeset),
			Config: changeset.DeployPrerequisiteConfig{
				Configs: prereqCfg,
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.DeployChainContractsChangeset),
			Config: changeset.DeployChainContractsConfig{
				HomeChainSelector:      homeChainSel,
				ContractParamsPerChain: contractParams,
			},
		},
	})
	require.NoError(t, err)

	// load onchain state
	state, err := changeset.LoadOnchainState(e)
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
}

func TestDeployCCIPContracts(t *testing.T) {
	t.Parallel()
	testhelpers.DeployCCIPContractsTest(t, 0)
}

func TestDeployStaticLinkToken(t *testing.T) {
	t.Parallel()
	e, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithStaticLink())
	// load onchain state
	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	for _, chain := range e.Env.AllChainSelectors() {
		require.NotNil(t, state.Chains[chain].StaticLinkToken)
	}
}
