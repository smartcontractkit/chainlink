package offchainreporting

import (
	"github.com/smartcontractkit/chainlink/offchainreporting/types"
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
}

func NewBootstrapNode(args BootstrapNodeArgs) (*BootstrapNode, error) {
	return &BootstrapNode{
		bootstrapArgs: args,
		started:       semaphore.NewWeighted(1),
	}, nil
}

func (b *BootstrapNode) Start() error {
	b.failIfAlreadyStarted()

	return nil
}

func (b *BootstrapNode) Close() error {
	return nil
}

func (b *BootstrapNode) failIfAlreadyStarted() {
	if !b.started.TryAcquire(1) {
		panic("can only start a BootstrapNode once")
	}
}
