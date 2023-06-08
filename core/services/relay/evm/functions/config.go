package functions

import (
	"google.golang.org/protobuf/proto"

	functions_config "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
)

func GetThresholdReportingPluginConfig(rawFunctionsPluginConfig []byte) (*functions_config.ThresholdReportingPluginConfig, error) {
	configProto := &functions_config.ReportingPluginConfig{}
	if err := proto.Unmarshal(rawFunctionsPluginConfig, configProto); err != nil {
		return nil, err
	}
	return configProto.ThresholdPluginConfig, nil
}

func EncodeThresholdPluginConfig(thConfig *functions_config.ThresholdReportingPluginConfig) ([]byte, error) {
	return proto.Marshal(thConfig)
}
