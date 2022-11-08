package client_test

import (
	"testing"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	evmmocks "github.com/smartcontractkit/chainlink/core/chains/evm/mocks"

	"github.com/stretchr/testify/assert"
)

func TestHighestHeadNodeSelector(t *testing.T) {
	t.Parallel()

	var nodes []evmclient.Node

	for i := 0; i < 3; i++ {
		node := evmmocks.NewNode(t)
		if i == 0 {
			// first node is out of sync
			node.On("StateAndLatestBlockNumber").Return(evmclient.NodeStateOutOfSync, int64(-1))
		} else if i == 1 {
			// second node is alive, LatestReceivedBlockNumber = 1
			node.On("StateAndLatestBlockNumber").Return(evmclient.NodeStateAlive, int64(1))
		} else {
			// third node is alive, LatestReceivedBlockNumber = 2 (best node)
			node.On("StateAndLatestBlockNumber").Return(evmclient.NodeStateAlive, int64(2))
		}
		nodes = append(nodes, node)
	}

	selector := evmclient.NewHighestHeadNodeSelector(nodes)
	assert.Equal(t, nodes[2], selector.Select())

	t.Run("stick to the same node", func(t *testing.T) {
		node := evmmocks.NewNode(t)
		// fourth node is alive, LatestReceivedBlockNumber = 2 (same as 3rd)
		node.On("StateAndLatestBlockNumber").Return(evmclient.NodeStateAlive, int64(2))
		nodes = append(nodes, node)

		selector := evmclient.NewHighestHeadNodeSelector(nodes)
		assert.Equal(t, nodes[2], selector.Select())
	})

	t.Run("another best node", func(t *testing.T) {
		node := evmmocks.NewNode(t)
		// fifth node is alive, LatestReceivedBlockNumber = 3 (better than 3rd and 4th)
		node.On("StateAndLatestBlockNumber").Return(evmclient.NodeStateAlive, int64(3))
		nodes = append(nodes, node)

		selector := evmclient.NewHighestHeadNodeSelector(nodes)
		assert.Equal(t, nodes[4], selector.Select())
	})

	t.Run("update lastBestNode", func(t *testing.T) {
		node1 := evmmocks.NewNode(t)
		node1.On("StateAndLatestBlockNumber").Return(evmclient.NodeStateAlive, int64(1)).Once()
		node1.On("StateAndLatestBlockNumber").Return(evmclient.NodeStateOutOfSync, int64(1)).Once()
		node1.On("StateAndLatestBlockNumber").Return(evmclient.NodeStateAlive, int64(1)).Once()
		node2 := evmmocks.NewNode(t)
		node2.On("StateAndLatestBlockNumber").Return(evmclient.NodeStateAlive, int64(1))
		nodes := []evmclient.Node{node1, node2}

		selector := evmclient.NewHighestHeadNodeSelector(nodes)
		// node1 would be set as lastBestNode (it is alive)
		assert.Equal(t, node1, selector.Select())

		// node1 is out of sync, node2 must be selected and set as the last best node
		assert.Equal(t, node2, selector.Select())

		// node1 is alive again, node2 must be returned as the last best node
		assert.Equal(t, node2, selector.Select())
	})

	t.Run("nodes never update latest block number", func(t *testing.T) {
		node1 := evmmocks.NewNode(t)
		node1.On("StateAndLatestBlockNumber").Return(evmclient.NodeStateAlive, int64(-1))
		node2 := evmmocks.NewNode(t)
		node2.On("StateAndLatestBlockNumber").Return(evmclient.NodeStateAlive, int64(-1))
		nodes := []evmclient.Node{node1, node2}

		selector := evmclient.NewHighestHeadNodeSelector(nodes)
		assert.Equal(t, node1, selector.Select())
	})
}

func TestHighestHeadNodeSelector_None(t *testing.T) {
	t.Parallel()

	var nodes []evmclient.Node

	for i := 0; i < 3; i++ {
		node := evmmocks.NewNode(t)
		if i == 0 {
			// first node is out of sync
			node.On("StateAndLatestBlockNumber").Return(evmclient.NodeStateOutOfSync, int64(-1))
		} else {
			// others are unreachable
			node.On("StateAndLatestBlockNumber").Return(evmclient.NodeStateUnreachable, int64(1))
		}
		nodes = append(nodes, node)
	}

	selector := evmclient.NewHighestHeadNodeSelector(nodes)
	assert.Nil(t, selector.Select())
}
