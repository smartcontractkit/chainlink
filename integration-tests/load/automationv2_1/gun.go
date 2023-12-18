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
	spamNPayload    []int
}

func NewLogTriggerUser(
	triggerContract contracts.LogEmitter,
	logger zerolog.Logger,
	numberOfEvents int,
	numberOfSpamEvents int,
	numberOfSpamNonMatchingEvents int,
) *LogTriggerGun {

	triggerPayload := make([]int, 0)
	spamPayload := make([]int, 0)
	spamNPayload := make([]int, 0)

	if numberOfEvents > 0 {
		for i := 0; i < numberOfEvents; i++ {
			triggerPayload = append(triggerPayload, 1)
		}
	}

	if numberOfSpamEvents > 0 {
		for i := 0; i < numberOfSpamEvents; i++ {
			spamPayload = append(spamPayload, 1)
		}
	}

	if numberOfSpamNonMatchingEvents > 0 {
		for i := 0; i < numberOfSpamNonMatchingEvents; i++ {
			spamNPayload = append(spamNPayload, 2)
		}
	}

	return &LogTriggerGun{
		triggerContract: triggerContract,
		logger:          logger,
		triggerPayload:  triggerPayload,
		spamPayload:     spamPayload,
		spamNPayload:    spamNPayload,
	}
}

func (m *LogTriggerGun) Call(_ *wasp.Generator) *wasp.Response {
	m.logger.Debug().Str("Trigger address", m.triggerContract.Address().String()).Msg("Triggering upkeep")

	if len(m.triggerPayload) > 0 {
		_, err := m.triggerContract.EmitLogIntsMultiIndexed(m.triggerPayload, 1)
		if err != nil {
			return &wasp.Response{Error: err.Error(), Failed: true}
		}
	}

	if len(m.spamPayload) > 0 {
		_, err := m.triggerContract.EmitLogIntsMultiIndexed(m.spamPayload, 2)
		if err != nil {
			return &wasp.Response{Error: err.Error(), Failed: true}
		}
	}

	if len(m.spamNPayload) > 0 {
		_, err := m.triggerContract.EmitLogIntsMultiIndexed(m.spamNPayload, 2)
		if err != nil {
			return &wasp.Response{Error: err.Error(), Failed: true}
		}
	}

	return &wasp.Response{}
}
