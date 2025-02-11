package changeset

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestRegisterNodesWithJD(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{Chains: 1, Nodes: 1})

	nodeP2pKey := e.NodeIDs[0]

	jobClient, ok := e.Offchain.(*memory.JobClient)

	require.True(t, ok, "expected Offchain to be of type *memory.JobClient")
	require.Lenf(t, jobClient.Nodes, 1, "expected exactly 1 node")
	require.Emptyf(t, jobClient.RegisteredNodes, "no registered nodes expected")

	csaKey := jobClient.Nodes[nodeP2pKey].Keys.CSA.PublicKeyString()

	e, err := changeset.ApplyChangesets(t, e, nil, []changeset.ChangesetApplication{
		{
			Changeset: changeset.WrapChangeSet(RegisterNodesWithJD),
			Config: RegisterNodesInput{
				EnvLabel:    "test-env",
				ProductName: "test-product",
				DONs: DONConfigMap{
					"don1": {
						Name: "don1",
						Nodes: []NodeCfg{
							{Name: "node1", CSAKey: csaKey},
						},
					},
				},
			},
		},
	})
	require.NoError(t, err)
	require.Lenf(t, jobClient.RegisteredNodes, 1, "1 registered node expected")
	require.NotNilf(t, jobClient.RegisteredNodes[csaKey], "expected node with csa key %s to be registered", csaKey)
}

func TestRegisterNodesInput_Validate(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {
		cfg := RegisterNodesInput{
			EnvLabel:    "test-env",
			ProductName: "test-product",
			DONs: DONConfigMap{
				"don1": {
					Name: "MyDON",
					Nodes: []NodeCfg{
						{Name: "node1", CSAKey: "0xabc", IsBootstrap: false},
					},
				},
			},
		}
		err := cfg.Validate()
		require.NoError(t, err, "expected valid config to pass validation")
	})

	t.Run("missing product name", func(t *testing.T) {
		cfg := RegisterNodesInput{
			EnvLabel: "test-env",
			DONs: DONConfigMap{
				"don2": {
					Name: "AnotherDON",
					Nodes: []NodeCfg{
						{Name: "node1", CSAKey: "0xdef"},
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
			DONs: DONConfigMap{
				"don3": {
					Name: "EmptyCSA",
					Nodes: []NodeCfg{
						{Name: "node1", CSAKey: "", IsBootstrap: true},
					},
				},
			},
		}
		err := cfg.Validate()
		require.Error(t, err, "expected an error when CSAKey is empty")
	})
}
