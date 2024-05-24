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
	type testNode struct {
		order int32
		state nodeState
	}
	type testCase struct {
		name   string
		nodes  []testNode
		expect []int // indexes of the nodes expected to be returned by Select
	}

	testCases := []testCase{
		{
			name: "TwoNodesSameOrder: Highest Allowed Order",
			nodes: []testNode{
				{order: 1, state: nodeStateAlive},
				{order: 1, state: nodeStateAlive},
			},
			expect: []int{0, 1, 0, 1, 0, 1},
		},
		{
			name: "TwoNodesSameOrder: Lowest Allowed Order",
			nodes: []testNode{
				{order: 100, state: nodeStateAlive},
				{order: 100, state: nodeStateAlive},
			},
			expect: []int{0, 1, 0, 1, 0, 1},
		},
		{
			name: "NoneAvailable",
			nodes: []testNode{
				{order: 1, state: nodeStateOutOfSync},
				{order: 1, state: nodeStateUnreachable},
				{order: 1, state: nodeStateUnreachable},
			},
			expect: []int{}, // no nodes should be selected
		},
		{
			name: "DifferentOrder",
			nodes: []testNode{
				{order: 1, state: nodeStateAlive},
				{order: 2, state: nodeStateAlive},
				{order: 3, state: nodeStateAlive},
			},
			expect: []int{0, 0}, // only the highest order node should be selected
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var nodes []Node[types.ID, Head, nodeClient]
			for _, tn := range tc.nodes {
				node := newMockNode[types.ID, Head, nodeClient](t)
				node.On("State").Return(tn.state)
				node.On("Order").Return(tn.order)
				nodes = append(nodes, node)
			}

			selector := newNodeSelector(NodeSelectionModePriorityLevel, nodes)
			for _, idx := range tc.expect {
				if idx >= len(nodes) {
					t.Fatalf("Invalid node index %d in test case '%s'", idx, tc.name)
				}
				assert.Same(t, nodes[idx], selector.Select())
			}

			// Check for nil selection if expected slice is empty
			if len(tc.expect) == 0 {
				assert.Nil(t, selector.Select())
			}
		})
	}
}
