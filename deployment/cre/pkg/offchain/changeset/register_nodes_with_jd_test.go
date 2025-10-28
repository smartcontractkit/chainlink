package changeset_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	cldf_offchain "github.com/smartcontractkit/chainlink-deployments-framework/offchain"
	"github.com/smartcontractkit/chainlink-deployments-framework/offchain/node"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"

	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain/changeset"
	operations2 "github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain/changeset/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/test"
)

func TestCsRegisterNodesWithJD_Apply(t *testing.T) {
	t.Parallel()

	zone := test.Zone

	t.Run("all ok", func(t *testing.T) {
		t.Parallel()
		env := test.SetupEnvV2(t, false)
		// Prepare input: one DON with two nodes
		input := changeset.CsRegisterNodesWithJDInput{
			Domain: "test-domain",
			DONs: []offchain.DONConfig{
				{
					Name: test.DONName,
					Nodes: []offchain.NodeCfg{
						{
							MinimalNodeCfg: node.MinimalNodeCfg{
								Name:   "node-1",
								CSAKey: "fake-csa-key-1",
							},
							P2PID: "fake-p2p-id",
							Zone:  zone,
						},
						{
							MinimalNodeCfg: node.MinimalNodeCfg{
								Name:   "node-2",
								CSAKey: "fake-csa-key-2",
							},
							P2PID: "fake-p2p-id",
							Zone:  zone,
						},
					},
				},
			},
		}
		cs := changeset.CsRegisterNodesWithJD{}

		// Apply changeset
		out, err := cs.Apply(*env.Env, input)
		require.NoError(t, err)

		// Validate output reports
		require.NotEmpty(t, out.Reports)
		for i, report := range out.Reports {
			// You can add more assertions here based on your report structure
			require.NotNil(t, report)
			// check that the report has the right node struct
			// need to cast it appropriately
			o := report.Output.(operations2.JDRegisterNodeOpOutput)
			assert.Equal(t, o.Node.Name, input.DONs[0].Nodes[i].Name)
			assert.Equal(t, o.Node.PublicKey, input.DONs[0].Nodes[i].CSAKey)
			checkLabels(t, o.Node.Labels, map[string]string{
				"product":             input.Domain,
				"environment":         env.Env.Name,
				"type":                "plugin",
				"zone":                zone,
				"don-" + test.DONName: "",
				"p2p_id":              "fake-p2p-id",
			})
		}
	})
	t.Run("register node error", func(t *testing.T) {
		t.Parallel()
		env := test.SetupEnvV2(t, false)

		env.Env.Offchain = testJDClient{
			env.TestJD,
		}

		// Prepare input: one DON with one node that will trigger an error
		input := changeset.CsRegisterNodesWithJDInput{
			Domain: "test-domain",
			DONs: []offchain.DONConfig{
				{
					Name: test.DONName,
					Nodes: []offchain.NodeCfg{
						{
							MinimalNodeCfg: node.MinimalNodeCfg{
								Name:   "register-node-error", // this name triggers the error in the test JD client
								CSAKey: "register-node-error-key-1",
							},
							Zone: zone,
						},
					},
				},
			},
		}

		// Apply changeset with test JD client that simulates error
		cs := changeset.CsRegisterNodesWithJD{}
		out, err := cs.Apply(*env.Env, input)
		require.Error(t, err)

		// Validate output reports
		require.NotEmpty(t, out.Reports)
		require.Len(t, out.Reports, 1)
		report := out.Reports[0]
		require.NotNil(t, report)
		o := report.Output.(operations2.JDRegisterNodeOpOutput)
		assert.NotEmpty(t, o.Error)
		assert.Contains(t, o.Error, "simulated register node error")
		assert.Empty(t, o.Node)
	})

	// test with mixed valid and error nodes
	t.Run("mixed valid and error nodes", func(t *testing.T) {
		t.Parallel()
		env := test.SetupEnvV2(t, false)
		nodesReps, err := env.Env.Offchain.ListNodes(t.Context(), &nodev1.ListNodesRequest{})
		require.NoError(t, err)
		nodes := nodesReps.GetNodes()

		env.Env.Offchain = testJDClient{
			env.TestJD,
		}
		// Prepare input: one DON with three nodes, one of which will trigger an error
		input := changeset.CsRegisterNodesWithJDInput{
			Domain: "cre",
			DONs: []offchain.DONConfig{
				{
					Name: test.DONName,
					Nodes: []offchain.NodeCfg{
						{
							MinimalNodeCfg: node.MinimalNodeCfg{
								Name:   nodes[0].Name,
								CSAKey: nodes[0].PublicKey,
							},
							Zone: zone,
						},
						{
							MinimalNodeCfg: node.MinimalNodeCfg{
								Name:   "register-node-error", // this name triggers the error in the test JD client
								CSAKey: "test-csa-key-1",
							},
						},
						{
							MinimalNodeCfg: node.MinimalNodeCfg{
								Name:   nodes[2].Name,
								CSAKey: nodes[2].PublicKey,
							},
							Zone: zone,
						},
					},
				},
			},
		}

		// Apply changeset with test JD client that simulates error
		cs := changeset.CsRegisterNodesWithJD{}
		out, err := cs.Apply(*env.Env, input)
		require.Error(t, err)

		// Validate output reports
		require.NotEmpty(t, out.Reports)
		require.Len(t, out.Reports, 3)

		// First node should be successful
		report1 := out.Reports[0]
		require.NotNil(t, report1)
		o1 := report1.Output.(operations2.JDRegisterNodeOpOutput)
		assert.Empty(t, o1.Error)
		require.NotNil(t, o1.Node)
		assert.Equal(t, o1.Node.Name, input.DONs[0].Nodes[0].Name)
		assert.Equal(t, o1.Node.PublicKey, input.DONs[0].Nodes[0].CSAKey)
		var p2pID string
		var typeLabel string
		for _, label := range nodes[0].Labels {
			if label.Key == "p2p_id" {
				p2pID = *label.Value
			}
			if label.Key == "type" {
				typeLabel = *label.Value
			}
		}
		checkLabels(t, o1.Node.Labels, map[string]string{
			"product":             input.Domain,
			"environment":         env.Env.Name,
			"type":                typeLabel,
			"zone":                zone,
			"don-" + test.DONName: test.DONName,
			"p2p_id":              p2pID,
		})

		// Second node should have an error
		report2 := out.Reports[1]
		require.NotNil(t, report2)
		o2 := report2.Output.(operations2.JDRegisterNodeOpOutput)
		assert.NotEmpty(t, o2.Error)
		assert.Contains(t, o2.Error, "simulated register node error")
		assert.Empty(t, o2.Node)

		for _, label := range nodes[2].Labels {
			if label.Key == "p2p_id" {
				p2pID = *label.Value
			}
			if label.Key == "type" {
				typeLabel = *label.Value
			}
		}

		// Third node should be successful
		report3 := out.Reports[2]
		require.NotNil(t, report3)
		o3 := report3.Output.(operations2.JDRegisterNodeOpOutput)
		assert.Empty(t, o3.Error)
		require.NotNil(t, o3.Node)
		assert.Equal(t, o3.Node.Name, input.DONs[0].Nodes[2].Name)
		assert.Equal(t, o3.Node.PublicKey, input.DONs[0].Nodes[2].CSAKey)
		checkLabels(t, o3.Node.Labels, map[string]string{
			"product":             input.Domain,
			"environment":         env.Env.Name,
			"type":                typeLabel,
			"zone":                zone,
			"don-" + test.DONName: test.DONName,
			"p2p_id":              p2pID,
		})
	})

	// test with valid node that gets updated
	t.Run("valid node that gets updated", func(t *testing.T) {
		t.Parallel()
		env := test.SetupEnvV2(t, false)

		nodes, err := env.Env.Offchain.ListNodes(t.Context(), &nodev1.ListNodesRequest{})
		require.NoError(t, err)
		firstNode := nodes.GetNodes()[0]

		n, err := env.Env.Offchain.GetNode(t.Context(), &nodev1.GetNodeRequest{
			PublicKey: &firstNode.PublicKey,
		})
		require.NoError(t, err)
		require.NotNil(t, n)
		require.NotNil(t, n.Node)
		require.Equal(t, n.Node.Name, firstNode.Name)
		require.Equal(t, n.Node.PublicKey, firstNode.PublicKey)

		// Prepare input: one DON with one node that will be registered and then updated
		input := changeset.CsRegisterNodesWithJDInput{
			Domain: "test-domain",
			DONs: []offchain.DONConfig{
				{
					Name: test.DONName,
					Nodes: []offchain.NodeCfg{
						{
							MinimalNodeCfg: node.MinimalNodeCfg{
								Name:   firstNode.Name,
								CSAKey: firstNode.PublicKey,
							},
							Zone: zone,
						},
					},
				},
			},
		}

		// Apply changeset with test JD client that simulates error
		cs := changeset.CsRegisterNodesWithJD{}
		out, err := cs.Apply(*env.Env, input)
		require.NoError(t, err)

		// Validate output reports
		require.Len(t, out.Reports, 1)
		require.NotEmpty(t, out.Reports)
		report := out.Reports[0]
		require.NotNil(t, report)
		o := report.Output.(operations2.JDRegisterNodeOpOutput)
		require.Empty(t, o.Error)
		require.NotNil(t, o.Node)
		assert.Equal(t, o.Node.Name, input.DONs[0].Nodes[0].Name)
		assert.Equal(t, o.Node.PublicKey, input.DONs[0].Nodes[0].CSAKey)
		var p2pID, typeLabel string
		for _, label := range firstNode.Labels {
			if label.Key == "p2p_id" {
				p2pID = *label.Value
			}
			if label.Key == "type" {
				typeLabel = *label.Value
			}
		}
		checkLabels(t, o.Node.Labels, map[string]string{
			"don-" + test.DONName: test.DONName, // the label already existed
			"product":             "cre",        // existing label should remain
			"type":                typeLabel,
			"environment":         env.Env.Name,
			"zone":                zone,
			"p2p_id":              p2pID,
		})
	})

	// test with update node error
	t.Run("update node error", func(t *testing.T) {
		t.Parallel()
		env := test.SetupEnvV2(t, false)

		nodes, err := env.Env.Offchain.ListNodes(t.Context(), &nodev1.ListNodesRequest{})
		require.NoError(t, err)
		firstNode := nodes.GetNodes()[0]

		n, err := env.Env.Offchain.GetNode(t.Context(), &nodev1.GetNodeRequest{
			PublicKey: &firstNode.PublicKey,
		})
		require.NoError(t, err)
		require.NotNil(t, n)
		require.NotNil(t, n.Node)
		require.Equal(t, n.Node.Name, firstNode.Name)
		require.Equal(t, n.Node.PublicKey, firstNode.PublicKey)

		_, err = env.Env.Offchain.UpdateNode(t.Context(), &nodev1.UpdateNodeRequest{
			Id:        firstNode.Id,
			Name:      "update-node-error",
			PublicKey: firstNode.PublicKey,
		})
		require.NoError(t, err)

		env.Env.Offchain = testJDClient{
			env.TestJD,
		}

		// Prepare input: one DON with one node that will trigger an update error
		input := changeset.CsRegisterNodesWithJDInput{
			Domain: "test-domain",
			DONs: []offchain.DONConfig{
				{
					Name: test.DONName,
					Nodes: []offchain.NodeCfg{
						{
							MinimalNodeCfg: node.MinimalNodeCfg{
								Name:   "update-node-error", // this name triggers the error in the test JD client
								CSAKey: firstNode.PublicKey,
							},
							Zone: zone,
						},
					},
				},
			},
		}

		// Apply changeset with test JD client that simulates error
		cs := changeset.CsRegisterNodesWithJD{}
		out, err := cs.Apply(*env.Env, input)
		require.Error(t, err)

		// Validate output reports
		require.NotEmpty(t, out.Reports)
		require.Len(t, out.Reports, 1)
		report := out.Reports[0]
		require.NotNil(t, report)
		o := report.Output.(operations2.JDRegisterNodeOpOutput)
		assert.NotEmpty(t, o.Error)
		assert.Contains(t, o.Error, "simulated update node error")
		assert.Empty(t, o.Node)
	})
}

func TestCsRegisterNodesWithJDV2_Apply(t *testing.T) {
	t.Parallel()

	t.Run("registers nodes for a DON", func(t *testing.T) {
		input := changeset.CsRegisterNodesWithJDInputV2{
			Domain: "cre",
			DONs: []offchain.DONConfig{
				{
					Name: test.DONName,
					Nodes: []offchain.NodeCfg{
						{
							MinimalNodeCfg: node.MinimalNodeCfg{Name: "node-1", CSAKey: "csa-key-1"},
							Zone:           test.Zone,
						},
						{
							MinimalNodeCfg: node.MinimalNodeCfg{Name: "node-2", CSAKey: "csa-key-2"},
							Zone:           test.Zone,
						},
					},
				},
			},
		}

		env := test.SetupEnvV2(t, false)
		cs := changeset.CsRegisterNodesWithJDV2{}

		out, err := cs.Apply(*env.Env, input)
		require.NoError(t, err)
		require.Len(t, out.Reports, 2)
		for i, report := range out.Reports {
			assert.NotNil(t, report)
			o := report.Output.(operations2.JDRegisterNodeOpOutput)
			assert.NotNil(t, o.Node)
			assert.Equal(t, input.DONs[0].Nodes[i].Name, o.Node.Name)
			assert.Equal(t, input.DONs[0].Nodes[i].CSAKey, o.Node.PublicKey)
			checkLabels(t, o.Node.Labels, map[string]string{
				"product":             input.Domain,
				"environment":         env.Env.Name,
				"type":                "plugin",
				"zone":                test.Zone,
				"don-" + test.DONName: "",
			})
		}
	})

	t.Run("fails with empty DONs", func(t *testing.T) {
		input := changeset.CsRegisterNodesWithJDInputV2{
			Domain: "cre",
			DONs:   []offchain.DONConfig{},
		}

		env := test.SetupEnvV2(t, false)
		cs := changeset.CsRegisterNodesWithJDV2{}
		err := cs.VerifyPreconditions(*env.Env, input)
		require.Error(t, err)
	})

	t.Run("fails with empty DON name", func(t *testing.T) {
		input := changeset.CsRegisterNodesWithJDInputV2{
			Domain: "cre",
			DONs: []offchain.DONConfig{
				{
					Name: "",
					Nodes: []offchain.NodeCfg{
						{
							MinimalNodeCfg: node.MinimalNodeCfg{Name: "node-1", CSAKey: "csa-key-1"},
							Zone:           test.Zone,
						},
					},
				},
			},
		}

		env := test.SetupEnvV2(t, false)
		cs := changeset.CsRegisterNodesWithJDV2{}
		err := cs.VerifyPreconditions(*env.Env, input)
		require.Error(t, err)
	})

	t.Run("fails with already registered node", func(t *testing.T) {
		t.Parallel()

		env := test.SetupEnvV2(t, false)

		nodesReps, err := env.Env.Offchain.ListNodes(t.Context(), &nodev1.ListNodesRequest{})
		require.NoError(t, err)
		nodes := nodesReps.GetNodes()

		input := changeset.CsRegisterNodesWithJDInputV2{
			Domain: "cre",
			DONs: []offchain.DONConfig{
				{
					Name: test.DONName,
					Nodes: []offchain.NodeCfg{
						{
							MinimalNodeCfg: node.MinimalNodeCfg{Name: nodes[0].Name, CSAKey: nodes[0].PublicKey},
							Zone:           test.Zone,
						},
					},
				},
			},
		}
		cs := changeset.CsRegisterNodesWithJDV2{}
		_, err = cs.Apply(*env.Env, input)
		require.Error(t, err)
		require.Contains(t, err.Error(), operations2.ErrNodeAlreadyExists.Error())
	})
}

func checkLabels(t *testing.T, labels []map[string]string, expected map[string]string) {
	t.Helper()
	assert.Len(t, labels, len(expected), "number of labels mismatch")
	for _, label := range labels {
		for k, v := range label {
			expectedV, ok := expected[k]
			assert.True(t, ok, "unexpected label key '%s'", k)
			assert.Equal(t, expectedV, v, "label value mismatch for key '%s'", k)
			delete(expected, k)
		}
	}
	assert.Empty(t, expected, "some expected labels were not found %v", expected)
}

// create and offchain client that overrides the RegisterNode method and UpdateNode method to simulate errors
type testJDClient struct {
	cldf_offchain.Client
}

func (t testJDClient) RegisterNode(ctx context.Context, in *nodev1.RegisterNodeRequest, opts ...grpc.CallOption) (*nodev1.RegisterNodeResponse, error) {
	if in.Name == "register-node-error" {
		return nil, errors.New("simulated register node error")
	}
	return t.Client.RegisterNode(ctx, in, opts...)
}

func (t testJDClient) UpdateNode(ctx context.Context, in *nodev1.UpdateNodeRequest, opts ...grpc.CallOption) (*nodev1.UpdateNodeResponse, error) {
	if in.Name == "update-node-error" {
		return nil, errors.New("simulated update node error")
	}
	return t.Client.UpdateNode(ctx, in, opts...)
}
