package networking

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"go.uber.org/multierr"
)

type bootstrapperV1V2 struct {
	stateMu   sync.Mutex
	state     bootstrapperState
	v2started bool

	v1 commontypes.Bootstrapper
	v2 commontypes.Bootstrapper

	logger loghelper.LoggerWithContext
}

func newBootstrapperV1V2(logger loghelper.LoggerWithContext, v1 commontypes.Bootstrapper, v2 commontypes.Bootstrapper) (*bootstrapperV1V2, error) {
	if v1 == nil || v2 == nil {
		return nil, errors.New("bootstrappers can't be nil")
	}
	return &bootstrapperV1V2{sync.Mutex{}, bootstrapperUnstarted, true, v1, v2, logger}, nil
}

// Start starts the underlying v1 and v2 bootstrappers. In case the v2
// bootstrapper fails, we log an error but do not return it. bootstrapperV1V2 is
// designed to be resilient against v2 failing to Start and will operate using
// only v1 if needed.
func (b *bootstrapperV1V2) Start() error {
	succeeded := false
	defer func() {
		if !succeeded {
			b.logger.Warn("BootstrapperV1V2: Start: errored, auto-closing", nil)
			b.Close()
		}
	}()

	b.stateMu.Lock()
	defer b.stateMu.Unlock()

	if b.state != bootstrapperUnstarted {
		return fmt.Errorf("cannot Start bootstrapperV1V2 that is in state %v", b.state)
	}
	b.state = bootstrapperStarted

	if err := b.v1.Start(); err != nil {
		b.logger.Warn("BootstrapperV1V2: Start: Failed to start v1", commontypes.LogFields{"err": err})
		return err
	}
	b.logger.Info("BootstrapperV1V2: Start: v1 started successfully", nil)
	if err := b.v2.Start(); err != nil {
		b.logger.Critical("BootstrapperV1V2: Start: Failed to start v2 bootstrapper as part of v1v2, operating only with v1", commontypes.LogFields{"error": err})
		b.v2started = false
	}
	succeeded = true
	return nil
}

func (b *bootstrapperV1V2) Close() error {
	b.logger.Debug("BootstrapperV1V2: Close", nil)
	b.stateMu.Lock()
	defer b.stateMu.Unlock()
	if b.state != bootstrapperStarted {
		return fmt.Errorf("cannot Close bootstrapperV1V2 that is in state %v", b.state)
	}
	b.state = bootstrapperClosed
	var allErrors error
	allErrors = multierr.Append(allErrors, b.v1.Close())
	if b.v2started {
		allErrors = multierr.Append(allErrors, b.v2.Close())
	}
	return allErrors
}
