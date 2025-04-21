package changeset

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func TestRegisterNodesWithJD(t *testing.T) {
	t.Parallel()

	ctx := testutils.Context(t)
	eWrapper := testutil.NewMemoryEnvV2(t, testutil.MemoryEnvConfig{
		NumNodes: 1,
	})

	e := eWrapper.Environment

	jobClient, ok := e.Offchain.(*memory.JobClient)
	require.True(t, ok, "expected Offchain to be of type *memory.JobClient")

	resp, err := jobClient.ListNodes(ctx, &nodev1.ListNodesRequest{})
	require.NoError(t, err)
	require.Lenf(t, resp.Nodes, 1, "expected exactly 1 node")
	require.Emptyf(t, jobClient.RegisteredNodes, "no registered nodes expected")

	csaKey := resp.Nodes[0].GetPublicKey()

	require.NoError(t, err)

	e, _, err = changeset.ApplyChangesetsV2(t, e,
		[]changeset.ConfiguredChangeSet{
			changeset.Configure(
				RegisterNodesWithJDChangeset,
				RegisterNodesInput{
					BaseLabels: map[string]string{
						"environment": "test-env",
						"product":     "test-product",
					},
					DONsList: []DONConfig{
						{
							ID:   1,
							Name: "don1",
							BootstrapNodes: []NodeCfg{
								{Name: "node1", CSAKey: csaKey},
							},
						},
					},
				},
			),
		},
	)

	require.NoError(t, err)
	require.Lenf(t, jobClient.RegisteredNodes, 1, "1 registered node expected")
	require.NotNilf(t, jobClient.RegisteredNodes[csaKey], "expected node with csa key %s to be registered", csaKey)
}

func TestRegisterNodesInput_Validate(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {
		cfg := RegisterNodesInput{
			BaseLabels: map[string]string{
				"environment": "test-env",
				"product":     "test-product",
			},
			DONsList: []DONConfig{
				{
					ID:   1,
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

	t.Run("empty map value is allowed", func(t *testing.T) {
		// empty map value can indicate the "presence" of something
		cfg := RegisterNodesInput{
			BaseLabels: map[string]string{
				"environment": "test-env",
				"product":     "",
			},
			DONsList: []DONConfig{
				{
					ID:   1,
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
		require.NoError(t, err)
	})

	t.Run("empty map key", func(t *testing.T) {
		cfg := RegisterNodesInput{
			BaseLabels: map[string]string{
				"environment": "test-env",
				"":            "test-product",
			},
			DONsList: []DONConfig{
				{
					ID:   1,
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
			BaseLabels: map[string]string{
				"environment": "test-env",
				"product":     "test-product",
			},
			DONsList: []DONConfig{
				{
					ID:   1,
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
			BaseLabels: map[string]string{
				"environment": "test-env",
				"product":     "test-product",
			},
			DONsList: []DONConfig{
				{
					ID:   1,
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
