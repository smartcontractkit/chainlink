package workflows

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGraph(t *testing.T) {
	g := &graph[int]{
		adjacencies: map[string]map[string]struct{}{
			"node1": {
				"node2": struct{}{},
				"node3": struct{}{},
			},
			"node2": {
				"node3": struct{}{},
			},
			"node3": {
				"node4": struct{}{},
			},
			"node4": {
				"node5": struct{}{},
			},
		},
		nodes: map[string]int{
			"node1": 1,
			"node2": 2,
			"node3": 3,
			"node4": 4,
			"node5": 5,
		},
	}

	inOrder := []int{}
	err := g.walkDo("node1", func(n int) error {
		inOrder = append(inOrder, n)
		return nil
	})
	require.NoError(t, err)
	expected := []int{
		1, 2, 3, 4, 5,
	}
	assert.ElementsMatch(t, expected, inOrder)

	got := g.adjacentNodes("node1")
	assert.Equal(t, []int{2, 3}, got)
}
