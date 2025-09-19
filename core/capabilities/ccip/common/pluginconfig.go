package common

import (
	"fmt"
	"maps"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
)

// PluginConfig holds the configuration for a plugin.
type PluginConfig struct {
	CommitPluginCodec          ccipocr3.CommitPluginCodec
	ExecutePluginCodec         ccipocr3.ExecutePluginCodec
	MessageHasher              ccipocr3.MessageHasher
	TokenDataEncoder           ccipocr3.TokenDataEncoder
	GasEstimateProvider        ccipocr3.EstimateProvider
	RMNCrypto                  ccipocr3.RMNCrypto
	ContractTransmitterFactory cctypes.ContractTransmitterFactory
	// PriceOnlyCommitFn optional method override for price only commit reports.
	PriceOnlyCommitFn     string
	ChainRW               ChainRWProvider
	AddressCodec          ChainSpecificAddressCodec
	ExtraDataCodec        SourceChainExtraDataCodec
	CCIPProviderSupported bool
}

// PluginServices aggregates services for a specific chain family.
type PluginServices struct {
	PluginConfig          PluginConfig
	ChainRW               MultiChainRW
	CCIPProviderSupported map[string]bool
}

// InitFunction defines a function to initialize a PluginConfig.
type InitFunction func(logger.Logger, ccipocr3.ExtraDataCodecBundle) PluginConfig

var registeredFactories = make(map[string]InitFunction)

// RegisterPluginConfig registers a plugin config factory for a chain family.
func RegisterPluginConfig(chainFamily string, factory InitFunction) {
	registeredFactories[chainFamily] = factory
}

// GetPluginServices initializes and returns PluginServices for a chain family.
func GetPluginServices(lggr logger.Logger, chainFamily string) (PluginServices, error) {
	_, exists := registeredFactories[chainFamily]
	if !exists {
		return PluginServices{}, fmt.Errorf("unsupported chain family: %s (available: %v)", chainFamily, maps.Keys(registeredFactories))
	}

	pluginServices := PluginServices{}
	extraDataCodecRegistry := GetExtraDataCodecRegistry()
	addressCodecRegistry := GetAddressCodecRegistry()
	chainRWProviderMap := make(map[string]ChainRWProvider)
	CCIPProviderSupported := make(map[string]bool)

	for family, initFunc := range registeredFactories {
		config := initFunc(lggr, GetExtraDataCodecRegistry())
		CCIPProviderSupported[family] = config.CCIPProviderSupported

		// Add all families to the registries. If the codecs are provided by the config, set them here, otherwise
		// ccipProvider will set them later in the oracle creator.
		extraDataCodecRegistry.RegisterFamily(family)
		addressCodecRegistry.RegisterFamily(family)

		if config.ExtraDataCodec != nil {
			extraDataCodecRegistry.RegisterCodec(family, config.ExtraDataCodec)
		}
		if config.AddressCodec != nil {
			addressCodecRegistry.RegisterCodec(family, config.AddressCodec)
		}
		if config.ChainRW != nil {
			chainRWProviderMap[family] = config.ChainRW
		}
		if family == chainFamily {
			pluginServices.PluginConfig = config
		}
	}

	pluginServices.ChainRW = NewCRCW(chainRWProviderMap)
	pluginServices.CCIPProviderSupported = CCIPProviderSupported
	return pluginServices, nil
}
