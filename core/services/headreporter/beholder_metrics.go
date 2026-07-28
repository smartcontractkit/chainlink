package headreporter

import (
	"context"
	"fmt"
	"strconv"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
)

const (
	metricLatestBlockNumber       = "head_reporter_latest_block_number"
	metricLatestBlockTimestamp    = "head_reporter_latest_block_timestamp_seconds"
	metricFinalizedBlockNumber    = "head_reporter_finalized_block_number"
	metricFinalizedBlockTimestamp = "head_reporter_finalized_block_timestamp_seconds"
	metricFinalityDepth           = "head_reporter_finality_depth_blocks"

	attrChainID       = "chain_id"
	attrNetwork       = "network"
	attrChainSelector = "chain_selector"
)

type (
	// finalizedBlock is the finalized counterpart of a headReport's latest block.
	finalizedBlock struct {
		number int64
		ts     uint64
	}

	// headReport is a family-neutral snapshot of chain head data, produced by both the
	// in-process EVM reporter and the generic relayer reporter.
	headReport struct {
		chainID       string
		network       string // chain family, e.g. "evm", "solana"
		chainSelector uint64
		hasSelector   bool
		latestNumber  int64
		latestTs      uint64
		finalized     *finalizedBlock // nil if not available
	}

	// HeadMetrics records head-report data as Beholder (OTel) metrics.
	HeadMetrics interface {
		RecordHeadReport(ctx context.Context, r headReport)
	}

	beholderHeadMetrics struct {
		latestNumber    metric.Int64Gauge
		latestTs        metric.Int64Gauge
		finalizedNumber metric.Int64Gauge
		finalizedTs     metric.Int64Gauge
		finalityDepth   metric.Int64Gauge
	}
)

// NewBeholderHeadMetrics registers the head-reporter gauge instruments once via
// beholder.GetMeter(). It is safe to call unconditionally: GetMeter() returns a no-op meter
// when Beholder telemetry is disabled.
func NewBeholderHeadMetrics() (HeadMetrics, error) {
	m := beholder.GetMeter()

	latestNumber, err := m.Int64Gauge(metricLatestBlockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to register %s gauge: %w", metricLatestBlockNumber, err)
	}
	latestTs, err := m.Int64Gauge(metricLatestBlockTimestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to register %s gauge: %w", metricLatestBlockTimestamp, err)
	}
	finalizedNumber, err := m.Int64Gauge(metricFinalizedBlockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to register %s gauge: %w", metricFinalizedBlockNumber, err)
	}
	finalizedTs, err := m.Int64Gauge(metricFinalizedBlockTimestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to register %s gauge: %w", metricFinalizedBlockTimestamp, err)
	}
	finalityDepth, err := m.Int64Gauge(metricFinalityDepth)
	if err != nil {
		return nil, fmt.Errorf("failed to register %s gauge: %w", metricFinalityDepth, err)
	}

	return &beholderHeadMetrics{
		latestNumber:    latestNumber,
		latestTs:        latestTs,
		finalizedNumber: finalizedNumber,
		finalizedTs:     finalizedTs,
		finalityDepth:   finalityDepth,
	}, nil
}

func (b *beholderHeadMetrics) RecordHeadReport(ctx context.Context, r headReport) {
	kvs := []attribute.KeyValue{
		attribute.String(attrChainID, r.chainID),
		attribute.String(attrNetwork, r.network),
	}
	if r.hasSelector {
		kvs = append(kvs, attribute.String(attrChainSelector, strconv.FormatUint(r.chainSelector, 10)))
	}
	attrs := metric.WithAttributes(kvs...)

	b.latestNumber.Record(ctx, r.latestNumber, attrs)
	b.latestTs.Record(ctx, int64(r.latestTs), attrs) //nolint:gosec // unix timestamps fit int64 until year 292277026596

	if r.finalized != nil {
		b.finalizedNumber.Record(ctx, r.finalized.number, attrs)
		b.finalizedTs.Record(ctx, int64(r.finalized.ts), attrs) //nolint:gosec // unix timestamps fit int64 until year 292277026596
		b.finalityDepth.Record(ctx, r.latestNumber-r.finalized.number, attrs)
	}
}
