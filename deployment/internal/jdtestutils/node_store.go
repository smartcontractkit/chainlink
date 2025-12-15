package jdtestutils

import (
	"errors"
	"fmt"
	"maps"
	"sync"

	"github.com/smartcontractkit/chainlink/deployment/utils/nodetestutils"
)

var (
	errNoExist = errors.New("does not exist")
)

// nodeStore is an interface for storing nodes.
type nodeStore interface {
	put(nodeID string, node *nodetestutils.Node) error
	get(nodeID string) (*nodetestutils.Node, error)
	list() []*nodetestutils.Node
	asMap() map[string]*nodetestutils.Node
	delete(nodeID string) error
}

var _ nodeStore = &mapNodeStore{}

type mapNodeStore struct {
	mu    sync.Mutex
	nodes map[string]*nodetestutils.Node
}

func newMapNodeStore(n map[string]*nodetestutils.Node) *mapNodeStore {
	return &mapNodeStore{
		nodes: n,
	}
}

func (m *mapNodeStore) put(nodeID string, node *nodetestutils.Node) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.nodes == nil {
		m.nodes = make(map[string]*nodetestutils.Node)
	}
	m.nodes[nodeID] = node
	return nil
}

func (m *mapNodeStore) get(nodeID string) (*nodetestutils.Node, error) {
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

func (m *mapNodeStore) list() []*nodetestutils.Node {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.nodes == nil {
		return nil
	}
	nodes := make([]*nodetestutils.Node, 0)
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

func (m *mapNodeStore) asMap() map[string]*nodetestutils.Node {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.nodes == nil {
		return nil
	}
	nodes := make(map[string]*nodetestutils.Node)
	maps.Copy(nodes, m.nodes)
	return nodes
}
