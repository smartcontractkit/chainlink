package functions

import (
	functions_config "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
	"google.golang.org/protobuf/proto"
)

func DecodeFunctionsReportingPluginConfig(raw []byte) (*functions_config.ReportingPluginConfig, error) {
	configProto := &functions_config.ReportingPluginConfig{}
	if err := proto.Unmarshal(raw, configProto); err != nil {
		return nil, err
	}
	return configProto, nil
}

func DecodeThresholdReportingPluginConfig(raw []byte) (*functions_config.ThresholdReportingPluginConfig, error) {
	configProto := &functions_config.ReportingPluginConfig{}
	if err := proto.Unmarshal(raw, configProto); err != nil {
		return nil, err
	}
	return configProto.ThresholdPluginConfig, nil
}

func EncodeFunctionsPluginConfig(rpConfig *functions_config.ReportingPluginConfigWrapper) ([]byte, error) {
	return proto.Marshal(rpConfig.Config)
}

func EncodeThresholdPluginConfig(thConfig *functions_config.ThresholdReportingPluginConfig) ([]byte, error) {
	return proto.Marshal(thConfig)
}
