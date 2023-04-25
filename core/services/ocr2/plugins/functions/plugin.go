package functions

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"

	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/ocr2dr_oracle"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/functions"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type FunctionsServicesConfig struct {
	Job            job.Job
	PipelineRunner pipeline.Runner
	JobORM         job.ORM
	OCR2JobConfig  validate.Config
	DB             *sqlx.DB
	Chain          evm.Chain
	ContractID     string
	Lggr           logger.Logger
	MailMon        *utils.MailboxMonitor
}

// Create all OCR2 plugin Oracles and all extra services needed to run a Functions job.
func NewFunctionsServices(sharedOracleArgs *libocr2.OracleArgs, conf *FunctionsServicesConfig) ([]job.ServiceCtx, error) {
	pluginORM := functions.NewORM(conf.DB, conf.Lggr, conf.OCR2JobConfig, common.HexToAddress(conf.ContractID))

	var pluginConfig config.PluginConfig
	err := json.Unmarshal(conf.Job.OCR2OracleSpec.PluginConfig.Bytes(), &pluginConfig)
	if err != nil {
		return nil, err
	}
	err = config.ValidatePluginConfig(pluginConfig)
	if err != nil {
		return nil, err
	}

	contractAddress := common.HexToAddress(conf.Job.OCR2OracleSpec.ContractID)
	oracleContract, err := ocr2dr_oracle.NewOCR2DROracle(contractAddress, conf.Chain.Client())
	if err != nil {
		return nil, errors.Wrapf(err, "Functions: failed to create a FunctionsOracle wrapper for address: %v", contractAddress)
	}
	svcLogger := conf.Lggr.Named("FunctionsListener").
		With(
			"contract", contractAddress,
			"jobName", conf.Job.PipelineSpec.JobName,
			"jobID", conf.Job.PipelineSpec.JobID,
			"externalJobID", conf.Job.ExternalJobID,
		)
	functionsListener := functions.NewFunctionsListener(oracleContract, conf.Job, conf.PipelineRunner, conf.JobORM, pluginORM, pluginConfig, conf.Chain.LogBroadcaster(), svcLogger, conf.MailMon)

	sharedOracleArgs.ReportingPluginFactory = FunctionsReportingPluginFactory{
		Logger:    sharedOracleArgs.Logger,
		PluginORM: pluginORM,
		JobID:     conf.Job.ExternalJobID,
	}
	functionsReportingPluginOracle, err := libocr2.NewOracle(*sharedOracleArgs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call NewOracle to create a Functions Reporting Plugin")
	}

	return []job.ServiceCtx{job.NewServiceAdapter(functionsReportingPluginOracle), functionsListener}, nil
}
