package directrequestocr

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/commontypes"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/ocr2dr_oracle"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/directrequestocr"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/directrequestocr/config"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

type DROracle struct {
	jb             job.Job
	pipelineRunner pipeline.Runner
	jobORM         job.ORM
	ocr2Provider   types.Plugin
	pluginConfig   config.PluginConfig
	pluginORM      directrequestocr.ORM
	chain          evm.Chain
	lggr           logger.Logger
	ocrLogger      commontypes.Logger
}

var _ plugins.OraclePlugin = &DROracle{}

func NewDROracle(jb job.Job, pipelineRunner pipeline.Runner, jobORM job.ORM, ocr2Provider types.Plugin, pluginORM directrequestocr.ORM, chain evm.Chain, lggr logger.Logger, ocrLogger commontypes.Logger) (*DROracle, error) {
	var pluginConfig config.PluginConfig
	err := json.Unmarshal(jb.OCR2OracleSpec.PluginConfig.Bytes(), &pluginConfig)
	if err != nil {
		return &DROracle{}, err
	}
	err = config.ValidatePluginConfig(pluginConfig)
	if err != nil {
		return &DROracle{}, err
	}

	return &DROracle{
		jb:             jb,
		pipelineRunner: pipelineRunner,
		jobORM:         jobORM,
		ocr2Provider:   ocr2Provider,
		pluginConfig:   pluginConfig,
		pluginORM:      pluginORM,
		chain:          chain,
		lggr:           lggr,
		ocrLogger:      ocrLogger,
	}, nil
}

func (o *DROracle) GetPluginFactory() (ocr2types.ReportingPluginFactory, error) {
	// TODO OCR reporting plugin: https://app.shortcut.com/chainlinklabs/story/54054/ocr-plugin-for-directrequest-ocr
	return nil, nil
}

func (o *DROracle) GetServices() ([]job.ServiceCtx, error) {
	contractAddress := common.HexToAddress(o.jb.OCR2OracleSpec.ContractID)
	oracle, err := ocr2dr_oracle.NewOCR2DROracle(contractAddress, o.chain.Client())
	if err != nil {
		return nil, errors.Wrapf(err, "OCR2DirectRequest: failed to create an OCR2DROracle wrapper for address: %v", contractAddress)
	}
	svcLogger := o.lggr.Named("DRListener").
		With(
			"contract", contractAddress,
			"jobName", o.jb.PipelineSpec.JobName,
			"jobID", o.jb.PipelineSpec.JobID,
			"externalJobID", o.jb.ExternalJobID,
		)
	logListener := directrequestocr.NewDRListener(oracle, o.jb, o.pipelineRunner, o.jobORM, o.pluginORM, o.pluginConfig, o.chain.LogBroadcaster(), svcLogger)
	var services []job.ServiceCtx
	services = append(services, logListener)
	return services, nil
}
