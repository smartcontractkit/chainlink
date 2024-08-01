package graph

import (
	"errors"
	"fmt"
	"sort"
)

// MinimumSpanningTree returns a minimum spanning tree within the given graph.
//
// The MST contains all vertices from the given graph as well as the required
// edges for building the MST. The original graph remains unchanged.
func MinimumSpanningTree[K comparable, T any](g Graph[K, T]) (Graph[K, T], error) {
	return spanningTree(g, false)
}

// MaximumSpanningTree returns a minimum spanning tree within the given graph.
//
// The MST contains all vertices from the given graph as well as the required
// edges for building the MST. The original graph remains unchanged.
func MaximumSpanningTree[K comparable, T any](g Graph[K, T]) (Graph[K, T], error) {
	return spanningTree(g, true)
}

func spanningTree[K comparable, T any](g Graph[K, T], maximum bool) (Graph[K, T], error) {
	if g.Traits().IsDirected {
		return nil, errors.New("spanning trees can only be determined for undirected graphs")
	}

	adjacencyMap, err := g.AdjacencyMap()
	if err != nil {
		return nil, fmt.Errorf("failed to get adjacency map: %w", err)
	}

	edges := make([]Edge[K], 0)
	subtrees := newUnionFind[K]()

	mst := NewLike(g)

	for v, adjacencies := range adjacencyMap {
		vertex, properties, err := g.VertexWithProperties(v) //nolint:govet
		if err != nil {
			return nil, fmt.Errorf("failed to get vertex %v: %w", v, err)
		}

		err = mst.AddVertex(vertex, copyVertexProperties(properties))
		if err != nil {
			return nil, fmt.Errorf("failed to add vertex %v: %w", v, err)
		}

		subtrees.add(v)

		for _, edge := range adjacencies {
			edges = append(edges, edge)
		}
	}

	if maximum {
		sort.Slice(edges, func(i, j int) bool {
			return edges[i].Properties.Weight > edges[j].Properties.Weight
		})
	} else {
		sort.Slice(edges, func(i, j int) bool {
			return edges[i].Properties.Weight < edges[j].Properties.Weight
		})
	}

	for _, edge := range edges {
		sourceRoot := subtrees.find(edge.Source)
		targetRoot := subtrees.find(edge.Target)

		if sourceRoot != targetRoot {
			subtrees.union(sourceRoot, targetRoot)

			if err = mst.AddEdge(copyEdge(edge)); err != nil {
				return nil, fmt.Errorf("failed to add edge (%v, %v): %w", edge.Source, edge.Target, err)
			}
		}
	}

	return mst, nil
}
