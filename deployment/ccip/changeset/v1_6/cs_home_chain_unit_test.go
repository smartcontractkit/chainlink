//go:build !integration

package v1_6_test

import (
	"testing"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/common/view/v1_0"
)

func TestDeployHomeChain(t *testing.T) {
	t.Parallel()

	homeChainSel := chain_selectors.TEST_90000001.Selector
	env, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{homeChainSel}),
		environment.WithLogger(logger.Test(t)),
	)
	require.NoError(t, err)

	testhelpers.RegisterNodes(t, env, 4, homeChainSel)

	e := *env

	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	require.NoError(t, err)
	p2pIds := nodes.NonBootstraps().PeerIDs()
	homeChainCfg := v1_6.DeployHomeChainConfig{
		HomeChainSel:     homeChainSel,
		RMNStaticConfig:  testhelpers.NewTestRMNStaticConfig(),
		RMNDynamicConfig: testhelpers.NewTestRMNDynamicConfig(),
		NodeOperators:    testhelpers.NewTestNodeOperator(e.BlockChains.EVMChains()[homeChainSel].DeployerKey.From),
		NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
			"NodeOperator": p2pIds,
		},
	}
	output, err := v1_6.DeployHomeChainChangeset(e, homeChainCfg)
	require.NoError(t, err)
	require.NoError(t, e.ExistingAddresses.Merge(output.AddressBook))
	state, err := stateview.LoadOnchainState(e)
	require.NoError(t, err)
	require.NotNil(t, state.Chains[homeChainSel].CapabilityRegistry)
	require.NotNil(t, state.Chains[homeChainSel].CCIPHome)
	require.NotNil(t, state.Chains[homeChainSel].RMNHome)
	view, err := state.View(&e, []uint64{homeChainSel})
	require.NoError(t, err)
	chainName := e.BlockChains.EVMChains()[homeChainSel].Name()
	_, ok := view.Chains[chainName]
	require.True(t, ok)
	capRegSnap, ok := view.Chains[chainName].CapabilityRegistry[state.Chains[homeChainSel].CapabilityRegistry.Address().String()]
	require.True(t, ok)
	require.NotNil(t, capRegSnap)
	require.Equal(t, []v1_0.NopView{
		{
			Admin: e.BlockChains.EVMChains()[homeChainSel].DeployerKey.From,
			Name:  "NodeOperator",
		},
	}, capRegSnap.Nops)
	require.Len(t, capRegSnap.Nodes, len(p2pIds))
}
