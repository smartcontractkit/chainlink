package changeset

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment"
	ccdeploy "github.com/smartcontractkit/chainlink/deployment/ccip"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestDeployChainContractsChangeset(t *testing.T) {
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     2,
		Nodes:      4,
	})
	selectors := e.AllChainSelectors()
	homeChainSel := selectors[0]
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	require.NoError(t, err)
	p2pIds := nodes.NonBootstraps().PeerIDs()
	// deploy home chain
	homeChainCfg := DeployHomeChainConfig{
		HomeChainSel:     homeChainSel,
		RMNStaticConfig:  ccdeploy.NewTestRMNStaticConfig(),
		RMNDynamicConfig: ccdeploy.NewTestRMNDynamicConfig(),
		NodeOperators:    ccdeploy.NewTestNodeOperator(e.Chains[homeChainSel].DeployerKey.From),
		NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
			"NodeOperator": p2pIds,
		},
	}
	output, err := DeployHomeChain(e, homeChainCfg)
	require.NoError(t, err)
	require.NoError(t, e.ExistingAddresses.Merge(output.AddressBook))

	// deploy pre-requisites
	prerequisites, err := DeployPrerequisites(e, DeployPrerequisiteConfig{
		ChainSelectors: selectors,
	})
	require.NoError(t, err)
	require.NoError(t, e.ExistingAddresses.Merge(prerequisites.AddressBook))

	// deploy ccip chain contracts
	output, err = DeployChainContracts(e, DeployChainContractsConfig{
		ChainSelectors:    selectors,
		HomeChainSelector: homeChainSel,
		MCMSCfg:           ccdeploy.NewTestMCMSConfig(t, e),
	})
	require.NoError(t, err)
	require.NoError(t, e.ExistingAddresses.Merge(output.AddressBook))

	// load onchain state
	state, err := ccdeploy.LoadOnchainState(e)
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
