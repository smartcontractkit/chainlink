package mercury

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2plus"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	relaymercuryv1 "github.com/smartcontractkit/chainlink-data-streams/mercury/v1"
	relaymercuryv2 "github.com/smartcontractkit/chainlink-data-streams/mercury/v2"
	relaymercuryv3 "github.com/smartcontractkit/chainlink-data-streams/mercury/v3"

	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/mercury/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/utils"
	mercuryv1 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v1"
	mercuryv2 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v2"
	mercuryv3 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v3"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

type Config interface {
	MaxSuccessfulRuns() uint64
	ResultWriteQueueDepth() uint64
	plugins.RegistrarConfig
}

// concrete implementation of MercuryConfig
type mercuryConfig struct {
	jobPipelineMaxSuccessfulRuns     uint64
	jobPipelineResultWriteQueueDepth uint64
	plugins.RegistrarConfig
}

func NewMercuryConfig(jobPipelineMaxSuccessfulRuns uint64, jobPipelineResultWriteQueueDepth uint64, pluginProcessCfg plugins.RegistrarConfig) Config {
	return &mercuryConfig{
		jobPipelineMaxSuccessfulRuns:     jobPipelineMaxSuccessfulRuns,
		jobPipelineResultWriteQueueDepth: jobPipelineResultWriteQueueDepth,
		RegistrarConfig:                  pluginProcessCfg,
	}
}

func (m *mercuryConfig) MaxSuccessfulRuns() uint64 {
	return m.jobPipelineMaxSuccessfulRuns
}

func (m *mercuryConfig) ResultWriteQueueDepth() uint64 {
	return m.jobPipelineResultWriteQueueDepth
}

func NewServices(
	jb job.Job,
	ocr2Provider commontypes.MercuryProvider,
	pipelineRunner pipeline.Runner,
	lggr logger.Logger,
	argsNoPlugin libocr2.MercuryOracleArgs,
	cfg Config,
	chEnhancedTelem chan ocrcommon.EnhancedTelemetryMercuryData,
	orm types.DataSourceORM,
	feedID utils.FeedID,
) ([]job.ServiceCtx, error) {
	if jb.PipelineSpec == nil {
		return nil, errors.New("expected job to have a non-nil PipelineSpec")
	}

	var pluginConfig config.PluginConfig
	err := json.Unmarshal(jb.OCR2OracleSpec.PluginConfig.Bytes(), &pluginConfig)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = config.ValidatePluginConfig(pluginConfig, feedID)
	if err != nil {
		return nil, err
	}
	lggr = lggr.Named("MercuryPlugin").With("jobID", jb.ID, "jobName", jb.Name.ValueOrZero())

	// encapsulate all the subservices and ensure we close them all if any fail to start
	srvs := []job.ServiceCtx{ocr2Provider}
	abort := func() {
		if cerr := services.MultiCloser(srvs).Close(); err != nil {
			lggr.Errorw("Error closing unused services", "err", cerr)
		}
	}
	saver := ocrcommon.NewResultRunSaver(pipelineRunner, lggr, cfg.MaxSuccessfulRuns(), cfg.ResultWriteQueueDepth())
	srvs = append(srvs, saver)

	loopEnabled, loopCmd := env.MercuryPlugin.Cmd.Get() != "", env.MercuryPlugin.Env.Get()
	// this is the factory that will be used to create the mercury plugin
	var factory ocr3types.MercuryPluginFactory
	switch feedID.Version() {
	case 1:
		ds := mercuryv1.NewDataSource(
			orm,
			pipelineRunner,
			jb,
			*jb.PipelineSpec,
			lggr,
			saver,
			chEnhancedTelem,
			ocr2Provider.MercuryChainReader(),
			ocr2Provider.MercuryServerFetcher(),
			pluginConfig.InitialBlockNumber.Ptr(),
			feedID,
		)

		if loopEnabled {
			mercuryLggr := lggr.Named("MercuryV1").Named(feedID.String())
			envVars, err2 := plugins.ParseEnvFile(env.MercuryPlugin.Env.Get())
			if err2 != nil {
				abort()
				return nil, fmt.Errorf("failed to parse mercury env file: %w", err2)
			}
			cmdFn, opts, err2 := cfg.RegisterLOOP(plugins.CmdConfig{
				ID:  mercuryLggr.Name(),
				Cmd: loopCmd,
				Env: envVars,
			})
			if err2 != nil {
				abort()
				return nil, fmt.Errorf("failed to register loop: %w", err2)
			}
			// in loopp mode, the factory is grpc server, and we need to handle the server lifecycle
			factoryServer := loop.NewMercuryV1Service(lggr, opts, cmdFn, ocr2Provider, ds)
			srvs = append(srvs, factoryServer)
			// adapt the grpc server to the vanilla mercury plugin factory interface used by the oracle
			factory = factoryServer
		} else {
			factory = relaymercuryv1.NewFactory(ds, lggr, ocr2Provider.OnchainConfigCodec(), ocr2Provider.ReportCodecV1())
		}
	case 2:
		ds := mercuryv2.NewDataSource(
			orm,
			pipelineRunner,
			jb,
			*jb.PipelineSpec,
			feedID,
			lggr,
			saver,
			chEnhancedTelem,
			ocr2Provider.MercuryServerFetcher(),
			*pluginConfig.LinkFeedID,
			*pluginConfig.NativeFeedID,
		)

		if loopEnabled {
			mercuryLggr := lggr.Named("MercuryV2").Named(feedID.String())
			envVars, err2 := plugins.ParseEnvFile(env.MercuryPlugin.Env.Get())
			if err2 != nil {
				abort()
				return nil, fmt.Errorf("failed to parse mercury env file: %w", err2)
			}
			cmdFn, opts, err2 := cfg.RegisterLOOP(plugins.CmdConfig{
				ID:  mercuryLggr.Name(),
				Cmd: loopCmd,
				Env: envVars,
			})
			if err2 != nil {
				abort()
				return nil, fmt.Errorf("failed to register loop: %w", err2)
			}
			// in loopp mode, the factory is grpc server, and we need to handle the server lifecycle
			factoryServer := loop.NewMercuryV2Service(lggr, opts, cmdFn, ocr2Provider, ds)
			srvs = append(srvs, factoryServer)
			// adapt the grpc server to the vanilla mercury plugin factory interface used by the oracle
			factory = factoryServer
		} else {
			factory = relaymercuryv2.NewFactory(ds, lggr, ocr2Provider.OnchainConfigCodec(), ocr2Provider.ReportCodecV2())
		}
	case 3:
		ds := mercuryv3.NewDataSource(
			orm,
			pipelineRunner,
			jb,
			*jb.PipelineSpec,
			feedID,
			lggr,
			saver,
			chEnhancedTelem,
			ocr2Provider.MercuryServerFetcher(),
			*pluginConfig.LinkFeedID,
			*pluginConfig.NativeFeedID,
		)

		if loopEnabled {
			mercuryLggr := lggr.Named("MercuryV3").Named(feedID.String())
			envVars, err2 := plugins.ParseEnvFile(env.MercuryPlugin.Env.Get())
			if err2 != nil {
				abort()
				return nil, fmt.Errorf("failed to parse mercury env file: %w", err2)
			}
			cmdFn, opts, err2 := cfg.RegisterLOOP(plugins.CmdConfig{
				ID:  mercuryLggr.Name(),
				Cmd: loopCmd,
				Env: envVars,
			})
			if err2 != nil {
				abort()
				return nil, fmt.Errorf("failed to register loop: %w", err2)
			}
			// in loopp mode, the factory is grpc server, and we need to handle the server lifecycle
			factoryServer := loop.NewMercuryV3Service(lggr, opts, cmdFn, ocr2Provider, ds)
			srvs = append(srvs, factoryServer)
			// adapt the grpc server to the vanilla mercury plugin factory interface used by the oracle
			factory = factoryServer
		} else {
			factory = relaymercuryv3.NewFactory(ds, lggr, ocr2Provider.OnchainConfigCodec(), ocr2Provider.ReportCodecV3())
		}

	default:
		return nil, errors.Errorf("unknown Mercury report schema version: %d", feedID.Version())
	}
	argsNoPlugin.MercuryPluginFactory = factory
	oracle, err := libocr2.NewOracle(argsNoPlugin)
	if err != nil {
		abort()
		return nil, errors.WithStack(err)
	}
	srvs = append(srvs, job.NewServiceAdapter(oracle))
	return srvs, nil
}
