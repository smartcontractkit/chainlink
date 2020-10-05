package offchainreporting

import (
	"context"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/internal/managed"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"
	"github.com/smartcontractkit/chainlink/libocr/subprocesses"

	"golang.org/x/sync/semaphore"
)

type OracleArgs struct {
	BinaryNetworkEndpointFactory types.BinaryNetworkEndpointFactory

	Bootstrappers []string

	ContractTransmitter types.ContractTransmitter

	ContractConfigTracker types.ContractConfigTracker

	Database types.Database

	Datasource types.DataSource

	LocalConfig types.LocalConfig

	Logger types.Logger

	MonitoringEndpoint types.MonitoringEndpoint

	PrivateKeys types.PrivateKeys
}

type Oracle struct {
	oracleArgs OracleArgs

	started *semaphore.Weighted

	subprocesses subprocesses.Subprocesses

	cancel context.CancelFunc
}

func NewOracle(args OracleArgs) (*Oracle, error) {
	if err := validateLocalConfig(args.LocalConfig); err != nil {
		return nil, errors.Wrapf(err, "bad local config while creating new oracle")
	}
	return &Oracle{
		oracleArgs: args,
		started:    semaphore.NewWeighted(1),
	}, nil
}

func (o *Oracle) Start() error {
	o.failIfAlreadyStarted()

	ctx, cancel := context.WithCancel(context.Background())
	o.cancel = cancel
	o.subprocesses.Go(func() {
		defer cancel()
		managed.RunManagedOracle(
			ctx,

			o.oracleArgs.Bootstrappers,
			o.oracleArgs.ContractConfigTracker,
			o.oracleArgs.ContractTransmitter,
			o.oracleArgs.Database,
			o.oracleArgs.Datasource,
			o.oracleArgs.LocalConfig,
			o.oracleArgs.Logger,
			o.oracleArgs.BinaryNetworkEndpointFactory,
			o.oracleArgs.PrivateKeys,
		)
	})
	return nil
}

func (o *Oracle) Close() error {
	if o.cancel != nil {
		o.cancel()
	}
	o.subprocesses.Wait()
	return nil
}

func (o *Oracle) failIfAlreadyStarted() {
	if !o.started.TryAcquire(1) {
		panic("can only start an Oracle once")
	}
}
