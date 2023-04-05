package mercury

import (
	"encoding/json"

	"github.com/pkg/errors"
	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2"

	relaymercury "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"
	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/mercury/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	mercury "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury"
)

type Config interface {
	JobPipelineMaxSuccessfulRuns() uint64
}

func NewServices(
	jb job.Job,
	ocr2Provider relaytypes.MercuryProvider,
	pipelineRunner pipeline.Runner,
	runResults chan pipeline.Run,
	lggr logger.Logger,
	argsNoPlugin libocr2.OracleArgs,
	cfg Config,
) ([]job.ServiceCtx, error) {
	if jb.PipelineSpec == nil {
		return nil, errors.New("expected job to have a non-nil PipelineSpec")
	}
	var pluginConfig config.PluginConfig
	err := json.Unmarshal(jb.OCR2OracleSpec.PluginConfig.Bytes(), &pluginConfig)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = config.ValidatePluginConfig(pluginConfig)
	if err != nil {
		return nil, err
	}
	lggr = lggr.Named("MercuryPlugin").With("jobID", jb.ID, "jobName", jb.Name.ValueOrZero())
	ds := mercury.NewDataSource(
		pipelineRunner,
		jb,
		*jb.PipelineSpec,
		lggr,
		runResults,
	)
	argsNoPlugin.ReportingPluginFactory = relaymercury.NewFactory(
		ds,
		lggr,
		ocr2Provider.OnchainConfigCodec(),
		ocr2Provider.ReportCodec(),
		ocr2Provider.ContractTransmitter(),
	)
	oracle, err := libocr2.NewOracle(argsNoPlugin)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return []job.ServiceCtx{ocr2Provider, ocrcommon.NewResultRunSaver(
		runResults,
		pipelineRunner,
		make(chan struct{}),
		lggr,
		cfg.JobPipelineMaxSuccessfulRuns(),
	),
		job.NewServiceAdapter(oracle)}, nil
}
