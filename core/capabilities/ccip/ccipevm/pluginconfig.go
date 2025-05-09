package ccipevm

import (
	chainsel "github.com/smartcontractkit/chain-selectors"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ocrimpls"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const defaultCommitGasLimit = 500_000

// pluginConfig is a struct that contains the configuration for a plugin.
type pluginConfig struct {
	extraDataCodec ccipcommon.ExtraDataCodec
}

// NewPluginConfig returns a new pluginConfig.
func NewPluginConfig(extraDataCodec ccipcommon.ExtraDataCodec) ccipcommon.OffChainPluginConfig {
	return pluginConfig{extraDataCodec: extraDataCodec}
}

// InitializePluginConfig returns a pluginConfig for EVM chains.
func (p pluginConfig) InitializePluginConfig(lggr logger.Logger) ccipcommon.PluginConfig {
	return ccipcommon.PluginConfig{
		CommitPluginCodec:          NewCommitPluginCodecV1(),
		ExecutePluginCodec:         NewExecutePluginCodecV1(p.extraDataCodec),
		MessageHasher:              NewMessageHasherV1(lggr.Named(chainsel.FamilyEVM).Named("MessageHasherV1"), p.extraDataCodec),
		TokenDataEncoder:           NewEVMTokenDataEncoder(),
		GasEstimateProvider:        NewGasEstimateProvider(p.extraDataCodec),
		RMNCrypto:                  NewEVMRMNCrypto(lggr.Named(chainsel.FamilyEVM).Named("RMNCrypto")),
		ContractTransmitterFactory: ocrimpls.NewEVMContractTransmitterFactory(p.extraDataCodec),
	}
}
