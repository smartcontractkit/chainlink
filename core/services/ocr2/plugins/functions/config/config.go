package config

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"

	decryptionPluginConfig "github.com/smartcontractkit/tdh2/go/ocr2/decryptionplugin/config"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink-common/pkg/assets"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions/allowlist"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions/subscriptions"
	s4PluginConfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/s4"
	"github.com/smartcontractkit/chainlink/v2/core/services/s4"
)

// This config is part of the job spec and is loaded only once on node boot/job creation.
type PluginConfig struct {
	EnableRequestSignatureCheck              bool                                      `json:"enableRequestSignatureCheck"`
	DONID                                    string                                    `json:"donID"`
	ContractVersion                          uint32                                    `json:"contractVersion"`
	MinRequestConfirmations                  uint32                                    `json:"minRequestConfirmations"`
	MinResponseConfirmations                 uint32                                    `json:"minResponseConfirmations"`
	MinIncomingConfirmations                 uint32                                    `json:"minIncomingConfirmations"`
	PastBlocksToPoll                         uint32                                    `json:"pastBlocksToPoll"`
	LogPollerCacheDurationSec                uint32                                    `json:"logPollerCacheDurationSec"` // Duration to cache previously detected request or response logs such that they can be filtered when calling logpoller_wrapper.LatestEvents()
	RequestTimeoutSec                        uint32                                    `json:"requestTimeoutSec"`
	RequestTimeoutCheckFrequencySec          uint32                                    `json:"requestTimeoutCheckFrequencySec"`
	RequestTimeoutBatchLookupSize            uint32                                    `json:"requestTimeoutBatchLookupSize"`
	PruneMaxStoredRequests                   uint32                                    `json:"pruneMaxStoredRequests"`
	PruneCheckFrequencySec                   uint32                                    `json:"pruneCheckFrequencySec"`
	PruneBatchSize                           uint32                                    `json:"pruneBatchSize"`
	ListenerEventHandlerTimeoutSec           uint32                                    `json:"listenerEventHandlerTimeoutSec"`
	ListenerEventsCheckFrequencyMillis       uint32                                    `json:"listenerEventsCheckFrequencyMillis"`
	ContractUpdateCheckFrequencySec          uint32                                    `json:"contractUpdateCheckFrequencySec"`
	MaxRequestSizeBytes                      uint32                                    `json:"maxRequestSizeBytes"`
	MaxRequestSizesList                      []uint32                                  `json:"maxRequestSizesList"`
	MaxSecretsSizesList                      []uint32                                  `json:"maxSecretsSizesList"`
	MinimumSubscriptionBalance               assets.Link                               `json:"minimumSubscriptionBalance"`
	AllowedHeartbeatInitiators               []string                                  `json:"allowedHeartbeatInitiators"`
	GatewayConnectorConfig                   *connector.ConnectorConfig                `json:"gatewayConnectorConfig"`
	OnchainAllowlist                         *allowlist.OnchainAllowlistConfig         `json:"onchainAllowlist"`
	OnchainSubscriptions                     *subscriptions.OnchainSubscriptionsConfig `json:"onchainSubscriptions"`
	RateLimiter                              *common.RateLimiterConfig                 `json:"rateLimiter"`
	S4Constraints                            *s4.Constraints                           `json:"s4Constraints"`
	DecryptionQueueConfig                    *DecryptionQueueConfig                    `json:"decryptionQueueConfig"`
	ExternalAdapterMaxRetries                *uint32                                   `json:"externalAdapterMaxRetries"`
	ExternalAdapterExponentialBackoffBaseSec *uint32                                   `json:"externalAdapterExponentialBackoffBaseSec"`
}

type DecryptionQueueConfig struct {
	MaxQueueLength           uint32 `json:"maxQueueLength"`
	MaxCiphertextBytes       uint32 `json:"maxCiphertextBytes"`
	MaxCiphertextIdLength    uint32 `json:"maxCiphertextIdLength"`
	CompletedCacheTimeoutSec uint32 `json:"completedCacheTimeoutSec"`
	DecryptRequestTimeoutSec uint32 `json:"decryptRequestTimeoutSec"`
}

func ValidatePluginConfig(config PluginConfig) error {
	if config.DecryptionQueueConfig != nil {
		if config.DecryptionQueueConfig.MaxQueueLength <= 0 {
			return errors.New("missing or invalid decryptionQueueConfig maxQueueLength")
		}
		if config.DecryptionQueueConfig.MaxCiphertextBytes <= 0 {
			return errors.New("missing or invalid decryptionQueueConfig maxCiphertextBytes")
		}
		if config.DecryptionQueueConfig.MaxCiphertextIdLength <= 0 {
			return errors.New("missing or invalid decryptionQueueConfig maxCiphertextIdLength")
		}
		if config.DecryptionQueueConfig.CompletedCacheTimeoutSec <= 0 {
			return errors.New("missing or invalid decryptionQueueConfig completedCacheTimeoutSec")
		}
		if config.DecryptionQueueConfig.DecryptRequestTimeoutSec <= 0 {
			return errors.New("missing or invalid decryptionQueueConfig decryptRequestTimeoutSec")
		}
	}
	return nil
}

// This config is stored in the Oracle contract (set via SetConfig()).
// Every SetConfig() call reloads the reporting plugin (FunctionsReportingPluginFactory.NewReportingPlugin())
type ReportingPluginConfigWrapper struct {
	Config *ReportingPluginConfig
}

func DecodeReportingPluginConfig(raw []byte) (*ReportingPluginConfigWrapper, error) {
	configProto := &ReportingPluginConfig{}
	err := proto.Unmarshal(raw, configProto)
	if err != nil {
		return nil, err
	}
	return &ReportingPluginConfigWrapper{Config: configProto}, nil
}

func EncodeReportingPluginConfig(rpConfig *ReportingPluginConfigWrapper) ([]byte, error) {
	return proto.Marshal(rpConfig.Config)
}

var _ decryptionPluginConfig.ConfigParser = &ThresholdConfigParser{}

type ThresholdConfigParser struct{}

func (ThresholdConfigParser) ParseConfig(config []byte) (*decryptionPluginConfig.ReportingPluginConfigWrapper, error) {
	reportingPluginConfigWrapper, err := DecodeReportingPluginConfig(config)
	if err != nil {
		return nil, errors.New("failed to decode Functions Threshold plugin config")
	}
	thresholdPluginConfig := reportingPluginConfigWrapper.Config.ThresholdPluginConfig

	if thresholdPluginConfig == nil {
		return nil, fmt.Errorf("PluginConfig bytes %x did not contain threshold plugin config", config)
	}

	return &decryptionPluginConfig.ReportingPluginConfigWrapper{
		Config: &decryptionPluginConfig.ReportingPluginConfig{
			MaxQueryLengthBytes:       thresholdPluginConfig.MaxQueryLengthBytes,
			MaxObservationLengthBytes: thresholdPluginConfig.MaxObservationLengthBytes,
			MaxReportLengthBytes:      thresholdPluginConfig.MaxReportLengthBytes,
			RequestCountLimit:         thresholdPluginConfig.RequestCountLimit,
			RequestTotalBytesLimit:    thresholdPluginConfig.RequestTotalBytesLimit,
			RequireLocalRequestCheck:  thresholdPluginConfig.RequireLocalRequestCheck,
			K:                         thresholdPluginConfig.K,
		},
	}, nil
}

func S4ConfigDecoder(config []byte) (*s4PluginConfig.PluginConfig, *types.ReportingPluginLimits, error) {
	reportingPluginConfigWrapper, err := DecodeReportingPluginConfig(config)
	if err != nil {
		return nil, nil, errors.New("failed to decode S4 plugin config")
	}

	pluginConfig := reportingPluginConfigWrapper.Config.S4PluginConfig
	if pluginConfig == nil {
		return nil, nil, fmt.Errorf("PluginConfig bytes %x did not contain s4 plugin config", config)
	}

	return &s4PluginConfig.PluginConfig{
			ProductName:             "functions",
			NSnapshotShards:         uint(pluginConfig.NSnapshotShards),
			MaxObservationEntries:   uint(pluginConfig.MaxObservationEntries),
			MaxReportEntries:        uint(pluginConfig.MaxReportEntries),
			MaxDeleteExpiredEntries: uint(pluginConfig.MaxDeleteExpiredEntries),
		},
		&types.ReportingPluginLimits{
			MaxQueryLength:       int(pluginConfig.MaxQueryLengthBytes),
			MaxObservationLength: int(pluginConfig.MaxObservationLengthBytes),
			MaxReportLength:      int(pluginConfig.MaxReportLengthBytes),
		},
		nil
}
