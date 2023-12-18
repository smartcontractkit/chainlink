package automationv2_1

import (
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/wasp"

	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

type LogTriggerGun struct {
	triggerContract contracts.LogEmitter
	logger          zerolog.Logger
	triggerPayload  []int
	spamPayload     []int
}

func NewLogTriggerUser(
	triggerContract contracts.LogEmitter,
	logger zerolog.Logger,
	numberOfEvents int,
	numberOfSpamEvents int,
) *LogTriggerGun {

	triggerPayload := make([]int, 0)
	spamPayload := make([]int, 0)

	if numberOfEvents > 0 {
		for i := 0; i < numberOfEvents; i++ {
			triggerPayload = append(triggerPayload, 1)
		}
	}

	if numberOfSpamEvents > 0 {
		for i := 0; i < numberOfSpamEvents; i++ {
			spamPayload = append(spamPayload, 0)
		}
	}

	return &LogTriggerGun{
		triggerContract: triggerContract,
		logger:          logger,
		triggerPayload:  triggerPayload,
		spamPayload:     spamPayload,
	}
}

func (m *LogTriggerGun) Call(_ *wasp.Generator) *wasp.Response {
	m.logger.Debug().Str("Trigger address", m.triggerContract.Address().String()).Msg("Triggering upkeep")

	if len(m.triggerPayload) > 0 {
		_, err := m.triggerContract.EmitLogIntsIndexed(m.triggerPayload)
		if err != nil {
			return &wasp.Response{Error: err.Error(), Failed: true}
		}
	}

	if len(m.spamPayload) > 0 {
		_, err := m.triggerContract.EmitLogIntsIndexed(m.spamPayload)
		if err != nil {
			return &wasp.Response{Error: err.Error(), Failed: true}
		}
	}

	return &wasp.Response{}
}
