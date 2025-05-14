package common

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/v2/core/logger"

	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

var (
	// RegisteredPluginConfigFactories is a map that holds the registered plugin config factories. It will be
	RegisteredPluginConfigFactories = make(map[string]func(lggr logger.Logger, extraDataCodec ExtraDataCodec) PluginConfig)

	// RegisteredCRCW is a map that holds the registered ChainRWProvider factories. It will be used to create
	RegisteredCRCW = make(map[string]ChainRWProvider)

	// RegisteredExtraDataCodec is a map that holds the registered SourceChainExtraDataCodec factories. It will be used to create
	RegisteredExtraDataCodec = make(map[string]SourceChainExtraDataCodec)

	// RegisteredAddressCodec is a map that holds the registered ChainSpecificAddressCodec factories. It will be used to create
	RegisteredAddressCodec = make(map[string]ChainSpecificAddressCodec)
)

// PluginConfig is a struct that contains the configuration for a plugin.
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
}

// PluginConfigFactory is a factory for creating PluginConfig instances.
type PluginConfigFactory struct {
	lggr           logger.Logger
	extraDataCodec ExtraDataCodec
}

// NewPluginConfigFactory is a constructor for PluginConfigFactory.
func NewPluginConfigFactory(lggr logger.Logger, extraDataCodec ExtraDataCodec) *PluginConfigFactory {
	return &PluginConfigFactory{
		lggr:           lggr,
		extraDataCodec: extraDataCodec,
	}
}

// CreatePluginConfig creates a PluginConfig instance based on the chain family.
func (f *PluginConfigFactory) CreatePluginConfig(chainFamily string) (PluginConfig, error) {
	pluginConfigFactory, exist := RegisteredPluginConfigFactories[chainFamily]
	if !exist {
		return PluginConfig{}, fmt.Errorf("unsupported chain family: %s", chainFamily)
	}

	return pluginConfigFactory(f.lggr, f.extraDataCodec), nil
}

// RegisterPluginConfig registers a plugin config factory for a specific chain family.
func RegisterPluginConfig(
	chainFamily string,
	pluginConfigFactory func(lggr logger.Logger, extraDataCodec ExtraDataCodec) PluginConfig,
	crw ChainRWProvider,
	extraDataCodec SourceChainExtraDataCodec,
	addressCodec ChainSpecificAddressCodec) {
	RegisteredExtraDataCodec[chainFamily] = extraDataCodec
	RegisteredPluginConfigFactories[chainFamily] = pluginConfigFactory
	RegisteredCRCW[chainFamily] = crw
	RegisteredAddressCodec[chainFamily] = addressCodec
}
