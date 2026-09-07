package headreporter

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	chainselector "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// relayerReportTimeout bounds each relay's LatestHead/FinalizedHead calls in
// relayerMetricsReporter.ReportPeriodic, so one unresponsive relay can't stall metrics
// collection for the rest (or delay the next periodic tick).
const relayerReportTimeout = 2 * time.Second

// chainSelector is the resolved chain_selector label for a chain, or the zero value with
// ok=false if it could not be resolved (e.g. chain not in the chain-selectors registry).
type chainSelector struct {
	value uint64
	ok    bool
}

// resolveChainSelector looks up the chain_selector label for (chainID, network) once. On
// failure it logs a warning and reports ok=false so callers omit the label rather than fail.
func resolveChainSelector(lggr logger.Logger, chainID, network string) chainSelector {
	details, err := chainselector.GetChainDetailsByChainIDAndFamily(chainID, strings.ToLower(network))
	if err != nil {
		lggr.Warnw("Failed to resolve chain selector for head metrics; omitting chain_selector label",
			"chainID", chainID, "network", network, "err", err)
		return chainSelector{}
	}
	return chainSelector{value: details.ChainSelector, ok: true}
}

// evmMetricsReporter records head/finality data for in-process EVM chains as Beholder metrics,
// driven by HeadTracker's new-head events.
type evmMetricsReporter struct {
	lggr      logger.Logger
	metrics   HeadMetrics
	selectors map[uint64]chainSelector // keyed by EVM chain ID
}

func NewEVMMetricsReporter(metrics HeadMetrics, lggr logger.Logger, chainIDs ...*big.Int) HeadReporter {
	lggr = lggr.Named("MetricsReporter")
	selectors := make(map[uint64]chainSelector, len(chainIDs))
	for _, chainID := range chainIDs {
		selectors[chainID.Uint64()] = resolveChainSelector(lggr, chainID.String(), chainselector.FamilyEVM)
	}
	return &evmMetricsReporter{lggr: lggr, metrics: metrics, selectors: selectors}
}

func (r *evmMetricsReporter) ReportNewHead(ctx context.Context, head *evmtypes.Head) error {
	sel := r.selectors[head.EVMChainID.ToInt().Uint64()]

	var finalized *finalizedBlock
	if latestFinalizedHead := head.LatestFinalizedHead(); latestFinalizedHead != nil {
		finalized = &finalizedBlock{
			number: latestFinalizedHead.BlockNumber(),
			ts:     utils.NonNegativeInt64ToUint64(latestFinalizedHead.GetTimestamp().UTC().Unix()),
		}
	}

	r.metrics.RecordHeadReport(ctx, headReport{
		chainID:       head.EVMChainID.String(),
		network:       chainselector.FamilyEVM,
		chainSelector: sel.value,
		hasSelector:   sel.ok,
		latestNumber:  head.Number,
		latestTs:      utils.NonNegativeInt64ToUint64(head.Timestamp.UTC().Unix()),
		finalized:     finalized,
	})
	return nil
}

// ReportPeriodic is unimplemented because head metrics for EVM are recorded on each new head.
func (r *evmMetricsReporter) ReportPeriodic(_ context.Context) error {
	return nil
}

// relayerMetricsReporter records head/finality data for generic (LOOP) relayers as Beholder
// metrics, polled on the periodic tick since there is no HeadTracker to subscribe to.
type relayerMetricsReporter struct {
	lggr      logger.Logger
	metrics   HeadMetrics
	relays    map[types.RelayID]loop.Relayer
	selectors map[types.RelayID]chainSelector
}

func NewRelayerMetricsReporter(metrics HeadMetrics, lggr logger.Logger, relayers map[types.RelayID]loop.Relayer) HeadReporter {
	if relayers == nil {
		return nil
	}
	lggr = lggr.Named("MetricsReporter")
	selectors := make(map[types.RelayID]chainSelector, len(relayers))
	for relayID := range relayers {
		selectors[relayID] = resolveChainSelector(lggr, relayID.ChainID, relayID.Network)
	}
	return &relayerMetricsReporter{lggr: lggr, metrics: metrics, relays: relayers, selectors: selectors}
}

// ReportNewHead is unimplemented because there is no HeadTracker to subscribe to.
func (r *relayerMetricsReporter) ReportNewHead(_ context.Context, _ *evmtypes.Head) error {
	return nil
}

func (r *relayerMetricsReporter) ReportPeriodic(ctx context.Context) error {
	for relayID, relay := range r.relays {
		if err := r.reportOne(ctx, relayID, relay); err != nil {
			return err
		}
	}
	return nil
}

func (r *relayerMetricsReporter) reportOne(ctx context.Context, relayID types.RelayID, relay loop.Relayer) error {
	ctx, cancel := context.WithTimeout(ctx, relayerReportTimeout)
	defer cancel()

	head, err := relay.LatestHead(ctx)
	if err != nil {
		return fmt.Errorf("failed getting head for chain %s: %w", relayID, err)
	}

	blockNum, err := parseGenericHeadBlockNumber(relayID, head)
	if err != nil {
		return err
	}

	var finalized *finalizedBlock
	finalizedHead, err := fetchFinalizedGenericHead(ctx, relayID, relay)
	if err != nil {
		r.lggr.Warnw("Failed to fetch finalized head", "chainID", relayID.ChainID, "network", relayID.Network, "err", err)
	} else if finalizedHead != nil {
		finalizedNum, parseErr := parseGenericHeadBlockNumber(relayID, *finalizedHead)
		if parseErr != nil {
			r.lggr.Warnw("Failed to parse finalized head block number", "chainID", relayID.ChainID, "network", relayID.Network, "err", parseErr)
		} else {
			finalized = &finalizedBlock{
				number: int64(finalizedNum), //nolint:gosec // block numbers fit int64
				ts:     finalizedHead.Timestamp,
			}
		}
	}

	sel := r.selectors[relayID]
	r.metrics.RecordHeadReport(ctx, headReport{
		chainID:       relayID.ChainID,
		network:       relayID.Network,
		chainSelector: sel.value,
		hasSelector:   sel.ok,
		latestNumber:  int64(blockNum), //nolint:gosec // block numbers fit int64
		latestTs:      head.Timestamp,
		finalized:     finalized,
	})
	return nil
}
