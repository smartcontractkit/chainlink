package client

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/common/types"
)

func TestHighestHeadNodeSelectorName(t *testing.T) {
	selector := newNodeSelector[types.ID, Head, NodeClient[types.ID, Head]](NodeSelectionModeHighestHead, nil)
	assert.Equal(t, selector.Name(), NodeSelectionModeHighestHead)
}

func TestHighestHeadNodeSelector(t *testing.T) {
	t.Parallel()

	type nodeClient NodeClient[types.ID, Head]

	var nodes []Node[types.ID, Head, nodeClient]

	for i := 0; i < 3; i++ {
		node := newMockNode[types.ID, Head, nodeClient](t)
		if i == 0 {
			// first node is out of sync
			node.On("StateAndLatest").Return(nodeStateOutOfSync, ChainInfo{BlockNumber: int64(-1)})
		} else if i == 1 {
			// second node is alive, LatestReceivedBlockNumber = 1
			node.On("StateAndLatest").Return(nodeStateAlive, ChainInfo{BlockNumber: int64(1)})
		} else {
			// third node is alive, LatestReceivedBlockNumber = 2 (best node)
			node.On("StateAndLatest").Return(nodeStateAlive, ChainInfo{BlockNumber: int64(2)})
		}
		node.On("Order").Maybe().Return(int32(1))
		nodes = append(nodes, node)
	}

	selector := newNodeSelector[types.ID, Head, nodeClient](NodeSelectionModeHighestHead, nodes)
	assert.Same(t, nodes[2], selector.Select())

	t.Run("stick to the same node", func(t *testing.T) {
		node := newMockNode[types.ID, Head, nodeClient](t)
		// fourth node is alive, LatestReceivedBlockNumber = 2 (same as 3rd)
		node.On("StateAndLatest").Return(nodeStateAlive, ChainInfo{BlockNumber: int64(2)})
		node.On("Order").Return(int32(1))
		nodes = append(nodes, node)

		selector := newNodeSelector(NodeSelectionModeHighestHead, nodes)
		assert.Same(t, nodes[2], selector.Select())
	})

	t.Run("another best node", func(t *testing.T) {
		node := newMockNode[types.ID, Head, nodeClient](t)
		// fifth node is alive, LatestReceivedBlockNumber = 3 (better than 3rd and 4th)
		node.On("StateAndLatest").Return(nodeStateAlive, ChainInfo{BlockNumber: int64(3)})
		node.On("Order").Return(int32(1))
		nodes = append(nodes, node)

		selector := newNodeSelector(NodeSelectionModeHighestHead, nodes)
		assert.Same(t, nodes[4], selector.Select())
	})

	t.Run("nodes never update latest block number", func(t *testing.T) {
		node1 := newMockNode[types.ID, Head, nodeClient](t)
		node1.On("StateAndLatest").Return(nodeStateAlive, ChainInfo{BlockNumber: int64(-1)})
		node1.On("Order").Return(int32(1))
		node2 := newMockNode[types.ID, Head, nodeClient](t)
		node2.On("StateAndLatest").Return(nodeStateAlive, ChainInfo{BlockNumber: int64(-1)})
		node2.On("Order").Return(int32(1))
		selector := newNodeSelector(NodeSelectionModeHighestHead, []Node[types.ID, Head, nodeClient]{node1, node2})
		assert.Same(t, node1, selector.Select())
	})
}

func TestHighestHeadNodeSelector_None(t *testing.T) {
	t.Parallel()

	type nodeClient NodeClient[types.ID, Head]
	var nodes []Node[types.ID, Head, nodeClient]

	for i := 0; i < 3; i++ {
		node := newMockNode[types.ID, Head, nodeClient](t)
		if i == 0 {
			// first node is out of sync
			node.On("StateAndLatest").Return(nodeStateOutOfSync, ChainInfo{BlockNumber: int64(-1)})
		} else {
			// others are unreachable
			node.On("StateAndLatest").Return(nodeStateUnreachable, ChainInfo{BlockNumber: int64(-1)})
		}
		nodes = append(nodes, node)
	}

	selector := newNodeSelector(NodeSelectionModeHighestHead, nodes)
	assert.Nil(t, selector.Select())
}

func TestHighestHeadNodeSelectorWithOrder(t *testing.T) {
	t.Parallel()

	type nodeClient NodeClient[types.ID, Head]
	var nodes []Node[types.ID, Head, nodeClient]

	t.Run("same head and order", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			node := newMockNode[types.ID, Head, nodeClient](t)
			node.On("StateAndLatest").Return(nodeStateAlive, ChainInfo{BlockNumber: int64(1)})
			node.On("Order").Return(int32(2))
			nodes = append(nodes, node)
		}
		selector := newNodeSelector(NodeSelectionModeHighestHead, nodes)
		//Should select the first node because all things are equal
		assert.Same(t, nodes[0], selector.Select())
	})

	t.Run("same head but different order", func(t *testing.T) {
		node1 := newMockNode[types.ID, Head, nodeClient](t)
		node1.On("StateAndLatest").Return(nodeStateAlive, ChainInfo{BlockNumber: int64(3)})
		node1.On("Order").Return(int32(3))

		node2 := newMockNode[types.ID, Head, nodeClient](t)
		node2.On("StateAndLatest").Return(nodeStateAlive, ChainInfo{BlockNumber: int64(3)})
		node2.On("Order").Return(int32(1))

		node3 := newMockNode[types.ID, Head, nodeClient](t)
		node3.On("StateAndLatest").Return(nodeStateAlive, ChainInfo{BlockNumber: int64(3)})
		node3.On("Order").Return(int32(2))

		nodes := []Node[types.ID, Head, nodeClient]{node1, node2, node3}
		selector := newNodeSelector(NodeSelectionModeHighestHead, nodes)
		//Should select the second node as it has the highest priority
		assert.Same(t, nodes[1], selector.Select())
	})

	t.Run("different head but same order", func(t *testing.T) {
		node1 := newMockNode[types.ID, Head, nodeClient](t)
		node1.On("StateAndLatest").Return(nodeStateAlive, ChainInfo{BlockNumber: int64(1)})
		node1.On("Order").Maybe().Return(int32(3))

		node2 := newMockNode[types.ID, Head, nodeClient](t)
		node2.On("StateAndLatest").Return(nodeStateAlive, ChainInfo{BlockNumber: int64(2)})
		node2.On("Order").Maybe().Return(int32(3))

		node3 := newMockNode[types.ID, Head, nodeClient](t)
		node3.On("StateAndLatest").Return(nodeStateAlive, ChainInfo{BlockNumber: int64(3)})
		node3.On("Order").Return(int32(3))

		nodes := []Node[types.ID, Head, nodeClient]{node1, node2, node3}
		selector := newNodeSelector(NodeSelectionModeHighestHead, nodes)
		//Should select the third node as it has the highest head
		assert.Same(t, nodes[2], selector.Select())
	})

	t.Run("different head and different order", func(t *testing.T) {
		node1 := newMockNode[types.ID, Head, nodeClient](t)
		node1.On("StateAndLatest").Return(nodeStateAlive, ChainInfo{BlockNumber: int64(10)})
		node1.On("Order").Maybe().Return(int32(3))

		node2 := newMockNode[types.ID, Head, nodeClient](t)
		node2.On("StateAndLatest").Return(nodeStateAlive, ChainInfo{BlockNumber: int64(11)})
		node2.On("Order").Maybe().Return(int32(4))

		node3 := newMockNode[types.ID, Head, nodeClient](t)
		node3.On("StateAndLatest").Return(nodeStateAlive, ChainInfo{BlockNumber: int64(11)})
		node3.On("Order").Maybe().Return(int32(3))

		node4 := newMockNode[types.ID, Head, nodeClient](t)
		node4.On("StateAndLatest").Return(nodeStateAlive, ChainInfo{BlockNumber: int64(10)})
		node4.On("Order").Maybe().Return(int32(1))

		nodes := []Node[types.ID, Head, nodeClient]{node1, node2, node3, node4}
		selector := newNodeSelector(NodeSelectionModeHighestHead, nodes)
		//Should select the third node as it has the highest head and will win the priority tie-breaker
		assert.Same(t, nodes[2], selector.Select())
	})
}
