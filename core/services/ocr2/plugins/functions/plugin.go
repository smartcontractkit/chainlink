package functions

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"
	"golang.org/x/exp/slices"

	"github.com/smartcontractkit/libocr/commontypes"
	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2plus"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/ocr2dr_oracle"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/functions"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	gwFunctions "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
	s4_plugin "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/s4"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/threshold"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/s4"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type FunctionsServicesConfig struct {
	Job               job.Job
	JobORM            job.ORM
	BridgeORM         bridges.ORM
	QConfig           pg.QConfig
	DB                *sqlx.DB
	Chain             evm.Chain
	ContractID        string
	Logger            logger.Logger
	MailMon           *utils.MailboxMonitor
	URLsMonEndpoint   commontypes.MonitoringEndpoint
	EthKeystore       keystore.Eth
	ThresholdKeyShare []byte
}

const (
	FunctionsBridgeName     string = "ea_bridge"
	FunctionsS4Namespace    string = "functions"
	MaxAdapterResponseBytes int64  = 1_000_000
)

// Create all OCR2 plugin Oracles and all extra services needed to run a Functions job.
func NewFunctionsServices(functionsOracleArgs, thresholdOracleArgs, s4OracleArgs *libocr2.OCR2OracleArgs, conf *FunctionsServicesConfig) ([]job.ServiceCtx, error) {
	pluginORM := functions.NewORM(conf.DB, conf.Logger, conf.QConfig, common.HexToAddress(conf.ContractID))
	s4ORM := s4.NewPostgresORM(conf.DB, conf.Logger, conf.QConfig, s4.SharedTableName, FunctionsS4Namespace)

	var pluginConfig config.PluginConfig
	if err := json.Unmarshal(conf.Job.OCR2OracleSpec.PluginConfig.Bytes(), &pluginConfig); err != nil {
		return nil, err
	}
	if err := config.ValidatePluginConfig(pluginConfig); err != nil {
		return nil, err
	}

	allServices := []job.ServiceCtx{}
	contractAddress := common.HexToAddress(conf.Job.OCR2OracleSpec.ContractID)
	oracleContract, err := ocr2dr_oracle.NewOCR2DROracle(contractAddress, conf.Chain.Client())
	if err != nil {
		return nil, errors.Wrapf(err, "Functions: failed to create a FunctionsOracle wrapper for address: %v", contractAddress)
	}

	var decryptor threshold.Decryptor
	// thresholdOracleArgs nil check will be removed once the Threshold plugin is fully integrated w/ Functions
	if len(conf.ThresholdKeyShare) > 0 && thresholdOracleArgs != nil && pluginConfig.DecryptionQueueConfig != nil {
		decryptionQueue := threshold.NewDecryptionQueue(
			int(pluginConfig.DecryptionQueueConfig.MaxQueueLength),
			int(pluginConfig.DecryptionQueueConfig.MaxCiphertextBytes),
			int(pluginConfig.DecryptionQueueConfig.MaxCiphertextIdLength),
			time.Duration(pluginConfig.DecryptionQueueConfig.CompletedCacheTimeoutSec)*time.Second,
			conf.Logger.Named("DecryptionQueue"),
		)
		decryptor = decryptionQueue
		thresholdServicesConfig := threshold.ThresholdServicesConfig{
			DecryptionQueue:    decryptionQueue,
			KeyshareWithPubKey: conf.ThresholdKeyShare,
			ConfigParser:       config.ThresholdConfigParser{},
		}
		thresholdService, err2 := threshold.NewThresholdService(thresholdOracleArgs, &thresholdServicesConfig)
		if err2 != nil {
			return nil, errors.Wrap(err2, "error calling NewThresholdServices")
		}
		allServices = append(allServices, thresholdService)
	} else {
		conf.Logger.Warn("Threshold configuration is incomplete. Threshold secrets decryption plugin is disabled.")
	}

	listenerLogger := conf.Logger.Named("FunctionsListener")
	bridgeAccessor := functions.NewBridgeAccessor(conf.BridgeORM, FunctionsBridgeName, MaxAdapterResponseBytes)
	functionsListener := functions.NewFunctionsListener(
		oracleContract,
		conf.Job,
		bridgeAccessor,
		pluginORM,
		pluginConfig,
		conf.Chain.LogBroadcaster(),
		listenerLogger,
		conf.MailMon,
		conf.URLsMonEndpoint,
		decryptor,
	)
	allServices = append(allServices, functionsListener)

	functionsOracleArgs.ReportingPluginFactory = FunctionsReportingPluginFactory{
		Logger:    functionsOracleArgs.Logger,
		PluginORM: pluginORM,
		JobID:     conf.Job.ExternalJobID,
	}
	functionsReportingPluginOracle, err := libocr2.NewOracle(*functionsOracleArgs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call NewOracle to create a Functions Reporting Plugin")
	}
	allServices = append(allServices, job.NewServiceAdapter(functionsReportingPluginOracle))

	if pluginConfig.GatewayConnectorConfig != nil && pluginConfig.S4Constraints != nil && pluginConfig.OnchainAllowlist != nil {
		allowlist, err2 := gwFunctions.NewOnchainAllowlist(conf.Chain.Client(), *pluginConfig.OnchainAllowlist, conf.Logger)
		if err2 != nil {
			return nil, errors.Wrap(err, "failed to call NewOnchainAllowlist while creating a Functions Reporting Plugin")
		}
		s4Storage := s4.NewStorage(conf.Logger, *pluginConfig.S4Constraints, s4ORM, utils.NewRealClock())
		connectorLogger := conf.Logger.Named("GatewayConnector").With("jobName", conf.Job.PipelineSpec.JobName)
		connector, err3 := NewConnector(pluginConfig.GatewayConnectorConfig, conf.EthKeystore, conf.Chain.ID(), s4Storage, allowlist, connectorLogger)
		if err3 != nil {
			return nil, errors.Wrap(err, "failed to create a GatewayConnector")
		}
		allServices = append(allServices, connector)
	} else {
		listenerLogger.Warn("No GatewayConnectorConfig, S4Constraints or OnchainAllowlist is found in the plugin config, GatewayConnector will not be enabled")
	}

	if s4OracleArgs != nil && pluginConfig.S4Constraints != nil {
		s4OracleArgs.ReportingPluginFactory = s4_plugin.S4ReportingPluginFactory{
			Logger:        s4OracleArgs.Logger,
			ORM:           s4ORM,
			ConfigDecoder: config.S4ConfigDecoder,
		}
		s4ReportingPluginOracle, err := libocr2.NewOracle(*s4OracleArgs)
		if err != nil {
			return nil, errors.Wrap(err, "failed to call NewOracle to create a S4 Reporting Plugin")
		}
		allServices = append(allServices, job.NewServiceAdapter(s4ReportingPluginOracle))
	} else {
		listenerLogger.Warn("s4OracleArgs is nil or S4Constraints are not configured. S4 plugin is disabled.")
	}

	return allServices, nil
}

func NewConnector(gwcCfg *connector.ConnectorConfig, ethKeystore keystore.Eth, chainID *big.Int, s4Storage s4.Storage, allowlist gwFunctions.OnchainAllowlist, lggr logger.Logger) (connector.GatewayConnector, error) {
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
	nodeAddress := enabledKeys[idx].ID()

	handler := functions.NewFunctionsConnectorHandler(nodeAddress, signerKey, s4Storage, allowlist, lggr)
	connector, err := connector.NewGatewayConnector(gwcCfg, handler, handler, utils.NewRealClock(), lggr)
	if err != nil {
		return nil, err
	}
	handler.SetConnector(connector)
	return connector, nil
}
