package automationv2_1

import (
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/wasp"
)

type LogTriggerGun struct {
	triggerContract contracts.LogEmitter
	upkeepContract  contracts.KeeperConsumer
	logger          zerolog.Logger
}

func NewLogTriggerUser(
	triggerContract contracts.LogEmitter,
	upkeepContract contracts.KeeperConsumer,
	logger zerolog.Logger,
) *LogTriggerGun {
	return &LogTriggerGun{
		triggerContract: triggerContract,
		upkeepContract:  upkeepContract,
		logger:          logger,
	}
}

func (m *LogTriggerGun) Call(l *wasp.Generator) *wasp.CallResult {
	m.logger.Debug().Str("Trigger address", m.triggerContract.Address().String()).Msg("Triggering upkeep")
	//initialCount, err := m.upkeepContract.Counter(context.Background())
	//m.logger.Debug().Int64("Initial count", initialCount.Int64()).Msg("Initial count")
	//if err != nil {
	//	return &wasp.CallResult{Error: err.Error(), Failed: true}
	//}
	_, err := m.triggerContract.EmitLogInt(1)
	if err != nil {
		return &wasp.CallResult{Error: err.Error(), Failed: true}
	}

	return &wasp.CallResult{}
}
