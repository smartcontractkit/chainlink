package offchainreporting2plus

import (
	"context"
	"fmt"
	"sync"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/managed"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/subprocesses"
)

type BootstrapperArgs struct {
	BootstrapperFactory    types.BootstrapperFactory
	V2Bootstrappers        []commontypes.BootstrapperLocator
	ContractConfigTracker  types.ContractConfigTracker
	Database               types.ConfigDatabase
	LocalConfig            types.LocalConfig
	Logger                 commontypes.Logger
	MonitoringEndpoint     commontypes.MonitoringEndpoint
	OffchainConfigDigester types.OffchainConfigDigester
}

type bootstrapperState int

const (
	bootstrapperStateUnstarted bootstrapperState = iota
	bootstrapperStateStarted
	bootstrapperStateClosed
)

// Bootstrapper connects to a particular feed and listens for config changes,
// but does not participate in the protocol. It merely acts as a bootstrap node
// for peer discovery.
type Bootstrapper struct {
	lock sync.Mutex

	state bootstrapperState

	bootstrapArgs BootstrapperArgs

	// subprocesses tracks completion of all go routines on Bootstrapper.Close()
	subprocesses subprocesses.Subprocesses

	// cancel sends a cancel message to all subprocesses, via a context.Context
	cancel context.CancelFunc
}

func NewBootstrapper(args BootstrapperArgs) (*Bootstrapper, error) {
	if err := SanityCheckLocalConfig(args.LocalConfig); err != nil {
		return nil, fmt.Errorf("bad local config while creating Bootstrapper: %w", err)
	}
	return &Bootstrapper{
		sync.Mutex{},
		bootstrapperStateUnstarted,
		args,
		subprocesses.Subprocesses{},
		nil,
	}, nil
}

// Start spins up a Bootstrapper.
func (b *Bootstrapper) Start() error {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.state != bootstrapperStateUnstarted {
		return fmt.Errorf("can only start Bootstrapper once")
	}
	b.state = bootstrapperStateStarted

	ctx, cancel := context.WithCancel(context.Background())
	b.cancel = cancel
	b.subprocesses.Go(func() {
		defer cancel()
		logger := loghelper.MakeRootLoggerWithContext(b.bootstrapArgs.Logger)
		managed.RunManagedBootstrapper(
			ctx,

			b.bootstrapArgs.BootstrapperFactory,
			b.bootstrapArgs.V2Bootstrappers,
			b.bootstrapArgs.ContractConfigTracker,
			b.bootstrapArgs.Database,
			b.bootstrapArgs.LocalConfig,
			logger,
			b.bootstrapArgs.OffchainConfigDigester,
		)
	})
	return nil
}

// Close shuts down a Bootstrapper. Can safely be called multiple times.
func (b *Bootstrapper) Close() error {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.state != bootstrapperStateStarted {
		return fmt.Errorf("can only close a started Bootstrapper")
	}
	b.state = bootstrapperStateClosed

	if b.cancel != nil {
		b.cancel()
	}
	// Wait for all subprocesses to shut down, before shutting down other resources.
	// (Wouldn't want anything to panic from attempting to use a closed resource.)
	b.subprocesses.Wait()
	return nil
}
