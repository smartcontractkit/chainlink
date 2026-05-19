package config

import (
	"github.com/aptos-labs/aptos-go-sdk"

	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
)

type MigrateOnRampDestChainConfigsToV2Config struct {
	ChainSelector         uint64
	DestChainSelectors    []uint64
	RouterModuleAddresses []aptos.AccountAddress
	MCMS                  *cldfproposalutils.TimelockConfig
}
