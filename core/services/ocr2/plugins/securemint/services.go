package securemint

// TODO(gg): maybe this should live in ocr3 instead of ocr2?

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	sm_ea "github.com/smartcontractkit/chainlink/v2/core/services/ocr3/securemint/external_adapter"
	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus"
	ocr2plus_types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/por_mock_ocr3plugin/por"
	sm_plugin "github.com/smartcontractkit/por_mock_ocr3plugin/por"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

type SecureMintConfig interface {
	JobPipelineMaxSuccessfulRuns() uint64
	JobPipelineResultWriteQueueDepth() uint64
	plugins.RegistrarConfig
}

// concrete implementation of SecureMintConfig
type secureMintConfig struct {
	jobPipelineMaxSuccessfulRuns     uint64
	jobPipelineResultWriteQueueDepth uint64
	plugins.RegistrarConfig
}

func NewSecureMintConfig(jobPipelineMaxSuccessfulRuns uint64, jobPipelineResultWriteQueueDepth uint64, pluginProcessCfg plugins.RegistrarConfig) SecureMintConfig {
	return &secureMintConfig{
		jobPipelineMaxSuccessfulRuns:     jobPipelineMaxSuccessfulRuns,
		jobPipelineResultWriteQueueDepth: jobPipelineResultWriteQueueDepth,
		RegistrarConfig:                  pluginProcessCfg,
	}
}

func (m *secureMintConfig) JobPipelineMaxSuccessfulRuns() uint64 {
	return m.jobPipelineMaxSuccessfulRuns
}

func (m *secureMintConfig) JobPipelineResultWriteQueueDepth() uint64 {
	return m.jobPipelineResultWriteQueueDepth
}

// Create all securemint plugin Oracles and all extra services needed to run a SecureMint job.
func NewSecureMintServices(ctx context.Context,
	jb job.Job,
	isNewlyCreatedJob bool,
	relayer loop.Relayer,
	kvStore job.KVStore,
	pipelineRunner pipeline.Runner,
	lggr logger.Logger,
	argsNoPlugin libocr.OCR3OracleArgs[por.ChainSelector],
	cfg SecureMintConfig,
	chEnhancedTelem chan ocrcommon.EnhancedTelemetryData,
	errorLog loop.ErrorLog,
) (srvs []job.ServiceCtx, err error) {
	var pluginConfig sm_plugin.PorOffchainConfig
	err = json.Unmarshal(jb.OCR2OracleSpec.PluginConfig.Bytes(), &pluginConfig)
	if err != nil {
		return
	}
	// TODO(gg): enable if validation exists
	// err = pluginConfig.Validate()
	// if err != nil {
	// 	return
	// }
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
		OracleSpecID:  *jb.OCR2OracleSpecID,
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
	srvs = append(srvs, provider)

	// TODO(gg): to be implemented when needed
	// secureMintProvider, ok := provider.(types.SecureMintProvider)
	// if !ok {
	// 	return nil, errors.New("could not coerce PluginProvider to SecureMintProvider")
	// }

	argsNoPlugin.ContractTransmitter = newStubContractTransmitter(lggr, ocr2plus_types.Account(spec.TransmitterID.String)) // TODO(gg): implement chain writing here
	argsNoPlugin.ContractConfigTracker = provider.ContractConfigTracker()
	argsNoPlugin.OffchainConfigDigester = provider.OffchainConfigDigester()

	abort := func() {
		if cerr := services.MultiCloser(srvs).Close(); err != nil {
			lggr.Errorw("Error closing unused services", "err", cerr)
		}
	}

	// dataSource := ocrcommon.NewDataSourceV2(pipelineRunner,
	// 	jb,
	// 	*jb.PipelineSpec,
	// 	lggr,
	// 	runSaver,
	// 	chEnhancedTelem)
	// lggr.Infof("Created data source %#v", dataSource)

	// juelsPerFeeCoinSource := ocrcommon.NewInMemoryDataSource(pipelineRunner, jb, pipeline.Spec{
	// 	ID:           jb.ID,
	// 	DotDagSource: pluginConfig.JuelsPerFeeCoinPipeline,
	// 	CreatedAt:    time.Now(),
	// }, lggr)

	// if pluginConfig.JuelsPerFeeCoinCache == nil || (pluginConfig.JuelsPerFeeCoinCache != nil && !pluginConfig.JuelsPerFeeCoinCache.Disable) {
	// 	lggr.Infof("juelsPerFeeCoin data source caching is enabled")
	// 	juelsPerFeeCoinSourceCache, err2 := ocrcommon.NewInMemoryDataSourceCache(juelsPerFeeCoinSource, kvStore, pluginConfig.JuelsPerFeeCoinCache)
	// 	if err2 != nil {
	// 		return nil, err2
	// 	}
	// 	juelsPerFeeCoinSource = juelsPerFeeCoinSourceCache
	// 	srvs = append(srvs, juelsPerFeeCoinSourceCache)
	// }

	// var gasPriceSubunitsDataSource libocr_median.DataSource
	// if pluginConfig.HasGasPriceSubunitsPipeline() {
	// 	gasPriceSubunitsDataSource = ocrcommon.NewInMemoryDataSource(pipelineRunner, jb, pipeline.Spec{
	// 		ID:           jb.ID,
	// 		DotDagSource: pluginConfig.GasPriceSubunitsPipeline,
	// 		CreatedAt:    time.Now(),
	// 	}, lggr)
	// } else {
	// 	gasPriceSubunitsDataSource = &median.ZeroDataSource{}
	// }

	if cmdName := env.SecureMintPlugin.Cmd.Get(); cmdName != "" {
		err = errors.New("loop for securemint plugin not implemented yet")
		abort()
		return
		// // use unique logger names so we can use it to register a loop
		// medianLggr := lggr.Named("Median").Named(spec.ContractID).Named(spec.GetID())
		// envVars, err2 := plugins.ParseEnvFile(env.MedianPlugin.Env.Get())
		// if err2 != nil {
		// 	err = fmt.Errorf("failed to parse median env file: %w", err2)
		// 	abort()
		// 	return
		// }
		// cmdFn, telem, err2 := cfg.RegisterLOOP(plugins.CmdConfig{
		// 	ID:  medianLggr.Name(),
		// 	Cmd: cmdName,
		// 	Env: envVars,
		// })
		// if err2 != nil {
		// 	err = fmt.Errorf("failed to register loop: %w", err2)
		// 	abort()
		// 	return
		// }
		// median := loop.NewMedianService(lggr, telem, cmdFn, medianProvider, spec.ContractID, dataSource, juelsPerFeeCoinSource, gasPriceSubunitsDataSource, errorLog, pluginConfig.DeviationFunctionDefinition)
		// argsNoPlugin.ReportingPluginFactory = median
		// srvs = append(srvs, median)
	} else {
		// TODO(gg): fill in params for the factory
		argsNoPlugin.ReportingPluginFactory = &sm_plugin.PorReportingPluginFactory{
			Logger:          argsNoPlugin.Logger,
			ExternalAdapter: sm_ea.NewExternalAdapter(pipelineRunner, jb, *jb.PipelineSpec, runSaver, lggr),
			ContractReader: &mockContractReader{
				// since we don't write to chain yet, we mock the contract reader which returns a the most recent config digest from the config contract
				getConfigDigestFunc: func() ([32]byte, error) {
					_, configDigest, err := argsNoPlugin.ContractConfigTracker.LatestConfigDetails(ctx)
					return configDigest, err
				},
			},
			ReportMarshaler: sm_plugin.NewMockReportMarshaler(),
			// ExternalAdapter: provider.ExternalAdapter(),
			// ContractReader:  provider.ContractReader(),
			// ReportMarshaler: provider.ReportMarshaler(),
		}
		if err != nil {
			err = fmt.Errorf("failed to create secure mint factory: %w", err)
			abort()
			return
		}
	}

	// TODO(gg): use promwrapper plugin to get ocr metrics?

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

// mockContractReader is a mock implementation of the ContractReader interface.
// It retrieves the latest config digest from the config contract and then uses that to return a mocked report.
// This is needed so that sm_plugin.ShouldTransmitAcceptedReport() does not fail (it checks the config digest).
type mockContractReader struct {
	getConfigDigestFunc func() ([32]byte, error)
}

func (m *mockContractReader) GetLatestTransmittedReportDetails(ctx context.Context, chainId por.ChainSelector) (sm_plugin.TransmittedReportDetails, error) {
	configDigest, err := m.getConfigDigestFunc()
	if err != nil {
		return sm_plugin.TransmittedReportDetails{}, fmt.Errorf("failed to get config digest: %w", err)
	}

	return sm_plugin.TransmittedReportDetails{
		ConfigDigest:    configDigest,
		SeqNr:           1,          // Mock sequence number
		LatestTimestamp: time.Now(), // Mock timestamp
	}, nil
}
