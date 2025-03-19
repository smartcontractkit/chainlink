package changeset

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-integrations/evm/testutils"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func TestRegisterNodesWithJD(t *testing.T) {
	t.Parallel()

	ctx := testutils.Context(t)
	e := testutil.NewMemoryEnv(t, false, 1)

	jobClient, ok := e.Offchain.(*memory.JobClient)
	require.True(t, ok, "expected Offchain to be of type *memory.JobClient")

	resp, err := jobClient.ListNodes(ctx, &nodev1.ListNodesRequest{})
	require.NoError(t, err)
	require.Lenf(t, resp.Nodes, 1, "expected exactly 1 node")
	require.Emptyf(t, jobClient.RegisteredNodes, "no registered nodes expected")

	csaKey := resp.Nodes[0].GetPublicKey()

	e, err = changeset.Apply(t, e, nil,
		changeset.Configure(
			deployment.CreateLegacyChangeSet(RegisterNodesWithJD),
			RegisterNodesInput{
				EnvLabel:    "test-env",
				ProductName: "test-product",
				DONsList: []DONConfig{
					{
						Name: "don1",
						BootstrapNodes: []NodeCfg{
							{Name: "node1", CSAKey: csaKey},
						},
					},
				},
			},
		),
	)
	require.NoError(t, err)
	require.Lenf(t, jobClient.RegisteredNodes, 1, "1 registered node expected")
	require.NotNilf(t, jobClient.RegisteredNodes[csaKey], "expected node with csa key %s to be registered", csaKey)
}

func TestRegisterNodesInput_Validate(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {
		cfg := RegisterNodesInput{
			EnvLabel:    "test-env",
			ProductName: "test-product",
			DONsList: []DONConfig{
				{
					Name: "MyDON",
					Nodes: []NodeCfg{
						{Name: "node1", CSAKey: "0xabc"},
					},
					BootstrapNodes: []NodeCfg{
						{Name: "bootstrap1", CSAKey: "0xdef"},
					},
				},
			},
		}
		err := cfg.Validate()
		require.NoError(t, err, "expected valid config to pass validation")
	})

	t.Run("missing product name", func(t *testing.T) {
		cfg := RegisterNodesInput{
			EnvLabel:    "test-env",
			ProductName: "",
			DONsList: []DONConfig{
				{
					Name: "AnotherDON",
					Nodes: []NodeCfg{
						{Name: "node1", CSAKey: "0xdef"},
					},
					BootstrapNodes: []NodeCfg{
						{Name: "node2", CSAKey: "0xabc"},
					},
				},
			},
		}
		err := cfg.Validate()
		require.Error(t, err, "expected an error when ProductName is empty")
	})

	t.Run("missing CSAKey", func(t *testing.T) {
		cfg := RegisterNodesInput{
			EnvLabel:    "test-env",
			ProductName: "test-product",
			DONsList: []DONConfig{
				{
					Name: "EmptyCSA",
					Nodes: []NodeCfg{
						{Name: "node1", CSAKey: ""},
					},
					BootstrapNodes: []NodeCfg{
						{Name: "bootstrap1", CSAKey: ""},
					},
				},
			},
		}
		err := cfg.Validate()
		require.Error(t, err, "expected an error when CSAKey is empty")
	})

	t.Run("missing BootstrapNode", func(t *testing.T) {
		cfg := RegisterNodesInput{
			EnvLabel:    "test-env",
			ProductName: "test-product",
			DONsList: []DONConfig{
				{
					Name: "EmptyCSA",
					Nodes: []NodeCfg{
						{Name: "node1", CSAKey: "0xaaa"},
					},
				},
			},
		}
		err := cfg.Validate()
		require.Error(t, err, "expected an error when BooststrapNodes is empty")
	})
}
