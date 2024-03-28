package ocr

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/wasp"

	client2 "github.com/smartcontractkit/chainlink-testing-framework/client"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

const (
	LokiTypeLabel = "ocr_v1_contracts_data"
)

// SimulateAndMonitorOCRAnswers polls all OCR contracts and pushes their answers to Loki every second
func SimulateAndMonitorOCRAnswers(
	logger zerolog.Logger,
	lc *wasp.LokiClient,
	eaChangeInterval time.Duration, // How often to change the EA responses
	ocrInstances []contracts.OffchainAggregator,
	workerNodes []*client.ChainlinkK8sClient,
	msClient *client2.MockserverClient,
	labels map[string]string,
) {
	go MonitorOCRRouds(logger, lc, ocrInstances, labels)
	go SimulateEAActivity(logger, eaChangeInterval, ocrInstances, workerNodes, msClient)
}

// SimulateEAActivity simulates activity on the EA by updating the mockserver responses
func SimulateEAActivity(
	l zerolog.Logger,
	eaChangeInterval time.Duration,
	ocrInstances []contracts.OffchainAggregator,
	workerNodes []*client.ChainlinkK8sClient,
	msClient *client2.MockserverClient,
) {
	l.Info().Msg("Simulating EA activity")
	defer l.Info().Msg("Stopped simulating EA activity")

	for {
		select {
		case <-time.After(eaChangeInterval):
			err := actions.SetAllAdapterResponsesToTheSameValue(rand.Intn(1000), ocrInstances, workerNodes, msClient)
			if err != nil {
				l.Error().Err(err).Msg("failed to update mockserver responses")
			}
		}

	}
}

// MonitorOCRRouds polls OCR contracts and pushes their round data to Loki every second
func MonitorOCRRouds(
	l zerolog.Logger,
	lc *wasp.LokiClient,
	ocrInstances []contracts.OffchainAggregator,
	labels map[string]string,
) {
	l.Info().Msg("Monitoring OCR rounds")
	defer l.Info().Msg("Stopped monitoring OCR rounds")
	for {
		select {
		case <-time.After(1 * time.Second):
			for _, ocr := range ocrInstances {
				metrics := GetOCRTestMetrics(ocr)
				SendMetricsToLoki(metrics, lc, labels)
			}
		}
	}
}

// UpdateLabels updates Loki labels for a log,
func UpdateLabels(t *testing.T, labels map[string]string) map[string]string {
	updatedLabels := make(map[string]string)
	for k, v := range labels {
		updatedLabels[k] = v
	}
	updatedLabels["type"] = LokiTypeLabel
	updatedLabels["go_test_name"] = t.Name()
	updatedLabels["gen_name"] = "performance"
	return updatedLabels
}

// SendMetricsToLoki sends new OCR round data to Loki
func SendMetricsToLoki(metrics *contracts.RoundData, lc *wasp.LokiClient, updatedLabels map[string]string) {
	if err := lc.HandleStruct(wasp.LabelsMapToModel(updatedLabels), time.Now(), metrics); err != nil {
		log.Error().Err(err).Msg("Error pushing OCR metrics to Loki")
	}
}

// GetOCRTestMetrics returns the latest OCR round data
func GetOCRTestMetrics(ocr contracts.OffchainAggregator) *contracts.RoundData {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	roundData, err := ocr.GetLatestRound(ctx)
	if err != nil {
		log.Error().Err(err).Str("Address", ocr.Address()).Msg("Error getting OCR contract round data")
	}
	return roundData
}
