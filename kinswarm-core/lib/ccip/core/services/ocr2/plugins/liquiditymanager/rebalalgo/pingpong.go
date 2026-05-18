package rebalalgo

import (
	"fmt"
	"math/big"
	"time"

	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/graph"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

// PingPong rebalancer keeps sending tokens between the chains without ever reaching balance. For testing purposes.
type PingPong struct{}

func NewPingPong() *PingPong {
	return &PingPong{}
}

func (p *PingPong) ComputeTransfersToBalance(g graph.Graph, unexecuted []UnexecutedTransfer) ([]models.ProposedTransfer, error) {
	newTransfers := make([]UnexecutedTransfer, 0)
	for _, netSel := range g.GetNetworks() {
		balance, err := g.GetLiquidity(netSel)
		if err != nil {
			return nil, fmt.Errorf("get %d liquidity: %w", netSel, err)
		}

		// subtract unexecuted transfers from the balance
		for _, tr := range unexecuted {
			if tr.FromNetwork() == netSel && tr.TransferStatus() != models.TransferStatusExecuted {
				balance = big.NewInt(0).Sub(balance, tr.TransferAmount())
			}
		}
		if balance.Cmp(big.NewInt(0)) <= 0 {
			continue
		}

		neighbors := p.eligibleNeighbors(g, netSel, append(unexecuted, newTransfers...))
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
			From:   tr.FromNetwork(),
			To:     tr.ToNetwork(),
			Amount: ubig.New(tr.TransferAmount()),
		}
	}
	return results, nil
}

// eligibleNeighbors returns the neighbors that:
// 1. Can transfer back (bidirectional graph connection).
// 2. There is no unexecuted transfer in either direction.
func (p *PingPong) eligibleNeighbors(g graph.Graph, netSel models.NetworkSelector, unexecuted []UnexecutedTransfer) []models.NetworkSelector {
	allNeighbors, exists := g.GetNeighbors(netSel, true)
	if !exists {
		panic(fmt.Errorf("critical internal graph issue: neighbors of %d not found", netSel))
	}

	targetNeighbors := make([]models.NetworkSelector, 0, len(allNeighbors))
	for _, neighborNetSel := range allNeighbors {
		if p.isUnexecutedBidirectionally(netSel, neighborNetSel, unexecuted) {
			continue
		}
		targetNeighbors = append(targetNeighbors, neighborNetSel)
	}

	return targetNeighbors
}

func (p *PingPong) isUnexecutedBidirectionally(from, to models.NetworkSelector, unexecuted []UnexecutedTransfer) bool {
	for _, tr := range unexecuted {
		if tr.TransferStatus() == models.TransferStatusExecuted {
			continue
		}
		if (tr.FromNetwork() == from && tr.ToNetwork() == to) ||
			(tr.FromNetwork() == to && tr.ToNetwork() == from) {
			return true
		}
	}
	return false
}
