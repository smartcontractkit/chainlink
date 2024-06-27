package loadvrfv2plus

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2plus"
	"github.com/smartcontractkit/wasp"
	"testing"
	"time"
)

/* Monitors on-chain stats of LoadConsumer and pushes them to Loki every second */

const (
	LokiTypeLabel = "vrfv2plus_contracts_load_summary"
	ErrMetrics    = "failed to get VRFv2Plus load test metrics"
	ErrLokiClient = "failed to create Loki client for monitoring"
	ErrLokiPush   = "failed to push monitoring metrics to Loki"
)

func MonitorLoadStats(lc *wasp.LokiClient, vrfv2PlusContracts *vrfv2plus.VRFV2_5Contracts, labels map[string]string) {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			SendLoadTestMetricsToLoki(vrfv2PlusContracts, lc, labels)
		}
	}()
}

func UpdateLabels(labels map[string]string, t *testing.T) map[string]string {
	updatedLabels := make(map[string]string)
	for k, v := range labels {
		updatedLabels[k] = v
	}
	updatedLabels["type"] = LokiTypeLabel
	updatedLabels["go_test_name"] = t.Name()
	updatedLabels["gen_name"] = "performance"
	return updatedLabels
}

func SendLoadTestMetricsToLoki(vrfv2PlusContracts *vrfv2plus.VRFV2_5Contracts, lc *wasp.LokiClient, updatedLabels map[string]string) {
	//todo - should work with multiple consumers and consumers having different keyhashes and wallets
	metrics, err := vrfv2PlusContracts.LoadTestConsumers[0].GetLoadTestMetrics(context.Background())
	if err != nil {
		log.Error().Err(err).Msg(ErrMetrics)
	}
	if err := lc.HandleStruct(wasp.LabelsMapToModel(updatedLabels), time.Now(), metrics); err != nil {
		log.Error().Err(err).Msg(ErrLokiPush)
	}
}
