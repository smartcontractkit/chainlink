package ccipnoop

import (
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ocrimpls"
)

// NewPluginConfig returns a pluginConfig .
func NewPluginConfig(lggr logger.Logger, extraDataCodec ccipcommon.ExtraDataCodec) ccipcommon.PluginConfig {
	return ccipcommon.PluginConfig{
		CommitPluginCodec:          NewCommitPluginCodecV1(),
		ExecutePluginCodec:         NewExecutePluginCodecV1(extraDataCodec),
		MessageHasher:              NewMessageHasherV1(lggr, extraDataCodec),
		TokenDataEncoder:           NewTokenDataEncoder(),
		GasEstimateProvider:        NewGasEstimateProvider(extraDataCodec),
		RMNCrypto:                  &NoopRMNCrypto{},
		AddressCodec:               AddressCodec{},
		ChainRW:                    chainRWProvider{},
		ContractTransmitterFactory: ocrimpls.NewContractTransmitterFactory(extraDataCodec),
		ExtraDataCodec:             extraDataDecoder{},
	}
}
