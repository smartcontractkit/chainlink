package graph

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

// GraphWriter provides write access to the liquidity graph.
type GraphWriter interface {
	// Add adds new data and connection to the graph.
	Add(from, to Data) error
	// SetLiquidity sets the liquidity of the provided network.
	SetLiquidity(n models.NetworkSelector, liquidity *big.Int) bool
	// SetTargetLiquidity sets the target liquidity of the provided network.
	SetTargetLiquidity(n models.NetworkSelector, liquidity *big.Int) bool
}

// NodeReader provides read access to the data saved in the graph nodes.
type NodeReader interface {
	// GetData returns the data associated with the provided network.
	GetData(n models.NetworkSelector) (Data, error)
	// GetLiquidity returns the liquidity of the provided network.
	GetLiquidity(n models.NetworkSelector) (*big.Int, error)
	// GetTokenAddress returns the token address of the provided network.
	GetTokenAddress(n models.NetworkSelector) (models.Address, error)
	// GetLiquidityManagerAddress returns the liquidity manager address of the provided network.
	GetLiquidityManagerAddress(n models.NetworkSelector) (models.Address, error)
	// GetXChainLiquidityManagerData returns the rebalancer data of the provided network.
	GetXChainLiquidityManagerData(n models.NetworkSelector) (map[models.NetworkSelector]XChainLiquidityManagerData, error)
}

// GraphReader provides read access to the graph data.
type GraphReader interface {
	NodeReader
	// GetNetworks returns the list of all the networks that appear on the graph.
	GetNetworks() []models.NetworkSelector
	// GetNeighbors returns the connected networks.
	GetNeighbors(from models.NetworkSelector, bidir bool) ([]models.NetworkSelector, bool)
	// GetEdges returns all the graph edges as a list of source/dest pairs.
	GetEdges() ([]models.Edge, error)
	// IsEmpty returns true when the graph does not contain any network.
	IsEmpty() bool
	// Len returns the number of vertices in the graph.
	Len() int
	// FindPath returns the path from the source to the destination network.
	// The iterator function is called for each node in the path with the data of the node.
	FindPath(from, to models.NetworkSelector, maxEdgesTraversed int, iterator func(nodes ...Data) bool) []models.NetworkSelector
}

// Graph contains graphs functionality for networks and liquidity
// data for a single token.
type Graph interface {
	GraphWriter
	GraphReader
	// Equals returns true if and only if the provided graph is equal to the receiver graph.
	Equals(other Graph) bool
	// String returns the string representation of the graph.
	String() string
	// Reset resets the graph to it's empty state.
	Reset()
	// Clone creates a deep copy of the graph.
	Clone() Graph
}

// GraphTest provides testing functionality for the graph.
type GraphTest interface {
	// AddNetwork adds a new network to the graph by initializing it as a node
	// and setting the initial data.
	AddNetwork(n models.NetworkSelector, data Data) bool
	// AddConnection adds a new directed graph edge.
	AddConnection(from, to models.NetworkSelector) error
	// HasNetwork returns true when the provided network exists on the graph.
	HasNetwork(n models.NetworkSelector) bool
}

var _ GraphTest = &liquidityGraph{}

func NewGraph() Graph {
	return &liquidityGraph{
		adj:  make(map[models.NetworkSelector][]models.NetworkSelector),
		data: make(map[models.NetworkSelector]Data),
	}
}

type liquidityGraph struct {
	adj  map[models.NetworkSelector][]models.NetworkSelector
	data map[models.NetworkSelector]Data
	lock sync.RWMutex
}

func (g *liquidityGraph) Equals(o Graph) bool {
	g.lock.RLock()
	defer g.lock.RUnlock()

	other := o.(*liquidityGraph)
	other.lock.RLock()
	defer other.lock.RUnlock()

	if g.len() != other.len() {
		return false
	}

	for _, n := range g.getNetworks() {
		if !other.hasNetwork(n) {
			return false
		}
		otherData, exist := other.getData(n)
		if !exist {
			return false
		}
		data, exist := g.getData(n)
		if !exist {
			return false
		}
		if !otherData.Equals(data) {
			return false
		}
		neibs, exists := g.getNeighbors(n, false)
		if !exists {
			return false
		}
		otherNeibs, exists := other.getNeighbors(n, false)
		if !exists {
			return false
		}
		if len(neibs) != len(otherNeibs) {
			return false
		}
		for i, neib := range neibs {
			if neib != otherNeibs[i] {
				return false
			}
		}
	}
	return true
}

func (g *liquidityGraph) String() string {
	g.lock.RLock()
	defer g.lock.RUnlock()

	type network struct {
		Selector models.NetworkSelector
		ChainID  uint64
	}
	adj := make([]network, 0, len(g.adj))
	for n := range g.adj {
		adj = append(adj, network{Selector: n, ChainID: n.ChainID()})
	}
	data := make(map[network]Data, len(g.data))
	for n, d := range g.data {
		data[network{Selector: n, ChainID: n.ChainID()}] = d
	}

	return fmt.Sprintf("Graph{graph: %+v, data: %+v}", adj, data)
}

func (g *liquidityGraph) Reset() {
	g.lock.Lock()
	defer g.lock.Unlock()

	g.adj = make(map[models.NetworkSelector][]models.NetworkSelector)
	g.data = make(map[models.NetworkSelector]Data)
}

func (g *liquidityGraph) Clone() Graph {
	g.lock.RLock()
	defer g.lock.RUnlock()

	clone := &liquidityGraph{
		adj:  make(map[models.NetworkSelector][]models.NetworkSelector, len(g.adj)),
		data: make(map[models.NetworkSelector]Data, len(g.data)),
	}

	for k, v := range g.adj {
		adjCopy := make([]models.NetworkSelector, len(v))
		copy(adjCopy, v)
		clone.adj[k] = adjCopy
	}

	for k, v := range g.data {
		clone.data[k] = v.Clone()
	}

	return clone
}
