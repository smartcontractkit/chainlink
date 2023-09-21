package client_test

import (
	"math/big"
	"testing"

	commonclient "github.com/smartcontractkit/chainlink/v2/common/chains/client"
	commonmocks "github.com/smartcontractkit/chainlink/v2/common/chains/client/mocks"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/mocks"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"

	"github.com/stretchr/testify/assert"
)

func TestRoundRobinNodeSelector(t *testing.T) {
	t.Parallel()

	var nodes []evmclient.Node

	for i := 0; i < 3; i++ {
		node := evmmocks.NewNode(t)
		if i == 0 {
			// first node is out of sync
			node.On("State").Return(evmclient.NodeStateOutOfSync)
		} else {
			// second & third nodes are alive
			node.On("State").Return(evmclient.NodeStateAlive)
		}
		nodes = append(nodes, node)
	}

	selector := evmclient.NewRoundRobinSelector(nodes)
	assert.Same(t, nodes[1], selector.Select())
	assert.Same(t, nodes[2], selector.Select())
	assert.Same(t, nodes[1], selector.Select())
	assert.Same(t, nodes[2], selector.Select())
}

func TestRoundRobinNodeSelector_None(t *testing.T) {
	t.Parallel()

	var nodes []evmclient.Node

	for i := 0; i < 3; i++ {
		node := evmmocks.NewNode(t)
		if i == 0 {
			// first node is out of sync
			node.On("State").Return(evmclient.NodeStateOutOfSync)
		} else {
			// others are unreachable
			node.On("State").Return(evmclient.NodeStateUnreachable)
		}
		nodes = append(nodes, node)
	}

	selector := evmclient.NewRoundRobinSelector(nodes)
	assert.Nil(t, selector.Select())
}

func TestCommonRoundRobinNodeSelector(t *testing.T) {
	t.Parallel()

	var nodes []commonclient.Node[*big.Int, *evmtypes.Head, evmRPC]

	for i := 0; i < 3; i++ {
		node := commonmocks.NewNode[*big.Int, *evmtypes.Head, evmRPC](t)
		if i == 0 {
			// first node is out of sync
			node.On("State").Return(commonclient.NodeStateOutOfSync)
		} else {
			// second & third nodes are alive
			node.On("State").Return(commonclient.NodeStateAlive)
		}
		nodes = append(nodes, node)
	}

	selector := commonclient.NewRoundRobinSelector[*big.Int, *evmtypes.Head, evmRPC](nodes)
	assert.Same(t, nodes[1], selector.Select())
	assert.Same(t, nodes[2], selector.Select())
	assert.Same(t, nodes[1], selector.Select())
	assert.Same(t, nodes[2], selector.Select())
}

func TestCommonRoundRobinNodeSelector_None(t *testing.T) {
	t.Parallel()

	var nodes []commonclient.Node[*big.Int, *evmtypes.Head, evmRPC]

	for i := 0; i < 3; i++ {
		node := commonmocks.NewNode[*big.Int, *evmtypes.Head, evmRPC](t)
		if i == 0 {
			// first node is out of sync
			node.On("State").Return(commonclient.NodeStateOutOfSync)
		} else {
			// others are unreachable
			node.On("State").Return(commonclient.NodeStateUnreachable)
		}
		nodes = append(nodes, node)
	}

	selector := commonclient.NewRoundRobinSelector(nodes)
	assert.Nil(t, selector.Select())
}
