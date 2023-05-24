package threshold

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/libocr/commontypes"
	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2"

	decryptionPlugin "github.com/smartcontractkit/tdh2/go/ocr2/decryptionplugin"
	decryptionPluginConfig "github.com/smartcontractkit/tdh2/go/ocr2/decryptionplugin/config"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type ThresholdServicesConfig struct {
	Job             job.Job
	PipelineRunner  pipeline.Runner
	JobORM          job.ORM
	OCR2JobConfig   validate.Config
	DB              *sqlx.DB
	Chain           evm.Chain
	ContractID      string
	Lggr            logger.Logger
	MailMon         *utils.MailboxMonitor
	URLsMonEndpoint commontypes.MonitoringEndpoint
	DecryptionQueue decryptionPlugin.DecryptionQueuingService
}

func NewThresholdService(sharedOracleArgs *libocr2.OracleArgs, conf *ThresholdServicesConfig) (job.ServiceCtx, error) {
	var pluginConfig decryptionPlugin.config.ReportingPluginConfig
	err := json.Unmarshal(conf.Job.OCR2OracleSpec.PluginConfig.Bytes(), &pluginConfig)
	if err != nil {
		return nil, err
	}
	err = config.ValidatePluginConfig(decryptionPlugin.pluginConfig)
	if err != nil {
		return nil, err
	}

	sharedOracleArgs.ReportingPluginFactory = decryptionPlugin.DecryptionReportingPluginFactory{
		DecryptionQueue: conf.DecryptionQueue,
		Logger:          sharedOracleArgs.Logger,
	}
	thresholdReportingPluginOracle, err := libocr2.NewOracle(*sharedOracleArgs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call NewOracle to create a Threshold Reporting Plugin")
	}

	return job.NewServiceAdapter(thresholdReportingPluginOracle), nil
}
