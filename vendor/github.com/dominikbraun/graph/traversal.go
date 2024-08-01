package graph

import "fmt"

// DFS performs a depth-first search on the graph, starting from the given vertex. The visit
// function will be invoked with the hash of the vertex currently visited. If it returns false, DFS
// will continue traversing the graph, and if it returns true, the traversal will be stopped. In
// case the graph is disconnected, only the vertices joined with the starting vertex are visited.
//
// This example prints all vertices of the graph in DFS-order:
//
//	g := graph.New(graph.IntHash)
//
//	_ = g.AddVertex(1)
//	_ = g.AddVertex(2)
//	_ = g.AddVertex(3)
//
//	_ = g.AddEdge(1, 2)
//	_ = g.AddEdge(2, 3)
//	_ = g.AddEdge(3, 1)
//
//	_ = graph.DFS(g, 1, func(value int) bool {
//		fmt.Println(value)
//		return false
//	})
//
// Similarly, if you have a graph of City vertices and the traversal should stop at London, the
// visit function would look as follows:
//
//	func(c City) bool {
//		return c.Name == "London"
//	}
//
// DFS is non-recursive and maintains a stack instead.
func DFS[K comparable, T any](g Graph[K, T], start K, visit func(K) bool) error {
	adjacencyMap, err := g.AdjacencyMap()
	if err != nil {
		return fmt.Errorf("could not get adjacency map: %w", err)
	}

	if _, ok := adjacencyMap[start]; !ok {
		return fmt.Errorf("could not find start vertex with hash %v", start)
	}

	stack := make([]K, 0)
	visited := make(map[K]bool)

	stack = append(stack, start)

	for len(stack) > 0 {
		currentHash := stack[len(stack)-1]

		stack = stack[:len(stack)-1]

		if _, ok := visited[currentHash]; !ok {
			// Stop traversing the graph if the visit function returns true.
			if stop := visit(currentHash); stop {
				break
			}
			visited[currentHash] = true

			for adjacency := range adjacencyMap[currentHash] {
				stack = append(stack, adjacency)
			}
		}
	}

	return nil
}

// BFS performs a breadth-first search on the graph, starting from the given vertex. The visit
// function will be invoked with the hash of the vertex currently visited. If it returns false, BFS
// will continue traversing the graph, and if it returns true, the traversal will be stopped. In
// case the graph is disconnected, only the vertices joined with the starting vertex are visited.
//
// This example prints all vertices of the graph in BFS-order:
//
//	g := graph.New(graph.IntHash)
//
//	_ = g.AddVertex(1)
//	_ = g.AddVertex(2)
//	_ = g.AddVertex(3)
//
//	_ = g.AddEdge(1, 2)
//	_ = g.AddEdge(2, 3)
//	_ = g.AddEdge(3, 1)
//
//	_ = graph.BFS(g, 1, func(value int) bool {
//		fmt.Println(value)
//		return false
//	})
//
// Similarly, if you have a graph of City vertices and the traversal should stop at London, the
// visit function would look as follows:
//
//	func(c City) bool {
//		return c.Name == "London"
//	}
//
// BFS is non-recursive and maintains a stack instead.
func BFS[K comparable, T any](g Graph[K, T], start K, visit func(K) bool) error {
	ignoreDepth := func(vertex K, _ int) bool {
		return visit(vertex)
	}
	return BFSWithDepth(g, start, ignoreDepth)
}

// BFSWithDepth works just as BFS and performs a breadth-first search on the graph, but its
// visit function is passed the current depth level as a second argument. Consequently, the
// current depth can be used for deciding whether or not to proceed past a certain depth.
//
//	_ = graph.BFSWithDepth(g, 1, func(value int, depth int) bool {
//		fmt.Println(value)
//		return depth > 3
//	})
//
// With the visit function from the example, the BFS traversal will stop once a depth greater
// than 3 is reached.
func BFSWithDepth[K comparable, T any](g Graph[K, T], start K, visit func(K, int) bool) error {
	adjacencyMap, err := g.AdjacencyMap()
	if err != nil {
		return fmt.Errorf("could not get adjacency map: %w", err)
	}

	if _, ok := adjacencyMap[start]; !ok {
		return fmt.Errorf("could not find start vertex with hash %v", start)
	}

	queue := make([]K, 0)
	visited := make(map[K]bool)

	visited[start] = true
	queue = append(queue, start)
	depth := 0

	for len(queue) > 0 {
		currentHash := queue[0]

		queue = queue[1:]
		depth++

		// Stop traversing the graph if the visit function returns true.
		if stop := visit(currentHash, depth); stop {
			break
		}

		for adjacency := range adjacencyMap[currentHash] {
			if _, ok := visited[adjacency]; !ok {
				visited[adjacency] = true
				queue = append(queue, adjacency)
			}
		}

	}

	return nil
}
