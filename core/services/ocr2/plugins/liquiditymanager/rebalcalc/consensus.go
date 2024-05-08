package rebalcalc

import (
	"fmt"
	"math/big"
	"sort"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

// quorum returns the number of observations required.
func quorum(f int) int {
	return 2*f + 1
}

// bft returns the number of Byzantine Fault Tolerant nodes for a given f,
// which is the minimum number of nodes required to reach consensus on an observed value.
func bft(f int) int {
	return f + 1
}

// MedianLiquidityPerChain returns the median liquidity per chain from the provided observations.
func MedianLiquidityPerChain(observations []models.Observation, f int) ([]models.NetworkLiquidity, error) {
	if len(observations) < quorum(f) {
		return nil, fmt.Errorf("need at least 2f+1 observations (ocr3types.QuorumTwoFPlusOne) to reach consensus, got: %d", len(observations))
	}

	liqObsPerChain := make(map[models.NetworkSelector][]*big.Int)
	for _, ob := range observations {
		for _, chainLiq := range ob.LiquidityPerChain {
			liqObsPerChain[chainLiq.Network] = append(liqObsPerChain[chainLiq.Network], chainLiq.Liquidity.ToInt())
		}
	}

	for chainID, liqs := range liqObsPerChain {
		if len(liqs) < bft(f) {
			// If we don't have enough observations for a chain, we remove it to excluded from the result
			delete(liqObsPerChain, chainID)
		}
	}

	medians := make([]models.NetworkLiquidity, 0, len(liqObsPerChain))
	for chainID, liqs := range liqObsPerChain {
		medians = append(medians, models.NewNetworkLiquidity(chainID, BigIntSortedMiddle(liqs)))
	}

	// sort by network id for deterministic results
	sort.Slice(medians, func(i, j int) bool {
		return medians[i].Network < medians[j].Network
	})

	return medians, nil
}

// PendingTransfersConsensus returns the pending transfers that have been observed by at least f+1 observers.
func PendingTransfersConsensus(observations []models.Observation, f int) ([]models.PendingTransfer, error) {
	if len(observations) < quorum(f) {
		return nil, fmt.Errorf("need at least 2f+1 observations (ocr3types.QuorumTwoFPlusOne) to reach consensus, got: %d", len(observations))
	}

	eventFromHash := make(map[[32]byte]models.PendingTransfer)
	counts := make(map[[32]byte]int)
	for _, obs := range observations {
		for _, tr := range obs.PendingTransfers {
			h, err := tr.Hash()
			if err != nil {
				return nil, fmt.Errorf("hash %v: %w", tr, err)
			}
			counts[h]++
			eventFromHash[h] = tr
		}
	}

	quorumEvents := make([]models.PendingTransfer, 0, len(counts))
	for h, count := range counts {
		if count >= bft(f) {
			ev, exists := eventFromHash[h]
			if !exists {
				return nil, fmt.Errorf("internal issue, event from hash %v not found", h)
			}
			quorumEvents = append(quorumEvents, ev)
		}
	}

	// sort by ID for deterministic results
	sort.Slice(quorumEvents, func(i, j int) bool {
		return quorumEvents[i].ID < quorumEvents[j].ID
	})

	return quorumEvents, nil
}

func InflightTransfersConsensus(observations []models.Observation, f int) ([]models.Transfer, error) {
	if len(observations) < quorum(f) {
		return nil, fmt.Errorf("need at least 2f+1 observations (ocr3types.QuorumTwoFPlusOne) to reach consensus, got: %d", len(observations))
	}

	key := func(tr models.Transfer) string {
		return fmt.Sprintf("%d-%d-%s", tr.From, tr.To, tr.Amount.String())
	}
	transferFromKey := make(map[string]models.Transfer)
	counts := make(map[string]int)
	for _, obs := range observations {
		for _, tr := range obs.InflightTransfers {
			k := key(tr)
			counts[k]++
			transferFromKey[k] = tr
		}
	}

	quorumEvents := make([]models.Transfer, 0, len(counts))
	for h, count := range counts {
		if count >= bft(f) {
			ev, exists := transferFromKey[h]
			if !exists {
				return nil, fmt.Errorf("internal issue, event from hash %v not found", h)
			}
			quorumEvents = append(quorumEvents, ev)
		}
	}

	// sort by network id for deterministic results
	sort.Slice(quorumEvents, func(i, j int) bool {
		return quorumEvents[i].From < quorumEvents[j].From
	})

	return quorumEvents, nil
}

// ConfigDigestsConsensus returns the config digests that have been observed by at least f+1 observers.
func ConfigDigestsConsensus(observations []models.Observation, f int) ([]models.ConfigDigestWithMeta, error) {
	if len(observations) < quorum(f) {
		return nil, fmt.Errorf("need at least 2f+1 observations (ocr3types.QuorumTwoFPlusOne) to reach consensus, got: %d", len(observations))
	}

	key := func(meta models.ConfigDigestWithMeta) string {
		return fmt.Sprintf("%d-%s", meta.NetworkSel, meta.Digest.Hex())
	}
	counts := make(map[string]int)
	cds := make(map[string]models.ConfigDigestWithMeta)
	for _, obs := range observations {
		for _, cd := range obs.ConfigDigests {
			k := key(cd)
			counts[k]++
			if counts[k] == 1 {
				cds[k] = cd
			}
		}
	}

	quorumCds := make([]models.ConfigDigestWithMeta, 0, len(counts))
	for k, count := range counts {
		if count >= bft(f) {
			cd, exists := cds[k]
			if !exists {
				return nil, fmt.Errorf("internal issue, config digest by key %s not found", k)
			}
			quorumCds = append(quorumCds, cd)
		}
	}

	// sort by network id for deterministic results
	sort.Slice(quorumCds, func(i, j int) bool {
		return quorumCds[i].NetworkSel < quorumCds[j].NetworkSel
	})

	return quorumCds, nil
}

// GraphEdgesConsensus returns the edges that have been observed by at least f+1 observers.
func GraphEdgesConsensus(observations []models.Observation, f int) ([]models.Edge, error) {
	if len(observations) < quorum(f) {
		return nil, fmt.Errorf("need at least 2f+1 observations (ocr3types.QuorumTwoFPlusOne) to reach consensus, got: %d", len(observations))
	}

	counts := make(map[models.Edge]int)
	for _, obs := range observations {
		for _, edge := range obs.Edges {
			counts[edge]++
		}
	}

	var quorumEdges []models.Edge
	for edge, count := range counts {
		if count >= bft(f) {
			quorumEdges = append(quorumEdges, edge)
		}
	}

	// sort for deterministic results
	sort.Slice(quorumEdges, func(i, j int) bool {
		if quorumEdges[i].Source == quorumEdges[j].Source {
			return quorumEdges[i].Dest < quorumEdges[j].Dest
		}
		return quorumEdges[i].Source < quorumEdges[j].Source
	})

	return quorumEdges, nil
}
