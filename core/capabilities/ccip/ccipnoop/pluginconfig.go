package ccipnoop

import (
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ocrimpls"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// InitializePluginConfig returns a pluginConfig for Solana chains.
func InitializePluginConfig(lggr logger.Logger, extraDataCodec ccipcommon.ExtraDataCodec) ccipcommon.PluginConfig {
	return ccipcommon.PluginConfig{
		CommitPluginCodec:          NewNoopCommitPluginCodecV1(),
		ExecutePluginCodec:         NewNoopExecutePluginCodecV1(extraDataCodec),
		MessageHasher:              NewNoopMessageHasherV1(lggr, extraDataCodec),
		TokenDataEncoder:           NewNoopTokenDataEncoder(),
		GasEstimateProvider:        NewNoopGasEstimateProvider(extraDataCodec),
		RMNCrypto:                  &NoopRMNCrypto{},
		AddressCodec:               NoopAddressCodec{},
		ChainRW:                    NoopChainRWProvider{},
		ContractTransmitterFactory: ocrimpls.NewNoopTransmitter(extraDataCodec),
		ExtraDataCodec:             NoopExtraDataDecoder{},
	}
}
