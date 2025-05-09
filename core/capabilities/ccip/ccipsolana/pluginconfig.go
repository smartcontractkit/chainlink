package ccipsolana

import (
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ocrimpls"

	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// pluginConfig is a struct that contains the configuration for a plugin.
type pluginConfig struct {
	extraDataCodec ccipcommon.ExtraDataCodec
}

// NewPluginConfig returns a new pluginConfig.
func NewPluginConfig(extraDataCodec ccipcommon.ExtraDataCodec) ccipcommon.OffChainPluginConfig {
	return &pluginConfig{extraDataCodec: extraDataCodec}
}

// InitializePluginConfig returns a pluginConfig for Solana chains.
func (p pluginConfig) InitializePluginConfig(lggr logger.Logger) ccipcommon.PluginConfig {
	return ccipcommon.PluginConfig{
		CommitPluginCodec:          NewCommitPluginCodecV1(),
		ExecutePluginCodec:         NewExecutePluginCodecV1(p.extraDataCodec),
		MessageHasher:              NewMessageHasherV1(lggr.Named(chainsel.FamilySolana).Named("MessageHasherV1"), p.extraDataCodec),
		TokenDataEncoder:           NewSolanaTokenDataEncoder(),
		GasEstimateProvider:        NewGasEstimateProvider(p.extraDataCodec),
		RMNCrypto:                  nil,
		ContractTransmitterFactory: ocrimpls.NewSVMContractTransmitterFactory(p.extraDataCodec),
	}
}
