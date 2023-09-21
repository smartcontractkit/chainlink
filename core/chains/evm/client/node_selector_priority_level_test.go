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

func TestCommonPriorityLevelNodeSelectorName(t *testing.T) {
	selector := commonclient.NewPriorityLevelNodeSelector[*big.Int, *evmtypes.Head, evmRPC](nil)
	assert.Equal(t, selector.Name(), commonclient.NodeSelectionMode_PriorityLevel)
}

func TestCommonPriorityLevelNodeSelector(t *testing.T) {
	t.Parallel()

	var nodes []commonclient.Node[*big.Int, *evmtypes.Head, evmRPC]
	n1 := commonmocks.NewNode[*big.Int, *evmtypes.Head, evmRPC](t)
	n1.On("State").Return(commonclient.NodeStateAlive)
	n1.On("Order").Return(int32(1))

	n2 := commonmocks.NewNode[*big.Int, *evmtypes.Head, evmRPC](t)
	n2.On("State").Return(commonclient.NodeStateAlive)
	n2.On("Order").Return(int32(1))

	n3 := commonmocks.NewNode[*big.Int, *evmtypes.Head, evmRPC](t)
	n3.On("State").Return(commonclient.NodeStateAlive)
	n3.On("Order").Return(int32(1))

	nodes = append(nodes, n1, n2, n3)
	selector := commonclient.NewPriorityLevelNodeSelector(nodes)
	assert.Same(t, nodes[0], selector.Select())
	assert.Same(t, nodes[1], selector.Select())
	assert.Same(t, nodes[2], selector.Select())
	assert.Same(t, nodes[0], selector.Select())
	assert.Same(t, nodes[1], selector.Select())
	assert.Same(t, nodes[2], selector.Select())
}

func TestCommonPriorityLevelNodeSelector_None(t *testing.T) {
	t.Parallel()

	var nodes []commonclient.Node[*big.Int, *evmtypes.Head, evmRPC]

	for i := 0; i < 3; i++ {
		node := commonmocks.NewNode[*big.Int, *evmtypes.Head, evmRPC](t)
		if i == 0 {
			// first node is out of sync
			node.On("State").Return(commonclient.NodeStateOutOfSync)
			node.On("Order").Return(int32(1))
		} else {
			// others are unreachable
			node.On("State").Return(commonclient.NodeStateUnreachable)
			node.On("Order").Return(int32(1))
		}
		nodes = append(nodes, node)
	}

	selector := commonclient.NewPriorityLevelNodeSelector(nodes)
	assert.Nil(t, selector.Select())
}

func TestCommonPriorityLevelNodeSelector_DifferentOrder(t *testing.T) {
	t.Parallel()

	var nodes []commonclient.Node[*big.Int, *evmtypes.Head, evmRPC]

	n1 := commonmocks.NewNode[*big.Int, *evmtypes.Head, evmRPC](t)
	n1.On("State").Return(commonclient.NodeStateAlive)
	n1.On("Order").Return(int32(1))

	n2 := commonmocks.NewNode[*big.Int, *evmtypes.Head, evmRPC](t)
	n2.On("State").Return(commonclient.NodeStateAlive)
	n2.On("Order").Return(int32(2))

	n3 := commonmocks.NewNode[*big.Int, *evmtypes.Head, evmRPC](t)
	n3.On("State").Return(commonclient.NodeStateAlive)
	n3.On("Order").Return(int32(3))

	nodes = append(nodes, n1, n2, n3)
	selector := commonclient.NewPriorityLevelNodeSelector(nodes)
	assert.Same(t, nodes[0], selector.Select())
	assert.Same(t, nodes[0], selector.Select())
}
