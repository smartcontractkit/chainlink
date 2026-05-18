package liquiditymanager

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"go.uber.org/multierr"

	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/bridge"
	evmliquiditymanager "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/chain/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/discoverer"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/graph"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/inflight"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/rebalalgo"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/rebalcalc"
)

type Plugin struct {
	f                       int
	closePluginTimeout      time.Duration
	liquidityManagerFactory evmliquiditymanager.Factory
	discoverer              discoverer.Discoverer
	bridgeFactory           bridge.Factory
	mu                      sync.RWMutex
	liquidityGraph          graph.Graph
	liquidityRebalancer     rebalalgo.RebalancingAlgo
	inflight                inflight.Container
	lggr                    logger.Logger
	reportCodec             evmliquiditymanager.OnchainReportCodec
}

func NewPlugin(
	f int,
	closePluginTimeout time.Duration,
	rootNetwork models.NetworkSelector,
	rootAddress models.Address,
	liquidityManagerFactory evmliquiditymanager.Factory,
	discoverer discoverer.Discoverer,
	bridgeFactory bridge.Factory,
	liquidityRebalancer rebalalgo.RebalancingAlgo,
	reportCodec evmliquiditymanager.OnchainReportCodec,
	lggr logger.Logger,
) *Plugin {
	return &Plugin{
		f:                       f,
		closePluginTimeout:      closePluginTimeout,
		liquidityManagerFactory: liquidityManagerFactory,
		bridgeFactory:           bridgeFactory,
		discoverer:              discoverer,
		liquidityGraph:          graph.NewGraph(),
		liquidityRebalancer:     liquidityRebalancer,
		inflight:                inflight.New(),
		reportCodec:             reportCodec,
		lggr:                    lggr,
		mu:                      sync.RWMutex{},
	}
}

func (p *Plugin) Query(_ context.Context, outcomeCtx ocr3types.OutcomeContext) (ocrtypes.Query, error) {
	p.lggr.Infow("in query", "seqNr", outcomeCtx.SeqNr)
	return ocrtypes.Query{}, nil
}

func (p *Plugin) Observation(ctx context.Context, outcomeCtx ocr3types.OutcomeContext, _ ocrtypes.Query) (ocrtypes.Observation, error) {
	lggr := p.lggr.With("seqNr", outcomeCtx.SeqNr, "phase", "Observation")
	lggr.Infow("in observation", "seqNr", outcomeCtx.SeqNr)

	if err := p.syncGraph(ctx); err != nil {
		return ocrtypes.Observation{}, fmt.Errorf("sync graph edges: %w", err)
	}

	networkLiquidities := make([]models.NetworkLiquidity, 0)
	for _, net := range p.liquidityGraph.GetNetworks() {
		liq, err := p.liquidityGraph.GetLiquidity(net)
		if err != nil {
			return ocrtypes.Observation{}, err
		}
		networkLiquidities = append(networkLiquidities, models.NewNetworkLiquidity(net, liq))
	}

	pendingTransfers, err := p.loadPendingTransfers(ctx, lggr)
	if err != nil {
		return ocrtypes.Observation{}, fmt.Errorf("load pending transfers: %w", err)
	}

	numExpired := p.inflight.Expire(pendingTransfers)
	inflightTransfers := p.inflight.GetAll()

	edges, err := p.liquidityGraph.GetEdges()
	if err != nil {
		return ocrtypes.Observation{}, fmt.Errorf("get edges: %w", err)
	}

	resolvedTransfers, err := p.resolveProposedTransfers(ctx, lggr, outcomeCtx)
	if err != nil {
		return ocrtypes.Observation{}, fmt.Errorf("resolve proposed transfers: %w", err)
	}

	configDigests := make([]models.ConfigDigestWithMeta, 0)
	for _, net := range p.liquidityGraph.GetNetworks() {
		data, err := p.liquidityGraph.GetData(net)
		if err != nil {
			return nil, fmt.Errorf("get rb %d data: %w", net, err)
		}
		configDigests = append(configDigests, models.ConfigDigestWithMeta{
			Digest:     data.ConfigDigest,
			NetworkSel: data.NetworkSelector,
		})
	}

	lggr.Infow("finished observing",
		"networkLiquidities", networkLiquidities,
		"pendingTransfers", pendingTransfers,
		"edges", edges,
		"resolvedTransfers", resolvedTransfers,
		"inflightTransfers", inflightTransfers,
		"numExpired", numExpired,
	)

	return models.NewObservation(
		networkLiquidities,
		resolvedTransfers,
		pendingTransfers,
		inflightTransfers,
		edges,
		configDigests).Encode()
}

func (p *Plugin) ObservationQuorum(outctx ocr3types.OutcomeContext, query ocrtypes.Query) (ocr3types.Quorum, error) {
	return ocr3types.QuorumTwoFPlusOne, nil
}

func (p *Plugin) Outcome(outctx ocr3types.OutcomeContext, query ocrtypes.Query, aos []ocrtypes.AttributedObservation) (ocr3types.Outcome, error) {
	lggr := p.lggr.With("seqNr", outctx.SeqNr, "numObservations", len(aos), "phase", "Outcome")
	lggr.Infow("in outcome")

	// Gather all the observations.
	observations := make([]models.Observation, 0, len(aos))
	for _, encodedObs := range aos {
		obs, err := models.DecodeObservation(encodedObs.Observation)
		if err != nil {
			return ocr3types.Outcome{}, fmt.Errorf("decode observation: %w", err)
		}
		lggr.Debugw("decoded observation from oracle", "observation", obs, "oracleID", encodedObs.Observer)
		observations = append(observations, obs)
	}

	// Come to a consensus based on the observations of all the different nodes.
	medianLiquidityPerChain, err := rebalcalc.MedianLiquidityPerChain(observations, p.f)
	if err != nil {
		return ocr3types.Outcome{}, fmt.Errorf("compute median liquidity per chain: %w", err)
	}

	graphEdges, err := rebalcalc.GraphEdgesConsensus(observations, p.f)
	if err != nil {
		return ocr3types.Outcome{}, fmt.Errorf("compute graph edges consensus: %w", err)
	}

	pendingTransfers, err := rebalcalc.PendingTransfersConsensus(observations, p.f)
	if err != nil {
		return ocr3types.Outcome{}, fmt.Errorf("compute pending transfers consensus: %w", err)
	}

	configDigests, err := rebalcalc.ConfigDigestsConsensus(observations, p.f)
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

	inflightTransfers, err := rebalcalc.InflightTransfersConsensus(observations, p.f)
	if err != nil {
		return nil, fmt.Errorf("compute inflight transfers consensus: %w", err)
	}

	lggr.Infow("computing transfers to reach balance",
		"pendingTransfers", pendingTransfers,
		"liquidityGraph", g,
		"resolvedTransfersQuorum", resolvedTransfersQuorum,
		"inflightTransfers", inflightTransfers,
	)
	proposedTransfers, err := p.liquidityRebalancer.ComputeTransfersToBalance(g, combinedUnexecutedTransfers(pendingTransfers, resolvedTransfersQuorum, inflightTransfers))
	if err != nil {
		return nil, fmt.Errorf("compute transfers to reach balance: %w", err)
	}

	lggr.Infow("finished computing outcome",
		"medianLiquidityPerChain", medianLiquidityPerChain,
		"pendingTransfers", pendingTransfers,
		"proposedTransfers", proposedTransfers,
		"resolvedTransfers", resolvedTransfersQuorum,
	)

	return models.NewOutcome(proposedTransfers, resolvedTransfersQuorum, pendingTransfers, configDigests).Encode()
}

func combinedUnexecutedTransfers(
	pendingTransfers []models.PendingTransfer,
	resolvedTransfersQuorum []models.Transfer,
	inflightTransfers []models.Transfer,
) []rebalalgo.UnexecutedTransfer {
	unexecuted := make([]rebalalgo.UnexecutedTransfer, 0, len(resolvedTransfersQuorum)+len(inflightTransfers)+len(pendingTransfers))
	for _, resolvedTransfer := range resolvedTransfersQuorum {
		unexecuted = append(unexecuted, resolvedTransfer)
	}
	for _, inflightTransfer := range inflightTransfers {
		unexecuted = append(unexecuted, inflightTransfer)
	}
	for _, pendingTransfer := range pendingTransfers {
		unexecuted = append(unexecuted, pendingTransfer)
	}
	return unexecuted
}

func (p *Plugin) Reports(seqNr uint64, outcome ocr3types.Outcome) ([]ocr3types.ReportWithInfo[models.Report], error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	lggr := p.lggr.With("seqNr", seqNr, "phase", "Reports")
	lggr.Infow("in reports")

	decodedOutcome, err := models.DecodeOutcome(outcome)
	if err != nil {
		return nil, fmt.Errorf("decode outcome: %w", err)
	}

	// get all incoming and outgoing transfers for each network
	// incoming transfers will need to be finalized
	// outgoing transfers will need to be executed
	incomingAndOutgoing := make(map[models.NetworkSelector][]models.Transfer)

	for _, outgoing := range decodedOutcome.ResolvedTransfers {
		incomingAndOutgoing[outgoing.From] = append(incomingAndOutgoing[outgoing.From], outgoing)
	}

	for _, incoming := range decodedOutcome.PendingTransfers {
		if incoming.Status == models.TransferStatusReady ||
			incoming.Status == models.TransferStatusFinalized {
			incomingAndOutgoing[incoming.To] = append(incomingAndOutgoing[incoming.To], incoming.Transfer)
		}
	}

	lggr = lggr.With(
		"incomingAndOutgoing", incomingAndOutgoing,
		"resolvedTransfers", decodedOutcome.ResolvedTransfers,
		"pendingTransfers", decodedOutcome.PendingTransfers,
		"proposedTransfers", decodedOutcome.ProposedTransfers,
	)
	lggr.Debugw("generated incoming and outgoing transfers")

	configDigestsMap := map[models.NetworkSelector]types.ConfigDigest{}
	for _, cd := range decodedOutcome.ConfigDigests {
		_, found := configDigestsMap[cd.NetworkSel]
		if found {
			return nil, fmt.Errorf("found duplicate config digest for %v", cd.NetworkSel)
		}
		configDigestsMap[cd.NetworkSel] = cd.Digest.ConfigDigest
	}

	var reports []ocr3types.ReportWithInfo[models.Report]
	for networkID, transfers := range incomingAndOutgoing {
		// todo: we shouldn't use plugin state
		rebalancerAddress, err := p.liquidityGraph.GetLiquidityManagerAddress(networkID)
		if err != nil {
			return nil, fmt.Errorf("liquidity manager for %v does not exist", networkID)
		}

		configDigest, found := configDigestsMap[networkID]
		if !found {
			return nil, fmt.Errorf("cannot find config digest for %v", networkID)
		}

		report := models.NewReport(transfers, rebalancerAddress, networkID, configDigest)
		encoded, err := p.reportCodec.Encode(report)
		if err != nil {
			return nil, fmt.Errorf("encode report metadata for onchain usage: %w", err)
		}
		reports = append(reports, ocr3types.ReportWithInfo[models.Report]{
			Report: encoded,
			Info:   report,
		})
	}
	sort.Slice(reports, func(i, j int) bool { return reports[i].Info.NetworkID < reports[j].Info.NetworkID })

	lggr.Infow("generated reports", "numReports", len(reports))
	return reports, nil
}

func (p *Plugin) ShouldAcceptAttestedReport(ctx context.Context, seqNr uint64, r ocr3types.ReportWithInfo[models.Report]) (bool, error) {
	lggr := p.lggr.With("seqNr", seqNr, "reportMeta", r.Info, "reportLen", len(r.Report), "phase", "ShouldAcceptAttestedReport")
	lggr.Infow("in should accept attested report")

	report, instructions, err := p.reportCodec.Decode(r.Info.NetworkID, r.Info.LiquidityManagerAddress, r.Report)
	if err != nil {
		return false, fmt.Errorf("failed to decode report: %w", err)
	}

	if report.IsEmpty() {
		lggr.Infow("report has no transfers, should not be transmitted")
		return false, nil
	}

	staleErr := p.stalenessValidation(ctx, lggr, seqNr, r.Info)
	if staleErr != nil {
		lggr.Infow("report is stale, should not be accepted", "staleErr", staleErr)
		return false, nil
	}

	// check if any of the transfers in the report are in-flight.
	for _, transfer := range report.Transfers {
		if p.inflight.IsInflight(transfer) {
			lggr.Infow("transfer is in-flight, should not be accepted", "transfer", transfer)
			return false, nil
		}
	}

	// add the transfers to the inflight container since none of them are inflight already.
	// need to range through the meta because it has the correct stage information.
	for _, transfer := range r.Info.Transfers {
		p.inflight.Add(transfer)
	}

	lggr.Infow("accepting report",
		"transfers", len(report.Transfers),
		"sendInstructions", instructions.SendLiquidityParams,
		"receiveInstructions", instructions.ReceiveLiquidityParams)

	return true, nil
}

func (p *Plugin) ShouldTransmitAcceptedReport(ctx context.Context, seqNr uint64, r ocr3types.ReportWithInfo[models.Report]) (bool, error) {
	lggr := p.lggr.With("seqNr", seqNr, "reportMeta", r.Info, "reportLen", len(r.Report), "phase", "ShouldTransmitAcceptedReport")
	lggr.Infow("in should transmit accepted report")

	report, instructions, err := p.reportCodec.Decode(r.Info.NetworkID, r.Info.LiquidityManagerAddress, r.Report)
	if err != nil {
		return false, fmt.Errorf("failed to decode report: %w", err)
	}

	if report.IsEmpty() {
		lggr.Infow("report has no transfers, should not be transmitted")
		return false, nil
	}

	lggr.Infow("in should transmit accepted report",
		"transfers", len(report.Transfers),
		"sendInstructions", instructions.SendLiquidityParams,
		"receiveInstructions", instructions.ReceiveLiquidityParams)

	staleErr := p.stalenessValidation(ctx, lggr, seqNr, r.Info)
	if staleErr != nil {
		lggr.Infow("report is stale, should not be transmitted", "staleErr", staleErr)
		return false, nil
	}

	lggr.Infow("transmitting accepted report")

	return true, nil
}

// validateReportStaleness performs various checks on the report to determine whether or not it is stale.
// A stale report should not be transmitted on-chain.
func (p *Plugin) stalenessValidation(
	ctx context.Context,
	lggr logger.Logger,
	seqNr uint64,
	report models.Report,
) error {
	// check sequence number to see if its already transmitted onchain.
	liquidityManager, err := p.liquidityManagerFactory.NewLiquidityManager(report.NetworkID, report.LiquidityManagerAddress)
	if err != nil {
		return fmt.Errorf("get liquidityManager: %w", err)
	}

	onchainSeqNr, err := liquidityManager.GetLatestSequenceNumber(ctx)
	if err != nil {
		return fmt.Errorf("get latest sequence number: %w", err)
	}

	if onchainSeqNr >= seqNr {
		return fmt.Errorf("report already transmitted onchain, report seqNr: %d, onchain seqNr: %d", seqNr, onchainSeqNr)
	}

	lggr.Debugw("onchain sequence number < current", "onchainSeqNr", onchainSeqNr)

	// check that the instructions will not cause failures onchain.
	// e.g send instructions when there is not enough liquidity.
	currentBalance, err := liquidityManager.GetBalance(ctx)
	if err != nil {
		lggr.Warnw("failed to get balance", "err", err)
		return fmt.Errorf("get liquidityManager liquidity: %w", err)
	}

	lggr.Debugw("checking if there is enough balance onchain to send", "currentBalance", currentBalance.String())

	for _, transfer := range report.Transfers {
		if transfer.From != report.NetworkID {
			continue
		}

		if currentBalance.Cmp(transfer.Amount.ToInt()) < 0 {
			return fmt.Errorf("not enough balance onchain to send, amount: %s, remaining: %s", transfer.Amount.String(), currentBalance.String())
		}
		currentBalance = currentBalance.Sub(currentBalance, transfer.Amount.ToInt())
	}

	if currentBalance.Cmp(big.NewInt(0)) < 0 {
		return fmt.Errorf("not enough balance onchain to send, remaining: %s", currentBalance.String())
	}

	lggr.Debugw("enough balance onchain to send", "currentBalance", currentBalance.String())

	return nil
}

func (p *Plugin) Close() error {
	p.lggr.Infow("closing plugin")
	ctx, cf := context.WithTimeout(context.Background(), p.closePluginTimeout)
	defer cf()

	var errs []error
	for _, networkID := range p.liquidityGraph.GetNetworks() {
		lggr := p.lggr.With("network", networkID, "chainID", networkID.ChainID())
		lggr.Infow("closing liquidityManager network")

		liquidityManagerAddress, err := p.liquidityGraph.GetLiquidityManagerAddress(networkID)
		if err != nil {
			errs = append(errs, fmt.Errorf("get liquidityManager address for %d: %w", networkID, err))
			continue
		}

		rb, err := p.liquidityManagerFactory.GetLiquidityManager(networkID, liquidityManagerAddress)
		if err != nil {
			errs = append(errs, fmt.Errorf("get liquidityManager (%d, %s): %w", networkID, liquidityManagerAddress.String(), err))
			continue
		}

		if err := rb.Close(ctx); err != nil {
			errs = append(errs, fmt.Errorf("close liquidityManager (%d, %s): %w", networkID, liquidityManagerAddress.String(), err))
			continue
		}

		lggr.Infow("finished closing liquidityManager network", "liquidityManager", liquidityManagerAddress.String())
	}

	return multierr.Combine(errs...)
}

func (p *Plugin) syncGraphEdges(ctx context.Context) error {
	p.lggr.Infow("syncing graph edges")
	// todo: discoverer factory is not required we can pass a discoverer instance to the plugin
	p.lggr.Infow("discovering liquidity managers")
	g, err := p.discoverer.Discover(ctx)
	if err != nil {
		return fmt.Errorf("discovering rebalancers: %w", err)
	}
	p.lggr.Infow("finished syncing graph edges", "graph", g.String())
	p.liquidityGraph = g
	return nil
}

func (p *Plugin) syncGraph(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.liquidityGraph.IsEmpty() {
		if err := p.syncGraphEdges(ctx); err != nil {
			return fmt.Errorf("sync graph edges: %w", err)
		}
	} else {
		p.lggr.Infow("syncing graph liquidities")
		if err := p.discoverer.DiscoverBalances(ctx, p.liquidityGraph); err != nil {
			return fmt.Errorf("discovering balances: %w", err)
		}
		p.lggr.Infow("finished syncing graph liquidities")
	}

	return nil
}

func (p *Plugin) loadPendingTransfers(ctx context.Context, lggr logger.Logger) ([]models.PendingTransfer, error) {
	lggr.Infow("loading pending transfers")

	pendingTransfers := make([]models.PendingTransfer, 0)
	edges, err := p.liquidityGraph.GetEdges()
	if err != nil {
		return nil, fmt.Errorf("get edges: %w", err)
	}
	for _, edge := range edges {
		logger := lggr.With("sourceNetwork", edge.Source, "sourceChainID", edge.Source.ChainID(), "destNetwork", edge.Dest, "destChainID", edge.Dest.ChainID())
		bridge, err := p.bridgeFactory.NewBridge(ctx, edge.Source, edge.Dest)
		if err != nil {
			return nil, fmt.Errorf("init bridge: %w", err)
		}

		if bridge == nil {
			logger.Warn("no bridge found for network pair")
			continue
		}

		localToken, err := p.liquidityGraph.GetTokenAddress(edge.Source)
		if err != nil {
			return nil, fmt.Errorf("get local token address for %v: %w", edge.Source, err)
		}
		remoteToken, err := p.liquidityGraph.GetTokenAddress(edge.Dest)
		if err != nil {
			return nil, fmt.Errorf("get remote token address for %v: %w", edge.Dest, err)
		}

		netPendingTransfers, err := bridge.GetTransfers(ctx, localToken, remoteToken)
		if err != nil {
			return nil, fmt.Errorf("get pending transfers: %w", err)
		}

		logger.Infow("loaded pending transfers for network", "pendingTransfers", netPendingTransfers)
		pendingTransfers = append(pendingTransfers, netPendingTransfers...)
	}

	return pendingTransfers, nil
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
			p.lggr.Debugw("median liquidity on network not found on edges quorum", "net", medianLiq.Network)
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

	quorumTransfers := make([]models.Transfer, 0)
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
				datesUnix = append(datesUnix, big.NewInt(tr.Date.UTC().Unix()))
			}
			medianizedNativeFee := rebalcalc.BigIntSortedMiddle(bridgeFees)
			medianizedDateUnix := rebalcalc.BigIntSortedMiddle(datesUnix)
			bridge, err := p.bridgeFactory.GetBridge(k.From, k.To)
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
				Date:               time.Unix(medianizedDateUnix.Int64(), 0).In(time.UTC), // medianized, not in the key
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

func (p *Plugin) resolveProposedTransfers(ctx context.Context, lggr logger.Logger, outcomeCtx ocr3types.OutcomeContext) ([]models.Transfer, error) {
	lggr.Infow("resolving proposed transfers", "prevSeqNr", outcomeCtx.SeqNr-1)

	if len(outcomeCtx.PreviousOutcome) == 0 {
		return []models.Transfer{}, nil
	}

	outcome, err := models.DecodeOutcome(outcomeCtx.PreviousOutcome)
	if err != nil {
		return nil, fmt.Errorf("decode previous outcome: %w", err)
	}

	resolvedTransfers := make([]models.Transfer, 0, len(outcome.ProposedTransfers))
	for _, proposedTransfer := range outcome.ProposedTransfers {
		bridge, err := p.bridgeFactory.NewBridge(ctx, proposedTransfer.From, proposedTransfer.To)
		if err != nil {
			return nil, fmt.Errorf("init bridge: %w", err)
		}

		fromNetRebalancer, err := p.liquidityGraph.GetLiquidityManagerAddress(proposedTransfer.From)
		if err != nil {
			return nil, fmt.Errorf("get liquidityManager address for %v: %w", proposedTransfer.From, err)
		}

		fromNetToken, err := p.liquidityGraph.GetTokenAddress(proposedTransfer.From)
		if err != nil {
			return nil, fmt.Errorf("get token address for %v: %w", proposedTransfer.From, err)
		}

		toNetRebalancer, err := p.liquidityGraph.GetLiquidityManagerAddress(proposedTransfer.To)
		if err != nil {
			return nil, fmt.Errorf("get liquidityManager address for %v: %w", proposedTransfer.To, err)
		}

		toNetToken, err := p.liquidityGraph.GetTokenAddress(proposedTransfer.To)
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
			lggr.Warnw("failed to get bridge payload and fee", "proposedTransfer", proposedTransfer, "err", err)
			continue
		}
		resolvedTransfer.BridgeData = bridgePayload
		resolvedTransfer.NativeBridgeFee = ubig.New(bridgeFee)
		resolvedTransfer.Stage = 0
		resolvedTransfers = append(resolvedTransfers, resolvedTransfer)
	}

	lggr.Infow("finished resolving proposed transfers", "resolvedTransfers", resolvedTransfers)

	return resolvedTransfers, nil
}
