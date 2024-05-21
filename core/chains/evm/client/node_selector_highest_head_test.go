package client_test

import (
	"testing"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/mocks"

	"github.com/stretchr/testify/assert"
)

func TestHighestHeadNodeSelectorName(t *testing.T) {
	selector := evmclient.NewHighestHeadNodeSelector(nil)
	assert.Equal(t, selector.Name(), evmclient.NodeSelectionMode_HighestHead)
}

func TestHighestHeadNodeSelector(t *testing.T) {
	t.Parallel()

	var nodes []evmclient.Node

	for i := 0; i < 3; i++ {
		node := evmmocks.NewNode(t)
		if i == 0 {
			// first node is out of sync
			node.On("StateAndLatest").Return(evmclient.NodeStateOutOfSync, int64(-1), nil)
		} else if i == 1 {
			// second node is alive, LatestReceivedBlockNumber = 1
			node.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(1), nil)
		} else {
			// third node is alive, LatestReceivedBlockNumber = 2 (best node)
			node.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(2), nil)
		}
		node.On("Order").Maybe().Return(int32(1))
		nodes = append(nodes, node)
	}

	selector := evmclient.NewHighestHeadNodeSelector(nodes)
	assert.Same(t, nodes[2], selector.Select())

	t.Run("stick to the same node", func(t *testing.T) {
		node := evmmocks.NewNode(t)
		// fourth node is alive, LatestReceivedBlockNumber = 2 (same as 3rd)
		node.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(2), nil)
		node.On("Order").Return(int32(1))
		nodes = append(nodes, node)

		selector := evmclient.NewHighestHeadNodeSelector(nodes)
		assert.Same(t, nodes[2], selector.Select())
	})

	t.Run("another best node", func(t *testing.T) {
		node := evmmocks.NewNode(t)
		// fifth node is alive, LatestReceivedBlockNumber = 3 (better than 3rd and 4th)
		node.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(3), nil)
		node.On("Order").Return(int32(1))
		nodes = append(nodes, node)

		selector := evmclient.NewHighestHeadNodeSelector(nodes)
		assert.Same(t, nodes[4], selector.Select())
	})

	t.Run("nodes never update latest block number", func(t *testing.T) {
		node1 := evmmocks.NewNode(t)
		node1.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(-1), nil)
		node1.On("Order").Return(int32(1))
		node2 := evmmocks.NewNode(t)
		node2.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(-1), nil)
		node2.On("Order").Return(int32(1))
		nodes := []evmclient.Node{node1, node2}

		selector := evmclient.NewHighestHeadNodeSelector(nodes)
		assert.Same(t, node1, selector.Select())
	})
}

func TestHighestHeadNodeSelector_None(t *testing.T) {
	t.Parallel()

	var nodes []evmclient.Node

	for i := 0; i < 3; i++ {
		node := evmmocks.NewNode(t)
		if i == 0 {
			// first node is out of sync
			node.On("StateAndLatest").Return(evmclient.NodeStateOutOfSync, int64(-1), nil)
		} else {
			// others are unreachable
			node.On("StateAndLatest").Return(evmclient.NodeStateUnreachable, int64(1), nil)
		}
		nodes = append(nodes, node)
	}

	selector := evmclient.NewHighestHeadNodeSelector(nodes)
	assert.Nil(t, selector.Select())
}

func TestHighestHeadNodeSelectorWithOrder(t *testing.T) {
	t.Parallel()

	var nodes []evmclient.Node

	t.Run("same head and order", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			node := evmmocks.NewNode(t)
			node.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(1), nil)
			node.On("Order").Return(int32(2))
			nodes = append(nodes, node)
		}
		selector := evmclient.NewHighestHeadNodeSelector(nodes)
		//Should select the first node because all things are equal
		assert.Same(t, nodes[0], selector.Select())
	})

	t.Run("same head but different order", func(t *testing.T) {
		node1 := evmmocks.NewNode(t)
		node1.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(3), nil)
		node1.On("Order").Return(int32(3))

		node2 := evmmocks.NewNode(t)
		node2.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(3), nil)
		node2.On("Order").Return(int32(1))

		node3 := evmmocks.NewNode(t)
		node3.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(3), nil)
		node3.On("Order").Return(int32(2))

		nodes := []evmclient.Node{node1, node2, node3}
		selector := evmclient.NewHighestHeadNodeSelector(nodes)
		//Should select the second node as it has the highest priority
		assert.Same(t, nodes[1], selector.Select())
	})

	t.Run("different head but same order", func(t *testing.T) {
		node1 := evmmocks.NewNode(t)
		node1.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(1), nil)
		node1.On("Order").Maybe().Return(int32(3))

		node2 := evmmocks.NewNode(t)
		node2.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(2), nil)
		node2.On("Order").Maybe().Return(int32(3))

		node3 := evmmocks.NewNode(t)
		node3.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(3), nil)
		node3.On("Order").Return(int32(3))

		nodes := []evmclient.Node{node1, node2, node3}
		selector := evmclient.NewHighestHeadNodeSelector(nodes)
		//Should select the third node as it has the highest head
		assert.Same(t, nodes[2], selector.Select())
	})

	t.Run("different head and different order", func(t *testing.T) {
		node1 := evmmocks.NewNode(t)
		node1.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(10), nil)
		node1.On("Order").Maybe().Return(int32(3))

		node2 := evmmocks.NewNode(t)
		node2.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(11), nil)
		node2.On("Order").Maybe().Return(int32(4))

		node3 := evmmocks.NewNode(t)
		node3.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(11), nil)
		node3.On("Order").Maybe().Return(int32(3))

		node4 := evmmocks.NewNode(t)
		node4.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(10), nil)
		node4.On("Order").Maybe().Return(int32(1))

		nodes := []evmclient.Node{node1, node2, node3, node4}
		selector := evmclient.NewHighestHeadNodeSelector(nodes)
		//Should select the third node as it has the highest head and will win the priority tie-breaker
		assert.Same(t, nodes[2], selector.Select())
	})
}
