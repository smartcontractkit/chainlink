package memory

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/smartcontractkit/chainlink/deployment/environment/test"
	"github.com/smartcontractkit/chainlink/v2/core/services/feeds"
)

type JobApprover interface {
	AutoApproveJob(ctx context.Context, p *feeds.ProposeJobArgs) error
}

type autoApprovalNode struct {
	*Node
}

var _ JobApprover = &autoApprovalNode{}

func (q *autoApprovalNode) AutoApproveJob(ctx context.Context, p *feeds.ProposeJobArgs) error {
	appProposalID, err := q.App.GetFeedsService().ProposeJob(ctx, p)
	if err != nil {
		return fmt.Errorf("failed to propose job: %w", err)
	}
	// auto approve
	proposedSpec, err := q.App.GetFeedsService().ListSpecsByJobProposalIDs(ctx, []int64{appProposalID})
	if err != nil {
		return fmt.Errorf("failed to list specs: %w", err)
	}
	// possible to have multiple specs for the same job proposal id; take the last one
	if len(proposedSpec) == 0 {
		return fmt.Errorf("no specs found for job proposal id: %d", appProposalID)
	}
	err = q.App.GetFeedsService().ApproveSpec(ctx, proposedSpec[len(proposedSpec)-1].ID, true)
	if err != nil {
		return fmt.Errorf("failed to approve job: %w", err)
	}
	return nil
}

type jobApproverGetter struct {
	s nodeStore
}

func (w *jobApproverGetter) Get(nodeID string) (test.JobApprover, error) {
	node, err := w.s.get(nodeID)
	if err != nil {
		return nil, err
	}
	return &autoApprovalNode{node}, nil
}

type ExternalJobIDExtractor struct {
	ExternalJobID string `toml:"externalJobID"`
}

var errNoExist = errors.New("does not exist")

// nodeStore is an interface for storing nodes.
type nodeStore interface {
	put(nodeID string, node *Node) error
	get(nodeID string) (*Node, error)
	list() []*Node
	asMap() map[string]*Node
	delete(nodeID string) error
}

var _ nodeStore = &mapNodeStore{}

type mapNodeStore struct {
	mu    sync.Mutex
	nodes map[string]*Node
}

func newMapNodeStore(n map[string]*Node) *mapNodeStore {
	return &mapNodeStore{
		nodes: n,
	}
}
func (m *mapNodeStore) put(nodeID string, node *Node) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.nodes == nil {
		m.nodes = make(map[string]*Node)
	}
	m.nodes[nodeID] = node
	return nil
}
func (m *mapNodeStore) get(nodeID string) (*Node, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.nodes == nil {
		return nil, fmt.Errorf("node not found: %s", nodeID)
	}
	node, ok := m.nodes[nodeID]
	if !ok {
		return nil, fmt.Errorf("%w: node not found: %s", errNoExist, nodeID)
	}
	return node, nil
}
func (m *mapNodeStore) list() []*Node {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.nodes == nil {
		return nil
	}
	nodes := make([]*Node, 0)
	for _, node := range m.nodes {
		nodes = append(nodes, node)
	}
	return nodes
}
func (m *mapNodeStore) delete(nodeID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.nodes == nil {
		return fmt.Errorf("node not found: %s", nodeID)
	}
	_, ok := m.nodes[nodeID]
	if !ok {
		return nil
	}
	delete(m.nodes, nodeID)
	return nil
}

func (m *mapNodeStore) asMap() map[string]*Node {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.nodes == nil {
		return nil
	}
	nodes := make(map[string]*Node)
	for k, v := range m.nodes {
		nodes[k] = v
	}
	return nodes
}
