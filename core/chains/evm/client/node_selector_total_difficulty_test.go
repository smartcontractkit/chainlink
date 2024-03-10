package client_test

import (
	"math/big"
	"testing"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/mocks"

	"github.com/stretchr/testify/assert"
)

func TestTotalDifficultyNodeSelectorName(t *testing.T) {
	selector := evmclient.NewTotalDifficultyNodeSelector(nil)
	assert.Equal(t, selector.Name(), evmclient.NodeSelectionMode_TotalDifficulty)
}

func TestTotalDifficultyNodeSelector(t *testing.T) {
	t.Parallel()

	var nodes []evmclient.Node

	for i := 0; i < 3; i++ {
		node := evmmocks.NewNode(t)
		if i == 0 {
			// first node is out of sync
			node.On("StateAndLatest").Return(evmclient.NodeStateOutOfSync, int64(-1), nil)
		} else if i == 1 {
			// second node is alive
			node.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(1), big.NewInt(7))
		} else {
			// third node is alive and best
			node.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(2), big.NewInt(8))
		}
		node.On("Order").Maybe().Return(int32(1))
		nodes = append(nodes, node)
	}

	selector := evmclient.NewTotalDifficultyNodeSelector(nodes)
	assert.Same(t, nodes[2], selector.Select())

	t.Run("stick to the same node", func(t *testing.T) {
		node := evmmocks.NewNode(t)
		// fourth node is alive (same as 3rd)
		node.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(2), big.NewInt(8))
		node.On("Order").Maybe().Return(int32(1))
		nodes = append(nodes, node)

		selector := evmclient.NewTotalDifficultyNodeSelector(nodes)
		assert.Same(t, nodes[2], selector.Select())
	})

	t.Run("another best node", func(t *testing.T) {
		node := evmmocks.NewNode(t)
		// fifth node is alive (better than 3rd and 4th)
		node.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(3), big.NewInt(11))
		node.On("Order").Maybe().Return(int32(1))
		nodes = append(nodes, node)

		selector := evmclient.NewTotalDifficultyNodeSelector(nodes)
		assert.Same(t, nodes[4], selector.Select())
	})

	t.Run("nodes never update latest block number", func(t *testing.T) {
		node1 := evmmocks.NewNode(t)
		node1.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(-1), nil)
		node1.On("Order").Maybe().Return(int32(1))
		node2 := evmmocks.NewNode(t)
		node2.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(-1), nil)
		node2.On("Order").Maybe().Return(int32(1))
		nodes := []evmclient.Node{node1, node2}

		selector := evmclient.NewTotalDifficultyNodeSelector(nodes)
		assert.Same(t, node1, selector.Select())
	})
}

func TestTotalDifficultyNodeSelector_None(t *testing.T) {
	t.Parallel()

	var nodes []evmclient.Node

	for i := 0; i < 3; i++ {
		node := evmmocks.NewNode(t)
		if i == 0 {
			// first node is out of sync
			node.On("StateAndLatest").Return(evmclient.NodeStateOutOfSync, int64(-1), nil)
		} else {
			// others are unreachable
			node.On("StateAndLatest").Return(evmclient.NodeStateUnreachable, int64(1), big.NewInt(7))
		}
		nodes = append(nodes, node)
	}

	selector := evmclient.NewTotalDifficultyNodeSelector(nodes)
	assert.Nil(t, selector.Select())
}

func TestTotalDifficultyNodeSelectorWithOrder(t *testing.T) {
	t.Parallel()

	var nodes []evmclient.Node

	t.Run("same td and order", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			node := evmmocks.NewNode(t)
			node.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(1), big.NewInt(10))
			node.On("Order").Return(int32(2))
			nodes = append(nodes, node)
		}
		selector := evmclient.NewTotalDifficultyNodeSelector(nodes)
		//Should select the first node because all things are equal
		assert.Same(t, nodes[0], selector.Select())
	})

	t.Run("same td but different order", func(t *testing.T) {
		node1 := evmmocks.NewNode(t)
		node1.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(3), big.NewInt(10))
		node1.On("Order").Return(int32(3))

		node2 := evmmocks.NewNode(t)
		node2.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(3), big.NewInt(10))
		node2.On("Order").Return(int32(1))

		node3 := evmmocks.NewNode(t)
		node3.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(3), big.NewInt(10))
		node3.On("Order").Return(int32(2))

		nodes := []evmclient.Node{node1, node2, node3}
		selector := evmclient.NewTotalDifficultyNodeSelector(nodes)
		//Should select the second node as it has the highest priority
		assert.Same(t, nodes[1], selector.Select())
	})

	t.Run("different td but same order", func(t *testing.T) {
		node1 := evmmocks.NewNode(t)
		node1.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(1), big.NewInt(10))
		node1.On("Order").Maybe().Return(int32(3))

		node2 := evmmocks.NewNode(t)
		node2.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(1), big.NewInt(11))
		node2.On("Order").Maybe().Return(int32(3))

		node3 := evmmocks.NewNode(t)
		node3.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(1), big.NewInt(12))
		node3.On("Order").Return(int32(3))

		nodes := []evmclient.Node{node1, node2, node3}
		selector := evmclient.NewTotalDifficultyNodeSelector(nodes)
		//Should select the third node as it has the highest td
		assert.Same(t, nodes[2], selector.Select())
	})

	t.Run("different head and different order", func(t *testing.T) {
		node1 := evmmocks.NewNode(t)
		node1.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(1), big.NewInt(100))
		node1.On("Order").Maybe().Return(int32(4))

		node2 := evmmocks.NewNode(t)
		node2.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(1), big.NewInt(110))
		node2.On("Order").Maybe().Return(int32(5))

		node3 := evmmocks.NewNode(t)
		node3.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(1), big.NewInt(110))
		node3.On("Order").Maybe().Return(int32(1))

		node4 := evmmocks.NewNode(t)
		node4.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(1), big.NewInt(105))
		node4.On("Order").Maybe().Return(int32(2))

		nodes := []evmclient.Node{node1, node2, node3, node4}
		selector := evmclient.NewTotalDifficultyNodeSelector(nodes)
		//Should select the third node as it has the highest td and will win the priority tie-breaker
		assert.Same(t, nodes[2], selector.Select())
	})
}
