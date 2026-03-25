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

type medianOnchainOtelInstruments struct {
	latestAnswer        metric.Float64Gauge
	latestTimestampUnix metric.Float64Gauge
	epoch               metric.Int64Gauge
	round               metric.Int64Gauge
}

var (
	medianOnchainOtelInst *medianOnchainOtelInstruments
	medianOnchainOtelOnce sync.Once
)

func loadMedianOnchainOtel() *medianOnchainOtelInstruments {
	medianOnchainOtelOnce.Do(func() {
		m := beholder.GetMeter()
		latestAnswer, err := m.Float64Gauge("ocr2_onchain_transmission_latest_answer")
		if err != nil {
			return
		}
		latestTimestampUnix, err := m.Float64Gauge("ocr2_onchain_transmission_latest_timestamp_unix")
		if err != nil {
			return
		}
		epoch, err := m.Int64Gauge("ocr2_onchain_transmission_epoch")
		if err != nil {
			return
		}
		round, err := m.Int64Gauge("ocr2_onchain_transmission_round")
		if err != nil {
			return
		}
		medianOnchainOtelInst = &medianOnchainOtelInstruments{
			latestAnswer:        latestAnswer,
			latestTimestampUnix: latestTimestampUnix,
			epoch:               epoch,
			round:               round,
		}
	})
	return medianOnchainOtelInst
}

type medianOnchainViewMetrics struct {
	chainID       string
	contract      string
	transmitterID string
}

func newMedianOnchainViewMetrics(chainID, contractAddress, transmitterID string) *medianOnchainViewMetrics {
	if transmitterID == "" {
		transmitterID = "unknown"
	}
	return &medianOnchainViewMetrics{
		chainID:       chainID,
		contract:      contractAddress,
		transmitterID: transmitterID,
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

	if ot := loadMedianOnchainOtel(); ot != nil {
		opts := m.onchainAttrs()
		ot.latestAnswer.Record(ctx, lv, opts)
		ot.latestTimestampUnix.Record(ctx, float64(updatedAt.Unix()), opts)
		ot.epoch.Record(ctx, int64(epoch), opts)
		ot.round.Record(ctx, int64(round), opts)
	}
}

func bigIntToPromGaugeValue(i *big.Int) float64 {
	if i == nil {
		return 0
	}
	f, _ := new(big.Float).SetInt(i).Float64()
	return f
}
