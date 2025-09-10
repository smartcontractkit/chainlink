package securemint

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	coretypes "github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core/securemint"
	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr3/promwrapper"
	sm_plugin_config "github.com/smartcontractkit/chainlink/v2/core/services/ocr3/securemint/config"
	sm_ea "github.com/smartcontractkit/chainlink/v2/core/services/ocr3/securemint/ea"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	evm_types "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	"github.com/smartcontractkit/chainlink/v2/plugins"
	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus"
)

var _ JobConfig = (*smJobConfig)(nil)

type JobConfig interface {
	JobPipelineMaxSuccessfulRuns() uint64
	JobPipelineResultWriteQueueDepth() uint64
	plugins.RegistrarConfig
}

// concrete implementation of JobConfig
type smJobConfig struct {
	jobPipelineMaxSuccessfulRuns     uint64
	jobPipelineResultWriteQueueDepth uint64
	plugins.RegistrarConfig
}

func NewJobConfig(jobPipelineMaxSuccessfulRuns uint64, jobPipelineResultWriteQueueDepth uint64, pluginProcessCfg plugins.RegistrarConfig) JobConfig {
	return &smJobConfig{
		jobPipelineMaxSuccessfulRuns:     jobPipelineMaxSuccessfulRuns,
		jobPipelineResultWriteQueueDepth: jobPipelineResultWriteQueueDepth,
		RegistrarConfig:                  pluginProcessCfg,
	}
}

func (m *smJobConfig) JobPipelineMaxSuccessfulRuns() uint64 {
	return m.jobPipelineMaxSuccessfulRuns
}

func (m *smJobConfig) JobPipelineResultWriteQueueDepth() uint64 {
	return m.jobPipelineResultWriteQueueDepth
}

// XXX_SingletonTransmitter is a hack to allow the secure mint integration test to access the transmitter in order to verify the sent reports.
var XXX_SingletonTransmitter atomic.Value // capabilities.TriggerCapability

// NewSecureMintServices creates all securemint plugin specific services.
func NewSecureMintServices(ctx context.Context,
	jb job.Job,
	isNewlyCreatedJob bool,
	relayer loop.Relayer,
	pipelineRunner pipeline.Runner,
	lggr logger.Logger,
	argsNoPlugin libocr.OCR3OracleArgs[securemint.ChainSelector],
	cfg JobConfig,
	capabilitiesRegistry coretypes.CapabilitiesRegistry,
) (srvs []job.ServiceCtx, err error) {
	// Parse and validate the secure mint plugin configuration
	secureMintPluginConfig, err := sm_plugin_config.Parse(jb.OCR2OracleSpec.PluginConfig.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to parse secure mint plugin config: %w", err)
	}

	if err = secureMintPluginConfig.Validate(); err != nil {
		return nil, fmt.Errorf("invalid secure mint plugin config: %#v, %w", secureMintPluginConfig, err)
	}

	spec := jb.OCR2OracleSpec

	// Get relay config to extract don ID
	relayConfig, err := evm_types.NewRelayOpts(types.RelayArgs{
		ExternalJobID: jb.ExternalJobID,
		JobID:         jb.ID,
		ContractID:    spec.ContractID,
		New:           isNewlyCreatedJob,
		RelayConfig:   spec.RelayConfig.Bytes(),
		ProviderType:  string(spec.PluginType),
	}).RelayConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get relay config: %w", err)
	}

	// Create result run saver for pipeline execution
	runSaver := ocrcommon.NewResultRunSaver(
		pipelineRunner,
		lggr,
		cfg.JobPipelineMaxSuccessfulRuns(),
		cfg.JobPipelineResultWriteQueueDepth(),
	)

	configProvider, err := relayer.NewConfigProvider(ctx, types.RelayArgs{
		ExternalJobID: jb.ExternalJobID,
		JobID:         jb.ID,
		OracleSpecID:  *jb.OCR2OracleSpecID,
		ContractID:    spec.ContractID,
		New:           isNewlyCreatedJob,
		RelayConfig:   spec.RelayConfig.Bytes(),
		ProviderType:  string(spec.PluginType),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create config provider: %w", err)
	}
	srvs = append(srvs, configProvider)

	argsNoPlugin.ContractConfigTracker = configProvider.ContractConfigTracker()
	argsNoPlugin.OffchainConfigDigester = configProvider.OffchainConfigDigester()

	// Create the new secure mint transmitter with trigger capabilities
	transmitterConfig := TransmitterConfig{
		Logger:                       lggr,
		CapabilitiesRegistry:         capabilitiesRegistry,
		DonID:                        relayConfig.LLODONID,
		TriggerCapabilityName:        secureMintPluginConfig.TriggerCapabilityName,
		TriggerCapabilityVersion:     secureMintPluginConfig.TriggerCapabilityVersion,
		TriggerTickerMinResolutionMs: secureMintPluginConfig.TriggerTickerMinResolutionMs,
		TriggerSendChannelBufferSize: secureMintPluginConfig.TriggerSendChannelBufferSize,
	}

	transmitter, err := transmitterConfig.NewTransmitter(spec.TransmitterID.String)
	if err != nil {
		return nil, fmt.Errorf("failed to create secure mint transmitter: %w", err)
	}
	srvs = append(srvs, transmitter)
	argsNoPlugin.ContractTransmitter = transmitter
	XXX_SingletonTransmitter.Store(transmitter)

	abort := func() {
		if cerr := services.MultiCloser(srvs).Close(); cerr != nil {
			lggr.Errorw("Error closing services", "err", cerr)
		}
	}

	// Create the reporting plugin factory
	cmdName := env.SecureMintPlugin.Cmd.Get()
	if cmdName == "" {
		abort()
		return nil, errors.New("secure mint plugin loop is not configured, non-loopp mode is not supported for secure mint")
	}
	lggr.Infof("Configuration indicates loopp usage for secure mint")

	// use unique logger names so we can use it to register a loop
	secureMintLggr := lggr.Named("SecureMint").Named(spec.ContractID).Named(spec.GetID())
	envVars, err := plugins.ParseEnvFile(env.SecureMintPlugin.Env.Get())
	if err != nil {
		err = fmt.Errorf("failed to parse secure mint env file: %w", err)
		abort()
		return
	}
	cmdFn, telem, err := cfg.RegisterLOOP(plugins.CmdConfig{
		ID:  secureMintLggr.Name(),
		Cmd: cmdName,
		Env: envVars,
	})
	if err != nil {
		err = fmt.Errorf("failed to register loop: %w", err)
		abort()
		return
	}

	ea, err := sm_ea.NewExternalAdapter(secureMintPluginConfig, pipelineRunner, jb, *jb.PipelineSpec, runSaver, lggr)
	if err != nil {
		return nil, fmt.Errorf("failed to create secure mint external adapter: %w", err)
	}

	secureMintPluginFactory := loop.NewPluginSecureMintService(lggr, telem, cmdFn, ea)
	srvs = append(srvs, secureMintPluginFactory)

	// Wrap the factory with prometheus metrics monitoring
	promPluginFactory := promwrapper.NewReportingPluginFactory(
		secureMintPluginFactory,
		lggr,
		"",
		spec.ChainID,
		"secure-mint",
	)
	argsNoPlugin.ReportingPluginFactory = promPluginFactory

	// Create the oracle
	var oracle libocr.Oracle
	oracle, err = libocr.NewOracle(argsNoPlugin)
	if err != nil {
		abort()
		return nil, fmt.Errorf("failed to create oracle: %w", err)
	}

	// Assemble all services
	srvs = append(srvs, runSaver, job.NewServiceAdapter(oracle))

	return srvs, nil
}
