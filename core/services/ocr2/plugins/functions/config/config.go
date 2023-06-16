package config

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	decryptionPluginConfig "github.com/smartcontractkit/tdh2/go/ocr2/decryptionplugin/config"
	"google.golang.org/protobuf/proto"
)

// This config is part of the job spec and is loaded only once on node boot/job creation.
type PluginConfig struct {
	MinIncomingConfirmations        uint32                     `json:"minIncomingConfirmations"`
	RequestTimeoutSec               uint32                     `json:"requestTimeoutSec"`
	RequestTimeoutCheckFrequencySec uint32                     `json:"requestTimeoutCheckFrequencySec"`
	RequestTimeoutBatchLookupSize   uint32                     `json:"requestTimeoutBatchLookupSize"`
	PruneMaxStoredRequests          uint32                     `json:"pruneMaxStoredRequests"`
	PruneCheckFrequencySec          uint32                     `json:"pruneCheckFrequencySec"`
	PruneBatchSize                  uint32                     `json:"pruneBatchSize"`
	ListenerEventHandlerTimeoutSec  uint32                     `json:"listenerEventHandlerTimeoutSec"`
	MaxRequestSizeBytes             uint32                     `json:"maxRequestSizeBytes"`
	GatewayConnectorConfig          *connector.ConnectorConfig `json:"gatewayConnectorConfig"`
}

func ValidatePluginConfig(config PluginConfig) error {
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
		},
	}, nil
}
