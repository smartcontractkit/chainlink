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
	AddrCodec             AddressCodec
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
	addressCodecMap := make(map[string]ChainSpecificAddressCodec)
	chainRWProviderMap := make(map[string]ChainRWProvider)
	CCIPProviderSupported := make(map[string]bool)

	for family, initFunc := range registeredFactories {
		config := initFunc(lggr, GetExtraDataCodecRegistry())
		CCIPProviderSupported[family] = config.CCIPProviderSupported
		if config.AddressCodec != nil {
			addressCodecMap[family] = config.AddressCodec
		}

		// Register all known families, this includes families whose SourceChainExtraDataCodec is provided
		// by the CCIPProvider
		extraDataCodecRegistry.RegisterFamily(family)
		if config.ExtraDataCodec != nil {
			// Register the actual codec for this family if we have it defined in core already
			extraDataCodecRegistry.RegisterCodec(family, config.ExtraDataCodec)
		}
		if config.ChainRW != nil {
			chainRWProviderMap[family] = config.ChainRW
		}
		if family == chainFamily {
			pluginServices.PluginConfig = config
		}
	}

	pluginServices.AddrCodec = NewAddressCodec(addressCodecMap)
	pluginServices.ChainRW = NewCRCW(chainRWProviderMap)
	pluginServices.CCIPProviderSupported = CCIPProviderSupported
	return pluginServices, nil
}
