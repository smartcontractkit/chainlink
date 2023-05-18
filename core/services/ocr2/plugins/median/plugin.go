package median

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
	libocr "github.com/smartcontractkit/libocr/offchainreporting2"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/v2/plugins"

	v2 "github.com/smartcontractkit/chainlink/v2/core/config/v2"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/median/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type MedianConfig interface {
	JobPipelineMaxSuccessfulRuns() uint64
	plugins.RegistrarConfig
}

// concrete implementation of MedianConfig
type medianConfig struct {
	jobPipelineMaxSuccessfulRuns uint64
	plugins.RegistrarConfig
}

func NewMedianConfig(jobPipelineMaxSuccessfulRuns uint64, pluginProcessCfg plugins.RegistrarConfig) MedianConfig {
	return &medianConfig{
		jobPipelineMaxSuccessfulRuns: jobPipelineMaxSuccessfulRuns,
		RegistrarConfig:              pluginProcessCfg,
	}
}

func (m *medianConfig) JobPipelineMaxSuccessfulRuns() uint64 {
	return m.jobPipelineMaxSuccessfulRuns
}

func NewMedianServices(ctx context.Context,
	jb job.Job,
	isNewlyCreatedJob bool,
	relayer loop.Relayer,
	pipelineRunner pipeline.Runner,
	runResults chan pipeline.Run,
	lggr logger.Logger,
	argsNoPlugin libocr.OracleArgs,
	cfg MedianConfig,
	chEnhancedTelem chan ocrcommon.EnhancedTelemetryData,
	errorLog loop.ErrorLog,

) (srvs []job.ServiceCtx, err error) {
	var pluginConfig config.PluginConfig
	err = json.Unmarshal(jb.OCR2OracleSpec.PluginConfig.Bytes(), &pluginConfig)
	if err != nil {
		return
	}
	err = config.ValidatePluginConfig(pluginConfig)
	if err != nil {
		return
	}
	spec := jb.OCR2OracleSpec

	provider, err := relayer.NewMedianProvider(ctx, types.RelayArgs{
		ExternalJobID: jb.ExternalJobID,
		JobID:         spec.ID,
		ContractID:    spec.ContractID,
		New:           isNewlyCreatedJob,
		RelayConfig:   spec.RelayConfig.Bytes(),
	}, types.PluginArgs{
		TransmitterID: spec.TransmitterID.String,
		PluginConfig:  spec.PluginConfig.Bytes(),
	})
	if err != nil {
		return
	}
	srvs = append(srvs, provider)
	argsNoPlugin.ContractTransmitter = provider.ContractTransmitter()
	argsNoPlugin.ContractConfigTracker = provider.ContractConfigTracker()
	argsNoPlugin.OffchainConfigDigester = provider.OffchainConfigDigester()

	abort := func() {
		var mc services.MultiClose
		for i := range srvs {
			mc = append(mc, srvs[i])
		}
		if cerr := mc.Close(); err != nil {
			lggr.Errorw("Error closing unused services", "err", cerr)
		}
	}

	dataSource, juelsPerFeeCoinSource := ocrcommon.NewDataSourceV2(pipelineRunner,
		jb,
		*jb.PipelineSpec,
		lggr,
		runResults,
		chEnhancedTelem,
	), ocrcommon.NewInMemoryDataSource(pipelineRunner, jb, pipeline.Spec{
		ID:           jb.ID,
		DotDagSource: pluginConfig.JuelsPerFeeCoinPipeline,
		CreatedAt:    time.Now(),
	}, lggr)

	if cmdName := v2.EnvMedianPluginCmd.Get(); cmdName != "" {
		medianLggr := lggr.Named("Median")
		var registeredLoop *plugins.RegisteredLoop
		// use logger name to ensure unique naming
		registeredLoop, err = cfg.RegisterLOOP(medianLggr.Name())
		if err != nil {
			err = fmt.Errorf("failed to register loop: %w", err)
			abort()
			return
		}
		median := loop.NewMedianService(lggr, func() *exec.Cmd {
			cmd := exec.Command(cmdName)
			plugins.SetCmdEnvFromConfig(cmd, registeredLoop.EnvCfg)
			return cmd
		}, provider, dataSource, juelsPerFeeCoinSource, errorLog)
		argsNoPlugin.ReportingPluginFactory = median
		srvs = append(srvs, median)
	} else {
		argsNoPlugin.ReportingPluginFactory, err = NewPlugin(lggr).NewMedianFactory(ctx, provider, dataSource, juelsPerFeeCoinSource, errorLog)
		if err != nil {
			err = fmt.Errorf("failed to create median factory: %w", err)
			abort()
			return
		}
	}

	var oracle *libocr.Oracle
	oracle, err = libocr.NewOracle(argsNoPlugin)
	if err != nil {
		abort()
		return
	}
	runSaver := ocrcommon.NewResultRunSaver(
		runResults,
		pipelineRunner,
		make(chan struct{}),
		lggr,
		cfg.JobPipelineMaxSuccessfulRuns(),
	)
	srvs = append(srvs, runSaver, job.NewServiceAdapter(oracle))
	if !jb.OCR2OracleSpec.CaptureEATelemetry {
		lggr.Infof("Enhanced EA telemetry is disabled for job %s", jb.Name.ValueOrZero())
	}
	return
}

type Plugin struct {
	lggr logger.Logger
	stop utils.StopChan
}

func NewPlugin(lggr logger.Logger) *Plugin {
	return &Plugin{lggr: lggr, stop: make(utils.StopChan)}
}

func (p *Plugin) NewMedianFactory(ctx context.Context, provider types.MedianProvider, dataSource, juelsPerFeeCoin median.DataSource, errorLog loop.ErrorLog) (loop.ReportingPluginFactory, error) {
	var ctxVals loop.ContextValues
	ctxVals.SetValues(ctx)
	lggr := p.lggr.With(ctxVals.Args()...)
	factory := median.NumericalMedianFactory{
		ContractTransmitter:       provider.MedianContract(),
		DataSource:                dataSource,
		JuelsPerFeeCoinDataSource: juelsPerFeeCoin,
		Logger: logger.NewOCRWrapper(lggr, true, func(msg string) {
			ctx, cancelFn := p.stop.NewCtx()
			defer cancelFn()
			if err := errorLog.SaveError(ctx, msg); err != nil {
				lggr.Errorw("Unable to save error", "err", msg)
			}
		}),
		OnchainConfigCodec: provider.OnchainConfigCodec(),
		ReportCodec:        provider.ReportCodec(),
	}
	return &reportingPluginFactoryService{lggr: p.lggr.Named("ReportingPluginFactory"), ReportingPluginFactory: factory}, nil
}

func (p *Plugin) Close() (err error) {
	close(p.stop)
	return
}

type reportingPluginFactoryService struct {
	utils.StartStopOnce
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
