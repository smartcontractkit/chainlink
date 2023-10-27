package logpoller

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/wasp"
)

/* LogEmitterGun is a gun that constantly emits logs from a contract  */
type LogEmitterGun struct {
	contract     *contracts.LogEmitter
	eventsToEmit []abi.Event
	logger       zerolog.Logger
	eventsPerTx  int
}

func NewLogEmitterGun(
	contract *contracts.LogEmitter,
	eventsToEmit []abi.Event,
	eventsPerTx int,
	logger zerolog.Logger,
) *LogEmitterGun {
	return &LogEmitterGun{
		contract:     contract,
		eventsToEmit: eventsToEmit,
		eventsPerTx:  eventsPerTx,
		logger:       logger,
	}
}

func (m *LogEmitterGun) Call(l *wasp.Generator) *wasp.CallResult {
	localCounter := 0
	logEmitter := (*m.contract)
	address := logEmitter.Address()
	for _, event := range m.eventsToEmit {
		m.logger.Debug().Str("Emitter address", address.String()).Str("Event type", event.Name).Msg("Emitting log from emitter")
		var err error
		switch event.Name {
		case "Log1":
			_, err = logEmitter.EmitLogInts(getIntSlice(m.eventsPerTx))
		case "Log2":
			_, err = logEmitter.EmitLogIntsIndexed(getIntSlice(m.eventsPerTx))
		case "Log3":
			_, err = logEmitter.EmitLogStrings(getStringSlice(m.eventsPerTx))
		default:
			err = fmt.Errorf("Unknown event name: %s", event.Name)
		}

		if err != nil {
			return &wasp.CallResult{Error: err.Error(), Failed: true}
		}
		localCounter += m.eventsPerTx * 3
	}

	return &wasp.CallResult{
		Data: localCounter,
	}
}
