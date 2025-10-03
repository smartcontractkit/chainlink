package ccipnoop

import (
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

// NewPluginConfig returns a pluginConfig .
func NewPluginConfig(lggr logger.Logger, extraDataCodec ccipocr3.ExtraDataCodecBundle) ccipcommon.PluginConfig {
	return ccipcommon.PluginConfig{
		CommitPluginCodec:   NewCommitPluginCodecV1(),
		ExecutePluginCodec:  NewExecutePluginCodecV1(extraDataCodec),
		MessageHasher:       NewMessageHasherV1(lggr, extraDataCodec),
		TokenDataEncoder:    NewTokenDataEncoder(),
		GasEstimateProvider: NewGasEstimateProvider(extraDataCodec),
		RMNCrypto:           &NoopRMNCrypto{},
		AddressCodec:        AddressCodec{},
		ChainRW:             chainRWProvider{},
		ExtraDataCodec:      extraDataDecoder{},
	}
}
