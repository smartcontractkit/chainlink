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
	AddNetwork(n models.NetworkID, v *big.Int) bool

	// GetNetworks returns the list of all the networks that appear on the graph.
	GetNetworks() []models.NetworkID

	// HasNetwork returns true when the provided network exists on the graph.
	HasNetwork(n models.NetworkID) bool

	// SetLiquidity sets the liquidity of the provided network.
	SetLiquidity(n models.NetworkID, v *big.Int) bool

	// GetLiquidity returns the liquidity of the provided network.
	GetLiquidity(n models.NetworkID) (*big.Int, error)

	// AddConnection adds a new directed graph edge.
	AddConnection(from, to models.NetworkID) bool

	// IsEmpty returns true when the graph does not contain any network.
	IsEmpty() bool

	// Reset resets the graph to it's empty state.
	Reset()
}

type Graph struct {
	networksGraph  map[models.NetworkID][]models.NetworkID
	networkBalance map[models.NetworkID]*big.Int
	mu             *sync.RWMutex
}

func NewGraph() *Graph {
	return &Graph{
		networksGraph:  make(map[models.NetworkID][]models.NetworkID),
		networkBalance: make(map[models.NetworkID]*big.Int),
		mu:             &sync.RWMutex{},
	}
}

func (g *Graph) AddNetwork(n models.NetworkID, liq *big.Int) bool {
	if g.HasNetwork(n) {
		return false
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	g.networksGraph[n] = make([]models.NetworkID, 0)
	g.networkBalance[n] = liq
	return true
}

func (g *Graph) GetNetworks() []models.NetworkID {
	g.mu.RLock()
	defer g.mu.RUnlock()

	networks := make([]models.NetworkID, 0, len(g.networksGraph))
	for networkID := range g.networksGraph {
		networks = append(networks, networkID)
	}

	// sort the results for deterministic output
	sort.Slice(networks, func(i, j int) bool { return networks[i] < networks[j] })
	return networks
}

func (g *Graph) HasNetwork(n models.NetworkID) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	_, exists := g.networksGraph[n]
	return exists
}

func (g *Graph) SetLiquidity(n models.NetworkID, v *big.Int) bool {
	if !g.HasNetwork(n) {
		return false
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	g.networkBalance[n] = v
	return true
}

func (g *Graph) GetLiquidity(n models.NetworkID) (*big.Int, error) {
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

func (g *Graph) AddConnection(from, to models.NetworkID) bool {
	if !g.HasNetwork(from) || !g.HasNetwork(to) {
		return false
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	g.networksGraph[from] = append(g.networksGraph[from], to)
	return true
}

func (g *Graph) IsEmpty() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return len(g.networksGraph) == 0
}

func (g *Graph) Reset() {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.networksGraph = make(map[models.NetworkID][]models.NetworkID)
	g.networkBalance = make(map[models.NetworkID]*big.Int)
}
