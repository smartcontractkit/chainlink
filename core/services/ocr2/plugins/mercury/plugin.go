package mercury

import (
	"encoding/json"

	"github.com/pkg/errors"

	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2plus"

	relaymercuryv1 "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury/v1"
	relaymercuryv2 "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury/v2"
	relaymercuryv3 "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury/v3"
	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/mercury/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/utils"
	mercuryv1 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v1"
	mercuryv2 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v2"
	mercuryv3 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v3"
)

type Config interface {
	MaxSuccessfulRuns() uint64
}

func NewServices(
	jb job.Job,
	ocr2Provider relaytypes.MercuryProvider,
	pipelineRunner pipeline.Runner,
	runResults chan *pipeline.Run,
	lggr logger.Logger,
	argsNoPlugin libocr2.MercuryOracleArgs,
	cfg Config,
	chEnhancedTelem chan ocrcommon.EnhancedTelemetryMercuryData,
	chainHeadTracker types.ChainHeadTracker,
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

	switch feedID.Version() {
	case 1:
		ds := mercuryv1.NewDataSource(
			orm,
			pipelineRunner,
			jb,
			*jb.PipelineSpec,
			lggr,
			runResults,
			chEnhancedTelem,
			chainHeadTracker,
			ocr2Provider.MercuryServerFetcher(),
			pluginConfig.InitialBlockNumber.Ptr(),
			feedID,
		)
		argsNoPlugin.MercuryPluginFactory = relaymercuryv1.NewFactory(
			ds,
			lggr,
			ocr2Provider.OnchainConfigCodec(),
			ocr2Provider.ReportCodecV1(),
		)
	case 2:
		ds := mercuryv2.NewDataSource(
			orm,
			pipelineRunner,
			jb,
			*jb.PipelineSpec,
			feedID,
			lggr,
			runResults,
			chEnhancedTelem,
			ocr2Provider.MercuryServerFetcher(),
			*pluginConfig.LinkFeedID,
			*pluginConfig.NativeFeedID,
		)
		argsNoPlugin.MercuryPluginFactory = relaymercuryv2.NewFactory(
			ds,
			lggr,
			ocr2Provider.OnchainConfigCodec(),
			ocr2Provider.ReportCodecV2(),
		)
	case 3:
		ds := mercuryv3.NewDataSource(
			orm,
			pipelineRunner,
			jb,
			*jb.PipelineSpec,
			feedID,
			lggr,
			runResults,
			chEnhancedTelem,
			ocr2Provider.MercuryServerFetcher(),
			*pluginConfig.LinkFeedID,
			*pluginConfig.NativeFeedID,
		)
		argsNoPlugin.MercuryPluginFactory = relaymercuryv3.NewFactory(
			ds,
			lggr,
			ocr2Provider.OnchainConfigCodec(),
			ocr2Provider.ReportCodecV3(),
		)
	default:
		return nil, errors.Errorf("unknown Mercury report schema version: %d", feedID.Version())
	}

	oracle, err := libocr2.NewOracle(argsNoPlugin)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	saver := ocrcommon.NewResultRunSaver(runResults, pipelineRunner, make(chan struct{}), lggr, cfg.MaxSuccessfulRuns())
	return []job.ServiceCtx{ocr2Provider, saver, job.NewServiceAdapter(oracle)}, nil
}
