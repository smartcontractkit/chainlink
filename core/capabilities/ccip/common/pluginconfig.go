package common

import (
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"

	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
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
	EVMPluginConfig    PluginConfig
	SolanaPluginConfig PluginConfig
}

// NewPluginConfigFactory is a constructor for PluginConfigFactory.
func NewPluginConfigFactory(evmPluginConfig, solanaPluginConfig PluginConfig) *PluginConfigFactory {
	return &PluginConfigFactory{
		EVMPluginConfig:    evmPluginConfig,
		SolanaPluginConfig: solanaPluginConfig,
	}
}

// CreatePluginConfig creates a PluginConfig instance based on the chain family.
func (f *PluginConfigFactory) CreatePluginConfig(chainFamily string) (PluginConfig, error) {
	switch chainFamily {
	case chainsel.FamilyEVM:
		return f.EVMPluginConfig, nil
	case chainsel.FamilySolana:
		return f.SolanaPluginConfig, nil
	default:
		return PluginConfig{}, fmt.Errorf("unsupported chain family: %s", chainFamily)
	}
}
