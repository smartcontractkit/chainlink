package ccipton

import (
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipnoop"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

// InitializePluginConfig returns a pluginConfig for TON chains.
func InitializePluginConfig(lggr logger.Logger, extraDataCodec ccipcommon.ExtraDataCodec) ccipcommon.PluginConfig {
	return ccipcommon.PluginConfig{
		ChainAccessorFactory: TONChainAccessorFactory{},
	}
}

// TON plugin is a noop implementation for now.
// This registers TON with a noop plugin via init() when this package is imported.
// This follows the pattern of dynamic plugin registration used in oraclecreator/plugin.go.
func init() {
	// Register the Noop plugin config factory for Ton
	ccipcommon.RegisterPluginConfig(chainsel.FamilyTon, ccipnoop.NewPluginConfig)
}
