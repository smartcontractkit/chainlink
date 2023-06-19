package functions

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"
	"golang.org/x/exp/slices"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2plus"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/ocr2dr_oracle"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/functions"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/s4"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	s4_orm "github.com/smartcontractkit/chainlink/v2/core/services/s4"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type FunctionsServicesConfig struct {
	Job             job.Job
	JobORM          job.ORM
	BridgeORM       bridges.ORM
	QConfig         pg.QConfig
	DB              *sqlx.DB
	Chain           evm.Chain
	ContractID      string
	Lggr            logger.Logger
	MailMon         *utils.MailboxMonitor
	URLsMonEndpoint commontypes.MonitoringEndpoint
	EthKeystore     keystore.Eth
}

const (
	FunctionsBridgeName     string = "ea_bridge"
	S4Namespace             string = "functions"
	MaxAdapterResponseBytes int64  = 1_000_000
)

// Create all OCR2 plugin Oracles and all extra services needed to run a Functions job.
func NewFunctionsServices(sharedOracleArgs *libocr2.OCR2OracleArgs, conf *FunctionsServicesConfig) ([]job.ServiceCtx, error) {
	pluginORM := functions.NewORM(conf.DB, conf.Lggr, conf.QConfig, common.HexToAddress(conf.ContractID))
	s4ORM := s4_orm.NewPostgresORM(conf.DB, conf.Lggr, conf.QConfig, s4_orm.SharedTableName, S4Namespace)

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
	listenerLogger := conf.Lggr.Named("FunctionsListener")
	bridgeAccessor := functions.NewBridgeAccessor(conf.BridgeORM, FunctionsBridgeName, MaxAdapterResponseBytes)
	functionsListener := functions.NewFunctionsListener(oracleContract, conf.Job, bridgeAccessor, pluginORM, pluginConfig, conf.Chain.LogBroadcaster(), listenerLogger, conf.MailMon, conf.URLsMonEndpoint)
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

	s4OracleArgs := *sharedOracleArgs
	s4OracleArgs.ReportingPluginFactory = s4.S4ReportingPluginFactory{
		Logger: sharedOracleArgs.Logger,
		ORM:    s4ORM,
		ConfigDecoder: func(configBytes []byte) (*s4.PluginConfig, *types.ReportingPluginLimits, error) {
			wrapper, err2 := config.DecodeReportingPluginConfig(configBytes)
			if err != nil {
				return nil, nil, err2
			}
			if wrapper.Config == nil || wrapper.Config.S4PluginConfig == nil {
				return nil, nil, errors.New("missing s4 plugin config")
			}
			return &s4.PluginConfig{
					ProductName:             "functions",
					NSnapshotShards:         uint(wrapper.Config.S4PluginConfig.NSnapshotShards),
					MaxObservationEntries:   uint(wrapper.Config.S4PluginConfig.MaxObservationEntries),
					MaxReportEntries:        uint(wrapper.Config.S4PluginConfig.MaxReportEntries),
					MaxDeleteExpiredEntries: uint(wrapper.Config.S4PluginConfig.MaxDeleteExpiredEntries),
				},
				&types.ReportingPluginLimits{
					MaxQueryLength:       int(wrapper.Config.S4PluginConfig.MaxQueryLengthBytes),
					MaxObservationLength: int(wrapper.Config.S4PluginConfig.MaxObservationLengthBytes),
					MaxReportLength:      int(wrapper.Config.S4PluginConfig.MaxReportLengthBytes),
				}, nil
		},
	}
	s4ReportingPluginOracle, err := libocr2.NewOracle(s4OracleArgs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call NewOracle to create a S4 Reporting Plugin")
	}
	allServices = append(allServices, job.NewServiceAdapter(s4ReportingPluginOracle))

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
	connector, err := connector.NewGatewayConnector(gwcCfg, handler, handler, utils.NewRealClock(), lggr)
	if err != nil {
		return nil, err
	}
	handler.SetConnector(connector)
	return connector, nil
}
