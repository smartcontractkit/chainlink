package client_test

import (
	"testing"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	evmmocks "github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/stretchr/testify/assert"
)

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
			node.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(1), utils.NewBigI(7))
		} else {
			// third node is alive and best
			node.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(2), utils.NewBigI(8))
		}
		nodes = append(nodes, node)
	}

	selector := evmclient.NewTotalDifficultyNodeSelector(nodes)
	assert.Same(t, nodes[2], selector.Select())

	t.Run("stick to the same node", func(t *testing.T) {
		node := evmmocks.NewNode(t)
		// fourth node is alive (same as 3rd)
		node.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(2), utils.NewBigI(8))
		nodes = append(nodes, node)

		selector := evmclient.NewTotalDifficultyNodeSelector(nodes)
		assert.Same(t, nodes[2], selector.Select())
	})

	t.Run("another best node", func(t *testing.T) {
		node := evmmocks.NewNode(t)
		// fifth node is alive (better than 3rd and 4th)
		node.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(3), utils.NewBigI(11))
		nodes = append(nodes, node)

		selector := evmclient.NewTotalDifficultyNodeSelector(nodes)
		assert.Same(t, nodes[4], selector.Select())
	})

	t.Run("nodes never update latest block number", func(t *testing.T) {
		node1 := evmmocks.NewNode(t)
		node1.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(-1), nil)
		node2 := evmmocks.NewNode(t)
		node2.On("StateAndLatest").Return(evmclient.NodeStateAlive, int64(-1), nil)
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
			node.On("StateAndLatest").Return(evmclient.NodeStateUnreachable, int64(1), utils.NewBigI(7))
		}
		nodes = append(nodes, node)
	}

	selector := evmclient.NewTotalDifficultyNodeSelector(nodes)
	assert.Nil(t, selector.Select())
}
