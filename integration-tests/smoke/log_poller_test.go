package smoke

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	logpoller "github.com/smartcontractkit/chainlink/integration-tests/universal/log_poller"
)

func TestLogPoller(t *testing.T) {
	cfg := logpoller.Config{
		General: &logpoller.General{
			Generator:   logpoller.GeneratorType_Looped,
			Contracts:   2,
			EventsPerTx: 5,
		},
		LoopedConfig: &logpoller.LoopedConfig{
			ContractConfig: logpoller.ContractConfig{
				ExecutionCount: 100,
			},
			FuzzConfig: logpoller.FuzzConfig{
				MinEmitWaitTimeMs: 200,
				MaxEmitWaitTimeMs: 500,
			},
		},
	}

	eventsToEmit := []abi.Event{}
	for _, event := range logpoller.EmitterABI.Events {
		eventsToEmit = append(eventsToEmit, event)
	}

	cfg.General.EventsToEmit = eventsToEmit

	logpoller.ExecuteBasicLogPollerTest(t, &cfg)
}
