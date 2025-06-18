package test

import (
	"context"
	"fmt"
	"sync"

	"google.golang.org/grpc"

	csav1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/csa"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"

	"github.com/smartcontractkit/chainlink/deployment"
	//	"github.com/smartcontractkit/chainlink/deployment/environment/memory"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

var _ nodev1.NodeServiceClient = (*JDNodeService)(nil)
var _ jobv1.JobServiceClient = (*JDNodeService)(nil)
var _ csav1.CSAServiceClient = (*JDNodeService)(nil)

// JDNodeService is a mock implementation of the JobDistributor that supports
// the Node methods
type JDNodeService struct {
	mu    sync.RWMutex
	store *store
	*UnimplementedJobServiceClient
	*UnimplementedCSAServiceClient
}

func NewJDService(nodes []deployment.Node) *JDNodeService {
	return &JDNodeService{
		//store: wrapAll(nodes),
		store: newStore(nodes),
	}
}

// NewJDServiceFromListNodes initializes the service from a ListNodesResponse
func NewJDServiceFromListNodes(resp *nodev1.ListNodesResponse) (*JDNodeService, error) {
	var nodes []deployment.Node
	for _, jdNodes := range resp.Nodes {
		n, err := newDeploymentNode(jdNodes)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, n)
	}
	return &JDNodeService{
		store: newStore(nodes),
	}, nil
}

func (s *JDNodeService) GetNode(ctx context.Context, req *nodev1.GetNodeRequest, opts ...grpc.CallOption) (*nodev1.GetNodeResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	w, err := s.store.getNode(req.Id)
	if err != nil {
		return nil, err
	}

	return &nodev1.GetNodeResponse{
		Node: newJDNode(w.Node),
	}, nil
}

func (s *JDNodeService) ListNodes(ctx context.Context, req *nodev1.ListNodesRequest, opts ...grpc.CallOption) (*nodev1.ListNodesResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var nodes []*nodev1.Node
	for _, w := range s.store.list() {
		n := newJDNode(w.Node)
		if ApplyNodeFilter(req.Filter, n) {
			nodes = append(nodes, n)
		}
	}

	return &nodev1.ListNodesResponse{
		Nodes: nodes,
	}, nil
}

func (s *JDNodeService) DisableNode(ctx context.Context, req *nodev1.DisableNodeRequest, opts ...grpc.CallOption) (*nodev1.DisableNodeResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	node, err := s.store.getNode(req.Id)
	if err != nil {
		return nil, err
	}

	node.enabled = false
	s.store.put(node)

	return &nodev1.DisableNodeResponse{}, nil
}

func (s *JDNodeService) EnableNode(ctx context.Context, req *nodev1.EnableNodeRequest, opts ...grpc.CallOption) (*nodev1.EnableNodeResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	node, err := s.store.getNode(req.Id)
	if err != nil {
		return nil, err
	}

	// Implement the logic to enable the node
	node.enabled = true
	s.store.put(node)

	return &nodev1.EnableNodeResponse{}, nil
}

func (s *JDNodeService) RegisterNode(ctx context.Context, req *nodev1.RegisterNodeRequest, opts ...grpc.CallOption) (*nodev1.RegisterNodeResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	n, _ := s.store.getNodeByCSA(req.PublicKey)
	if n != nil {
		return nil, fmt.Errorf("node already registered with CSA key %s", req.PublicKey)
	}

	w, err := newWrapperFromRegister(req)
	if err != nil {
		return nil, err
	}
	s.store.put(w)

	return &nodev1.RegisterNodeResponse{}, nil
}

func (s *JDNodeService) ListNodeChainConfigs(ctx context.Context, req *nodev1.ListNodeChainConfigsRequest, opts ...grpc.CallOption) (*nodev1.ListNodeChainConfigsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// the chain config filter is a strict subset of the node filter
	var filter *nodev1.ListNodesRequest_Filter
	if req.Filter != nil {
		filter = &nodev1.ListNodesRequest_Filter{
			Ids: req.Filter.NodeIds,
		}
	}
	var out []*nodev1.ChainConfig
	for _, w := range s.store.list() {
		if ApplyNodeFilter(filter, w.toJDNode()) {
			cc, err := w.Node.ChainConfigs()
			if err != nil {
				return nil, err
			}
			out = append(out, cc...)
		}
	}
	return &nodev1.ListNodeChainConfigsResponse{
		ChainConfigs: out,
	}, nil
}

func newWrapperFromRegister(req *nodev1.RegisterNodeRequest) (*wrappedNode, error) {
	return nil, nil
}

func (s *JDNodeService) UpdateNode(ctx context.Context, req *nodev1.UpdateNodeRequest, opts ...grpc.CallOption) (*nodev1.UpdateNodeResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.store.getNodeByP2P(p2pKey(req.Id))
	if err != nil {
		return nil, fmt.Errorf("node not found for p2p %s", req.Id)
	}

	w, err := newWrapperFromUpdate(req)
	if err != nil {
		return nil, err
	}

	s.store.put(w)
	return &nodev1.UpdateNodeResponse{}, nil
}

func newWrapperFromUpdate(req *nodev1.UpdateNodeRequest) (*wrappedNode, error) {
	return nil, nil
}

func newJDNode(n deployment.Node) *nodev1.Node {
	out := nodev1.Node{
		Id:          n.NodeID,
		Labels:      n.Labels,
		Name:        n.Name,
		PublicKey:   n.CSAKey,
		WorkflowKey: &n.WorkflowKey,
	}

	return &out
}

func newDeploymentNode(n *nodev1.Node) (deployment.Node, error) {
	out := deployment.Node{
		NodeID:      n.Id,
		Labels:      n.Labels,
		Name:        n.Name,
		CSAKey:      n.PublicKey,
		WorkflowKey: n.GetWorkflowKey(),
	}
	for _, label := range n.Labels {
		if p, err := p2pkey.MakePeerID(*label.Value); err == nil {
			out.PeerID = p
		}
	}
	return out, nil
}

// wrappedNode is a wrapper around deployment.Node that adds some state
type wrappedNode struct {
	deployment.Node
	enabled bool
}

func newWrapper(n deployment.Node) *wrappedNode {
	return &wrappedNode{
		Node:    n,
		enabled: true,
	}
}

func (w *wrappedNode) toJDNode() *nodev1.Node {
	return newJDNode(w.Node)
}

// p2pKey is a wrapper around string to make it easier to read
type p2pKey string

func (p p2pKey) String() string {
	return string(p)
}
func (p p2pKey) Validate() error {
	_, err := p2pkey.MakePeerID(p.String())
	return err
}

// csaKey is a wrapper around string to make it easier to read
type csaKey = string

// store is a thread-safe store for wrappedNode
// it is indexed by both p2p key and csa key
type store struct {
	mu  sync.RWMutex
	db2 map[string]*wrappedNode

	p2pToID map[p2pKey]string
	csaToID map[csaKey]string
}

func newStore(node []deployment.Node) *store {
	s := &store{
		db2:     make(map[string]*wrappedNode),
		csaToID: make(map[csaKey]string),
		p2pToID: make(map[p2pKey]string),
	}
	for _, v := range node {
		w := newWrapper(v)
		s.db2[v.NodeID] = w
		s.p2pToID[p2pKey(w.Node.PeerID.String())] = v.NodeID
		s.csaToID[w.Node.CSAKey] = v.NodeID
	}
	return s
}

func (s *store) getNode(id string) (*wrappedNode, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n, ok := s.db2[id]
	if !ok {
		return nil, fmt.Errorf("node not found for id %s", id)
	}
	return n, nil
}

func (s *store) getNodeByP2P(p2p p2pKey) (*wrappedNode, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	id, ok := s.p2pToID[p2p]
	if !ok {
		return nil, fmt.Errorf("node not found for p2p %s", p2p)
	}
	return s.getNode(id)
}

func (s *store) getNodeByCSA(csa csaKey) (*wrappedNode, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	id, ok := s.csaToID[csa]
	if !ok {
		return nil, fmt.Errorf("node not found for csa key %s", csa)
	}
	return s.getNode(id)
}

func (s *store) list() []*wrappedNode {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []*wrappedNode
	for _, v := range s.db2 {
		out = append(out, v)
	}
	return out
}

func (s *store) put(n *wrappedNode) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.db2[n.Node.NodeID] = n
	s.csaToID[n.Node.CSAKey] = n.NodeID
	s.p2pToID[p2pKey(n.Node.PeerID.String())] = n.NodeID
}

type UnimplementedJobServiceClient struct{}

func (s *UnimplementedJobServiceClient) BatchProposeJob(ctx context.Context, in *jobv1.BatchProposeJobRequest, opts ...grpc.CallOption) (*jobv1.BatchProposeJobResponse, error) {
	// TODO CCIP-3108  implement me
	panic("implement me")
}

func (s *UnimplementedJobServiceClient) DeleteJob(ctx context.Context, in *jobv1.DeleteJobRequest, opts ...grpc.CallOption) (*jobv1.DeleteJobResponse, error) {
	panic("unimplemented")
}

func (s *UnimplementedJobServiceClient) UpdateJob(ctx context.Context, in *jobv1.UpdateJobRequest, opts ...grpc.CallOption) (*jobv1.UpdateJobResponse, error) {
	panic("unimplemented")
}

// GetJob implements job.JobServiceClient.
func (s *UnimplementedJobServiceClient) GetJob(ctx context.Context, in *jobv1.GetJobRequest, opts ...grpc.CallOption) (*jobv1.GetJobResponse, error) {
	panic("unimplemented")
}

// GetProposal implements job.JobServiceClient.
func (s *UnimplementedJobServiceClient) GetProposal(ctx context.Context, in *jobv1.GetProposalRequest, opts ...grpc.CallOption) (*jobv1.GetProposalResponse, error) {
	panic("unimplemented")
}

// ListJobs implements job.JobServiceClient.
func (s *UnimplementedJobServiceClient) ListJobs(ctx context.Context, in *jobv1.ListJobsRequest, opts ...grpc.CallOption) (*jobv1.ListJobsResponse, error) {
	// returns blank response to ensure the nops view can be generated using the mock client
	return &jobv1.ListJobsResponse{}, nil
}

// ListProposals implements job.JobServiceClient.
func (s *UnimplementedJobServiceClient) ListProposals(ctx context.Context, in *jobv1.ListProposalsRequest, opts ...grpc.CallOption) (*jobv1.ListProposalsResponse, error) {
	panic("unimplemented")
}

// ProposeJob implements job.JobServiceClient.
func (s *UnimplementedJobServiceClient) ProposeJob(ctx context.Context, in *jobv1.ProposeJobRequest, opts ...grpc.CallOption) (*jobv1.ProposeJobResponse, error) {
	panic("unimplemented")
}

// RevokeJob implements job.JobServiceClient.
func (s *UnimplementedJobServiceClient) RevokeJob(ctx context.Context, in *jobv1.RevokeJobRequest, opts ...grpc.CallOption) (*jobv1.RevokeJobResponse, error) {
	panic("unimplemented")
}

type UnimplementedCSAServiceClient struct{}

func (s *UnimplementedCSAServiceClient) GetKeypair(ctx context.Context, in *csav1.GetKeypairRequest, opts ...grpc.CallOption) (*csav1.GetKeypairResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s *UnimplementedCSAServiceClient) ListKeypairs(ctx context.Context, in *csav1.ListKeypairsRequest, opts ...grpc.CallOption) (*csav1.ListKeypairsResponse, error) {
	panic("unimplemented")
}
