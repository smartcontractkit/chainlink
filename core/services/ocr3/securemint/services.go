package securemint

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr3/promwrapper"
	sm_plugin_config "github.com/smartcontractkit/chainlink/v2/core/services/ocr3/securemint/config"
	sm_ea "github.com/smartcontractkit/chainlink/v2/core/services/ocr3/securemint/ea"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/plugins"
	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus"
	ocr2plus_types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/por_mock_ocr3plugin/por"
	sm_plugin "github.com/smartcontractkit/por_mock_ocr3plugin/por"
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

// TODO(gg): add unit tests to chain reading & writing
// NewSecureMintServices creates all securemint plugin specific services.
func NewSecureMintServices(ctx context.Context,
	jb job.Job,
	isNewlyCreatedJob bool,
	relayer loop.Relayer,
	pipelineRunner pipeline.Runner,
	lggr logger.Logger,
	argsNoPlugin libocr.OCR3OracleArgs[por.ChainSelector],
	cfg JobConfig,
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

	rawForwarderAddress, ok := spec.RelayConfig["keystoneForwarderAddress"]
	if !ok {
		return nil, fmt.Errorf("keystoneForwarderAddress not found in plugin config")
	}

	// Validate that the address is a string (for future use if needed)
	forwarderAddress, ok := rawForwarderAddress.(string)
	if !ok {
		return nil, fmt.Errorf("keystoneForwarderAddress in plugin config must be a string, got %T", rawForwarderAddress)
	}

	// Extract contract writer from the relayer
	// evmRelayer, ok := relayer.(interface {
	// 	GetChain() interface {
	// 		ContractWriter() types.ContractWriter
	// 	}
	// })
	// if !ok {
	// 	return nil, fmt.Errorf("relayer does not support EVM operations")
	// }

	// contractWriter := evmRelayer.GetChain().ContractWriter()
	// if contractWriter == nil {
	// 	return nil, fmt.Errorf("failed to get contract writer from relayer")
	// }

	// Get the eth key for the transmitter
	rid, err := spec.RelayID()
	if err != nil {
		return nil, fmt.Errorf("failed to get relay ID: %w", err)
	}

	chainID, err := strconv.ParseUint(rid.ChainID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse chain ID %s as uint64", rid.ChainID)
	}

	// TODO: These are placeholders. In a real implementation, extract these from config or relayer.
	receiverAddress := "" // Set to actual DF Cache receiver address as string
	fromAccount := ocr2plus_types.Account(spec.TransmitterID.String)
	chainSelector := por.ChainSelector(chainID)

	// r := relayer.(evm.Relayer)
	// legacyevm.NewLegacyChains(relayer)
	// evm.NewWriteTarget(ctx, nil, chain, *wCfg.GasLimitDefault(), lggr)
	// transmitter, err := evm.NewContractTransmitter(lggr, forwarderAddress, receiverAddress, fromAccount, chainSelector, backend)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create contract transmitter: %w", err)
	// }

	contractWriter, err := relayer.NewContractWriter(ctx, spec.RelayConfig.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to create contract writer: %w", err)
	}

	keystoneForwarderTransmitter, err := createKeystoneForwarderContractTransmitter(lggr, contractWriter, fromAccount, forwarderAddress, receiverAddress, chainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to create keystone forwarder transmitter: %w", err)
	}
	argsNoPlugin.ContractTransmitter = keystoneForwarderTransmitter

	abort := func() {
		if cerr := services.MultiCloser(srvs).Close(); cerr != nil {
			lggr.Errorw("Error closing services", "err", cerr)
		}
	}

	// Create the reporting plugin factory
	if cmdName := env.SecureMintPlugin.Cmd.Get(); cmdName != "" {
		abort()
		return nil, errors.New("LOOPP for securemint plugin not implemented yet")
	}

	// Create the original SecureMint plugin factory
	smPluginFactory := &sm_plugin.PorReportingPluginFactory{
		Logger:          argsNoPlugin.Logger,
		ExternalAdapter: sm_ea.NewExternalAdapter(secureMintPluginConfig, pipelineRunner, jb, *jb.PipelineSpec, runSaver, lggr),
		ContractReader:  newStubContractReader(argsNoPlugin.ContractConfigTracker), // since we don't write to chain yet, we mock the contract reader which returns the most recent config digest from the config contract
		ReportMarshaler: sm_plugin.NewMockReportMarshaler(),
	}

	// Wrap the factory with prometheus metrics monitoring
	argsNoPlugin.ReportingPluginFactory = promwrapper.NewReportingPluginFactory(
		smPluginFactory,
		lggr,
		rid.ChainID,
		"secure-mint",
	)

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
