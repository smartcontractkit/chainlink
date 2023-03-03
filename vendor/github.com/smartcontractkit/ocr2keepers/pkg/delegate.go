package ocr2keepers

import (
	"context"
	"fmt"
	"log"

	offchainreporting "github.com/smartcontractkit/libocr/offchainreporting2"
	"github.com/smartcontractkit/ocr2keepers/internal/keepers"
)

// Delegate is a container struct for an Oracle plugin. This struct provides
// the ability to start and stop underlying services associated with the
// plugin instance.
type Delegate struct {
	keeper *offchainreporting.Oracle
}

// NewDelegate provides a new Delegate from a provided config. A new logger
// is defined that wraps the configured logger with a default Go logger.
// The plugin uses a *log.Logger by default so all log output from the
// built-in logger are written to the provided logger as Debug logs prefaced
// with '[keepers-plugin] ' and a short file name.
func NewDelegate(c DelegateConfig) (*Delegate, error) {
	wrapper := &logWriter{l: c.Logger}
	l := log.New(wrapper, "[keepers-plugin] ", log.Lshortfile)

	// set some defaults
	conf := keepers.ReportingFactoryConfig{
		CacheExpiration:       DefaultCacheExpiration,
		CacheEvictionInterval: DefaultCacheClearInterval,
		MaxServiceWorkers:     DefaultMaxServiceWorkers,
		ServiceQueueLength:    DefaultServiceQueueLength,
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

	l.Printf("creating oracle with reporting factory config: %+v", conf)

	// create the oracle from config values
	keeper, err := offchainreporting.NewOracle(offchainreporting.OracleArgs{
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
		ReportingPluginFactory: keepers.NewReportingPluginFactory(
			c.HeadSubscriber,
			c.Registry,
			c.PerformLogProvider,
			c.ReportEncoder,
			l,
			conf,
		),
	})

	if err != nil {
		return nil, fmt.Errorf("%w: failed to create new OCR oracle", err)
	}

	return &Delegate{keeper: keeper}, nil
}

// Start starts the OCR oracle and any associated services
func (d *Delegate) Start(_ context.Context) error {
	if err := d.keeper.Start(); err != nil {
		return fmt.Errorf("%w: failed to start keeper oracle", err)
	}
	return nil
}

// Close stops the OCR oracle and any associated services
func (d *Delegate) Close() error {
	if err := d.keeper.Close(); err != nil {
		return fmt.Errorf("%w: failed to close keeper oracle", err)
	}
	return nil
}
