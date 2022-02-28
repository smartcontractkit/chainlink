package median

import (
	"encoding/json"
	"time"

	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/median/config"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins"
	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/relay/types"
)

// The Median struct holds parameters needed to run a Median plugin.
type Median struct {
	jb             job.Job
	ocr2Provider   types.OCR2ProviderCtx
	pipelineRunner pipeline.Runner
	runResults     chan pipeline.Run
	lggr           logger.Logger
	ocrLogger      commontypes.Logger

	pluginConfig config.PluginConfig
}

var _ plugins.OraclePlugin = &Median{}

// NewMedian parses the arguments and returns a new Median struct.
func NewMedian(jb job.Job, ocr2Provider types.OCR2ProviderCtx, pipelineRunner pipeline.Runner, runResults chan pipeline.Run, lggr logger.Logger, ocrLogger commontypes.Logger) (*Median, error) {
	var pluginConfig config.PluginConfig
	err := json.Unmarshal(jb.OCR2OracleSpec.PluginConfig.Bytes(), &pluginConfig)
	if err != nil {
		return &Median{}, err
	}
	err = config.ValidatePluginConfig(pluginConfig)
	if err != nil {
		return &Median{}, err
	}

	return &Median{
		jb:             jb,
		ocr2Provider:   ocr2Provider,
		pipelineRunner: pipelineRunner,
		runResults:     runResults,
		lggr:           lggr,
		ocrLogger:      ocrLogger,
		pluginConfig:   pluginConfig,
	}, nil
}

// GetPluginFactory return a median.NumericalMedianFactory.
func (m *Median) GetPluginFactory() (ocr2types.ReportingPluginFactory, error) {
	juelsPerFeeCoinPipelineSpec := pipeline.Spec{
		ID:           m.jb.ID,
		DotDagSource: m.pluginConfig.JuelsPerFeeCoinPipeline,
		CreatedAt:    time.Now(),
	}
	return median.NumericalMedianFactory{
		ContractTransmitter: m.ocr2Provider.MedianContract(),
		DataSource: ocrcommon.NewDataSourceV2(m.pipelineRunner,
			m.jb,
			*m.jb.PipelineSpec,
			m.lggr,
			m.runResults,
		),
		JuelsPerFeeCoinDataSource: ocrcommon.NewInMemoryDataSource(m.pipelineRunner, m.jb, juelsPerFeeCoinPipelineSpec, m.lggr),
		ReportCodec:               m.ocr2Provider.ReportCodec(),
		Logger:                    m.ocrLogger,
	}, nil
}

// GetServices return an empty Service slice because Median does not need any services besides the generic OCR2 ones
// supplied in the OCR2 delegate. This method exists to satisfy the plugins.OraclePlugin interface.
func (m *Median) GetServices() ([]job.ServiceCtx, error) {
	return []job.ServiceCtx{}, nil
}
