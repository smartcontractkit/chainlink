package view

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/test"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

func TestGenerateNopsView(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	// Create 3 node IDs
	nodeIDs := []string{"node1", "node2", "node3"}

	// Set up mock nodes with different configurations
	p2pIDs := []string{}
	csaKeys := []string{}
	deploymentNodes := []deployment.Node{}

	for i, id := range nodeIDs {
		// Create unique P2P IDs and CSA keys
		p2pKey := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(int64(i)))
		p2pIDs = append(p2pIDs, p2pKey.ID())
		csaKey := "csa_key_" + id
		csaKeys = append(csaKeys, csaKey)

		// Create a node
		node := deployment.Node{
			NodeID:      id,
			Name:        "Node " + id,
			PeerID:      p2pKey.PeerID(),
			IsBootstrap: i == 0, // Make the first node a bootstrap node
			AdminAddr:   "0x" + id,
			CSAKey:      csaKey,
			WorkflowKey: "workflow_" + id,
			Labels: []*ptypes.Label{
				&ptypes.Label{
					Key:   "role",
					Value: ptr("tester")},
				&ptypes.Label{
					Key:   "p2p",
					Value: ptr(p2pIDs[i])},
			},
		}
		deploymentNodes = append(deploymentNodes, node)
	}

	// Create mock JD service
	jdService := test.NewJDService(deploymentNodes)

	t.Run("successful view generation", func(t *testing.T) {
		// Generate view
		nopsView, err := GenerateNopsView(lggr, nodeIDs, jdService)
		require.NoError(t, err)

		// Check that we have all 3 nodes in the view
		require.Len(t, nopsView, 3)

		// Check each node's properties
		for i, id := range nodeIDs {
			nodeName := "Node " + id
			node, exists := nopsView[nodeName]
			require.True(t, exists, "Node %s should exist in the view", nodeName)

			assert.Equal(t, id, node.NodeID)

			assert.Equal(t, csaKeys[i], node.CSAKey)
			assert.Equal(t, "workflow_"+id, node.WorkflowKey)

			// Check labels
			require.Len(t, node.Labels, 2)
			assertLabelExists(t, node.Labels, "role", "tester")
			assertLabelExists(t, node.Labels, "p2p", p2pIDs[i])

			// Empty jobspecs is expected as our mock returns empty responses
			assert.Empty(t, node.ApprovedJobspecs)
		}
	})

	t.Run("node not found in JD", func(t *testing.T) {
		v, err := GenerateNopsView(lggr, []string{"unknown_node"}, jdService)
		require.NoError(t, err)
		assert.Empty(t, v)
	})

	t.Run("error from ListNodes", func(t *testing.T) {
		// Create a custom JD service that returns an error for ListNodes
		errorJDService := &customJDService{
			NodeServiceClient: jdService,
			listNodesError:    errors.New("failed to list nodes from JD"),
		}
		// Should return the error from ListNodes
		_, err := GenerateNopsView(lggr, nodeIDs, errorJDService)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to list nodes from JD")
	})
}

// Helper function to check if a label with expected key/value exists
func assertLabelExists(t *testing.T, labels []LabelView, key, expectedValue string) {
	t.Helper()
	for _, label := range labels {
		if label.Key == key {
			require.NotNil(t, label.Value)
			assert.Equal(t, expectedValue, *label.Value)
			return
		}
	}
	t.Errorf("Label with key %s not found", key)
}

// Custom JD service implementation for error testing
type customJDService struct {
	nodev1.NodeServiceClient
	listNodesError error
	*test.UnimplementedCSAServiceClient
	*test.UnimplementedJobServiceClient
}

func (s *customJDService) ListNodes(ctx context.Context, req *nodev1.ListNodesRequest, opts ...grpc.CallOption) (*nodev1.ListNodesResponse, error) {
	if s.listNodesError != nil {
		return nil, s.listNodesError
	}
	return s.NodeServiceClient.ListNodes(ctx, req, opts...)
}

func ptr[T any](t T) *T {
	return &t
}
