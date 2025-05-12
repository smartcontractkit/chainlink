package txmgr

import (
	"context"
	"fmt"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-framework/metrics"
)

var (
	promNumSuccessfulTxs = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "tx_manager_num_successful_transactions",
		Help: "Total number of successful transactions. Note that this can err to be too high since transactions are counted on each confirmation, which can happen multiple times per transaction in the case of re-orgs",
	}, []string{"chainID"})
	promRevertedTxCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "tx_manager_num_tx_reverted",
		Help: "Number of times a transaction reverted on-chain. Note that this can err to be too high since transactions are counted on each confirmation, which can happen multiple times per transaction in the case of re-orgs",
	}, []string{"chainID"})
	promFwdTxCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "tx_manager_fwd_tx_count",
		Help: "The number of forwarded transaction attempts labeled by status",
	}, []string{"chainID", "successful"})
	promTxAttemptCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "tx_manager_tx_attempt_count",
		Help: "The number of transaction attempts that are currently being processed by the transaction manager",
	}, []string{"chainID"})
	promNumFinalizedTxs = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "tx_manager_num_finalized_transactions",
		Help: "Total number of finalized transactions",
	}, []string{"chainID"})
)

type evmTxmMetrics struct {
	metrics.GenericTXMMetrics
	chainID          string
	numSuccessfulTxs metric.Int64Counter
	numRevertedTxs   metric.Int64Counter
	fwdTxCount       metric.Int64Counter
	txAttemptCount   metric.Float64Gauge
	numFinalizedTxs  metric.Int64Counter
}

func NewEVMTxmMetrics(chainID string) (*evmTxmMetrics, error) {
	genericTXMMetrics, err := metrics.NewGenericTxmMetrics(chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize generic TXM metrics: %w", err)
	}

	numSuccessfulTxs, err := beholder.GetMeter().Int64Counter("tx_manager_num_successful_transactions")
	if err != nil {
		return nil, fmt.Errorf("failed to register number of successful transactions metric: %w", err)
	}

	numRevertedTxs, err := beholder.GetMeter().Int64Counter("tx_manager_num_tx_reverted")
	if err != nil {
		return nil, fmt.Errorf("failed to register number of reverted transactions metric: %w", err)
	}

	fwdTxCount, err := beholder.GetMeter().Int64Counter("tx_manager_fwd_tx_count")
	if err != nil {
		return nil, fmt.Errorf("failed to register forward transaction count metric: %w", err)
	}

	txAttemptCount, err := beholder.GetMeter().Float64Gauge("tx_manager_tx_attempt_count")
	if err != nil {
		return nil, fmt.Errorf("failed to register transaction attempt count metric: %w", err)
	}

	numFinalizedTxs, err := beholder.GetMeter().Int64Counter("tx_manager_num_finalized_transactions")
	if err != nil {
		return nil, fmt.Errorf("failed to register number of finalized transactions metric: %w", err)
	}

	return &evmTxmMetrics{
		chainID:           chainID,
		GenericTXMMetrics: genericTXMMetrics,
		numSuccessfulTxs:  numSuccessfulTxs,
		numRevertedTxs:    numRevertedTxs,
		fwdTxCount:        fwdTxCount,
		txAttemptCount:    txAttemptCount,
		numFinalizedTxs:   numFinalizedTxs,
	}, nil
}

func (m *evmTxmMetrics) IncrementNumSuccessfulTxs(ctx context.Context) {
	promNumSuccessfulTxs.WithLabelValues(m.chainID).Add(float64(1))
	m.numSuccessfulTxs.Add(ctx, 1, metric.WithAttributes(attribute.String("chainID", m.chainID)))
}

func (m *evmTxmMetrics) IncrementNumRevertedTxs(ctx context.Context) {
	promRevertedTxCount.WithLabelValues(m.chainID).Add(float64(1))
	m.numRevertedTxs.Add(ctx, 1, metric.WithAttributes(attribute.String("chainID", m.chainID)))
}

func (m *evmTxmMetrics) IncrementFwdTxCount(ctx context.Context, successful bool) {
	promFwdTxCount.WithLabelValues(m.chainID, strconv.FormatBool(successful)).Add(float64(1))
	m.fwdTxCount.Add(ctx, 1, metric.WithAttributes(attribute.String("chainID", m.chainID), attribute.Bool("successful", successful)))
}

func (m *evmTxmMetrics) RecordTxAttemptCount(ctx context.Context, value float64) {
	promTxAttemptCount.WithLabelValues(m.chainID).Set(value)
	m.txAttemptCount.Record(ctx, value, metric.WithAttributes(attribute.String("chainID", m.chainID)))
}

func (m *evmTxmMetrics) IncrementNumFinalizedTxs(ctx context.Context) {
	promNumFinalizedTxs.WithLabelValues(m.chainID).Add(float64(1))
	m.numFinalizedTxs.Add(ctx, 1, metric.WithAttributes(attribute.String("chainID", m.chainID)))
}
