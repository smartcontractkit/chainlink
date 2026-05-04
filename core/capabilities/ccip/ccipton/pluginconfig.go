package ccipton

import (
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

// InitializePluginConfig returns a minimal PluginConfig for TON chains.
// Codec fields are sourced from CCIPProvider.Codec() at runtime.
func InitializePluginConfig(_ logger.Logger, _ ccipocr3.ExtraDataCodecBundle) ccipcommon.PluginConfig {
	return ccipcommon.PluginConfig{
		CCIPProviderSupported: true,
	}
}

func init() {
	ccipcommon.RegisterPluginConfig(chainsel.FamilyTon, InitializePluginConfig)
}
