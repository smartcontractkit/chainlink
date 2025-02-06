package test

import (
	"context"
	"errors"
	"sync"

	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	//"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/csa"

	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink/deployment"
)

var _ node.NodeServiceServer = (*JDNodeService)(nil)

type wrapper struct {
	deployment.Node
	enabled bool
}

func newWrapper(n deployment.Node) *wrapper {
	return &wrapper{
		Node:    n,
		enabled: true,
	}
}

func (w *wrapper) toNode() *nodev1.Node {
	return newJDNode(w.Node)
}

func wrapAll(m map[string]deployment.Node) map[string]*wrapper {
	w := make(map[string]*wrapper)
	for k, v := range m {
		w[k] = newWrapper(v)
	}
	return w
}

// JDNodeService is a mock implementation of the JobDistributor that supports
// the Node methods
type JDNodeService struct {
	mu    sync.RWMutex
	nodes map[string]*wrapper
	node.UnimplementedNodeServiceServer
}

func NewJDService(nodes map[string]deployment.Node) *JDNodeService {
	return &JDNodeService{
		nodes: wrapAll(nodes),
	}
}

func newWrapperFromJDNode(n *nodev1.Node) *wrapper {
}

// NewJDServiceFromListNodes initializes the service from a ListNodesResponse
func NewJDServiceFromListNodes(resp *nodev1.ListNodesResponse) *JDNodeService {
	nodes := make(map[string]deployment.Node)
	for _, node := range resp.Nodes {
		nodes[node.Id] = deployment.Node{
			// Populate other fields as needed
		}
	}
	return &JDNodeService{
		nodes: nodes,
	}
}

func (s *JDNodeService) GetNode(ctx context.Context, req *nodev1.GetNodeRequest) (*nodev1.GetNodeResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	node, ok := s.nodes[req.Id]
	if !ok {
		return nil, errors.New("node not found")
	}

	return &nodev1.GetNodeResponse{
		Node: newJDNode(node),
	}, nil
}

func (s *JDNodeService) ListNodes(ctx context.Context, req *nodev1.ListNodesRequest) (*nodev1.ListNodesResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var nodes []*nodev1.Node
	for _, node := range s.nodes {
		nodes = append(nodes, newJDNode(node))
	}

	return &nodev1.ListNodesResponse{
		Nodes: nodes,
	}, nil
}

func (s *JDNodeService) DisableNode(ctx context.Context, req *nodev1.DisableNodeRequest) (*nodev1.DisableNodeResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	node, ok := s.nodes[req.Id]
	if !ok {
		return nil, errors.New("node not found")
	}

	// Implement the logic to disable the node
	node.Disabled = true
	s.nodes[req.Id] = node

	return &nodev1.DisableNodeResponse{}, nil
}

func (s *JDNodeService) EnableNode(ctx context.Context, req *nodev1.EnableNodeRequest) (*nodev1.EnableNodeResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	node, ok := s.nodes[req.Id]
	if !ok {
		return nil, errors.New("node not found")
	}

	// Implement the logic to enable the node
	node.Disabled = false
	s.nodes[req.Id] = node

	return &nodev1.EnableNodeResponse{}, nil
}

func (s *JDNodeService) RegisterNode(ctx context.Context, req *nodev1.RegisterNodeRequest) (*nodev1.RegisterNodeResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.nodes[req.Node.Id]; exists {
		return nil, errors.New("node already exists")
	}

	// Implement the logic to register the node
	s.nodes[req.Node.Id] = deployment.Node{
		// Populate other fields as needed
	}

	return &nodev1.RegisterNodeResponse{}, nil
}

func (s *JDNodeService) UpdateNode(ctx context.Context, req *nodev1.UpdateNodeRequest) (*nodev1.UpdateNodeResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	node, ok := s.nodes[req.Node.Id]
	if !ok {
		return nil, errors.New("node not found")
	}

	// Implement the logic to update the node
	node.Name = req.Node.Name
	node.Labels = req.Node.Labels
	// Update other fields as needed
	s.nodes[req.Node.Id] = node

	return &nodev1.UpdateNodeResponse{}, nil
}

func newJDNode(n deployment.Node) *nodev1.Node {
	out := nodev1.Node{}
	out.Id = n.PeerID.String()
	out.Labels = n.Labels
	out.Name = n.Name
	out.PublicKey = n.CSAKey
	return &out
}

func newDeploymentNode(n *nodev1.Node) deployment.Node {
	return deployment.Node{
		NodeID: n.Id,
		Labels: n.Labels,
		Name:   n.Name,
		CSAKey: n.PublicKey,
	}
}
