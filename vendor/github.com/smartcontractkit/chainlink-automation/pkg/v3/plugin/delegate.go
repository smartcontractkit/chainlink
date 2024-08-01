package plugin

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
	offchainreporting "github.com/smartcontractkit/libocr/offchainreporting2plus"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocr2plustypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-automation/pkg/v3/config"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/runner"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/telemetry"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/types"
	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"
)

var (
	newOracleFn = offchainreporting.NewOracle
)

type oracle interface {
	Start() error
	Close() error
}

// DelegateConfig provides a single configuration struct for all options
// to be passed to the oracle, oracle factory, and underlying plugin/services.
type DelegateConfig struct {
	BinaryNetworkEndpointFactory ocr2plustypes.BinaryNetworkEndpointFactory
	V2Bootstrappers              []commontypes.BootstrapperLocator
	ContractConfigTracker        ocr2plustypes.ContractConfigTracker
	ContractTransmitter          ocr3types.ContractTransmitter[AutomationReportInfo]
	KeepersDatabase              ocr3types.Database
	Logger                       commontypes.Logger
	MonitoringEndpoint           commontypes.MonitoringEndpoint
	OffchainConfigDigester       ocr2plustypes.OffchainConfigDigester
	OffchainKeyring              ocr2plustypes.OffchainKeyring
	OnchainKeyring               ocr3types.OnchainKeyring[AutomationReportInfo]
	LocalConfig                  ocr2plustypes.LocalConfig
	MetricsRegisterer            prometheus.Registerer

	// LogProvider allows reads on the latest log events ready to be processed
	LogProvider ocr2keepers.LogEventProvider

	// EventProvider allows reads on latest transmit events
	EventProvider types.TransmitEventProvider

	// Runnable is a check pipeline runner
	Runnable types.Runnable

	// Encoder provides methods to encode/decode reports
	Encoder ocr2keepers.Encoder

	// BlockSubscriber provides subscribe/unsubscribe methods for block source
	// data
	BlockSubscriber ocr2keepers.BlockSubscriber

	// RecoverableProvider provides recoverable payloads to be proposed to the
	// network
	RecoverableProvider ocr2keepers.RecoverableProvider

	// PayloadBuilder provides a utility method to build an upkeep payload from
	// a trigger id and block key
	PayloadBuilder ocr2keepers.PayloadBuilder

	// UpkeepProvider ...
	UpkeepProvider ocr2keepers.ConditionalUpkeepProvider

	// UpkeepStateUpdater
	UpkeepStateUpdater ocr2keepers.UpkeepStateUpdater

	// Methods passed from core
	UpkeepTypeGetter types.UpkeepTypeGetter
	WorkIDGenerator  types.WorkIDGenerator

	// CacheExpiration is the duration of time a cached key is available. Use
	// this value to balance memory usage and RPC calls. A new set of keys is
	// generated with every block so a good setting might come from block time
	// times number of blocks of history to support not replaying reports.
	CacheExpiration time.Duration

	// CacheEvictionInterval is a parameter for how often the cache attempts to
	// evict expired keys. This value should be short enough to ensure key
	// eviction doesn't block for too long, and long enough that it doesn't
	// cause frequent blocking.
	CacheEvictionInterval time.Duration

	// MaxServiceWorkers is the total number of go-routines allowed to make RPC
	// simultaneous calls on behalf of the sampling operation. This parameter
	// is 10x the number of available CPUs by default. The RPC calls are memory
	// heavy as opposed to CPU heavy as most of the work involves waiting on
	// network responses.
	MaxServiceWorkers int

	// ServiceQueueLength is the buffer size for the RPC service queue. Fewer
	// workers or slower RPC responses will cause this queue to build up.
	// Adding new items to the queue will block if the queue becomes full.
	ServiceQueueLength int
}

// Delegate is a container struct for an Oracle plugin. This struct provides
// the ability to start and stop underlying services associated with the
// plugin instance.
type Delegate struct {
	keeper oracle
	logger *log.Logger
}

// NewDelegate provides a new Delegate from a provided config. A new logger
// is defined that wraps the configured logger with a default Go logger.
// The plugin uses a *log.Logger by default so all log output from the
// built-in logger are written to the provided logger as Debug logs prefaced
// with '[keepers-plugin] ' and a short file name.
func NewDelegate(c DelegateConfig) (*Delegate, error) {
	// set some defaults
	conf := config.ReportingFactoryConfig{
		CacheExpiration:       config.DefaultCacheExpiration,
		CacheEvictionInterval: config.DefaultCacheClearInterval,
		MaxServiceWorkers:     config.DefaultMaxServiceWorkers,
		ServiceQueueLength:    config.DefaultServiceQueueLength,
	}

	// override if set in config
	if c.CacheExpiration != 0 {
		conf.CacheExpiration = c.CacheExpiration
	}

	if c.CacheEvictionInterval != 0 {
		conf.CacheEvictionInterval = c.CacheEvictionInterval
	}

	if c.MaxServiceWorkers != 0 {
		conf.MaxServiceWorkers = c.MaxServiceWorkers
	}

	if c.ServiceQueueLength != 0 {
		conf.ServiceQueueLength = c.ServiceQueueLength
	}

	// the log wrapper is to be able to use a log.Logger everywhere instead of
	// a variety of logger ocr2plustypes. all logs write to the Debug method.
	wrapper := &logWriter{l: c.Logger}
	l := log.New(wrapper, fmt.Sprintf("[%s] ", telemetry.ServiceName), log.Lshortfile)

	l.Printf("creating oracle with reporting factory config: %+v", conf)

	// create the oracle from config values
	keeper, err := newOracleFn(offchainreporting.OCR3OracleArgs[AutomationReportInfo]{
		BinaryNetworkEndpointFactory: c.BinaryNetworkEndpointFactory,
		V2Bootstrappers:              c.V2Bootstrappers,
		ContractConfigTracker:        c.ContractConfigTracker,
		ContractTransmitter:          c.ContractTransmitter,
		Database:                     c.KeepersDatabase,
		LocalConfig:                  c.LocalConfig,
		Logger:                       c.Logger,
		MonitoringEndpoint:           c.MonitoringEndpoint,
		OffchainConfigDigester:       c.OffchainConfigDigester,
		OffchainKeyring:              c.OffchainKeyring,
		OnchainKeyring:               c.OnchainKeyring,
		MetricsRegisterer:            c.MetricsRegisterer,

		ReportingPluginFactory: NewReportingPluginFactory(
			c.LogProvider,
			c.EventProvider,
			c.BlockSubscriber,
			c.RecoverableProvider,
			c.PayloadBuilder,
			c.UpkeepProvider,
			c.Runnable,
			runner.RunnerConfig{
				Workers:           conf.MaxServiceWorkers,
				WorkerQueueLength: conf.ServiceQueueLength,
				CacheExpire:       conf.CacheExpiration,
				CacheClean:        conf.CacheEvictionInterval,
			},
			c.Encoder,
			c.UpkeepTypeGetter,
			c.WorkIDGenerator,
			c.UpkeepStateUpdater,
			l,
		),
	})

	if err != nil {
		return nil, fmt.Errorf("%w: failed to create new OCR oracle", err)
	}

	return &Delegate{
		keeper: keeper,
		logger: l,
	}, nil
}

// Start starts the OCR oracle and any associated services
func (d *Delegate) Start(_ context.Context) error {
	d.logger.Println("starting oracle")

	if err := d.keeper.Start(); err != nil {
		return fmt.Errorf("%w: failed to start keeper oracle", err)
	}

	return nil
}

// Close stops the OCR oracle and any associated services
func (d *Delegate) Close() error {
	d.logger.Println("stopping oracle")

	if err := d.keeper.Close(); err != nil {
		return fmt.Errorf("%w: failed to close keeper oracle", err)
	}

	return nil
}

type logWriter struct {
	l commontypes.Logger
}

func (l *logWriter) Write(p []byte) (n int, err error) {
	l.l.Debug(string(p), nil)
	n = len(p)
	return n, nil
}
