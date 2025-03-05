package keeper

import (
	"github.com/smartcontractkit/chainlink-integrations/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_1"
)

var Registry1_1ABI = types.MustGetABI(keeper_registry_wrapper1_1.KeeperRegistryABI)

type RegistryGasChecker interface {
	CheckGasOverhead() uint32
	PerformGasOverhead() uint32
	MaxPerformDataSize() uint32
}
