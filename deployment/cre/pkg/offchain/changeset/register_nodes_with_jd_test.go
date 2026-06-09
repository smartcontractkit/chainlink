package changeset_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"
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

		h := test.NewTestHarness(t)

		// Prepare input: one DON with two nodes
		input := changeset.CsRegisterNodesWithJDInput{
			Domain:      "test-domain",
			Environment: test.EnvironmentName,
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

		// Apply changeset
		task := runtime.ChangesetTask(changeset.CsRegisterNodesWithJD{}, input)
		err := h.Runtime.Exec(task)
		require.NoError(t, err)
		out := h.Runtime.State().Outputs[task.ID()]
		require.NotNil(t, out, "changeset output should not be nil")

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
				"environment":         h.Runtime.Environment().Name,
				"type":                "plugin",
				"zone":                zone,
				"don-" + test.DONName: "",
				"p2p_id":              "fake-p2p-id",
			})
		}
	})

	t.Run("register node error", func(t *testing.T) {
		t.Parallel()

		var (
			h   = test.NewTestHarness(t)
			env = h.Runtime.Environment()
		)

		env.Offchain = testJDClient{
			h.TestJD,
		}

		// Prepare input: one DON with one node that will trigger an error
		input := changeset.CsRegisterNodesWithJDInput{
			Domain:      "test-domain",
			Environment: test.EnvironmentName,
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
		//
		// We cannot use the runtime Exec method here the runtime does not store output for failed changesets executions
		// and the test is asserting on the output.
		cs := changeset.CsRegisterNodesWithJD{}
		out, err := cs.Apply(env, input)
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

		var (
			h   = test.NewTestHarness(t)
			env = h.Runtime.Environment()
		)

		nodesReps, err := env.Offchain.ListNodes(t.Context(), &nodev1.ListNodesRequest{})
		require.NoError(t, err)
		nodes := nodesReps.GetNodes()

		env.Offchain = testJDClient{
			h.TestJD,
		}

		// Prepare input: one DON with three nodes, one of which will trigger an error
		input := changeset.CsRegisterNodesWithJDInput{
			Domain:      "cre",
			Environment: test.EnvironmentName,
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
		out, err := cs.Apply(env, input)
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
			"environment":         h.Runtime.Environment().Name,
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
			"environment":         h.Runtime.Environment().Name,
			"type":                typeLabel,
			"zone":                zone,
			"don-" + test.DONName: test.DONName,
			"p2p_id":              p2pID,
		})
	})

	// test with valid node that gets updated
	t.Run("valid node that gets updated", func(t *testing.T) {
		t.Parallel()

		h := test.NewTestHarness(t)

		nodes, err := h.Runtime.Environment().Offchain.ListNodes(t.Context(), &nodev1.ListNodesRequest{})
		require.NoError(t, err)
		firstNode := nodes.GetNodes()[0]

		n, err := h.Runtime.Environment().Offchain.GetNode(t.Context(), &nodev1.GetNodeRequest{
			PublicKey: &firstNode.PublicKey,
		})
		require.NoError(t, err)
		require.NotNil(t, n)
		require.NotNil(t, n.Node)
		require.Equal(t, n.Node.Name, firstNode.Name)
		require.Equal(t, n.Node.PublicKey, firstNode.PublicKey)

		// Prepare input: one DON with one node that will be registered and then updated
		input := changeset.CsRegisterNodesWithJDInput{
			Domain:      "test-domain",
			Environment: test.EnvironmentName,
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
		task := runtime.ChangesetTask(changeset.CsRegisterNodesWithJD{}, input)
		err = h.Runtime.Exec(task)
		require.NoError(t, err)
		out := h.Runtime.State().Outputs[task.ID()]
		require.NotNil(t, out, "changeset output should not be nil")

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
			"environment":         h.Runtime.Environment().Name,
			"zone":                zone,
			"p2p_id":              p2pID,
		})
	})

	// test with update node error
	t.Run("update node error", func(t *testing.T) {
		t.Parallel()

		var (
			h   = test.NewTestHarness(t)
			env = h.Runtime.Environment()
		)

		nodes, err := h.Runtime.Environment().Offchain.ListNodes(t.Context(), &nodev1.ListNodesRequest{})
		require.NoError(t, err)
		firstNode := nodes.GetNodes()[0]

		n, err := env.Offchain.GetNode(t.Context(), &nodev1.GetNodeRequest{
			PublicKey: &firstNode.PublicKey,
		})
		require.NoError(t, err)
		require.NotNil(t, n)
		require.NotNil(t, n.Node)
		require.Equal(t, n.Node.Name, firstNode.Name)
		require.Equal(t, n.Node.PublicKey, firstNode.PublicKey)

		_, err = env.Offchain.UpdateNode(t.Context(), &nodev1.UpdateNodeRequest{
			Id:        firstNode.Id,
			Name:      "update-node-error",
			PublicKey: firstNode.PublicKey,
		})
		require.NoError(t, err)

		env.Offchain = testJDClient{
			h.TestJD,
		}

		// Prepare input: one DON with one node that will trigger an update error
		input := changeset.CsRegisterNodesWithJDInput{
			Domain:      "test-domain",
			Environment: test.EnvironmentName,
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
		out, err := cs.Apply(env, input)
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
			Domain:      "cre",
			Environment: test.EnvironmentName,
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

		h := test.NewTestHarness(t)

		task := runtime.ChangesetTask(changeset.CsRegisterNodesWithJDV2{}, input)
		err := h.Runtime.Exec(task)
		require.NoError(t, err)

		out := h.Runtime.State().Outputs[task.ID()]
		require.NotNil(t, out, "changeset output should not be nil")
		require.Len(t, out.Reports, 2)
		for i, report := range out.Reports {
			assert.NotNil(t, report)
			o := report.Output.(operations2.JDRegisterNodeOpOutput)
			assert.NotNil(t, o.Node)
			assert.Equal(t, input.DONs[0].Nodes[i].Name, o.Node.Name)
			assert.Equal(t, input.DONs[0].Nodes[i].CSAKey, o.Node.PublicKey)
			checkLabels(t, o.Node.Labels, map[string]string{
				"product":             input.Domain,
				"environment":         h.Runtime.Environment().Name,
				"type":                "plugin",
				"zone":                test.Zone,
				"don-" + test.DONName: "",
			})
		}
	})

	t.Run("fails with empty DONs", func(t *testing.T) {
		input := changeset.CsRegisterNodesWithJDInputV2{
			Domain:      "cre",
			Environment: test.EnvironmentName,
			DONs:        []offchain.DONConfig{},
		}

		h := test.NewTestHarness(t)
		cs := changeset.CsRegisterNodesWithJDV2{}
		err := cs.VerifyPreconditions(h.Runtime.Environment(), input)
		require.Error(t, err)
	})

	t.Run("fails with empty DON name", func(t *testing.T) {
		input := changeset.CsRegisterNodesWithJDInputV2{
			Domain:      "cre",
			Environment: test.EnvironmentName,
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

		h := test.NewTestHarness(t)
		cs := changeset.CsRegisterNodesWithJDV2{}
		err := cs.VerifyPreconditions(h.Runtime.Environment(), input)
		require.Error(t, err)
	})

	t.Run("fails with already registered node", func(t *testing.T) {
		t.Parallel()

		h := test.NewTestHarness(t)

		nodesReps, err := h.Runtime.Environment().Offchain.ListNodes(t.Context(), &nodev1.ListNodesRequest{})
		require.NoError(t, err)
		nodes := nodesReps.GetNodes()

		input := changeset.CsRegisterNodesWithJDInputV2{
			Domain:      "cre",
			Environment: test.EnvironmentName,
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

		task := runtime.ChangesetTask(changeset.CsRegisterNodesWithJDV2{}, input)
		err = h.Runtime.Exec(task)
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
