package changeset

import (
	"testing"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment"
	ccdeploy "github.com/smartcontractkit/chainlink/deployment/ccip"
	"github.com/smartcontractkit/chainlink/deployment/common/view/v1_0"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestDeployHomeChain(t *testing.T) {
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     2,
		Nodes:      4,
	})
	homeChainSel := e.AllChainSelectors()[0]
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	require.NoError(t, err)
	p2pIds := nodes.NonBootstraps().PeerIDs()
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
	state, err := ccdeploy.LoadOnchainState(e)
	require.NoError(t, err)
	require.NotNil(t, state.Chains[homeChainSel].CapabilityRegistry)
	require.NotNil(t, state.Chains[homeChainSel].CCIPHome)
	require.NotNil(t, state.Chains[homeChainSel].RMNHome)
	snap, err := state.View([]uint64{homeChainSel})
	require.NoError(t, err)
	chainid, err := chainsel.ChainIdFromSelector(homeChainSel)
	require.NoError(t, err)
	chainName, err := chainsel.NameFromChainId(chainid)
	require.NoError(t, err)
	_, ok := snap[chainName]
	require.True(t, ok)
	capRegSnap, ok := snap[chainName].CapabilityRegistry[state.Chains[homeChainSel].CapabilityRegistry.Address().String()]
	require.True(t, ok)
	require.NotNil(t, capRegSnap)
	require.Equal(t, capRegSnap.Nops, []v1_0.NopView{
		{
			Admin: e.Chains[homeChainSel].DeployerKey.From,
			Name:  "NodeOperator",
		},
	})
	require.Len(t, capRegSnap.Nodes, len(p2pIds))
}
