package rebalalgo

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
	lggr   logger.Logger
	config models.PluginConfig
}

func NewTargetMinBalancer(lggr logger.Logger, config models.PluginConfig) *TargetMinBalancer {
	return &TargetMinBalancer{
		lggr:   lggr.With("service", "TargetMinBalancer"),
		config: config,
	}
}

type determineTransfersFunc func(graphFuture graph.Graph, targetNetwork models.NetworkSelector, networkFunds map[models.NetworkSelector]*Funds) ([]models.ProposedTransfer, error)

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
	graphFuture, err := r.getExpectedGraph(graphNow, nonExecutedTransfers)
	if err != nil {
		return nil, fmt.Errorf("get expected graph: %w", err)
	}

	r.lggr.Debugf("finding networks that require funding")
	networksRequiringFunding, networkFunds, err := r.findNetworksRequiringFunding(graphNow, graphFuture)
	if err != nil {
		return nil, fmt.Errorf("find networks that require funding: %w", err)
	}
	r.lggr.Debugf("networks requiring funding: %v", networksRequiringFunding)
	r.lggr.Debugf("network funds: %+v", networkFunds)

	proposedTransfers := make([]models.ProposedTransfer, 0)
	for _, net := range networksRequiringFunding {
		r.lggr.Debugf("finding transfers for network %v", net)
		potentialTransfers, err1 := transfersFunc(graphFuture, net, networkFunds)
		if err1 != nil {
			return nil, fmt.Errorf("finding transfers for network %v: %w", net, err1)
		}

		dataFuture, err2 := graphFuture.GetData(net)
		if err2 != nil {
			return nil, fmt.Errorf("get future data of net %v: %w", net, err2)
		}
		liqDiffFuture := new(big.Int).Sub(dataFuture.TargetLiquidity, dataFuture.Liquidity)
		netProposedTransfers, err3 := r.applyProposedTransfers(graphFuture, potentialTransfers, liqDiffFuture)
		if err3 != nil {
			return nil, fmt.Errorf("applying transfers: %w", err3)
		}
		proposedTransfers = append(proposedTransfers, netProposedTransfers...)
	}

	return proposedTransfers, nil
}

// getExpectedGraph returns a copy of the graph instance with all the non-executed transfers applied and targets set for each network.
func (r *TargetMinBalancer) getExpectedGraph(g graph.Graph, nonExecutedTransfers []UnexecutedTransfer) (graph.Graph, error) {
	expG := g.Clone()

	for _, tr := range nonExecutedTransfers {
		liqTo, err := expG.GetLiquidity(tr.ToNetwork())
		if err != nil {
			return nil, err
		}
		expG.SetLiquidity(tr.ToNetwork(), big.NewInt(0).Add(liqTo, tr.TransferAmount()))

		// we only subtract from the sender if the transfer is still in progress, otherwise the source value would have already been updated
		switch tr.TransferStatus() {
		case models.TransferStatusProposed, models.TransferStatusInflight:
			liqFrom, err := expG.GetLiquidity(tr.FromNetwork())
			if err != nil {
				return nil, err
			}
			expG.SetLiquidity(tr.FromNetwork(), big.NewInt(0).Sub(liqFrom, tr.TransferAmount()))
		}
	}

	for _, selector := range expG.GetNetworks() {
		target := r.config.RebalancerConfig.DefaultTarget
		override, ok := r.config.RebalancerConfig.NetworkTargetOverrides[selector]
		if ok {
			target = override
		}
		expG.SetTargetLiquidity(selector, target)
	}

	return expG, nil
}

func (r *TargetMinBalancer) findNetworksRequiringFunding(graphNow, graphFuture graph.Graph) ([]models.NetworkSelector, map[models.NetworkSelector]*Funds, error) {
	mapNetworkFunds := make(map[models.NetworkSelector]*Funds)
	liqDiffsFuture := make(map[models.NetworkSelector]*big.Int)

	res := make([]models.NetworkSelector, 0)
	for _, net := range graphNow.GetNetworks() {
		//use min here for transferable. because we don't know when the transfers will complete and want to avoid issues
		transferableAmount, ataErr := availableTransferableAmount(graphNow, graphFuture, net)
		if ataErr != nil {
			return nil, nil, fmt.Errorf("getting available transferrable amount for net %d: %v", net, ataErr)
		}
		dataFuture, err := graphFuture.GetData(net)
		if err != nil {
			return nil, nil, fmt.Errorf("get future data of net %v: %w", net, err)
		}
		liqDiffFuture := new(big.Int).Sub(dataFuture.TargetLiquidity, dataFuture.Liquidity)
		mapNetworkFunds[net] = &Funds{
			AvailableAmount: transferableAmount,
		}
		if liqDiffFuture.Cmp(big.NewInt(0)) <= 0 {
			continue
		}

		// only if we are below 5% else it's too little to worry about now
		fivePercent := big.NewInt(5)
		hundred := big.NewInt(100)
		value := new(big.Int).Mul(dataFuture.TargetLiquidity, fivePercent)
		value.Div(value, hundred)
		if liqDiffFuture.Cmp(value) < 0 {
			continue
		}
		liqDiffsFuture[net] = liqDiffFuture
		res = append(res, net)
	}

	sort.Slice(res, func(i, j int) bool { return liqDiffsFuture[res[i]].Cmp(liqDiffsFuture[res[j]]) > 0 })
	return res, mapNetworkFunds, nil
}

func (r *TargetMinBalancer) oneHopTransfers(graphFuture graph.Graph, targetNetwork models.NetworkSelector, networkFunds map[models.NetworkSelector]*Funds) ([]models.ProposedTransfer, error) {
	zero := big.NewInt(0)
	potentialTransfers := make([]models.ProposedTransfer, 0)

	neighbors, exist := graphFuture.GetNeighbors(targetNetwork, false)
	if !exist {
		r.lggr.Debugf("no neighbors found for %v", targetNetwork)
		return nil, nil
	}
	targetFuture, err := graphFuture.GetData(targetNetwork)
	if err != nil {
		return nil, fmt.Errorf("get data Future of net %v: %w", targetNetwork, err)
	}

	for _, source := range neighbors {
		transferAmount := new(big.Int).Sub(targetFuture.TargetLiquidity, targetFuture.Liquidity)
		r.lggr.Debugf("checking transfer from %v to %v for amount %v", source, targetNetwork, transferAmount)

		//source network available transferable amount
		srcData, dErr := graphFuture.GetData(source)
		if dErr != nil {
			return nil, fmt.Errorf("error during GetData for %v in graphFuture: %v", source, dErr)
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
func (r *TargetMinBalancer) twoHopTransfers(graphFuture graph.Graph, targetNetwork models.NetworkSelector, networkFunds map[models.NetworkSelector]*Funds) ([]models.ProposedTransfer, error) {
	zero := big.NewInt(0)
	iterator := func(nodes ...graph.Data) bool { return true }
	potentialTransfers := make([]models.ProposedTransfer, 0)

	for _, source := range graphFuture.GetNetworks() {
		if source == targetNetwork {
			continue
		}
		path := graphFuture.FindPath(source, targetNetwork, 2, iterator)
		if len(path) != 2 {
			continue
		}
		middle := path[0]

		targetData, err := graphFuture.GetData(targetNetwork)
		if err != nil {
			return nil, fmt.Errorf("error during GetData for %v in graphFuture: %v", targetNetwork, err)
		}
		transferAmount := new(big.Int).Sub(targetData.TargetLiquidity, targetData.Liquidity)
		r.lggr.Debugf("checking transfer from %v -> %v -> %v for amount %v", source, middle, targetNetwork, transferAmount)

		//source network available transferable amount
		srcData, dErr := graphFuture.GetData(source)
		if dErr != nil {
			return nil, fmt.Errorf("error during GetData for %v in graphFuture: %v", source, dErr)
		}
		srcAvailableAmount := new(big.Int).Sub(srcData.Liquidity, srcData.MinimumLiquidity)
		srcAmountToTarget := new(big.Int).Sub(srcData.Liquidity, srcData.TargetLiquidity)

		//middle network available transferable amount
		middleData, dErr := graphFuture.GetData(middle)
		if dErr != nil {
			return nil, fmt.Errorf("error during GetData for %v in graphFuture: %v", middle, dErr)
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
func (r *TargetMinBalancer) applyProposedTransfers(graphFuture graph.Graph, potentialTransfers []models.ProposedTransfer, requiredAmount *big.Int) ([]models.ProposedTransfer, error) {
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

		senderData, err := graphFuture.GetData(d.From)
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

		liqBefore, err := graphFuture.GetLiquidity(d.To)
		if err != nil {
			return nil, fmt.Errorf("get liquidity of transfer receiver %v: %w", d.To, err)
		}
		graphFuture.SetLiquidity(d.To, big.NewInt(0).Add(liqBefore, d.Amount.ToInt()))

		liqBefore, err = graphFuture.GetLiquidity(d.From)
		if err != nil {
			return nil, fmt.Errorf("get liquidity of sender %v: %w", d.From, err)
		}
		graphFuture.SetLiquidity(d.From, big.NewInt(0).Sub(liqBefore, d.Amount.ToInt()))

		if fundsRaised.Cmp(requiredAmount) >= 0 {
			r.lggr.Debugf("all funds raised skipping further transfers")
			skip = true
		}
	}

	return proposedTransfers, nil
}
