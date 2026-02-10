package mockjd

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	csav1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/csa"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
)

type Server struct {
	grpcServer *grpc.Server
	listener   net.Listener
	service    *MockService
}

func NewServer(csaKey string) (*Server, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("failed to create listener: %w", err)
	}

	grpcServer := grpc.NewServer()
	service := &MockService{
		csaKey:       csaKey,
		nodes:        make(map[string]*nodeData),
		chainConfigs: make(map[string][]*nodev1.ChainConfig),
		proposals:    make(map[string]*jobv1.Proposal),
	}

	nodev1.RegisterNodeServiceServer(grpcServer, service)
	jobv1.RegisterJobServiceServer(grpcServer, service)
	csav1.RegisterCSAServiceServer(grpcServer, service)

	return &Server{
		grpcServer: grpcServer,
		listener:   listener,
		service:    service,
	}, nil
}

func (s *Server) Start() error {
	go func() {
		_ = s.grpcServer.Serve(s.listener)
	}()
	return nil
}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
}

func (s *Server) Addr() string {
	return s.listener.Addr().String()
}

type nodeData struct {
	id        string
	name      string
	publicKey string
	labels    []*ptypes.Label
	connected bool
}

type MockService struct {
	nodev1.UnimplementedNodeServiceServer
	jobv1.UnimplementedJobServiceServer
	csav1.UnimplementedCSAServiceServer

	mu           sync.RWMutex
	csaKey       string
	nodes        map[string]*nodeData
	chainConfigs map[string][]*nodev1.ChainConfig
	proposals    map[string]*jobv1.Proposal
	proposalSeq  int
}

func (m *MockService) ListKeypairs(ctx context.Context, req *csav1.ListKeypairsRequest) (*csav1.ListKeypairsResponse, error) {
	return &csav1.ListKeypairsResponse{
		Keypairs: []*csav1.Keypair{
			{PublicKey: m.csaKey},
		},
	}, nil
}

func (m *MockService) GetKeypair(ctx context.Context, req *csav1.GetKeypairRequest) (*csav1.GetKeypairResponse, error) {
	return &csav1.GetKeypairResponse{
		Keypair: &csav1.Keypair{PublicKey: m.csaKey},
	}, nil
}

func (m *MockService) RegisterNode(ctx context.Context, req *nodev1.RegisterNodeRequest) (*nodev1.RegisterNodeResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, n := range m.nodes {
		if n.publicKey == req.PublicKey {
			return &nodev1.RegisterNodeResponse{
				Node: m.toProtoNode(n),
			}, nil
		}
	}

	id := uuid.New().String()
	node := &nodeData{
		id:        id,
		name:      req.Name,
		publicKey: req.PublicKey,
		labels:    convertLabels(req.Labels),
		connected: true,
	}
	m.nodes[id] = node

	return &nodev1.RegisterNodeResponse{
		Node: m.toProtoNode(node),
	}, nil
}

func (m *MockService) UpdateNode(ctx context.Context, req *nodev1.UpdateNodeRequest) (*nodev1.UpdateNodeResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	node, ok := m.nodes[req.Id]
	if !ok {
		for _, n := range m.nodes {
			if n.publicKey == req.PublicKey {
				node = n
				break
			}
		}
	}
	if node == nil {
		return nil, fmt.Errorf("node not found: %s", req.Id)
	}

	node.name = req.Name
	node.labels = convertLabels(req.Labels)

	return &nodev1.UpdateNodeResponse{
		Node: m.toProtoNode(node),
	}, nil
}

func (m *MockService) GetNode(ctx context.Context, req *nodev1.GetNodeRequest) (*nodev1.GetNodeResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if req.PublicKey != nil {
		for _, n := range m.nodes {
			if n.publicKey == *req.PublicKey {
				return &nodev1.GetNodeResponse{Node: m.toProtoNode(n)}, nil
			}
		}
		return nil, fmt.Errorf("node not found for public key")
	}

	node, ok := m.nodes[req.Id]
	if !ok {
		return nil, fmt.Errorf("node not found: %s", req.Id)
	}
	return &nodev1.GetNodeResponse{Node: m.toProtoNode(node)}, nil
}

func (m *MockService) ListNodes(ctx context.Context, req *nodev1.ListNodesRequest) (*nodev1.ListNodesResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var nodes []*nodev1.Node
	for _, n := range m.nodes {
		nodes = append(nodes, m.toProtoNode(n))
	}
	return &nodev1.ListNodesResponse{Nodes: nodes}, nil
}

func (m *MockService) ListNodeChainConfigs(ctx context.Context, req *nodev1.ListNodeChainConfigsRequest) (*nodev1.ListNodeChainConfigsResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var configs []*nodev1.ChainConfig
	if req.Filter != nil && len(req.Filter.NodeIds) > 0 {
		for _, nodeID := range req.Filter.NodeIds {
			if cc, ok := m.chainConfigs[nodeID]; ok {
				configs = append(configs, cc...)
			}
		}
	}
	return &nodev1.ListNodeChainConfigsResponse{ChainConfigs: configs}, nil
}

func (m *MockService) ProposeJob(ctx context.Context, req *jobv1.ProposeJobRequest) (*jobv1.ProposeJobResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.proposalSeq++
	id := fmt.Sprintf("proposal-%d", m.proposalSeq)

	proposal := &jobv1.Proposal{
		Id:             id,
		Status:         jobv1.ProposalStatus_PROPOSAL_STATUS_PENDING,
		DeliveryStatus: jobv1.ProposalDeliveryStatus_PROPOSAL_DELIVERY_STATUS_DELIVERED,
		Spec:           req.Spec,
	}
	m.proposals[id] = proposal

	return &jobv1.ProposeJobResponse{Proposal: proposal}, nil
}

func (m *MockService) GetProposal(ctx context.Context, req *jobv1.GetProposalRequest) (*jobv1.GetProposalResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	proposal, ok := m.proposals[req.Id]
	if !ok {
		return nil, fmt.Errorf("proposal not found: %s", req.Id)
	}
	return &jobv1.GetProposalResponse{Proposal: proposal}, nil
}

func (m *MockService) ListJobs(ctx context.Context, req *jobv1.ListJobsRequest) (*jobv1.ListJobsResponse, error) {
	return &jobv1.ListJobsResponse{}, nil
}

func (m *MockService) toProtoNode(n *nodeData) *nodev1.Node {
	return &nodev1.Node{
		Id:          n.id,
		Name:        n.name,
		PublicKey:   n.publicKey,
		Labels:      n.labels,
		IsEnabled:   true,
		IsConnected: n.connected,
	}
}

func convertLabels(in []*ptypes.Label) []*ptypes.Label {
	if in == nil {
		return nil
	}
	out := make([]*ptypes.Label, len(in))
	copy(out, in)
	return out
}
