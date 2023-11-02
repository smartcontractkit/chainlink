package reorg

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	logpoller "github.com/smartcontractkit/chainlink/integration-tests/universal/log_poller"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

func TestLogPollerFinalityTag(t *testing.T) {
	cfg := logpoller.Config{
		General: &logpoller.General{
			Generator:      logpoller.GeneratorType_WASP,
			Contracts:      2,
			EventsPerTx:    40,
			UseFinalityTag: true,
		},
		Wasp: &logpoller.WaspConfig{
			Load: &logpoller.Load{
				LPS:                   200,
				RateLimitUnitDuration: models.MustNewDuration(2 * time.Second),
				Duration:              models.MustNewDuration(20 * time.Minute),
				CallTimeout:           models.MustNewDuration(3 * time.Minute),
			},
		},
	}

	eventsToEmit := []abi.Event{}
	for _, event := range logpoller.EmitterABI.Events {
		eventsToEmit = append(eventsToEmit, event)
	}

	cfg.General.EventsToEmit = eventsToEmit

	logpoller.ExecuteCILogPollerTest(t, &cfg, func(_ int64, endBlock int64) (int64, error) {
		return endBlock + 10, nil
	})
}

func TestLogPollerFixedFinaltyDepth(t *testing.T) {
	cfg := logpoller.Config{
		General: &logpoller.General{
			Generator:      logpoller.GeneratorType_Looped,
			Contracts:      2,
			EventsPerTx:    10,
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

	logpoller.ExecuteCILogPollerTest(t, &cfg, func(chainId int64, endBlock int64) (int64, error) {
		finalityDepth, err := logpoller.GetFinalityDepth(chainId)
		if err != nil {
			return 0, err
		}

		return endBlock + finalityDepth, nil
	})
}
