package medianpoc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/reportingplugins"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

func NewPlugin(lggr logger.Logger) *Plugin {
	return &Plugin{
		Plugin:               loop.Plugin{Logger: lggr},
		MedianProviderServer: reportingplugins.MedianProviderServer{},
		stop:                 make(services.StopChan),
	}
}

type Plugin struct {
	loop.Plugin
	stop services.StopChan
	reportingplugins.MedianProviderServer
}

type PipelineNotFoundError struct {
	Key string
}

func (e *PipelineNotFoundError) Error() string {
	return fmt.Sprintf("no pipeline found for %s", e.Key)
}

func (p *Plugin) NewValidationService(ctx context.Context) (core.ValidationService, error) {
	s := &reportingPluginValidationService{lggr: p.Logger}
	p.SubService(s)
	return s, nil
}

type pipelineSpec struct {
	Name string `json:"name"`
	Spec string `json:"spec"`
}

type jsonConfig struct {
	Pipelines []pipelineSpec `json:"pipelines"`
}

func (j jsonConfig) defaultPipeline() (string, error) {
	return j.getPipeline("__DEFAULT_PIPELINE__")
}

func (j jsonConfig) getPipeline(key string) (string, error) {
	for _, v := range j.Pipelines {
		if v.Name == key {
			return v.Spec, nil
		}
	}
	return "", &PipelineNotFoundError{key}
}

func (p *Plugin) NewReportingPluginFactory(
	ctx context.Context,
	config core.ReportingPluginServiceConfig,
	provider types.MedianProvider,
	pipelineRunner core.PipelineRunnerService,
	telemetry core.TelemetryClient,
	errorLog core.ErrorLog,
	keyValueStore core.KeyValueStore,
	relayerSet core.RelayerSet,
) (types.ReportingPluginFactory, error) {
	f, err := p.newFactory(ctx, config, provider, pipelineRunner, telemetry, errorLog)
	if err != nil {
		return nil, err
	}
	s := &reportingPluginFactoryService{lggr: p.Logger, ReportingPluginFactory: f}
	p.SubService(s)
	return s, nil
}

func (p *Plugin) newFactory(ctx context.Context, config core.ReportingPluginServiceConfig, provider types.MedianProvider, pipelineRunner core.PipelineRunnerService, telemetry core.TelemetryClient, errorLog core.ErrorLog) (*median.NumericalMedianFactory, error) {
	jc := &jsonConfig{}
	err := json.Unmarshal([]byte(config.PluginConfig), jc)
	if err != nil {
		return nil, err
	}

	dp, err := jc.defaultPipeline()
	if err != nil {
		return nil, err
	}
	ds := &DataSource{
		pipelineRunner: pipelineRunner,
		spec:           dp,
		lggr:           p.Logger,
	}

	jfp, err := jc.getPipeline("juelsPerFeeCoinPipeline")
	if err != nil {
		return nil, err
	}
	jds := &DataSource{
		pipelineRunner: pipelineRunner,
		spec:           jfp,
		lggr:           p.Logger,
	}

	var gds median.DataSource
	gp, err := jc.getPipeline("gasPriceSubunitsPipeline")

	var pnf *PipelineNotFoundError
	pipelineNotFound := errors.As(err, &pnf)
	if !pipelineNotFound && err != nil {
		return nil, err
	}

	// We omit gas price in observation to maintain backwards compatibility in libocr (with older nodes).
	// Once all chainlink nodes have updated to libocr version >= fd3cab206b2c
	// the IncludeGasPriceSubunitsInObservation field can be removed

	var includeGasPriceSubunitsInObservation bool
	if pipelineNotFound {
		gds = &ZeroDataSource{}
		includeGasPriceSubunitsInObservation = false
	} else {
		gds = &DataSource{
			pipelineRunner: pipelineRunner,
			spec:           gp,
			lggr:           p.Logger,
		}
		includeGasPriceSubunitsInObservation = true
	}

	factory := &median.NumericalMedianFactory{
		ContractTransmitter:                  provider.MedianContract(),
		DataSource:                           ds,
		JuelsPerFeeCoinDataSource:            jds,
		GasPriceSubunitsDataSource:           gds,
		IncludeGasPriceSubunitsInObservation: includeGasPriceSubunitsInObservation,
		Logger: logger.NewOCRWrapper(
			p.Logger,
			true,
			func(msg string) {},
		),
		OnchainConfigCodec: provider.OnchainConfigCodec(),
		ReportCodec:        provider.ReportCodec(),
	}
	return factory, nil
}

type reportingPluginFactoryService struct {
	services.StateMachine
	lggr logger.Logger
	ocrtypes.ReportingPluginFactory
}

func (r *reportingPluginFactoryService) Name() string { return r.lggr.Name() }

func (r *reportingPluginFactoryService) Start(ctx context.Context) error {
	return r.StartOnce("ReportingPluginFactory", func() error { return nil })
}

func (r *reportingPluginFactoryService) Close() error {
	return r.StopOnce("ReportingPluginFactory", func() error { return nil })
}

func (r *reportingPluginFactoryService) HealthReport() map[string]error {
	return map[string]error{r.Name(): r.Healthy()}
}

type reportingPluginValidationService struct {
	services.StateMachine
	lggr logger.Logger
}

func (r *reportingPluginValidationService) ValidateConfig(ctx context.Context, config map[string]interface{}) error {
	tt, ok := config["telemetryType"]
	if !ok {
		return fmt.Errorf("expected telemtry type")
	}
	telemetryType, ok := tt.(string)
	if !ok {
		return fmt.Errorf("expected telemtry type to be of type string but got %T", tt)
	}
	if telemetryType != "median" {
		return fmt.Errorf("expected telemtry type to be median but got %q", telemetryType)
	}

	return nil
}
func (r *reportingPluginValidationService) Name() string { return r.lggr.Name() }

func (r *reportingPluginValidationService) Start(ctx context.Context) error {
	return r.StartOnce("ValidationService", func() error { return nil })
}

func (r *reportingPluginValidationService) Close() error {
	return r.StopOnce("ValidationService", func() error { return nil })
}

func (r *reportingPluginValidationService) HealthReport() map[string]error {
	return map[string]error{r.Name(): r.Healthy()}
}
