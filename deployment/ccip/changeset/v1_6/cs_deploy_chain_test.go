package v1_6_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
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
	cfg := make(map[uint64]commontypes.MCMSWithTimelockConfigV2)
	contractParams := make(map[uint64]v1_6.ChainContractParams)
	for _, chain := range e.AllChainSelectors() {
		cfg[chain] = proposalutils.SingleGroupTimelockConfigV2(t)
		contractParams[chain] = v1_6.ChainContractParams{
			FeeQuoterParams: v1_6.DefaultFeeQuoterParams(),
			OffRampParams:   v1_6.DefaultOffRampParams(),
		}
	}
	prereqCfg := make([]changeset.DeployPrerequisiteConfigPerChain, 0)
	for _, chain := range e.AllChainSelectors() {
		prereqCfg = append(prereqCfg, changeset.DeployPrerequisiteConfigPerChain{
			ChainSelector: chain,
		})
	}

	e, err = commonchangeset.Apply(t, e, nil,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.DeployHomeChainChangeset),
			v1_6.DeployHomeChainConfig{
				HomeChainSel:     homeChainSel,
				RMNStaticConfig:  testhelpers.NewTestRMNStaticConfig(),
				RMNDynamicConfig: testhelpers.NewTestRMNDynamicConfig(),
				NodeOperators:    testhelpers.NewTestNodeOperator(e.Chains[homeChainSel].DeployerKey.From),
				NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
					"NodeOperator": p2pIds,
				},
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(commonchangeset.DeployLinkToken),
			evmSelectors,
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
			cfg,
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(changeset.DeployPrerequisitesChangeset),
			changeset.DeployPrerequisiteConfig{
				Configs: prereqCfg,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.DeployChainContractsChangeset),
			v1_6.DeployChainContractsConfig{
				HomeChainSelector:      homeChainSel,
				ContractParamsPerChain: contractParams,
			},
		),
	)
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
		require.NotNil(t, state.Chains[sel].RegistryModules1_6)
		require.NotNil(t, state.Chains[sel].Router)
		require.NotNil(t, state.Chains[sel].RMNRemote)
		require.NotNil(t, state.Chains[sel].TestRouter)
		require.NotNil(t, state.Chains[sel].NonceManager)
		require.NotNil(t, state.Chains[sel].FeeQuoter)
		require.NotNil(t, state.Chains[sel].OffRamp)
		require.NotNil(t, state.Chains[sel].OnRamp)
	}
	// remove feequoter from address book
	// and deploy again, it should deploy a new feequoter
	ab := deployment.NewMemoryAddressBook()
	for _, sel := range evmSelectors {
		require.NoError(t, ab.Save(sel, state.Chains[sel].FeeQuoter.Address().Hex(),
			deployment.NewTypeAndVersion(changeset.FeeQuoter, deployment.Version1_6_0)))
	}
	//nolint:staticcheck //SA1019 ignoring deprecated
	require.NoError(t, e.ExistingAddresses.Remove(ab))

	// try to deploy chain contracts again and it should not deploy any new contracts except feequoter
	// but should not error
	e, err = commonchangeset.Apply(t, e, nil, commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(v1_6.DeployChainContractsChangeset),
		v1_6.DeployChainContractsConfig{
			HomeChainSelector:      homeChainSel,
			ContractParamsPerChain: contractParams,
		},
	))
	require.NoError(t, err)
	// verify all contracts populated
	postState, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)
	for _, sel := range evmSelectors {
		require.Equal(t, state.Chains[sel].RMNRemote, postState.Chains[sel].RMNRemote)
		require.Equal(t, state.Chains[sel].Router, postState.Chains[sel].Router)
		require.Equal(t, state.Chains[sel].TestRouter, postState.Chains[sel].TestRouter)
		require.Equal(t, state.Chains[sel].NonceManager, postState.Chains[sel].NonceManager)
		require.NotEqual(t, state.Chains[sel].FeeQuoter, postState.Chains[sel].FeeQuoter)
		require.NotEmpty(t, postState.Chains[sel].FeeQuoter)
		require.Equal(t, state.Chains[sel].OffRamp, postState.Chains[sel].OffRamp)
		require.Equal(t, state.Chains[sel].OnRamp, postState.Chains[sel].OnRamp)
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
