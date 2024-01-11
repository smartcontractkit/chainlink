package rebalancer

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/liquiditygraph"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/liquiditymanager"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/liquidityrebalancer"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
)

type Plugin struct {
	f                       int
	closePluginTimeout      time.Duration
	rootNetwork             models.NetworkID
	liquidityManagers       *liquiditymanager.Registry
	liquidityManagerFactory liquiditymanager.Factory
	liquidityGraph          liquiditygraph.LiquidityGraph
	liquidityRebalancer     liquidityrebalancer.Rebalancer
	pendingTransfers        *PendingTransfersCache
	lggr                    logger.Logger
}

func NewPlugin(
	f int,
	closePluginTimeout time.Duration,
	liquidityManagerNetwork models.NetworkID,
	liquidityManagerAddress models.Address,
	liquidityManagerFactory liquiditymanager.Factory,
	liquidityGraph liquiditygraph.LiquidityGraph,
	liquidityRebalancer liquidityrebalancer.Rebalancer,
	lggr logger.Logger,
) *Plugin {

	liquidityManagers := liquiditymanager.NewRegistry()
	liquidityManagers.Add(liquidityManagerNetwork, liquidityManagerAddress)

	return &Plugin{
		f:                       f,
		closePluginTimeout:      closePluginTimeout,
		rootNetwork:             liquidityManagerNetwork,
		liquidityManagers:       liquidityManagers,
		liquidityManagerFactory: liquidityManagerFactory,
		liquidityGraph:          liquidityGraph,
		liquidityRebalancer:     liquidityRebalancer,
		pendingTransfers:        NewPendingTransfersCache(),
		lggr:                    lggr,
	}
}

func (p *Plugin) Query(_ context.Context, outcomeCtx ocr3types.OutcomeContext) (ocrtypes.Query, error) {
	p.lggr.Infow("in query", "seqNr", outcomeCtx.SeqNr)
	return ocrtypes.Query{}, nil
}

func (p *Plugin) Observation(ctx context.Context, outcomeCtx ocr3types.OutcomeContext, _ ocrtypes.Query) (ocrtypes.Observation, error) {
	p.lggr.Infow("in observation", "seqNr", outcomeCtx.SeqNr)

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

	p.lggr.Infow("finished observing", "networkLiquidities", networkLiquidities, "pendingTransfers", pendingTransfers)

	return models.NewObservation(networkLiquidities, pendingTransfers).Encode(), nil
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
	lggr := p.lggr.With("seqNr", outctx.SeqNr)
	lggr.Infow("in outcome", "seqNr", outctx.SeqNr, "numObs", len(aos))

	observations := make([]models.Observation, 0, len(aos))
	for _, encodedObs := range aos {
		obs, err := models.DecodeObservation(encodedObs.Observation)
		if err != nil {
			return ocr3types.Outcome{}, fmt.Errorf("decode observation: %w", err)
		}
		observations = append(observations, obs)
	}

	medianLiquidityPerChain := p.computeMedianLiquidityPerChain(observations)

	pendingTransfers, err := p.computePendingTransfersConsensus(observations)
	if err != nil {
		return ocr3types.Outcome{}, fmt.Errorf("compute pending transfers consensus: %w", err)
	}

	lggr.Infow("finished computing outcome", "medianLiquidityPerChain", medianLiquidityPerChain, "pendingTransfers", pendingTransfers)

	return models.NewObservation(medianLiquidityPerChain, pendingTransfers).Encode(), nil
}

func (p *Plugin) Reports(seqNr uint64, outcome ocr3types.Outcome) ([]ocr3types.ReportWithInfo[models.ReportMetadata], error) {
	lggr := p.lggr.With("seqNr", seqNr)
	lggr.Infow("in reports", "seqNr", seqNr)
	obs, err := models.DecodeObservation(outcome)
	if err != nil {
		return nil, fmt.Errorf("decode outcome: %w", err)
	}

	lggr.Infow("computing transfers to reach balance",
		"pendingTransfers", obs.PendingTransfers,
		"liquidityGraph", p.liquidityGraph,
		"liquidityPerChain", obs.LiquidityPerChain)
	transfersToReachBalance, err := p.liquidityRebalancer.ComputeTransfersToBalance(
		p.liquidityGraph, obs.PendingTransfers, obs.LiquidityPerChain)
	if err != nil {
		return nil, fmt.Errorf("compute transfers to reach balance: %w", err)
	}

	// group transfers by source chain
	transfersBySourceNet := make(map[models.NetworkID][]models.Transfer)
	for _, tr := range transfersToReachBalance {
		transfersBySourceNet[tr.From] = append(transfersBySourceNet[tr.From], tr)
	}

	lggr.Infow("finished computing transfers to reach balance", "transfersBySourceNet", transfersBySourceNet)

	var reports []ocr3types.ReportWithInfo[models.ReportMetadata]
	for sourceNet, transfers := range transfersBySourceNet {
		lmAddress, exists := p.liquidityManagers.Get(sourceNet)
		if !exists {
			return nil, fmt.Errorf("liquidity manager for %v does not exist", sourceNet)
		}

		reportMeta := models.NewReportMetadata(transfers, lmAddress, sourceNet)
		encoded, err := reportMeta.OnchainEncode()
		if err != nil {
			return nil, fmt.Errorf("encode report metadata for onchain usage: %w", err)
		}
		reports = append(reports, ocr3types.ReportWithInfo[models.ReportMetadata]{
			Report: encoded,
			Info:   reportMeta,
		})
	}

	lggr.Infow("generated reports", "numReports", len(reports))

	return reports, nil
}

func (p *Plugin) ShouldAcceptAttestedReport(ctx context.Context, seqNr uint64, r ocr3types.ReportWithInfo[models.ReportMetadata]) (bool, error) {
	p.lggr.Infow("in should accept attested report", "seqNr", seqNr, "reportMeta", r.Info, "reportJSON", string(r.Report))

	report, err := models.DecodeReport(r.Report)
	if err != nil {
		return false, fmt.Errorf("decode report metadata: %w", err)
	}

	fmt.Println(report.Transfers)
	// todo: check if reportMeta.transfers are valid

	return true, nil
}

func (p *Plugin) ShouldTransmitAcceptedReport(ctx context.Context, seqNr uint64, r ocr3types.ReportWithInfo[models.ReportMetadata]) (bool, error) {
	p.lggr.Infow("in should transmit accepted report", "seqNr", seqNr, "reportMeta", r.Info)

	newPendingTransfers := make([]models.PendingTransfer, 0, len(r.Info.Transfers))
	for _, tr := range r.Info.Transfers {
		if p.pendingTransfers.ContainsTransfer(tr) {
			return false, nil
		}
		newPendingTransfers = append(newPendingTransfers, models.NewPendingTransfer(tr))
	}

	p.pendingTransfers.Add(newPendingTransfers)
	return true, nil
}

func (p *Plugin) Close() error {
	p.lggr.Infow("closing plugin")

	ctx, cf := context.WithTimeout(context.Background(), p.closePluginTimeout)
	defer cf()

	for networkID, lmAddr := range p.liquidityManagers.GetAll() {
		// todo: lmCloser := liquidityManagerFactory.NewLiquidityManagerCloser(); lmCloser.Close()
		lm, err := p.liquidityManagerFactory.NewLiquidityManager(networkID, lmAddr)
		if err != nil {
			return err
		}

		if err := lm.Close(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (p *Plugin) syncGraphEdges(ctx context.Context) error {
	// todo: if there wasn't any change to the graph stop earlier
	p.lggr.Infow("syncing graph edges")

	rootLM, exists := p.liquidityManagers.Get(p.rootNetwork)
	if !exists {
		return fmt.Errorf("root lm %v not found", p.rootNetwork)
	}

	lm, err := p.liquidityManagerFactory.NewLiquidityManager(p.rootNetwork, rootLM)
	if err != nil {
		return fmt.Errorf("init liquidity manager: %w", err)
	}

	lms, g, err := lm.Discover(ctx, p.liquidityManagerFactory)
	if err != nil {
		return fmt.Errorf("discover lms: %w", err)
	}

	p.liquidityGraph = g      // todo: thread safe
	p.liquidityManagers = lms // todo: thread safe

	return nil
}

func (p *Plugin) syncGraphBalances(ctx context.Context) ([]models.NetworkLiquidity, error) {
	networks := p.liquidityGraph.GetNetworks()
	networkLiquidities := make([]models.NetworkLiquidity, 0, len(networks))

	for _, networkID := range networks {
		lmAddr, exists := p.liquidityManagers.Get(networkID)
		if !exists {
			return nil, fmt.Errorf("liquidity manager for network %v was not found", networkID)
		}

		lm, err := p.liquidityManagerFactory.NewLiquidityManager(networkID, lmAddr)
		if err != nil {
			return nil, fmt.Errorf("init liquidity manager: %w", err)
		}

		balance, err := lm.GetBalance(ctx)
		if err != nil {
			return nil, fmt.Errorf("get %v balance: %w", networkID, err)
		}

		p.liquidityGraph.SetLiquidity(networkID, balance)
		networkLiquidities = append(networkLiquidities, models.NewNetworkLiquidity(networkID, balance))
	}

	return networkLiquidities, nil
}

func (p *Plugin) loadPendingTransfers(ctx context.Context) ([]models.PendingTransfer, error) {
	// todo: do not load pending transfers all the time
	p.lggr.Infow("loading pending transfers")

	pendingTransfers := make([]models.PendingTransfer, 0)
	for networkID, lmAddress := range p.liquidityManagers.GetAll() {
		lm, err := p.liquidityManagerFactory.NewLiquidityManager(networkID, lmAddress)
		if err != nil {
			return nil, fmt.Errorf("init liquidity manager: %w", err)
		}

		netPendingTransfers, err := lm.GetPendingTransfers(ctx)
		if err != nil {
			return nil, fmt.Errorf("get pending %v transfers: %w", networkID, err)
		}

		pendingTransfers = append(pendingTransfers, netPendingTransfers...)
	}

	p.pendingTransfers.Set(pendingTransfers)
	return pendingTransfers, nil
}

func (p *Plugin) computeMedianLiquidityPerChain(observations []models.Observation) []models.NetworkLiquidity {
	liqObsPerChain := make(map[models.NetworkID][]*big.Int)
	for _, ob := range observations {
		for _, chainLiq := range ob.LiquidityPerChain {
			liqObsPerChain[chainLiq.Network] = append(liqObsPerChain[chainLiq.Network], chainLiq.Liquidity)
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
