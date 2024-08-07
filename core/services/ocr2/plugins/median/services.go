package median

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	libocr_median "github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-feeds/median"
	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/median/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

type MedianConfig interface {
	JobPipelineMaxSuccessfulRuns() uint64
	JobPipelineResultWriteQueueDepth() uint64
	plugins.RegistrarConfig
}

// concrete implementation of MedianConfig
type medianConfig struct {
	jobPipelineMaxSuccessfulRuns     uint64
	jobPipelineResultWriteQueueDepth uint64
	plugins.RegistrarConfig
}

func NewMedianConfig(jobPipelineMaxSuccessfulRuns uint64, jobPipelineResultWriteQueueDepth uint64, pluginProcessCfg plugins.RegistrarConfig) MedianConfig {
	return &medianConfig{
		jobPipelineMaxSuccessfulRuns:     jobPipelineMaxSuccessfulRuns,
		jobPipelineResultWriteQueueDepth: jobPipelineResultWriteQueueDepth,
		RegistrarConfig:                  pluginProcessCfg,
	}
}

func (m *medianConfig) JobPipelineMaxSuccessfulRuns() uint64 {
	return m.jobPipelineMaxSuccessfulRuns
}

func (m *medianConfig) JobPipelineResultWriteQueueDepth() uint64 {
	return m.jobPipelineResultWriteQueueDepth
}

func NewMedianServices(ctx context.Context,
	jb job.Job,
	isNewlyCreatedJob bool,
	relayer loop.Relayer,
	kvStore job.KVStore,
	pipelineRunner pipeline.Runner,
	lggr logger.Logger,
	argsNoPlugin libocr.OCR2OracleArgs,
	cfg MedianConfig,
	chEnhancedTelem chan ocrcommon.EnhancedTelemetryData,
	errorLog loop.ErrorLog,

) (srvs []job.ServiceCtx, err error) {
	var pluginConfig config.PluginConfig
	err = json.Unmarshal(jb.OCR2OracleSpec.PluginConfig.Bytes(), &pluginConfig)
	if err != nil {
		return
	}
	err = pluginConfig.ValidatePluginConfig()
	if err != nil {
		return
	}
	spec := jb.OCR2OracleSpec

	runSaver := ocrcommon.NewResultRunSaver(
		pipelineRunner,
		lggr,
		cfg.JobPipelineMaxSuccessfulRuns(),
		cfg.JobPipelineResultWriteQueueDepth(),
	)

	provider, err := relayer.NewPluginProvider(ctx, types.RelayArgs{
		ExternalJobID: jb.ExternalJobID,
		JobID:         jb.ID,
		ContractID:    spec.ContractID,
		New:           isNewlyCreatedJob,
		RelayConfig:   spec.RelayConfig.Bytes(),
		ProviderType:  string(spec.PluginType),
	}, types.PluginArgs{
		TransmitterID: spec.TransmitterID.String,
		PluginConfig:  spec.PluginConfig.Bytes(),
	})
	if err != nil {
		return
	}

	medianProvider, ok := provider.(types.MedianProvider)
	if !ok {
		return nil, errors.New("could not coerce PluginProvider to MedianProvider")
	}

	srvs = append(srvs, provider)
	argsNoPlugin.ContractTransmitter = provider.ContractTransmitter()
	argsNoPlugin.ContractConfigTracker = provider.ContractConfigTracker()
	argsNoPlugin.OffchainConfigDigester = provider.OffchainConfigDigester()

	abort := func() {
		if cerr := services.MultiCloser(srvs).Close(); err != nil {
			lggr.Errorw("Error closing unused services", "err", cerr)
		}
	}

	dataSource := ocrcommon.NewDataSourceV2(pipelineRunner,
		jb,
		*jb.PipelineSpec,
		lggr,
		runSaver,
		chEnhancedTelem)

	juelsPerFeeCoinSource := ocrcommon.NewInMemoryDataSource(pipelineRunner, jb, pipeline.Spec{
		ID:           jb.ID,
		DotDagSource: pluginConfig.JuelsPerFeeCoinPipeline,
		CreatedAt:    time.Now(),
	}, lggr)

	if pluginConfig.JuelsPerFeeCoinCache == nil || (pluginConfig.JuelsPerFeeCoinCache != nil && !pluginConfig.JuelsPerFeeCoinCache.Disable) {
		lggr.Infof("juelsPerFeeCoin data source caching is enabled")
		juelsPerFeeCoinSourceCache, err2 := ocrcommon.NewInMemoryDataSourceCache(juelsPerFeeCoinSource, kvStore, pluginConfig.JuelsPerFeeCoinCache)
		if err2 != nil {
			return nil, err2
		}
		juelsPerFeeCoinSource = juelsPerFeeCoinSourceCache
		srvs = append(srvs, juelsPerFeeCoinSourceCache)
	}

	var gasPriceSubunitsDataSource libocr_median.DataSource
	if pluginConfig.HasGasPriceSubunitsPipeline() {
		gasPriceSubunitsDataSource = ocrcommon.NewInMemoryDataSource(pipelineRunner, jb, pipeline.Spec{
			ID:           jb.ID,
			DotDagSource: pluginConfig.GasPriceSubunitsPipeline,
			CreatedAt:    time.Now(),
		}, lggr)
	} else {
		gasPriceSubunitsDataSource = &median.ZeroDataSource{}
	}

	if cmdName := env.MedianPlugin.Cmd.Get(); cmdName != "" {
		// use unique logger names so we can use it to register a loop
		medianLggr := lggr.Named("Median").Named(spec.ContractID).Named(spec.GetID())
		envVars, err2 := plugins.ParseEnvFile(env.MedianPlugin.Env.Get())
		if err2 != nil {
			err = fmt.Errorf("failed to parse median env file: %w", err2)
			abort()
			return
		}
		cmdFn, telem, err2 := cfg.RegisterLOOP(plugins.CmdConfig{
			ID:  medianLggr.Name(),
			Cmd: cmdName,
			Env: envVars,
		})
		if err2 != nil {
			err = fmt.Errorf("failed to register loop: %w", err2)
			abort()
			return
		}
		median := loop.NewMedianService(lggr, telem, cmdFn, medianProvider, spec.ContractID, dataSource, juelsPerFeeCoinSource, gasPriceSubunitsDataSource, errorLog)
		argsNoPlugin.ReportingPluginFactory = median
		srvs = append(srvs, median)
	} else {
		argsNoPlugin.ReportingPluginFactory, err = median.NewPlugin(lggr).NewMedianFactory(ctx, medianProvider, spec.ContractID, dataSource, juelsPerFeeCoinSource, gasPriceSubunitsDataSource, errorLog)
		if err != nil {
			err = fmt.Errorf("failed to create median factory: %w", err)
			abort()
			return
		}
	}

	var oracle libocr.Oracle
	oracle, err = libocr.NewOracle(argsNoPlugin)
	if err != nil {
		abort()
		return
	}
	srvs = append(srvs, runSaver, job.NewServiceAdapter(oracle))
	if !jb.OCR2OracleSpec.CaptureEATelemetry {
		lggr.Infof("Enhanced EA telemetry is disabled for job %s", jb.Name.ValueOrZero())
	}
	return
}
