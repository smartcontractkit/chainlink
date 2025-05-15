package common

import (
	"fmt"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// PluginConfig holds the configuration for a plugin.
type PluginConfig struct {
	CommitPluginCodec          cciptypes.CommitPluginCodec
	ExecutePluginCodec         cciptypes.ExecutePluginCodec
	MessageHasher              cciptypes.MessageHasher
	TokenDataEncoder           cciptypes.TokenDataEncoder
	GasEstimateProvider        cciptypes.EstimateProvider
	RMNCrypto                  cciptypes.RMNCrypto
	ContractTransmitterFactory cctypes.ContractTransmitterFactory
	// PriceOnlyCommitFn optional method override for price only commit reports.
	PriceOnlyCommitFn string
	ChainRW           ChainRWProvider
	AddressCodec      ChainSpecificAddressCodec
}

// PluginServices aggregates services for a specific chain family.
type PluginServices struct {
	PluginConfig   PluginConfig
	AddrCodec      AddressCodec
	ExtraDataCodec ExtraDataCodec
	ChainRW        MultiChainRW
}

// InitFunction defines a function to initialize a PluginConfig.
type InitFunction func(logger.Logger, ExtraDataCodec) PluginConfig

var registeredFactories = make(map[string]InitFunction)
var registeredExtraDataCodec = make(map[string]SourceChainExtraDataCodec)

// RegisterPluginConfig registers a plugin config factory for a chain family.
func RegisterPluginConfig(chainFamily string, factory InitFunction, extraDataCodec SourceChainExtraDataCodec) {
	registeredExtraDataCodec[chainFamily] = extraDataCodec
	registeredFactories[chainFamily] = factory
}

// GetPluginServices initializes and returns PluginServices for a chain family.
func GetPluginServices(lggr logger.Logger, chainFamily string) (PluginServices, error) {
	_, exists := registeredFactories[chainFamily]
	if !exists {
		return PluginServices{}, fmt.Errorf("unsupported chain family: %s", chainFamily)
	}

	pluginServices := PluginServices{}
	addressCodecMap := make(map[string]ChainSpecificAddressCodec)
	chainRWProviderMap := make(map[string]ChainRWProvider)
	pluginServices.ExtraDataCodec = NewExtraDataCodec(registeredExtraDataCodec)

	for family, initFunc := range registeredFactories {
		config := initFunc(lggr, pluginServices.ExtraDataCodec)
		addressCodecMap[family] = config.AddressCodec
		chainRWProviderMap[family] = config.ChainRW
		if family == chainFamily {
			pluginServices.PluginConfig = config
		}
	}

	pluginServices.AddrCodec = NewAddressCodec(addressCodecMap)
	pluginServices.ChainRW = NewCRCW(chainRWProviderMap)
	return pluginServices, nil
}
