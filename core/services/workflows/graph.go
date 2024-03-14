package workflows

import "fmt"

type graph[T any] struct {
	// ref -> refs
	adjacencies map[string]map[string]struct{}
	// ref -> nodes
	nodes map[string]T
}

func (g *graph[T]) walkDo(startingRef string, f func(n T) error) error {
	nodesToVisit := []string{startingRef}
	for adj := range g.adjacencies[startingRef] {
		nodesToVisit = append(nodesToVisit, adj)
	}

	visited := map[string]struct{}{}
	for {
		if len(nodesToVisit) == 0 {
			return nil
		}

		curr := nodesToVisit[0]
		nodesToVisit = nodesToVisit[1:]
		if _, found := visited[curr]; found {
			continue
		}

		n, ok := g.nodes[curr]
		if !ok {
			return fmt.Errorf("could not find node with ref %s", curr)
		}
		visited[curr] = struct{}{}

		for adj := range g.adjacencies[curr] {
			nodesToVisit = append(nodesToVisit, adj)
		}

		err := f(n)
		if err != nil {
			return err
		}
	}
}

func (g *graph[T]) adjacentNodes(ref string) []T {
	refs, ok := g.adjacencies[ref]
	if !ok {
		return []T{}
	}

	nodes := []T{}
	for adjacent := range refs {
		n, ok := g.nodes[adjacent]
		if ok {
			nodes = append(nodes, n)
		}
	}

	return nodes
}
