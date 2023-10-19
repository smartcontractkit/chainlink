package loadfunctions

import (
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/wasp"
	"testing"
	"time"
)

/* Monitors on-chain stats of LoadConsumer and pushes them to Loki every second */

const (
	LokiTypeLabel = "functions_contracts_load_summary"
	ErrMetrics    = "failed to get Functions load test metrics"
	ErrLokiClient = "failed to create Loki client for monitoring"
	ErrLokiPush   = "failed to push monitoring metrics to Loki"
)

type LoadStats struct {
	Succeeded uint32
	Errored   uint32
	Empty     uint32
}

func MonitorLoadStats(t *testing.T, ft *FunctionsTest, labels map[string]string) {
	go func() {
		updatedLabels := make(map[string]string)
		for k, v := range labels {
			updatedLabels[k] = v
		}
		updatedLabels["type"] = LokiTypeLabel
		updatedLabels["go_test_name"] = t.Name()
		updatedLabels["gen_name"] = "performance"
		lc, err := wasp.NewLokiClient(wasp.NewEnvLokiConfig())
		if err != nil {
			log.Error().Err(err).Msg(ErrLokiClient)
			return
		}
		if err := ft.LoadTestClient.ResetStats(); err != nil {
			log.Error().Err(err).Msg("failed to reset load test client stats")
		}
		for {
			time.Sleep(5 * time.Second)
			stats, err := ft.LoadTestClient.GetStats()
			if err != nil {
				log.Error().Err(err).Msg(ErrMetrics)
			}
			log.Info().
				Hex("LastReqID", []byte(stats.LastRequestID)).
				Str("LastResponse", stats.LastResponse).
				Str("LastError", stats.LastError).
				Uint32("Total", stats.Total).
				Uint32("Succeeded", stats.Succeeded).
				Uint32("Errored", stats.Errored).
				Uint32("Empty", stats.Empty).Msg("On-chain stats for load test client")
			if err := lc.HandleStruct(wasp.LabelsMapToModel(updatedLabels), time.Now(), stats); err != nil {
				log.Error().Err(err).Msg(ErrLokiPush)
			}
		}
	}()
}
