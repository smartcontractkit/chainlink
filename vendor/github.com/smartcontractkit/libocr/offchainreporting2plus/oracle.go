package offchainreporting2plus

import (
	"context"
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/managed"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/subprocesses"
)

type OracleArgs interface {
	oracleArgsMarker()
	localConfig() types.LocalConfig
	runManaged(ctx context.Context)
}

// OCR2OracleArgs contains the configuration and services a caller must provide, in
// order to run the offchainreporting protocol.
type OCR2OracleArgs struct {
	// A factory for producing network endpoints. A network endpoints consists of
	// networking methods a consumer must implement to allow a node to
	// communicate with other participating nodes.
	BinaryNetworkEndpointFactory types.BinaryNetworkEndpointFactory

	// V2Bootstrappers is the list of bootstrap node addresses and IDs for the v2 stack.
	V2Bootstrappers []commontypes.BootstrapperLocator

	// Tracks configuration changes.
	ContractConfigTracker types.ContractConfigTracker

	// Interfaces with the OCR2Aggregator smart contract's transmission related logic.
	ContractTransmitter types.ContractTransmitter

	// Database provides persistent storage.
	Database types.Database

	// LocalConfig contains oracle-specific configuration details which are not
	// mandated by the on-chain configuration specification via OffchainAggregatoo.SetConfig.
	LocalConfig types.LocalConfig

	// Logger logs stuff.
	Logger commontypes.Logger

	// Used to send logs to a monitor.
	MonitoringEndpoint commontypes.MonitoringEndpoint

	// Enables adding metrics to track. This may be nil.
	MetricsRegisterer prometheus.Registerer

	// Computes a config digest using purely offchain logic.
	OffchainConfigDigester types.OffchainConfigDigester

	// OffchainKeyring contains the secret keys needed for the OCR protocol, and methods
	// which use those keys without exposing them to the rest of the application.
	OffchainKeyring types.OffchainKeyring

	// OnchainKeyring is used to sign reports that can be validated
	// offchain and by the target contract.
	OnchainKeyring types.OnchainKeyring

	// ReportingPluginFactory creates ReportingPlugins that determine the
	// "application logic" used in a OCR2 protocol instance.
	ReportingPluginFactory types.ReportingPluginFactory
}

func (OCR2OracleArgs) oracleArgsMarker() {}

func (args OCR2OracleArgs) localConfig() types.LocalConfig { return args.LocalConfig }

func (args OCR2OracleArgs) runManaged(ctx context.Context) {
	logger := loghelper.MakeRootLoggerWithContext(args.Logger)

	managed.RunManagedOCR2Oracle(
		ctx,

		args.V2Bootstrappers,
		args.ContractConfigTracker,
		args.ContractTransmitter,
		args.Database,
		args.LocalConfig,
		logger,
		args.MetricsRegisterer,
		args.MonitoringEndpoint,
		args.BinaryNetworkEndpointFactory,
		args.OffchainConfigDigester,
		args.OffchainKeyring,
		args.OnchainKeyring,
		args.ReportingPluginFactory,
	)
}

// MercuryOracleArgs is used for OCR3 protocol instances that only use
// the Mercury-specific subset of OCR3 features.
type MercuryOracleArgs struct {
	// A factory for producing network endpoints. A network endpoints consists of
	// networking methods a consumer must implement to allow a node to
	// communicate with other participating nodes.
	BinaryNetworkEndpointFactory types.BinaryNetworkEndpointFactory

	// V2Bootstrappers is the list of bootstrap node addresses and IDs for the v2 stack.
	V2Bootstrappers []commontypes.BootstrapperLocator

	// Tracks configuration changes.
	ContractConfigTracker types.ContractConfigTracker

	// Interfaces with the OCR2Aggregator smart contract's transmission related logic.
	ContractTransmitter types.ContractTransmitter

	// Database provides persistent storage.
	Database ocr3types.Database

	// LocalConfig contains oracle-specific configuration details which are not
	// mandated by the on-chain configuration specification via OffchainAggregatoo.SetConfig.
	LocalConfig types.LocalConfig

	// Logger logs stuff.
	Logger commontypes.Logger

	// Enables adding metrics to track. This may be nil.
	MetricsRegisterer prometheus.Registerer

	// Used to send logs to a monitor.
	MonitoringEndpoint commontypes.MonitoringEndpoint

	// Computes a config digest using purely offchain logic.
	OffchainConfigDigester types.OffchainConfigDigester

	// OffchainKeyring contains the secret keys needed for the OCR protocol, and methods
	// which use those keys without exposing them to the rest of the application.
	OffchainKeyring types.OffchainKeyring

	// OnchainKeyring is used to sign reports that can be validated
	// offchain and by the target contract.
	OnchainKeyring types.OnchainKeyring

	// ReportingPluginFactory creates ReportingPlugins that determine the
	// "application logic" used in an OCR protocol instance.
	MercuryPluginFactory ocr3types.MercuryPluginFactory
}

func (MercuryOracleArgs) oracleArgsMarker() {}

func (args MercuryOracleArgs) localConfig() types.LocalConfig { return args.LocalConfig }

func (args MercuryOracleArgs) runManaged(ctx context.Context) {
	logger := loghelper.MakeRootLoggerWithContext(args.Logger)

	managed.RunManagedMercuryOracle(
		ctx,

		args.V2Bootstrappers,
		args.ContractConfigTracker,
		args.ContractTransmitter,
		args.Database,
		args.LocalConfig,
		logger,
		args.MetricsRegisterer,
		args.MonitoringEndpoint,
		args.BinaryNetworkEndpointFactory,
		args.OffchainConfigDigester,
		args.OffchainKeyring,
		args.OnchainKeyring,
		args.MercuryPluginFactory,
	)
}

type OCR3OracleArgs[RI any] struct {
	// A factory for producing network endpoints. A network endpoints consists of
	// networking methods a consumer must implement to allow a node to
	// communicate with other participating nodes.
	BinaryNetworkEndpointFactory types.BinaryNetworkEndpointFactory

	// V2Bootstrappers is the list of bootstrap node addresses and IDs for the v2 stack.
	V2Bootstrappers []commontypes.BootstrapperLocator

	// Tracks configuration changes.
	ContractConfigTracker types.ContractConfigTracker

	// Transmit reports to the targeted system (e.g. a blockchain)
	ContractTransmitter ocr3types.ContractTransmitter[RI]

	// Database provides persistent storage.
	Database ocr3types.Database

	// LocalConfig contains oracle-specific configuration details which are not
	// mandated by the on-chain configuration specification via OffchainAggregatoo.SetConfig.
	LocalConfig types.LocalConfig

	// Logger logs stuff.
	Logger commontypes.Logger

	// Enables adding metrics to track. This may be nil.
	MetricsRegisterer prometheus.Registerer

	// Used to send logs to a monitor.
	MonitoringEndpoint commontypes.MonitoringEndpoint

	// Computes a config digest using purely offchain logic.
	OffchainConfigDigester types.OffchainConfigDigester

	// OffchainKeyring contains the secret keys needed for the OCR protocol, and methods
	// which use those keys without exposing them to the rest of the application.
	OffchainKeyring types.OffchainKeyring

	// OnchainKeyring is used to sign reports that can be validated
	// offchain and by the target contract.
	OnchainKeyring ocr3types.OnchainKeyring[RI]

	// PluginFactory creates Plugins that determine the "application logic" used
	// in a protocol instance.
	ReportingPluginFactory ocr3types.ReportingPluginFactory[RI]
}

func (OCR3OracleArgs[RI]) oracleArgsMarker() {}

func (args OCR3OracleArgs[RI]) localConfig() types.LocalConfig { return args.LocalConfig }

func (args OCR3OracleArgs[RI]) runManaged(ctx context.Context) {
	logger := loghelper.MakeRootLoggerWithContext(args.Logger)

	managed.RunManagedOCR3Oracle(
		ctx,

		args.V2Bootstrappers,
		args.ContractConfigTracker,
		args.ContractTransmitter,
		args.Database,
		args.LocalConfig,
		logger,
		args.MetricsRegisterer,
		args.MonitoringEndpoint,
		args.BinaryNetworkEndpointFactory,
		args.OffchainConfigDigester,
		args.OffchainKeyring,
		args.OnchainKeyring,
		args.ReportingPluginFactory,
	)
}

type oracleState int

const (
	oracleStateUnstarted oracleState = iota
	oracleStateStarted
	oracleStateClosed
)

type Oracle interface {
	Start() error
	Close() error
}

type oracle struct {
	lock sync.Mutex

	state oracleState

	oracleArgs OracleArgs

	// subprocesses tracks completion of all go routines on Oracle.Close()
	subprocesses subprocesses.Subprocesses

	// cancel sends a cancel message to all subprocesses, via a context.Context
	cancel context.CancelFunc
}

// NewOracle returns a newly initialized Oracle using the provided services
// and configuration.
func NewOracle(args OracleArgs) (Oracle, error) {
	if err := SanityCheckLocalConfig(args.localConfig()); err != nil {
		return nil, fmt.Errorf("bad local config while creating new oracle: %w", err)
	}
	return &oracle{
		sync.Mutex{},
		oracleStateUnstarted,
		args,
		subprocesses.Subprocesses{},
		nil,
	}, nil
}

// Start spins up a Oracle.
func (o *oracle) Start() error {
	o.lock.Lock()
	defer o.lock.Unlock()

	if o.state != oracleStateUnstarted {
		return fmt.Errorf("can only start Oracle once")
	}
	o.state = oracleStateStarted

	ctx, cancel := context.WithCancel(context.Background())
	o.cancel = cancel
	o.subprocesses.Go(func() {
		defer cancel()

		o.oracleArgs.runManaged(ctx)
	})
	return nil
}

// Close shuts down an oracle. Can safely be called multiple times.
func (o *oracle) Close() error {
	o.lock.Lock()
	defer o.lock.Unlock()

	if o.state != oracleStateStarted {
		return fmt.Errorf("can only close a started oracle")
	}
	o.state = oracleStateClosed

	if o.cancel != nil {
		o.cancel()
	}
	// Wait for all subprocesses to shut down, before shutting down other resources.
	// (Wouldn't want anything to panic from attempting to use a closed resource.)
	o.subprocesses.Wait()
	return nil
}
