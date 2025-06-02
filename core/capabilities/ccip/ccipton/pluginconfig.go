package ccipton

import (
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipnoop"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

func init() {
	// Register the Noop plugin config factory for Ton
	ccipcommon.RegisterPluginConfig(chainsel.FamilyTon, ccipnoop.InitializePluginConfig)
}
