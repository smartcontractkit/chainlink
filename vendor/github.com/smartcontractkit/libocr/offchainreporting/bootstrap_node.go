package offchainreporting

import (
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting/internal/managed"
	"github.com/smartcontractkit/libocr/offchainreporting/types"
	"github.com/smartcontractkit/libocr/subprocesses"
)

type BootstrapNodeArgs struct {
	BootstrapperFactory   types.BootstrapperFactory
	V2Bootstrappers       []commontypes.BootstrapperLocator
	ContractConfigTracker types.ContractConfigTracker
	Database              types.Database
	LocalConfig           types.LocalConfig
	Logger                commontypes.Logger
	MonitoringEndpoint    commontypes.MonitoringEndpoint
}

type bootstrapNodeState int

const (
	bootstrapNodeStateUnstarted bootstrapNodeState = iota
	bootstrapNodeStateStarted
	bootstrapNodeStateClosed
)

// BootstrapNode connects to a particular feed and listens for config changes,
// but does not participate in the protocol. It merely acts as a bootstrap node
// for peer discovery.
type BootstrapNode struct {
	lock sync.Mutex

	state bootstrapNodeState

	bootstrapArgs BootstrapNodeArgs

	// subprocesses tracks completion of all go routines on BootstrapNode.Close()
	subprocesses subprocesses.Subprocesses

	// cancel sends a cancel message to all subprocesses, via a context.Context
	cancel context.CancelFunc
}

func NewBootstrapNode(args BootstrapNodeArgs) (*BootstrapNode, error) {
	if err := SanityCheckLocalConfig(args.LocalConfig); err != nil {
		return nil, errors.Wrapf(err,
			"bad local config while creating BootstrapNode")
	}
	return &BootstrapNode{
		sync.Mutex{},
		bootstrapNodeStateUnstarted,
		args,
		subprocesses.Subprocesses{},
		nil,
	}, nil
}

// Start spins up a BootstrapNode.
func (b *BootstrapNode) Start() error {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.state != bootstrapNodeStateUnstarted {
		return fmt.Errorf("can only start BootstrapNode once")
	}
	b.state = bootstrapNodeStateStarted

	ctx, cancel := context.WithCancel(context.Background())
	b.cancel = cancel
	b.subprocesses.Go(func() {
		defer cancel()
		logger := loghelper.MakeRootLoggerWithContext(b.bootstrapArgs.Logger)
		managed.RunManagedBootstrapNode(
			ctx,

			b.bootstrapArgs.BootstrapperFactory,
			b.bootstrapArgs.V2Bootstrappers,
			b.bootstrapArgs.ContractConfigTracker,
			b.bootstrapArgs.Database,
			b.bootstrapArgs.LocalConfig,
			logger,
		)
	})
	return nil
}

// Close shuts down a BootstrapNode. Can safely be called multiple times.
func (b *BootstrapNode) Close() error {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.state != bootstrapNodeStateStarted {
		return fmt.Errorf("can only close a started BootstrapNode")
	}
	b.state = bootstrapNodeStateClosed

	if b.cancel != nil {
		b.cancel()
	}
	// Wait for all subprocesses to shut down, before shutting down other resources.
	// (Wouldn't want anything to panic from attempting to use a closed resource.)
	b.subprocesses.Wait()
	return nil
}
