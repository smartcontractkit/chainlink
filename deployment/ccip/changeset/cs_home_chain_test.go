package changeset_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/common/view/v1_0"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestDeployHomeChain(t *testing.T) {
	t.Parallel()
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
	homeChainCfg := changeset.DeployHomeChainConfig{
		HomeChainSel:     homeChainSel,
		RMNStaticConfig:  testhelpers.NewTestRMNStaticConfig(),
		RMNDynamicConfig: testhelpers.NewTestRMNDynamicConfig(),
		NodeOperators:    testhelpers.NewTestNodeOperator(e.Chains[homeChainSel].DeployerKey.From),
		NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
			"NodeOperator": p2pIds,
		},
	}
	output, err := changeset.DeployHomeChainChangeset(e, homeChainCfg)
	require.NoError(t, err)
	require.NoError(t, e.ExistingAddresses.Merge(output.AddressBook))
	state, err := changeset.LoadOnchainState(e)
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
	require.Equal(t, []v1_0.NopView{
		{
			Admin: e.Chains[homeChainSel].DeployerKey.From,
			Name:  "NodeOperator",
		},
	}, capRegSnap.Nops)
	require.Len(t, capRegSnap.Nodes, len(p2pIds))
}

func TestRemoveDonsValidate(t *testing.T) {
	e, _ := testhelpers.NewMemoryEnvironment(t)
	s, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	homeChain := s.Chains[e.HomeChainSel]
	var tt = []struct {
		name      string
		config    changeset.RemoveDONsConfig
		expectErr bool
	}{
		{
			name: "invalid home",
			config: changeset.RemoveDONsConfig{
				HomeChainSel: 0,
				DonIDs:       []uint32{1},
			},
			expectErr: true,
		},
		{
			name: "invalid dons",
			config: changeset.RemoveDONsConfig{
				HomeChainSel: e.HomeChainSel,
				DonIDs:       []uint32{1377},
			},
			expectErr: true,
		},
		{
			name: "no dons",
			config: changeset.RemoveDONsConfig{
				HomeChainSel: e.HomeChainSel,
				DonIDs:       []uint32{},
			},
			expectErr: true,
		},
		{
			name: "success",
			config: changeset.RemoveDONsConfig{
				HomeChainSel: e.HomeChainSel,
				DonIDs:       []uint32{1},
			},
			expectErr: false,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.Validate(homeChain)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRemoveDons(t *testing.T) {
	e, _ := testhelpers.NewMemoryEnvironment(t)
	s, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	homeChain := s.Chains[e.HomeChainSel]

	// Remove a don w/o MCMS
	donsBefore, err := homeChain.CapabilityRegistry.GetDONs(nil)
	require.NoError(t, err)
	e.Env, err = commoncs.ApplyChangesets(t, e.Env, nil, []commoncs.ChangesetApplication{
		{
			Changeset: commoncs.WrapChangeSet(changeset.RemoveDONs),
			Config: changeset.RemoveDONsConfig{
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
	donsBefore, err = homeChain.CapabilityRegistry.GetDONs(nil)
	require.NoError(t, err)
	e.Env, err = commoncs.ApplyChangesets(t, e.Env, map[uint64]*proposalutils.TimelockExecutionContracts{
		e.HomeChainSel: {
			Timelock:  s.Chains[e.HomeChainSel].Timelock,
			CallProxy: s.Chains[e.HomeChainSel].CallProxy,
		},
	}, []commoncs.ChangesetApplication{
		{
			Changeset: commoncs.WrapChangeSet(commoncs.TransferToMCMSWithTimelock),
			Config: commoncs.TransferToMCMSWithTimelockConfig{
				ContractsByChain: map[uint64][]common.Address{
					e.HomeChainSel: {homeChain.CapabilityRegistry.Address()},
				},
				MinDelay: 0,
			},
		},
		{
			Changeset: commoncs.WrapChangeSet(changeset.RemoveDONs),
			Config: changeset.RemoveDONsConfig{
				HomeChainSel: e.HomeChainSel,
				DonIDs:       []uint32{donsBefore[0].Id},
				MCMS:         &changeset.MCMSConfig{MinDelay: 0},
			},
		},
	})
	require.NoError(t, err)
	donsAfter, err = homeChain.CapabilityRegistry.GetDONs(nil)
	require.NoError(t, err)
	require.Len(t, donsAfter, len(donsBefore)-1)
}
