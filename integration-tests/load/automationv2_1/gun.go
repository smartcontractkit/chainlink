package automationv2_1

import (
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/wasp"

	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

type LogTriggerGun struct {
	triggerContract               contracts.LogEmitter
	logger                        zerolog.Logger
	numberOfEvents                int
	numberOfSpamMatchingEvents    int
	numberOfSpamNonMatchingEvents int
}

func NewLogTriggerUser(
	triggerContract contracts.LogEmitter,
	logger zerolog.Logger,
	numberOfEvents int,
	numberOfSpamMatchingEvents int,
	numberOfSpamNonMatchingEvents int,
) *LogTriggerGun {

	return &LogTriggerGun{
		triggerContract:               triggerContract,
		logger:                        logger,
		numberOfEvents:                numberOfEvents,
		numberOfSpamMatchingEvents:    numberOfSpamMatchingEvents,
		numberOfSpamNonMatchingEvents: numberOfSpamNonMatchingEvents,
	}
}

func (m *LogTriggerGun) Call(_ *wasp.Generator) *wasp.Response {
	m.logger.Debug().Str("Trigger address", m.triggerContract.Address().String()).Msg("Triggering upkeep")

	if m.numberOfEvents > 0 {
		_, err := m.triggerContract.EmitLogIntMultiIndexed(1, 1, m.numberOfEvents)
		if err != nil {
			return &wasp.Response{Error: err.Error(), Failed: true}
		}
	}

	if m.numberOfSpamMatchingEvents > 0 {
		_, err := m.triggerContract.EmitLogIntMultiIndexed(1, 2, m.numberOfSpamMatchingEvents)
		if err != nil {
			return &wasp.Response{Error: err.Error(), Failed: true}
		}
	}

	if m.numberOfSpamNonMatchingEvents > 0 {
		_, err := m.triggerContract.EmitLogIntMultiIndexed(2, 2, m.numberOfSpamNonMatchingEvents)
		if err != nil {
			return &wasp.Response{Error: err.Error(), Failed: true}
		}
	}

	return &wasp.Response{}
}
