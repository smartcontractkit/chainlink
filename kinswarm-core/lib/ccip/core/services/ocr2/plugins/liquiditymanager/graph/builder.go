package graph

import (
	"context"
	"fmt"

	mapset "github.com/deckarep/golang-set/v2"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

func NewGraphFromEdges(edges []models.Edge) (Graph, error) {
	g := NewGraph().(*liquidityGraph)
	for _, edge := range edges {
		_ = g.addNetwork(edge.Source, Data{NetworkSelector: edge.Source})
		_ = g.addNetwork(edge.Dest, Data{NetworkSelector: edge.Dest})
		if err := g.addConnection(edge.Source, edge.Dest); err != nil {
			return nil, fmt.Errorf("add connection %d -> %d: %w", edge.Source, edge.Dest, err)
		}
	}
	return g, nil
}

type DataGetter func(
	ctx context.Context,
	v Vertex,
) (Data, []Vertex, error)

func NewGraphWithData(ctx context.Context, start Vertex, dataGetter DataGetter) (Graph, error) {
	g := NewGraph().(*liquidityGraph)

	seen := mapset.NewSet[Vertex]()
	queue := mapset.NewSet[Vertex]()

	queue.Add(start)
	seen.Add(start)

	for queue.Cardinality() > 0 {
		v, ok := queue.Pop()
		if !ok {
			return nil, fmt.Errorf("unexpected internal error")
		}

		val, neighbors, err := dataGetter(ctx, v)
		if err != nil {
			return nil, fmt.Errorf("could not get value for vertex(selector=%d;addr:%s): %w", v.NetworkSelector, v.LiquidityManager.String(), err)
		}
		g.addNetwork(v.NetworkSelector, val)

		for _, neighbor := range neighbors {
			if !g.hasNetwork(neighbor.NetworkSelector) {
				val2, _, err := dataGetter(ctx, neighbor)
				if err != nil {
					return nil, fmt.Errorf("could not get value for neighbor vertex(selector=%d;addr:%s): %w", v.NetworkSelector, v.LiquidityManager.String(), err)
				}
				g.addNetwork(neighbor.NetworkSelector, val2)
			}

			if err := g.addConnection(v.NetworkSelector, neighbor.NetworkSelector); err != nil {
				return nil, fmt.Errorf("error adding connection from %+v to %+v: %w", v, neighbor, err)
			}

			if !seen.Contains(neighbor) {
				queue.Add(neighbor)
				seen.Add(neighbor)
			}
		}
	}

	return g, nil
}
