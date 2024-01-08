package smoke

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"

	logpoller "github.com/smartcontractkit/chainlink/integration-tests/universal/log_poller"
)

// consistency test with no network disruptions with approximate emission of 1500-1600 logs per second for ~110-120 seconds
// 6 filters are registered
func TestLogPollerFewFiltersFixedDepth(t *testing.T) {
	cfg := logpoller.Config{
		General: &logpoller.General{
			Generator:      logpoller.GeneratorType_Looped,
			Contracts:      2,
			EventsPerTx:    4,
			UseFinalityTag: false,
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

func TestLogPollerFewFiltersFinalityTag(t *testing.T) {
	cfg := logpoller.Config{
		General: &logpoller.General{
			Generator:      logpoller.GeneratorType_Looped,
			Contracts:      2,
			EventsPerTx:    4,
			UseFinalityTag: true,
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

// consistency test with no network disruptions with approximate emission of 1000-1100 logs per second for ~110-120 seconds
// 900 filters are registered
func TestLogManyFiltersPollerFixedDepth(t *testing.T) {
	cfg := logpoller.Config{
		General: &logpoller.General{
			Generator:      logpoller.GeneratorType_Looped,
			Contracts:      300,
			EventsPerTx:    3,
			UseFinalityTag: false,
		},
		LoopedConfig: &logpoller.LoopedConfig{
			ContractConfig: logpoller.ContractConfig{
				ExecutionCount: 30,
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

func TestLogManyFiltersPollerFinalityTag(t *testing.T) {
	cfg := logpoller.Config{
		General: &logpoller.General{
			Generator:      logpoller.GeneratorType_Looped,
			Contracts:      300,
			EventsPerTx:    3,
			UseFinalityTag: true,
		},
		LoopedConfig: &logpoller.LoopedConfig{
			ContractConfig: logpoller.ContractConfig{
				ExecutionCount: 30,
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

// consistency test that introduces random distruptions by pausing either Chainlink or Postgres containers for random interval of 5-20 seconds
// with approximate emission of 520-550 logs per second for ~110 seconds
// 6 filters are registered
func TestLogPollerWithChaosFixedDepth(t *testing.T) {
	cfg := logpoller.Config{
		General: &logpoller.General{
			Generator:      logpoller.GeneratorType_Looped,
			Contracts:      2,
			EventsPerTx:    100,
			UseFinalityTag: false,
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
		ChaosConfig: &logpoller.ChaosConfig{
			ExperimentCount: 10,
		},
	}

	eventsToEmit := []abi.Event{}
	for _, event := range logpoller.EmitterABI.Events {
		eventsToEmit = append(eventsToEmit, event)
	}

	cfg.General.EventsToEmit = eventsToEmit

	logpoller.ExecuteBasicLogPollerTest(t, &cfg)
}

func TestLogPollerWithChaosFinalityTag(t *testing.T) {
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
				MinEmitWaitTimeMs: 200,
				MaxEmitWaitTimeMs: 500,
			},
		},
		ChaosConfig: &logpoller.ChaosConfig{
			ExperimentCount: 10,
		},
	}

	eventsToEmit := []abi.Event{}
	for _, event := range logpoller.EmitterABI.Events {
		eventsToEmit = append(eventsToEmit, event)
	}

	cfg.General.EventsToEmit = eventsToEmit

	logpoller.ExecuteBasicLogPollerTest(t, &cfg)
}

// consistency test that registers filters after events were emitted and then triggers replay via API
// unfortunately there is no way to make sure that logs that are indexed are only picked up by replay
// and not by backup poller
// with approximate emission of 24 logs per second for ~110 seconds
// 6 filters are registered
func TestLogPollerReplayFixedDepth(t *testing.T) {
	cfg := logpoller.Config{
		General: &logpoller.General{
			Generator:      logpoller.GeneratorType_Looped,
			Contracts:      2,
			EventsPerTx:    4,
			UseFinalityTag: false,
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
	consistencyTimeout := "5m"

	logpoller.ExecuteLogPollerReplay(t, &cfg, consistencyTimeout)
}

func TestLogPollerReplayFinalityTag(t *testing.T) {
	cfg := logpoller.Config{
		General: &logpoller.General{
			Generator:      logpoller.GeneratorType_Looped,
			Contracts:      2,
			EventsPerTx:    4,
			UseFinalityTag: false,
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
	consistencyTimeout := "5m"

	logpoller.ExecuteLogPollerReplay(t, &cfg, consistencyTimeout)
}
