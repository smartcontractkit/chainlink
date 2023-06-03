package functions

import (
	"encoding/json"
	"math/big"
	"net/url"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"
	"golang.org/x/exp/slices"

	"github.com/smartcontractkit/libocr/commontypes"
	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/ocr2dr_oracle"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/functions"
	gw_common "github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type FunctionsServicesConfig struct {
	Job             job.Job
	JobORM          job.ORM
	BridgeORM       bridges.ORM
	OCR2JobConfig   validate.Config
	DB              *sqlx.DB
	Chain           evm.Chain
	ContractID      string
	Lggr            logger.Logger
	MailMon         *utils.MailboxMonitor
	URLsMonEndpoint commontypes.MonitoringEndpoint
	EthKeystore     keystore.Eth
}

const (
	FunctionsBridgeName     bridges.BridgeName = "ea_bridge"
	MaxAdapterResponseBytes int64              = 1_000_000
)

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

	allServices := []job.ServiceCtx{}
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

	bridge, err := conf.BridgeORM.FindBridge(FunctionsBridgeName)
	if err != nil {
		return nil, errors.Wrap(err, "Functions: unable to find bridge")
	}
	eaClient := functions.NewExternalAdapterClient(url.URL(bridge.URL), MaxAdapterResponseBytes)
	functionsListener := functions.NewFunctionsListener(oracleContract, conf.Job, eaClient, pluginORM, pluginConfig, conf.Chain.LogBroadcaster(), svcLogger, conf.MailMon, conf.URLsMonEndpoint)
	allServices = append(allServices, functionsListener)

	sharedOracleArgs.ReportingPluginFactory = FunctionsReportingPluginFactory{
		Logger:    sharedOracleArgs.Logger,
		PluginORM: pluginORM,
		JobID:     conf.Job.ExternalJobID,
	}
	functionsReportingPluginOracle, err := libocr2.NewOracle(*sharedOracleArgs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call NewOracle to create a Functions Reporting Plugin")
	}
	allServices = append(allServices, job.NewServiceAdapter(functionsReportingPluginOracle))

	if pluginConfig.GatewayConnectorConfig != nil {
		connectorLogger := conf.Lggr.Named("GatewayConnector").With("jobName", conf.Job.PipelineSpec.JobName)
		connector, err := NewConnector(pluginConfig.GatewayConnectorConfig, conf.EthKeystore, conf.Chain.ID(), connectorLogger)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create a GatewayConnector")
		}
		allServices = append(allServices, connector)
	}

	return allServices, nil
}

func NewConnector(gwcCfg *connector.ConnectorConfig, ethKeystore keystore.Eth, chainID *big.Int, lggr logger.Logger) (connector.GatewayConnector, error) {
	enabledKeys, err := ethKeystore.EnabledKeysForChain(chainID)
	if err != nil {
		return nil, err
	}
	configuredNodeAddress := common.HexToAddress(gwcCfg.NodeAddress)
	idx := slices.IndexFunc(enabledKeys, func(key ethkey.KeyV2) bool { return key.Address == configuredNodeAddress })
	if idx == -1 {
		return nil, errors.New("key for configured node address not found")
	}
	signerKey := enabledKeys[idx].ToEcdsaPrivKey()

	handler := functions.NewFunctionsConnectorHandler(signerKey, lggr)
	connector, err := connector.NewGatewayConnector(gwcCfg, handler, handler, gw_common.NewRealClock(), lggr)
	if err != nil {
		return nil, err
	}
	handler.SetConnector(connector)
	return connector, nil
}
