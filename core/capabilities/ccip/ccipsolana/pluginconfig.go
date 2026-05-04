package ccipsolana

import (
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip/consts"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ocrimpls"
)

// InitializePluginConfig returns a minimal PluginConfig for Solana chains.
// Codec fields are sourced from CCIPProvider.Codec() at runtime.
func InitializePluginConfig(_ logger.Logger, extraDataCodec ccipocr3.ExtraDataCodecBundle) ccipcommon.PluginConfig {
	return ccipcommon.PluginConfig{
		ContractTransmitterFactory: ocrimpls.NewSVMContractTransmitterFactory(extraDataCodec),
		PriceOnlyCommitFn:          consts.MethodCommitPriceOnly,
		CCIPProviderSupported:      true,
	}
}

func init() {
	// Register the Solana plugin config factory
	ccipcommon.RegisterPluginConfig(chainsel.FamilySolana, InitializePluginConfig)
}
