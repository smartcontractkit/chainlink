package logpoller

import (
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/wasp"

	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

/* LogEmitterGun is a gun that constantly emits logs from a contract  */
type LogEmitterGun struct {
	contract     *contracts.LogEmitter
	eventsToEmit []abi.Event
	logger       zerolog.Logger
	eventsPerTx  int
}

type Counter struct {
	mu    *sync.Mutex
	value int
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

func (m *LogEmitterGun) Call(l *wasp.Generator) *wasp.Response {
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
			err = fmt.Errorf("unknown event name: %s", event.Name)
		}

		if err != nil {
			return &wasp.Response{Error: err.Error(), Failed: true}
		}
		localCounter++
	}

	// I don't think that will work as expected, I should atomically read the value and save it, so maybe just a mutex?
	if counter, ok := l.InputSharedData().(*Counter); ok {
		counter.mu.Lock()
		defer counter.mu.Unlock()
		counter.value += localCounter
	} else {
		return &wasp.Response{
			Error:  "SharedData did not contain a Counter",
			Failed: true,
		}
	}

	return &wasp.Response{}
}
