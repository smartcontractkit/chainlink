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
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

type Capability struct {
	loop.Plugin
	reportingplugins.PluginProviderServer
	config Config
}

type Config struct {
	RequestTimeout *time.Duration
	BatchSize      int
	Logger         logger.Logger
	EncoderFactory EncoderFactory

	store      *store
	capability *capability
	clock      clockwork.Clock
}

const (
	defaultRequestExpiry time.Duration = 1 * time.Hour
	defaultBatchSize                   = 1000
)

type EncoderFactory func(config *values.Map) (types.Encoder, error)

func NewOCR3(config Config) *Capability {
	if config.RequestTimeout == nil {
		dre := defaultRequestExpiry
		config.RequestTimeout = &dre
	}

	if config.BatchSize == 0 {
		config.BatchSize = defaultBatchSize
	}

	if config.clock == nil {
		config.clock = clockwork.NewRealClock()
	}

	if config.store == nil {
		config.store = newStore()
	}

	if config.capability == nil {
		ci := newCapability(config.store, config.clock, *config.RequestTimeout, config.EncoderFactory, config.Logger)
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

func (o *Capability) NewReportingPluginFactory(ctx context.Context, cfg commontypes.ReportingPluginServiceConfig, provider commontypes.PluginProvider, pipelineRunner commontypes.PipelineRunnerService, telemetry commontypes.TelemetryClient, errorLog commontypes.ErrorLog, capabilityRegistry commontypes.CapabilitiesRegistry, keyValueStore commontypes.KeyValueStore) (commontypes.OCR3ReportingPluginFactory, error) {
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

func (o *Capability) NewValidationService(ctx context.Context) (commontypes.ValidationService, error) {
	s := &validationService{lggr: o.Logger}
	o.SubService(s)
	return s, nil
}
