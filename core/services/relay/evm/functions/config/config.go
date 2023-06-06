package config

import (
	"google.golang.org/protobuf/proto"
)

func DecodeMultiPluginOffchainConfig(raw []byte) (functionsPluginConfig, thresholdPluginConfig []byte, err error) {
	configProto := &MultiPluginConfig{}
	if err := proto.Unmarshal(raw, configProto); err != nil {
		return nil, nil, err
	}
	return configProto.FunctionsPluginConfig, configProto.ThresholdPluginConfig, nil
}

func EncodeMultiPluginConfig(functionsPluginConfig, thresholdPluginConfig []byte) ([]byte, error) {
	config := &MultiPluginConfig{
		FunctionsPluginConfig: functionsPluginConfig,
		ThresholdPluginConfig: thresholdPluginConfig,
	}
	return proto.Marshal(config)
}
