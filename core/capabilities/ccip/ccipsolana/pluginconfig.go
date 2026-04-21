package ccipsolana

import (
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip/consts"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/v2/core/config/env"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana/ccip/codec"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ocrimpls"
)

// InitializePluginConfig returns a pluginConfig for Solana chains.
func InitializePluginConfig(lggr logger.Logger, extraDataCodec ccipocr3.ExtraDataCodecBundle) ccipcommon.PluginConfig {
	return ccipcommon.PluginConfig{
		CommitPluginCodec:          codec.NewCommitPluginCodecV1(),
		ExecutePluginCodec:         codec.NewExecutePluginCodecV1(extraDataCodec),
		MessageHasher:              codec.NewMessageHasherV1(logger.Sugared(lggr).Named(chainsel.FamilySolana).Named("MessageHasherV1"), extraDataCodec),
		TokenDataEncoder:           codec.NewSolanaTokenDataEncoder(),
		GasEstimateProvider:        NewGasEstimateProvider(extraDataCodec),
		RMNCrypto:                  nil,
		ContractTransmitterFactory: ocrimpls.NewSVMContractTransmitterFactory(extraDataCodec),
		AddressCodec:               codec.NewAddressCodec(),
		ChainRW:                    ChainRWProvider{},
		ExtraDataCodec:             codec.NewExtraDataDecoder(),
		PriceOnlyCommitFn:          consts.MethodCommitPriceOnly,
		CCIPProviderSupported:      env.SolanaPlugin.Cmd.Get() != "",
	}
}

func init() {
	// Register the Solana plugin config factory
	ccipcommon.RegisterPluginConfig(chainsel.FamilySolana, InitializePluginConfig)
}
