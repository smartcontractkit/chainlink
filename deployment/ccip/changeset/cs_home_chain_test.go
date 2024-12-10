package changeset

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
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
		RMNStaticConfig:  NewTestRMNStaticConfig(),
		RMNDynamicConfig: NewTestRMNDynamicConfig(),
		NodeOperators:    NewTestNodeOperator(e.Chains[homeChainSel].DeployerKey.From),
		NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
			"NodeOperator": p2pIds,
		},
	}
	output, err := DeployHomeChain(e, homeChainCfg)
	require.NoError(t, err)
	require.NoError(t, e.ExistingAddresses.Merge(output.AddressBook))
	state, err := LoadOnchainState(e)
	require.NoError(t, err)
	require.NotNil(t, state.Chains[homeChainSel].CapabilityRegistry)
	require.NotNil(t, state.Chains[homeChainSel].CCIPHome)
	require.NotNil(t, state.Chains[homeChainSel].RMNHome)
	snap, err := state.View([]uint64{homeChainSel})
	require.NoError(t, err)
	chainName := e.Chains[homeChainSel].Name()
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

func TestRemoveDons(t *testing.T) {
	e := NewMemoryEnvironment(t)
	s, err := LoadOnchainState(e.Env)
	require.NoError(t, err)
	homeChain := s.Chains[e.HomeChainSel]

	// Invalid dons
	a := RemoveDONsConfig{
		HomeChainSel: e.HomeChainSel,
		DonIDs:       []uint32{},
	}
	require.Error(t, a.Validate(homeChain))

	// Remove a don w/o MCMS
	donsBefore, err := homeChain.CapabilityRegistry.GetDONs(nil)
	require.NoError(t, err)
	e.Env, err = commoncs.ApplyChangesets(t, e.Env, map[uint64]*commoncs.TimelockExecutionContracts{
		e.HomeChainSel: {
			Timelock:  s.Chains[e.HomeChainSel].Timelock,
			CallProxy: s.Chains[e.HomeChainSel].CallProxy,
		},
	}, []commoncs.ChangesetApplication{
		{
			Changeset: commoncs.WrapChangeSet(RemoveDONs),
			Config: RemoveDONsConfig{
				HomeChainSel: e.HomeChainSel,
				DonIDs:       []uint32{donsBefore[0].Id},
			},
		},
	})
	require.NoError(t, err)
	donsAfter, err := homeChain.CapabilityRegistry.GetDONs(nil)
	require.NoError(t, err)
	require.Len(t, donsAfter, len(donsBefore)-1)

	// Remove a don w/ MCMS
}
