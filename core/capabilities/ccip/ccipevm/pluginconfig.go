package ccipevm

import (
	chainsel "github.com/smartcontractkit/chain-selectors"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ocrimpls"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const defaultCommitGasLimit = 500_000

// PluginConfig is a struct that contains the configuration for a plugin.
type PluginConfig struct {
	extraDataCodec ccipcommon.ExtraDataCodec
}

// NewPluginConfig returns a new PluginConfig.
func NewPluginConfig(extraDataCodec ccipcommon.ExtraDataCodec) *PluginConfig {
	return &PluginConfig{extraDataCodec: extraDataCodec}
}

// InitializePluginConfig returns a PluginConfig for EVM chains.
func (p PluginConfig) InitializePluginConfig(lggr logger.Logger) ccipcommon.PluginConfig {
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
