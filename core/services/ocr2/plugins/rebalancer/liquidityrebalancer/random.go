package liquidityrebalancer

import (
	"fmt"
	"math/big"
	mathrand "math/rand"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/liquiditygraph"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
)

var _ Rebalancer = &randomRebalancer{}

type randomRebalancer struct {
	maxNumTransfers      int
	checkSourceDestEqual bool
	lggr                 logger.Logger
}

func NewRandomRebalancer(maxNumTransfers int, checkSourceDestEqual bool, lggr logger.Logger) Rebalancer {
	return &randomRebalancer{
		maxNumTransfers:      maxNumTransfers,
		checkSourceDestEqual: checkSourceDestEqual,
		lggr:                 lggr.Named("RandomRebalancer"),
	}
}

// ComputeTransfersToBalance implements Rebalancer.
func (r *randomRebalancer) ComputeTransfersToBalance(
	g liquiditygraph.LiquidityGraph,
	inflightTransfers []models.PendingTransfer,
) ([]models.Transfer, error) {
	liquidities := make([]models.NetworkLiquidity, 0)
	for _, net := range g.GetNetworks() {
		liq, err := g.GetLiquidity(net)
		if err != nil {
			return nil, fmt.Errorf("internal graph err: %w", err)
		}
		liquidities = append(liquidities, models.NewNetworkLiquidity(net, liq))
	}

	if len(liquidities) == 0 {
		return nil, nil
	}

	// seed the randomness source so that all rebalancers produce the same output
	// for the same input
	r.lggr.Infow("RandomRebalancer: using median liquidity as seed", "medianLiquidity1", liquidities[0].Liquidity.String())
	source := mathrand.NewSource(liquidities[0].Liquidity.Int64()) //nolint:gosec
	rng := mathrand.New(source)                                    //nolint:gosec
	numTransfers := rng.Int63n(int64(r.maxNumTransfers))
	r.lggr.Infow("RandomRebalancer: generated random number of transfers", "numTransfers", numTransfers)
	var transfers []models.Transfer
	for i := 0; i < int(numTransfers); i++ {
		randSourceChain := pickRandom(rng, g.GetNetworks())
		neighbors, exist := g.GetNeighbors(randSourceChain)
		if !exist {
			return nil, fmt.Errorf("chain selector %d not found in graph %+v", randSourceChain, g)
		}
		randDestChain := pickRandom(rng, neighbors)
		r.lggr.Infow("RandomRebalancer: generated random transfer source and dest", "sourceChain", randSourceChain, "destChain", randDestChain)
		if r.checkSourceDestEqual && randSourceChain == randDestChain {
			continue
		}
		// use median liquidity to generate random amount
		var liqSource *big.Int
		for _, medianLiq := range liquidities {
			if medianLiq.Network == randSourceChain {
				liqSource = medianLiq.Liquidity
				break
			}
		}
		if liqSource == nil {
			return nil, fmt.Errorf("failed to find median liquidity for source chain %v", randSourceChain)
		}
		amount := rng.Int63n(liqSource.Int64())
		r.lggr.Infow("RandomRebalancer: generated random transfer amount", "amount", amount)

		transfers = append(transfers, models.Transfer{
			From:   randSourceChain,
			To:     randDestChain,
			Amount: big.NewInt(amount),
		})
	}
	r.lggr.Info("RandomRebalancer: generated random transfers", "transfers", transfers)
	return transfers, nil
}

func pickRandom(rng *mathrand.Rand, networks []models.NetworkSelector) models.NetworkSelector {
	randIndex := rng.Intn(len(networks))
	return networks[randIndex]
}
