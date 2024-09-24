package mercury

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/pkg/errors"

	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2plus"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	relaymercuryv1 "github.com/smartcontractkit/chainlink-data-streams/mercury/v1"
	relaymercuryv2 "github.com/smartcontractkit/chainlink-data-streams/mercury/v2"
	relaymercuryv3 "github.com/smartcontractkit/chainlink-data-streams/mercury/v3"
	relaymercuryv4 "github.com/smartcontractkit/chainlink-data-streams/mercury/v4"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

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
	mercuryv4 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v4"
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
	enableTriggerCapability bool,
) ([]job.ServiceCtx, error) {
	if jb.PipelineSpec == nil {
		return nil, errors.New("expected job to have a non-nil PipelineSpec")
	}

	var err error
	var pluginConfig config.PluginConfig
	if len(jb.OCR2OracleSpec.PluginConfig) == 0 {
		if !enableTriggerCapability {
			return nil, fmt.Errorf("at least one transmission option must be configured")
		}
	} else {
		err = json.Unmarshal(jb.OCR2OracleSpec.PluginConfig.Bytes(), &pluginConfig)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		err = config.ValidatePluginConfig(pluginConfig, feedID)
		if err != nil {
			return nil, err
		}
	}

	lggr = lggr.Named("MercuryPlugin").With("jobID", jb.ID, "jobName", jb.Name.ValueOrZero())

	// encapsulate all the subservices and ensure we close them all if any fail to start
	srvs := []job.ServiceCtx{ocr2Provider}
	abort := func() {
		if err = services.MultiCloser(srvs).Close(); err != nil {
			lggr.Errorw("Error closing unused services", "err", err)
		}
	}
	saver := ocrcommon.NewResultRunSaver(pipelineRunner, lggr, cfg.MaxSuccessfulRuns(), cfg.ResultWriteQueueDepth())
	srvs = append(srvs, saver)

	// this is the factory that will be used to create the mercury plugin
	var (
		factory         ocr3types.MercuryPluginFactory
		factoryServices []job.ServiceCtx
	)
	fCfg := factoryCfg{
		orm:                   orm,
		pipelineRunner:        pipelineRunner,
		jb:                    jb,
		lggr:                  lggr,
		saver:                 saver,
		chEnhancedTelem:       chEnhancedTelem,
		ocr2Provider:          ocr2Provider,
		reportingPluginConfig: pluginConfig,
		cfg:                   cfg,
		feedID:                feedID,
	}
	switch feedID.Version() {
	case 1:
		factory, factoryServices, err = newv1factory(fCfg)
		if err != nil {
			abort()
			return nil, fmt.Errorf("failed to create mercury v1 factory: %w", err)
		}
		srvs = append(srvs, factoryServices...)
	case 2:
		factory, factoryServices, err = newv2factory(fCfg)
		if err != nil {
			abort()
			return nil, fmt.Errorf("failed to create mercury v2 factory: %w", err)
		}
		srvs = append(srvs, factoryServices...)
	case 3:
		factory, factoryServices, err = newv3factory(fCfg)
		if err != nil {
			abort()
			return nil, fmt.Errorf("failed to create mercury v3 factory: %w", err)
		}
		srvs = append(srvs, factoryServices...)
	case 4:
		factory, factoryServices, err = newv4factory(fCfg)
		if err != nil {
			abort()
			return nil, fmt.Errorf("failed to create mercury v4 factory: %w", err)
		}
		srvs = append(srvs, factoryServices...)
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

type factoryCfg struct {
	orm                   types.DataSourceORM
	pipelineRunner        pipeline.Runner
	jb                    job.Job
	lggr                  logger.Logger
	saver                 *ocrcommon.RunResultSaver
	chEnhancedTelem       chan ocrcommon.EnhancedTelemetryMercuryData
	ocr2Provider          commontypes.MercuryProvider
	reportingPluginConfig config.PluginConfig
	cfg                   Config
	feedID                utils.FeedID
}

func getPluginFeedIDs(pluginConfig config.PluginConfig) (linkFeedID utils.FeedID, nativeFeedID utils.FeedID) {
	if pluginConfig.LinkFeedID != nil {
		linkFeedID = *pluginConfig.LinkFeedID
	}
	if pluginConfig.NativeFeedID != nil {
		nativeFeedID = *pluginConfig.NativeFeedID
	}
	return linkFeedID, nativeFeedID
}

func newv4factory(factoryCfg factoryCfg) (ocr3types.MercuryPluginFactory, []job.ServiceCtx, error) {
	var factory ocr3types.MercuryPluginFactory
	srvs := make([]job.ServiceCtx, 0)

	linkFeedID, nativeFeedID := getPluginFeedIDs(factoryCfg.reportingPluginConfig)

	ds := mercuryv4.NewDataSource(
		factoryCfg.orm,
		factoryCfg.pipelineRunner,
		factoryCfg.jb,
		*factoryCfg.jb.PipelineSpec,
		factoryCfg.feedID,
		factoryCfg.lggr,
		factoryCfg.saver,
		factoryCfg.chEnhancedTelem,
		factoryCfg.ocr2Provider.MercuryServerFetcher(),
		linkFeedID,
		nativeFeedID,
	)

	loopCmd := env.MercuryPlugin.Cmd.Get()
	loopEnabled := loopCmd != ""

	if loopEnabled {
		cmdFn, opts, mercuryLggr, err := initLoop(loopCmd, factoryCfg.cfg, factoryCfg.feedID, factoryCfg.lggr)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to init loop for feed %s: %w", factoryCfg.feedID, err)
		}
		// in loop mode, the factory is grpc server, and we need to handle the server lifecycle
		factoryServer := loop.NewMercuryV4Service(mercuryLggr, opts, cmdFn, factoryCfg.ocr2Provider, ds)
		srvs = append(srvs, factoryServer)
		// adapt the grpc server to the vanilla mercury plugin factory interface used by the oracle
		factory = factoryServer
	} else {
		factory = relaymercuryv4.NewFactory(ds, factoryCfg.lggr, factoryCfg.ocr2Provider.OnchainConfigCodec(), factoryCfg.ocr2Provider.ReportCodecV4())
	}
	return factory, srvs, nil
}

func newv3factory(factoryCfg factoryCfg) (ocr3types.MercuryPluginFactory, []job.ServiceCtx, error) {
	var factory ocr3types.MercuryPluginFactory
	srvs := make([]job.ServiceCtx, 0)

	linkFeedID, nativeFeedID := getPluginFeedIDs(factoryCfg.reportingPluginConfig)

	ds := mercuryv3.NewDataSource(
		factoryCfg.orm,
		factoryCfg.pipelineRunner,
		factoryCfg.jb,
		*factoryCfg.jb.PipelineSpec,
		factoryCfg.feedID,
		factoryCfg.lggr,
		factoryCfg.saver,
		factoryCfg.chEnhancedTelem,
		factoryCfg.ocr2Provider.MercuryServerFetcher(),
		linkFeedID,
		nativeFeedID,
	)

	loopCmd := env.MercuryPlugin.Cmd.Get()
	loopEnabled := loopCmd != ""

	if loopEnabled {
		cmdFn, opts, mercuryLggr, err := initLoop(loopCmd, factoryCfg.cfg, factoryCfg.feedID, factoryCfg.lggr)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to init loop for feed %s: %w", factoryCfg.feedID, err)
		}
		// in loopp mode, the factory is grpc server, and we need to handle the server lifecycle
		factoryServer := loop.NewMercuryV3Service(mercuryLggr, opts, cmdFn, factoryCfg.ocr2Provider, ds)
		srvs = append(srvs, factoryServer)
		// adapt the grpc server to the vanilla mercury plugin factory interface used by the oracle
		factory = factoryServer
	} else {
		factory = relaymercuryv3.NewFactory(ds, factoryCfg.lggr, factoryCfg.ocr2Provider.OnchainConfigCodec(), factoryCfg.ocr2Provider.ReportCodecV3())
	}
	return factory, srvs, nil
}

func newv2factory(factoryCfg factoryCfg) (ocr3types.MercuryPluginFactory, []job.ServiceCtx, error) {
	var factory ocr3types.MercuryPluginFactory
	srvs := make([]job.ServiceCtx, 0)

	linkFeedID, nativeFeedID := getPluginFeedIDs(factoryCfg.reportingPluginConfig)

	ds := mercuryv2.NewDataSource(
		factoryCfg.orm,
		factoryCfg.pipelineRunner,
		factoryCfg.jb,
		*factoryCfg.jb.PipelineSpec,
		factoryCfg.feedID,
		factoryCfg.lggr,
		factoryCfg.saver,
		factoryCfg.chEnhancedTelem,
		factoryCfg.ocr2Provider.MercuryServerFetcher(),
		linkFeedID,
		nativeFeedID,
	)

	loopCmd := env.MercuryPlugin.Cmd.Get()
	loopEnabled := loopCmd != ""

	if loopEnabled {
		cmdFn, opts, mercuryLggr, err := initLoop(loopCmd, factoryCfg.cfg, factoryCfg.feedID, factoryCfg.lggr)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to init loop for feed %s: %w", factoryCfg.feedID, err)
		}
		// in loopp mode, the factory is grpc server, and we need to handle the server lifecycle
		factoryServer := loop.NewMercuryV2Service(mercuryLggr, opts, cmdFn, factoryCfg.ocr2Provider, ds)
		srvs = append(srvs, factoryServer)
		// adapt the grpc server to the vanilla mercury plugin factory interface used by the oracle
		factory = factoryServer
	} else {
		factory = relaymercuryv2.NewFactory(ds, factoryCfg.lggr, factoryCfg.ocr2Provider.OnchainConfigCodec(), factoryCfg.ocr2Provider.ReportCodecV2())
	}
	return factory, srvs, nil
}

func newv1factory(factoryCfg factoryCfg) (ocr3types.MercuryPluginFactory, []job.ServiceCtx, error) {
	var factory ocr3types.MercuryPluginFactory
	srvs := make([]job.ServiceCtx, 0)

	ds := mercuryv1.NewDataSource(
		factoryCfg.orm,
		factoryCfg.pipelineRunner,
		factoryCfg.jb,
		*factoryCfg.jb.PipelineSpec,
		factoryCfg.lggr,
		factoryCfg.saver,
		factoryCfg.chEnhancedTelem,
		factoryCfg.ocr2Provider.MercuryChainReader(),
		factoryCfg.ocr2Provider.MercuryServerFetcher(),
		factoryCfg.reportingPluginConfig.InitialBlockNumber.Ptr(),
		factoryCfg.feedID,
	)

	loopCmd := env.MercuryPlugin.Cmd.Get()
	loopEnabled := loopCmd != ""

	if loopEnabled {
		cmdFn, opts, mercuryLggr, err := initLoop(loopCmd, factoryCfg.cfg, factoryCfg.feedID, factoryCfg.lggr)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to init loop for feed %s: %w", factoryCfg.feedID, err)
		}
		// in loopp mode, the factory is grpc server, and we need to handle the server lifecycle
		factoryServer := loop.NewMercuryV1Service(mercuryLggr, opts, cmdFn, factoryCfg.ocr2Provider, ds)
		srvs = append(srvs, factoryServer)
		// adapt the grpc server to the vanilla mercury plugin factory interface used by the oracle
		factory = factoryServer
	} else {
		factory = relaymercuryv1.NewFactory(ds, factoryCfg.lggr, factoryCfg.ocr2Provider.OnchainConfigCodec(), factoryCfg.ocr2Provider.ReportCodecV1())
	}
	return factory, srvs, nil
}

func initLoop(cmd string, cfg Config, feedID utils.FeedID, lggr logger.Logger) (func() *exec.Cmd, loop.GRPCOpts, logger.Logger, error) {
	lggr.Debugw("Initializing Mercury loop", "command", cmd)
	mercuryLggr := lggr.Named(fmt.Sprintf("MercuryV%d", feedID.Version())).Named(feedID.String())
	envVars, err := plugins.ParseEnvFile(env.MercuryPlugin.Env.Get())
	if err != nil {
		return nil, loop.GRPCOpts{}, nil, fmt.Errorf("failed to parse mercury env file: %w", err)
	}
	cmdFn, opts, err := cfg.RegisterLOOP(plugins.CmdConfig{
		ID:  mercuryLggr.Name(),
		Cmd: cmd,
		Env: envVars,
	})
	if err != nil {
		return nil, loop.GRPCOpts{}, nil, fmt.Errorf("failed to register loop: %w", err)
	}
	return cmdFn, opts, mercuryLggr, nil
}
