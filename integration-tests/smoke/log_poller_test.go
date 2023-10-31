package smoke

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	logpoller "github.com/smartcontractkit/chainlink/integration-tests/universal/log_poller"
)

// consistency test with no network disruptions with approximate emission of 1500-1600 logs per second for ~110-120 seconds
// 6 filters are registered
func TestLogPoller(t *testing.T) {
	t.Skip()
	cfg := logpoller.Config{
		General: &logpoller.General{
			Generator:   logpoller.GeneratorType_Looped,
			Contracts:   2,
			EventsPerTx: 300,
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
func TestLogManyFiltersPoller(t *testing.T) {
	// t.Skip("test me first")
	cfg := logpoller.Config{
		General: &logpoller.General{
			Generator:   logpoller.GeneratorType_Looped,
			Contracts:   300,
			EventsPerTx: 3,
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

// FIX VALUES
// consistency test that introduces random distruptions by pausing either Chainlink or Postgres containers for random interval of 5-20 seconds
// with approximate emission of 1500-1600 logs per second for ~110-120 seconds
// 6 filters are registered
func TestLogPollerWithChaos(t *testing.T) {
	t.Skip()
	cfg := logpoller.Config{
		General: &logpoller.General{
			Generator:   logpoller.GeneratorType_Looped,
			Contracts:   2,
			EventsPerTx: 100,
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

// consistency test that waits for log poller to process last block after event emission, then it registers filters
// and then waits for backup poller process to kick in and process the past logs
func TestLogPollerBackup(t *testing.T) {
	t.Skip()
	cfg := logpoller.Config{
		General: &logpoller.General{
			Generator:   logpoller.GeneratorType_Looped,
			Contracts:   2,
			EventsPerTx: 200,
		},
		LoopedConfig: &logpoller.LoopedConfig{
			ContractConfig: logpoller.ContractConfig{
				ExecutionCount: 14,
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
	backupPollerTimeout := "5m"

	logpoller.ExecuteBackupLogPollerTest(t, &cfg, backupPollerTimeout)
}

// consistency test that registers filters after events were emitted and then triggers replay via API
// unfortunately there is no way to make sure that logs that are indexed are only picked up by replay
// and not by backup poller
func TestLogPollerReplay(t *testing.T) {
	t.Skip()
	cfg := logpoller.Config{
		General: &logpoller.General{
			Generator:   logpoller.GeneratorType_Looped,
			Contracts:   2,
			EventsPerTx: 4,
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

	logpoller.ExecuteBackupLogPollerReplay(t, &cfg, consistencyTimeout)
}
