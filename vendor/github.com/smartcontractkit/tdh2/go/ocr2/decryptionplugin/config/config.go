package config

import (
	"google.golang.org/protobuf/proto"
)

// This config is stored in the Oracle contract (set via SetConfig()).
// Every SetConfig() call reloads the reporting plugin (DirectRequestReportingPluginFactory.NewReportingPlugin())
type ReportingPluginConfigWrapper struct {
	Config *ReportingPluginConfig
}

func DecodeReportingPluginConfig(raw []byte) (*ReportingPluginConfigWrapper, error) {
	configProto := &ReportingPluginConfig{}
	if err := proto.Unmarshal(raw, configProto); err != nil {
		return nil, err
	}
	return &ReportingPluginConfigWrapper{Config: configProto}, nil
}

func EncodeReportingPluginConfig(rpConfig *ReportingPluginConfigWrapper) ([]byte, error) {
	return proto.Marshal(rpConfig.Config)
}

//go:generate mockery --quiet --name ConfigParser --output ./mocks/ --case=underscore
type ConfigParser interface {
	ParseConfig(offchainConfig []byte) (*ReportingPluginConfigWrapper, error)
}

type DefaultConfigParser struct {
}

func (p *DefaultConfigParser) ParseConfig(offchainConfig []byte) (*ReportingPluginConfigWrapper, error) {
	return DecodeReportingPluginConfig(offchainConfig)
}
