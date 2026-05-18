package graph

import (
	"fmt"
	"math/big"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

func (g *liquidityGraph) Add(from, to Data) error {
	g.lock.Lock()
	defer g.lock.Unlock()

	_ = g.addNetwork(from.NetworkSelector, from)
	_ = g.addNetwork(to.NetworkSelector, to)

	if err := g.addConnection(from.NetworkSelector, to.NetworkSelector); err != nil {
		return fmt.Errorf("add connection %d -> %d: %w", from.NetworkSelector, to.NetworkSelector, err)
	}
	return nil
}

func (g *liquidityGraph) SetLiquidity(n models.NetworkSelector, liquidity *big.Int) bool {
	g.lock.Lock()
	defer g.lock.Unlock()

	if !g.hasNetwork(n) {
		return false
	}

	prev := g.data[n]
	g.data[n] = Data{
		Liquidity:               liquidity,
		TokenAddress:            prev.TokenAddress,
		LiquidityManagerAddress: prev.LiquidityManagerAddress,
		ConfigDigest:            prev.ConfigDigest,
		NetworkSelector:         prev.NetworkSelector,
		MinimumLiquidity:        prev.MinimumLiquidity,
		TargetLiquidity:         prev.TargetLiquidity,
	}
	return true
}

func (g *liquidityGraph) SetTargetLiquidity(n models.NetworkSelector, target *big.Int) bool {
	g.lock.Lock()
	defer g.lock.Unlock()

	if !g.hasNetwork(n) {
		return false
	}

	prev := g.data[n]
	g.data[n] = Data{
		Liquidity:               prev.Liquidity,
		TokenAddress:            prev.TokenAddress,
		LiquidityManagerAddress: prev.LiquidityManagerAddress,
		ConfigDigest:            prev.ConfigDigest,
		NetworkSelector:         prev.NetworkSelector,
		MinimumLiquidity:        prev.MinimumLiquidity,
		TargetLiquidity:         target,
	}
	return true
}

func (g *liquidityGraph) AddNetwork(n models.NetworkSelector, data Data) bool {
	g.lock.Lock()
	defer g.lock.Unlock()

	return g.addNetwork(n, data)
}

func (g *liquidityGraph) AddConnection(from, to models.NetworkSelector) error {
	g.lock.Lock()
	defer g.lock.Unlock()

	return g.addConnection(from, to)
}

func (g *liquidityGraph) addNetwork(n models.NetworkSelector, data Data) bool {
	if g.hasNetwork(n) {
		return false
	}
	g.adj[n] = make([]models.NetworkSelector, 0)
	g.data[n] = data
	return true
}

func (g *liquidityGraph) addConnection(from, to models.NetworkSelector) error {
	if !g.hasNetwork(from) {
		return fmt.Errorf("network %d not found", from)
	}
	if !g.hasNetwork(to) {
		return fmt.Errorf("network %d not found", to)
	}
	g.adj[from] = append(g.adj[from], to)
	return nil
}
