package ccipton

import (
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipnoop"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

// InitializePluginConfig returns a pluginConfig for TON chains.
func InitializePluginConfig(_ logger.Logger, extraDataCodec ccipocr3.ExtraDataCodecBundle) ccipcommon.PluginConfig {
	return ccipcommon.PluginConfig{
		GasEstimateProvider:   ccipnoop.NewGasEstimateProvider(extraDataCodec), // TODO: implement
		TokenDataEncoder:      ccipnoop.NewTokenDataEncoder(),                  // TODO: implement
		CCIPProviderSupported: true,
	}
}

func init() {
	ccipcommon.RegisterPluginConfig(chainsel.FamilyTon, InitializePluginConfig)
}
