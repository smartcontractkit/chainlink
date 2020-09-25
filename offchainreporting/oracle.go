package offchainreporting

import (
	"github.com/smartcontractkit/chainlink/offchainreporting/types"

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
}

func NewOracle(args OracleArgs) (*Oracle, error) {
	return &Oracle{
		oracleArgs: args,
		started:    semaphore.NewWeighted(1),
	}, nil
}

func (o *Oracle) Start() error {
	o.failIfAlreadyStarted()

	return nil
}

func (o *Oracle) Close() error {
	return nil
}

func (o *Oracle) failIfAlreadyStarted() {
	if !o.started.TryAcquire(1) {
		panic("can only start an Oracle once")
	}
}
