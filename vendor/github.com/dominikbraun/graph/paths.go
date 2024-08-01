package graph

import (
	"errors"
	"fmt"
	"math"
)

var ErrTargetNotReachable = errors.New("target vertex not reachable from source")

// CreatesCycle determines whether adding an edge between the two given vertices
// would introduce a cycle in the graph. CreatesCycle will not create an edge.
//
// A potential edge would create a cycle if the target vertex is also a parent
// of the source vertex. In order to determine this, CreatesCycle runs a DFS.
func CreatesCycle[K comparable, T any](g Graph[K, T], source, target K) (bool, error) {
	if _, err := g.Vertex(source); err != nil {
		return false, fmt.Errorf("could not get vertex with hash %v: %w", source, err)
	}

	if _, err := g.Vertex(target); err != nil {
		return false, fmt.Errorf("could not get vertex with hash %v: %w", target, err)
	}

	if source == target {
		return true, nil
	}

	predecessorMap, err := g.PredecessorMap()
	if err != nil {
		return false, fmt.Errorf("failed to get predecessor map: %w", err)
	}

	stack := make([]K, 0)
	visited := make(map[K]bool)

	stack = append(stack, source)

	for len(stack) > 0 {
		currentHash := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if _, ok := visited[currentHash]; !ok {
			// If the adjacent vertex also is the target vertex, the target is a
			// parent of the source vertex. An edge would introduce a cycle.
			if currentHash == target {
				return true, nil
			}

			visited[currentHash] = true

			for adjacency := range predecessorMap[currentHash] {
				stack = append(stack, adjacency)
			}
		}
	}

	return false, nil
}

// ShortestPath computes the shortest path between a source and a target vertex
// under consideration of the edge weights. It returns a slice of hash values of
// the vertices forming that path.
//
// The returned path includes the source and target vertices. If the target is
// not reachable from the source, ErrTargetNotReachable will be returned. Should
// there be multiple shortest paths, and arbitrary one will be returned.
//
// ShortestPath has a time complexity of O(|V|+|E|log(|V|)).
func ShortestPath[K comparable, T any](g Graph[K, T], source, target K) ([]K, error) {
	weights := make(map[K]float64)
	visited := make(map[K]bool)

	weights[source] = 0
	visited[target] = true

	queue := newPriorityQueue[K]()
	adjacencyMap, err := g.AdjacencyMap()
	if err != nil {
		return nil, fmt.Errorf("could not get adjacency map: %w", err)
	}

	for hash := range adjacencyMap {
		if hash != source {
			weights[hash] = math.Inf(1)
			visited[hash] = false
		}

		queue.Push(hash, weights[hash])
	}

	// bestPredecessors stores the cheapest or least-weighted predecessor for
	// each vertex. Given an edge AC with weight=4 and an edge BC with weight=2,
	// the cheapest predecessor for C is B.
	bestPredecessors := make(map[K]K)

	for queue.Len() > 0 {
		vertex, _ := queue.Pop()
		hasInfiniteWeight := math.IsInf(weights[vertex], 1)

		for adjacency, edge := range adjacencyMap[vertex] {
			edgeWeight := edge.Properties.Weight

			// Setting the weight to 1 is required for unweighted graphs whose
			// edge weights are 0. Otherwise, all paths would have a sum of 0
			// and a random path would be returned.
			if !g.Traits().IsWeighted {
				edgeWeight = 1
			}

			weight := weights[vertex] + float64(edgeWeight)

			if weight < weights[adjacency] && !hasInfiniteWeight {
				weights[adjacency] = weight
				bestPredecessors[adjacency] = vertex
				queue.UpdatePriority(adjacency, weight)
			}
		}
	}

	path := []K{target}
	current := target

	for current != source {
		// If the current vertex is not present in bestPredecessors, current is
		// set to the zero value of K. Without this check, this would lead to an
		// endless prepending of zero values to the path. Also, the target would
		// not be reachable from one of the preceding vertices.
		if _, ok := bestPredecessors[current]; !ok {
			return nil, ErrTargetNotReachable
		}
		current = bestPredecessors[current]
		path = append([]K{current}, path...)
	}

	return path, nil
}

type sccState[K comparable] struct {
	adjacencyMap map[K]map[K]Edge[K]
	components   [][]K
	stack        []K
	onStack      map[K]bool
	visited      map[K]struct{}
	lowlink      map[K]int
	index        map[K]int
	time         int
}

// StronglyConnectedComponents detects all strongly connected components within
// the graph and returns the hashes of the vertices shaping these components, so
// each component is represented by a []K.
//
// StronglyConnectedComponents can only run on directed graphs.
func StronglyConnectedComponents[K comparable, T any](g Graph[K, T]) ([][]K, error) {
	if !g.Traits().IsDirected {
		return nil, errors.New("SCCs can only be detected in directed graphs")
	}

	adjacencyMap, err := g.AdjacencyMap()
	if err != nil {
		return nil, fmt.Errorf("could not get adjacency map: %w", err)
	}

	state := &sccState[K]{
		adjacencyMap: adjacencyMap,
		components:   make([][]K, 0),
		stack:        make([]K, 0),
		onStack:      make(map[K]bool),
		visited:      make(map[K]struct{}),
		lowlink:      make(map[K]int),
		index:        make(map[K]int),
	}

	for hash := range state.adjacencyMap {
		if _, ok := state.visited[hash]; !ok {
			findSCC(hash, state)
		}
	}

	return state.components, nil
}

func findSCC[K comparable](vertexHash K, state *sccState[K]) {
	state.stack = append(state.stack, vertexHash)
	state.onStack[vertexHash] = true
	state.visited[vertexHash] = struct{}{}
	state.index[vertexHash] = state.time
	state.lowlink[vertexHash] = state.time

	state.time++

	for adjacency := range state.adjacencyMap[vertexHash] {
		if _, ok := state.visited[adjacency]; !ok {
			findSCC(adjacency, state)

			smallestLowlink := math.Min(
				float64(state.lowlink[vertexHash]),
				float64(state.lowlink[adjacency]),
			)
			state.lowlink[vertexHash] = int(smallestLowlink)
		} else {
			// If the adjacent vertex already is on the stack, the edge joining
			// the current and the adjacent vertex is a back ege. Therefore, the
			// lowlink value of the vertex has to be updated to the index of the
			// adjacent vertex if it is smaller than the current lowlink value.
			if state.onStack[adjacency] {
				smallestLowlink := math.Min(
					float64(state.lowlink[vertexHash]),
					float64(state.index[adjacency]),
				)
				state.lowlink[vertexHash] = int(smallestLowlink)
			}
		}
	}

	// If the lowlink value of the vertex is equal to its DFS value, this is the
	// head vertex of a strongly connected component that's shaped by the vertex
	// and all vertices on the stack.
	if state.lowlink[vertexHash] == state.index[vertexHash] {
		var hash K
		var component []K

		for hash != vertexHash {
			hash = state.stack[len(state.stack)-1]
			state.stack = state.stack[:len(state.stack)-1]
			state.onStack[hash] = false

			component = append(component, hash)
		}

		state.components = append(state.components, component)
	}
}

// AllPathsBetween computes and returns all paths between two given vertices. A
// path is represented as a slice of vertex hashes. The returned slice contains
// these paths.
//
// AllPathsBetween utilizes a non-recursive, stack-based implementation. It has
// an estimated runtime complexity of O(n^2) where n is the number of vertices.
func AllPathsBetween[K comparable, T any](g Graph[K, T], start, end K) ([][]K, error) {
	adjacencyMap, err := g.AdjacencyMap()
	if err != nil {
		return nil, err
	}

	// The algorithm used relies on stacks instead of recursion. It is described
	// here: https://boycgit.github.io/all-paths-between-two-vertex/
	mainStack := newStack[K]()
	viceStack := newStack[stack[K]]()

	checkEmpty := func() error {
		if mainStack.isEmpty() || viceStack.isEmpty() {
			return errors.New("empty stack")
		}
		return nil
	}

	buildLayer := func(element K) {
		mainStack.push(element)

		newElements := newStack[K]()
		for e := range adjacencyMap[element] {
			var contains bool
			mainStack.forEach(func(k K) {
				if e == k {
					contains = true
				}
			})
			if contains {
				continue
			}
			newElements.push(e)
		}
		viceStack.push(newElements)
	}

	buildStack := func() error {
		if err = checkEmpty(); err != nil {
			return fmt.Errorf("unable to build stack: %w", err)
		}

		elements, _ := viceStack.top()

		for !elements.isEmpty() {
			element, _ := elements.pop()
			buildLayer(element)
			elements, _ = viceStack.top()
		}

		return nil
	}

	removeLayer := func() error {
		if err = checkEmpty(); err != nil {
			return fmt.Errorf("unable to remove layer: %w", err)
		}

		if e, _ := viceStack.top(); !e.isEmpty() {
			return errors.New("the top element of vice-stack is not empty")
		}

		_, _ = mainStack.pop()
		_, _ = viceStack.pop()

		return nil
	}

	buildLayer(start)

	allPaths := make([][]K, 0)

	for !mainStack.isEmpty() {
		v, _ := mainStack.top()
		adjs, _ := viceStack.top()

		if adjs.isEmpty() {
			if v == end {
				path := make([]K, 0)
				mainStack.forEach(func(k K) {
					path = append(path, k)
				})
				allPaths = append(allPaths, path)
			}

			err = removeLayer()
			if err != nil {
				return nil, err
			}
		} else {
			if err = buildStack(); err != nil {
				return nil, err
			}
		}
	}

	return allPaths, nil
}
