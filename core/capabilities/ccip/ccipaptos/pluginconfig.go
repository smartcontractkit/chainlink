package ccipaptos

import (
	chainsel "github.com/smartcontractkit/chain-selectors"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ocrimpls"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// InitializePluginConfig returns a PluginConfig for EVM chains.
func InitializePluginConfig(lggr logger.Logger, extraDataCodec ccipcommon.ExtraDataCodec) ccipcommon.PluginConfig {
	return ccipcommon.PluginConfig{
		CommitPluginCodec:          NewCommitPluginCodecV1(),
		ExecutePluginCodec:         NewExecutePluginCodecV1(extraDataCodec),
		MessageHasher:              NewMessageHasherV1(lggr.Named(chainsel.FamilyAptos).Named("MessageHasherV1"), extraDataCodec),
		TokenDataEncoder:           NewAptosTokenDataEncoder(),
		GasEstimateProvider:        NewGasEstimateProvider(),
		RMNCrypto:                  nil,
		ContractTransmitterFactory: ocrimpls.NewAptosContractTransmitterFactory(extraDataCodec),
	}
}
