package ccipsolana

import (
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ocrimpls"

	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// PluginConfig is a struct that contains the configuration for a plugin.
type PluginConfig struct {
	extraDataCodec ccipcommon.ExtraDataCodec
}

// NewPluginConfig returns a new PluginConfig.
func NewPluginConfig(extraDataCodec ccipcommon.ExtraDataCodec) *PluginConfig {
	return &PluginConfig{extraDataCodec: extraDataCodec}
}

// InitializePluginConfig returns a PluginConfig for Solana chains.
func (p PluginConfig) InitializePluginConfig(lggr logger.Logger) ccipcommon.PluginConfig {
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
