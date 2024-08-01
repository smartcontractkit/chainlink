package graph

import (
	"fmt"
)

// Union combines two given graphs into a new graph. The vertex hashes in both
// graphs are expected to be unique. The two input graphs will remain unchanged.
//
// Both graphs should be either directed or undirected. All traits for the new
// graph will be derived from g.
func Union[K comparable, T any](g, h Graph[K, T]) (Graph[K, T], error) {
	union, err := g.Clone()
	if err != nil {
		return union, fmt.Errorf("failed to clone g: %w", err)
	}

	adjacencyMap, err := h.AdjacencyMap()
	if err != nil {
		return union, fmt.Errorf("failed to get adjacency map: %w", err)
	}

	addedEdges := make(map[K]map[K]struct{})

	for currentHash := range adjacencyMap {
		vertex, properties, err := h.VertexWithProperties(currentHash) //nolint:govet
		if err != nil {
			return union, fmt.Errorf("failed to get vertex %v: %w", currentHash, err)
		}

		err = union.AddVertex(vertex, copyVertexProperties(properties))
		if err != nil {
			return union, fmt.Errorf("failed to add vertex %v: %w", currentHash, err)
		}
	}

	for _, adjacencies := range adjacencyMap {
		for _, edge := range adjacencies {
			if _, sourceOK := addedEdges[edge.Source]; sourceOK {
				if _, targetOK := addedEdges[edge.Source][edge.Target]; targetOK {
					// If the edge addedEdges[source][target] exists, the edge
					// has already been created and thus can be skipped here.
					continue
				}
			}

			err = union.AddEdge(copyEdge(edge))
			if err != nil {
				return union, fmt.Errorf("failed to add edge (%v, %v): %w", edge.Source, edge.Target, err)
			}

			if _, ok := addedEdges[edge.Source]; !ok {
				addedEdges[edge.Source] = make(map[K]struct{})
			}
			addedEdges[edge.Source][edge.Target] = struct{}{}
		}
	}

	return union, nil
}

// unionFind implements a union-find or disjoint set data structure that works
// with vertex hashes as vertices. It's an internal helper type at the moment,
// but could perhaps be exposed publicly in the future.
//
// unionFind is not related to the Union function.
type unionFind[K comparable] struct {
	parents map[K]K
}

func newUnionFind[K comparable](vertices ...K) *unionFind[K] {
	u := &unionFind[K]{
		parents: make(map[K]K, len(vertices)),
	}

	for _, vertex := range vertices {
		u.parents[vertex] = vertex
	}

	return u
}

func (u *unionFind[K]) add(vertex K) {
	u.parents[vertex] = vertex
}

func (u *unionFind[K]) union(vertex1, vertex2 K) {
	root1 := u.find(vertex1)
	root2 := u.find(vertex2)

	if root1 == root2 {
		return
	}

	u.parents[root2] = root1
}

func (u *unionFind[K]) find(vertex K) K {
	root := vertex

	for u.parents[root] != root {
		root = u.parents[root]
	}

	// Perform a path compression in order to optimize of future find calls.
	current := vertex

	for u.parents[current] != root {
		parent := u.parents[vertex]
		u.parents[vertex] = root
		current = parent
	}

	return root
}

func copyVertexProperties(source VertexProperties) func(*VertexProperties) {
	return func(p *VertexProperties) {
		for k, v := range source.Attributes {
			p.Attributes[k] = v
		}
		p.Weight = source.Weight
	}
}
