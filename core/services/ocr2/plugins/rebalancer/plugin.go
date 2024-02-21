package rebalancer

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"go.uber.org/multierr"

	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/bridge"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/discoverer"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/graph"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/liquiditymanager"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/liquidityrebalancer"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
)

type Plugin struct {
	f                       int
	rootNetwork             models.NetworkSelector
	rootAddress             models.Address
	closePluginTimeout      time.Duration
	liquidityManagerFactory liquiditymanager.Factory
	discovererFactory       discoverer.Factory
	bridgeFactory           bridge.Factory
	mu                      sync.RWMutex
	rebalancerGraph         graph.Graph
	liquidityRebalancer     liquidityrebalancer.Rebalancer
	pendingTransfers        *PendingTransfersCache
	lggr                    logger.Logger
}

func NewPlugin(
	f int,
	closePluginTimeout time.Duration,
	rootNetwork models.NetworkSelector,
	rootAddress models.Address,
	liquidityManagerFactory liquiditymanager.Factory,
	discovererFactory discoverer.Factory,
	bridgeFactory bridge.Factory,
	liquidityRebalancer liquidityrebalancer.Rebalancer,
	lggr logger.Logger,
) *Plugin {
	return &Plugin{
		f:                       f,
		rootNetwork:             rootNetwork,
		rootAddress:             rootAddress,
		closePluginTimeout:      closePluginTimeout,
		liquidityManagerFactory: liquidityManagerFactory,
		discovererFactory:       discovererFactory,
		bridgeFactory:           bridgeFactory,
		rebalancerGraph:         graph.NewGraph(),
		liquidityRebalancer:     liquidityRebalancer,
		pendingTransfers:        NewPendingTransfersCache(),
		lggr:                    lggr,
		mu:                      sync.RWMutex{},
	}
}

func (p *Plugin) Query(_ context.Context, outcomeCtx ocr3types.OutcomeContext) (ocrtypes.Query, error) {
	p.lggr.Infow("in query", "seqNr", outcomeCtx.SeqNr)
	return ocrtypes.Query{}, nil
}

func (p *Plugin) Observation(ctx context.Context, outcomeCtx ocr3types.OutcomeContext, _ ocrtypes.Query) (ocrtypes.Observation, error) {
	lggr := p.lggr.With("seqNr", outcomeCtx.SeqNr)
	lggr.Infow("in observation", "seqNr", outcomeCtx.SeqNr)

	if err := p.syncGraphEdges(ctx); err != nil {
		return ocrtypes.Observation{}, fmt.Errorf("sync graph edges: %w", err)
	}

	networkLiquidities, err := p.syncGraphBalances(ctx)
	if err != nil {
		return ocrtypes.Observation{}, fmt.Errorf("sync graph balances: %w", err)
	}

	pendingTransfers, err := p.loadPendingTransfers(ctx)
	if err != nil {
		return ocrtypes.Observation{}, fmt.Errorf("load pending transfers: %w", err)
	}

	edges, err := p.rebalancerGraph.GetEdges()
	if err != nil {
		return ocrtypes.Observation{}, fmt.Errorf("get edges: %w", err)
	}

	resolvedTransfers, err := p.resolveProposedTransfers(ctx, outcomeCtx)
	if err != nil {
		return ocrtypes.Observation{}, fmt.Errorf("resolve proposed transfers: %w", err)
	}

	configDigests := make([]models.ConfigDigestWithMeta, 0)
	for _, net := range p.rebalancerGraph.GetNetworks() {
		data, err := p.rebalancerGraph.GetData(net)
		if err != nil {
			return nil, fmt.Errorf("get rb %d data: %w", net, err)
		}
		configDigests = append(configDigests, models.ConfigDigestWithMeta{
			Digest:         data.ConfigDigest,
			NetworkSel:     data.NetworkSelector,
			RebalancerAddr: data.RebalancerAddress,
		})
	}

	lggr.Infow("finished observing",
		"networkLiquidities", networkLiquidities,
		"pendingTransfers", pendingTransfers,
		"edges", edges,
		"resolvedTransfers", resolvedTransfers,
	)

	return models.NewObservation(networkLiquidities, resolvedTransfers, pendingTransfers, edges, configDigests).Encode(), nil
}

func (p *Plugin) ValidateObservation(outctx ocr3types.OutcomeContext, query ocrtypes.Query, ao ocrtypes.AttributedObservation) error {
	p.lggr.Infow("in validate observation", "seqNr", outctx.SeqNr)

	_, err := models.DecodeObservation(ao.Observation)
	if err != nil {
		return fmt.Errorf("invalid observation: %w", err)
	}

	// todo: consider adding more validations

	return nil
}

func (p *Plugin) ObservationQuorum(outctx ocr3types.OutcomeContext, query ocrtypes.Query) (ocr3types.Quorum, error) {
	return ocr3types.QuorumTwoFPlusOne, nil
}

func (p *Plugin) Outcome(outctx ocr3types.OutcomeContext, query ocrtypes.Query, aos []ocrtypes.AttributedObservation) (ocr3types.Outcome, error) {
	lggr := p.lggr.With("seqNr", outctx.SeqNr, "numObservations", len(aos))
	lggr.Infow("in outcome")

	// Gather all the observations.
	observations := make([]models.Observation, 0, len(aos))
	for _, encodedObs := range aos {
		obs, err := models.DecodeObservation(encodedObs.Observation)
		if err != nil {
			return ocr3types.Outcome{}, fmt.Errorf("decode observation: %w", err)
		}
		lggr.Debugw("decoded observation", "observation", obs, "oracleID", encodedObs.Observer)
		observations = append(observations, obs)
	}

	// Come to a consensus based on the observations of all the different nodes.
	medianLiquidityPerChain := p.computeMedianLiquidityPerChain(observations)
	graphEdges := p.computeGraphEdgesConsensus(observations)
	pendingTransfers, err := p.computePendingTransfersConsensus(observations)
	if err != nil {
		return ocr3types.Outcome{}, fmt.Errorf("compute pending transfers consensus: %w", err)
	}
	configDigests, err := p.computeConfigDigestsConsensus(observations)
	if err != nil {
		return ocr3types.Outcome{}, fmt.Errorf("compute config digests consensus: %w", err)
	}

	// Compute a new graph with the median liquidities and the edges of the quorum of nodes.
	g, err := p.computeMedianGraph(graphEdges, medianLiquidityPerChain)
	if err != nil {
		return nil, fmt.Errorf("compute median graph: %w", err)
	}

	resolvedTransfersQuorum, err := p.computeResolvedTransfersQuorum(observations)
	if err != nil {
		return nil, fmt.Errorf("compute resolved transfers quorum: %w", err)
	}

	lggr.Infow("computing transfers to reach balance",
		"pendingTransfers", pendingTransfers,
		"liquidityGraph", g,
		"resolvedTransfersQuorum", resolvedTransfersQuorum,
	)
	inflightTransfers := append(pendingTransfers, toPending(resolvedTransfersQuorum)...)
	proposedTransfers, err := p.liquidityRebalancer.ComputeTransfersToBalance(g, inflightTransfers)
	if err != nil {
		return nil, fmt.Errorf("compute transfers to reach balance: %w", err)
	}

	lggr.Infow("finished computing outcome",
		"medianLiquidityPerChain", medianLiquidityPerChain,
		"pendingTransfers", pendingTransfers,
		"proposedTransfers", proposedTransfers,
		"resolvedTransfers", resolvedTransfersQuorum,
	)

	return models.NewOutcome(proposedTransfers, resolvedTransfersQuorum, pendingTransfers, configDigests).Encode(), nil
}

func (p *Plugin) Reports(seqNr uint64, outcome ocr3types.Outcome) ([]ocr3types.ReportWithInfo[models.Report], error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	lggr := p.lggr.With("seqNr", seqNr)
	lggr.Infow("in reports")

	decodedOutcome, err := models.DecodeOutcome(outcome)
	if err != nil {
		return nil, fmt.Errorf("decode outcome: %w", err)
	}

	// get all incoming and outgoing transfers for each network
	// incoming transfers will need to be finalized
	// outgoing transfers will need to be executed
	incomingAndOutgoing := make(map[models.NetworkSelector][]models.Transfer)
	for _, networkID := range p.rebalancerGraph.GetNetworks() {
		for _, outgoing := range decodedOutcome.ResolvedTransfers {
			if outgoing.From == networkID {
				incomingAndOutgoing[networkID] = append(incomingAndOutgoing[networkID], outgoing)
			}
		}
		for _, incoming := range decodedOutcome.PendingTransfers {
			if incoming.To == networkID &&
				(incoming.Status == models.TransferStatusReady ||
					incoming.Status == models.TransferStatusFinalized) {
				incomingAndOutgoing[networkID] = append(incomingAndOutgoing[networkID], incoming.Transfer)
			}
		}
	}

	lggr = lggr.With(
		"incomingAndOutgoing", incomingAndOutgoing,
		"resolvedTransfers", decodedOutcome.ResolvedTransfers,
		"pendingTransfers", decodedOutcome.PendingTransfers,
		"proposedTransfers", decodedOutcome.ProposedTransfers,
	)
	lggr.Debugw("got incoming and outgoing transfers")

	configDigestsMap := map[models.NetworkSelector]map[models.Address]types.ConfigDigest{}
	for _, cd := range decodedOutcome.ConfigDigests {
		_, found := configDigestsMap[cd.NetworkSel]
		if found {
			return nil, fmt.Errorf("found duplicate config digest for %v", cd.NetworkSel)
		}
		configDigestsMap[cd.NetworkSel] = map[models.Address]types.ConfigDigest{
			cd.RebalancerAddr: cd.Digest.ConfigDigest,
		}
	}

	var reports []ocr3types.ReportWithInfo[models.Report]
	for networkID, transfers := range incomingAndOutgoing {
		lmAddress, err := p.rebalancerGraph.GetRebalancerAddress(networkID)
		if err != nil {
			return nil, fmt.Errorf("liquidity manager for %v does not exist", networkID)
		}

		configDigests, found := configDigestsMap[networkID]
		if !found {
			return nil, fmt.Errorf("cannot find config digest for %v", networkID)
		}
		configDigest, found := configDigests[lmAddress]
		if !found {
			return nil, fmt.Errorf("cannot find config digest for %v:%s", networkID, lmAddress)
		}

		report := models.NewReport(transfers, lmAddress, networkID, configDigest)
		encoded, err := report.OnchainEncode()
		if err != nil {
			return nil, fmt.Errorf("encode report metadata for onchain usage: %w", err)
		}
		reports = append(reports, ocr3types.ReportWithInfo[models.Report]{
			Report: encoded,
			Info:   report,
		})
	}

	lggr.Infow("generated reports", "numReports", len(reports))
	return reports, nil
}

func (p *Plugin) ShouldAcceptAttestedReport(ctx context.Context, seqNr uint64, r ocr3types.ReportWithInfo[models.Report]) (bool, error) {
	lggr := p.lggr.With("seqNr", seqNr, "reportMeta", r.Info, "reportHex", hexutil.Encode(r.Report), "reportLen", len(r.Report))
	lggr.Infow("in should accept attested report")

	report, instructions, err := models.DecodeReport(p.rootNetwork, p.rootAddress, r.Report)
	if err != nil {
		return false, fmt.Errorf("failed to decode report: %w", err)
	}

	lggr.Infow("accepting report",
		"transfers", len(report.Transfers),
		"sendInstructions", instructions.SendLiquidityParams,
		"receiveInstructions", instructions.ReceiveLiquidityParams)
	// todo: check if reportMeta.transfers are valid

	return true, nil
}

func (p *Plugin) ShouldTransmitAcceptedReport(ctx context.Context, seqNr uint64, r ocr3types.ReportWithInfo[models.Report]) (bool, error) {
	lggr := p.lggr.With("seqNr", seqNr, "reportMeta", r.Info, "reportHex", hexutil.Encode(r.Report), "reportLen", len(r.Report))
	lggr.Infow("in should transmit accepted report")

	report, instructions, err := models.DecodeReport(r.Info.NetworkID, r.Info.LiquidityManagerAddress, r.Report)
	if err != nil {
		return false, fmt.Errorf("failed to decode report: %w", err)
	}

	lggr.Infow("in should transmit accepted report",
		"transfers", len(report.Transfers),
		"sendInstructions", instructions.SendLiquidityParams,
		"receiveInstructions", instructions.ReceiveLiquidityParams)

	rebalancer, err := p.liquidityManagerFactory.NewRebalancer(r.Info.NetworkID, r.Info.LiquidityManagerAddress)
	if err != nil {
		return false, fmt.Errorf("init liquidity manager: %w", err)
	}

	// check sequence number to see if its already transmitted onchain
	latestSeqNr, err := rebalancer.GetLatestSequenceNumber(ctx)
	if err != nil {
		return false, fmt.Errorf("get latest sequence number: %w", err)
	}

	if latestSeqNr >= seqNr {
		lggr.Debugw("report already transmitted onchain, returning false", "latestSeqNr", latestSeqNr)
		return false, nil
	}

	lggr.Infow("transmitting accepted report", "latestSeqNr", latestSeqNr)

	return true, nil
}

func (p *Plugin) Close() error {
	p.lggr.Infow("closing plugin")
	ctx, cf := context.WithTimeout(context.Background(), p.closePluginTimeout)
	defer cf()

	var errs []error
	for _, networkID := range p.rebalancerGraph.GetNetworks() {
		rebalancerAddress, err := p.rebalancerGraph.GetRebalancerAddress(networkID)
		if err != nil {
			errs = append(errs, fmt.Errorf("get rebalancer address for %v: %w", networkID, err))
			continue
		}

		rb, err := p.liquidityManagerFactory.GetRebalancer(networkID, rebalancerAddress)
		if err != nil {
			errs = append(errs, fmt.Errorf("get rebalancer (%d, %v): %w", networkID, rebalancerAddress, err))
			continue
		}

		if err := rb.Close(ctx); err != nil {
			errs = append(errs, fmt.Errorf("close rebalancer (%d, %v): %w", networkID, rebalancerAddress, err))
			continue
		}
	}

	return multierr.Combine(errs...)
}

func (p *Plugin) syncGraphEdges(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// todo: if there wasn't any change to the graph stop earlier
	p.lggr.Infow("syncing graph edges")

	p.lggr.Infow("discovering rebalancers")
	discoverer, err := p.discovererFactory.NewDiscoverer(p.rootNetwork, p.rootAddress)
	if err != nil {
		return fmt.Errorf("init discoverer: %w", err)
	}

	g, err := discoverer.Discover(ctx)
	if err != nil {
		return fmt.Errorf("discovering rebalancers: %w", err)
	}

	p.rebalancerGraph = g

	p.lggr.Infow("finished syncing graph edges", "graph", g.String())

	return nil
}

func (p *Plugin) syncGraphBalances(ctx context.Context) ([]models.NetworkLiquidity, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	networks := p.rebalancerGraph.GetNetworks()
	p.lggr.Infow("syncing graph balances", "networks", networks)

	networkLiquidities := make([]models.NetworkLiquidity, 0, len(networks))
	for _, networkID := range networks {
		lmAddr, err := p.rebalancerGraph.GetRebalancerAddress(networkID)
		if err != nil {
			return nil, fmt.Errorf("liquidity manager for network %v was not found", networkID)
		}

		lm, err := p.liquidityManagerFactory.NewRebalancer(networkID, lmAddr)
		if err != nil {
			return nil, fmt.Errorf("init liquidity manager: %w", err)
		}

		balance, err := lm.GetBalance(ctx)
		if err != nil {
			return nil, fmt.Errorf("get %v balance: %w", networkID, err)
		}

		p.rebalancerGraph.SetLiquidity(networkID, balance)
		networkLiquidities = append(networkLiquidities, models.NewNetworkLiquidity(networkID, balance))
	}

	return networkLiquidities, nil
}

func (p *Plugin) loadPendingTransfers(ctx context.Context) ([]models.PendingTransfer, error) {
	p.lggr.Infow("loading pending transfers")

	pendingTransfers := make([]models.PendingTransfer, 0)
	for _, networkID := range p.rebalancerGraph.GetNetworks() {
		neighbors, ok := p.rebalancerGraph.GetNeighbors(networkID)
		if !ok {
			p.lggr.Warnw("no neighbors found for network", "network", networkID)
			continue
		}

		// todo: figure out what to do with this
		// dateToStartLookingFrom := time.Now().Add(-10 * 24 * time.Hour)

		// if mostRecentTransfer, exists := p.pendingTransfers.LatestNetworkTransfer(networkID); exists {
		// 	dateToStartLookingFrom = mostRecentTransfer.Date
		// }

		for _, neighbor := range neighbors {
			bridge, err := p.bridgeFactory.NewBridge(networkID, neighbor)
			if err != nil {
				return nil, fmt.Errorf("init bridge: %w", err)
			}

			if bridge == nil {
				p.lggr.Warnw("no bridge found for network pair", "sourceNetwork", networkID, "destNetwork", neighbor)
				continue
			}

			localToken, err := p.rebalancerGraph.GetTokenAddress(networkID)
			if err != nil {
				return nil, fmt.Errorf("get local token address for %v: %w", networkID, err)
			}
			remoteToken, err := p.rebalancerGraph.GetTokenAddress(neighbor)
			if err != nil {
				return nil, fmt.Errorf("get remote token address for %v: %w", neighbor, err)
			}

			netPendingTransfers, err := bridge.GetTransfers(ctx, localToken, remoteToken)
			if err != nil {
				return nil, fmt.Errorf("get pending transfers: %w", err)
			}

			p.lggr.Infow("loaded pending transfers", "network", networkID, "pendingTransfers", netPendingTransfers)
			pendingTransfers = append(pendingTransfers, netPendingTransfers...)
		}
	}

	p.pendingTransfers.Add(pendingTransfers)
	return pendingTransfers, nil
}

func (p *Plugin) computeMedianLiquidityPerChain(observations []models.Observation) []models.NetworkLiquidity {
	liqObsPerChain := make(map[models.NetworkSelector][]*big.Int)
	for _, ob := range observations {
		for _, chainLiq := range ob.LiquidityPerChain {
			liqObsPerChain[chainLiq.Network] = append(liqObsPerChain[chainLiq.Network], chainLiq.Liquidity.ToInt())
		}
	}

	medians := make([]models.NetworkLiquidity, 0, len(liqObsPerChain))
	for chainID, liqs := range liqObsPerChain {
		medians = append(medians, models.NewNetworkLiquidity(chainID, bigIntSortedMiddle(liqs)))
	}
	// sort by network id for deterministic results
	sort.Slice(medians, func(i, j int) bool {
		return medians[i].Network < medians[j].Network
	})
	return medians
}

func (p *Plugin) computePendingTransfersConsensus(observations []models.Observation) ([]models.PendingTransfer, error) {
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

	var quorumEvents []models.PendingTransfer
	for h, count := range counts {
		if count >= p.f+1 {
			ev, exists := eventFromHash[h]
			if !exists {
				return nil, fmt.Errorf("internal issue, event from hash %v not found", h)
			}
			quorumEvents = append(quorumEvents, ev)
		}
	}

	return quorumEvents, nil
}

func (p *Plugin) computeConfigDigestsConsensus(observations []models.Observation) ([]models.ConfigDigestWithMeta, error) {
	key := func(meta models.ConfigDigestWithMeta) string {
		return fmt.Sprintf("%d-%s-%s", meta.NetworkSel, meta.RebalancerAddr, meta.Digest.Hex())
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

	var quorumCds []models.ConfigDigestWithMeta
	for k, count := range counts {
		if count >= p.f+1 {
			cd, exists := cds[k]
			if !exists {
				return nil, fmt.Errorf("internal issue, config digest by key %s not found", k)
			}
			quorumCds = append(quorumCds, cd)
		}
	}

	return quorumCds, nil
}

func (p *Plugin) computeGraphEdgesConsensus(observations []models.Observation) []models.Edge {
	counts := make(map[models.Edge]int)
	for _, obs := range observations {
		for _, edge := range obs.Edges {
			counts[edge]++
		}
	}

	var quorumEdges []models.Edge
	for edge, count := range counts {
		if count >= p.f+1 {
			quorumEdges = append(quorumEdges, edge)
		}
	}

	return quorumEdges
}

// computeMedianGraph computes a graph with the provided median liquidities per chain and edges that quorum agreed on.
func (p *Plugin) computeMedianGraph(
	edges []models.Edge, medianLiquidities []models.NetworkLiquidity) (graph.Graph, error) {

	g, err := graph.NewGraphFromEdges(edges)
	if err != nil {
		return nil, fmt.Errorf("new graph from edges: %w", err)
	}

	for _, medianLiq := range medianLiquidities {
		if !g.SetLiquidity(medianLiq.Network, medianLiq.Liquidity.ToInt()) {
			p.lggr.Errorw("median liquidity on network not found on edges quorum", "net", medianLiq.Network)
		}
	}

	return g, nil
}

func (p *Plugin) computeResolvedTransfersQuorum(observations []models.Observation) ([]models.Transfer, error) {
	// assumption: there shouldn't be more than 1 transfer for a (from, to) pair from a single oracle's observation.
	// otherwise they can be collapsed into a single transfer.
	// TODO: we can check for this in ValidateObservation
	type key struct {
		From               models.NetworkSelector
		To                 models.NetworkSelector
		AmountString       string
		Sender             models.Address
		Receiver           models.Address
		LocalTokenAddress  models.Address
		RemoteTokenAddress models.Address
	}
	counts := make(map[key][]models.Transfer)
	for _, obs := range observations {
		p.lggr.Debugw("observed transfers", "transfers", obs.ResolvedTransfers)
		for _, tr := range obs.ResolvedTransfers {
			p.lggr.Debugw("inserting resolved transfer into mapping", "transfer", tr)
			k := key{
				From:               tr.From,
				To:                 tr.To,
				AmountString:       tr.Amount.String(),
				Sender:             tr.Sender,
				Receiver:           tr.Receiver,
				LocalTokenAddress:  tr.LocalTokenAddress,
				RemoteTokenAddress: tr.RemoteTokenAddress,
			}
			counts[k] = append(counts[k], tr)
		}
	}

	p.lggr.Debugw("resolved transfers counts", "counts", len(counts))

	var quorumTransfers []models.Transfer
	for k, transfers := range counts {
		if len(transfers) >= p.f+1 {
			p.lggr.Debugw("quorum reached on transfer", "transfer", k, "votes", len(transfers))
			// need to compute the "medianized" bridge payload
			// only the bridge knows how to do this so we need to delegate it to them
			// the native bridge fee can also be medianized, no need for the bridge to do that
			var (
				bridgeFees     []*big.Int
				bridgePayloads [][]byte
				datesUnix      []*big.Int
			)
			for _, tr := range transfers {
				bridgeFees = append(bridgeFees, tr.NativeBridgeFee.ToInt())
				bridgePayloads = append(bridgePayloads, tr.BridgeData)
				datesUnix = append(datesUnix, big.NewInt(tr.Date.Unix()))
			}
			medianizedNativeFee := bigIntSortedMiddle(bridgeFees)
			medianizedDateUnix := bigIntSortedMiddle(datesUnix)
			bridge, err := p.bridgeFactory.NewBridge(k.From, k.To)
			if err != nil {
				return nil, fmt.Errorf("init bridge: %w", err)
			}
			quorumizedBridgePayload, err := bridge.QuorumizedBridgePayload(bridgePayloads, p.f)
			if err != nil {
				return nil, fmt.Errorf("quorumized bridge payload: %w", err)
			}
			quorumTransfer := models.Transfer{
				From:               k.From,
				To:                 k.To,
				Amount:             transfers[0].Amount,
				Date:               time.Unix(medianizedDateUnix.Int64(), 0), // medianized, not in the key
				Sender:             transfers[0].Sender,
				Receiver:           transfers[0].Receiver,
				LocalTokenAddress:  transfers[0].LocalTokenAddress,
				RemoteTokenAddress: transfers[0].RemoteTokenAddress,
				BridgeData:         quorumizedBridgePayload,       // "quorumized", not in the key
				NativeBridgeFee:    ubig.New(medianizedNativeFee), // medianized, not in the key
			}
			quorumTransfers = append(quorumTransfers, quorumTransfer)
		} else {
			p.lggr.Debugw("dropping transfer, not enough votes on it", "transfer", k, "votes", len(transfers))
		}
	}

	return quorumTransfers, nil
}

func (p *Plugin) resolveProposedTransfers(ctx context.Context, outcomeCtx ocr3types.OutcomeContext) ([]models.Transfer, error) {
	p.lggr.Infow("resolving proposed transfers", "seqNr", outcomeCtx.SeqNr, "prevSeqNr", outcomeCtx.SeqNr-1)

	if len(outcomeCtx.PreviousOutcome) == 0 {
		return nil, nil
	}

	outcome, err := models.DecodeOutcome(outcomeCtx.PreviousOutcome)
	if err != nil {
		return nil, fmt.Errorf("decode previous outcome: %w", err)
	}

	var resolvedTransfers []models.Transfer
	for _, proposedTransfer := range outcome.ProposedTransfers {
		bridge, err := p.bridgeFactory.NewBridge(proposedTransfer.From, proposedTransfer.To)
		if err != nil {
			return nil, fmt.Errorf("init bridge: %w", err)
		}

		fromNetRebalancer, err := p.rebalancerGraph.GetRebalancerAddress(proposedTransfer.From)
		if err != nil {
			return nil, fmt.Errorf("get rebalancer address for %v: %w", proposedTransfer.From, err)
		}

		fromNetToken, err := p.rebalancerGraph.GetTokenAddress(proposedTransfer.From)
		if err != nil {
			return nil, fmt.Errorf("get token address for %v: %w", proposedTransfer.From, err)
		}

		toNetRebalancer, err := p.rebalancerGraph.GetRebalancerAddress(proposedTransfer.To)
		if err != nil {
			return nil, fmt.Errorf("get rebalancer address for %v: %w", proposedTransfer.To, err)
		}

		toNetToken, err := p.rebalancerGraph.GetTokenAddress(proposedTransfer.To)
		if err != nil {
			return nil, fmt.Errorf("get token address for %v: %w", proposedTransfer.To, err)
		}

		resolvedTransfer := models.Transfer{
			From:               proposedTransfer.From,
			To:                 proposedTransfer.To,
			Amount:             proposedTransfer.Amount,
			Sender:             fromNetRebalancer,
			Receiver:           toNetRebalancer,
			LocalTokenAddress:  fromNetToken,
			RemoteTokenAddress: toNetToken,
			// BridgeData: nil, // will be filled in below
			// NativeBridgeFee: big.NewInt(0), // will be filled in below
		}

		bridgePayload, bridgeFee, err := bridge.GetBridgePayloadAndFee(ctx, resolvedTransfer)
		if err != nil {
			p.lggr.Warnw("failed to get bridge payload and fee", "proposedTransfer", proposedTransfer, "err", err)
			// return nil, fmt.Errorf("get bridge payload and fee: %w", err)
			continue
		}
		resolvedTransfer.BridgeData = bridgePayload
		resolvedTransfer.NativeBridgeFee = ubig.New(bridgeFee)
		resolvedTransfers = append(resolvedTransfers, resolvedTransfer)
	}

	p.lggr.Infow("finished resolving proposed transfers", "resolvedTransfers", resolvedTransfers)

	return resolvedTransfers, nil
}

// bigIntSortedMiddle returns the middle number after sorting the provided numbers. nil is returned if the provided slice is empty.
// If length of the provided slice is even, the right-hand-side value of the middle 2 numbers is returned.
// The objective of this function is to always pick within the range of values reported by honest nodes when we have 2f+1 values.
// todo: move to libs
func bigIntSortedMiddle(vals []*big.Int) *big.Int {
	if len(vals) == 0 {
		return nil
	}

	valsCopy := make([]*big.Int, len(vals))
	copy(valsCopy[:], vals[:])
	sort.Slice(valsCopy, func(i, j int) bool {
		return valsCopy[i].Cmp(valsCopy[j]) == -1
	})
	return valsCopy[len(valsCopy)/2]
}

func toPending(ts []models.Transfer) []models.PendingTransfer {
	var pts []models.PendingTransfer
	for _, t := range ts {
		pts = append(pts, models.NewPendingTransfer(t))
	}
	return pts
}
