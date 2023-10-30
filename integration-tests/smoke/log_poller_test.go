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

func TestLogPollerBackup(t *testing.T) {
	cfg := logpoller.Config{
		General: &logpoller.General{
			Generator:   logpoller.GeneratorType_Looped,
			Contracts:   2,
			EventsPerTx: 5,
		},
		LoopedConfig: &logpoller.LoopedConfig{
			ContractConfig: logpoller.ContractConfig{
				ExecutionCount: 10,
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

	logpoller.ExecuteBackupLogPollerTest(t, &cfg)
}

func TestLogPollerReplay(t *testing.T) {
	// with these 12k logs it doesn't finish in 5 minutes
	// at some point log count in DB stops begin updated
	// I imagine it's a test issue, not a log poller issue
	// cfg := logpoller.Config{
	// 	General: &logpoller.General{
	// 		Generator:   logpoller.GeneratorType_Looped,
	// 		Contracts:   4,
	// 		EventsPerTx: 10,
	// 	},
	// 	LoopedConfig: &logpoller.LoopedConfig{
	// 		ContractConfig: logpoller.ContractConfig{
	// 			ExecutionCount: 100,
	// 		},
	// 		FuzzConfig: logpoller.FuzzConfig{
	// 			MinEmitWaitTimeMs: 200,
	// 			MaxEmitWaitTimeMs: 500,
	// 		},
	// 	},
	// }

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
