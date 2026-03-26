package evm

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
)

// Prometheus gauges for the node's view of on-chain OCR2 median aggregator state (DF-22676).
// latest_answer is exported as float64; values above ~2^53 may lose integer precision (documented in HELP).
var (
	promOCR2OnchainLatestAnswer = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ocr2_onchain_transmission_latest_answer",
		Help: "Latest median answer from the node's on-chain read of LatestTransmissionDetails (OCR2 aggregator). Float64 may not represent full int256 precision.",
	}, []string{"chain_id", "contract_address", "transmitter_id"})
	promOCR2OnchainLatestTimestampUnix = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ocr2_onchain_transmission_latest_timestamp_unix",
		Help: "Unix seconds of latest on-chain transmission timestamp from the node's LatestTransmissionDetails read.",
	}, []string{"chain_id", "contract_address", "transmitter_id"})
	promOCR2OnchainEpoch = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ocr2_onchain_transmission_epoch",
		Help: "Epoch from the node's on-chain LatestTransmissionDetails read.",
	}, []string{"chain_id", "contract_address", "transmitter_id"})
	promOCR2OnchainRound = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ocr2_onchain_transmission_round",
		Help: "Round from the node's on-chain LatestTransmissionDetails read.",
	}, []string{"chain_id", "contract_address", "transmitter_id"})
)

var (
	medianOnchainOtelOnce sync.Once
	medianOtelLatestAnswer        metric.Float64Gauge
	medianOtelLatestTimestampUnix metric.Float64Gauge
	medianOtelEpoch               metric.Int64Gauge
	medianOtelRound               metric.Int64Gauge
)

func initMedianOnchainOtel() {
	medianOnchainOtelOnce.Do(func() {
		meter := beholder.GetMeter()
		la, err := meter.Float64Gauge("ocr2_onchain_transmission_latest_answer")
		if err != nil {
			return
		}
		ltu, err := meter.Float64Gauge("ocr2_onchain_transmission_latest_timestamp_unix")
		if err != nil {
			return
		}
		e, err := meter.Int64Gauge("ocr2_onchain_transmission_epoch")
		if err != nil {
			return
		}
		r, err := meter.Int64Gauge("ocr2_onchain_transmission_round")
		if err != nil {
			return
		}
		medianOtelLatestAnswer = la
		medianOtelLatestTimestampUnix = ltu
		medianOtelEpoch = e
		medianOtelRound = r
	})
}

// medianOnchainViewMetrics records the same on-chain OCR2 view series to Prometheus and OpenTelemetry (Beholder), matching the txmMetrics pattern (package-level Prom vecs, OTel instruments on the metrics type).
type medianOnchainViewMetrics struct {
	chainID       string
	contract      string
	transmitterID string

	latestAnswer        metric.Float64Gauge
	latestTimestampUnix metric.Float64Gauge
	epoch               metric.Int64Gauge
	round               metric.Int64Gauge
}

func newMedianOnchainViewMetrics(chainID, contractAddress, transmitterID string) *medianOnchainViewMetrics {
	if transmitterID == "" {
		transmitterID = "unknown"
	}
	initMedianOnchainOtel()
	return &medianOnchainViewMetrics{
		chainID:               chainID,
		contract:              contractAddress,
		transmitterID:         transmitterID,
		latestAnswer:          medianOtelLatestAnswer,
		latestTimestampUnix:   medianOtelLatestTimestampUnix,
		epoch:                 medianOtelEpoch,
		round:                 medianOtelRound,
	}
}

func (m *medianOnchainViewMetrics) onchainAttrs() metric.MeasurementOption {
	return metric.WithAttributes(
		attribute.String("chain_id", m.chainID),
		attribute.String("contract_address", m.contract),
		attribute.String("transmitter_id", m.transmitterID),
	)
}

func (m *medianOnchainViewMetrics) record(ctx context.Context, latestAnswer *big.Int, epoch uint32, round uint8, updatedAt time.Time) {
	lv := bigIntToPromGaugeValue(latestAnswer)
	promOCR2OnchainLatestAnswer.WithLabelValues(m.chainID, m.contract, m.transmitterID).Set(lv)
	promOCR2OnchainLatestTimestampUnix.WithLabelValues(m.chainID, m.contract, m.transmitterID).Set(float64(updatedAt.Unix()))
	promOCR2OnchainEpoch.WithLabelValues(m.chainID, m.contract, m.transmitterID).Set(float64(epoch))
	promOCR2OnchainRound.WithLabelValues(m.chainID, m.contract, m.transmitterID).Set(float64(round))

	if m.latestAnswer == nil {
		return
	}
	opts := m.onchainAttrs()
	m.latestAnswer.Record(ctx, lv, opts)
	m.latestTimestampUnix.Record(ctx, float64(updatedAt.Unix()), opts)
	m.epoch.Record(ctx, int64(epoch), opts)
	m.round.Record(ctx, int64(round), opts)
}

func bigIntToPromGaugeValue(i *big.Int) float64 {
	if i == nil {
		return 0
	}
	f, _ := new(big.Float).SetInt(i).Float64()
	return f
}
