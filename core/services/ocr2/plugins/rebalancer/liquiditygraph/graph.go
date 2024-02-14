package liquiditygraph

import (
	"fmt"
	"math/big"
	"sort"
	"sync"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
)

// LiquidityGraph contains graphs functionality for networks and liquidity.
// Graph operations of the implementations should be thread-safe.
type LiquidityGraph interface {
	// AddNetwork adds a new network to the graph by initializing it as a node and setting the initial liquidity.
	AddNetwork(n models.NetworkSelector, v *big.Int) bool

	// GetNetworks returns the list of all the networks that appear on the graph.
	GetNetworks() []models.NetworkSelector

	// HasNetwork returns true when the provided network exists on the graph.
	HasNetwork(n models.NetworkSelector) bool

	// SetLiquidity sets the liquidity of the provided network.
	SetLiquidity(n models.NetworkSelector, v *big.Int) bool

	// GetLiquidity returns the liquidity of the provided network.
	GetLiquidity(n models.NetworkSelector) (*big.Int, error)

	// AddConnection adds a new directed graph edge.
	AddConnection(from, to models.NetworkSelector) error

	// HasConnection returns true if a connection from/to the provided network exist.
	HasConnection(from, to models.NetworkSelector) bool

	// GetNeighbors returns the neighboring network selectors.
	GetNeighbors(from models.NetworkSelector) ([]models.NetworkSelector, bool)

	// GetEdges returns all the graph edges as a list of source/dest pairs.
	GetEdges() ([]models.Edge, error)

	// IsEmpty returns true when the graph does not contain any network.
	IsEmpty() bool

	// Reset resets the graph to it's empty state.
	Reset()

	// String returns the string representation of the graph.
	String() string
}

type Graph struct {
	networksGraph  map[models.NetworkSelector][]models.NetworkSelector
	networkBalance map[models.NetworkSelector]*big.Int
	mu             *sync.RWMutex
}

func NewGraph() *Graph {
	return &Graph{
		networksGraph:  make(map[models.NetworkSelector][]models.NetworkSelector),
		networkBalance: make(map[models.NetworkSelector]*big.Int),
		mu:             &sync.RWMutex{},
	}
}

func NewGraphFromEdges(edges []models.Edge) (*Graph, error) {
	g := NewGraph()
	for _, edge := range edges {
		g.AddNetwork(edge.Source, big.NewInt(0))
		g.AddNetwork(edge.Dest, big.NewInt(0))
		if err := g.AddConnection(edge.Source, edge.Dest); err != nil {
			return nil, fmt.Errorf("add connection %d -> %d: %w", edge.Source, edge.Dest, err)
		}
	}
	return g, nil
}

func (g *Graph) AddNetwork(n models.NetworkSelector, liq *big.Int) bool {
	if g.HasNetwork(n) {
		return false
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	g.networksGraph[n] = make([]models.NetworkSelector, 0)
	g.networkBalance[n] = liq
	return true
}

func (g *Graph) GetNetworks() []models.NetworkSelector {
	g.mu.RLock()
	defer g.mu.RUnlock()

	networks := make([]models.NetworkSelector, 0, len(g.networksGraph))
	for networkID := range g.networksGraph {
		networks = append(networks, networkID)
	}

	// sort the results for deterministic output
	sort.Slice(networks, func(i, j int) bool { return networks[i] < networks[j] })
	return networks
}

func (g *Graph) HasNetwork(n models.NetworkSelector) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	_, exists := g.networksGraph[n]
	return exists
}

func (g *Graph) SetLiquidity(n models.NetworkSelector, v *big.Int) bool {
	if !g.HasNetwork(n) {
		return false
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	g.networkBalance[n] = v
	return true
}

func (g *Graph) GetLiquidity(n models.NetworkSelector) (*big.Int, error) {
	if !g.HasNetwork(n) {
		return nil, fmt.Errorf("network not found")
	}

	g.mu.RLock()
	defer g.mu.RUnlock()

	w, exists := g.networkBalance[n]
	if !exists {
		return nil, fmt.Errorf("graph internal error, network balance not found")
	}

	return w, nil
}

func (g *Graph) AddConnection(from, to models.NetworkSelector) error {
	if !g.HasNetwork(from) {
		return fmt.Errorf("network %d not found", from)
	}
	if !g.HasNetwork(to) {
		return fmt.Errorf("network %d not found", to)
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	g.networksGraph[from] = append(g.networksGraph[from], to)
	return nil
}

func (g *Graph) HasConnection(from, to models.NetworkSelector) bool {
	if !g.HasNetwork(from) || !g.HasNetwork(to) {
		return false
	}

	g.mu.RLock()
	defer g.mu.RUnlock()

	neibs, exist := g.networksGraph[from]
	if !exist {
		return false
	}

	for _, net := range neibs {
		if net == to {
			return true
		}
	}
	return false
}

func (g *Graph) GetNeighbors(from models.NetworkSelector) ([]models.NetworkSelector, bool) {
	if !g.HasNetwork(from) {
		return nil, false
	}

	g.mu.RLock()
	defer g.mu.RUnlock()

	neibs, exist := g.networksGraph[from]
	if !exist {
		return nil, false
	}

	sort.Slice(neibs, func(i, j int) bool { return neibs[i] < neibs[j] })
	return neibs, exist
}

func (g *Graph) GetEdges() ([]models.Edge, error) {
	edges := make([]models.Edge, 0)
	for _, sourceNet := range g.GetNetworks() {
		destNetworks, ok := g.GetNeighbors(sourceNet)
		if !ok {
			return nil, fmt.Errorf("internal graph error %d not found", sourceNet)
		}
		for _, destNet := range destNetworks {
			edges = append(edges, models.NewEdge(sourceNet, destNet))
		}
	}
	return edges, nil
}

func (g *Graph) IsEmpty() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return len(g.networksGraph) == 0
}

func (g *Graph) Reset() {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.networksGraph = make(map[models.NetworkSelector][]models.NetworkSelector)
	g.networkBalance = make(map[models.NetworkSelector]*big.Int)
}

func (g *Graph) String() string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return fmt.Sprintf("Graph{networksGraph: %+v, networkBalance: %+v}", g.networksGraph, g.networkBalance)
}
