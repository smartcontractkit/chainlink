package automationv2_1

import (
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/wasp"

	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

type LogTriggerGun struct {
	triggerContract contracts.LogEmitter
	upkeepContract  contracts.KeeperConsumer
	logger          zerolog.Logger
	numberOfEvents  int
}

func NewLogTriggerUser(
	triggerContract contracts.LogEmitter,
	upkeepContract contracts.KeeperConsumer,
	logger zerolog.Logger,
	numberOfEvents int,
) *LogTriggerGun {
	return &LogTriggerGun{
		triggerContract: triggerContract,
		upkeepContract:  upkeepContract,
		logger:          logger,
		numberOfEvents:  numberOfEvents,
	}
}

func (m *LogTriggerGun) Call(l *wasp.Generator) *wasp.CallResult {
	m.logger.Debug().Str("Trigger address", m.triggerContract.Address().String()).Msg("Triggering upkeep")
	payload := make([]int, 0)
	for i := 0; i < m.numberOfEvents; i++ {
		payload = append(payload, 1)
	}
	_, err := m.triggerContract.EmitLogInts(payload)
	if err != nil {
		return &wasp.CallResult{Error: err.Error(), Failed: true}
	}

	return &wasp.CallResult{}
}
