package test

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	//"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/csa"

	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
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
	mu sync.RWMutex
	//	store map[string]*wrapper
	store *store
	node.UnimplementedNodeServiceServer
}

type p2pKey string

func (p p2pKey) String() string {
	return string(p)
}
func (p p2pKey) Validate() error {
	_, err := p2pkey.MakePeerID(p.String())
	return err
}

type csaKey = string
type store struct {
	mu      sync.RWMutex
	db      map[p2pKey]*wrapper
	csa2p2p map[csaKey]p2pKey
}

func (s *store) getNode(p2p p2pKey) (*wrapper, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n, ok := s.db[p2p]
	if !ok {
		return nil, fmt.Errorf("node not found for p2p %s", p2p)
	}
	return n, nil
}

func (s *store) getNodeByCSA(csa csaKey) (*wrapper, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p2p, ok := s.csa2p2p[csa]
	if !ok {
		return nil, fmt.Errorf("node not found for csa key %s", csa)
	}
	return s.getNode(p2p)
}

func (s *store) list() []*wrapper {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []*wrapper
	for _, v := range s.db {
		out = append(out, v)
	}
	return out
}

func (s *store) put(n *wrapper) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.db[p2pKey(n.Node.PeerID.String())] = n
	s.csa2p2p[csaKey(n.Node.CSAKey)] = p2pKey(n.Node.PeerID.String())
}

func newStore(node []deployment.Node) *store {
	s := &store{
		db:      make(map[p2pKey]*wrapper),
		csa2p2p: make(map[csaKey]p2pKey),
	}
	for _, v := range node {
		w := newWrapper(v)
		s.db[p2pKey(w.Node.PeerID.String())] = w
		s.csa2p2p[csaKey(w.Node.CSAKey)] = p2pKey(w.Node.PeerID.String())
	}
	return s
}

func NewJDService(nodes []deployment.Node) *JDNodeService {
	return &JDNodeService{
		//store: wrapAll(nodes),
		store: newStore(nodes),
	}
}

func newWrapperFromJDNode(n *nodev1.Node) *wrapper {
	return nil
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

func (s *JDNodeService) GetNode(ctx context.Context, req *nodev1.GetNodeRequest) (*nodev1.GetNodeResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	p2p := p2pKey(req.Id)
	if err := p2p.Validate(); err != nil {
		return nil, fmt.Errorf("request ID is not a valid peer ID: %w", err)
	}
	w, err := s.store.getNode(p2p)
	if err != nil {
		return nil, err
	}

	return &nodev1.GetNodeResponse{
		Node: newJDNode(w.Node),
	}, nil
}

func (s *JDNodeService) ListNodes(ctx context.Context, req *nodev1.ListNodesRequest) (*nodev1.ListNodesResponse, error) {
	include := func(node *nodev1.Node) bool {
		if req.Filter == nil {
			return true
		}
		if len(req.Filter.Ids) > 0 {
			idx := slices.IndexFunc(req.Filter.Ids, func(id string) bool {
				return node.Id == id
			})
			if idx < 0 {
				return false
			}
		}
		for _, selector := range req.Filter.Selectors {
			idx := slices.IndexFunc(node.Labels, func(label *ptypes.Label) bool {
				return label.Key == selector.Key
			})
			if idx < 0 {
				return false
			}
			label := node.Labels[idx]

			switch selector.Op {
			case ptypes.SelectorOp_IN:
				values := strings.Split(*selector.Value, ",")
				found := slices.Contains(values, *label.Value)
				if !found {
					return false
				}
			default:
				panic("unimplemented selector")
			}
		}
		return true
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	var nodes []*nodev1.Node
	for _, w := range s.store.list() {
		n := newJDNode(w.Node)
		if include(n) {
			nodes = append(nodes, n)
		}
	}

	return &nodev1.ListNodesResponse{
		Nodes: nodes,
	}, nil
}

func (s *JDNodeService) DisableNode(ctx context.Context, req *nodev1.DisableNodeRequest) (*nodev1.DisableNodeResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	p2p := p2pKey(req.Id)
	if err := p2p.Validate(); err != nil {
		return nil, fmt.Errorf("request ID is not a valid peer ID: %w", err)
	}

	node, err := s.store.getNode(p2p)
	if err != nil {
		return nil, err
	}

	// Implement the logic to disable the node
	node.enabled = false
	s.store.put(node)

	return &nodev1.DisableNodeResponse{}, nil
}

func (s *JDNodeService) EnableNode(ctx context.Context, req *nodev1.EnableNodeRequest) (*nodev1.EnableNodeResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	p2p := p2pKey(req.Id)
	if err := p2p.Validate(); err != nil {
		return nil, fmt.Errorf("request ID is not a valid peer ID: %w", err)
	}

	node, err := s.store.getNode(p2p)
	if err != nil {
		return nil, err
	}

	// Implement the logic to enable the node
	node.enabled = true
	s.store.put(node)

	return &nodev1.EnableNodeResponse{}, nil
}

func (s *JDNodeService) RegisterNode(ctx context.Context, req *nodev1.RegisterNodeRequest) (*nodev1.RegisterNodeResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	n, _ := s.store.getNodeByCSA(csaKey(req.PublicKey))
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

func (s *JDNodeService) ListNodeChainConfigs(ctx context.Context, req *nodev1.ListNodeChainConfigsRequest) (*nodev1.ListNodeChainConfigsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var out []*nodev1.ChainConfig
	for _, w := range s.store.list() {
		cc, err := w.Node.ChainConfigs()
		if err != nil {
			return nil, err
		}
		out = append(out, cc...)
	}
	return &nodev1.ListNodeChainConfigsResponse{
		ChainConfigs: out,
	}, nil
}

func newWrapperFromRegister(req *nodev1.RegisterNodeRequest) (*wrapper, error) {
	return nil, nil
}

func (s *JDNodeService) UpdateNode(ctx context.Context, req *nodev1.UpdateNodeRequest) (*nodev1.UpdateNodeResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.store.getNode(p2pKey(req.Id))
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

func newWrapperFromUpdate(req *nodev1.UpdateNodeRequest) (*wrapper, error) {
	return nil, nil
}

func newJDNode(n deployment.Node) *nodev1.Node {
	out := nodev1.Node{
		Id:        n.PeerID.String(),
		Labels:    n.Labels,
		Name:      n.Name,
		PublicKey: n.CSAKey,
	}

	return &out
}

func newDeploymentNode(n *nodev1.Node) (deployment.Node, error) {
	p, err := p2pkey.MakePeerID(n.Id)
	if err != nil {
		return deployment.Node{}, fmt.Errorf("only support jd nodes with id as peer id. is %s is not peer id: %w", n.Id, err)
	}
	return deployment.Node{
		NodeID: p.String(),
		PeerID: p,
		Labels: n.Labels,
		Name:   n.Name,
		CSAKey: n.PublicKey,
	}, nil
}
