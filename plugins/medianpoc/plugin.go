package medianpoc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/reportingplugins"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
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
	return "", fmt.Errorf("no pipeline found for %s", key)
}

func (p *Plugin) NewReportingPluginFactory(
	ctx context.Context,
	config types.ReportingPluginServiceConfig,
	provider types.MedianProvider,
	pipelineRunner types.PipelineRunnerService,
	telemetry types.TelemetryClient,
	errorLog types.ErrorLog,
) (types.ReportingPluginFactory, error) {
	f, err := p.newFactory(ctx, config, provider, pipelineRunner, telemetry, errorLog)
	if err != nil {
		return nil, err
	}
	s := &reportingPluginFactoryService{lggr: p.Logger, ReportingPluginFactory: f}
	p.SubService(s)
	return s, nil
}

func (p *Plugin) newFactory(ctx context.Context, config types.ReportingPluginServiceConfig, provider types.MedianProvider, pipelineRunner types.PipelineRunnerService, telemetry types.TelemetryClient, errorLog types.ErrorLog) (*median.NumericalMedianFactory, error) {
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
	factory := &median.NumericalMedianFactory{
		ContractTransmitter:       provider.MedianContract(),
		DataSource:                ds,
		JuelsPerFeeCoinDataSource: jds,
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
