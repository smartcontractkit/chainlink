package graph

import (
	"fmt"
	"math/big"
	"sort"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

func (g *liquidityGraph) GetData(n models.NetworkSelector) (Data, error) {
	g.lock.RLock()
	defer g.lock.RUnlock()

	data, exists := g.getData(n)
	if !exists {
		return Data{}, fmt.Errorf("network %d not found", n)
	}
	return data, nil
}

func (g *liquidityGraph) GetLiquidity(n models.NetworkSelector) (*big.Int, error) {
	g.lock.RLock()
	defer g.lock.RUnlock()

	if !g.hasNetwork(n) {
		return nil, fmt.Errorf("network not found")
	}

	d, exists := g.getData(n)
	if !exists {
		return nil, fmt.Errorf("graph internal error, network not found")
	}

	return d.Liquidity, nil
}

func (g *liquidityGraph) GetTokenAddress(n models.NetworkSelector) (models.Address, error) {
	g.lock.RLock()
	defer g.lock.RUnlock()

	if !g.hasNetwork(n) {
		return models.Address{}, fmt.Errorf("network not found")
	}

	d, exists := g.getData(n)
	if !exists {
		return models.Address{}, fmt.Errorf("graph internal error, network not found")
	}

	return d.TokenAddress, nil
}

func (g *liquidityGraph) GetLiquidityManagerAddress(n models.NetworkSelector) (models.Address, error) {
	g.lock.RLock()
	defer g.lock.RUnlock()

	if !g.hasNetwork(n) {
		return models.Address{}, fmt.Errorf("network not found")
	}

	d, exists := g.getData(n)
	if !exists {
		return models.Address{}, fmt.Errorf("graph internal error, network not found")
	}

	return d.LiquidityManagerAddress, nil
}

func (g *liquidityGraph) GetXChainLiquidityManagerData(n models.NetworkSelector) (map[models.NetworkSelector]XChainLiquidityManagerData, error) {
	g.lock.RLock()
	defer g.lock.RUnlock()

	if !g.hasNetwork(n) {
		return nil, fmt.Errorf("network not found")
	}

	w, exists := g.getData(n)
	if !exists {
		return nil, fmt.Errorf("graph internal error, network balance not found")
	}

	return w.XChainLiquidityManagers, nil
}

func (g *liquidityGraph) HasNetwork(n models.NetworkSelector) bool {
	g.lock.RLock()
	defer g.lock.RUnlock()

	return g.hasNetwork(n)
}

func (g *liquidityGraph) GetNetworks() []models.NetworkSelector {
	g.lock.RLock()
	defer g.lock.RUnlock()

	return g.getNetworks()
}

func (g *liquidityGraph) GetNeighbors(from models.NetworkSelector, bidir bool) ([]models.NetworkSelector, bool) {
	g.lock.RLock()
	defer g.lock.RUnlock()

	return g.getNeighbors(from, bidir)
}

func (g *liquidityGraph) GetEdges() ([]models.Edge, error) {
	g.lock.RLock()
	defer g.lock.RUnlock()

	edges := make([]models.Edge, 0)
	for _, sourceNet := range g.getNetworks() {
		destNetworks, ok := g.getNeighbors(sourceNet, false)
		if !ok {
			return nil, fmt.Errorf("internal graph error %d not found", sourceNet)
		}
		for _, destNet := range destNetworks {
			edges = append(edges, models.NewEdge(sourceNet, destNet))
		}
	}
	return edges, nil
}

func (g *liquidityGraph) IsEmpty() bool {
	g.lock.RLock()
	defer g.lock.RUnlock()

	return g.len() == 0
}

func (g *liquidityGraph) Len() int {
	g.lock.RLock()
	defer g.lock.RUnlock()

	return g.len()
}

// FindPath finds a path from the source network to the destination network with the given number of edges that are allow to be traversed.
// It calls the iterator function with each individual node in the path.
// It returns the list of network selectors representing the path from source to destination (including the destination node).
func (g *liquidityGraph) FindPath(from, to models.NetworkSelector, maxEdgesTraversed int, iterator func(nodes ...Data) bool) []models.NetworkSelector {
	g.lock.RLock()
	defer g.lock.RUnlock()

	return g.findPath(from, to, maxEdgesTraversed, iterator)
}

func (g *liquidityGraph) findPath(from, to models.NetworkSelector, maxEdgesTraversed int, iterator func(nodes ...Data) bool) []models.NetworkSelector {
	if maxEdgesTraversed == 0 {
		return []models.NetworkSelector{}
	}
	neibs, exist := g.adj[from]
	if !exist {
		return []models.NetworkSelector{}
	}
	for _, n := range neibs {
		if n == to {
			if !iterator(g.data[to]) {
				continue
			}
			return []models.NetworkSelector{n}
		}
	}
	for _, n := range neibs {
		if p := g.findPath(n, to, maxEdgesTraversed-1, iterator); len(p) > 0 {
			data := []Data{g.data[n]}
			for _, d := range p {
				data = append(data, g.data[d])
			}
			if !iterator(data...) {
				continue
			}
			return append([]models.NetworkSelector{n}, p...)
		}
	}
	return []models.NetworkSelector{}
}

func (g *liquidityGraph) getData(n models.NetworkSelector) (Data, bool) {
	data, exists := g.data[n]
	return data, exists
}

func (g *liquidityGraph) getNetworks() []models.NetworkSelector {
	networks := make([]models.NetworkSelector, 0, len(g.adj))
	for networkID := range g.adj {
		networks = append(networks, networkID)
	}
	// sort the results for deterministic output
	sort.Slice(networks, func(i, j int) bool { return networks[i] < networks[j] })
	return networks
}

// getNeighbors returns the neighbors of a network, if bidir is true it returns only the bidirectional connections
func (g *liquidityGraph) getNeighbors(from models.NetworkSelector, bidir bool) ([]models.NetworkSelector, bool) {
	if !g.hasNetwork(from) {
		return nil, false
	}

	neibs, exist := g.adj[from]
	if !exist {
		return nil, false
	}

	if bidir {
		bineibs := make([]models.NetworkSelector, 0)
		for _, neib := range neibs {
			if g.hasConnection(neib, from) {
				bineibs = append(bineibs, neib)
			}
		}
		neibs = bineibs
	}
	sort.Slice(neibs, func(i, j int) bool { return neibs[i] < neibs[j] })
	return neibs, exist
}

func (g *liquidityGraph) hasNetwork(n models.NetworkSelector) bool {
	_, exists := g.adj[n]
	return exists
}

func (g *liquidityGraph) hasConnection(from, to models.NetworkSelector) bool {
	if !g.hasNetwork(from) || !g.hasNetwork(to) {
		return false
	}
	neibs, exist := g.adj[from]
	if !exist {
		// has no connections
		return false
	}
	for _, net := range neibs {
		if net == to {
			return true
		}
	}
	return false
}

func (g *liquidityGraph) len() int {
	return len(g.adj)
}
