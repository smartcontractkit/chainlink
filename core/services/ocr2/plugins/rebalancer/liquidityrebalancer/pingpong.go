package liquidityrebalancer

import (
	"fmt"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/graph"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
)

// PingPong rebalancer keeps sending tokens between the chains without ever reaching balance. For testing purposes.
type PingPong struct{}

func NewPingPong() *PingPong {
	return &PingPong{}
}

func (p *PingPong) ComputeTransfersToBalance(g graph.Graph, inflightTransfers []models.PendingTransfer) ([]models.ProposedTransfer, error) {
	newTransfers := make([]models.PendingTransfer, 0)
	for _, netSel := range g.GetNetworks() {
		balance, err := g.GetLiquidity(netSel)
		if err != nil {
			return nil, fmt.Errorf("get %d liquidity: %w", netSel, err)
		}

		// subtract inflight transfers from the balance
		for _, tr := range inflightTransfers {
			if tr.From == netSel && tr.Status != models.TransferStatusExecuted {
				balance = big.NewInt(0).Sub(balance, tr.Amount.ToInt())
			}
		}
		if balance.Cmp(big.NewInt(0)) <= 0 {
			continue
		}

		neighbors := p.eligibleNeighbors(g, netSel, append(inflightTransfers, newTransfers...))
		if len(neighbors) == 0 {
			continue
		}

		// Equally distribute the balance to each eligible neighbor.
		// If balance is not enough for everyone then start skipping neighbors.
		amountToSend := big.NewInt(0).Div(balance, big.NewInt(int64(len(neighbors))))
		for amountToSend.Cmp(big.NewInt(0)) <= 0 && len(neighbors) > 1 {
			neighbors = neighbors[:len(neighbors)-1]
			amountToSend = big.NewInt(0).Div(balance, big.NewInt(int64(len(neighbors))))
		}

		for _, neighborNetSel := range neighbors {
			newTransfer := models.NewTransfer(netSel, neighborNetSel, amountToSend, time.Now().UTC(), nil)
			newTransfers = append(newTransfers, models.NewPendingTransfer(newTransfer))
		}
	}

	results := make([]models.ProposedTransfer, len(newTransfers))
	for i, tr := range newTransfers {
		results[i] = models.ProposedTransfer{
			From:   tr.From,
			To:     tr.To,
			Amount: tr.Amount,
		}
	}
	return results, nil
}

// eligibleNeighbors returns the neighbors that:
// 1. Can transfer back (bidirectional graph connection).
// 2. There is no inflight transfer in either direction.
func (p *PingPong) eligibleNeighbors(g graph.Graph, netSel models.NetworkSelector, inflight []models.PendingTransfer) []models.NetworkSelector {
	allNeighbors, exists := g.GetNeighbors(netSel)
	if !exists {
		panic(fmt.Errorf("critical internal graph issue: neighbors of %d not found", netSel))
	}

	targetNeighbors := make([]models.NetworkSelector, 0, len(allNeighbors))
	for _, neighborNetSel := range allNeighbors {
		if !g.HasConnection(neighborNetSel, netSel) {
			continue
		}
		if p.isInflightBidirectionally(netSel, neighborNetSel, inflight) {
			continue
		}
		targetNeighbors = append(targetNeighbors, neighborNetSel)
	}

	return targetNeighbors
}

func (p *PingPong) isInflightBidirectionally(from, to models.NetworkSelector, inflight []models.PendingTransfer) bool {
	for _, tr := range inflight {
		if tr.Status == models.TransferStatusExecuted {
			continue
		}
		if (tr.From == from && tr.To == to) || (tr.From == to && tr.To == from) {
			return true
		}
	}
	return false
}
