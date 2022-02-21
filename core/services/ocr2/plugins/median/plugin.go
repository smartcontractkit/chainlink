package median

import (
	"encoding/json"
	"time"

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

type Median struct {
	jobSpec        job.Job
	ocr2Provider   types.OCR2Provider
	pipelineRunner pipeline.Runner
	runResults     chan pipeline.Run
	lggr           logger.Logger
	ocrLogger      commontypes.Logger

	config PluginConfig
}

var _ plugins.OraclePlugin = Median{}

func NewMedian(jobSpec job.Job, ocr2Provider types.OCR2Provider, pipelineRunner pipeline.Runner, runResults chan pipeline.Run, lggr logger.Logger, ocrLogger commontypes.Logger) (Median, error) {
	var config PluginConfig
	err := json.Unmarshal(jobSpec.OCR2OracleSpec.PluginConfig.Bytes(), &config)
	if err != nil {
		return Median{}, err
	}
	err = validatePluginConfig(config)
	if err != nil {
		return Median{}, err
	}

	return Median{
		jobSpec:        jobSpec,
		ocr2Provider:   ocr2Provider,
		pipelineRunner: pipelineRunner,
		runResults:     runResults,
		lggr:           lggr,
		ocrLogger:      ocrLogger,
		config:         config,
	}, nil
}

func (m Median) GetPluginFactory() (plugin ocr2types.ReportingPluginFactory, err error) {
	juelsPerFeeCoinPipelineSpec := pipeline.Spec{
		ID:           m.jobSpec.ID,
		DotDagSource: m.config.JuelsPerFeeCoinPipeline,
		CreatedAt:    time.Now(),
	}
	numericalMedianFactory := median.NumericalMedianFactory{
		ContractTransmitter: m.ocr2Provider.MedianContract(),
		DataSource: ocrcommon.NewDataSourceV2(m.pipelineRunner,
			m.jobSpec,
			*m.jobSpec.PipelineSpec,
			m.lggr,
			m.runResults,
		),
		JuelsPerFeeCoinDataSource: ocrcommon.NewInMemoryDataSource(m.pipelineRunner, m.jobSpec, juelsPerFeeCoinPipelineSpec, m.lggr),
		ReportCodec:               m.ocr2Provider.ReportCodec(),
		Logger:                    m.ocrLogger,
	}
	return numericalMedianFactory, nil
}

func (m Median) GetServices() (services []job.Service, err error) {
	return []job.Service{}, nil
}
