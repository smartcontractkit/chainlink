package median

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	mediantypes "github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"

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

// This wrapper avoids the need to modify the signature of NewMedianFactory in all of the non-evm
// relay repos as well as its primary definition in chainlink-common. Once ChainReader is implemented
// and working on all 4 blockchain families, we can remove the original MedianContract() method from
// MedianProvider and pass medianContract as a separate param to NewMedianFactory
type medianProviderWrapper struct {
	types.MedianProvider
	contract mediantypes.MedianContract
}

// Override relay's implementation of MedianContract with product plugin's implementation of
// MedianContract, making use of product-agnostic ChainReader to read the contract instead of relay MedianContract
func (m medianProviderWrapper) MedianContract() mediantypes.MedianContract {
	return m.contract
}

func NewMedianServices(ctx context.Context,
	jb job.Job,
	isNewlyCreatedJob bool,
	relayer loop.Relayer,
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
	err = config.ValidatePluginConfig(pluginConfig)
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

	dataSource, juelsPerFeeCoinSource := ocrcommon.NewDataSourceV2(pipelineRunner,
		jb,
		*jb.PipelineSpec,
		lggr,
		runSaver,
		chEnhancedTelem,
	), ocrcommon.NewInMemoryDataSource(pipelineRunner, jb, pipeline.Spec{
		ID:           jb.ID,
		DotDagSource: pluginConfig.JuelsPerFeeCoinPipeline,
		CreatedAt:    time.Now(),
	}, lggr)

	medianPluginCmd := env.MedianPluginCmd.Get()
	medianLoopEnabled := medianPluginCmd != ""

	// TODO BCF-2821 handle this properly as this blocks Solana chain reader dev
	if !medianLoopEnabled && medianProvider.ChainReader() != nil {
		lggr.Info("Chain Reader enabled")
		medianProvider = medianProviderWrapper{
			medianProvider, // attach newer MedianContract which uses ChainReader
			newMedianContract(provider.ChainReader(), common.HexToAddress(spec.ContractID)),
		}
	} else {
		lggr.Info("Chain Reader disabled")
	}

	if medianLoopEnabled {
		// use unique logger names so we can use it to register a loop
		medianLggr := lggr.Named("Median").Named(spec.ContractID).Named(spec.GetID())
		cmdFn, telem, err2 := cfg.RegisterLOOP(medianLggr.Name(), medianPluginCmd)
		if err2 != nil {
			err = fmt.Errorf("failed to register loop: %w", err2)
			abort()
			return
		}
		median := loop.NewMedianService(lggr, telem, cmdFn, medianProvider, dataSource, juelsPerFeeCoinSource, errorLog)
		argsNoPlugin.ReportingPluginFactory = median
		srvs = append(srvs, median)
	} else {
		argsNoPlugin.ReportingPluginFactory, err = median.NewPlugin(lggr).NewMedianFactory(ctx, medianProvider, dataSource, juelsPerFeeCoinSource, errorLog)
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

type medianContract struct {
	chainReader types.ChainReader
	contract    types.BoundContract
}

type latestTransmissionDetailsResponse struct {
	configDigest    ocr2types.ConfigDigest
	epoch           uint32
	round           uint8
	latestAnswer    *big.Int
	latestTimestamp time.Time
}

type latestRoundRequested struct {
	configDigest ocr2types.ConfigDigest
	epoch        uint32
	round        uint8
}

func (m *medianContract) LatestTransmissionDetails(ctx context.Context) (configDigest ocr2types.ConfigDigest, epoch uint32, round uint8, latestAnswer *big.Int, latestTimestamp time.Time, err error) {
	var resp latestTransmissionDetailsResponse

	err = m.chainReader.GetLatestValue(ctx, m.contract, "LatestTransmissionDetails", nil, &resp)
	if err != nil {
		return
	}

	return resp.configDigest, resp.epoch, resp.round, resp.latestAnswer, resp.latestTimestamp, err
}

func (m *medianContract) LatestRoundRequested(ctx context.Context, lookback time.Duration) (configDigest ocr2types.ConfigDigest, epoch uint32, round uint8, err error) {
	var resp latestRoundRequested

	err = m.chainReader.GetLatestValue(ctx, m.contract, "LatestRoundReported", map[string]string{}, &resp)
	if err != nil {
		return
	}

	return resp.configDigest, resp.epoch, resp.round, err
}

func newMedianContract(chainReader types.ChainReader, address common.Address) *medianContract {
	contract := types.BoundContract{Address: address.String(), Name: "median", Pending: true}
	return &medianContract{chainReader, contract}
}
