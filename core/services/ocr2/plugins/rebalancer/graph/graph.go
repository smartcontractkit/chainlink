package graph

import (
	"fmt"
	"math/big"
	"sort"
	"sync"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
)

// Graph contains graphs functionality for networks and liquidity.
// Graph operations of the implementations should be thread-safe.
type Graph interface {
	// AddNetwork adds a new network to the graph by initializing it as a node
	// and setting the initial data.
	AddNetwork(n models.NetworkSelector, data Data) bool

	// GetNetworks returns the list of all the networks that appear on the graph.
	GetNetworks() []models.NetworkSelector

	// HasNetwork returns true when the provided network exists on the graph.
	HasNetwork(n models.NetworkSelector) bool

	// SetLiquidity sets the liquidity of the provided network.
	SetLiquidity(n models.NetworkSelector, liquidity *big.Int) bool

	// GetLiquidity returns the liquidity of the provided network.
	GetLiquidity(n models.NetworkSelector) (*big.Int, error)

	GetTokenAddress(n models.NetworkSelector) (models.Address, error)

	GetRebalancerAddress(n models.NetworkSelector) (models.Address, error)

	GetXChainRebalancerData(n models.NetworkSelector) (map[models.NetworkSelector]XChainRebalancerData, error)

	// GetData returns the data of the provided network.
	// TODO: remove redundant methods, e.g. GetLiquidity can be replaced with GetData
	GetData(n models.NetworkSelector) (Data, error)

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

	Len() int
}

func NewGraphFromEdges(edges []models.Edge) (Graph, error) {
	g := NewGraph()
	for _, edge := range edges {
		g.AddNetwork(edge.Source, Data{})
		g.AddNetwork(edge.Dest, Data{})
		if err := g.AddConnection(edge.Source, edge.Dest); err != nil {
			return nil, fmt.Errorf("add connection %d -> %d: %w", edge.Source, edge.Dest, err)
		}
	}
	return g, nil
}

type XChainRebalancerData struct {
	RemoteRebalancerAddress   models.Address
	LocalBridgeAdapterAddress models.Address
	RemoteTokenAddress        models.Address
}

type Data struct {
	Liquidity         *big.Int
	TokenAddress      models.Address
	RebalancerAddress models.Address
	XChainRebalancers map[models.NetworkSelector]XChainRebalancerData
	ConfigDigest      models.ConfigDigest
	NetworkSelector   models.NetworkSelector
}

type gph struct {
	adj  map[models.NetworkSelector][]models.NetworkSelector
	data map[models.NetworkSelector]Data
	mu   *sync.RWMutex
}

func NewGraph() Graph {
	return &gph{
		adj:  make(map[models.NetworkSelector][]models.NetworkSelector),
		data: make(map[models.NetworkSelector]Data),
		mu:   &sync.RWMutex{},
	}
}

func (g *gph) AddNetwork(n models.NetworkSelector, data Data) bool {
	if g.HasNetwork(n) {
		return false
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	g.adj[n] = make([]models.NetworkSelector, 0)
	g.data[n] = data
	return true
}

func (g *gph) GetNetworks() []models.NetworkSelector {
	g.mu.RLock()
	defer g.mu.RUnlock()

	networks := make([]models.NetworkSelector, 0, len(g.adj))
	for networkID := range g.adj {
		networks = append(networks, networkID)
	}

	// sort the results for deterministic output
	sort.Slice(networks, func(i, j int) bool { return networks[i] < networks[j] })
	return networks
}

func (g *gph) HasNetwork(n models.NetworkSelector) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	_, exists := g.adj[n]
	return exists
}

func (g *gph) SetLiquidity(n models.NetworkSelector, liquidity *big.Int) bool {
	if !g.HasNetwork(n) {
		return false
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	prev := g.data[n]
	g.data[n] = Data{
		Liquidity:         liquidity,
		TokenAddress:      prev.TokenAddress,
		RebalancerAddress: prev.RebalancerAddress,
		ConfigDigest:      prev.ConfigDigest,
		NetworkSelector:   prev.NetworkSelector,
	}
	return true
}

func (g *gph) GetLiquidity(n models.NetworkSelector) (*big.Int, error) {
	if !g.HasNetwork(n) {
		return nil, fmt.Errorf("network not found")
	}

	g.mu.RLock()
	defer g.mu.RUnlock()

	w, exists := g.data[n]
	if !exists {
		return nil, fmt.Errorf("graph internal error, network balance not found")
	}

	return w.Liquidity, nil
}

func (g *gph) GetData(n models.NetworkSelector) (Data, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	data, exists := g.data[n]
	if !exists {
		return Data{}, fmt.Errorf("network %d not found", n)
	}
	return data, nil
}

func (g *gph) GetTokenAddress(n models.NetworkSelector) (models.Address, error) {
	if !g.HasNetwork(n) {
		return models.Address{}, fmt.Errorf("network not found")
	}

	g.mu.RLock()
	defer g.mu.RUnlock()

	w, exists := g.data[n]
	if !exists {
		return models.Address{}, fmt.Errorf("graph internal error, network balance not found")
	}

	return w.TokenAddress, nil
}

func (g *gph) GetRebalancerAddress(n models.NetworkSelector) (models.Address, error) {
	if !g.HasNetwork(n) {
		return models.Address{}, fmt.Errorf("network not found")
	}

	g.mu.RLock()
	defer g.mu.RUnlock()

	w, exists := g.data[n]
	if !exists {
		return models.Address{}, fmt.Errorf("graph internal error, network balance not found")
	}

	return w.RebalancerAddress, nil
}

func (g *gph) AddConnection(from, to models.NetworkSelector) error {
	if !g.HasNetwork(from) {
		return fmt.Errorf("network %d not found", from)
	}
	if !g.HasNetwork(to) {
		return fmt.Errorf("network %d not found", to)
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	g.adj[from] = append(g.adj[from], to)
	return nil
}

func (g *gph) HasConnection(from, to models.NetworkSelector) bool {
	if !g.HasNetwork(from) || !g.HasNetwork(to) {
		return false
	}

	g.mu.RLock()
	defer g.mu.RUnlock()

	neibs, exist := g.adj[from]
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

func (g *gph) GetNeighbors(from models.NetworkSelector) ([]models.NetworkSelector, bool) {
	if !g.HasNetwork(from) {
		return nil, false
	}

	g.mu.RLock()
	defer g.mu.RUnlock()

	neibs, exist := g.adj[from]
	if !exist {
		return nil, false
	}

	sort.Slice(neibs, func(i, j int) bool { return neibs[i] < neibs[j] })
	return neibs, exist
}

func (g *gph) GetEdges() ([]models.Edge, error) {
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

func (g *gph) IsEmpty() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return len(g.adj) == 0
}

func (g *gph) Reset() {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.adj = make(map[models.NetworkSelector][]models.NetworkSelector)
	g.data = make(map[models.NetworkSelector]Data)
}

func (g *gph) String() string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return fmt.Sprintf("Graph{networksGraph: %+v, networkBalance: %+v}", g.adj, g.data)
}

func (g *gph) Len() int {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return len(g.adj)
}

// GetXChainRebalancerData implements Graph.
func (g *gph) GetXChainRebalancerData(n models.NetworkSelector) (map[models.NetworkSelector]XChainRebalancerData, error) {
	if !g.HasNetwork(n) {
		return nil, fmt.Errorf("network not found")
	}

	g.mu.RLock()
	defer g.mu.RUnlock()

	w, exists := g.data[n]
	if !exists {
		return nil, fmt.Errorf("graph internal error, network balance not found")
	}

	return w.XChainRebalancers, nil
}
