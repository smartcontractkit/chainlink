package ccipton

import (
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ton/pkg/ccip/codec"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipnoop"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

// InitializePluginConfig returns a pluginConfig for TON chains.
func InitializePluginConfig(lggr logger.Logger, extraDataCodec ccipocr3.ExtraDataCodecBundle) ccipcommon.PluginConfig {
	return ccipcommon.PluginConfig{
		AddressCodec:       codec.NewAddressCodec(),
		CommitPluginCodec:  codec.NewCommitPluginCodecV1(),
		ExecutePluginCodec: codec.NewExecutePluginCodecV1(extraDataCodec),
		// TODO(EVM2TON): this is a temp fix for nil msgHasher access, should be using CCIPProvider msgHasher instead
		MessageHasher:         codec.NewMessageHasherV1(logger.Sugared(lggr).Named(chainsel.FamilyTon).Named("MessageHasherV1"), extraDataCodec),
		ExtraDataCodec:        codec.NewExtraDataDecoder(),
		GasEstimateProvider:   ccipnoop.NewGasEstimateProvider(extraDataCodec), // TODO: implement
		TokenDataEncoder:      ccipnoop.NewTokenDataEncoder(),                  // TODO: implement
		CCIPProviderSupported: true,
	}
}

func init() {
	ccipcommon.RegisterPluginConfig(chainsel.FamilyTon, InitializePluginConfig)
}
