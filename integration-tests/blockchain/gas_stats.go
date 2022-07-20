package blockchain

import (
	"sort"

	"github.com/rs/zerolog/log"
)

const (
	// GWei one giga-wei used for gas calculations
	GWei = 1e9
	// ETH one eth in wei
	ETH = 1e18
)

// GasStats helper struct to determine gas metrics across all txs of a test
type GasStats struct {
	NodeID       int
	TotalGasUsed int64
	SeenHashes   map[string]bool
	ClientTXs    []TXGasData
}

// TXGasData transaction gas data
type TXGasData struct {
	TXHash            string
	Value             uint64
	GasLimit          uint64
	GasUsed           uint64
	GasPrice          uint64
	CumulativeGasUsed uint64
}

// NewGasStats creates new gas stats collector
func NewGasStats(nodeID int) *GasStats {
	return &GasStats{
		NodeID:     nodeID,
		SeenHashes: make(map[string]bool),
		ClientTXs:  make([]TXGasData, 0),
	}
}

// AddClientTXData adds client TX data
func (g *GasStats) AddClientTXData(data TXGasData) {
	if _, ok := g.SeenHashes[data.TXHash]; !ok {
		g.ClientTXs = append(g.ClientTXs, data)
	}
}

func (g *GasStats) txCost(data TXGasData) float64 {
	fee := float64(data.GasPrice * data.GasUsed)
	val := float64(data.Value)
	return (fee + val) / ETH
}

func (g *GasStats) maxGasUsage() uint64 {
	if len(g.ClientTXs) == 0 {
		return 0
	}
	sort.Slice(g.ClientTXs, func(i, j int) bool {
		return g.ClientTXs[i].GasUsed > g.ClientTXs[j].GasUsed
	})
	return g.ClientTXs[0].GasUsed
}

func (g *GasStats) totalCost() float64 {
	var total float64
	for _, tx := range g.ClientTXs {
		total += g.txCost(tx)
	}
	return total
}

// PrintStats prints gas stats and total TXs cost
func (g *GasStats) PrintStats() {
	log.Info().Msg("---------- Start Gas Stats ----------")
	log.Info().Int("Node", g.NodeID).Uint64("Gas (GWei)", g.maxGasUsage()).Msg("Max gas used")
	for _, tx := range g.ClientTXs {
		log.Info().
			Int("Node", g.NodeID).
			Float64("Value (ETH)", float64(tx.Value)/ETH).
			Uint64("Gas used", tx.GasUsed).
			Float64("Suggested gas price (GWei)", float64(tx.GasPrice)/GWei).
			Uint64("Gas Limit", tx.GasLimit).
			Float64("Cost (ETH)", g.txCost(tx)).
			Msg("TX Cost")
	}
	log.Info().
		Float64("ETH", g.totalCost()).
		Msg("Total TXs cost")
	log.Info().Msg("---------------- End ---------------")
}
