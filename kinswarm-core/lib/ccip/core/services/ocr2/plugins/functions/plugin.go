package functions

import (
	"context"
	"encoding/json"
	"math/big"
	"slices"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jonboulle/clockwork"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/libocr/commontypes"
	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2plus"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/functions"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	hc "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
	hf "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions"
	gwAllowlist "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions/allowlist"
	gwSubscriptions "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions/subscriptions"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
	s4_plugin "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/s4"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/threshold"
	evmrelayTypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/s4"
)

type FunctionsServicesConfig struct {
	Job               job.Job
	JobORM            job.ORM
	BridgeORM         bridges.ORM
	DS                sqlutil.DataSource
	Chain             legacyevm.Chain
	ContractID        string
	Logger            logger.Logger
	MailMon           *mailbox.Monitor
	URLsMonEndpoint   commontypes.MonitoringEndpoint
	EthKeystore       keystore.Eth
	ThresholdKeyShare []byte
	LogPollerWrapper  evmrelayTypes.LogPollerWrapper
}

const (
	FunctionsBridgeName                   string = "ea_bridge"
	FunctionsS4Namespace                  string = "functions"
	MaxAdapterResponseBytes               int64  = 1_000_000
	DefaultOffchainTransmitterChannelSize uint32 = 1000
	DefaultMaxAdapterRetry                int    = 3
	DefaultExponentialBackoffBase                = 5 * time.Second
)

// Create all OCR2 plugin Oracles and all extra services needed to run a Functions job.
func NewFunctionsServices(ctx context.Context, functionsOracleArgs, thresholdOracleArgs, s4OracleArgs *libocr2.OCR2OracleArgs, conf *FunctionsServicesConfig) ([]job.ServiceCtx, error) {
	pluginORM := functions.NewORM(conf.DS, common.HexToAddress(conf.ContractID))
	s4ORM := s4.NewCachedORMWrapper(s4.NewPostgresORM(conf.DS, s4.SharedTableName, FunctionsS4Namespace), conf.Logger)

	var pluginConfig config.PluginConfig
	if err := json.Unmarshal(conf.Job.OCR2OracleSpec.PluginConfig.Bytes(), &pluginConfig); err != nil {
		return nil, err
	}
	if err := config.ValidatePluginConfig(pluginConfig); err != nil {
		return nil, err
	}

	allServices := []job.ServiceCtx{}

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

	var s4Storage s4.Storage
	if pluginConfig.S4Constraints != nil {
		s4Storage = s4.NewStorage(conf.Logger, *pluginConfig.S4Constraints, s4ORM, clockwork.NewRealClock())
	}

	offchainTransmitter := functions.NewOffchainTransmitter(DefaultOffchainTransmitterChannelSize)
	listenerLogger := conf.Logger.Named("FunctionsListener")

	var maxRetries int
	if pluginConfig.ExternalAdapterMaxRetries != nil {
		maxRetries = int(*pluginConfig.ExternalAdapterMaxRetries)
	} else {
		maxRetries = DefaultMaxAdapterRetry
	}
	conf.Logger.Debugf("external adapter maxRetries configured to: %d", maxRetries)

	var exponentialBackoffBase time.Duration
	if pluginConfig.ExternalAdapterExponentialBackoffBaseSec != nil {
		exponentialBackoffBase = time.Duration(*pluginConfig.ExternalAdapterExponentialBackoffBaseSec) * time.Second
	} else {
		exponentialBackoffBase = DefaultExponentialBackoffBase
	}
	conf.Logger.Debugf("external adapter exponentialBackoffBase configured to: %g sec", exponentialBackoffBase.Seconds())

	bridgeAccessor := functions.NewBridgeAccessor(conf.BridgeORM, FunctionsBridgeName, MaxAdapterResponseBytes, maxRetries, exponentialBackoffBase)
	functionsListener := functions.NewFunctionsListener(
		conf.Job,
		conf.Chain.Client(),
		conf.Job.OCR2OracleSpec.ContractID,
		bridgeAccessor,
		pluginORM,
		pluginConfig,
		s4Storage,
		listenerLogger,
		conf.URLsMonEndpoint,
		decryptor,
		conf.LogPollerWrapper,
	)
	allServices = append(allServices, functionsListener)

	functionsOracleArgs.ReportingPluginFactory = FunctionsReportingPluginFactory{
		Logger:              functionsOracleArgs.Logger,
		PluginORM:           pluginORM,
		JobID:               conf.Job.ExternalJobID,
		ContractVersion:     pluginConfig.ContractVersion,
		OffchainTransmitter: offchainTransmitter,
	}
	functionsReportingPluginOracle, err := libocr2.NewOracle(*functionsOracleArgs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call NewOracle to create a Functions Reporting Plugin")
	}
	allServices = append(allServices, job.NewServiceAdapter(functionsReportingPluginOracle))

	if pluginConfig.GatewayConnectorConfig != nil && s4Storage != nil && pluginConfig.OnchainAllowlist != nil && pluginConfig.RateLimiter != nil && pluginConfig.OnchainSubscriptions != nil {
		allowlistORM, err := gwAllowlist.NewORM(conf.DS, conf.Logger, pluginConfig.OnchainAllowlist.ContractAddress)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create allowlist ORM")
		}
		allowlist, err2 := gwAllowlist.NewOnchainAllowlist(conf.Chain.Client(), *pluginConfig.OnchainAllowlist, allowlistORM, conf.Logger)
		if err2 != nil {
			return nil, errors.Wrap(err, "failed to create OnchainAllowlist")
		}
		rateLimiter, err2 := hc.NewRateLimiter(*pluginConfig.RateLimiter)
		if err2 != nil {
			return nil, errors.Wrap(err, "failed to create a RateLimiter")
		}
		subscriptionsORM, err := gwSubscriptions.NewORM(conf.DS, conf.Logger, pluginConfig.OnchainSubscriptions.ContractAddress)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create subscriptions ORM")
		}
		subscriptions, err2 := gwSubscriptions.NewOnchainSubscriptions(conf.Chain.Client(), *pluginConfig.OnchainSubscriptions, subscriptionsORM, conf.Logger)
		if err2 != nil {
			return nil, errors.Wrap(err, "failed to create a OnchainSubscriptions")
		}
		connectorLogger := conf.Logger.Named("GatewayConnector").With("jobName", conf.Job.PipelineSpec.JobName)
		connector, handler, err2 := NewConnector(ctx, &pluginConfig, conf.EthKeystore, conf.Chain.ID(), s4Storage, allowlist, rateLimiter, subscriptions, functionsListener, offchainTransmitter, connectorLogger)
		if err2 != nil {
			return nil, errors.Wrap(err, "failed to create a GatewayConnector")
		}
		allServices = append(allServices, connector)
		allServices = append(allServices, handler)
	} else {
		listenerLogger.Warn("Insufficient config, GatewayConnector will not be enabled")
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

func NewConnector(ctx context.Context, pluginConfig *config.PluginConfig, ethKeystore keystore.Eth, chainID *big.Int, s4Storage s4.Storage, allowlist gwAllowlist.OnchainAllowlist, rateLimiter *hc.RateLimiter, subscriptions gwSubscriptions.OnchainSubscriptions, listener functions.FunctionsListener, offchainTransmitter functions.OffchainTransmitter, lggr logger.Logger) (connector.GatewayConnector, connector.GatewayConnectorHandler, error) {
	enabledKeys, err := ethKeystore.EnabledKeysForChain(ctx, chainID)
	if err != nil {
		return nil, nil, err
	}
	configuredNodeAddress := common.HexToAddress(pluginConfig.GatewayConnectorConfig.NodeAddress)
	idx := slices.IndexFunc(enabledKeys, func(key ethkey.KeyV2) bool { return key.Address == configuredNodeAddress })
	if idx == -1 {
		return nil, nil, errors.New("key for configured node address not found")
	}
	signerKey := enabledKeys[idx].ToEcdsaPrivKey()
	if enabledKeys[idx].ID() != pluginConfig.GatewayConnectorConfig.NodeAddress {
		return nil, nil, errors.New("node address mismatch")
	}

	handler, err := functions.NewFunctionsConnectorHandler(pluginConfig, signerKey, s4Storage, allowlist, rateLimiter, subscriptions, listener, offchainTransmitter, lggr)
	if err != nil {
		return nil, nil, err
	}
	// handler acts as a signer here
	connector, err := connector.NewGatewayConnector(pluginConfig.GatewayConnectorConfig, handler, clockwork.NewRealClock(), lggr)
	if err != nil {
		return nil, nil, err
	}
	err = connector.AddHandler([]string{hf.MethodSecretsSet, hf.MethodSecretsList, hf.MethodHeartbeat}, handler)
	if err != nil {
		return nil, nil, err
	}
	handler.SetConnector(connector)
	return connector, handler, nil
}
