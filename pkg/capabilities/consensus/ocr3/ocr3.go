package ocr3

import (
	"context"
	"time"

	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/reportingplugins"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

type Capability struct {
	loop.Plugin
	reportingplugins.PluginProviderServer
	config Config
}

type Config struct {
	RequestTimeout    *time.Duration
	BatchSize         int
	Logger            logger.Logger
	AggregatorFactory types.AggregatorFactory
	EncoderFactory    types.EncoderFactory
	SendBufferSize    int

	store      *store
	capability *capability
	clock      clockwork.Clock
}

const (
	defaultRequestExpiry  time.Duration = 20 * time.Second
	defaultBatchSize                    = 20
	defaultSendBufferSize               = 10
)

func NewOCR3(config Config) *Capability {
	if config.RequestTimeout == nil {
		dre := defaultRequestExpiry
		config.RequestTimeout = &dre
	}

	if config.BatchSize == 0 {
		config.BatchSize = defaultBatchSize
	}

	if config.SendBufferSize == 0 {
		config.SendBufferSize = defaultSendBufferSize
	}

	if config.clock == nil {
		config.clock = clockwork.NewRealClock()
	}

	if config.store == nil {
		config.store = newStore()
	}

	if config.capability == nil {
		ci := newCapability(config.store, config.clock, *config.RequestTimeout, config.AggregatorFactory, config.EncoderFactory, config.Logger,
			config.SendBufferSize)
		config.capability = ci
	}

	cp := &Capability{
		Plugin:               loop.Plugin{Logger: config.Logger},
		PluginProviderServer: reportingplugins.PluginProviderServer{},
		config:               config,
	}

	cp.SubService(config.capability)
	return cp
}

func (o *Capability) NewReportingPluginFactory(ctx context.Context, cfg core.ReportingPluginServiceConfig,
	provider commontypes.PluginProvider, pipelineRunner core.PipelineRunnerService, telemetry core.TelemetryClient,
	errorLog core.ErrorLog, capabilityRegistry core.CapabilitiesRegistry, keyValueStore core.KeyValueStore,
	relayerSet core.RelayerSet) (core.OCR3ReportingPluginFactory, error) {
	factory, err := newFactory(o.config.store, o.config.capability, o.config.BatchSize, o.config.Logger)
	if err != nil {
		return nil, err
	}

	err = capabilityRegistry.Add(ctx, o.config.capability)
	if err != nil {
		return nil, err
	}

	return factory, err
}

func (o *Capability) NewValidationService(ctx context.Context) (core.ValidationService, error) {
	s := &validationService{lggr: o.Logger}
	o.SubService(s)
	return s, nil
}
