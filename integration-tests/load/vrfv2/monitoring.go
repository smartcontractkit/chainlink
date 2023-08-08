package loadvrfv2

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2_actions"
	"github.com/smartcontractkit/wasp"
	"testing"
	"time"
)

/* Monitors on-chain stats of LoadConsumer and pushes them to Loki every second */

const (
	LokiTypeLabel = "vrfv2_contracts_load_summary"
	ErrMetrics    = "failed to get VRFv2 load test metrics"
	ErrLokiClient = "failed to create Loki client for monitoring"
	ErrLokiPush   = "failed to push monitoring metrics to Loki"
)

func MonitorLoadStats(t *testing.T, vrfv2Contracts *vrfv2_actions.VRFV2Contracts, labels map[string]string) {
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
		for {
			time.Sleep(1 * time.Second)
			metrics, err := vrfv2Contracts.LoadTestConsumer.GetLoadTestMetrics(context.Background())
			if err != nil {
				log.Error().Err(err).Msg(ErrMetrics)
			}
			if err := lc.HandleStruct(wasp.LabelsMapToModel(updatedLabels), time.Now(), metrics); err != nil {
				log.Error().Err(err).Msg(ErrLokiPush)
			}
		}
	}()
}
