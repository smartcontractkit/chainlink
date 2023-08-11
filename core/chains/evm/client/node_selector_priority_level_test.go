package client_test

import (
	"testing"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/mocks"

	"github.com/stretchr/testify/assert"
)

func TestPriorityLevelNodeSelectorName(t *testing.T) {
	selector := evmclient.NewPriorityLevelNodeSelector(nil)
	assert.Equal(t, selector.Name(), evmclient.NodeSelectionMode_PriorityLevel)
}

func TestPriorityLevelNodeSelector(t *testing.T) {
	t.Parallel()

	var nodes []evmclient.Node
	n1 := evmmocks.NewNode(t)
	n1.On("State").Return(evmclient.NodeStateAlive)
	n1.On("Order").Return(int32(1))

	n2 := evmmocks.NewNode(t)
	n2.On("State").Return(evmclient.NodeStateAlive)
	n2.On("Order").Return(int32(1))

	n3 := evmmocks.NewNode(t)
	n3.On("State").Return(evmclient.NodeStateAlive)
	n3.On("Order").Return(int32(1))

	nodes = append(nodes, n1, n2, n3)
	selector := evmclient.NewPriorityLevelNodeSelector(nodes)
	assert.Same(t, nodes[0], selector.Select())
	assert.Same(t, nodes[1], selector.Select())
	assert.Same(t, nodes[2], selector.Select())
	assert.Same(t, nodes[0], selector.Select())
	assert.Same(t, nodes[1], selector.Select())
	assert.Same(t, nodes[2], selector.Select())
}

func TestPriorityLevelNodeSelector_None(t *testing.T) {
	t.Parallel()

	var nodes []evmclient.Node

	for i := 0; i < 3; i++ {
		node := evmmocks.NewNode(t)
		if i == 0 {
			// first node is out of sync
			node.On("State").Return(evmclient.NodeStateOutOfSync)
			node.On("Order").Return(int32(1))
		} else {
			// others are unreachable
			node.On("State").Return(evmclient.NodeStateUnreachable)
			node.On("Order").Return(int32(1))
		}
		nodes = append(nodes, node)
	}

	selector := evmclient.NewPriorityLevelNodeSelector(nodes)
	assert.Nil(t, selector.Select())
}

func TestPriorityLevelNodeSelector_DifferentOrder(t *testing.T) {
	t.Parallel()

	var nodes []evmclient.Node
	n1 := evmmocks.NewNode(t)
	n1.On("State").Return(evmclient.NodeStateAlive)
	n1.On("Order").Return(int32(1))

	n2 := evmmocks.NewNode(t)
	n2.On("State").Return(evmclient.NodeStateAlive)
	n2.On("Order").Return(int32(2))

	n3 := evmmocks.NewNode(t)
	n3.On("State").Return(evmclient.NodeStateAlive)
	n3.On("Order").Return(int32(3))

	nodes = append(nodes, n1, n2, n3)
	selector := evmclient.NewPriorityLevelNodeSelector(nodes)
	assert.Same(t, nodes[0], selector.Select())
	assert.Same(t, nodes[0], selector.Select())
}
