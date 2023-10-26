package client

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/common/types"
)

func TestRoundRobinNodeSelectorName(t *testing.T) {
	selector := NewRoundRobinSelector[types.ID, Head, NodeClient[types.ID, Head]](nil)
	assert.Equal(t, selector.Name(), NodeSelectionMode_RoundRobin)
}

func TestRoundRobinNodeSelector(t *testing.T) {
	t.Parallel()

	type nodeClient NodeClient[types.ID, Head]
	var nodes []Node[types.ID, Head, nodeClient]

	for i := 0; i < 3; i++ {
		node := newMockNode[types.ID, Head, nodeClient](t)
		if i == 0 {
			// first node is out of sync
			node.On("State").Return(nodeStateOutOfSync)
		} else {
			// second & third nodes are alive
			node.On("State").Return(nodeStateAlive)
		}
		nodes = append(nodes, node)
	}

	selector := NewRoundRobinSelector(nodes)
	assert.Same(t, nodes[1], selector.Select())
	assert.Same(t, nodes[2], selector.Select())
	assert.Same(t, nodes[1], selector.Select())
	assert.Same(t, nodes[2], selector.Select())
}

func TestRoundRobinNodeSelector_None(t *testing.T) {
	t.Parallel()

	type nodeClient NodeClient[types.ID, Head]
	var nodes []Node[types.ID, Head, nodeClient]

	for i := 0; i < 3; i++ {
		node := newMockNode[types.ID, Head, nodeClient](t)
		if i == 0 {
			// first node is out of sync
			node.On("State").Return(nodeStateOutOfSync)
		} else {
			// others are unreachable
			node.On("State").Return(nodeStateUnreachable)
		}
		nodes = append(nodes, node)
	}

	selector := NewRoundRobinSelector(nodes)
	assert.Nil(t, selector.Select())
}
