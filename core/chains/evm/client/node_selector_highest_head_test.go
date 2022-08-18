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
			node.On("State").Return(evmclient.NodeStateOutOfSync)
			node.On("LatestReceivedBlockNumber").Return(int64(-1))
		} else if i == 1 {
			// second node is alive, LatestReceivedBlockNumber = 1
			node.On("State").Return(evmclient.NodeStateAlive)
			node.On("LatestReceivedBlockNumber").Return(int64(1))
		} else {
			// third node is alive, LatestReceivedBlockNumber = 2 (best node)
			node.On("State").Return(evmclient.NodeStateAlive)
			node.On("LatestReceivedBlockNumber").Return(int64(2))
		}
		nodes = append(nodes, node)
	}

	selector := evmclient.NewHighestHeadNodeSelector(nodes)
	assert.Equal(t, nodes[2], selector.Select())

	t.Run("stick to the same node", func(t *testing.T) {
		node := evmmocks.NewNode(t)
		// fourth node is alive, LatestReceivedBlockNumber = 2 (same as 3rd)
		node.On("State").Return(evmclient.NodeStateAlive)
		node.On("LatestReceivedBlockNumber").Return(int64(2))
		nodes = append(nodes, node)

		selector := evmclient.NewHighestHeadNodeSelector(nodes)
		assert.Equal(t, nodes[2], selector.Select())
	})

	t.Run("another best node", func(t *testing.T) {
		node := evmmocks.NewNode(t)
		// fifth node is alive, LatestReceivedBlockNumber = 3 (better than 3rd and 4th)
		node.On("State").Return(evmclient.NodeStateAlive).Maybe()
		node.On("LatestReceivedBlockNumber").Return(int64(3))
		nodes = append(nodes, node)

		selector := evmclient.NewHighestHeadNodeSelector(nodes)
		assert.Equal(t, nodes[4], selector.Select())
	})
}

func TestHighestHeadNodeSelector_None(t *testing.T) {
	t.Parallel()

	var nodes []evmclient.Node

	for i := 0; i < 3; i++ {
		node := evmmocks.NewNode(t)
		if i == 0 {
			// first node is out of sync
			node.On("State").Return(evmclient.NodeStateOutOfSync)
			node.On("LatestReceivedBlockNumber").Return(int64(-1))
		} else {
			// others are unreachable
			node.On("State").Return(evmclient.NodeStateUnreachable)
			node.On("LatestReceivedBlockNumber").Return(int64(1))
		}
		nodes = append(nodes, node)
	}

	selector := evmclient.NewHighestHeadNodeSelector(nodes)
	assert.Nil(t, selector.Select())
}
