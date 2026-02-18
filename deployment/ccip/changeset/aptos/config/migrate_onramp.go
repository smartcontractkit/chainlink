package config

import (
	"github.com/aptos-labs/aptos-go-sdk"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

type MigrateOnRampDestChainConfigsToV2Config struct {
	ChainSelector         uint64
	DestChainSelectors    []uint64
	RouterModuleAddresses []aptos.AccountAddress
	MCMS                  *proposalutils.TimelockConfig
}
