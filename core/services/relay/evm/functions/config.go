package functions

import (
	functions_config "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
	"google.golang.org/protobuf/proto"
)

func DecodeThresholdReportingPluginConfig(raw []byte) (*functions_config.ThresholdReportingPluginConfig, error) {
	configProto := &functions_config.ReportingPluginConfig{}
	if err := proto.Unmarshal(raw, configProto); err != nil {
		return nil, err
	}
	return configProto.ThresholdPluginConfig, nil
}

func EncodeThresholdPluginConfig(thConfig *functions_config.ThresholdReportingPluginConfig) ([]byte, error) {
	return proto.Marshal(thConfig)
}
