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
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/threshold"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
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
	Lggr              logger.Logger
	MailMon           *utils.MailboxMonitor
	URLsMonEndpoint   commontypes.MonitoringEndpoint
	EthKeystore       keystore.Eth
	ThresholdKeyShare []byte
}

const (
	FunctionsBridgeName     string = "ea_bridge"
	MaxAdapterResponseBytes int64  = 1_000_000
)

// Create all OCR2 plugin Oracles and all extra services needed to run a Functions job.
func NewFunctionsServices(functionsOracleArgs, thresholdOracleArgs *libocr2.OCR2OracleArgs, conf *FunctionsServicesConfig) ([]job.ServiceCtx, error) {
	pluginORM := functions.NewORM(conf.DB, conf.Lggr, conf.QConfig, common.HexToAddress(conf.ContractID))

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

	var decryptor threshold.Decryptor
	// thresholdOracleArgs nil check will be removed once the Threshold plugin is fully integrated w/ Functions
	if len(conf.ThresholdKeyShare) > 0 && thresholdOracleArgs != nil {
		dqc := getDecryptionQueueConfig(pluginConfig)
		decryptionQueue := threshold.NewDecryptionQueue(
			int(dqc.MaxDecryptionQueueLength),
			int(dqc.MaxCiphertextBytes),
			int(dqc.MaxCiphertextIdLength),
			time.Duration(dqc.CompletedDecryptionCacheTimeoutSec)*time.Second,
			conf.Lggr.Named("DecryptionQueue"),
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
		conf.Lggr.Warn("ThresholdKeyShare is empty. Threshold secrets decryption plugin is disabled.")
	}

	listenerLogger := conf.Lggr.Named("FunctionsListener")
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

func getDecryptionQueueConfig(pluginConfig config.PluginConfig) config.DecryptionQueueConfig {
	dqc := config.DecryptionQueueConfig{
		MaxDecryptionQueueLength:           100,
		MaxCiphertextBytes:                 10_000,
		MaxCiphertextIdLength:              100,
		CompletedDecryptionCacheTimeoutSec: 300,
	}

	if pluginConfig.MaxDecryptionQueueLength > 0 {
		dqc.MaxDecryptionQueueLength = pluginConfig.MaxDecryptionQueueLength
	}
	if pluginConfig.MaxCiphertextBytes > 0 {
		dqc.MaxCiphertextBytes = pluginConfig.MaxCiphertextBytes
	}
	if pluginConfig.MaxCiphertextIdLength > 0 {
		dqc.MaxCiphertextIdLength = pluginConfig.MaxCiphertextIdLength
	}
	if pluginConfig.CompletedDecryptionCacheTimeoutSec > 0 {
		dqc.CompletedDecryptionCacheTimeoutSec = pluginConfig.CompletedDecryptionCacheTimeoutSec
	}

	return dqc
}
