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
	factory func(lggr logger.Logger, extraDataCodec ExtraDataCodec) PluginConfig) error {
	if _, exists := RegisteredPluginConfigFactories[chainFamily]; exists {
		return fmt.Errorf("plugin config factory for chain family %s already registered", chainFamily)
	}

	RegisteredPluginConfigFactories[chainFamily] = factory
	return nil
}

// RegisterCRCW registers a ChainRWProvider for a specific chain family.
func RegisterCRCW(
	chainFamily string,
	factory ChainRWProvider,
) error {
	if _, exists := RegisteredCRCW[chainFamily]; exists {
		return fmt.Errorf("CRCW factory for chain family %s already registered", chainFamily)
	}

	RegisteredCRCW[chainFamily] = factory
	return nil
}

// RegisterExtraDataCodec registers a SourceChainExtraDataCodec for a specific chain family.
func RegisterExtraDataCodec(
	chainFamily string,
	extraDataCodec SourceChainExtraDataCodec,
) error {
	if _, exists := RegisteredExtraDataCodec[chainFamily]; exists {
		return fmt.Errorf("extra data codec for chain family %s already registered", chainFamily)
	}

	RegisteredExtraDataCodec[chainFamily] = extraDataCodec
	return nil
}

// RegisterAddressCodec registers a ChainSpecificAddressCodec for a specific chain family.
func RegisterAddressCodec(
	chainFamily string,
	addressCodec ChainSpecificAddressCodec,
) error {
	if _, exists := RegisteredAddressCodec[chainFamily]; exists {
		return fmt.Errorf("address codec for chain family %s already registered", chainFamily)
	}

	RegisteredAddressCodec[chainFamily] = addressCodec
	return nil
}
