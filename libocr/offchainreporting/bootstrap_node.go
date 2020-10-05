package offchainreporting

import (
	"context"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/internal/managed"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"
	"github.com/smartcontractkit/chainlink/libocr/subprocesses"
	"golang.org/x/sync/semaphore"
)

type BootstrapNodeArgs struct {
	BootstrapperFactory   types.BootstrapperFactory
	Bootstrappers         []string
	ContractConfigTracker types.ContractConfigTracker
	Database              types.Database
	LocalConfig           types.LocalConfig
	Logger                types.Logger
}

type BootstrapNode struct {
	bootstrapArgs BootstrapNodeArgs

	started *semaphore.Weighted

	subprocesses subprocesses.Subprocesses

	cancel context.CancelFunc
}

func NewBootstrapNode(args BootstrapNodeArgs) (*BootstrapNode, error) {
	if err := validateLocalConfig(args.LocalConfig); err != nil {
		return nil, errors.Wrapf(err,
			"bad local config while creating bootstrap node")
	}
	return &BootstrapNode{
		bootstrapArgs: args,
		started:       semaphore.NewWeighted(1),
	}, nil
}

func (b *BootstrapNode) Start() error {
	b.failIfAlreadyStarted()

	ctx, cancel := context.WithCancel(context.Background())
	b.cancel = cancel
	b.subprocesses.Go(func() {
		defer cancel()
		managed.RunManagedBootstrapNode(
			ctx,

			b.bootstrapArgs.BootstrapperFactory,
			b.bootstrapArgs.Bootstrappers,
			b.bootstrapArgs.ContractConfigTracker,
			b.bootstrapArgs.Database,
			b.bootstrapArgs.LocalConfig,
			b.bootstrapArgs.Logger,
		)
	})
	return nil
}

func (b *BootstrapNode) Close() error {
	if b.cancel != nil {
		b.cancel()
	}
	b.subprocesses.Wait()
	return nil
}

func (b *BootstrapNode) failIfAlreadyStarted() {
	if !b.started.TryAcquire(1) {
		panic("can only start a BootstrapNode once")
	}
}
