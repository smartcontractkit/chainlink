package median

import (
	"encoding/json"
	"time"

	"github.com/smartcontractkit/chainlink-relay/pkg/types"
	"github.com/smartcontractkit/libocr/commontypes"
	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/median/config"
	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

// NewMedian parses the arguments and returns a new Median struct.
func NewMedianServices(jb job.Job,
	chainSet evm.ChainSet,
	ocr2Provider types.MedianProvider,
	pipelineRunner pipeline.Runner,
	runResults chan pipeline.Run,
	lggr logger.Logger,
	ocrLogger commontypes.Logger,
	new bool,
	argsNoPlugin libocr2.OracleArgs,
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
	//chainID := big.NewInt(0).SetUint64(pluginConfig.ChainID)
	//chain, err := chainSet.Get(chainID)
	//if err != nil {
	//	return nil, err
	//}
	juelsPerFeeCoinPipelineSpec := pipeline.Spec{
		ID:           jb.ID,
		DotDagSource: pluginConfig.JuelsPerFeeCoinPipeline,
		CreatedAt:    time.Now(),
	}
	argsNoPlugin.ReportingPluginFactory = median.NumericalMedianFactory{
		ContractTransmitter: ocr2Provider.MedianContract(),
		DataSource: ocrcommon.NewDataSourceV2(pipelineRunner,
			jb,
			*jb.PipelineSpec,
			lggr,
			runResults,
		),
		JuelsPerFeeCoinDataSource: ocrcommon.NewInMemoryDataSource(pipelineRunner, jb, juelsPerFeeCoinPipelineSpec, lggr),
		OnchainConfigCodec:        ocr2Provider.OnchainConfigCodec(),
		ReportCodec:               ocr2Provider.ReportCodec(),
		Logger:                    ocrLogger,
	}
	oracle, err := libocr2.NewOracle(argsNoPlugin)
	if err != nil {
		return nil, err
	}
	if new {
		// If this is a brand-new job, then we make use of the start blocks. If not then we're rebooting and log poller will pick up where we left off.
		// Median doesn't need specific fromBlocks before starting the plugin.
		return []job.ServiceCtx{ocr2Provider, ocrcommon.NewResultRunSaver(
			runResults,
			pipelineRunner,
			make(chan struct{}),
			lggr),
			ocrcommon.NewBackfilledOracle(lggr, []ocrcommon.BackfilledPoller{}, job.NewServiceAdapter(oracle))}, nil
	}
	return []job.ServiceCtx{}, nil
}
