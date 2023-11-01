package client

import (
	"testing"

	"github.com/smartcontractkit/chainlink/v2/common/types"

	"github.com/stretchr/testify/assert"
)

func TestPriorityLevelNodeSelectorName(t *testing.T) {
	selector := newNodeSelector[types.ID, Head, NodeClient[types.ID, Head]](NodeSelectionModePriorityLevel, nil)
	assert.Equal(t, selector.Name(), NodeSelectionModePriorityLevel)
}

func TestPriorityLevelNodeSelector(t *testing.T) {
	t.Parallel()

	type nodeClient NodeClient[types.ID, Head]
	var nodes []Node[types.ID, Head, nodeClient]
	n1 := newMockNode[types.ID, Head, nodeClient](t)
	n1.On("State").Return(nodeStateAlive)
	n1.On("Order").Return(int32(1))

	n2 := newMockNode[types.ID, Head, nodeClient](t)
	n2.On("State").Return(nodeStateAlive)
	n2.On("Order").Return(int32(1))

	n3 := newMockNode[types.ID, Head, nodeClient](t)
	n3.On("State").Return(nodeStateAlive)
	n3.On("Order").Return(int32(1))

	nodes = append(nodes, n1, n2, n3)
	selector := newNodeSelector(NodeSelectionModePriorityLevel, nodes)
	assert.Same(t, nodes[0], selector.Select())
	assert.Same(t, nodes[1], selector.Select())
	assert.Same(t, nodes[2], selector.Select())
	assert.Same(t, nodes[0], selector.Select())
	assert.Same(t, nodes[1], selector.Select())
	assert.Same(t, nodes[2], selector.Select())
}

func TestPriorityLevelNodeSelector_None(t *testing.T) {
	t.Parallel()

	type nodeClient NodeClient[types.ID, Head]
	var nodes []Node[types.ID, Head, nodeClient]

	for i := 0; i < 3; i++ {
		node := newMockNode[types.ID, Head, nodeClient](t)
		if i == 0 {
			// first node is out of sync
			node.On("State").Return(nodeStateOutOfSync)
			node.On("Order").Return(int32(1))
		} else {
			// others are unreachable
			node.On("State").Return(nodeStateUnreachable)
			node.On("Order").Return(int32(1))
		}
		nodes = append(nodes, node)
	}

	selector := newNodeSelector(NodeSelectionModePriorityLevel, nodes)
	assert.Nil(t, selector.Select())
}

func TestPriorityLevelNodeSelector_DifferentOrder(t *testing.T) {
	t.Parallel()

	type nodeClient NodeClient[types.ID, Head]
	var nodes []Node[types.ID, Head, nodeClient]
	n1 := newMockNode[types.ID, Head, nodeClient](t)
	n1.On("State").Return(nodeStateAlive)
	n1.On("Order").Return(int32(1))

	n2 := newMockNode[types.ID, Head, nodeClient](t)
	n2.On("State").Return(nodeStateAlive)
	n2.On("Order").Return(int32(2))

	n3 := newMockNode[types.ID, Head, nodeClient](t)
	n3.On("State").Return(nodeStateAlive)
	n3.On("Order").Return(int32(3))

	nodes = append(nodes, n1, n2, n3)
	selector := newNodeSelector(NodeSelectionModePriorityLevel, nodes)
	assert.Same(t, nodes[0], selector.Select())
	assert.Same(t, nodes[0], selector.Select())
}
