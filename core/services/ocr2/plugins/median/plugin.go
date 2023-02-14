package median

import (
	"context"
	"encoding/json"
	"time"

	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/median/config"
	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

type MedianConfig interface {
	JobPipelineMaxSuccessfulRuns() uint64
}

func NewMedianServices(ctx context.Context,
	jb job.Job,
	isNewlyCreatedJob bool,
	relayer loop.ChainRelayer,
	pipelineRunner pipeline.Runner,
	runResults chan pipeline.Run,
	lggr logger.Logger,
	argsNoPlugin libocr2.OracleArgs,
	cfg MedianConfig,
) ([]job.ServiceCtx, error) {
	var pluginConfig config.PluginConfig
	err := json.Unmarshal(jb.OCR2OracleSpec.PluginConfig.Bytes(), &pluginConfig)
	if err != nil {
		return nil, err
	}
	err = config.ValidatePluginConfig(pluginConfig)
	if err != nil {
		return nil, err
	}
	spec := jb.OCR2OracleSpec
	//TODO need to refresh this call on certain errors...push func() relayer down to this level
	// https://smartcontract-it.atlassian.net/browse/BCF-2112
	provider, err := relayer.NewMedianPluginProvider(ctx, types.RelayArgs{
		ExternalJobID: jb.ExternalJobID,
		JobID:         spec.ID,
		ContractID:    spec.ContractID,
		New:           isNewlyCreatedJob,
		RelayConfig:   spec.RelayConfig.Bytes(),
	}, types.PluginArgs{
		TransmitterID: spec.TransmitterID.String,
		PluginConfig:  spec.PluginConfig.Bytes(),
	}, ocrcommon.NewDataSourceV2(pipelineRunner,
		jb,
		*jb.PipelineSpec,
		lggr,
		runResults,
	), ocrcommon.NewInMemoryDataSource(pipelineRunner, jb, pipeline.Spec{
		ID:           jb.ID,
		DotDagSource: pluginConfig.JuelsPerFeeCoinPipeline,
		CreatedAt:    time.Now(),
	}, lggr))
	if err != nil {
		return nil, err
	}
	argsNoPlugin.ContractTransmitter = provider.ContractTransmitter()
	argsNoPlugin.ContractConfigTracker = provider.ContractConfigTracker()
	argsNoPlugin.OffchainConfigDigester = provider.OffchainConfigDigester()
	argsNoPlugin.ReportingPluginFactory = provider
	oracle, err := libocr2.NewOracle(argsNoPlugin)
	if err != nil {
		return nil, err
	}
	runSaver := ocrcommon.NewResultRunSaver(
		runResults,
		pipelineRunner,
		make(chan struct{}),
		lggr,
		cfg.JobPipelineMaxSuccessfulRuns(),
	)
	return []job.ServiceCtx{provider, runSaver, job.NewServiceAdapter(oracle)}, nil
}
