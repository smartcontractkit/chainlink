package liquidityrebalancer

import (
	"fmt"
	"math/big"
	"sort"

	big2 "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/graph"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

// TargetMinBalancer tries to reach balance using a target and minimum liquidity that is configured on each network.
type TargetMinBalancer struct {
	lggr logger.Logger
}

func NewTargetMinBalancer(lggr logger.Logger) *TargetMinBalancer {
	return &TargetMinBalancer{
		lggr: lggr.With("service", "TargetMinBalancer"),
	}
}

type determineTransfersFunc func(graphLater graph.Graph, targetNetwork models.NetworkSelector, networkFunds map[models.NetworkSelector]*Funds) ([]models.ProposedTransfer, error)

func (r *TargetMinBalancer) ComputeTransfersToBalance(graphNow graph.Graph, nonExecutedTransfers []UnexecutedTransfer) ([]models.ProposedTransfer, error) {
	nonExecutedTransfers = filterUnexecutedTransfers(nonExecutedTransfers)

	var proposedTransfers []models.ProposedTransfer
	// 4 rounds of rebalancing alternate between 1 hop and 2 hop transfers.
	// this allows us to have multistep transaction initiated at the same time.
	for i := 0; i < 5; i++ {
		r.lggr.Debugf("Round %d: nonExecutedTransfers: %v", i, nonExecutedTransfers)
		var currentProposed []models.ProposedTransfer
		transfersFunc := r.oneHopTransfers
		if i%2 != 0 {
			transfersFunc = r.twoHopTransfers
		}
		currentProposed, err := r.rebalancingRound(graphNow, nonExecutedTransfers, transfersFunc)
		if err != nil {
			return nil, err
		}
		r.lggr.Debugf("Round %d: proposed transfers: %v", i, currentProposed)
		for _, t := range currentProposed {
			// put proposed in nonExecutedTransfers to carryover to next round
			nonExecutedTransfers = append(nonExecutedTransfers, t)
		}
		proposedTransfers = append(proposedTransfers, currentProposed...)
	}

	r.lggr.Debugf("merging proposed transfers")
	proposedTransfers = mergeProposedTransfers(proposedTransfers)
	r.lggr.Debugf("sorting proposed transfers for determinism")
	sort.Sort(models.ProposedTransfers(proposedTransfers))

	return proposedTransfers, nil
}

func (r *TargetMinBalancer) rebalancingRound(graphNow graph.Graph, nonExecutedTransfers []UnexecutedTransfer, transfersFunc determineTransfersFunc) ([]models.ProposedTransfer, error) {
	var err error
	graphLater, err := getExpectedGraph(graphNow, nonExecutedTransfers)
	if err != nil {
		return nil, fmt.Errorf("get expected graph: %w", err)
	}

	r.lggr.Debugf("finding networks that require funding")
	networksRequiringFunding, networkFunds, err := r.findNetworksRequiringFunding(graphNow, graphLater)
	if err != nil {
		return nil, fmt.Errorf("find networks that require funding: %w", err)
	}
	r.lggr.Debugf("networks requiring funding: %v", networksRequiringFunding)
	r.lggr.Debugf("network funds: %+v", networkFunds)

	proposedTransfers := make([]models.ProposedTransfer, 0)
	for _, net := range networksRequiringFunding {
		r.lggr.Debugf("finding transfers for network %v", net)
		potentialTransfers, err1 := transfersFunc(graphLater, net, networkFunds)
		if err1 != nil {
			return nil, fmt.Errorf("finding transfers for network %v: %w", net, err1)
		}

		dataLater, err2 := graphLater.GetData(net)
		if err2 != nil {
			return nil, fmt.Errorf("get data later of net %v: %w", net, err2)
		}
		liqDiffLater := new(big.Int).Sub(dataLater.TargetLiquidity, dataLater.Liquidity)
		netProposedTransfers, err3 := r.applyProposedTransfers(graphLater, potentialTransfers, liqDiffLater)
		if err3 != nil {
			return nil, fmt.Errorf("applying transfers: %w", err3)
		}
		proposedTransfers = append(proposedTransfers, netProposedTransfers...)
	}

	return proposedTransfers, nil
}

func (r *TargetMinBalancer) findNetworksRequiringFunding(graphNow, graphLater graph.Graph) ([]models.NetworkSelector, map[models.NetworkSelector]*Funds, error) {
	mapNetworkFunds := make(map[models.NetworkSelector]*Funds)
	liqDiffsLater := make(map[models.NetworkSelector]*big.Int)

	//TODO: LM-23 Create minTokenTransfer config to filter-out small rebalance txs
	// check that the transfer is not tiny, we should only transfer if it is significant. What is too tiny?
	// we could prevent this by only making a network requiring funding if its below X% of the target

	res := make([]models.NetworkSelector, 0)
	for _, net := range graphNow.GetNetworks() {
		//use min here for transferable. because we don't know when the transfers will complete and want to avoid issues
		transferableAmount, ataErr := availableTransferableAmount(graphNow, graphLater, net)
		if ataErr != nil {
			return nil, nil, fmt.Errorf("getting available transferrable amount for net %d: %v", net, ataErr)
		}
		dataLater, err := graphLater.GetData(net)
		if err != nil {
			return nil, nil, fmt.Errorf("get data later of net %v: %w", net, err)
		}
		liqDiffLater := new(big.Int).Sub(dataLater.TargetLiquidity, dataLater.Liquidity)
		mapNetworkFunds[net] = &Funds{
			AvailableAmount: transferableAmount,
		}
		if liqDiffLater.Cmp(big.NewInt(0)) <= 0 {
			continue
		}
		liqDiffsLater[net] = liqDiffLater
		res = append(res, net)
	}

	sort.Slice(res, func(i, j int) bool { return liqDiffsLater[res[i]].Cmp(liqDiffsLater[res[j]]) > 0 })
	return res, mapNetworkFunds, nil
}

func (r *TargetMinBalancer) oneHopTransfers(graphLater graph.Graph, targetNetwork models.NetworkSelector, networkFunds map[models.NetworkSelector]*Funds) ([]models.ProposedTransfer, error) {
	zero := big.NewInt(0)
	potentialTransfers := make([]models.ProposedTransfer, 0)

	neighbors, exist := graphLater.GetNeighbors(targetNetwork, false)
	if !exist {
		r.lggr.Debugf("no neighbors found for %v", targetNetwork)
		return nil, nil
	}
	targetLater, err := graphLater.GetData(targetNetwork)
	if err != nil {
		return nil, fmt.Errorf("get data later of net %v: %w", targetNetwork, err)
	}

	for _, source := range neighbors {
		transferAmount := new(big.Int).Sub(targetLater.TargetLiquidity, targetLater.Liquidity)
		r.lggr.Debugf("checking transfer from %v to %v for amount %v", source, targetNetwork, transferAmount)

		//source network available transferable amount
		srcData, dErr := graphLater.GetData(source)
		if dErr != nil {
			return nil, fmt.Errorf("error during GetData for %v in graphLater: %v", source, dErr)
		}
		srcAvailableAmount := new(big.Int).Sub(srcData.Liquidity, srcData.MinimumLiquidity)
		srcAmountToTarget := new(big.Int).Sub(srcData.Liquidity, srcData.TargetLiquidity)

		if srcAmountToTarget.Cmp(zero) <= 0 || srcAvailableAmount.Cmp(zero) <= 0 {
			r.lggr.Debugf("source network %v does not have a surplus to transfer so skipping transfer, source amount to target %v", source, srcAmountToTarget)
			continue
		}

		if transferAmount.Cmp(srcAmountToTarget) > 0 {
			// if transferAmount > srcAmountToTarget take less
			r.lggr.Debugf("source network %v does not have the desired amount, desired amount %v taking available %v", source, transferAmount, srcAmountToTarget)
			transferAmount = srcAmountToTarget
		}

		newAmount := new(big.Int).Sub(networkFunds[source].AvailableAmount, transferAmount)
		if newAmount.Cmp(zero) < 0 {
			r.lggr.Debugf("source network %v doesn't have enough available liquidity so skipping transfer, desired amount %v but only have %v available", source, transferAmount, networkFunds[source].AvailableAmount)
			continue
		}
		networkFunds[source].AvailableAmount = newAmount

		potentialTransfers = append(potentialTransfers, newTransfer(source, targetNetwork, transferAmount))
	}

	return potentialTransfers, nil
}

// twoHopTransfers finds networks that can increase liquidity of the target network with an intermediate network.
func (r *TargetMinBalancer) twoHopTransfers(graphLater graph.Graph, targetNetwork models.NetworkSelector, networkFunds map[models.NetworkSelector]*Funds) ([]models.ProposedTransfer, error) {
	zero := big.NewInt(0)
	iterator := func(nodes ...graph.Data) bool { return true }
	potentialTransfers := make([]models.ProposedTransfer, 0)

	for _, source := range graphLater.GetNetworks() {
		if source == targetNetwork {
			continue
		}
		path := graphLater.FindPath(source, targetNetwork, 2, iterator)
		if len(path) != 2 {
			continue
		}
		middle := path[0]

		targetData, err := graphLater.GetData(targetNetwork)
		if err != nil {
			return nil, fmt.Errorf("error during GetData for %v in graphLater: %v", targetNetwork, err)
		}
		transferAmount := new(big.Int).Sub(targetData.TargetLiquidity, targetData.Liquidity)
		r.lggr.Debugf("checking transfer from %v -> %v -> %v for amount %v", source, middle, targetNetwork, transferAmount)

		//source network available transferable amount
		srcData, dErr := graphLater.GetData(source)
		if dErr != nil {
			return nil, fmt.Errorf("error during GetData for %v in graphLater: %v", source, dErr)
		}
		srcAvailableAmount := new(big.Int).Sub(srcData.Liquidity, srcData.MinimumLiquidity)
		srcAmountToTarget := new(big.Int).Sub(srcData.Liquidity, srcData.TargetLiquidity)

		//middle network available transferable amount
		middleData, dErr := graphLater.GetData(middle)
		if dErr != nil {
			return nil, fmt.Errorf("error during GetData for %v in graphLater: %v", middle, dErr)
		}
		middleAvailableAmount := new(big.Int).Sub(middleData.Liquidity, middleData.MinimumLiquidity)
		middleAmountToTarget := new(big.Int).Sub(middleData.Liquidity, middleData.TargetLiquidity)

		if transferAmount.Cmp(srcAmountToTarget) > 0 {
			// if transferAmount > srcAmountToTarget take less
			transferAmount = srcAmountToTarget
		}
		if transferAmount.Cmp(zero) <= 0 {
			continue
		}

		if srcAmountToTarget.Cmp(transferAmount) < 0 || srcAvailableAmount.Cmp(transferAmount) < 0 {
			continue
		}
		if middleAvailableAmount.Cmp(transferAmount) < 0 {
			// middle hop doesn't have enough available liquidity
			r.lggr.Debugf("middle network %v liquidity too low, skipping transfer: middleAmountToTarget %v, middleAvailableAmount %v", middle, middleAmountToTarget, middleAvailableAmount)
			continue
		}

		newAmount := new(big.Int).Sub(networkFunds[source].AvailableAmount, transferAmount)
		if newAmount.Cmp(zero) < 0 {
			r.lggr.Debugf("source network %v doesn't have enough available liquidity so skipping transfer, desired amount %v but only have %v available", source, transferAmount, networkFunds[source].AvailableAmount)
			continue
		}
		networkFunds[source].AvailableAmount = newAmount

		potentialTransfers = append(potentialTransfers, newTransfer(source, middle, transferAmount))
	}

	return potentialTransfers, nil
}

// applyProposedTransfers applies the proposed transfers to the graph.
// increments the raised funds and gives a refund to the sender if more funds have been raised than the required amount.
// It updates the liquidity of the sender and receiver networks in the graph. It stops further transfers if all funds have been raised.
func (r *TargetMinBalancer) applyProposedTransfers(graphLater graph.Graph, potentialTransfers []models.ProposedTransfer, requiredAmount *big.Int) ([]models.ProposedTransfer, error) {
	// sort by amount,sender,receiver
	sort.Slice(potentialTransfers, func(i, j int) bool {
		if potentialTransfers[i].Amount.Cmp(potentialTransfers[j].Amount) == 0 {
			if potentialTransfers[i].From == potentialTransfers[j].From {
				return potentialTransfers[i].To < potentialTransfers[j].To
			}
			return potentialTransfers[i].From < potentialTransfers[j].From
		}
		return potentialTransfers[i].Amount.Cmp(potentialTransfers[j].Amount) > 0
	})

	fundsRaised := big.NewInt(0)
	proposedTransfers := make([]models.ProposedTransfer, 0, len(potentialTransfers))
	skip := false
	for _, d := range potentialTransfers {
		if skip {
			r.lggr.Debugf("skipping transfer: %s", d)
			continue
		}

		senderData, err := graphLater.GetData(d.From)
		if err != nil {
			return nil, fmt.Errorf("get liquidity of sender %v: %w", d.From, err)
		}
		availableAmount := big.NewInt(0).Sub(senderData.Liquidity, senderData.MinimumLiquidity)
		if availableAmount.Cmp(big.NewInt(0)) <= 0 {
			r.lggr.Debugf("no more tokens to transfer, skipping transfer: %s", d)
			continue
		}

		if availableAmount.Cmp(d.Amount.ToInt()) < 0 {
			d.Amount = big2.New(availableAmount)
			r.lggr.Debugf("reducing transfer amount since sender balance has dropped: %s", d)
		}

		// increment the raised funds
		fundsRaised = big.NewInt(0).Add(fundsRaised, d.Amount.ToInt())

		// in case we raised more than target amount give refund to the sender
		if refund := big.NewInt(0).Sub(fundsRaised, requiredAmount); refund.Cmp(big.NewInt(0)) > 0 {
			d.Amount = big2.New(big.NewInt(0).Sub(d.Amount.ToInt(), refund))
			fundsRaised = big.NewInt(0).Sub(fundsRaised, refund)
		}
		r.lggr.Debugf("applying transfer: %v", d)
		proposedTransfers = append(proposedTransfers, models.ProposedTransfer{From: d.From, To: d.To, Amount: d.Amount})

		liqBefore, err := graphLater.GetLiquidity(d.To)
		if err != nil {
			return nil, fmt.Errorf("get liquidity of transfer receiver %v: %w", d.To, err)
		}
		graphLater.SetLiquidity(d.To, big.NewInt(0).Add(liqBefore, d.Amount.ToInt()))

		liqBefore, err = graphLater.GetLiquidity(d.From)
		if err != nil {
			return nil, fmt.Errorf("get liquidity of sender %v: %w", d.From, err)
		}
		graphLater.SetLiquidity(d.From, big.NewInt(0).Sub(liqBefore, d.Amount.ToInt()))

		if fundsRaised.Cmp(requiredAmount) >= 0 {
			r.lggr.Debugf("all funds raised skipping further transfers")
			skip = true
		}
	}

	return proposedTransfers, nil
}
