package automationv2_1

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/wasp"
	"time"
)

type LogTriggerGun struct {
	triggerContract *contracts.LogEmitter
	upkeepContract  *contracts.KeeperConsumer
	logger          zerolog.Logger
}

func NewLogTriggerUser(
	triggerContract *contracts.LogEmitter,
	upkeepContract *contracts.KeeperConsumer,
	logger zerolog.Logger,
) *LogTriggerGun {
	return &LogTriggerGun{
		triggerContract: triggerContract,
		upkeepContract:  upkeepContract,
		logger:          logger,
	}
}

func (m *LogTriggerGun) Call(l *wasp.Generator) *wasp.CallResult {
	logTrigger := *m.triggerContract
	upkeepCounter := *m.upkeepContract
	address := logTrigger.Address()
	m.logger.Debug().Str("Trigger address", address.String()).Msg("Triggering upkeep")
	initialCount, err := upkeepCounter.Counter(context.Background())
	m.logger.Debug().Int64("Initial count", initialCount.Int64()).Msg("Initial count")
	startTime := time.Now()
	if err != nil {
		return &wasp.CallResult{Error: err.Error(), Failed: true}
	}
	_, err = logTrigger.EmitLogInt(1)
	if err != nil {
		return &wasp.CallResult{Error: err.Error(), Failed: true}
	}
	for {
		count, err := upkeepCounter.Counter(context.Background())
		//m.logger.Debug().Int64("Count", count.Int64()).Msg("Count")
		if err != nil {
			return &wasp.CallResult{Error: err.Error(), Failed: true}
		}
		if count.Int64() >= initialCount.Int64()+1 {
			endTime := time.Now()
			duration := int(endTime.Sub(startTime).Seconds())
			m.logger.Info().Int("Duration", duration).Msg("Duration")
			break
		}
	}

	return &wasp.CallResult{}
}
