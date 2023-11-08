package reorg

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"

	logpoller "github.com/smartcontractkit/chainlink/integration-tests/universal/log_poller"
)

func TestLogPollerFromEnv(t *testing.T) {
	cfg := logpoller.Config{
		General: &logpoller.General{
			Generator:      logpoller.GeneratorType_Looped,
			Contracts:      2,
			EventsPerTx:    100,
			UseFinalityTag: true,
		},
		LoopedConfig: &logpoller.LoopedConfig{
			ContractConfig: logpoller.ContractConfig{
				ExecutionCount: 100,
			},
			FuzzConfig: logpoller.FuzzConfig{
				MinEmitWaitTimeMs: 400,
				MaxEmitWaitTimeMs: 600,
			},
		},
	}

	eventsToEmit := []abi.Event{}
	for _, event := range logpoller.EmitterABI.Events {
		eventsToEmit = append(eventsToEmit, event)
	}

	cfg.General.EventsToEmit = eventsToEmit
	err := cfg.OverrideFromEnv()
	if err != nil {
		t.Errorf("failed to override config from env: %v", err)
		t.FailNow()
	}

	logpoller.ExecuteCILogPollerTest(t, &cfg)
}
